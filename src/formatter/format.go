package formatter

import (
	"fhirLSP/src/lexer"
	"fhirLSP/src/util"
	"fmt"
	"os"
	"strings"
)

type Formatter struct {
	lastPosition lexer.Position
	tokenIndex   int
	isGroup      bool
	identLevel   int
	isStatement  bool
	tokens       []lexer.TokenInfo
}

func NewFormatter() *Formatter {
	return &Formatter{
		identLevel: 0,
	}
}

var symbols = []lexer.Token{
	lexer.FULL_STOP,
	lexer.SEMI,
	lexer.COLON,
	lexer.COMMA,
	lexer.TRANSFORM,
	lexer.OPENING_PAREN,
	lexer.CLOSING_PAREN,
	lexer.CLOSING_PAREN,
	lexer.QUERY_VARIALBE,
}

func (f *Formatter) Format() {
	var builder strings.Builder
	lex := lexer.NewLexer("./src/example.map")
	defer lex.Close()
	f.tokens = lex.AccumTokens()
	newLine := false

	for idx, token := range f.tokens {
		pos, tok, lit := token.Position, token.Token, token.Literal
		fmt.Printf("%d:%d\t%s\t%s\n", pos.Line, pos.Column, tok, lit)
		value := lit

		if f.lastPosition.Line != pos.Line {
			builder.WriteString("\n")
			newLine = true
		}

		hasSpace := true

		switch tok {
		case lexer.IDENT:
			switch lit {
			case lexer.IDENT_GROUP:
				builder.WriteString("\n")
				f.isGroup = true
			case lexer.IDENT_THEN:
				f.isStatement = true
			default:
				{
					if idx < len(f.tokens) {
						nextToken := f.tokens[idx+1]

						for _, symbol := range symbols {
							if nextToken.Token == symbol {
								hasSpace = false
								break
							}
						}
					}
				}
			}
		case lexer.OPENING_BRACE:
			{
				f.identLevel++
			}
		case lexer.CLOSING_BRACE:
			{
				f.identLevel--
			}
		case lexer.COMMENT:
			{
				fmt.Println(pos.Line, " ", f.lastPosition.Line)
				value = fmt.Sprintf("// %s", lit)
				hasSpace = false
			}
		case lexer.STRING:
			{
				if idx < len(f.tokens) {
					nextToken := f.tokens[idx+1]

					if nextToken.Token == lexer.SEMI {
						hasSpace = false
					}
				}

				value = fmt.Sprintf("\"%s\"", lit)
			}
		case lexer.SEMI:
			{
				hasSpace = false
				f.isStatement = false
			}
		case lexer.TRANSFORM:
			{
				hasSpace = false
				f.isStatement = true
			}
		default:
			{
				if util.Contains(symbols, tok) {
					if tok == lexer.COMMA {
						hasSpace = true
					} else if tok == lexer.CLOSING_PAREN {
						found, nextToken := f.GetNextToken(idx)
						if found {
							if nextToken.Token == lexer.COMMA {
								hasSpace = false
							}
						} else {
							hasSpace = true
						}
					} else {
						hasSpace = false
					}
				}
			}
		}

		if f.identLevel > 0 && newLine {
			for i := 0; i < f.identLevel; i++ {
				builder.WriteString("\t")
			}
			if f.isStatement {
				builder.WriteString("\t")
			}
			newLine = false
		}

		builder.WriteString(value)

		if hasSpace {
			builder.WriteString(" ")
		}
		f.lastPosition = pos
	}

	d1 := []byte(builder.String())
	err := os.WriteFile("./output/example.map", d1, 0644)
	if err != nil {
		panic(err)
	}
}

func (f *Formatter) GetNextToken(index int) (bool, lexer.TokenInfo) {
	if index < len(f.tokens) {
		return true, f.tokens[index+1]
	}

	return false, lexer.TokenInfo{}
}

func (f *Formatter) NextToken() (lexer.TokenInfo, error) {
	if f.tokenIndex < len(f.tokens) {
		f.tokenIndex++
		return f.tokens[f.tokenIndex], nil
	}

	return lexer.TokenInfo{}, fmt.Errorf("no more tokens")
}

func (f *Formatter) PreviousToken() (lexer.TokenInfo, error) {
	if f.tokenIndex > 0 {
		f.tokenIndex--
		return f.tokens[f.tokenIndex], nil
	}

	return lexer.TokenInfo{}, fmt.Errorf("no more tokens")
}

func (f *Formatter) PeekToken() (lexer.TokenInfo, error) {
	if f.tokenIndex < len(f.tokens) {
		return f.tokens[f.tokenIndex], nil
	}

	return lexer.TokenInfo{}, fmt.Errorf("no more tokens")
}