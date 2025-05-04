package main

var (
	containerFTypeByServer = map[string]string{
		"incus": "incus_combined.tar.gz",
		"lxd":   "lxd_combined.tar.gz",
	}
	vmMetadataFTypeByServer = map[string]string{
		"incus": "incus.tar.xz",
		"lxd":   "lxd.tar.xz",
	}
	vmRootfsFTypeByServer = map[string]string{
		"incus": "disk-kvm.img",
		"lxd":   "disk1.img",
	}
)
