package expression_test

import (
	"testing"

	"github.com/renatopp/pipelang/test/common"
)

func TestSingleAssignments(t *testing.T) {
	common.AssertCode(t, `a := 1`, `1`)
	common.AssertCode(t, `a := 1; a`, `1`)
	common.AssertCode(t, `a := (5, 3, 2)`, `5`)
	common.AssertCode(t, `a := [5, 3, 2]`, `[5, 3, 2]`)
	common.AssertCode(t, `a := 1; a := 're'`, `re`)

	common.AssertCodeError(t, `a := 1; a = '3'`)
	common.AssertCodeError(t, `a = 3`)
}

func TestTupledAssignments(t *testing.T) {
	common.AssertCode(t, `a := (1, 2)`, `1`)
	common.AssertCode(t, `(a, b) := (1, 2)`, `(1, 2)`)
	common.AssertCode(t, `(a, b) := (1, 2, 3); (a, b)`, `(1, 2)`)
	common.AssertCode(t, `(a, b, c) := (3, 2, 1); (a, b, c)`, `(3, 2, 1)`)
	common.AssertCode(t, `(a, (b, d), c) := (3, (2, 4, 6), 1, 3); (a, b, c, d)`, `(3, 2, 1, 4)`)

	common.AssertCodeError(t, `(a, b) := 1`)
	common.AssertCodeError(t, `(a, b) := [1, 2, 3]`)
	common.AssertCodeError(t, `(a, b, c) := (2, 3)`)
	common.AssertCodeError(t, `(a, b) = (2, 3)`)
	common.AssertCodeError(t, `(a, b) := (2, 3); (a, b) = (3, 'invalid')`)
}

func TestSpreadAssignments(t *testing.T) {
	common.AssertCode(t, `a := [1,2,3,4]...`, `1`)
	common.AssertCode(t, `(a, b, c, d) := [1,2,3,4]...; (a, b, c, d)`, `(1, 2, 3, 4)`)
	common.AssertCode(t, `(a, ...c, d) := [1,2,3,4]...; (a, c, d)`, `(1, [2, 3], 4)`)
	common.AssertCode(t, `(a, ...b) := [1,2,3,4]...; (a, b)`, `(1, [2, 3, 4])`)
	common.AssertCode(t, `(a, ...b) := (1)...; (a, b)`, `(1, [])`)
	common.AssertCode(t, `(a, ...b, c) := (1, 2)...; (a, b, c)`, `(1, [], 2)`)

	common.AssertCodeError(t, `(a, ...b) := []...`)
	common.AssertCodeError(t, `(a, ...b, c) := (1, )`)
}

func TestOpenTuples(t *testing.T) {
	common.AssertCode(t, `a, b := 1, 2; (a, b)`, `(1, 2)`)
}

func TestAs(t *testing.T) {
	common.AssertCode(t, `1, 2 as (a, b);`, `(1, 2)`)
}
