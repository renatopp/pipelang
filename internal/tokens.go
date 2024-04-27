package internal

import "github.com/renatopp/langtools/tokens"

var (
	T_EOF     = tokens.EOF
	T_INVALID = tokens.INVALID

	T_EOE        tokens.TokenType = "eoe" // end of expression
	T_KEYWORD    tokens.TokenType = "keyword"
	T_OPERATOR   tokens.TokenType = "operator"
	T_ASSIGNMENT tokens.TokenType = "assignment"
	T_IDENTIFIER tokens.TokenType = "identifier"
	T_NUMBER     tokens.TokenType = "number"
	T_HEX_NUMBER tokens.TokenType = "hex_number"
	T_BIN_NUMBER tokens.TokenType = "bin_number"
	T_OCT_NUMBER tokens.TokenType = "oct_number"
	T_STRING     tokens.TokenType = "string"
	T_BOOLEAN    tokens.TokenType = "boolean"
	T_COMMENT    tokens.TokenType = "comment"

	T_SPREAD tokens.TokenType = "spread"
	T_COMMA  tokens.TokenType = "comma"
	T_UNWRAP tokens.TokenType = "unwrap"
	T_WRAP   tokens.TokenType = "wrap"
	T_ACCESS tokens.TokenType = "access"
	T_LAMBDA tokens.TokenType = "lambda"
	T_PIPE   tokens.TokenType = "pipe"

	T_LPAREN tokens.TokenType = "lparen"
	T_RPAREN tokens.TokenType = "rparen"
	T_LBRACK tokens.TokenType = "lbrack"
	T_RBRACK tokens.TokenType = "rbrack"
	T_LBRACE tokens.TokenType = "lbrace"
	T_RBRACE tokens.TokenType = "rbrace"
)
