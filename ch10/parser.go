package main

import (
	"fmt"
	"io/ioutil"
	"os"
	"strconv"
	"strings"
)

// --------------------------- TOKENS & LEXER ---------------------------
type TokenType string

type Token struct {
	Type    TokenType
	Literal string
}

const (
	// special
	TokenIllegal TokenType = "ILLEGAL"
	TokenEOF     TokenType = "EOF"
	// literals
	TokenIdent  TokenType = "IDENT"
	TokenNumber TokenType = "NUMBER"
	TokenString TokenType = "STRING"
	// keywords
	TokenLet    TokenType = "LET"
	TokenIf     TokenType = "IF"
	TokenElse   TokenType = "ELSE"
	TokenReturn TokenType = "RETURN"
	// operators
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
	// delimiters
	TokenSemicolon TokenType = ";"
	TokenLParen    TokenType = "("
	TokenRParen    TokenType = ")"
	TokenLBrace    TokenType = "{"
	TokenRBrace    TokenType = "}"
)

type Lexer struct {
	input        string
	position     int // current index in input (points to current char)
	readPosition int // next index in input (points to next char)
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
			// skip comment
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
		// read string literal (without surrounding quotes in Literal)
		tok.Type = TokenString
		tok.Literal = l.readString()
	case 0:
		tok = Token{Type: TokenEOF, Literal: ""}
	default:
		if isLetter(l.ch) {
			lit := l.readIdentifier()
			tok = Token{Type: lookupIdent(lit), Literal: lit}
			return tok
		} else if isDigit(l.ch) {
			lit := l.readNumber()
			tok = Token{Type: TokenNumber, Literal: lit}
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
		// handle simple escape for \" by skipping the slash
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
	return ('a' <= ch && ch <= 'z') || ('A' <= ch && ch <= 'Z') || ch == '_'
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

// --------------------------- AST ---------------------------
type Node interface {
	TokenLiteral() string
	String() string
}

type Statement interface {
	Node
	statementNode()
}

type Expression interface {
	Node
	expressionNode()
}

type Program struct {
	Statements []Statement
}

func (p *Program) TokenLiteral() string {
	if len(p.Statements) > 0 {
		return p.Statements[0].TokenLiteral()
	}
	return ""
}
func (p *Program) String() string {
	var b strings.Builder
	for _, s := range p.Statements {
		b.WriteString(s.String())
		b.WriteString("\n")
	}
	return b.String()
}

// Statements
type LetStatement struct {
	Token Token // 'let'
	Name  *Identifier
	Value Expression
}

func (ls *LetStatement) statementNode()       {}
func (ls *LetStatement) TokenLiteral() string { return ls.Token.Literal }
func (ls *LetStatement) String() string {
	if ls.Value == nil {
		return fmt.Sprintf("let %s = ;", ls.Name.String())
	}
	return fmt.Sprintf("let %s = %s;", ls.Name.String(), ls.Value.String())
}

type ReturnStatement struct {
	Token       Token // 'return'
	ReturnValue Expression
}

func (rs *ReturnStatement) statementNode()       {}
func (rs *ReturnStatement) TokenLiteral() string { return rs.Token.Literal }
func (rs *ReturnStatement) String() string {
	if rs.ReturnValue == nil {
		return "return ;"
	}
	return "return " + rs.ReturnValue.String() + ";"
}

type ExpressionStatement struct {
	Token      Token // first token of expression
	Expression Expression
}

func (es *ExpressionStatement) statementNode()       {}
func (es *ExpressionStatement) TokenLiteral() string { return es.Token.Literal }
func (es *ExpressionStatement) String() string {
	if es.Expression == nil {
		return ""
	}
	return es.Expression.String() + ";"
}

// Expressions
type Identifier struct {
	Token Token
	Value string
}

func (i *Identifier) expressionNode()      {}
func (i *Identifier) TokenLiteral() string { return i.Token.Literal }
func (i *Identifier) String() string       { return i.Value }

type IntegerLiteral struct {
	Token Token
	Value int64
}

func (il *IntegerLiteral) expressionNode()      {}
func (il *IntegerLiteral) TokenLiteral() string { return il.Token.Literal }
func (il *IntegerLiteral) String() string       { return il.Token.Literal }

type StringLiteral struct {
	Token Token
	Value string
}

func (sl *StringLiteral) expressionNode()      {}
func (sl *StringLiteral) TokenLiteral() string { return sl.Token.Literal }
func (sl *StringLiteral) String() string       { return sl.Value }

type PrefixExpression struct {
	Token    Token // operator token, e.g. '-'
	Operator string
	Right    Expression
}

func (pe *PrefixExpression) expressionNode()      {}
func (pe *PrefixExpression) TokenLiteral() string { return pe.Token.Literal }
func (pe *PrefixExpression) String() string {
	return "(" + pe.Operator + pe.Right.String() + ")"
}

type InfixExpression struct {
	Token    Token // operator token
	Left     Expression
	Operator string
	Right    Expression
}

func (ie *InfixExpression) expressionNode()      {}
func (ie *InfixExpression) TokenLiteral() string { return ie.Token.Literal }
func (ie *InfixExpression) String() string {
	return "(" + ie.Left.String() + " " + ie.Operator + " " + ie.Right.String() + ")"
}

type IfExpression struct {
	Token       Token // 'if'
	Condition   Expression
	Consequence *BlockStatement
	Alternative *BlockStatement
}

func (ie *IfExpression) expressionNode()      {}
func (ie *IfExpression) TokenLiteral() string { return ie.Token.Literal }
func (ie *IfExpression) String() string {
	out := "if "
	if ie.Condition != nil {
		out += ie.Condition.String() + " "
	}
	if ie.Consequence != nil {
		out += ie.Consequence.String()
	}
	if ie.Alternative != nil {
		out += " else " + ie.Alternative.String()
	}
	return out
}

type BlockStatement struct {
	Token      Token // '{'
	Statements []Statement
}

func (bs *BlockStatement) statementNode()       {}
func (bs *BlockStatement) TokenLiteral() string { return bs.Token.Literal }
func (bs *BlockStatement) String() string {
	var b strings.Builder
	b.WriteString("{ ")
	for _, s := range bs.Statements {
		b.WriteString(s.String())
		b.WriteString(" ")
	}
	b.WriteString("}")
	return b.String()
}

// --------------------------- PARSER (Pratt) ---------------------------
type (
	prefixParseFn func() Expression
	infixParseFn  func(Expression) Expression
)

type Parser struct {
	l *Lexer

	curToken  Token
	peekToken Token

	errors []string

	prefixParseFns map[TokenType]prefixParseFn
	infixParseFns  map[TokenType]infixParseFn
}

const (
	_ int = iota
	LOWEST
	EQUALS      // == or !=
	LESSGREATER // < > <= >=
	SUM         // + -
	PRODUCT     // * /
	PREFIX      // -X
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

func NewParser(l *Lexer) *Parser {
	p := &Parser{
		l:      l,
		errors: []string{},
	}

	// Initialize maps
	p.prefixParseFns = make(map[TokenType]prefixParseFn)
	p.infixParseFns = make(map[TokenType]infixParseFn)

	// register prefix parse fns
	p.registerPrefix(TokenIdent, p.parseIdentifier)
	p.registerPrefix(TokenNumber, p.parseIntegerLiteral)
	p.registerPrefix(TokenString, p.parseStringLiteral)
	p.registerPrefix(TokenMinus, p.parsePrefixExpression)
	p.registerPrefix(TokenLParen, p.parseGroupedExpression)
	p.registerPrefix(TokenIf, p.parseIfExpression)

	// register infix parse fns
	for _, t := range []TokenType{
		TokenPlus, TokenMinus, TokenSlash, TokenAsterisk,
		TokenEqual, TokenNotEqual, TokenLT, TokenGT, TokenLE, TokenGE,
	} {
		p.registerInfix(t, p.parseInfixExpression)
	}

	// read two tokens, so curToken and peekToken are set
	p.nextToken()
	p.nextToken()

	return p
}

// registration helpers
func (p *Parser) registerPrefix(t TokenType, fn prefixParseFn) {
	p.prefixParseFns[t] = fn
}
func (p *Parser) registerInfix(t TokenType, fn infixParseFn) {
	p.infixParseFns[t] = fn
}

func (p *Parser) Errors() []string { return p.errors }

func (p *Parser) nextToken() {
	p.curToken = p.peekToken
	p.peekToken = p.l.NextToken()
}

func (p *Parser) curTokenIs(t TokenType) bool {
	return p.curToken.Type == t
}
func (p *Parser) peekTokenIs(t TokenType) bool {
	return p.peekToken.Type == t
}
func (p *Parser) expectPeek(t TokenType) bool {
	if p.peekTokenIs(t) {
		p.nextToken()
		return true
	}
	p.peekError(t)
	return false
}

func (p *Parser) peekError(t TokenType) {
	msg := fmt.Sprintf("expected next token to be %s, got %s instead", t, p.peekToken.Type)
	p.errors = append(p.errors, msg)
}

// Parsing entry point
func (p *Parser) ParseProgram() *Program {
	program := &Program{}
	program.Statements = []Statement{}

	for p.curToken.Type != TokenEOF {
		stmt := p.parseStatement()
		if stmt != nil {
			program.Statements = append(program.Statements, stmt)
		}
		p.nextToken()
	}
	return program
}

func (p *Parser) parseStatement() Statement {
	switch p.curToken.Type {
	case TokenLet:
		return p.parseLetStatement()
	case TokenReturn:
		return p.parseReturnStatement()
	default:
		return p.parseExpressionStatement()
	}
}

func (p *Parser) parseLetStatement() *LetStatement {
	stmt := &LetStatement{Token: p.curToken}

	if !p.expectPeek(TokenIdent) {
		return nil
	}

	stmt.Name = &Identifier{Token: p.curToken, Value: p.curToken.Literal}

	if !p.expectPeek(TokenAssign) {
		return nil
	}

	// move to expression
	p.nextToken()
	stmt.Value = p.parseExpression(LOWEST)

	// optional semicolon
	if p.peekTokenIs(TokenSemicolon) {
		p.nextToken()
	}
	return stmt
}

func (p *Parser) parseReturnStatement() *ReturnStatement {
	stmt := &ReturnStatement{Token: p.curToken}
	p.nextToken()
	stmt.ReturnValue = p.parseExpression(LOWEST)
	if p.peekTokenIs(TokenSemicolon) {
		p.nextToken()
	}
	return stmt
}

func (p *Parser) parseExpressionStatement() *ExpressionStatement {
	stmt := &ExpressionStatement{Token: p.curToken}
	stmt.Expression = p.parseExpression(LOWEST)
	if p.peekTokenIs(TokenSemicolon) {
		p.nextToken()
	}
	return stmt
}

// Pratt parse core
func (p *Parser) parseExpression(precedence int) Expression {
	prefix := p.prefixParseFns[p.curToken.Type]
	if prefix == nil {
		p.errors = append(p.errors, "no prefix parse function for "+string(p.curToken.Type))
		return nil
	}
	leftExp := prefix()

	for !p.peekTokenIs(TokenSemicolon) && precedence < p.peekPrecedence() {
		infix := p.infixParseFns[p.peekToken.Type]
		if infix == nil {
			return leftExp
		}
		p.nextToken()
		leftExp = infix(leftExp)
	}
	return leftExp
}

// prefix parse functions
func (p *Parser) parseIdentifier() Expression {
	return &Identifier{Token: p.curToken, Value: p.curToken.Literal}
}

func (p *Parser) parseIntegerLiteral() Expression {
	lit := &IntegerLiteral{Token: p.curToken}
	val, err := strconv.ParseInt(p.curToken.Literal, 10, 64)
	if err != nil {
		p.errors = append(p.errors, "could not parse integer "+p.curToken.Literal)
		return nil
	}
	lit.Value = val
	return lit
}

func (p *Parser) parseStringLiteral() Expression {
	return &StringLiteral{Token: p.curToken, Value: p.curToken.Literal}
}

func (p *Parser) parsePrefixExpression() Expression {
	expr := &PrefixExpression{Token: p.curToken, Operator: p.curToken.Literal}
	p.nextToken()
	expr.Right = p.parseExpression(PREFIX)
	return expr
}

func (p *Parser) parseGroupedExpression() Expression {
	// current token is '('
	p.nextToken()
	exp := p.parseExpression(LOWEST)
	if !p.expectPeek(TokenRParen) {
		return nil
	}
	return exp
}

func (p *Parser) parseIfExpression() Expression {
	expr := &IfExpression{Token: p.curToken}

	// support both styles: if (cond) { ... }  or  if cond { ... }
	if p.peekTokenIs(TokenLParen) {
		p.nextToken() // advance to '('
		p.nextToken() // advance to first token of condition
		expr.Condition = p.parseExpression(LOWEST)
		if !p.expectPeek(TokenRParen) {
			return nil
		}
	} else {
		p.nextToken()
		expr.Condition = p.parseExpression(LOWEST)
	}

	// expect '{'
	if !p.expectPeek(TokenLBrace) {
		return nil
	}
	expr.Consequence = p.parseBlockStatement()

	// optional else
	if p.peekTokenIs(TokenElse) {
		p.nextToken() // move to 'else'
		// allow 'else { ... }' only
		if !p.expectPeek(TokenLBrace) {
			return nil
		}
		expr.Alternative = p.parseBlockStatement()
	}

	return expr
}

func (p *Parser) parseBlockStatement() *BlockStatement {
	block := &BlockStatement{Token: p.curToken}
	block.Statements = []Statement{}

	// move into first token of block
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

// infix parse functions
func (p *Parser) parseInfixExpression(left Expression) Expression {
	expr := &InfixExpression{
		Token:    p.curToken,
		Left:     left,
		Operator: p.curToken.Literal,
	}
	prec := p.curPrecedence()
	p.nextToken()
	expr.Right = p.parseExpression(prec)
	return expr
}

// precedence helpers
func (p *Parser) peekPrecedence() int {
	if prec, ok := precedences[p.peekToken.Type]; ok {
		return prec
	}
	return LOWEST
}

func (p *Parser) curPrecedence() int {
	if prec, ok := precedences[p.curToken.Type]; ok {
		return prec
	}
	return LOWEST
}

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
		if n.Value != nil {
			printNode(n.Value, indent+1)
		}
	case *ReturnStatement:
		fmt.Println(prefix + "ReturnStatement")
		if n.ReturnValue != nil {
			printNode(n.ReturnValue, indent+1)
		}
	case *ExpressionStatement:
		fmt.Println(prefix + "ExpressionStatement")
		if n.Expression != nil {
			printNode(n.Expression, indent+1)
		}
	case *Identifier:
		fmt.Println(prefix+"Identifier:", n.Value)
	case *IntegerLiteral:
		fmt.Println(prefix+"IntegerLiteral:", n.Value)
	case *StringLiteral:
		fmt.Println(prefix+"StringLiteral:", n.Value)
	case *PrefixExpression:
		fmt.Println(prefix+"PrefixExpression:", n.Operator)
		printNode(n.Right, indent+1)
	case *InfixExpression:
		fmt.Println(prefix+"InfixExpression:", n.Operator)
		printNode(n.Left, indent+1)
		printNode(n.Right, indent+1)
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
	var src string

	if len(os.Args) > 1 {
		filename := os.Args[1]
		b, err := ioutil.ReadFile(filename)
		if err != nil {
			fmt.Println("Error reading file:", err)
			return
		}
		src = string(b)
	} else {
		src = `
        // demo program
        let x = 10;
        let msg = "Hello \"world\"!";
        if (x >= 10) {
            return x + 5;
        }
        if (x != 5) {
            return msg;
        }
        `
		fmt.Println("No filename provided — using embedded demo program.")
	}

	lex := NewLexer(src)
	parser := NewParser(lex)
	program := parser.ParseProgram()

	if len(parser.Errors()) != 0 {
		fmt.Println("Parser errors:")
		for _, e := range parser.Errors() {
			fmt.Println("  -", e)
		}
		return
	}

	fmt.Println("=== Pretty-printed AST ===")
	printNode(program, 0)
}
