package cmds

import (
	"fmt"
	"os"

	pipe "github.com/renatopp/pipelang"
)

func Run() {
	file := os.Args[2]
	_, err := pipe.RunFile(file)
	if err != nil {
		fmt.Println(err)
	}

	// fmt.Println(res)
}
