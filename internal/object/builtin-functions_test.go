package object_test

import (
	"testing"

	"github.com/renatopp/pipelang/test/common"
)

func TestFunction_Filter(t *testing.T) {
	common.AssertCode(t, ` [1,2,3] | filter x: x.IsOdd() | List`, `[1, 3]`)
}

func TestFunction_Each(t *testing.T) {
	common.AssertCode(t, ` [1,2,3] | each x: x*2 | List`, `[1, 2, 3]`)
}

func TestFunction_Map(t *testing.T) {
	common.AssertCode(t, ` [1,2,3] | map x: x*2 | List`, `[2, 4, 6]`)
}

func TestFunction_Reduce(t *testing.T) {
	common.AssertCode(t, ` [1,2,3] | reduce 0, (acc, x): acc + x`, `6`)
	// common.AssertCode(t, ` [1,2,3] | reduce 0, Number.Add`, `6`)
}

func TestFunction_Sum(t *testing.T) {
	common.AssertCode(t, ` [1,2,3] | sum`, `6`)
}

func TestFunction_SumBy(t *testing.T) {
	common.AssertCode(t, ` [1,2,3] | sumBy x: x*2`, `12`)
}

func TestFunction_Count(t *testing.T) {
	common.AssertCode(t, ` [1,2,3] | count`, `3`)
}

func TestFunction_CountBy(t *testing.T) {
	common.AssertCode(t, ` [1,2,3] | countBy x: 2`, `6`)
}
