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

	ASSIGN // =

	TRANSFORM // ->
	COLON     // :
	COMMENT   // // (comment)
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

	ASSIGN:    "=",
	TRANSFORM: "->",
	COLON:     ":",
	COMMENT:   "//",
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
		default:
			if unicode.IsSpace(r) {
				continue
			} else if unicode.IsLetter(r) {
				return l.ReadIdentifier(r)
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
			l.resetPosition()
			break
		} else {
			rawString += string(r)
		}
	}
	return l.pos, COMMENT, rawString
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
		if unicode.IsDigit(r) || unicode.IsLetter(r) {
			identifier += string(r)
		} else {
			l.r.UnreadRune()
			break
		}
	}
	return l.pos, IDENT, identifier
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
