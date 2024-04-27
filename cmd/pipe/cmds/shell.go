package cmds

import (
	"fmt"
	"os"
	"os/exec"
	gor "runtime"

	"github.com/peterh/liner"
	pipe "github.com/renatopp/pipelang"
	"github.com/renatopp/pipelang/internal/ast"
	"github.com/renatopp/pipelang/internal/object"
	"github.com/renatopp/pipelang/internal/runtime"
)

func Shell() {
	clearConsole()
	rt := runtime.New()
	line := liner.NewLiner()
	defer line.Close()

	line.SetCtrlCAborts(true)
	line.SetMultiLineMode(true)

	fmt.Printf("| PIPELANG v%s\n", pipe.Version())
	fmt.Printf("| Type 'clear' to clear the console\n")
	fmt.Printf("| Type 'exit' to quit\n")
	fmt.Printf("\n")

main:
	for {
		cmd, err := line.Prompt("pipe> ")
		if err != nil {
			break
		}

		switch cmd {
		case "clear":
			clearConsole()
			continue
		case "exit":
			return
		}

		line.AppendHistory(cmd)
		node, err := rt.LoadAst([]byte(cmd))
		if err != nil {
			println(err.Error())
			continue
		}

		block := node.(*ast.Block)
		var res object.Object
		for _, exp := range block.Expressions {
			res, err = rt.RunAst(exp)
			if err != nil {
				println(err.Error())
				continue main
			}
		}
		println(res.AsString())
	}
}

func clearConsole() {
	var cmd *exec.Cmd
	if gor.GOOS == "windows" {
		cmd = exec.Command("cmd", "/c", "cls")
	} else {
		cmd = exec.Command("clear")
	}
	cmd.Stdout = os.Stdout
	cmd.Run()
}
