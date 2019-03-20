package lexer

import (
	"github.com/radlinskii/interpreter/token"
)

// Lexer is a struct representing the lexical analyzer.
type Lexer struct {
	input        string
	position     int
	nextPosition int
	ch           byte
	RowNum       int
}

// New creates new instance of the Lexer.
func New(input string) *Lexer {
	l := &Lexer{input: input, RowNum: 1}
	l.readChar()
	return l
}

// Reads next char from the input.
// Increments values of position and nextPositon and advances the current character.
func (l *Lexer) readChar() {
	if l.nextPosition >= len(l.input) {
		l.ch = 0
	} else {
		l.ch = l.input[l.nextPosition]
	}
	l.position = l.nextPosition
	l.nextPosition++
}

// Returns next character from the input.
func (l *Lexer) peekChar() byte {
	if l.nextPosition >= len(l.input) {
		return 0
	}
	return l.input[l.nextPosition]
}

func (l *Lexer) skipWhitespace() {
	for l.ch == ' ' || l.ch == '\t' || l.ch == '\n' || l.ch == '\r' {
		if l.ch == '\n' || l.ch == '\r' {
			l.RowNum++
		}
		l.readChar()
	}
}

func (l *Lexer) skipOneLineComment() {
	for l.ch != '\n' && l.ch != '\r' && l.ch != 0 {
		l.readChar()
	}
}

func (l *Lexer) skipMultipleLineComment() {
	// skipping '/*'
	l.readChar()
	l.readChar()

	for l.ch != 0 {
		if l.ch == '*' {
			if l.peekChar() == '/' {
				l.readChar()
				l.readChar()
				return
			}
		}

		if l.ch == '\n' || l.ch == '\r' {
			l.RowNum++
		}
		l.readChar()
	}
	// todo: error if comment is not finished
}

// NextToken analyzes text and returns the first token it founds.
func (l *Lexer) NextToken() (tok token.Token) {
	l.skipWhitespace()

	switch l.ch {
	case '=':
		if l.peekChar() == '=' {
			l.readChar()
			tok = token.Token{Type: token.EQ, Literal: "=="}
		} else {
			tok = newToken(token.ASSIGN, l.ch, l.RowNum)
		}
	case '+':
		tok = newToken(token.PLUS, l.ch, l.RowNum)
	case '-':
		tok = newToken(token.MINUS, l.ch, l.RowNum)
	case '!':
		if l.peekChar() == '=' {
			l.readChar()
			tok = token.Token{Type: token.NEQ, Literal: "!="}
		} else {
			tok = newToken(token.BANG, l.ch, l.RowNum)
		}
	case '*':
		tok = newToken(token.ASTERISK, l.ch, l.RowNum)
	case '/':
		if l.peekChar() == '/' {
			l.skipOneLineComment()
			return l.NextToken()
		} else if l.peekChar() == '*' {
			l.skipMultipleLineComment()
			return l.NextToken()
		}
		tok = newToken(token.SLASH, l.ch, l.RowNum)
	case '<':
		if l.peekChar() == '=' {
			l.readChar()
			tok = token.Token{Type: token.LTE, Literal: "<="}
		} else {
			tok = newToken(token.LT, l.ch, l.RowNum)
		}
	case '>':
		if l.peekChar() == '=' {
			l.readChar()
			tok = token.Token{Type: token.GTE, Literal: ">="}
		} else {
			tok = newToken(token.GT, l.ch, l.RowNum)
		}
	case ',':
		tok = newToken(token.COMMA, l.ch, l.RowNum)
	case ';':
		tok = newToken(token.SEMICOLON, l.ch, l.RowNum)
	case '(':
		tok = newToken(token.LPAREN, l.ch, l.RowNum)
	case ')':
		tok = newToken(token.RPAREN, l.ch, l.RowNum)
	case '{':
		tok = newToken(token.LBRACE, l.ch, l.RowNum)
	case '}':
		tok = newToken(token.RBRACE, l.ch, l.RowNum)
	case '[':
		tok = newToken(token.LBRACKET, l.ch, l.RowNum)
	case ']':
		tok = newToken(token.RBRACKET, l.ch, l.RowNum)
	case ':':
		tok = newToken(token.COLON, l.ch, l.RowNum)
	case '"':
		tok.Type = token.STRING
		tok.Literal = l.readString()
	case 0:
		tok.Literal = ""
		tok.Type = token.EOF
	default:
		if isLetter(l.ch) {
			tok.Literal = l.readIdent()
			// check if the read identifier is a keyword
			tok.Type = token.LookUpIdent(tok.Literal)
			return tok
		} else if isDigit(l.ch) {
			tok.Type = token.INT
			tok.Literal = l.readNumber()
			return tok
		} else {
			tok = newToken(token.ILLEGAL, l.ch, l.RowNum)
		}
	}
	l.readChar()
	return tok
}

// Keep reading input as long as it's a word.
func (l *Lexer) readIdent() string {
	position := l.position
	for isLetter(l.ch) {
		l.readChar()
	}
	return l.input[position:l.position]
}

// Keep reading as long as the input's a number.
func (l *Lexer) readNumber() string {
	position := l.position
	for isDigit(l.ch) {
		l.readChar()
	}
	return l.input[position:l.position]
}

func (l *Lexer) readString() string {
	position := l.position + 1
	for {
		l.readChar()
		if l.ch == '"' {
			break
		}
	}

	return l.input[position:l.position]
}

func isLetter(ch byte) bool {
	return 'a' <= ch && ch <= 'z' || 'A' <= ch && ch <= 'Z' || ch == '_'
}

func isDigit(ch byte) bool {
	return '0' <= ch && ch <= '9'
}

// create new token with given values
func newToken(tokenType token.Type, ch byte, lineNum int) token.Token {
	return token.Token{Type: tokenType, Literal: string(ch), LineNumber: lineNum}
}
