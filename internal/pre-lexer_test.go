package internal_test

import (
	"testing"

	"github.com/renatopp/langtools/tokens"
	"github.com/renatopp/pipelang/internal"
	"github.com/stretchr/testify/assert"
)

func TestBasic(t *testing.T) {
	var input = `
	fn main() {} [] 1_000_000, | : :=
	`

	expected := []tokens.TokenType{
		internal.T_EOE,
		internal.T_KEYWORD,
		internal.T_IDENTIFIER,
		internal.T_LPAREN,
		internal.T_RPAREN,
		internal.T_LBRACE,
		internal.T_RBRACE,
		internal.T_LBRACK,
		internal.T_RBRACK,
		internal.T_NUMBER,
		internal.T_COMMA,
		internal.T_PIPE,
		internal.T_LAMBDA,
		internal.T_ASSIGNMENT,
		internal.T_EOE,
		internal.T_EOE,
		internal.T_EOF,
	}
	lexer := internal.NewPreLexer([]byte(input))
	result := lexer.All()

	for i, token := range result {
		assert.Equal(t, expected[i], token.Type)
	}
}

func TestNumbers(t *testing.T) {
	var input = `
	1 1.0 1.0e10 1.0e-10 1.0e+10 0x1 0b1 0o1
	`

	expected := []tokens.TokenType{
		internal.T_EOE,
		internal.T_NUMBER,
		internal.T_NUMBER,
		internal.T_NUMBER,
		internal.T_NUMBER,
		internal.T_NUMBER,
		internal.T_HEX_NUMBER,
		internal.T_BIN_NUMBER,
		internal.T_OCT_NUMBER,
		internal.T_EOE,
		internal.T_EOE,
		internal.T_EOF,
	}
	lexer := internal.NewPreLexer([]byte(input))
	result := lexer.All()

	for i, token := range result {
		assert.Equal(t, expected[i], token.Type)
	}
}
