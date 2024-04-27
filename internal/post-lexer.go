package internal

import (
	"github.com/renatopp/langtools/lexers"
	"github.com/renatopp/langtools/tokens"
)

type PostLexer struct {
	tokens []*tokens.Token
	cursor int
}

func NewPostLexer(tokens []*tokens.Token) *PostLexer {
	p := &PostLexer{}
	p.tokens = tokens
	return p
}

// INTERFACE ------------------------------------------------------------------
func (p *PostLexer) Errors() []lexers.LexerError {
	return []lexers.LexerError{}
}

func (p *PostLexer) HasErrors() bool {
	return false
}

func (p *PostLexer) EatToken() *tokens.Token {
	if p.cursor >= len(p.tokens) {
		return p.tokens[len(p.tokens)-1]
	}
	token := p.tokens[p.cursor]
	p.cursor++
	return token
}

func (p *PostLexer) PeekToken() *tokens.Token {
	return p.PeekTokenAt(0)
}

func (p *PostLexer) PrevToken() *tokens.Token {
	return p.tokens[max(p.cursor-1, 0)]
}

func (p *PostLexer) PrevTokenAt(i int) *tokens.Token {
	return p.tokens[max(p.cursor-i, 0)]
}

func (p *PostLexer) PeekTokenAt(i int) *tokens.Token {
	i = max(i, 0)
	if p.cursor+i >= len(p.tokens) {
		return p.tokens[len(p.tokens)-1]
	}
	return p.tokens[p.cursor+i]
}

func (p *PostLexer) Next() (token *tokens.Token, eof bool) {
	t := p.EatToken()
	if t.IsType(T_EOF) {
		return t, true
	}
	return t, false
}

// OPTIMIZATION ---------------------------------------------------------------
func (p *PostLexer) All() []*tokens.Token {
	return p.tokens
}

func (p *PostLexer) Optimize() error {
	result := make([]*tokens.Token, 0)

	lastComment := ""
	var firstCommentToken *tokens.Token
	var lastCommentToken *tokens.Token
	for i := 0; i < len(p.tokens); i++ {
		cur := p.tokens[i]

		// ------------------------------------------------------------------------
		// Merge comments to be used in documentations
		// ------------------------------------------------------------------------
		if cur.IsType(T_COMMENT) {
			if firstCommentToken == nil {
				lastComment = p.getComment(cur)
				firstCommentToken = cur
			} else {
				lastComment += "\n" + p.getComment(cur)
			}

			lastCommentToken = cur
			// Move the cursor to skip the next EOE
			i++
			continue

		} else if firstCommentToken != nil {
			fromLine, fromColumn, _, _ := firstCommentToken.Range()
			_, _, toLine, toColumn := lastCommentToken.Range()

			tokens.NewToken(T_COMMENT, lastComment).
				WithRange(fromLine, fromColumn, toLine, toColumn)

			// result = append(result, t) // remove while its not used
			lastComment = ""
			firstCommentToken = nil
		}

		// ------------------------------------------------------------------------
		// Remove unnecessary EOE tokens
		// ------------------------------------------------------------------------
		if cur.IsType(T_EOE) {
			for p.PeekTokenAt(i + 1).IsType(T_EOE) {
				i++
				continue
			}

			if !p.PeekTokenAt(i + 1).IsType(T_PIPE) {
				result = append(result, cur)
			}

			continue
		}

		// ------------------------------------------------------------------------
		// Merge and organize pipes
		// ------------------------------------------------------------------------
		if cur.IsType(T_PIPE) {
			// TODO: remove previous EOE if any
			// TODO: add () to call

			for p.PeekTokenAt(i+1).IsOneOfTypes(T_PIPE, T_EOE) {
				i++
				continue
			}

			result = append(result, cur)
			continue
		}

		result = append(result, cur)
	}

	p.tokens = result

	return nil
}

func (p *PostLexer) getComment(token *tokens.Token) string {
	if token.Literal[2] == ' ' {
		return token.Literal[3:]
	}

	return token.Literal[2:]
}
