package utils

import (
	"fmt"
	"testing"

	. "github.com/onsi/gomega"
)

func TestIPRange(t *testing.T) {
	g := NewWithT(t)

	r, err := ParseIPRange("10.100.10.100-10.100.10.120")
	g.Expect(err).ToNot(HaveOccurred())

	for addr := range r.Iterate() {
		fmt.Println(addr)

		if addr == "10.100.10.114" {
			// break
		}
	}
}
