package loadbalancer

import (
	"context"
	"fmt"
	"sync"
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

	as := make([]string, 15)
	errs := make([]error, 15)
	wg := sync.WaitGroup{}
	for i := range 15 {
		wg.Add(1)
		go func(i int) {
			as[i], errs[i] = (&ipamAllocator{
				lxcClient:        lxcClient,
				clusterName:      fmt.Sprintf("c-%d", i),
				clusterNamespace: "default",

				networkName: "testbr0",

				rangesKey:   "user.capn.vip.ranges",
				volatileKey: func(s string) string { return fmt.Sprintf("user.capn.vip.volatile.%s", s) },
			}).Allocate(context.TODO())
			wg.Done()
		}(i)
	}

	wg.Wait()

	for i := range 15 {
		fmt.Println(as[i], errs[i])
	}
}
