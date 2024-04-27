package expression_test

import (
	"testing"

	"github.com/renatopp/pipelang/test/common"
)

func TestMatch(t *testing.T) {
	fizzbuzz := `
	fn fizzbuzz(x) {
		match (x%3, x%5) {
			(0, 0): "FizzBuzz"
			(0, _): "Fizz"
			(_, 0): "Buzz"
			(_ as  f,  _ as s): [x, f, s]
		}
	}
	`
	common.AssertCode(t, fizzbuzz+`fizzbuzz(1)`, `[1, 1, 1]`)
	common.AssertCode(t, fizzbuzz+`fizzbuzz(3)`, `Fizz`)
	common.AssertCode(t, fizzbuzz+`fizzbuzz(5)`, `Buzz`)
	common.AssertCode(t, fizzbuzz+`fizzbuzz(15)`, `FizzBuzz`)

}
