package expression_test

import (
	"testing"

	"github.com/renatopp/pipelang/test/common"
)

func TestArithmeticExpressions(t *testing.T) {
	common.AssertCode(t, `1 + 1`, `2`)
	common.AssertCode(t, `1 + 2 * 3`, `7`)
	common.AssertCode(t, `1 * 2 + 3`, `5`)
	common.AssertCode(t, `1 * (2 + 3)`, `5`)
	common.AssertCode(t, `(1 + 2) * 3`, `9`)
	common.AssertCode(t, `2/5 * 4`, `1.600000`)
	common.AssertCode(t, `10%2`, `0`)
	common.AssertCode(t, `11%2`, `1`)
	common.AssertCode(t, `2^10`, `1024`)
	common.AssertCode(t, `+1`, `1`)
	common.AssertCode(t, `-1`, `-1`)

	common.AssertCodeError(t, `1 +`)
	common.AssertCodeError(t, `/1`)
	common.AssertCodeError(t, `1 + true`)
}

func TestComparisonExpressions(t *testing.T) {
	common.AssertCode(t, `1 == 1`, `true`)
	common.AssertCode(t, `1 != 1`, `false`)
	common.AssertCode(t, `1 < 1`, `false`)
	common.AssertCode(t, `1 <= 1`, `true`)
	common.AssertCode(t, `1 > 1`, `false`)
	common.AssertCode(t, `1 >= 1`, `true`)

	common.AssertCode(t, `1 == 2`, `false`)
	common.AssertCode(t, `1 != 2`, `true`)
	common.AssertCode(t, `1 < 2`, `true`)
	common.AssertCode(t, `1 <= 2`, `true`)
	common.AssertCode(t, `1 > 2`, `false`)
	common.AssertCode(t, `1 >= 2`, `false`)
	common.AssertCode(t, `1 == false`, `false`)

	common.AssertCodeError(t, `1 ==`)
}

func TestLogicalExpressions(t *testing.T) {
	common.AssertCode(t, `true and true`, `true`)
	common.AssertCode(t, `true and false`, `false`)
	common.AssertCode(t, `false and true`, `false`)
	common.AssertCode(t, `false and false`, `false`)

	common.AssertCode(t, `true or true`, `true`)
	common.AssertCode(t, `true or false`, `true`)
	common.AssertCode(t, `false or true`, `true`)
	common.AssertCode(t, `false or false`, `false`)

	common.AssertCode(t, `true xor true`, `false`)
	common.AssertCode(t, `true xor false`, `true`)
	common.AssertCode(t, `false xor true`, `true`)
	common.AssertCode(t, `false xor false`, `false`)

	common.AssertCode(t, `not true`, `false`)
	common.AssertCode(t, `not false`, `true`)

	common.AssertCodeError(t, `true and`)
	common.AssertCodeError(t, `true or`)
}

func TestConcat(t *testing.T) {
	common.AssertCode(t, `"a" .. "b"`, `ab`)
	common.AssertCode(t, `"a" .. 1`, `a1`)
	common.AssertCode(t, `1 .. "a"`, `1a`)
	common.AssertCode(t, `1 .. 2`, `12`)
}
