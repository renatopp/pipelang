package expression_test

import (
	"testing"

	"github.com/renatopp/pipelang/test/common"
)

func TestErrors(t *testing.T) {
	common.AssertCodeError(t, `
		fn explode {
			raise 'error'
		}
		explode()
	`)

	common.AssertCode(t, `
	fn explode {
		raise 'error'
	}
	explode()?!
	`, `(error, false)`)

	common.AssertCode(t, `
	fn explode {
		raise 'error'
	}
	explode() ?? 'default'
	`, `default`)
}
