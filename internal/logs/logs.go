package logs

import (
	"log"
	"os"
)

var logger *log.Logger

func init() {
	logger = log.New(os.Stderr, "", 0)
}

func isEnabled() bool {
	return os.Getenv("PIPE_VERBOSE") == "1"
}

func Print(msg string, v ...any) {
	if !isEnabled() {
		return
	}

	logger.Printf(msg+"\n", v...)
}
