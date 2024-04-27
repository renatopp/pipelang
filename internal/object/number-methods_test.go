package object_test

import (
	"testing"

	"github.com/renatopp/pipelang/test/common"
)

func TestString_Abs(t *testing.T) {
	common.AssertCode(t, `(-5).Abs()`, `5`)
	common.AssertCode(t, `(0).Abs()`, `0`)
	common.AssertCode(t, `(5).Abs()`, `5`)
	common.AssertCode(t, `(-5.32).Abs()`, `5.320000`)
}

func TestString_Ceil(t *testing.T) {
	common.AssertCode(t, `(-5).Ceil()`, `-5`)
	common.AssertCode(t, `(0).Ceil()`, `0`)
	common.AssertCode(t, `(5).Ceil()`, `5`)
	common.AssertCode(t, `(-5.32).Ceil()`, `-5`)
	common.AssertCode(t, `(-5.62).Ceil()`, `-5`)
}

func TestString_Floor(t *testing.T) {
	common.AssertCode(t, `(-5).Floor()`, `-5`)
	common.AssertCode(t, `(0).Floor()`, `0`)
	common.AssertCode(t, `(5).Floor()`, `5`)
	common.AssertCode(t, `(-5.32).Floor()`, `-6`)
	common.AssertCode(t, `(-5.62).Floor()`, `-6`)
}

func TestString_Round(t *testing.T) {
	common.AssertCode(t, `(-5).Round()`, `-5`)
	common.AssertCode(t, `(0).Round()`, `0`)
	common.AssertCode(t, `(5).Round()`, `5`)
	common.AssertCode(t, `(-5.32).Round()`, `-5`)
	common.AssertCode(t, `(-5.62).Round()`, `-6`)
}

func TestString_RoundToEven(t *testing.T) {
	common.AssertCode(t, `(-5).RoundToEven()`, `-5`)
	common.AssertCode(t, `(0).RoundToEven()`, `0`)
	common.AssertCode(t, `(5).RoundToEven()`, `5`)
	common.AssertCode(t, `(-5.32).RoundToEven()`, `-5`)
	common.AssertCode(t, `(-5.62).RoundToEven()`, `-6`)
}

func TestString_Sign(t *testing.T) {
	common.AssertCode(t, `(-5).Sign()`, `-1`)
	common.AssertCode(t, `(0).Sign()`, `1`)
	common.AssertCode(t, `(5).Sign()`, `1`)
	common.AssertCode(t, `(-5.32).Sign()`, `-1`)
	common.AssertCode(t, `(-5.62).Sign()`, `-1`)
}

func TestString_CopySign(t *testing.T) {
	common.AssertCode(t, `(-5).CopySign(10)`, `5`)
	common.AssertCode(t, `(0).CopySign(10)`, `0`)
	common.AssertCode(t, `(5).CopySign(-10)`, `-5`)
	common.AssertCode(t, `(-5.32).CopySign(10)`, `5.320000`)
	common.AssertCode(t, `(-5.62).CopySign(10)`, `5.620000`)
}

func TestString_Truncate(t *testing.T) {
	common.AssertCode(t, `(-5).Truncate()`, `-5`)
	common.AssertCode(t, `(0).Truncate()`, `0`)
	common.AssertCode(t, `(5).Truncate()`, `5`)
	common.AssertCode(t, `(-5.32).Truncate()`, `-5`)
	common.AssertCode(t, `(-5.62).Truncate()`, `-5`)
}

func TestString_Clamp(t *testing.T) {
	common.AssertCode(t, `(-5).Clamp(-10, 10)`, `-5`)
	common.AssertCode(t, `(0).Clamp(-10, 10)`, `0`)
	common.AssertCode(t, `(5).Clamp(-10, 10)`, `5`)
	common.AssertCode(t, `(-5.32).Clamp(-5, 5)`, `-5`)
	common.AssertCode(t, `(-5.62).Clamp(-5, 5)`, `-5`)
}

func TestString_Remainder(t *testing.T) {
	common.AssertCode(t, `(-5).Remainder(2)`, `-1`)
	common.AssertCode(t, `(0).Remainder(2)`, `0`)
	common.AssertCode(t, `(5).Remainder(2)`, `1`)
	common.AssertCode(t, `(-5.32).Remainder(2)`, `0.680000`)
	common.AssertCode(t, `(-5.62).Remainder(2)`, `0.380000`)
}

func TestString_Min(t *testing.T) {
	common.AssertCode(t, `(0).Min(1, 2, 0, -1)`, `-1`)
	common.AssertCode(t, `(-5).Min(2)`, `-5`)
	common.AssertCode(t, `(0).Min(2)`, `0`)
	common.AssertCode(t, `(5).Min(2)`, `2`)
	common.AssertCode(t, `(-5.32).Min(2)`, `-5.320000`)
	common.AssertCode(t, `(-5.62).Min(2)`, `-5.620000`)
}

func TestString_Max(t *testing.T) {
	common.AssertCode(t, `(-5).Max(2)`, `2`)
	common.AssertCode(t, `(0).Max(2)`, `2`)
	common.AssertCode(t, `(5).Max(2)`, `5`)
	common.AssertCode(t, `(-5.32).Max(2)`, `2`)
	common.AssertCode(t, `(-5.62).Max(2)`, `2`)
}

func TestIsOdd(t *testing.T) {
	common.AssertCode(t, `(-5).IsOdd()`, `true`)
	common.AssertCode(t, `(0).IsOdd()`, `false`)
	common.AssertCode(t, `(5).IsOdd()`, `true`)
	common.AssertCode(t, `(6).IsOdd()`, `false`)
}

func TestIsEven(t *testing.T) {
	common.AssertCode(t, `(-5).IsEven()`, `false`)
	common.AssertCode(t, `(0).IsEven()`, `true`)
	common.AssertCode(t, `(5).IsEven()`, `false`)
	common.AssertCode(t, `(6).IsEven()`, `true`)
}
