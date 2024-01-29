package lexer

import (
	"bufio"
	"io"
	"log"
	"os"
	"unicode"
)

type Token int

const (
	EOF = iota
	ILLEGAL
	IDENT
	INT
	STRING
	SEMI // ;

	// Infix ops
	ADD // +
	SUB // -
	MUL // *
	DIV // /

	GREATER_THAN             // >
	LESS_THAN                // <
	GREATER_THAN_OR_EQUAL_TO // >=
	LESS_THAN_OR_EQUAL_TO    // <=
	NOT_EQUAL_TO             // !=

	ASSIGN    // =
	TRANSFORM // ->
	COLON     // :
	COMMA     // ,
	COMMENT   // // (comment)

	OPENING_BRACE // {
	CLOSING_BRACE // }
	OPENING_PAREN // (
	CLOSING_PAREN // )
	FULL_STOP     // .

	QUERY_VARIALBE // $ (query variable)
)

var tokens = []string{
	EOF:     "EOF",
	ILLEGAL: "ILLEGAL",
	IDENT:   "IDENT",
	INT:     "INT",
	STRING:  "STRING",
	SEMI:    ";",

	// Infix ops
	ADD: "+",
	SUB: "-",
	MUL: "*",
	DIV: "/",

	GREATER_THAN:             ">",
	LESS_THAN:                "<",
	GREATER_THAN_OR_EQUAL_TO: ">=",
	LESS_THAN_OR_EQUAL_TO:    "<=",
	NOT_EQUAL_TO:             "!=",

	ASSIGN:    "=",
	TRANSFORM: "->",
	COLON:     ":",
	COMMA:     ",",
	COMMENT:   "//",

	OPENING_BRACE: "{",
	CLOSING_BRACE: "}",
	OPENING_PAREN: "(",
	CLOSING_PAREN: ")",
	FULL_STOP:     ".",

	QUERY_VARIALBE: "$",
}

type Position struct {
	Line   int
	Column int
}

type Lexer struct {
	r    *bufio.Reader
	pos  Position
	file *os.File
}

type TokenInfo struct {
	Position Position
	Token    Token
	Literal  string
}

func (t Token) String() string {
	return tokens[t]
}

func NewLexer(path string) *Lexer {
	file, err := os.Open(path)
	if err != nil {
		log.Fatal(err)
	}

	return &Lexer{
		r:    bufio.NewReader(file),
		pos:  Position{Line: 1, Column: 0},
		file: file,
	}
}

func (l *Lexer) Close() {
	l.file.Close()
}

func (l *Lexer) ReadTokens() (Position, Token, string) {
	for {
		r, _, err := l.r.ReadRune()
		if err != nil {
			if err == io.EOF {
				return l.pos, EOF, ""
			}
			panic(err)
		}

		l.pos.Column++
		switch r {
		case '\n':
			l.resetPosition()
		case ';':
			return l.pos, SEMI, ";"
		case '"', '\'':
			return l.ReadString()
		case '+':
			return l.pos, ADD, "+"
		case '-':
			nextRune := l.Peek()
			if nextRune == '>' {
				l.Next()
				return l.pos, TRANSFORM, "->"
			}
			return l.pos, SUB, "-"
		case ':':
			return l.pos, COLON, ":"
		case '/':
			nextRune := l.Peek()
			if nextRune == '/' {
				l.Next()
				return l.ReadComment()
			}
			return l.pos, DIV, "/"
		case '{':
			return l.pos, OPENING_BRACE, "{"
		case '}':
			return l.pos, CLOSING_BRACE, "}"
		case '(':
			return l.pos, OPENING_PAREN, "("
		case ')':
			return l.pos, CLOSING_PAREN, ")"
		case '.':
			return l.pos, FULL_STOP, "."
		case ',':
			return l.pos, COMMA, ","
		case '=':
			return l.pos, ASSIGN, "="
		case '$':
			return l.pos, QUERY_VARIALBE, "$"
		case '>':
			nextRune := l.Peek()
			if nextRune == '=' {
				l.Next()
				return l.pos, GREATER_THAN_OR_EQUAL_TO, ">="
			}
			return l.pos, GREATER_THAN, ">"
		case '<':
			nextRune := l.Peek()
			if nextRune == '=' {
				l.Next()
				return l.pos, LESS_THAN_OR_EQUAL_TO, "<="
			}
			return l.pos, LESS_THAN, "<"
		case '!':
			nextRune := l.Peek()
			if nextRune == '=' {
				l.Next()
				return l.pos, NOT_EQUAL_TO, "!="
			}
		default:
			if unicode.IsSpace(r) {
				continue
			} else if unicode.IsLetter(r) {
				return l.ReadIdentifier(r)
			} else if unicode.IsDigit(r) {
				return l.ReadNumber(r)
			}
			return l.pos, STRING, string(r)
		}
	}
}

func (l *Lexer) Next() rune {
	r, _, err := l.r.ReadRune()
	if err != nil {
		return EOF
	}
	l.pos.Column++
	return r
}

func (l *Lexer) Peek() rune {
	r, _, _ := l.r.ReadRune()
	l.r.UnreadRune()
	return r
}

func (l *Lexer) ReadComment() (Position, Token, string) {
	rawString := ""
	var newPos Position
	for {
		r, _, err := l.r.ReadRune()
		if err != nil {
			if err == io.EOF {
				return l.pos, EOF, ""
			}
			panic(err)
		}
		l.pos.Column++
		if r == '\n' {
			newPos = l.pos
			l.resetPosition()
			break
		} else {
			rawString += string(r)
		}
	}
	return newPos, COMMENT, rawString
}

func (l *Lexer) ReadString() (Position, Token, string) {
	rawString := ""
	for {
		r, _, err := l.r.ReadRune()
		if err != nil {
			if err == io.EOF {
				return l.pos, EOF, ""
			}
			panic(err)
		}
		l.pos.Column++
		if r == '"' || r == '\'' {
			break
		} else {
			rawString += string(r)
		}
	}
	return l.pos, STRING, rawString
}

func (l *Lexer) ReadNumber(current rune) (Position, Token, string) {
	number := string(current)
	for {
		r, _, err := l.r.ReadRune()
		if err != nil {
			if err == io.EOF {
				return l.pos, EOF, ""
			}
			panic(err)
		}
		l.pos.Column++
		if unicode.IsDigit(r) {
			number += string(r)
		} else {
			l.r.UnreadRune()
			break
		}
	}
	return l.pos, INT, number
}

func (l *Lexer) ReadIdentifier(current rune) (Position, Token, string) {
	identifier := string(current)
	for {
		r, _, err := l.r.ReadRune()
		if err != nil {
			if err == io.EOF {
				return l.pos, EOF, ""
			}
			panic(err)
		}
		l.pos.Column++
		if validIdentifierSymbol(r) {
			identifier += string(r)
		} else {
			l.r.UnreadRune()
			break
		}
	}
	return l.pos, IDENT, identifier
}

func validIdentifierSymbol(symbol rune) bool {
	return unicode.IsLetter(symbol) || unicode.IsDigit(symbol) || symbol == '_'
}

func (l *Lexer) resetPosition() {
	l.pos.Line++
	l.pos.Column = 0
}

func (l *Lexer) AccumTokens() []TokenInfo {
	var tokens []TokenInfo
	for {
		pos, tok, lit := l.ReadTokens()
		if tok == EOF {
			break
		}
		info := TokenInfo{
			Position: pos,
			Token:    tok,
			Literal:  lit,
		}
		tokens = append(tokens, info)
	}

	return tokens
}
