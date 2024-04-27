package expression_test

import (
	"testing"

	"github.com/renatopp/pipelang/test/common"
)

func TestFors(t *testing.T) {
	common.AssertCode(t, `
	i := 0
	for {
		i += 1
		if i == 5 { break }
	}
	i
	`, `5`)

	common.AssertCode(t, `
	i := 0
	t := 0
	for i < 10 {
		i += 1
		if i % 2 == 0 {
			continue
		}
		t += i
	}
	t
	`, `25`)

	common.AssertCode(t, `
	fn OneTwoThree() {
		yield 1
		yield 2
		yield 3
	}

	i := 0
	for x in OneTwoThree() {
		i += x
	}
	i
	`, `6`)

	common.AssertCodeError(t, `y in x`)
	common.AssertCodeError(t, `for x in 1`)
	common.AssertCodeError(t, `for x`)
	common.AssertCodeError(t, `for for {}`)
}
