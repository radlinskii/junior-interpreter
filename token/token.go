package token

type Type string

type Token struct {
	Type    Type
	Literal string
}

const (
	ILLEGAL = "ILLEGAL"

	EOF = "EOF"

	IDENT = "IDENT"

	INT = "INT"

	// OPERATORS
	ASSIGN   = "="
	PLUS     = "+"
	MINUS    = "-"
	BANG     = "!"
	ASTERISK = "*"
	SLASH    = "/"

	LT  = "<"
	GT  = ">"
	LTE = "<="
	GTE = ">="
	EQ  = "=="
	NEQ = "!="

	COMMA     = ","
	SEMICOLON = ";"

	LPAREN = "("
	RPAREN = ")"
	LBRACE = "{"
	RBRACE = "}"

	FUNCTION = "FUNCTION"
	RETURN   = "RETURN"
	VAR      = "VAR"
	TRUE     = "TRUE"
	FALSE    = "FALSE"
	IF       = "IF"
	ELSE     = "ELSE"
)

var keywords = map[string]Type{
	"fun":    FUNCTION,
	"var":    VAR,
	"return": RETURN,
	"true":   TRUE,
	"false":  FALSE,
	"if":     IF,
	"else":   ELSE,
}

func LookUpIdent(ident string) Type {
	if tok, ok := keywords[ident]; ok {
		return tok
	}
	return IDENT
}
