package formatter

import (
	"fhirLSP/src/lexer"
	"fmt"
)

type Formatter struct {
	lastPosition lexer.Position
}

func NewFormatter() *Formatter {
	return &Formatter{}
}

func (f *Formatter) Format() {
	lex := lexer.NewLexer("./src/example.map")
	defer lex.Close()
	tokens := lex.AccumTokens()
	for _, token := range tokens {
		pos, tok, lit := token.Position, token.Token, token.Literal
		fmt.Printf("%d:%d\t%s\t%s\n", pos.Line, pos.Column, tok, lit)
		// if f.lastPosition.Line != pos.Line {
		// 	fmt.Print("\n")
		// }

		// fmt.Print(lit)
		// fmt.Print(" ")
		f.lastPosition = pos
	}
}
