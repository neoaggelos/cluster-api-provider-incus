package loadbalancer

import (
	"context"
	"fmt"
	"testing"

	"github.com/lxc/cluster-api-provider-incus/internal/lxc"

	. "github.com/onsi/gomega"
)

func TestAllocate(t *testing.T) {
	g := NewWithT(t)

	opts, err := lxc.ConfigurationFromLocal("", "", false)
	g.Expect(err).ToNot(HaveOccurred())

	lxcClient, err := lxc.New(context.TODO(), opts)
	g.Expect(err).ToNot(HaveOccurred())

	a := &ipamAllocator{
		lxcClient:        lxcClient,
		clusterName:      "c6",
		clusterNamespace: "default",

		networkName:       "testbr0",
		rangesKey:         "user.capn.vip.ranges",
		volatilePrefixKey: "user.capn.vip.volatile",
	}

	addr, err := a.Allocate(context.TODO())
	g.Expect(err).ToNot(HaveOccurred())

	fmt.Println(addr)
}
