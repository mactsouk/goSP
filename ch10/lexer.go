package main

import (
	"fmt"
)

// --- Token definitions ---
type TokenType string

const (
	// Special tokens
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

	TokenLT TokenType = "<"
	TokenGT TokenType = ">"

	// Delimiters
	TokenSemicolon TokenType = ";"
	TokenLParen    TokenType = "("
	TokenRParen    TokenType = ")"
	TokenLBrace    TokenType = "{"
	TokenRBrace    TokenType = "}"
)

// --- Token ---
type Token struct {
	Type    TokenType
	Literal string
}

// --- Lexer ---
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

// --- Main lexing logic ---
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
			return Token{Type: lookupIdent(literal), Literal: literal}
		} else if isDigit(l.ch) {
			return Token{Type: TokenNumber, Literal: l.readNumber()}
		} else {
			tok = Token{Type: TokenIllegal, Literal: string(l.ch)}
		}
	}

	l.readChar()
	return tok
}

// --- Helpers ---
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
			l.readChar() // skip escape char
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

// --- Demo ---
func main() {
	input := `
    // Sample program with strings and operators
    let x = 10;
    let msg = "Hello \"world\"!";
    if x >= 10 {
        return x + 5;
    }
    if x != 5 {
        return msg;
    }
    `
	lexer := NewLexer(input)

	for tok := lexer.NextToken(); tok.Type != TokenEOF; tok = lexer.NextToken() {
		fmt.Printf("%+v\n", tok)
	}
}

