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
	isGroup      bool
	identLevel   int
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
	tokens := lex.AccumTokens()
	newLine := false

	for idx, token := range tokens {
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
			default:
				{
					if idx < len(tokens) {
						nextToken := tokens[idx+1]

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
				if idx < len(tokens) {
					nextToken := tokens[idx+1]

					if nextToken.Token == lexer.SEMI {
						hasSpace = false
					}
				}

				value = fmt.Sprintf("\"%s\"", lit)
			}
		default:
			{
				if util.Contains(symbols, tok) {
					if tok == lexer.COMMA || tok == lexer.CLOSING_PAREN {
						hasSpace = true
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
