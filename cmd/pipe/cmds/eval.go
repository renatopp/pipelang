package cmds

import (
	"fmt"
	"os"
	"strings"

	pipe "github.com/renatopp/pipelang"
)

func Eval() {
	program := strings.Join(os.Args[2:], ";")
	res, err := pipe.RunCode([]byte(program))
	if err != nil {
		fmt.Println(err)
	}

	fmt.Println(res)
}
