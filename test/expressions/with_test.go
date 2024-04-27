package expression_test

import (
	"testing"

	"github.com/renatopp/pipelang/test/common"
)

func TestWith(t *testing.T) {
	common.AssertCode(t, `
	x := 1
	with 5 as x {
		x + 1
	}
	x
	`, `1`)

	common.AssertCode(t, `
	with 5 as x {
		x + 1
	}
	`, `6`)

	common.AssertCodeError(t, `
	with 5 as x {
		x + 1
	}
	x
	`)
}
