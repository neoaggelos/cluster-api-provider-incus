package types_test

import (
	"fmt"
	"testing"

	"github.com/lxc/cluster-api-provider-incus/internal/types"
	. "github.com/onsi/gomega"
)

func TestTerminalError(t *testing.T) {
	for _, tc := range []struct {
		name           string
		err            error
		expectTerminal bool
	}{
		{name: "NilErrorIsNot"},
		{name: "SomeErrorIsNot", err: fmt.Errorf("some error")},
		{name: "TerminalErrorIs", err: types.TerminalError(fmt.Errorf("terminal error")), expectTerminal: true},
		{name: "WrappedTerminalErrorIs", err: fmt.Errorf("wrapped: %w", types.TerminalError(fmt.Errorf("some error"))), expectTerminal: true},
	} {
		t.Run(tc.name, func(t *testing.T) {
			g := NewWithT(t)
			if tc.expectTerminal {
				g.Expect(types.IsTerminalError(tc.err)).To(BeTrue(), "must be a terminal error")
			} else {
				g.Expect(types.IsTerminalError(tc.err)).To(BeFalse(), "must not be a terminal error")
			}
		})
	}
}
