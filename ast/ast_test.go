package ast

import (
	"testing"

	"github.com/radlinskii/interpreter/token"
)

func TestString(t *testing.T) {
	program := &Program{
		Statements: []Statement{
			&ConstStatement{
				Token: token.Token{Type: token.CONST, Literal: "const"},
				Name: &Identifier{
					Token: token.Token{Type: token.IDENT, Literal: "myConst"},
					Value: "myConst",
				},
				Value: &Identifier{
					Token: token.Token{Type: token.IDENT, Literal: "anotherConst"},
					Value: "anotherConst",
				},
			},
		},
	}

	if program.String() != "const myConst = anotherConst;" {
		t.Errorf("program.String() wrong. got=%q", program.String())
	}
}
