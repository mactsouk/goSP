package main

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"
)

// //////////////////////////////////////////////////////
// -------------------- LEXER ------------------------
// //////////////////////////////////////////////////////
type TokenType string

const (
	TokenEOF     TokenType = "EOF"
	TokenIllegal TokenType = "ILLEGAL"

	// Identifiers + literals
	TokenIdent  TokenType = "IDENT"
	TokenNumber TokenType = "NUMBER"
	TokenString TokenType = "STRING"

	// Keywords
	TokenLet    TokenType = "LET"
	TokenIf     TokenType = "IF"
	TokenElse   TokenType = "ELSE"
	TokenReturn TokenType = "RETURN"

	// Operators
	TokenAssign   TokenType = "="
	TokenPlus     TokenType = "+"
	TokenMinus    TokenType = "-"
	TokenAsterisk TokenType = "*"
	TokenSlash    TokenType = "/"
	TokenEqual    TokenType = "=="
	TokenNotEqual TokenType = "!="
	TokenLE       TokenType = "<="
	TokenGE       TokenType = ">="
	TokenLT       TokenType = "<"
	TokenGT       TokenType = ">"

	// Delimiters
	TokenSemicolon TokenType = ";"
	TokenLParen    TokenType = "("
	TokenRParen    TokenType = ")"
	TokenLBrace    TokenType = "{"
	TokenRBrace    TokenType = "}"
)

type Token struct {
	Type    TokenType
	Literal string
}

type Lexer struct {
	input        string
	position     int
	readPosition int
	ch           byte
}

func NewLexer(input string) *Lexer {
	l := &Lexer{input: input}
	l.readChar()
	return l
}

func (l *Lexer) readChar() {
	if l.readPosition >= len(l.input) {
		l.ch = 0
	} else {
		l.ch = l.input[l.readPosition]
	}
	l.position = l.readPosition
	l.readPosition++
}

func (l *Lexer) peekChar() byte {
	if l.readPosition >= len(l.input) {
		return 0
	}
	return l.input[l.readPosition]
}

func (l *Lexer) NextToken() Token {
	l.skipWhitespace()
	var tok Token

	switch l.ch {
	case '=':
		if l.peekChar() == '=' {
			ch := l.ch
			l.readChar()
			tok = Token{Type: TokenEqual, Literal: string(ch) + string(l.ch)}
		} else {
			tok = Token{Type: TokenAssign, Literal: string(l.ch)}
		}
	case '!':
		if l.peekChar() == '=' {
			ch := l.ch
			l.readChar()
			tok = Token{Type: TokenNotEqual, Literal: string(ch) + string(l.ch)}
		} else {
			tok = Token{Type: TokenIllegal, Literal: string(l.ch)}
		}
	case '<':
		if l.peekChar() == '=' {
			ch := l.ch
			l.readChar()
			tok = Token{Type: TokenLE, Literal: string(ch) + string(l.ch)}
		} else {
			tok = Token{Type: TokenLT, Literal: string(l.ch)}
		}
	case '>':
		if l.peekChar() == '=' {
			ch := l.ch
			l.readChar()
			tok = Token{Type: TokenGE, Literal: string(ch) + string(l.ch)}
		} else {
			tok = Token{Type: TokenGT, Literal: string(l.ch)}
		}
	case '+':
		tok = Token{Type: TokenPlus, Literal: string(l.ch)}
	case '-':
		tok = Token{Type: TokenMinus, Literal: string(l.ch)}
	case '*':
		tok = Token{Type: TokenAsterisk, Literal: string(l.ch)}
	case '/':
		if l.peekChar() == '/' {
			l.readChar()
			l.skipComment()
			return l.NextToken()
		} else {
			tok = Token{Type: TokenSlash, Literal: string(l.ch)}
		}
	case ';':
		tok = Token{Type: TokenSemicolon, Literal: string(l.ch)}
	case '(':
		tok = Token{Type: TokenLParen, Literal: string(l.ch)}
	case ')':
		tok = Token{Type: TokenRParen, Literal: string(l.ch)}
	case '{':
		tok = Token{Type: TokenLBrace, Literal: string(l.ch)}
	case '}':
		tok = Token{Type: TokenRBrace, Literal: string(l.ch)}
	case '"':
		tok.Type = TokenString
		tok.Literal = l.readString()
	case 0:
		tok = Token{Type: TokenEOF, Literal: ""}
	default:
		if isLetter(l.ch) {
			literal := l.readIdentifier()
			tok.Type = lookupIdent(literal)
			tok.Literal = literal
			return tok
		} else if isDigit(l.ch) {
			tok.Type = TokenNumber
			tok.Literal = l.readNumber()
			return tok
		} else {
			tok = Token{Type: TokenIllegal, Literal: string(l.ch)}
		}
	}
	l.readChar()
	return tok
}

func (l *Lexer) skipWhitespace() {
	for l.ch == ' ' || l.ch == '\t' || l.ch == '\n' || l.ch == '\r' {
		l.readChar()
	}
}

func (l *Lexer) skipComment() {
	for l.ch != '\n' && l.ch != 0 {
		l.readChar()
	}
}

func (l *Lexer) readIdentifier() string {
	start := l.position
	for isLetter(l.ch) {
		l.readChar()
	}
	return l.input[start:l.position]
}

func (l *Lexer) readNumber() string {
	start := l.position
	for isDigit(l.ch) {
		l.readChar()
	}
	return l.input[start:l.position]
}

func (l *Lexer) readString() string {
	l.readChar() // skip opening quote
	start := l.position
	for l.ch != '"' && l.ch != 0 {
		if l.ch == '\\' && l.peekChar() == '"' {
			l.readChar()
		}
		l.readChar()
	}
	s := l.input[start:l.position]
	l.readChar() // skip closing quote
	return s
}

func isLetter(ch byte) bool {
	return 'a' <= ch && ch <= 'z' || 'A' <= ch && ch <= 'Z' || ch == '_'
}

func isDigit(ch byte) bool {
	return '0' <= ch && ch <= '9'
}

func lookupIdent(ident string) TokenType {
	switch ident {
	case "let":
		return TokenLet
	case "if":
		return TokenIf
	case "else":
		return TokenElse
	case "return":
		return TokenReturn
	default:
		return TokenIdent
	}
}

// //////////////////////////////////////////////////////
// -------------------- AST --------------------------
// //////////////////////////////////////////////////////
type Node interface{}
type Statement interface{ Node }
type Expression interface{ Node }

type Program struct {
	Statements []Statement
}

type LetStatement struct {
	Name  *Identifier
	Value Expression
}

type ReturnStatement struct {
	ReturnValue Expression
}

type ExpressionStatement struct {
	Expression Expression
}

type Identifier struct {
	Value string
}

type IntegerLiteral struct {
	Value int64
}

type StringLiteral struct {
	Value string
}

type PrefixExpression struct {
	Operator string
	Right    Expression
}

type InfixExpression struct {
	Left     Expression
	Operator string
	Right    Expression
}

type IfExpression struct {
	Condition   Expression
	Consequence *BlockStatement
	Alternative *BlockStatement
}

type BlockStatement struct {
	Statements []Statement
}

// //////////////////////////////////////////////////////
// -------------------- PARSER -----------------------
// //////////////////////////////////////////////////////

// Order of operations (Precedence)
const (
	_ int = iota
	LOWEST
	EQUALS      // ==
	LESSGREATER // > or <
	SUM         // +
	PRODUCT     // *
	PREFIX      // -X or !X
	CALL        // myFunction(X)
)

var precedences = map[TokenType]int{
	TokenEqual:    EQUALS,
	TokenNotEqual: EQUALS,
	TokenLT:       LESSGREATER,
	TokenGT:       LESSGREATER,
	TokenLE:       LESSGREATER,
	TokenGE:       LESSGREATER,
	TokenPlus:     SUM,
	TokenMinus:    SUM,
	TokenSlash:    PRODUCT,
	TokenAsterisk: PRODUCT,
}

type (
	prefixParseFn func() Expression
	infixParseFn  func(Expression) Expression
)

type Parser struct {
	l       *Lexer
	curTok  Token
	peekTok Token
	errors  []string

	prefixParseFns map[TokenType]prefixParseFn
	infixParseFns  map[TokenType]infixParseFn
}

func NewParser(l *Lexer) *Parser {
	p := &Parser{
		l:      l,
		errors: []string{},
	}

	p.prefixParseFns = make(map[TokenType]prefixParseFn)
	p.registerPrefix(TokenIdent, p.parseIdentifier)
	p.registerPrefix(TokenNumber, p.parseIntegerLiteral)
	p.registerPrefix(TokenString, p.parseStringLiteral)
	p.registerPrefix(TokenMinus, p.parsePrefixExpression) // Handle -5
	p.registerPrefix(TokenLParen, p.parseGroupedExpression)
	p.registerPrefix(TokenIf, p.parseIfExpression)

	p.infixParseFns = make(map[TokenType]infixParseFn)
	p.registerInfix(TokenPlus, p.parseInfixExpression)
	p.registerInfix(TokenMinus, p.parseInfixExpression)
	p.registerInfix(TokenSlash, p.parseInfixExpression)
	p.registerInfix(TokenAsterisk, p.parseInfixExpression)
	p.registerInfix(TokenEqual, p.parseInfixExpression)
	p.registerInfix(TokenNotEqual, p.parseInfixExpression)
	p.registerInfix(TokenLT, p.parseInfixExpression)
	p.registerInfix(TokenGT, p.parseInfixExpression)
	p.registerInfix(TokenLE, p.parseInfixExpression)
	p.registerInfix(TokenGE, p.parseInfixExpression)

	// Read two tokens, so curTok and peekTok are set
	p.nextToken()
	p.nextToken()

	return p
}

func (p *Parser) registerPrefix(tokenType TokenType, fn prefixParseFn) {
	p.prefixParseFns[tokenType] = fn
}

func (p *Parser) registerInfix(tokenType TokenType, fn infixParseFn) {
	p.infixParseFns[tokenType] = fn
}

func (p *Parser) nextToken() {
	p.curTok = p.peekTok
	p.peekTok = p.l.NextToken()
}

func (p *Parser) ParseProgram() *Program {
	program := &Program{}
	program.Statements = []Statement{}

	for p.curTok.Type != TokenEOF {
		stmt := p.parseStatement()
		if stmt != nil {
			program.Statements = append(program.Statements, stmt)
		}
		p.nextToken()
	}
	return program
}

func (p *Parser) parseStatement() Statement {
	switch p.curTok.Type {
	case TokenLet:
		return p.parseLetStatement()
	case TokenReturn:
		return p.parseReturnStatement()
	default:
		return p.parseExpressionStatement()
	}
}

func (p *Parser) parseLetStatement() *LetStatement {
	stmt := &LetStatement{}
	// Expect 'let' (current)
	if !p.expectPeek(TokenIdent) {
		return nil
	}
	stmt.Name = &Identifier{Value: p.curTok.Literal}

	if !p.expectPeek(TokenAssign) {
		return nil
	}
	p.nextToken()
	stmt.Value = p.parseExpression(LOWEST)

	if p.peekTok.Type == TokenSemicolon {
		p.nextToken()
	}
	return stmt
}

func (p *Parser) parseReturnStatement() *ReturnStatement {
	stmt := &ReturnStatement{}
	p.nextToken()

	stmt.ReturnValue = p.parseExpression(LOWEST)

	if p.peekTok.Type == TokenSemicolon {
		p.nextToken()
	}
	return stmt
}

func (p *Parser) parseExpressionStatement() *ExpressionStatement {
	stmt := &ExpressionStatement{}
	stmt.Expression = p.parseExpression(LOWEST)

	if p.peekTok.Type == TokenSemicolon {
		p.nextToken()
	}
	return stmt
}

// --- CORE PRATT PARSER LOGIC ---

func (p *Parser) parseExpression(precedence int) Expression {
	prefix := p.prefixParseFns[p.curTok.Type]
	if prefix == nil {
		// Error: No prefix function found
		return nil
	}
	leftExp := prefix()

	for p.peekTok.Type != TokenSemicolon && precedence < p.peekPrecedence() {
		infix := p.infixParseFns[p.peekTok.Type]
		if infix == nil {
			return leftExp
		}
		p.nextToken()
		leftExp = infix(leftExp)
	}
	return leftExp
}

func (p *Parser) parseIdentifier() Expression {
	return &Identifier{Value: p.curTok.Literal}
}

func (p *Parser) parseIntegerLiteral() Expression {
	val, _ := strconv.ParseInt(p.curTok.Literal, 10, 64)
	return &IntegerLiteral{Value: val}
}

func (p *Parser) parseStringLiteral() Expression {
	return &StringLiteral{Value: p.curTok.Literal}
}

func (p *Parser) parsePrefixExpression() Expression {
	expression := &PrefixExpression{
		Operator: p.curTok.Literal,
	}
	p.nextToken()
	expression.Right = p.parseExpression(PREFIX)
	return expression
}

func (p *Parser) parseInfixExpression(left Expression) Expression {
	expression := &InfixExpression{
		Left:     left,
		Operator: p.curTok.Literal,
	}
	precedence := p.curPrecedence()
	p.nextToken()
	expression.Right = p.parseExpression(precedence)
	return expression
}

func (p *Parser) parseGroupedExpression() Expression {
	p.nextToken()
	exp := p.parseExpression(LOWEST)
	if !p.expectPeek(TokenRParen) {
		return nil
	}
	return exp
}

func (p *Parser) parseIfExpression() Expression {
	expression := &IfExpression{}
	// We allow `if x` or `if (x)`
	if p.peekTok.Type == TokenLParen {
		p.nextToken() // eat (
		p.nextToken() // move to start of expr
		expression.Condition = p.parseExpression(LOWEST)
		if !p.expectPeek(TokenRParen) {
			return nil
		}
	} else {
		p.nextToken()
		expression.Condition = p.parseExpression(LOWEST)
	}

	if !p.expectPeek(TokenLBrace) {
		return nil
	}
	expression.Consequence = p.parseBlockStatement()

	if p.peekTok.Type == TokenElse {
		p.nextToken()
		if !p.expectPeek(TokenLBrace) {
			return nil
		}
		expression.Alternative = p.parseBlockStatement()
	}
	return expression
}

func (p *Parser) parseBlockStatement() *BlockStatement {
	block := &BlockStatement{}
	block.Statements = []Statement{}
	p.nextToken()

	for p.curTok.Type != TokenRBrace && p.curTok.Type != TokenEOF {
		stmt := p.parseStatement()
		if stmt != nil {
			block.Statements = append(block.Statements, stmt)
		}
		p.nextToken()
	}
	return block
}

// --- Helpers ---

func (p *Parser) expectPeek(t TokenType) bool {
	if p.peekTok.Type == t {
		p.nextToken()
		return true
	}
	return false
}

func (p *Parser) peekPrecedence() int {
	if p, ok := precedences[p.peekTok.Type]; ok {
		return p
	}
	return LOWEST
}

func (p *Parser) curPrecedence() int {
	if p, ok := precedences[p.curTok.Type]; ok {
		return p
	}
	return LOWEST
}

// //////////////////////////////////////////////////////
// -------------------- EVALUATOR ---------------------
// //////////////////////////////////////////////////////
type ObjectType string

const (
	IntegerObj = "INTEGER"
	StringObj  = "STRING"
	NullObj    = "NULL"
	ReturnObj  = "RETURN"
)

type Object interface {
	Type() ObjectType
	Inspect() string
}

type Integer struct{ Value int64 }

func (i *Integer) Type() ObjectType { return IntegerObj }
func (i *Integer) Inspect() string  { return fmt.Sprintf("%d", i.Value) }

type String struct{ Value string }

func (s *String) Type() ObjectType { return StringObj }
func (s *String) Inspect() string  { return s.Value }

type Null struct{}

func (n *Null) Type() ObjectType { return NullObj }
func (n *Null) Inspect() string  { return "null" }

type ReturnValue struct{ Value Object }

func (rv *ReturnValue) Type() ObjectType { return ReturnObj }
func (rv *ReturnValue) Inspect() string  { return rv.Value.Inspect() }

type Env struct {
	store map[string]Object
}

func NewEnv() *Env {
	return &Env{store: make(map[string]Object)}
}

func (e *Env) Get(name string) (Object, bool) {
	obj, ok := e.store[name]
	return obj, ok
}

func (e *Env) Set(name string, val Object) Object {
	e.store[name] = val
	return val
}

// //////////////////////////////////////////////////////
// -------------------- EVALUATOR ---------------------
// //////////////////////////////////////////////////////

// Eval is the main entry point
func Eval(node Node, env *Env) Object {
	switch n := node.(type) {

	// Statements
	case *Program:
		return evalProgram(n, env)

	case *BlockStatement:
		return evalBlockStatement(n, env)

	case *ExpressionStatement:
		return Eval(n.Expression, env)

	case *LetStatement:
		val := Eval(n.Value, env)
		if isError(val) {
			return val
		}
		env.Set(n.Name.Value, val)
		return val

	case *ReturnStatement:
		val := Eval(n.ReturnValue, env)
		if isError(val) {
			return val
		}
		return &ReturnValue{Value: val}

	// Expressions
	case *IntegerLiteral:
		return &Integer{Value: n.Value}

	case *StringLiteral:
		return &String{Value: n.Value}

	case *Identifier:
		return evalIdentifier(n, env)

	case *PrefixExpression:
		right := Eval(n.Right, env)
		if isError(right) {
			return right
		}
		return evalPrefixExpression(n.Operator, right)

	case *InfixExpression:
		left := Eval(n.Left, env)
		if isError(left) {
			return left
		}
		right := Eval(n.Right, env)
		if isError(right) {
			return right
		}
		return evalInfixExpression(n.Operator, left, right)

	case *IfExpression:
		return evalIfExpression(n, env)
	}

	return nil
}

// --- Statement Evaluators ---

func evalProgram(program *Program, env *Env) Object {
	var result Object

	for _, stmt := range program.Statements {
		result = Eval(stmt, env)

		switch result := result.(type) {
		case *ReturnValue:
			return result.Value
		case *Null: // continue
		}
	}
	return result
}

func evalBlockStatement(block *BlockStatement, env *Env) Object {
	var result Object

	for _, stmt := range block.Statements {
		result = Eval(stmt, env)

		if result != nil {
			rt := result.Type()
			if rt == ReturnObj {
				return result
			}
		}
	}
	return result
}

// --- Expression Evaluators ---

func evalIdentifier(node *Identifier, env *Env) Object {
	if val, ok := env.Get(node.Value); ok {
		return val
	}
	return &Null{} // Defined as "not found"
}

func evalPrefixExpression(operator string, right Object) Object {
	switch operator {
	case "-":
		return evalMinusPrefixOperatorExpression(right)
	default:
		return &Null{}
	}
}

func evalMinusPrefixOperatorExpression(right Object) Object {
	if right.Type() != IntegerObj {
		return &Null{}
	}
	value := right.(*Integer).Value
	return &Integer{Value: -value}
}

func evalInfixExpression(operator string, left, right Object) Object {
	// Handle Integer operations
	if left.Type() == IntegerObj && right.Type() == IntegerObj {
		return evalIntegerInfixExpression(operator, left, right)
	}
	// Handle String concatenation
	if left.Type() == StringObj && right.Type() == StringObj {
		return evalStringInfixExpression(operator, left, right)
	}

	return &Null{}
}

func evalIntegerInfixExpression(operator string, left, right Object) Object {
	leftVal := left.(*Integer).Value
	rightVal := right.(*Integer).Value

	switch operator {
	case "+":
		return &Integer{Value: leftVal + rightVal}
	case "-":
		return &Integer{Value: leftVal - rightVal}
	case "*":
		return &Integer{Value: leftVal * rightVal}
	case "/":
		return &Integer{Value: leftVal / rightVal}
	// Comparisons return 1 for true, 0 for false
	case "<":
		return nativeBoolToInteger(leftVal < rightVal)
	case ">":
		return nativeBoolToInteger(leftVal > rightVal)
	case "<=":
		return nativeBoolToInteger(leftVal <= rightVal)
	case ">=":
		return nativeBoolToInteger(leftVal >= rightVal)
	case "==":
		return nativeBoolToInteger(leftVal == rightVal)
	case "!=":
		return nativeBoolToInteger(leftVal != rightVal)
	default:
		return &Null{}
	}
}

func evalStringInfixExpression(operator string, left, right Object) Object {
	if operator != "+" {
		return &Null{}
	}
	leftVal := left.(*String).Value
	rightVal := right.(*String).Value
	return &String{Value: leftVal + rightVal}
}

func evalIfExpression(ie *IfExpression, env *Env) Object {
	condition := Eval(ie.Condition, env)

	if isTruthy(condition) {
		return Eval(ie.Consequence, env)
	} else if ie.Alternative != nil {
		return Eval(ie.Alternative, env)
	} else {
		return &Null{}
	}
}

// --- Helpers ---

func isTruthy(obj Object) bool {
	switch obj {
	case nil:
		return false
	}

	switch i := obj.(type) {
	case *Integer:
		return i.Value != 0
	case *String:
		return i.Value != ""
	case *Null:
		return false
	default:
		return true
	}
}

func nativeBoolToInteger(input bool) *Integer {
	if input {
		return &Integer{Value: 1}
	}
	return &Integer{Value: 0}
}

func isError(obj Object) bool {
	if obj != nil {
		return obj.Type() == "ERROR" // If you add ErrorObj later
	}
	return false
}

// -------------------- PRETTY PRINT -------------------
func printNode(node Node, indent int) {
	prefix := strings.Repeat("  ", indent)
	switch n := node.(type) {
	case *Program:
		fmt.Println(prefix + "Program")
		for _, stmt := range n.Statements {
			printNode(stmt, indent+1)
		}
	case *LetStatement:
		fmt.Println(prefix+"LetStatement:", n.Name.Value)
		printNode(n.Value, indent+1)
	case *ReturnStatement:
		fmt.Println(prefix + "ReturnStatement")
		printNode(n.ReturnValue, indent+1)
	case *ExpressionStatement:
		fmt.Println(prefix + "ExpressionStatement")
		printNode(n.Expression, indent+1)
	case *Identifier:
		fmt.Println(prefix+"Identifier:", n.Value)
	case *IntegerLiteral:
		fmt.Println(prefix+"IntegerLiteral:", n.Value)
	case *StringLiteral:
		fmt.Println(prefix+"StringLiteral:", n.Value)
	case *IfExpression:
		fmt.Println(prefix + "IfExpression")
		fmt.Println(prefix + "  Condition:")
		printNode(n.Condition, indent+2)
		fmt.Println(prefix + "  Consequence:")
		printNode(n.Consequence, indent+2)
		if n.Alternative != nil {
			fmt.Println(prefix + "  Alternative:")
			printNode(n.Alternative, indent+2)
		}
	case *BlockStatement:
		fmt.Println(prefix + "BlockStatement")
		for _, stmt := range n.Statements {
			printNode(stmt, indent+1)
		}
	default:
		fmt.Println(prefix + "Unknown node type")
	}
}

func main() {
	var input string

	if len(os.Args) > 1 {
		// read from file
		data, err := os.ReadFile(os.Args[1])
		if err != nil {
			fmt.Println("Error reading file:", err)
			return
		}
		input = string(data)
	} else {
		// read from stdin
		fmt.Println("Enter your code (finish input with empty line):")
		scanner := bufio.NewScanner(os.Stdin)
		for scanner.Scan() {
			line := scanner.Text()
			if line == "" {
				break
			}
			input += line + "\n"
		}
	}

	lexer := NewLexer(input)
	parser := NewParser(lexer)
	program := parser.ParseProgram()

	// fmt.Println("\n=== Pretty-printed AST ===")
	// printNode(program, 0)
	env := NewEnv()
	result := Eval(program, env)
	fmt.Println("=== Program Result ===")
	fmt.Println(result.Inspect())
}
