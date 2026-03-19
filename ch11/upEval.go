package main

import (
	"fmt"
	"io"
	"os"
	"strconv"
	"strings"
	"unicode"
)

// ==================== 1. LEXER ====================

type TokenType string

const (
	TokenEOF       TokenType = "EOF"
	TokenIllegal   TokenType = "ILLEGAL"
	TokenIdent     TokenType = "IDENT"
	TokenNumber    TokenType = "NUMBER"
	TokenString    TokenType = "STRING"
	TokenBool      TokenType = "BOOLEAN"
	TokenLet       TokenType = "LET"
	TokenIf        TokenType = "IF"
	TokenElse      TokenType = "ELSE"
	TokenReturn    TokenType = "RETURN"
	TokenFunc      TokenType = "FUNC"
	TokenFor       TokenType = "FOR"
	TokenAssign    TokenType = "="
	TokenPlus      TokenType = "+"
	TokenMinus     TokenType = "-"
	TokenAsterisk  TokenType = "*"
	TokenSlash     TokenType = "/"
	TokenEqual     TokenType = "=="
	TokenNotEqual  TokenType = "!="
	TokenLE        TokenType = "<="
	TokenGE        TokenType = ">="
	TokenLT        TokenType = "<"
	TokenGT        TokenType = ">"
	TokenBang      TokenType = "!"
	TokenAnd       TokenType = "&&"
	TokenOr        TokenType = "||"
	TokenSemicolon TokenType = ";"
	TokenLParen    TokenType = "("
	TokenRParen    TokenType = ")"
	TokenLBrace    TokenType = "{"
	TokenRBrace    TokenType = "}"
	TokenComma     TokenType = ","
)

type Token struct {
	Type    TokenType
	Literal string
	Line    int
	Col     int
}

type Lexer struct {
	input        []rune
	position     int
	readPosition int
	ch           rune
	line         int
	col          int
}

func NewLexer(input string) *Lexer {
	l := &Lexer{input: []rune(input), line: 1, col: 0}
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
	l.col++
	if l.ch == '\n' {
		l.line++
		l.col = 0
	}
}

func (l *Lexer) peekChar() rune {
	if l.readPosition >= len(l.input) {
		return 0
	}
	return l.input[l.readPosition]
}

func (l *Lexer) skipWhitespace() {
	for unicode.IsSpace(l.ch) {
		l.readChar()
	}
}

func (l *Lexer) readIdentifier() string {
	start := l.position
	for unicode.IsLetter(l.ch) || unicode.IsDigit(l.ch) || l.ch == '_' {
		l.readChar()
	}
	return string(l.input[start:l.position])
}

func (l *Lexer) readNumber() string {
	start := l.position
	for unicode.IsDigit(l.ch) {
		l.readChar()
	}
	return string(l.input[start:l.position])
}

func (l *Lexer) readString() string {
	l.readChar()
	var sb strings.Builder
	for l.ch != '"' && l.ch != 0 {
		if l.ch == '\\' {
			l.readChar()
			switch l.ch {
			case 'n':
				sb.WriteRune('\n')
			case 't':
				sb.WriteRune('\t')
			case '\\':
				sb.WriteRune('\\')
			case '"':
				sb.WriteRune('"')
			default:
				sb.WriteRune(l.ch)
			}
		} else {
			sb.WriteRune(l.ch)
		}
		l.readChar()
	}
	l.readChar()
	return sb.String()
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
	case "func":
		return TokenFunc
	case "for":
		return TokenFor
	case "true", "false":
		return TokenBool
	default:
		return TokenIdent
	}
}

func (l *Lexer) NextToken() Token {
	l.skipWhitespace()
	var tok Token
	line, col := l.line, l.col

	switch l.ch {
	case '=':
		if l.peekChar() == '=' {
			l.readChar()
			tok = Token{Type: TokenEqual, Literal: "=="}
		} else {
			tok = Token{Type: TokenAssign, Literal: "="}
		}
	case '!':
		if l.peekChar() == '=' {
			l.readChar()
			tok = Token{Type: TokenNotEqual, Literal: "!="}
		} else {
			tok = Token{Type: TokenBang, Literal: "!"}
		}
	case '<':
		if l.peekChar() == '=' {
			l.readChar()
			tok = Token{Type: TokenLE, Literal: "<="}
		} else {
			tok = Token{Type: TokenLT, Literal: "<"}
		}
	case '>':
		if l.peekChar() == '=' {
			l.readChar()
			tok = Token{Type: TokenGE, Literal: ">="}
		} else {
			tok = Token{Type: TokenGT, Literal: ">"}
		}
	case '&':
		if l.peekChar() == '&' {
			l.readChar()
			tok = Token{Type: TokenAnd, Literal: "&&"}
		} else {
			tok = Token{Type: TokenIllegal, Literal: string(l.ch)}
		}
	case '|':
		if l.peekChar() == '|' {
			l.readChar()
			tok = Token{Type: TokenOr, Literal: "||"}
		} else {
			tok = Token{Type: TokenIllegal, Literal: string(l.ch)}
		}
	case '+':
		tok = Token{Type: TokenPlus, Literal: "+"}
	case '-':
		tok = Token{Type: TokenMinus, Literal: "-"}
	case '*':
		tok = Token{Type: TokenAsterisk, Literal: "*"}
	case '/':
		if l.peekChar() == '/' {
			l.readChar()
			for l.ch != '\n' && l.ch != 0 {
				l.readChar()
			}
			return l.NextToken()
		} else {
			tok = Token{Type: TokenSlash, Literal: "/"}
		}
	case ';':
		tok = Token{Type: TokenSemicolon, Literal: ";"}
	case ',':
		tok = Token{Type: TokenComma, Literal: ","}
	case '(':
		tok = Token{Type: TokenLParen, Literal: "("}
	case ')':
		tok = Token{Type: TokenRParen, Literal: ")"}
	case '{':
		tok = Token{Type: TokenLBrace, Literal: "{"}
	case '}':
		tok = Token{Type: TokenRBrace, Literal: "}"}
	case '"':
		tok.Type = TokenString
		tok.Literal = l.readString()
		tok.Line, tok.Col = line, col
		return tok
	case 0:
		tok.Literal = ""
		tok.Type = TokenEOF
		tok.Line, tok.Col = line, col
		return tok
	default:
		if unicode.IsLetter(l.ch) {
			tok.Literal = l.readIdentifier()
			tok.Type = lookupIdent(tok.Literal)
			tok.Line, tok.Col = line, col
			return tok
		} else if unicode.IsDigit(l.ch) {
			tok.Type = TokenNumber
			tok.Literal = l.readNumber()
			tok.Line, tok.Col = line, col
			return tok
		} else {
			tok = Token{Type: TokenIllegal, Literal: string(l.ch)}
		}
	}
	l.readChar()
	tok.Line, tok.Col = line, col
	return tok
}

// ==================== 2. AST ====================

type Node interface{ String() string }
type Statement interface {
	Node
	statementNode()
}
type Expression interface {
	Node
	expressionNode()
}

type Program struct{ Statements []Statement }

func (p *Program) String() string {
	var sb strings.Builder
	for _, s := range p.Statements {
		sb.WriteString(s.String() + "\n")
	}
	return sb.String()
}

type LetStatement struct {
	Name  string
	Value Expression
}

func (ls *LetStatement) statementNode() {}
func (ls *LetStatement) String() string {
	return fmt.Sprintf("let %s = %s;", ls.Name, ls.Value.String())
}

type ReturnStatement struct{ Value Expression }

func (rs *ReturnStatement) statementNode() {}
func (rs *ReturnStatement) String() string { return fmt.Sprintf("return %s;", rs.Value.String()) }

type ExpressionStatement struct{ Expression Expression }

func (es *ExpressionStatement) statementNode() {}
func (es *ExpressionStatement) String() string { return es.Expression.String() + ";" }

type BlockStatement struct{ Statements []Statement }

func (bs *BlockStatement) statementNode() {}
func (bs *BlockStatement) String() string {
	var sb strings.Builder
	sb.WriteString("{ ")
	for _, s := range bs.Statements {
		sb.WriteString(s.String() + " ")
	}
	sb.WriteString("}")
	return sb.String()
}

type IfStatement struct {
	Condition   Expression
	Consequence *BlockStatement
	Alternative *BlockStatement
}

func (is *IfStatement) statementNode() {}
func (is *IfStatement) String() string { return "if/else" }

type ForStatement struct {
	Init              Statement
	Condition, Update Expression
	Body              *BlockStatement
}

func (fs *ForStatement) statementNode() {}
func (fs *ForStatement) String() string { return "for statement" }

type FunctionStatement struct {
	Name       string
	Parameters []string
	Body       *BlockStatement
}

func (fs *FunctionStatement) statementNode() {}
func (fs *FunctionStatement) String() string {
	return fmt.Sprintf("func %s(%s) ...", fs.Name, strings.Join(fs.Parameters, ", "))
}

type Identifier struct{ Value string }

func (i *Identifier) expressionNode() {}
func (i *Identifier) String() string  { return i.Value }

type IntegerLiteral struct{ Value int64 }

func (il *IntegerLiteral) expressionNode() {}
func (il *IntegerLiteral) String() string  { return fmt.Sprintf("%d", il.Value) }

type StringLiteral struct{ Value string }

func (sl *StringLiteral) expressionNode() {}
func (sl *StringLiteral) String() string  { return fmt.Sprintf("\"%s\"", sl.Value) }

type BooleanLiteral struct{ Value bool }

func (b *BooleanLiteral) expressionNode() {}
func (b *BooleanLiteral) String() string  { return fmt.Sprintf("%t", b.Value) }

type PrefixExpression struct {
	Operator string
	Right    Expression
}

func (pe *PrefixExpression) expressionNode() {}
func (pe *PrefixExpression) String() string {
	return fmt.Sprintf("(%s%s)", pe.Operator, pe.Right.String())
}

type InfixExpression struct {
	Left     Expression
	Operator string
	Right    Expression
}

func (ie *InfixExpression) expressionNode() {}
func (ie *InfixExpression) String() string {
	return fmt.Sprintf("(%s %s %s)", ie.Left.String(), ie.Operator, ie.Right.String())
}

type AssignExpression struct {
	Name  *Identifier
	Value Expression
}

func (ae *AssignExpression) expressionNode() {}
func (ae *AssignExpression) String() string {
	return fmt.Sprintf("%s = %s", ae.Name.String(), ae.Value.String())
}

type CallExpression struct {
	Function  Expression
	Arguments []Expression
}

func (ce *CallExpression) expressionNode() {}
func (ce *CallExpression) String() string  { return "function call" }

type FunctionLiteral struct {
	Parameters []string
	Body       *BlockStatement
}

func (fl *FunctionLiteral) expressionNode() {}
func (fl *FunctionLiteral) String() string  { return "func literal" }

// ==================== 3. PARSER ====================

const (
	_ int = iota
	LOWEST
	ASSIGN
	LOGICAL
	EQUALS
	LESSGREATER
	SUM
	PRODUCT
	PREFIX
	CALL
)

var precedences = map[TokenType]int{
	TokenAssign: ASSIGN,
	TokenOr:     LOGICAL, TokenAnd: LOGICAL,
	TokenEqual: EQUALS, TokenNotEqual: EQUALS,
	TokenLT: LESSGREATER, TokenGT: LESSGREATER, TokenLE: LESSGREATER, TokenGE: LESSGREATER,
	TokenPlus: SUM, TokenMinus: SUM,
	TokenSlash: PRODUCT, TokenAsterisk: PRODUCT,
	TokenLParen: CALL,
}

type Parser struct {
	lexer     *Lexer
	curToken  Token
	peekToken Token
	errors    []string
}

func NewParser(lexer *Lexer) *Parser {
	p := &Parser{lexer: lexer, errors: []string{}}
	p.nextToken()
	p.nextToken()
	return p
}

func (p *Parser) Errors() []string { return p.errors }

func (p *Parser) nextToken() {
	p.curToken = p.peekToken
	p.peekToken = p.lexer.NextToken()
}

func (p *Parser) curTokenIs(t TokenType) bool  { return p.curToken.Type == t }
func (p *Parser) peekTokenIs(t TokenType) bool { return p.peekToken.Type == t }
func (p *Parser) expectPeek(t TokenType) bool {
	if p.peekTokenIs(t) {
		p.nextToken()
		return true
	}
	p.errors = append(p.errors, fmt.Sprintf("Error Line %d Col %d: expected %s, got %s", p.peekToken.Line, p.peekToken.Col, t, p.peekToken.Type))
	return false
}
func (p *Parser) peekPrecedence() int {
	if p, ok := precedences[p.peekToken.Type]; ok {
		return p
	}
	return LOWEST
}
func (p *Parser) curPrecedence() int {
	if p, ok := precedences[p.curToken.Type]; ok {
		return p
	}
	return LOWEST
}

func (p *Parser) ParseProgram() *Program {
	program := &Program{Statements: []Statement{}}
	for !p.curTokenIs(TokenEOF) {
		stmt := p.parseStatement()
		if stmt != nil {
			program.Statements = append(program.Statements, stmt)
		}
		p.nextToken()
	}
	return program
}

func (p *Parser) parseStatement() Statement {
	// FIX: Handle TokenFunc ambiguity.
	// If it's "func name...", it's a Statement. If "func(...", it's an Expression.
	if p.curTokenIs(TokenFunc) {
		if p.peekTokenIs(TokenIdent) {
			if s := p.parseFunctionStatement(); s != nil {
				return s
			}
			return nil
		}
		// Fallthrough to default to parse as ExpressionStatement
	}

	switch p.curToken.Type {
	case TokenLet:
		if s := p.parseLetStatement(); s != nil {
			return s
		}
		return nil
	case TokenReturn:
		if s := p.parseReturnStatement(); s != nil {
			return s
		}
		return nil
	case TokenIf:
		if s := p.parseIfStatement(); s != nil {
			return s
		}
		return nil
	case TokenFor:
		if s := p.parseForStatement(); s != nil {
			return s
		}
		return nil
	default:
		if s := p.parseExpressionStatement(); s != nil {
			return s
		}
		return nil
	}
}

func (p *Parser) parseLetStatement() *LetStatement {
	stmt := &LetStatement{}
	if !p.expectPeek(TokenIdent) {
		return nil
	}
	stmt.Name = p.curToken.Literal
	if !p.expectPeek(TokenAssign) {
		return nil
	}
	p.nextToken()
	stmt.Value = p.parseExpression(LOWEST)
	if p.peekTokenIs(TokenSemicolon) {
		p.nextToken()
	}
	return stmt
}

func (p *Parser) parseReturnStatement() *ReturnStatement {
	stmt := &ReturnStatement{}
	p.nextToken()
	stmt.Value = p.parseExpression(LOWEST)
	if p.peekTokenIs(TokenSemicolon) {
		p.nextToken()
	}
	return stmt
}

func (p *Parser) parseExpressionStatement() *ExpressionStatement {
	stmt := &ExpressionStatement{Expression: p.parseExpression(LOWEST)}
	if p.peekTokenIs(TokenSemicolon) {
		p.nextToken()
	}
	return stmt
}

func (p *Parser) parseBlockStatement() *BlockStatement {
	block := &BlockStatement{Statements: []Statement{}}
	p.nextToken()
	for !p.curTokenIs(TokenRBrace) && !p.curTokenIs(TokenEOF) {
		stmt := p.parseStatement()
		if stmt != nil {
			block.Statements = append(block.Statements, stmt)
		}
		p.nextToken()
	}
	return block
}

func (p *Parser) parseIfStatement() *IfStatement {
	stmt := &IfStatement{}
	p.nextToken()
	stmt.Condition = p.parseExpression(LOWEST)
	if !p.expectPeek(TokenLBrace) {
		return nil
	}
	stmt.Consequence = p.parseBlockStatement()
	if p.peekTokenIs(TokenElse) {
		p.nextToken()
		if !p.expectPeek(TokenLBrace) {
			return nil
		}
		stmt.Alternative = p.parseBlockStatement()
	}
	return stmt
}

func (p *Parser) parseForStatement() *ForStatement {
	stmt := &ForStatement{}
	p.nextToken() // eat 'for'

	if p.curTokenIs(TokenLParen) {
		p.nextToken()
	}

	if !p.curTokenIs(TokenLet) {
		p.errors = append(p.errors, "for loop must start with let")
		return nil
	}
	stmt.Init = p.parseLetStatement()
	p.nextToken()
	stmt.Condition = p.parseExpression(LOWEST)
	if !p.expectPeek(TokenSemicolon) {
		return nil
	}
	p.nextToken()
	stmt.Update = p.parseExpression(LOWEST)

	if p.peekTokenIs(TokenRParen) {
		p.nextToken()
	}

	if !p.expectPeek(TokenLBrace) {
		return nil
	}
	stmt.Body = p.parseBlockStatement()
	return stmt
}

func (p *Parser) parseFunctionStatement() *FunctionStatement {
	stmt := &FunctionStatement{}
	if !p.expectPeek(TokenIdent) {
		return nil
	}
	stmt.Name = p.curToken.Literal
	if !p.expectPeek(TokenLParen) {
		return nil
	}
	stmt.Parameters = p.parseFunctionParameters()
	if !p.expectPeek(TokenLBrace) {
		return nil
	}
	stmt.Body = p.parseBlockStatement()
	return stmt
}

func (p *Parser) parseFunctionParameters() []string {
	identifiers := []string{}
	if p.peekTokenIs(TokenRParen) {
		p.nextToken()
		return identifiers
	}
	p.nextToken()
	identifiers = append(identifiers, p.curToken.Literal)
	for p.peekTokenIs(TokenComma) {
		p.nextToken()
		p.nextToken()
		identifiers = append(identifiers, p.curToken.Literal)
	}
	if !p.expectPeek(TokenRParen) {
		return nil
	}
	return identifiers
}

func (p *Parser) parseExpression(precedence int) Expression {
	var leftExp Expression
	switch p.curToken.Type {
	case TokenIdent:
		leftExp = &Identifier{Value: p.curToken.Literal}
	case TokenNumber:
		val, _ := strconv.ParseInt(p.curToken.Literal, 0, 64)
		leftExp = &IntegerLiteral{Value: val}
	case TokenString:
		leftExp = &StringLiteral{Value: p.curToken.Literal}
	case TokenBool:
		leftExp = &BooleanLiteral{Value: p.curToken.Literal == "true"}
	case TokenBang, TokenMinus:
		exp := &PrefixExpression{Operator: p.curToken.Literal}
		p.nextToken()
		exp.Right = p.parseExpression(PREFIX)
		leftExp = exp
	case TokenLParen:
		p.nextToken()
		leftExp = p.parseExpression(LOWEST)
		if !p.expectPeek(TokenRParen) {
			return nil
		}
	case TokenFunc:
		leftExp = p.parseFunctionLiteral()
	default:
		p.errors = append(p.errors, fmt.Sprintf("no prefix parse function for %s", p.curToken.Type))
		return nil
	}

	for !p.peekTokenIs(TokenSemicolon) && !p.peekTokenIs(TokenRParen) && precedence < p.peekPrecedence() {
		switch p.peekToken.Type {
		case TokenAssign:
			p.nextToken() // eat '='
			p.nextToken() // move to value
			asExp := &AssignExpression{Name: leftExp.(*Identifier)}
			asExp.Value = p.parseExpression(LOWEST)
			leftExp = asExp
		case TokenPlus, TokenMinus, TokenSlash, TokenAsterisk, TokenEqual, TokenNotEqual,
			TokenLT, TokenGT, TokenLE, TokenGE, TokenAnd, TokenOr:
			p.nextToken()
			exp := &InfixExpression{Left: leftExp, Operator: p.curToken.Literal}
			precedence := p.curPrecedence()
			p.nextToken()
			exp.Right = p.parseExpression(precedence)
			leftExp = exp
		case TokenLParen:
			p.nextToken()
			exp := &CallExpression{Function: leftExp}
			exp.Arguments = p.parseCallArguments()
			leftExp = exp
		default:
			return leftExp
		}
	}
	return leftExp
}

func (p *Parser) parseFunctionLiteral() Expression {
	lit := &FunctionLiteral{}
	if !p.expectPeek(TokenLParen) {
		return nil
	}
	lit.Parameters = p.parseFunctionParameters()
	if !p.expectPeek(TokenLBrace) {
		return nil
	}
	lit.Body = p.parseBlockStatement()
	return lit
}

func (p *Parser) parseCallArguments() []Expression {
	args := []Expression{}
	if p.peekTokenIs(TokenRParen) {
		p.nextToken()
		return args
	}
	p.nextToken()
	args = append(args, p.parseExpression(LOWEST))
	for p.peekTokenIs(TokenComma) {
		p.nextToken()
		p.nextToken()
		args = append(args, p.parseExpression(LOWEST))
	}
	if !p.expectPeek(TokenRParen) {
		return nil
	}
	return args
}

// ==================== 4. OBJECT SYSTEM & EVALUATOR ====================

type ObjectType string

const (
	INTEGER_OBJ      = "INTEGER"
	BOOLEAN_OBJ      = "BOOLEAN"
	STRING_OBJ       = "STRING"
	NULL_OBJ         = "NULL"
	RETURN_VALUE_OBJ = "RETURN_VALUE"
	ERROR_OBJ        = "ERROR"
	FUNCTION_OBJ     = "FUNCTION"
	BUILTIN_OBJ      = "BUILTIN"
)

type Object interface {
	Type() ObjectType
	Inspect() string
}
type Integer struct{ Value int64 }

func (i *Integer) Type() ObjectType { return INTEGER_OBJ }
func (i *Integer) Inspect() string  { return fmt.Sprintf("%d", i.Value) }

type Boolean struct{ Value bool }

func (b *Boolean) Type() ObjectType { return BOOLEAN_OBJ }
func (b *Boolean) Inspect() string  { return fmt.Sprintf("%t", b.Value) }

type String struct{ Value string }

func (s *String) Type() ObjectType { return STRING_OBJ }
func (s *String) Inspect() string  { return s.Value }

type Null struct{}

func (n *Null) Type() ObjectType { return NULL_OBJ }
func (n *Null) Inspect() string  { return "null" }

type ReturnValue struct{ Value Object }

func (rv *ReturnValue) Type() ObjectType { return RETURN_VALUE_OBJ }
func (rv *ReturnValue) Inspect() string  { return rv.Value.Inspect() }

type Error struct{ Message string }

func (e *Error) Type() ObjectType { return ERROR_OBJ }
func (e *Error) Inspect() string  { return "ERROR: " + e.Message }

type Function struct {
	Parameters []string
	Body       *BlockStatement
	Env        *Environment
}

func (f *Function) Type() ObjectType { return FUNCTION_OBJ }
func (f *Function) Inspect() string  { return "function" }

type BuiltinFunction func(args ...Object) Object
type Builtin struct{ Fn BuiltinFunction }

func (b *Builtin) Type() ObjectType { return BUILTIN_OBJ }
func (b *Builtin) Inspect() string  { return "builtin function" }

type Environment struct {
	store map[string]Object
	outer *Environment
}

func NewEnvironment() *Environment { return &Environment{store: make(map[string]Object)} }
func NewEnclosedEnvironment(outer *Environment) *Environment {
	return &Environment{store: make(map[string]Object), outer: outer}
}
func (e *Environment) Get(name string) (Object, bool) {
	if v, ok := e.store[name]; ok {
		return v, true
	}
	if e.outer != nil {
		return e.outer.Get(name)
	}
	return nil, false
}
func (e *Environment) Set(name string, val Object) Object { e.store[name] = val; return val }
func (e *Environment) Update(name string, val Object) (Object, bool) {
	if _, ok := e.store[name]; ok {
		e.store[name] = val
		return val, true
	}
	if e.outer != nil {
		return e.outer.Update(name, val)
	}
	return nil, false
}

var (
	TRUE  = &Boolean{Value: true}
	FALSE = &Boolean{Value: false}
	NULL  = &Null{}
)

func Eval(node Node, env *Environment) Object {
	switch node := node.(type) {
	case *Program:
		return evalProgram(node, env)
	case *BlockStatement:
		return evalBlockStatement(node, env)
	case *ExpressionStatement:
		return Eval(node.Expression, env)
	case *ReturnStatement:
		val := Eval(node.Value, env)
		if isError(val) {
			return val
		}
		return &ReturnValue{Value: val}
	case *LetStatement:
		val := Eval(node.Value, env)
		if isError(val) {
			return val
		}
		env.Set(node.Name, val)
		return NULL
	case *FunctionStatement:
		env.Set(node.Name, &Function{Parameters: node.Parameters, Body: node.Body, Env: env})
		return NULL
	case *ForStatement:
		return evalForStatement(node, env)
	case *IfStatement:
		return evalIfStatement(node, env)
	case *Identifier:
		return evalIdentifier(node, env)
	case *IntegerLiteral:
		return &Integer{Value: node.Value}
	case *BooleanLiteral:
		return nativeBoolToBoolean(node.Value)
	case *StringLiteral:
		return &String{Value: node.Value}
	case *PrefixExpression:
		right := Eval(node.Right, env)
		if isError(right) {
			return right
		}
		return evalPrefixExpression(node.Operator, right)
	case *InfixExpression:
		left := Eval(node.Left, env)
		if isError(left) {
			return left
		}
		right := Eval(node.Right, env)
		if isError(right) {
			return right
		}
		return evalInfixExpression(node.Operator, left, right)
	case *AssignExpression:
		val := Eval(node.Value, env)
		if isError(val) {
			return val
		}
		if _, ok := env.Update(node.Name.Value, val); !ok {
			return newError("identifier not found for assignment: %s", node.Name.Value)
		}
		return val
	case *FunctionLiteral:
		return &Function{Parameters: node.Parameters, Body: node.Body, Env: env}
	case *CallExpression:
		fn := Eval(node.Function, env)
		if isError(fn) {
			return fn
		}
		args := evalExpressions(node.Arguments, env)
		if len(args) == 1 && isError(args[0]) {
			return args[0]
		}
		return applyFunction(fn, args)
	}
	return nil
}

func evalProgram(program *Program, env *Environment) Object {
	var result Object
	for _, stmt := range program.Statements {
		result = Eval(stmt, env)
		switch result := result.(type) {
		case *ReturnValue:
			return result.Value
		case *Error:
			return result
		}
	}
	return result
}

func evalBlockStatement(block *BlockStatement, env *Environment) Object {
	var result Object
	for _, stmt := range block.Statements {
		result = Eval(stmt, env)
		if result != nil {
			if result.Type() == RETURN_VALUE_OBJ || result.Type() == ERROR_OBJ {
				return result
			}
		}
	}
	return result
}

func evalIfStatement(ie *IfStatement, env *Environment) Object {
	cond := Eval(ie.Condition, env)
	if isError(cond) {
		return cond
	}
	if isTruthy(cond) {
		return Eval(ie.Consequence, env)
	} else if ie.Alternative != nil {
		return Eval(ie.Alternative, env)
	}
	return NULL
}

func evalForStatement(fs *ForStatement, env *Environment) Object {
	loopEnv := NewEnclosedEnvironment(env)
	if fs.Init != nil {
		if err := Eval(fs.Init, loopEnv); isError(err) {
			return err
		}
	}
	for {
		if fs.Condition != nil {
			cond := Eval(fs.Condition, loopEnv)
			if isError(cond) {
				return cond
			}
			if !isTruthy(cond) {
				break
			}
		}
		res := Eval(fs.Body, loopEnv)
		if res != nil && (res.Type() == RETURN_VALUE_OBJ || res.Type() == ERROR_OBJ) {
			return res
		}
		if fs.Update != nil {
			if err := Eval(fs.Update, loopEnv); isError(err) {
				return err
			}
		}
	}
	return NULL
}

func evalIdentifier(node *Identifier, env *Environment) Object {
	if val, ok := env.Get(node.Value); ok {
		return val
	}
	if builtin, ok := builtins[node.Value]; ok {
		return builtin
	}
	return newError("identifier not found: %s", node.Value)
}

func evalExpressions(exps []Expression, env *Environment) []Object {
	var result []Object
	for _, e := range exps {
		evaluated := Eval(e, env)
		if isError(evaluated) {
			return []Object{evaluated}
		}
		result = append(result, evaluated)
	}
	return result
}

func evalPrefixExpression(operator string, right Object) Object {
	switch operator {
	case "!":
		return evalBangOperatorExpression(right)
	case "-":
		return evalMinusPrefixOperatorExpression(right)
	default:
		return newError("unknown operator: %s%s", operator, right.Type())
	}
}

func evalInfixExpression(operator string, left, right Object) Object {
	switch {
	case left.Type() == INTEGER_OBJ && right.Type() == INTEGER_OBJ:
		return evalIntegerInfixExpression(operator, left, right)
	case left.Type() == STRING_OBJ && right.Type() == STRING_OBJ:
		return evalStringInfixExpression(operator, left, right)
	case operator == "==":
		return nativeBoolToBoolean(left == right)
	case operator == "!=":
		return nativeBoolToBoolean(left != right)
	case operator == "&&":
		return nativeBoolToBoolean(isTruthy(left) && isTruthy(right))
	case operator == "||":
		return nativeBoolToBoolean(isTruthy(left) || isTruthy(right))
	case left.Type() != right.Type():
		return newError("type mismatch: %s %s %s", left.Type(), operator, right.Type())
	default:
		return newError("unknown operator: %s %s %s", left.Type(), operator, right.Type())
	}
}

func evalIntegerInfixExpression(operator string, left, right Object) Object {
	lv, rv := left.(*Integer).Value, right.(*Integer).Value
	switch operator {
	case "+":
		return &Integer{Value: lv + rv}
	case "-":
		return &Integer{Value: lv - rv}
	case "*":
		return &Integer{Value: lv * rv}
	case "/":
		if rv == 0 {
			return newError("division by zero")
		}
		return &Integer{Value: lv / rv}
	case "<":
		return nativeBoolToBoolean(lv < rv)
	case ">":
		return nativeBoolToBoolean(lv > rv)
	case "<=":
		return nativeBoolToBoolean(lv <= rv)
	case ">=":
		return nativeBoolToBoolean(lv >= rv)
	case "==":
		return nativeBoolToBoolean(lv == rv)
	case "!=":
		return nativeBoolToBoolean(lv != rv)
	default:
		return newError("unknown operator: %s %s %s", left.Type(), operator, right.Type())
	}
}

func evalStringInfixExpression(operator string, left, right Object) Object {
	if operator != "+" {
		return newError("unknown operator: %s %s %s", left.Type(), operator, right.Type())
	}
	return &String{Value: left.(*String).Value + right.(*String).Value}
}

func evalBangOperatorExpression(right Object) Object {
	switch right {
	case TRUE:
		return FALSE
	case FALSE:
		return TRUE
	case NULL:
		return TRUE
	default:
		return FALSE
	}
}

func evalMinusPrefixOperatorExpression(right Object) Object {
	if right.Type() != INTEGER_OBJ {
		return newError("unknown operator: -%s", right.Type())
	}
	return &Integer{Value: -right.(*Integer).Value}
}

func applyFunction(fn Object, args []Object) Object {
	switch fn := fn.(type) {
	case *Function:
		extendedEnv := NewEnclosedEnvironment(fn.Env)
		for i, param := range fn.Parameters {
			if i < len(args) {
				extendedEnv.Set(param, args[i])
			}
		}
		evaluated := Eval(fn.Body, extendedEnv)
		if returnValue, ok := evaluated.(*ReturnValue); ok {
			return returnValue.Value
		}
		return evaluated
	case *Builtin:
		return fn.Fn(args...)
	default:
		return newError("not a function: %s", fn.Type())
	}
}

func nativeBoolToBoolean(input bool) *Boolean {
	if input {
		return TRUE
	}
	return FALSE
}
func isTruthy(obj Object) bool {
	if obj == NULL || obj == FALSE {
		return false
	}
	return true
}
func newError(format string, a ...interface{}) *Error {
	return &Error{Message: fmt.Sprintf(format, a...)}
}
func isError(obj Object) bool {
	if obj != nil {
		return obj.Type() == ERROR_OBJ
	}
	return false
}

var builtins = map[string]*Builtin{
	"print": {Fn: func(args ...Object) Object {
		for _, arg := range args {
			fmt.Print(arg.Inspect())
		}
		return NULL
	}},
	"println": {Fn: func(args ...Object) Object {
		for _, arg := range args {
			fmt.Print(arg.Inspect())
		}
		fmt.Println()
		return NULL
	}},
}

// ==================== 5. MAIN ====================

func main() {
	var input string
	if len(os.Args) > 1 {
		data, err := os.ReadFile(os.Args[1])
		if err != nil {
			fmt.Println(err)
			return
		}
		input = string(data)
	} else {
		data, err := io.ReadAll(os.Stdin)
		if err != nil {
			fmt.Println(err)
			return
		}
		input = string(data)
	}
	l := NewLexer(input)
	p := NewParser(l)
	program := p.ParseProgram()
	if len(p.Errors()) > 0 {
		fmt.Println("Parser Errors:")
		for _, msg := range p.Errors() {
			fmt.Println("  " + msg)
		}
		os.Exit(1)
	}
	env := NewEnvironment()
	Eval(program, env)
}
