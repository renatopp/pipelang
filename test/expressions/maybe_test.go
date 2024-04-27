package expression_test

import (
	"testing"

	"github.com/renatopp/pipelang/test/common"
)

func TestMaybe(t *testing.T) {
	common.AssertCode(t, `
		maybe := Maybe(3)
		[maybe.Ok(), maybe.Error(), maybe.Value(), maybe.Result()]
	`, `[true, false, 3, 3]`)

	common.AssertCode(t, `
	maybe := Maybe(Error('error'))
	[maybe.Ok(), maybe.Error(), maybe.Value(), maybe.Result()]
	`, `[false, error, false, error]`)

	common.AssertCode(t, `
		maybe := Maybe(Maybe(Maybe(3)))
		[maybe.Ok(), maybe.Error(), maybe.Value(), maybe.Result()]
	`, `[true, false, 3, 3]`)

	common.AssertCode(t, `
		maybe := Maybe(Maybe(Maybe(Error('error'))))
		[maybe.Ok(), maybe.Error(), maybe.Value(), maybe.Result()]
	`, `[false, error, false, error]`)

	common.AssertCode(t, `
		maybe := Maybe(Maybe(Maybe(Error('error'))))
		(err, val) := maybe!
		[err, val]
	`, `[error, false]`)

	common.AssertCode(t, `
		maybe := Maybe(34)
		(err, val) := maybe!
		[err, val]
	`, `[false, 34]`)
}
