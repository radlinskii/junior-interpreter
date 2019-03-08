package token

// Type is a token type.
type Type string

// Token is the lexical symbol that gets returned after performing lexical analysis.
type Token struct {
	Type    Type
	Literal string
}

const (
	// ILLEGAL token is created when symbols not belonging to our language are found.
	ILLEGAL = "ILLEGAL"

	// EOF - end of file
	EOF = "EOF"

	// IDENT - identifier
	IDENT = "IDENT"

	// INT - integer literal
	INT = "INT"

	// ASSIGN - assign operator
	ASSIGN = "="
	// PLUS - sum / concatenation
	PLUS = "+"
	// MINUS - subtraction / negate number
	MINUS = "-"
	// BANG - negate logical expression or value
	BANG = "!"
	// ASTERISK - multiplication
	ASTERISK = "*"
	// SLASH - division
	SLASH = "/"

	// LT - lower than
	LT = "<"
	// GT - greater than
	GT = ">"
	// LTE - lower than or equal
	LTE = "<="
	// GTE - greater than or equal
	GTE = ">="
	// EQ - equal
	EQ = "=="
	// NEQ - not equal
	NEQ = "!="

	// COMMA - values delimeter
	COMMA = ","
	// SEMICOLON - separates expressions
	SEMICOLON = ";"

	// LPAREN - function calls, binding expressions
	LPAREN = "("
	// RPAREN - function calls, binding expressions
	RPAREN = ")"
	// LBRACE - starts block statement
	LBRACE = "{"
	// RBRACE = ends block statements
	RBRACE = "}"

	// FUNCTION keyword "fun"
	FUNCTION = "FUNCTION"
	// RETURN keyword "return"
	RETURN = "RETURN"
	// VAR keyword "var"
	VAR = "VAR"
	// TRUE keyword "true"
	TRUE = "TRUE"
	// FALSE keyword "false"
	FALSE = "FALSE"
	// IF keyword "if"
	IF = "IF"
	// ELSE keyword "else"
	ELSE = "ELSE"
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

// LookUpIdent checks if identifier exists in the map of keywords.
func LookUpIdent(ident string) Type {
	if tok, ok := keywords[ident]; ok {
		return tok
	}
	return IDENT
}
