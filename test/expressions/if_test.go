package expression_test

import (
	"testing"

	"github.com/renatopp/pipelang/test/common"
)

func TestIfs(t *testing.T) {
	common.AssertCode(t, `if true { 'ok' } else { 'fail' }`, `ok`)
	common.AssertCode(t, `if false { 'ok' } else { 'fail' }`, `fail`)
	common.AssertCode(t, `if x := 2; false { x*5 } else { x*2 }`, `4`)
	common.AssertCode(t, `if x := 2; true { x*5 } else { x*2 }`, `10`)

	common.AssertCodeError(t, `if {}`)
	common.AssertCodeError(t, `if`)
	common.AssertCodeError(t, `if
	{}`)
}
