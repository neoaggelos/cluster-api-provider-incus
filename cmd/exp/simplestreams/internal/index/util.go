package index

import (
	"bytes"
	"context"
	"crypto/sha256"
	"fmt"
	"io"
	"os"
	"os/exec"
	"path/filepath"

	"sigs.k8s.io/controller-runtime/pkg/log"
)

type containerImageInfo struct {
	Sha256 string
	Size   int64
}

func getContainerImageInfo(path string) (containerImageInfo, error) {
	stat, err := os.Stat(path)
	if err != nil {
		return containerImageInfo{}, fmt.Errorf("failed to stat: %w", err)
	}

	// get the image sha256
	f, err := os.Open(path)
	if err != nil {
		return containerImageInfo{}, fmt.Errorf("failed to open: %w", err)
	}
	defer func() {
		_ = f.Close()
	}()

	hash := sha256.New()
	if _, err = io.Copy(hash, f); err != nil {
		return containerImageInfo{}, fmt.Errorf("failed to calculate sha256 sum: %w", err)
	}

	return containerImageInfo{
		Size:   stat.Size(),
		Sha256: fmt.Sprintf("%x", hash.Sum(nil)),
	}, nil
}

type virtualMachineImageInfo struct {
	MetaSize   int64
	MetaSha256 string

	RootSize   int64
	RootSha256 string

	CombinedSha256 string
}

func getVirtualMachineImageInfo(metadata []byte, rootfs []byte) (virtualMachineImageInfo, error) {
	info := virtualMachineImageInfo{
		MetaSize: int64(len(metadata)),
		RootSize: int64(len(rootfs)),
	}
	hash := sha256.New()
	if _, err := hash.Write(metadata); err != nil {
		return virtualMachineImageInfo{}, fmt.Errorf("failed to calculate metadata sha256 sum: %w", err)
	}
	info.MetaSha256 = fmt.Sprintf("%x", hash.Sum(nil))

	if _, err := hash.Write(rootfs); err != nil {
		return virtualMachineImageInfo{}, fmt.Errorf("failed to calculate combined sha256 sum: %w", err)
	}
	info.CombinedSha256 = fmt.Sprintf("%x", hash.Sum(nil))

	hash.Reset()
	if _, err := hash.Write(rootfs); err != nil {
		return virtualMachineImageInfo{}, fmt.Errorf("failed to calculate rootfs sha256 sum: %w", err)
	}
	info.RootSha256 = fmt.Sprintf("%x", hash.Sum(nil))

	return info, nil
}

func copyFile(source, destination string) error {
	f, err := os.Open(source)
	if err != nil {
		return fmt.Errorf("failed to open source: %w", err)
	}
	defer func() {
		_ = f.Close()
	}()

	if err := os.MkdirAll(filepath.Dir(destination), 0755); err != nil {
		return fmt.Errorf("failed to create directory for destination: %w", err)
	}
	of, err := os.Create(destination)
	if err != nil {
		return fmt.Errorf("failed to open destination: %w", err)
	}

	if n, err := io.Copy(of, f); err != nil {
		return fmt.Errorf("failed to copy file data: %w", err)
	} else if err := of.Truncate(n); err != nil {
		return fmt.Errorf("failed to truncate output file: %w", err)
	}

	if err := of.Close(); err != nil {
		return fmt.Errorf("failed to flush file: %w", err)
	}

	return nil
}

// qemuCompressImage accepts raw bytes of a qcow2 image, compresses with `qemu-img convert -c` and returns the raw bytes of the compressed image
func qemuCompressImage(ctx context.Context, raw []byte) ([]byte, error) {
	log.FromContext(ctx).Info("Compressing rootfs.img, this might take a while", "uncompressed", len(raw))
	tmpDir, err := os.MkdirTemp("", "")
	if err != nil {
		return nil, fmt.Errorf("failed to create temporary directory: %w", err)
	}
	defer func() {
		_ = os.RemoveAll(tmpDir)
	}()
	if err := os.WriteFile(filepath.Join(tmpDir, "uncompressed.qcow2"), raw, 0644); err != nil {
		return nil, fmt.Errorf("failed to write uncompressed rootfs to temporary file: %w", err)
	}

	// attempt to use qemu-img from path, or fallbcak to /opt/incus/bin/qemu-img
	qemuImg, err := exec.LookPath("qemu-img")
	if err != nil {
		qemuImg = "/opt/incus/bin/qemu-img"
	}

	var stderr bytes.Buffer
	cmd := exec.CommandContext(ctx, qemuImg, "convert", "-O", "qcow2", "-c", filepath.Join(tmpDir, "uncompressed.qcow2"), filepath.Join(tmpDir, "compressed.qcow2"))
	cmd.Stderr = &stderr

	if err := cmd.Run(); err != nil {
		return nil, fmt.Errorf("qemu-img convert -c command failed with stderr=%q: %w", stderr.String(), err)
	}

	b, err := os.ReadFile(filepath.Join(tmpDir, "compressed.qcow2"))
	if err != nil {
		return nil, fmt.Errorf("failed to read rootfs after compression: %w", err)
	}

	log.FromContext(ctx).Info("Compressed rootfs.img", "uncompressed", len(raw), "compressed", len(b))
	return b, nil
}
