package main

import (
	"fmt"
	"io"
	"os"
	"strings"
	"unicode"
)

// ==================== TOKEN DEFINITIONS ====================

type TokenType string

const (
	TokenEOF     TokenType = "EOF"
	TokenIllegal TokenType = "ILLEGAL"

	// Identifiers + Literals
	TokenIdent  TokenType = "IDENT"
	TokenNumber TokenType = "NUMBER"
	TokenString TokenType = "STRING"
	TokenBool   TokenType = "BOOLEAN"

	// Keywords
	TokenLet    TokenType = "LET"
	TokenIf     TokenType = "IF"
	TokenElse   TokenType = "ELSE"
	TokenReturn TokenType = "RETURN"
	TokenFunc   TokenType = "FUNC"
	TokenFor    TokenType = "FOR"

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
	TokenBang     TokenType = "!"
	TokenAnd      TokenType = "&&"
	TokenOr       TokenType = "||"

	// Delimiters
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
}

func (t Token) String() string {
	return fmt.Sprintf("{Type:%s Literal:%q}", t.Type, t.Literal)
}

// ==================== LEXER IMPLEMENTATION ====================

type Lexer struct {
	input        []rune // Rune slice for UTF-8 safety
	position     int    // current position in input (points to current char)
	readPosition int    // current reading position in input (after current char)
	ch           rune   // current char under examination
}

func NewLexer(input string) *Lexer {
	l := &Lexer{input: []rune(input)}
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
	l.readChar() // skip the opening quote
	var sb strings.Builder
	for l.ch != '"' && l.ch != 0 {
		if l.ch == '\\' {
			l.readChar() // skip the backslash
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
	l.readChar() // skip the closing quote
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
		// BUG FIX: Return immediately. readString() already advanced past the closing quote.
		// If we break here, the l.readChar() at the bottom will skip the next character.
		return tok
	case 0:
		tok.Literal = ""
		tok.Type = TokenEOF
		// BUG FIX: Return immediately to avoid reading past EOF.
		return tok
	default:
		if unicode.IsLetter(l.ch) {
			tok.Literal = l.readIdentifier()
			tok.Type = lookupIdent(tok.Literal)
			return tok
		} else if unicode.IsDigit(l.ch) {
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

// ==================== MAIN ====================

func main() {
	var input string

	if len(os.Args) > 1 {
		filename := os.Args[1]
		data, err := os.ReadFile(filename)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error reading file '%s': %v\n", filename, err)
			os.Exit(1)
		}
		input = string(data)
		fmt.Printf("Lexing file: %s\n", filename)
	} else {
		fmt.Println("Reading from Standard Input (Press Ctrl+D to finish)...")
		data, err := io.ReadAll(os.Stdin)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error reading from stdin: %v\n", err)
			os.Exit(1)
		}
		input = string(data)
	}

	fmt.Println("--------------------------------")

	l := NewLexer(input)

	for {
		tok := l.NextToken()
		fmt.Printf("%s\n", tok)
		if tok.Type == TokenEOF {
			break
		}
	}
}
