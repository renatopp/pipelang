package expression_test

import (
	"testing"

	"github.com/renatopp/pipelang/test/common"
)

func TestData(t *testing.T) {
	common.AssertCode(t, `
	data Sample {
		attr = 1

		fn method(this) {
			return this.attr * 2
		}
	}

	s := Sample()
	s.method()
	`, `2`)
}
