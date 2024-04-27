package expression_test

import (
	"testing"

	"github.com/renatopp/pipelang/test/common"
)

func TestBasicFunctions(t *testing.T) {
	common.AssertCode(t, `fn named() { 2 }; named()`, `2`)
	common.AssertCode(t, `fn named() { 4; return 2 }; named()`, `2`)
	common.AssertCode(t, `unnamed := fn() { 5 }; unnamed()`, `5`)
	common.AssertCode(t, `fn multBy(x) { fn(y) { x * y} }; multBy(5)(2)`, `10`)
	common.AssertCode(t, `fn noparam { 'ok' }; noparam()`, `ok`)
	common.AssertCode(t, `noparam := fn { 'ok' }; noparam()`, `ok`)
	common.AssertCode(t, `fn ht(a, ...b, c){ [a, b, c] }; ht(1, 2, 3, 4)`, `[1, [2, 3], 4]`)
	common.AssertCode(t, `fn p(a, b){ [a, b] }; p(1, 2, 3)`, `[1, 2]`)

	common.AssertCodeError(t, `fn invalid`)
	common.AssertCodeError(t, `fn e { raise 'error' }; e()`)
	common.AssertCodeError(t, `fn x(a, b) {}; x(2)`)
}

func TestLambdas(t *testing.T) {
	common.AssertCode(t, `f := x : x + 1;f(2)`, `3`)
	common.AssertCode(t, `x:=4;f :=:x*2;f()`, `8`)

	common.AssertCodeError(t, `:`)
}

func TestGenerators(t *testing.T) {
	def := `
	fn OneTwoThree() {
		yield 1
		yield 2
		yield 3
	}
	s := OneTwoThree()
	[s.Finished(), s.Next().Result(), s.Next().Result(), s.Next().Result(), s.Next().Ok(), s.Finished()]	`
	common.AssertCode(t, def, `[false, 1, 2, 3, false, true]`)

	def = `
	fn OneTwoThree() {
		yield 1
		yield break
		yield 2
		yield 3
	}
	s := OneTwoThree()
	[s.Finished(), s.Next().Result(), s.Next().Ok(), s.Finished()]
	`
	common.AssertCode(t, def, `[false, 1, false, true]`)

	def = `
	fn OneTwoThree() {
		yield 1
		return 0
	}
	`
	common.AssertCodeError(t, def)
}
