package expression_test

import (
	"testing"

	"github.com/renatopp/pipelang/test/common"
)

func TestPipes(t *testing.T) {
	fun := `
	fn sum(x, y) {
		x + y
	}

	fn mult(x, y) {
		x * y
	}

	`

	common.AssertCode(t, fun+`25 | sum 5 | mult 2`, `60`)
	common.AssertCode(t, fun+`10 | mult 2 | sum 5`, `25`)
}
