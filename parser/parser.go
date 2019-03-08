package parser

import (
	"fmt"
	"strconv"

	"github.com/radlinskii/interpreter/ast"
	"github.com/radlinskii/interpreter/lexer"
	"github.com/radlinskii/interpreter/token"
)

const (
	_ int = iota
	LOWEST
	EQUALS      // ==
	LESSGREATER // > | < | <= | >=
	SUM         // + | -
	PRODUCT     // / | *
	PREFIX      // -X | !X
	CALL        // someFun(X)
)

type prefixParseFunc func() ast.Expression
type infixParseFunc func(ast.Expression) ast.Expression

type Parser struct {
	lexer *lexer.Lexer

	curToken  token.Token
	peekToken token.Token
	errors    []string

	prefixParseFuncs map[token.Type]prefixParseFunc
	infixParseFuncs  map[token.Type]infixParseFunc
}

func New(l *lexer.Lexer) *Parser {
	p := &Parser{lexer: l, errors: []string{}}

	// read two tokens so curToken and peekToken are both set
	p.nextToken()
	p.nextToken()

	p.prefixParseFuncs = make(map[token.Type]prefixParseFunc)
	p.registerPrefix(token.IDENT, p.parseIdentifier)
	p.registerPrefix(token.INT, p.parseIntegerLiteral)
	p.registerPrefix(token.BANG, p.parsePrefixExpression)
	p.registerPrefix(token.MINUS, p.parsePrefixExpression)

	return p
}

func (p *Parser) Errors() []string {
	return p.errors
}

func (p *Parser) nextToken() {
	p.curToken = p.peekToken
	p.peekToken = p.lexer.NextToken()
}

func (p *Parser) registerPrefix(tokenType token.Type, fn prefixParseFunc) {
	p.prefixParseFuncs[tokenType] = fn
}

func (p *Parser) registerInfix(tokenType token.Type, fn infixParseFunc) {
	p.infixParseFuncs[tokenType] = fn
}

func (p *Parser) ParseProgram() *ast.Program {
	program := &ast.Program{}
	program.Statements = []ast.Statement{}

	for !p.curTokenIs(token.EOF) {
		stmnt := p.parseStatement()
		if stmnt != nil {
			program.Statements = append(program.Statements, stmnt)
		}
		p.nextToken()
	}

	return program
}

func (p *Parser) parseIdentifier() ast.Expression {
	return &ast.Identifier{Token: p.curToken, Value: p.curToken.Literal}
}

func (p *Parser) parseStatement() ast.Statement {
	switch p.curToken.Type {
	case token.VAR:
		return p.parseVarStatement()
	case token.RETURN:
		return p.parseReturnStatement()
	default:
		return p.parseExpressionStatement()
	}
}

func (p *Parser) parseVarStatement() ast.Statement {
	stmnt := &ast.VarStatement{Token: p.curToken}

	if !p.expectPeek(token.IDENT) {
		return nil
	}

	stmnt.Name = &ast.Identifier{Token: p.curToken, Value: p.curToken.Literal}

	if !p.expectPeek(token.ASSIGN) {
		return nil
	}

	// TODO skipping the expression until we hit semicolon
	for !p.curTokenIs(token.SEMICOLON) {
		p.nextToken()
	}

	return stmnt
}

func (p *Parser) curTokenIs(t token.Type) bool {
	return p.curToken.Type == t
}

func (p *Parser) peekTokenIs(t token.Type) bool {
	return p.peekToken.Type == t
}

func (p *Parser) expectPeek(t token.Type) bool {
	if p.peekTokenIs(t) {
		p.nextToken()
		return true
	}
	p.peekError(t)
	return false
}

func (p *Parser) peekError(t token.Type) {
	msg := fmt.Sprintf("Unexpected token %s on line %d, expected %s.\n", p.peekToken.Type, p.lexer.RowNum, t)

	p.errors = append(p.errors, msg)
}

func (p *Parser) parseReturnStatement() *ast.ReturnStatement {
	stmnt := &ast.ReturnStatement{Token: p.curToken}

	p.nextToken()

	// TODO skipping the expression until we hit semicolon
	for !p.curTokenIs(token.SEMICOLON) {
		p.nextToken()
	}

	return stmnt
}

func (p *Parser) parseExpressionStatement() *ast.ExpressionStatement {
	stmnt := &ast.ExpressionStatement{Token: p.curToken}

	stmnt.Expression = p.parseExpression(LOWEST)
	p.nextToken()

	if p.peekTokenIs(token.SEMICOLON) {
		p.nextToken()
	}

	return stmnt
}

func (p *Parser) parseExpression(precendence int) ast.Expression {
	prefixFunc := p.prefixParseFuncs[p.curToken.Type]
	if prefixFunc == nil {
		p.noPrefixParseFuncError(p.curToken.Type)
		return nil
	}
	leftExp := prefixFunc()

	return leftExp
}

func (p *Parser) noPrefixParseFuncError(t token.Type) {
	msg := fmt.Sprintf("no prefix parse function for %s found", t)
	p.errors = append(p.errors, msg)
}

func (p *Parser) parsePrefixExpression() ast.Expression {
	expression := &ast.PrefixExpression{
		Token:    p.curToken,
		Operator: p.curToken.Literal,
	}

	p.nextToken()

	expression.Right = p.parseExpression(PREFIX)

	return expression
}

func (p *Parser) parseIntegerLiteral() ast.Expression {
	lit := &ast.IntegerLiteral{Token: p.curToken}

	value, err := strconv.ParseInt(p.curToken.Literal, 0, 64)
	if err != nil {
		msg := fmt.Sprintf("could not parse %q as integer", p.curToken.Literal)
		p.errors = append(p.errors, msg)

		return nil
	}

	lit.Value = value

	return lit
}
