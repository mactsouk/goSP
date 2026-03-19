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
	l.readChar() // skip quote
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
	l.readChar() // skip closing quote
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
		return tok // RETURN IMMEDIATELY
	case 0:
		tok.Literal = ""
		tok.Type = TokenEOF
		tok.Line, tok.Col = line, col
		return tok // RETURN IMMEDIATELY
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
	TokenLT: LESSGREATER, TokenGT: LESSGREATER,
	TokenLE: LESSGREATER, TokenGE: LESSGREATER,
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
	switch p.curToken.Type {
	case TokenLet:
		return p.parseLetStatement()
	case TokenReturn:
		return p.parseReturnStatement()
	case TokenIf:
		return p.parseIfStatement()
	case TokenFor:
		return p.parseForStatement()
	case TokenFunc:
		return p.parseFunctionStatement()
	default:
		return p.parseExpressionStatement()
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
	p.nextToken()
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

// ==================== 4. MAIN ====================

func printAST(node Node, indent string) {
	switch n := node.(type) {
	case *Program:
		fmt.Println("Program")
		for _, s := range n.Statements {
			printAST(s, indent+"  ")
		}
	case *FunctionStatement:
		fmt.Printf("%sFunctionStatement: %s(%s)\n", indent, n.Name, strings.Join(n.Parameters, ", "))
		printAST(n.Body, indent+"  ")
	case *BlockStatement:
		fmt.Printf("%sBlockStatement\n", indent)
		for _, s := range n.Statements {
			printAST(s, indent+"  ")
		}
	case *IfStatement:
		fmt.Printf("%sIfStatement\n", indent)
		fmt.Printf("%s  Condition:\n", indent)
		printAST(n.Condition, indent+"    ")
		fmt.Printf("%s  Consequence:\n", indent)
		printAST(n.Consequence, indent+"    ")
		if n.Alternative != nil {
			fmt.Printf("%s  Alternative:\n", indent)
			printAST(n.Alternative, indent+"    ")
		}
	case *ForStatement:
		fmt.Printf("%sForStatement\n", indent)
		if n.Init != nil {
			fmt.Printf("%s  Init:\n", indent)
			printAST(n.Init, indent+"    ")
		}
		fmt.Printf("%s  Condition:\n", indent)
		printAST(n.Condition, indent+"    ")
		fmt.Printf("%s  Update:\n", indent)
		printAST(n.Update, indent+"    ")
		fmt.Printf("%s  Body:\n", indent)
		printAST(n.Body, indent+"    ")
	case *ReturnStatement:
		fmt.Printf("%sReturnStatement\n", indent)
		printAST(n.Value, indent+"  ")
	case *LetStatement:
		fmt.Printf("%sLetStatement: %s\n", indent, n.Name)
		printAST(n.Value, indent+"  ")
	case *ExpressionStatement:
		fmt.Printf("%sExpressionStatement\n", indent)
		printAST(n.Expression, indent+"  ")
	case *AssignExpression:
		fmt.Printf("%sAssignExpression\n", indent)
		fmt.Printf("%s  Name: %s\n", indent, n.Name.Value)
		fmt.Printf("%s  Value:\n", indent)
		printAST(n.Value, indent+"    ")
	case *CallExpression:
		fmt.Printf("%sCallExpression\n", indent)
		fmt.Printf("%s  Function: %s\n", indent, n.Function.String())
		fmt.Printf("%s  Arguments:\n", indent)
		for _, arg := range n.Arguments {
			printAST(arg, indent+"    ")
		}
	case *InfixExpression:
		fmt.Printf("%sInfixExpression (%s)\n", indent, n.Operator)
		fmt.Printf("%s  Left: %s\n", indent, n.Left.String())
		fmt.Printf("%s  Right: %s\n", indent, n.Right.String())
	default:
		if node != nil {
			fmt.Printf("%s%s\n", indent, node.String())
		}
	}
}

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

	fmt.Println("AST Structure:")
	printAST(program, "")
}
