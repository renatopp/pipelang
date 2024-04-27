package common

import (
	"testing"

	"github.com/renatopp/pipelang/internal/object"
	"github.com/renatopp/pipelang/internal/runtime"
)

func Run(program string) (*object.Scope, object.Object, error) {
	r := runtime.New()
	obj, err := r.RunCode([]byte(program))
	return r.GlobalScope(), obj, err
}

func AssertCode(t *testing.T, input string, output string) {
	_, obj, err := Run(input)
	if err != nil {
		t.Errorf("unexpected error: %v", err)
		return
	}

	if obj.AsString() != output {
		t.Errorf("expected %s, got %s", output, obj.AsString())
	}
}

func AssertCodeError(t *testing.T, input string) {
	_, _, err := Run(input)
	if err == nil {
		t.Errorf("expected error, got nil")
	}
}
