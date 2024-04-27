package internal

import (
	"fmt"
	"io"
	"log"
	"slices"
	"strconv"
	"strings"

	"github.com/renatopp/langtools/parsers"
	"github.com/renatopp/langtools/tokens"
	"github.com/renatopp/pipelang/internal/ast"
)

type PrefixFn func() ast.Node
type InfixFn func(left ast.Node) ast.Node
type PostfixFn func(left ast.Node) ast.Node

type PipeParser struct {
	*parsers.BaseParser
	Lexer *PostLexer
	Log   *log.Logger

	prefixFns  map[tokens.TokenType]PrefixFn
	infixFns   map[tokens.TokenType]InfixFn
	postfixFns map[tokens.TokenType]PostfixFn

	yieldStack *Stack[bool] // to check if fn is generator
	// TODO: how to handle feature enabling/disabling
	openTupleLock *Stack[bool] // lock to disable open tuples, use true to lock
	conditionLock *Stack[bool] // lock to disable data instantiation and dict creation, use true to lock
	lambdaLock    *Stack[bool] // lock to disable lambdas, use true to lock
	pipeLock      *Stack[bool] // lock to disable pipes, use true to lock
}

func NewPipeParser(lexer *PostLexer) *PipeParser {
	p := &PipeParser{}
	p.BaseParser = parsers.NewBaseParser(lexer)
	p.Lexer = lexer
	p.Log = log.New(io.Discard, "", 0)
	p.prefixFns = make(map[tokens.TokenType]PrefixFn)
	p.infixFns = make(map[tokens.TokenType]InfixFn)
	p.postfixFns = make(map[tokens.TokenType]PostfixFn)
	p.yieldStack = NewStack[bool]()

	p.openTupleLock = NewStack[bool]()
	p.conditionLock = NewStack[bool]()
	p.lambdaLock = NewStack[bool]()
	p.pipeLock = NewStack[bool]()

	p.registerPrefixFn(T_NUMBER, p.prefixNumber)
	p.registerPrefixFn(T_HEX_NUMBER, p.prefixHexNumber)
	p.registerPrefixFn(T_BIN_NUMBER, p.prefixBinNumber)
	p.registerPrefixFn(T_OCT_NUMBER, p.prefixOctNumber)
	p.registerPrefixFn(T_STRING, p.prefixString)
	p.registerPrefixFn(T_BOOLEAN, p.prefixBoolean)
	p.registerPrefixFn(T_IDENTIFIER, p.prefixIdentifier)
	p.registerPrefixFn(T_OPERATOR, p.prefixOperator)
	p.registerPrefixFn(T_LPAREN, p.prefixParenthesis)
	p.registerPrefixFn(T_SPREAD, p.prefixSpread)
	p.registerPrefixFn(T_LBRACK, p.prefixBracket)
	p.registerPrefixFn(T_LBRACE, p.prefixBrace)
	p.registerPrefixFn(T_LAMBDA, p.prefixLambda)
	p.registerPrefixFn(T_KEYWORD, p.prefixKeyword)

	p.registerInfixFn(T_OPERATOR, p.infixOperator)
	p.registerInfixFn(T_LPAREN, p.infixParenthesis)
	p.registerInfixFn(T_ASSIGNMENT, p.infixAssignment)
	p.registerInfixFn(T_LAMBDA, p.infixLambda)
	p.registerInfixFn(T_ACCESS, p.infixAccess)
	p.registerInfixFn(T_KEYWORD, p.infixKeyword)
	p.registerInfixFn(T_COMMA, p.infixComma)
	p.registerInfixFn(T_PIPE, p.infixPipe)
	p.registerInfixFn(T_LBRACE, p.infixBrace)
	p.registerInfixFn(T_LBRACK, p.infixBracket)

	p.registerPostfixFn(T_SPREAD, p.postfixSpread)
	p.registerPostfixFn(T_WRAP, p.postfixWrap)
	p.registerPostfixFn(T_UNWRAP, p.postfixUnwrap)

	return p
}

// ----------------------------------------------------------------------------
// Interface
// ----------------------------------------------------------------------------

func (p *PipeParser) isEndOfBlock(t *tokens.Token) bool {
	return t.IsType(T_RBRACE) || t.IsType(T_EOF)
}

func (p *PipeParser) isEndOfExpr(token *tokens.Token) bool {
	return token.IsType(T_EOE)
}

func (p *PipeParser) registerPrefixFn(tokenType tokens.TokenType, fn PrefixFn) {
	p.prefixFns[tokenType] = fn
}

func (p *PipeParser) registerInfixFn(tokenType tokens.TokenType, fn InfixFn) {
	p.infixFns[tokenType] = fn
}

func (p *PipeParser) registerPostfixFn(tokenType tokens.TokenType, fn PostfixFn) {
	p.postfixFns[tokenType] = fn
}

func (p *PipeParser) precedence(t *tokens.Token) int {
	switch {
	case t.IsType(T_ASSIGNMENT):
		return 2

	case t.IsType(T_PIPE):
		return 3

	case t.IsType(T_KEYWORD) && t.IsOneOfLiterals("as"):
		return 4

	case t.IsType(T_COMMA):
		return 10

	case t.IsType(T_LAMBDA):
		return 20

	case t.IsType(T_OPERATOR) && t.IsOneOfLiterals(".."):
		return 28

	case t.IsType(T_OPERATOR) && t.IsOneOfLiterals("??"):
		return 29

	case t.IsType(T_OPERATOR) && t.IsOneOfLiterals("or", "xor"):
		return 30

	case t.IsType(T_OPERATOR) && t.IsOneOfLiterals("and"):
		return 31

	case t.IsType(T_OPERATOR) && t.IsOneOfLiterals("not"):
		return 32

	case t.IsType(T_OPERATOR) && t.IsOneOfLiterals("==", "!="):
		return 40

	case t.IsType(T_OPERATOR) && t.IsOneOfLiterals("<", ">", "<=", ">=", "<=>"):
		return 41

	case t.IsType(T_OPERATOR) && t.IsOneOfLiterals("+", "-"):
		return 50

	case t.IsType(T_OPERATOR) && t.IsOneOfLiterals("*", "/"):
		return 51

	case t.IsType(T_OPERATOR) && t.IsOneOfLiterals("%"):
		return 52

	case t.IsType(T_OPERATOR) && t.IsOneOfLiterals("^"):
		return 53

	case t.IsType(T_UNWRAP):
		return 60

	case t.IsType(T_WRAP):
		return 61

	case t.IsType(T_LPAREN):
		return 62

	case t.IsType(T_LBRACE):
		return 62

	case t.IsType(T_SPREAD):
		return 63

	case t.IsType(T_ACCESS):
		return 70

	case t.IsType(T_LBRACK):
		return 70

	default:
		return 0
	}
}

func (p *PipeParser) Parse() ast.Node {
	return p.parseBlock()
}

// ----------------------------------------------------------------------------
// Statements functions
// ----------------------------------------------------------------------------
func (p *PipeParser) parseBlock() *ast.Block {
	p.openTupleLock.Push(false)
	defer p.openTupleLock.Pop()
	p.conditionLock.Push(false)
	defer p.conditionLock.Pop()
	p.lambdaLock.Push(false)
	defer p.lambdaLock.Pop()
	p.pipeLock.Push(false)
	defer p.pipeLock.Pop()

	// Removes the { token
	cur := p.Lexer.PeekToken()
	braced := cur.IsType(T_LBRACE)
	if braced {
		p.Lexer.EatToken()
	}

	// Parses the block
	first := p.Lexer.PeekToken()
	expressions := []ast.Node{}
	for {
		// Skip unused EOE tokens
		p.skipEoes()

		// Check if the block is finished
		cur := p.Lexer.PeekToken()
		if p.isEndOfBlock(cur) || p.Lexer.HasErrors() || p.HasErrors() {
			break
		}

		// Parse the statement
		expr := p.parseOptionalExpression()
		if expr != nil {
			expressions = append(expressions, expr)
		}

		p.ExpectTypes(T_EOE, T_EOF, T_RBRACE)
	}

	// Checks for the } token
	if braced {
		p.ExpectType(T_RBRACE)
		p.Lexer.EatToken()
	}

	return &ast.Block{
		Token:       first,
		Expressions: expressions,
		Scoped:      true,
	}
}

// Parse single expressions eg: `f(20)`. MAY RETURN NIL
func (p *PipeParser) parseRequiredExpression() ast.Node {
	expr := p.parseExpression(0)
	if expr == nil {
		t := p.Lexer.PrevToken()
		p.RegisterErrorWithToken(fmt.Sprintf("expected expression after '%s'", escapeError(t.Literal)), p.Lexer.PeekToken())
	}

	return expr
}

// Parse single expressions eg: `f(20)`. MAY RETURN NIL!
func (p *PipeParser) parseOptionalExpression() ast.Node {
	return p.parseExpression(0)
}

// Parse expressions separated by ;, eg: `1 ; f() ; 2+3; ...`. MAY RETURN NIL!
func (p *PipeParser) parseExpressionStatements() []ast.Node {
	expressions := []ast.Node{}

	for {
		p.skipEoes()

		if p.Lexer.HasErrors() || p.HasErrors() {
			break
		}

		expr := p.parseOptionalExpression()
		if expr == nil {
			break
		}
		expressions = append(expressions, expr)

		p.skipEoes()
		if p.Lexer.PeekToken().IsLiteral(";") {
			p.Lexer.EatToken()
		}
	}

	return expressions
}

// Parse expressions separated by commas, eg: `1, 2+3, f(), ...`
func (p *PipeParser) parseExpressionList(precedence int) []ast.Node {
	expressions := []ast.Node{}

	for {
		if p.Lexer.HasErrors() || p.HasErrors() {
			break
		}

		expr := p.parseExpression(precedence)
		if expr == nil {
			break
		}
		expressions = append(expressions, expr)

		if !p.Lexer.PeekToken().IsType(T_COMMA) {
			break
		}
		p.skipEoes()
		p.Lexer.EatToken()
	}

	return expressions
}

func (p *PipeParser) parseExpression(precedence int) ast.Node {
	// println("...parseExpression", p.Lexer.PeekToken().Literal, precedence)
	prefix := p.prefixFns[p.Lexer.PeekToken().Type]
	if prefix == nil {
		return nil
	}
	left := prefix()

	cur := p.Lexer.PeekToken()
	for {
		starting := cur

		if !p.isEndOfExpr(cur) && precedence < p.precedence(cur) {
			infix := p.infixFns[cur.Type]
			if infix != nil {
				n := infix(left)
				if n != nil {
					left = n
				}
				cur = p.Lexer.PeekToken()
			}
		}

		for {
			postfix := p.postfixFns[cur.Type]
			if postfix == nil {
				break
			}

			newLeft := postfix(left)
			if newLeft == nil {
				break
			}
			left = newLeft
			cur = p.Lexer.PeekToken()
		}

		// Didn't find any infix or postfix function
		if starting == cur {
			break
		}
	}

	return left
}

// ----------------------------------------------------------------------------
// Prefix functions
// ----------------------------------------------------------------------------

func (p *PipeParser) prefixNumber() ast.Node {
	cur := p.Lexer.EatToken()
	value, err := strconv.ParseFloat(cur.Literal, 64)
	if err != nil {
		p.RegisterErrorWithToken("invalid number literal", cur)
	}

	return &ast.Number{
		Token: cur,
		Value: value,
	}
}

func (p *PipeParser) prefixHexNumber() ast.Node {
	cur := p.Lexer.EatToken()
	value, err := strconv.ParseInt(cur.Literal, 16, 64)
	if err != nil {
		p.RegisterErrorWithToken("invalid hexadecimal literal", cur)
	}

	return &ast.Number{
		Token: cur,
		Value: float64(value),
	}
}

func (p *PipeParser) prefixOctNumber() ast.Node {
	cur := p.Lexer.EatToken()
	value, err := strconv.ParseInt(cur.Literal, 8, 64)
	if err != nil {
		p.RegisterErrorWithToken("invalid octal literal", cur)
	}

	return &ast.Number{
		Token: cur,
		Value: float64(value),
	}
}

func (p *PipeParser) prefixBinNumber() ast.Node {
	cur := p.Lexer.EatToken()
	value, err := strconv.ParseInt(cur.Literal, 2, 64)
	if err != nil {
		p.RegisterErrorWithToken("invalid binary literal", cur)
	}

	return &ast.Number{
		Token: cur,
		Value: float64(value),
	}
}

func (p *PipeParser) prefixString() ast.Node {
	cur := p.Lexer.EatToken()
	return &ast.String{
		Token: cur,
		Value: cur.Literal,
	}
}

func (p *PipeParser) prefixBoolean() ast.Node {
	cur := p.Lexer.EatToken()
	value := cur.Literal == "true"
	return &ast.Boolean{
		Token: cur,
		Value: value,
	}
}

func (p *PipeParser) prefixIdentifier() ast.Node {
	cur := p.Lexer.EatToken()
	return &ast.Identifier{
		Token: cur,
		Value: cur.Literal,
	}
}

func (p *PipeParser) prefixOperator() ast.Node {
	cur := p.Lexer.EatToken()
	right := p.parseExpression(p.precedence(cur))

	if right == nil {
		p.RegisterErrorWithToken("expected expression", cur)
		return nil
	}

	return &ast.PrefixOperator{
		Token:    cur,
		Operator: cur.Literal,
		Right:    right,
	}
}

func (p *PipeParser) prefixParenthesis() ast.Node {
	p.openTupleLock.Push(true)
	defer p.openTupleLock.Pop()

	cur := p.Lexer.EatToken()
	p.conditionLock.Push(false)
	expr := p.parseExpressionList(0)
	p.conditionLock.Pop()
	p.ExpectType(T_RPAREN)
	p.Lexer.EatToken()

	if len(expr) == 0 {
		// TODO: what to do with an empty tuple?
		p.RegisterErrorWithToken("expected expression: tuples cannot be empty", cur)
		return nil
	}

	if len(expr) == 1 {
		return expr[0]
	}

	return &ast.Tuple{
		Token:    cur,
		Elements: expr,
	}
}

func (p *PipeParser) prefixSpread() ast.Node {
	cur := p.Lexer.EatToken()
	right := p.parseExpression(p.precedence(cur))

	if right == nil {
		p.RegisterErrorWithToken("expected expression after a spread", cur)
		return nil
	}

	return &ast.Spread{
		Token:  cur,
		Target: right,
		In:     true,
	}
}

func (p *PipeParser) prefixBracket() ast.Node {
	p.openTupleLock.Push(true)
	defer p.openTupleLock.Pop()

	cur := p.Lexer.EatToken()
	expr := p.parseExpressionList(0)
	p.ExpectType(T_RBRACK)
	p.Lexer.EatToken()

	return &ast.List{
		Token:    cur,
		Elements: expr,
	}
}

// Parse dict definitions {a=1, b=2, c=3}
func (p *PipeParser) prefixBrace() ast.Node {
	if p.conditionLock.PeekOr(true) {
		return nil
	}

	p.openTupleLock.Push(true)
	defer p.openTupleLock.Pop()

	cur := p.Lexer.EatToken()
	elements := []ast.Node{}
	for {
		p.skipEoes()

		if p.Lexer.PeekToken().IsType(T_RBRACE) {
			break
		}

		cur := p.Lexer.EatToken()
		if !cur.IsOneOfTypes(T_IDENTIFIER, T_STRING, T_NUMBER) {
			p.RegisterErrorWithToken("expected key for dictionary", cur)
			break
		}
		key := &ast.String{Token: cur, Value: cur.Literal}

		if !p.Lexer.EatToken().IsLiteral("=") {
			p.RegisterErrorWithToken("expected '=' after dictionary key", key.GetToken())
			break
		}

		value := p.parseRequiredExpression()
		if p.Lexer.PeekToken().IsType(T_COMMA) {
			p.Lexer.EatToken()
		}

		p.skipEoes()
		elements = append(elements, key, value)
	}

	p.ExpectType(T_RBRACE)
	p.Lexer.EatToken()

	return &ast.Dict{
		Token:    cur,
		Elements: elements,
	}
}

func (p *PipeParser) prefixLambda() ast.Node {
	cur := p.Lexer.PeekToken()
	params := &ast.Tuple{Token: cur, Elements: []ast.Node{}}
	return p.infixLambda(params)
}

func (p *PipeParser) prefixKeyword() ast.Node {
	cur := p.Lexer.PeekToken()

	switch cur.Literal {
	case "fn":
		return p.prefixFunction()

	case "return":
		return p.prefixReturn()

	case "raise":
		return p.prefixRaise()

	case "yield":
		return p.prefixYield()

	case "break":
		return p.prefixBreak()

	case "continue":
		return p.prefixContinue()

	case "if":
		return p.prefixIf()

	case "for":
		return p.prefixFor()

	case "with":
		return p.prefixWith()

	case "data":
		return p.prefixData()

	case "match":
		return p.prefixMatch()

	case "in":
		return nil
	}

	p.RegisterErrorWithToken("unexpected prefix keyword", cur)
	return nil
}

func (p *PipeParser) prefixFunction() ast.Node {
	first := p.Lexer.EatToken()

	name := ""
	params := []ast.Node{}
	var body *ast.Block

	// Parse the function name
	cur := p.Lexer.PeekToken()
	if cur.IsType(T_IDENTIFIER) {
		name = cur.Literal
		p.Lexer.EatToken()
	}

	// Parse the parameters (...)
	p.conditionLock.Push(true)
	cur = p.Lexer.PeekToken()
	if cur.IsType(T_LPAREN) {
		p.openTupleLock.Push(true)
		defer p.openTupleLock.Pop()

		p.Lexer.EatToken()
		params = p.parseExpressionList(0)
		p.skipEoes()
		p.ExpectType(T_RPAREN)
		p.Lexer.EatToken()
	}
	p.conditionLock.Pop()

	// Parse the body {...}
	p.yieldStack.Push(false)
	defer p.yieldStack.Pop()
	cur = p.Lexer.PeekToken()
	if cur.IsType(T_LBRACE) {
		body = p.parseBlock()
	}

	if body == nil {
		p.RegisterErrorWithToken("expected function body", cur)
		return nil
	}

	p.validateFunctionParameter(params)
	p.validateFunctionBody(body)
	return &ast.FunctionDef{
		Token:      first,
		Name:       name,
		Parameters: params,
		Body:       body,
		Generator:  p.yieldStack.Peek(),
	}
}

func (p *PipeParser) prefixReturn() ast.Node {
	cur := p.Lexer.EatToken()
	expr := p.parseOptionalExpression()
	if expr == nil {
		expr = &ast.Boolean{Token: cur, Value: false}
	}
	return &ast.Return{
		Token:      cur,
		Expression: expr,
	}
}

func (p *PipeParser) prefixRaise() ast.Node {
	cur := p.Lexer.EatToken()
	expr := p.parseOptionalExpression()
	if expr == nil {
		expr = &ast.Boolean{Token: cur, Value: false}
	}
	return &ast.Raise{
		Token:      cur,
		Expression: expr,
	}
}

func (p *PipeParser) prefixYield() ast.Node {
	cur := p.Lexer.EatToken()

	if next := p.Lexer.PeekToken(); next.IsLiteral("break") {
		p.Lexer.EatToken()
		return &ast.Yield{
			Token: cur,
			Break: true,
		}
	}

	expr := p.parseRequiredExpression()
	if expr == nil {
		expr = &ast.Boolean{Token: cur, Value: false}
	}

	p.yieldStack.Set(true)
	return &ast.Yield{
		Token:      cur,
		Expression: expr,
	}
}

func (p *PipeParser) prefixBreak() ast.Node {
	cur := p.Lexer.EatToken()

	return &ast.Break{
		Token: cur,
	}
}

func (p *PipeParser) prefixContinue() ast.Node {
	cur := p.Lexer.EatToken()

	return &ast.Continue{
		Token: cur,
	}
}

func (p *PipeParser) prefixIf() ast.Node {
	cur := p.Lexer.EatToken()
	if !p.NotExpectTypes(T_LBRACE, T_EOF, T_EOE) {
		return nil
	}

	p.conditionLock.Push(true)
	conditions := p.parseExpressionStatements()
	p.conditionLock.Pop()
	if len(conditions) == 0 {
		p.RegisterErrorWithToken("expected at least one condition", cur)
		return nil
	}

	p.ExpectTypes(T_LBRACE)
	trueExpression := p.parseBlock()

	var falseExpression ast.Node = nil
	if p.Lexer.PeekToken().IsLiteral("else") {
		p.Lexer.EatToken()

		if p.Lexer.PeekToken().IsLiteral("if") {
			falseExpression = p.prefixIf()
		} else {
			falseExpression = p.parseBlock()
		}
	} else {
		falseExpression = &ast.Block{
			Token:       cur,
			Expressions: []ast.Node{},
			Scoped:      true,
		}
	}

	return &ast.If{
		Token:           cur,
		Conditions:      conditions,
		TrueExpression:  trueExpression,
		FalseExpression: falseExpression,
	}
}

func (p *PipeParser) prefixFor() ast.Node {
	cur := p.Lexer.EatToken()

	p.conditionLock.Push(true)
	conditions := p.parseExpressionStatements()
	if len(conditions) == 0 {
		conditions = []ast.Node{&ast.Boolean{Token: cur, Value: true}}
	}

	var in ast.Node
	if p.Lexer.PeekToken().IsLiteral("in") {
		p.Lexer.EatToken()
		in = p.parseRequiredExpression()
	}
	p.conditionLock.Pop()

	p.ExpectTypes(T_LBRACE)
	expression := p.parseBlock()

	if in != nil {
		if len(conditions) == 0 {
			p.RegisterErrorWithToken("expected left side of `in` expression", cur)
		} else {
			p.validateAssignmentLeft(conditions[len(conditions)-1])
		}
	}

	return &ast.For{
		Token:        cur,
		Conditions:   conditions,
		InExpression: in,
		Expression:   expression,
	}
}

func (p *PipeParser) prefixWith() ast.Node {
	cur := p.Lexer.EatToken()
	if !p.NotExpectTypes(T_EOF, T_EOE) {
		return nil
	}

	p.conditionLock.Push(true)
	condition := p.parseOptionalExpression()
	p.conditionLock.Pop()
	p.ExpectTypes(T_LBRACE)
	expression := p.parseBlock()

	return &ast.With{
		Token:      cur,
		Condition:  condition,
		Expression: expression,
	}
}

func (p *PipeParser) prefixData() ast.Node {
	first := p.Lexer.EatToken()

	name := ""
	extensions := []ast.Node{}
	var body *ast.Block

	// Parse the function name
	cur := p.Lexer.PeekToken()
	if cur.IsType(T_IDENTIFIER) {
		name = cur.Literal
		p.Lexer.EatToken()
	}

	// Parse the parameters (...)

	p.conditionLock.Push(true)
	cur = p.Lexer.PeekToken()
	if cur.IsType(T_LPAREN) {
		p.openTupleLock.Push(true)
		defer p.openTupleLock.Pop()

		p.Lexer.EatToken()
		extensions = p.parseExpressionList(0)
		p.skipEoes()
		p.ExpectType(T_RPAREN)
		p.Lexer.EatToken()
	}
	p.conditionLock.Pop()

	// Parse the body {...}
	p.yieldStack.Push(false)
	defer p.yieldStack.Pop()
	cur = p.Lexer.PeekToken()
	if cur.IsType(T_LBRACE) {
		body = p.parseBlock()
	}

	if body == nil {
		p.RegisterErrorWithToken("expected function body", cur)
		return nil
	}

	p.validateDataBody(body)
	attributes := map[string]ast.Node{}
	methods := map[string]ast.Node{}
	for _, expr := range body.Expressions {
		switch expr := expr.(type) {
		case *ast.Assignment:
			name := expr.Left.(*ast.Identifier).Value
			attributes[name] = expr.Right

		case *ast.FunctionDef:
			methods[expr.Name] = expr
		}
	}

	return &ast.DataDef{
		Token:      first,
		Name:       name,
		Extensions: extensions,
		Attributes: attributes,
		Methods:    methods,
	}
}

func (p *PipeParser) prefixMatch() ast.Node {
	cur := p.Lexer.EatToken()

	p.conditionLock.Push(true)
	compare := p.parseOptionalExpression()
	p.conditionLock.Pop()
	if compare == nil {
		compare = &ast.Boolean{Token: cur, Value: true}
	}

	p.ExpectTypes(T_LBRACE)
	p.Lexer.EatToken()

	p.lambdaLock.Push(true)
	defer p.lambdaLock.Pop()

	cases := []ast.Node{}
	for {
		if p.Lexer.PeekToken().IsType(T_RBRACE) {
			p.Lexer.EatToken()
			break
		}

		p.skipEoes()
		match := p.parseRequiredExpression()
		if match == nil {
			break
		}
		p.ExpectType(T_LAMBDA)
		p.Lexer.EatToken()

		var expression ast.Node
		if p.Lexer.PeekToken().IsType(T_LBRACE) {
			expression = p.parseBlock()
		} else {
			expression = p.parseRequiredExpression()

		}
		p.skipEoes()

		cases = append(cases, match, expression)
	}

	return &ast.Match{
		Token:      cur,
		Expression: compare,
		Cases:      cases,
	}
}

// ----------------------------------------------------------------------------
// Infix functions
// ----------------------------------------------------------------------------

func (p *PipeParser) infixOperator(left ast.Node) ast.Node {
	cur := p.Lexer.EatToken()
	right := p.parseExpression(p.precedence(cur))

	if right == nil {
		p.RegisterErrorWithToken("expected expression", cur)
		return nil
	}

	return &ast.InfixOperator{
		Token:    cur,
		Operator: cur.Literal,
		Left:     left,
		Right:    right,
	}
}

func (p *PipeParser) infixParenthesis(left ast.Node) ast.Node {
	p.openTupleLock.Push(true)
	defer p.openTupleLock.Pop()

	cur := p.Lexer.EatToken()
	right := p.parseExpressionList(0)
	p.ExpectType(T_RPAREN)
	p.Lexer.EatToken()

	return &ast.Call{
		Token:     cur,
		Target:    left,
		Arguments: right,
	}
}

func (p *PipeParser) infixAssignment(left ast.Node) ast.Node {
	cur := p.Lexer.EatToken()

	p.conditionLock.Push(false)
	right := p.parseExpression(p.precedence(cur))
	p.conditionLock.Pop()

	if right == nil {
		p.RegisterErrorWithToken("expected expression", cur)
		return nil
	}

	p.validateAssignmentLeft(left)

	return &ast.Assignment{
		Token:    cur,
		Operator: cur.Literal,
		Left:     left,
		Right:    right,
	}
}

func (p *PipeParser) infixLambda(left ast.Node) ast.Node {
	if p.lambdaLock.PeekOr(true) {
		return nil
	}

	cur := p.Lexer.EatToken()

	params := []ast.Node{}
	switch left := left.(type) {
	case *ast.Tuple:
		params = left.Elements

	case *ast.Identifier:
		params = append(params, left)

	default:
		p.RegisterErrorWithToken("expected tuple or identifier as lambda parameters", left.GetToken())
		return nil
	}

	p.yieldStack.Push(false)
	defer p.yieldStack.Pop()

	var body ast.Node
	if p.Lexer.PeekToken().IsType(T_LBRACE) {
		body = p.parseBlock()
	} else {
		body = p.parseRequiredExpression()
	}

	p.validateFunctionParameter(params)
	p.validateFunctionBody(body)
	return &ast.FunctionDef{
		Token:      cur,
		Name:       "",
		Parameters: params,
		Body:       body,
		Generator:  p.yieldStack.Peek(),
	}
}

func (p *PipeParser) infixAccess(left ast.Node) ast.Node {
	cur := p.Lexer.EatToken()
	right := p.parseExpression(p.precedence(cur))

	if right == nil {
		p.RegisterErrorWithToken("expected expression", cur)
		return nil
	} else if _, ok := right.(*ast.Identifier); !ok {
		p.RegisterErrorWithToken("expected identifier as right side of an access", right.GetToken())
		return nil
	}

	return &ast.Access{
		Token: cur,
		Left:  left,
		Right: right,
	}
}

func (p *PipeParser) infixKeyword(left ast.Node) ast.Node {
	cur := p.Lexer.PeekToken()
	switch cur.Literal {
	case "as":
		return p.infixAs(left)

	case "is":
		return p.infixIs(left)

	case "in":
		return nil
	}

	p.RegisterErrorWithToken("unexpected keyword", cur)
	return nil
}

func (p *PipeParser) infixComma(left ast.Node) ast.Node {
	locked := p.openTupleLock.PeekOr(false)
	if locked {
		return nil
	}

	p.openTupleLock.Push(true)
	defer p.openTupleLock.Pop()

	cur := p.Lexer.EatToken()
	right := p.parseExpressionList(p.precedence(cur))

	return &ast.Tuple{
		Token:    left.GetToken(),
		Elements: append([]ast.Node{left}, right...),
	}
}

func (p *PipeParser) infixAs(left ast.Node) ast.Node {
	cur := p.Lexer.EatToken()

	right := p.parseRequiredExpression()
	if right == nil {
		return nil
	}

	p.validateAssignmentLeft(right)
	return &ast.Assignment{
		Token:    cur,
		Operator: ":=",
		Left:     right,
		Right:    left,
	}
}

func (p *PipeParser) infixPipe(left ast.Node) ast.Node {
	if p.pipeLock.PeekOr(false) {
		return nil
	}

	p.pipeLock.Push(true)
	defer p.pipeLock.Pop()

	cur := p.Lexer.EatToken()

	p.openTupleLock.Push(true)
	defer p.openTupleLock.Pop()

	p.ExpectType(T_IDENTIFIER)
	id := p.Lexer.EatToken()
	var target ast.Node = &ast.Identifier{Token: id, Value: id.Literal}
	for p.Lexer.PeekToken().IsType(T_ACCESS) {
		target = p.infixAccess(target)
	}

	args := p.parseExpressionList(p.precedence(cur))

	return &ast.Call{
		Token:     cur,
		Target:    target,
		Arguments: slices.Concat([]ast.Node{left}, args),
	}
}

func (p *PipeParser) infixIs(left ast.Node) ast.Node {
	p.Lexer.EatToken()
	return nil
}

func (p *PipeParser) infixBrace(left ast.Node) ast.Node {
	if p.conditionLock.PeekOr(true) {
		return nil
	}

	d := p.prefixBrace().(*ast.Dict)

	return &ast.Instantiate{
		Token:    d.Token,
		Target:   left,
		Elements: d.Elements,
	}
}

func (p *PipeParser) infixBracket(left ast.Node) ast.Node {
	cur := p.Lexer.EatToken()
	index := p.parseRequiredExpression()
	if index == nil {
		return nil
	}

	p.ExpectType(T_RBRACK)
	p.Lexer.EatToken()

	return &ast.Index{
		Token:  cur,
		Target: left,
		Index:  index,
	}
}

// ----------------------------------------------------------------------------
// Postfix functions
// ----------------------------------------------------------------------------

func (p *PipeParser) postfixSpread(left ast.Node) ast.Node {
	cur := p.Lexer.EatToken()
	return &ast.Spread{
		Token:  cur,
		Target: left,
		In:     false,
	}
}

func (p *PipeParser) postfixWrap(left ast.Node) ast.Node {
	cur := p.Lexer.EatToken()
	return &ast.Wrap{
		Token:  cur,
		Target: left,
	}
}

func (p *PipeParser) postfixUnwrap(left ast.Node) ast.Node {
	cur := p.Lexer.EatToken()
	return &ast.Unwrap{
		Token:  cur,
		Target: left,
	}
}

// ----------------------------------------------------------------------------
// Utilities
// ----------------------------------------------------------------------------
func (p *PipeParser) skipEoes() {
	for p.Lexer.PeekToken().IsType(T_EOE) {
		p.Lexer.EatToken()
	}
}

// Checks if the next token is the given types.
func (p *PipeParser) ExpectTypes(expected ...tokens.TokenType) bool {
	cur := p.Lexer.PeekToken()
	for _, t := range expected {
		if cur.IsType(t) {
			return true
		}
	}

	e := make([]string, len(expected))
	for i, t := range expected {
		e[i] = string(t)
	}

	p.RegisterErrorWithToken(fmt.Sprintf("expected one of the following tokens: [%s]", strings.Join(e, ", ")), cur)
	return false
}

// Checks if the next token is the given types.
func (p *PipeParser) NotExpectTypes(unexpected ...tokens.TokenType) bool {
	cur := p.Lexer.PeekToken()
	for _, t := range unexpected {
		if cur.IsType(t) {
			p.RegisterErrorWithToken(fmt.Sprintf("unexpected token '%s'", escapeError(cur.Literal)), cur)
			return false
		}
	}

	return true
}

func (p *PipeParser) validateAssignmentLeft(left ast.Node) bool {
	hasSpread := false

	switch left := left.(type) {
	case *ast.Identifier:
		// ok

	case *ast.Index:
		if !p.validateAssignmentLeft(left.Target) {
			return false
		}

	case *ast.Access:
		if _, ok := left.Right.(*ast.Identifier); !ok {
			p.RegisterErrorWithToken("expected identifier as right side of an access", left.GetToken())
		}

		if !p.validateAssignmentLeft(left.Left) {
			return false
		}

	case *ast.Spread:
		if !left.In {
			p.RegisterErrorWithToken("spread-out operator (a...) cannot be used in left-side of assignments. Use spread-in (...a) instead", left.GetToken())
			return false
		}

		p.validateAssignmentLeft(left.Target)

	case *ast.Tuple:
		for _, el := range left.Elements {
			if _, ok := el.(*ast.Spread); ok {
				if hasSpread {
					p.RegisterErrorWithToken("only one spread-in (...a) operator is allowed in the left side of an assignment", el.GetToken())
					return false
				}
				hasSpread = true
			}

			if !p.validateAssignmentLeft(el) {
				return false
			}
		}

	default:
		p.RegisterErrorWithToken(fmt.Sprintf("expected identifier, received %s instead", escapeError(left.GetToken().Literal)), left.GetToken())
		return false
	}

	return true
}

func (p *PipeParser) validateFunctionParameter(params []ast.Node) bool {
	hasSpread := false

	for _, param := range params {
		switch param := param.(type) {
		case *ast.Identifier:
			// ok

		case *ast.Spread:
			if !param.In {
				p.RegisterErrorWithToken("spread-out operator (a...) cannot be used in function parameters. Use spread-in (...a) instead", param.GetToken())
				return false
			}

			if hasSpread {
				p.RegisterErrorWithToken("only one spread-in (...a) operator is allowed in the function parameters", param.GetToken())
				return false
			}

			if _, ok := param.Target.(*ast.Identifier); !ok {
				p.RegisterErrorWithToken("expected identifier after spread-in operator (...)", param.GetToken())
				return false
			}

			hasSpread = true

		default:
			p.RegisterErrorWithToken(fmt.Sprintf("expected identifier, received %s instead", escapeError(param.GetToken().Literal)), param.GetToken())
			return false
		}
	}

	return true
}

func (p *PipeParser) validateFunctionBody(body ast.Node) bool {
	if body == nil {
		p.RegisterErrorWithToken("expected function body", p.Lexer.PeekToken())
		return false
	}

	returns := []*ast.Return{}
	yields := []*ast.Yield{}

	ast.Traverse(body, func(depth int, node ast.Node) {
		switch node := node.(type) {
		case *ast.Return:
			returns = append(returns, node)

		case *ast.Yield:
			yields = append(yields, node)
		}
	})

	// Generator functions cannot have returns
	if len(yields) > 0 {
		for _, ret := range returns {
			p.RegisterErrorWithToken("generator functions cannot return elements. Use `yield break` to stop the function instead.", ret.GetToken())
		}
	}

	return true
}

func (p *PipeParser) validateDataBody(body ast.Node) bool {
	if body == nil {
		p.RegisterErrorWithToken("expected data body", p.Lexer.PeekToken())
		return false
	}

	declarations := []string{}

	for _, node := range body.(*ast.Block).Expressions {
		switch node := node.(type) {
		case *ast.FunctionDef:
			if node.Name == "" {
				p.RegisterErrorWithToken("data definitions cannot have anonymous methods", node.GetToken())
			}

			if len(node.Parameters) < 1 || node.Parameters[0].(*ast.Identifier).Value != "this" {
				p.RegisterErrorWithToken("data definitions methods must have a `this` parameter", node.GetToken())
			}

			if slices.Contains(declarations, node.Name) {
				p.RegisterErrorWithToken(fmt.Sprintf("data definition with duplicated declaration '%s'", node.Name), node.GetToken())
			}

			declarations = append(declarations, node.Name)

		case *ast.Assignment:
			if id, ok := node.Left.(*ast.Identifier); !ok {
				p.RegisterErrorWithToken("data definitions can only have assignments to identifiers", node.GetToken())
			} else {
				if slices.Contains(declarations, id.Value) {
					p.RegisterErrorWithToken(fmt.Sprintf("data definition with duplicated declaration '%s'", id.Value), node.GetToken())
				}

				declarations = append(declarations, id.Value)
			}

			if node.Operator != "=" {
				p.RegisterErrorWithToken("data definitions can only have assignments with the `=` operator", node.GetToken())
			}

		}

	}

	return true
}

func escapeError(literal string) string {
	return strings.ReplaceAll(literal, "\n", "\\n")
}
