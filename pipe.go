package pipe

import (
	_ "embed"

	"github.com/renatopp/pipelang/internal/runtime"
)

//go:embed VERSION
var version string

func RunCode(code []byte) (string, error) {
	rt := runtime.New()
	obj, err := rt.RunCode(code)
	if err != nil {
		return "", err
	}

	return obj.AsString(), nil
}

func RunFile(path string) (string, error) {
	rt := runtime.New()
	obj, err := rt.RunFile(path)
	if err != nil {
		return "", err
	}

	return obj.AsString(), nil
}

func Version() string {
	return version
}
