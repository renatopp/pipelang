package main

import (
	"bufio"
	"flag"
	"fmt"
	"log"
	"os"

	pipe "github.com/renatopp/pipelang"
	"github.com/renatopp/pipelang/cmd/pipe/cmds"
)

func init() {
	var verbose bool
	flag.BoolVar(&verbose, "v", false, "enable compiler debug messages")
	flag.BoolVar(&verbose, "verbose", false, "enable compiler debug messages")
	flag.Parse()

	if verbose {
		os.Setenv("PIPE_VERBOSE", "1")
	}
}

func main() {
	if len(os.Args) < 2 {
		shell()
		return
	}

	switch os.Args[1] {

	case "version":
		version()

	case "help":
		help()

	case "-":
		stdin()

	case "run":
		run()

	case "eval":
		eval()

	case "shell":
		shell()

	case "debug":
		debug()

	default:
		help()
	}
}

func version() {
	fmt.Println(pipe.Version())
}

func help() {
	fmt.Println("usage: pipe [command] [args]")
	fmt.Println("")
	fmt.Println("commands:")
	fmt.Println("  version        Print the version")
	fmt.Println("  help           Print this help")
	fmt.Println("  shell          Start the REPL")
	fmt.Println("  run [file]     Run a file")
	fmt.Println("  eval [command] Evaluate a string")
	fmt.Println("  -              Read from stdin")
}

func shell() {
	cmds.Shell()
}

func run() {
	cmds.Run()
}

func eval() {
	cmds.Eval()
}

func debug() {
	cmds.Debug()
}

func stdin() {
	stat, _ := os.Stdin.Stat()
	if (stat.Mode() & os.ModeCharDevice) != 0 {
		return
	}

	var stdin []byte
	scanner := bufio.NewScanner(os.Stdin)
	for scanner.Scan() {
		stdin = append(stdin, scanner.Bytes()...)
	}
	if err := scanner.Err(); err != nil {
		log.Fatal(err)
	}
	fmt.Printf("stdin = %s\n", stdin)
}
