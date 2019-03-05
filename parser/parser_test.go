package parser

import (
	"testing"

	"../ast"
	"../lexer"
)

func TestVarStatements(t *testing.T) {
	input := `
	var x = 5;
	var y = 10;
	var foo = 9999;
	`

	l := lexer.New(input)
	p := New(l)

	program := p.ParseProgram()
	checkParserErrors(t, p)

	if program == nil {
		t.Fatalf("ParseProgram returned nil")
	}
	if len(program.Statements) != 3 {
		t.Fatalf("program.Statements doesn't contain 3 statements. got = %d", len(program.Statements))
	}

	tests := []struct {
		expectedIdentifier string
	}{
		{"x"},
		{"y"},
		{"foo"},
	}

	for i, tt := range tests {
		stmnt := program.Statements[i]
		if !testVarStatement(t, stmnt, tt.expectedIdentifier) {
			return
		}
	}
}

func testVarStatement(t *testing.T, s ast.Statement, name string) bool {
	t.Log(name)
	if s.TokenLiteral() != "var" {
		t.Errorf("s.TokenLiteral not 'var'. got='%q'", s.TokenLiteral())
		return false
	}

	varStmnt, ok := s.(*ast.VarStatement)
	if !ok {
		t.Errorf("varStmnt not *ast.VarStatement. got=%T", s)
		return false
	}

	if varStmnt.Name.Value != name {
		t.Errorf("varStmnt.Name.Value not '%s'. got=%s", name, varStmnt.Name.Value)
		return false
	}

	if varStmnt.Name.TokenLiteral() != name {
		t.Errorf("s.Name not %s. got=%s", name, varStmnt.Name)
		return false
	}

	return true
}

func checkParserErrors(t *testing.T, p *Parser) {
	errors := p.Errors()
	if len(errors) == 0 {
		return
	}

	t.Errorf("parser encountered %d errors", len(errors))
	for _, msg := range errors {
		t.Errorf("parser error: %q", msg)
	}

	t.FailNow()
}
