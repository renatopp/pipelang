package errfmt

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/renatopp/langtools/errors"
	"github.com/renatopp/langtools/lexers"
	"github.com/renatopp/langtools/parsers"
	"github.com/renatopp/langtools/utils"
	"github.com/renatopp/pipelang/internal"
)

type Error struct {
	Category string
	Errors   []errors.Error
	Source   []byte
	Path     string
	Message  string
	// Stack
}

func (e *Error) Error() string {
	return e.Message
}

func FormatLexerErrors(err []lexers.LexerError, source []byte, path string) error {
	var errs []errors.Error = make([]errors.Error, len(err))
	for i, e := range err {
		errs[i] = e
	}

	return formatErrorWithinSource("lexer", errs, source, path)
}

func FormatParserErrors(err []parsers.ParserError, source []byte, path string) error {
	var errs []errors.Error = make([]errors.Error, len(err))
	for i, e := range err {
		errs[i] = e
	}

	return formatErrorWithinSource("parser", errs, source, path)
}

func FormatEvaluationError(err *internal.Error, source []byte, path string) error {
	return formatErrorWithinSource("runtime", []errors.Error{err}, source, path)
}

func formatErrorWithinSource(category string, errs []errors.Error, source []byte, path string) error {
	result := ""
	for i, e := range errs {
		// Main error
		if i == 0 {
			fromLine, fromCol, toLine, toCol := e.Range()
			if fromLine == 0 {
				result += fmt.Sprintf("[%s error] %s\n", category, e.Error())
				continue
			}

			fromCol, toCol = min(fromCol, toCol), max(fromCol, toCol) // multi row errors
			lines := utils.GetSourceLines(source, fromLine, fromLine)
			highlight := utils.HighlightChars(fromCol, toCol)
			digits := len(strconv.Itoa(toLine)) + 1

			result += fmt.Sprintf("Error at file [%s], line %d, column %d:\n", path, fromLine, fromCol)
			for i, line := range lines {
				number := utils.PadLeft(strconv.Itoa(fromLine+i), digits)
				result += fmt.Sprintf("%s| %s\n", number, strings.ReplaceAll(line, "\t", " "))
			}
			result += fmt.Sprintf("%s| %s\n", utils.PadLeft("", digits), highlight)
			result += fmt.Sprintf("[%s error] %s\n", category, e.Error())

			// Secondary Errors
		} else {
			if i == 1 {
				result += "\nOther possible errors:\n"
			}

			fromLine, fromCol := e.At()
			result += fmt.Sprintf("  %d, %d: %s\n", fromLine, fromCol, e.Error())
		}
	}

	return &Error{
		Category: category,
		Errors:   errs,
		Source:   source,
		Path:     path,
		Message:  result,
	}
}
