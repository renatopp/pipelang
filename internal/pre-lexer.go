package internal

import (
	"fmt"
	"slices"

	"github.com/renatopp/langtools/lexers"
	"github.com/renatopp/langtools/runes"
	"github.com/renatopp/langtools/tokens"
)

type PreLexer struct {
	*lexers.BaseLexer
	queue    []*tokens.Token
	finished bool
}

func NewPreLexer(input []byte) *PreLexer {
	p := &PreLexer{}

	p.BaseLexer = lexers.NewBaseLexer(input, func(bl *lexers.BaseLexer) *tokens.Token {
		return p.tokenizer()
	})
	p.queue = make([]*tokens.Token, 0)
	return p
}

func (p *PreLexer) tokenizer() *tokens.Token {
	for {
		c0 := p.PeekCharAt(0)
		c1 := p.PeekCharAt(1)
		c2 := p.PeekCharAt(2)
		c3 := p.PeekCharAt(3)

		s1 := string(c0.Rune)
		s2 := s1 + string(c1.Rune)
		s3 := s2 + string(c2.Rune)

		if p.HasTooManyErrors() || p.finished {
			return tokens.NewToken(T_EOF, "").WithRangeChars(c0, c1)
		}

		if len(p.queue) > 0 {
			token := p.queue[0]
			p.queue = p.queue[1:]
			return token
		}

		switch {
		// skip spaces
		case c0.IsOneOf(' ', '\t', '\r'):
			p.EatSpaces()
			continue

		// comment
		case s2 == "--":
			return p.EatUntilEndOfLine().WithType(T_COMMENT)

		// end of expression
		case c0.IsOneOf('\n', ';', 0):
			p.EatChar()
			if c0.Rune == 0 {
				p.finished = true
			}
			return tokens.NewToken(T_EOE, s1).WithRangeChars(c0, c1)

		case runes.IsAlpha(c0.Rune) || c0.Is('_'): // || c0.Is('$') || c0.Is('@'):
			token := p.EatIdentifierWith('_') //, '$', '@')

			switch {
			case token.Literal == "true" || token.Literal == "false":
				return token.WithType(T_BOOLEAN)

			// keywords
			case slices.Contains(Keywords, token.Literal):
				return token.WithType(T_KEYWORD)

			// operators
			case slices.Contains(Operators, token.Literal):
				return token.WithType(T_OPERATOR)

			// identifiers
			default:
				return token.WithType(T_IDENTIFIER)
			}

		// numbers
		case runes.IsNumeric(c0.Rune) || (c0.Is('.') && runes.IsNumeric(c1.Rune)):
			if c0.Is('0') && c1.IsOneOf('x', 'X') {
				return p.EatHexadecimal().WithType(T_HEX_NUMBER)
			}

			if c0.Is('0') && c1.IsOneOf('b', 'B') {
				return p.EatBinary().WithType(T_BIN_NUMBER)
			}

			if c0.Is('0') && c1.IsOneOf('o', 'O') {
				return p.EatOctal().WithType(T_OCT_NUMBER)
			}

			return p.EatNumber().WithType(T_NUMBER)

		// strings
		case c0.IsOneOf('"', '\''):
			return p.EatString().WithType(T_STRING)

		// raw strings
		case c0.Is('`'):
			return p.EatRawString().WithType(T_STRING)

			// spread
		case s3 == "...":
			p.EatChars(3)
			return tokens.NewToken(T_SPREAD, s3).WithRangeChars(c0, c3)

		// 3-length assignments
		case slices.Contains(Assignments, s3):
			p.EatChars(3)
			return tokens.NewToken(T_ASSIGNMENT, s3).WithRangeChars(c0, c3)

		// 3-length operators
		case slices.Contains(InfixOperators, s3):
			p.EatChars(3)
			return tokens.NewToken(T_OPERATOR, s3).WithRangeChars(c0, c3)

			// 2-length assignments
		case slices.Contains(Assignments, s2):
			p.EatChars(2)
			return tokens.NewToken(T_ASSIGNMENT, s2).WithRangeChars(c0, c2)

		// 2-length operators
		case slices.Contains(InfixOperators, s2):
			p.EatChars(2)
			return tokens.NewToken(T_OPERATOR, s2).WithRangeChars(c0, c2)

			// 1-length assignments
		case slices.Contains(Assignments, s1):
			p.EatChar()
			return tokens.NewToken(T_ASSIGNMENT, s1).WithRangeChars(c0, c1)

		// 1-length operators
		case slices.Contains(InfixOperators, s1):
			p.EatChar()
			return tokens.NewToken(T_OPERATOR, s1).WithRangeChars(c0, c1)

		// comma
		case c0.Is(','):
			p.EatChar()
			return tokens.NewToken(T_COMMA, s1).WithRangeChars(c0, c1)

		// unwrap
		case c0.Is('!'):
			p.EatChar()
			return tokens.NewToken(T_UNWRAP, s1).WithRangeChars(c0, c1)

		// wrap
		case c0.Is('?'):
			p.EatChar()
			return tokens.NewToken(T_WRAP, s1).WithRangeChars(c0, c1)

		// lambda
		case c0.Is(':'):
			p.EatChar()
			return tokens.NewToken(T_LAMBDA, s1).WithRangeChars(c0, c1)

		// pipe
		case c0.Is('|'):
			p.EatChar()
			return tokens.NewToken(T_PIPE, s1).WithRangeChars(c0, c1)

		// access
		case c0.Is('.'):
			p.EatChar()
			return tokens.NewToken(T_ACCESS, s1).WithRangeChars(c0, c1)

		// lparen
		case c0.Is('('):
			p.EatChar()
			return tokens.NewToken(T_LPAREN, s1).WithRangeChars(c0, c1)

		// rparen
		case c0.Is(')'):
			p.EatChar()
			return tokens.NewToken(T_RPAREN, s1).WithRangeChars(c0, c1)

		// lbrace
		case c0.Is('{'):
			p.EatChar()
			return tokens.NewToken(T_LBRACE, s1).WithRangeChars(c0, c1)

		// rbrace
		case c0.Is('}'):
			p.EatChar()
			return tokens.NewToken(T_RBRACE, s1).WithRangeChars(c0, c1)

			// lbrack
		case c0.Is('['):
			p.EatChar()
			return tokens.NewToken(T_LBRACK, s1).WithRangeChars(c0, c1)

			// rbrack
		case c0.Is(']'):
			p.EatChar()
			return tokens.NewToken(T_RBRACK, s1).WithRangeChars(c0, c1)

		// pipe
		case c0.Is('|'):
			p.EatChar()
			return tokens.NewToken(T_PIPE, s1).WithRangeChars(c0, c1)

		default:
			p.EatChar()
			p.RegisterError(fmt.Sprintf("invalid character '%s'", s1))
			return tokens.NewToken(T_INVALID, s1).WithRangeChars(c0, c1)
		}
	}
}

func (p *PreLexer) EatChars(n int) []tokens.Char {
	result := make([]tokens.Char, 0)
	for i := 0; i < n; i++ {
		result = append(result, p.EatChar())
	}
	return result
}
