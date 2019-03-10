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
	// LOWEST == 1 default precedence
	LOWEST
	// EQUALS == 2 precedence for operators [==,!=]
	EQUALS
	// LESSGREATER == 3 precedence for operators [>,<,>=,<=]
	LESSGREATER
	// SUM == 4 precedence for operators [+,"infixed" -]
	SUM
	// PRODUCT == 5 precedence for operators [*,/]
	PRODUCT
	// PREFIX == 6 precedence for operators ["prefixed" -,!]
	PREFIX
	// CALL == 7 precedence for operator (
	CALL
)

var precedences = map[token.Type]int{
	token.EQ:       EQUALS,
	token.NEQ:      EQUALS,
	token.LTE:      LESSGREATER,
	token.GTE:      LESSGREATER,
	token.LT:       LESSGREATER,
	token.GT:       LESSGREATER,
	token.PLUS:     SUM,
	token.MINUS:    SUM,
	token.SLASH:    PRODUCT,
	token.ASTERISK: PRODUCT,
	token.LPAREN:   CALL,
}

type prefixParseFunc func() ast.Expression
type infixParseFunc func(ast.Expression) ast.Expression

// Parser structure represents the semantic analyzer.
type Parser struct {
	lexer *lexer.Lexer

	curToken  token.Token
	peekToken token.Token
	errors    []string

	prefixParseFuncs map[token.Type]prefixParseFunc
	infixParseFuncs  map[token.Type]infixParseFunc
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

// New creates new Parser with given lexical analyzer object.
func New(l *lexer.Lexer) *Parser {
	p := &Parser{lexer: l, errors: []string{}}

	// read two tokens so curToken and peekToken are both set
	p.nextToken()
	p.nextToken()

	p.prefixParseFuncs = make(map[token.Type]prefixParseFunc)
	p.infixParseFuncs = make(map[token.Type]infixParseFunc)

	p.registerPrefix(token.IDENT, p.parseIdentifier)
	p.registerPrefix(token.IF, p.parseIfExpression)

	p.registerPrefix(token.INT, p.parseIntegerLiteral)
	p.registerPrefix(token.TRUE, p.parseBooleanLiteral)
	p.registerPrefix(token.FALSE, p.parseBooleanLiteral)
	p.registerPrefix(token.FUNCTION, p.parseFunctionLiteral)
	p.registerPrefix(token.LPAREN, p.parseGroupedExpression)

	p.registerPrefix(token.BANG, p.parsePrefixExpression)
	p.registerPrefix(token.MINUS, p.parsePrefixExpression)

	p.registerInfix(token.LPAREN, p.parseCallExpression)

	p.registerInfix(token.EQ, p.parseInfixExpression)
	p.registerInfix(token.NEQ, p.parseInfixExpression)
	p.registerInfix(token.LTE, p.parseInfixExpression)
	p.registerInfix(token.GTE, p.parseInfixExpression)
	p.registerInfix(token.LT, p.parseInfixExpression)
	p.registerInfix(token.GT, p.parseInfixExpression)
	p.registerInfix(token.PLUS, p.parseInfixExpression)
	p.registerInfix(token.MINUS, p.parseInfixExpression)
	p.registerInfix(token.SLASH, p.parseInfixExpression)
	p.registerInfix(token.ASTERISK, p.parseInfixExpression)

	return p
}

// ParseProgram starts the actual analysis of Lexer's program.
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

// Errors returns the errors that occurred during the semantic analysis.
func (p *Parser) Errors() []string {
	return p.errors
}

// returns Identifier AST node created from current token
func (p *Parser) parseIdentifier() ast.Expression {
	return &ast.Identifier{Token: p.curToken, Value: p.curToken.Literal}
}

// returns Statement AST node created from current and following tokens.
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

// checks if current token is of given type
func (p *Parser) curTokenIs(t token.Type) bool {
	return p.curToken.Type == t
}

// checks if next token is of given type
func (p *Parser) peekTokenIs(t token.Type) bool {
	return p.peekToken.Type == t
}

// checks if next token is of given type
// creates a peekError if next token isn't the expected one
func (p *Parser) expectPeek(t token.Type) bool {
	if p.peekTokenIs(t) {
		p.nextToken()
		return true
	}
	p.peekError(t)
	return false
}

// creates an error and adds it to the parser errors list
func (p *Parser) peekError(t token.Type) {
	msg := fmt.Sprintf("Unexpected token %s on line %d, expected %s.\n", p.peekToken.Type, p.lexer.RowNum, t)

	p.errors = append(p.errors, msg)
}

// parses production of var statement --> "var" <ident> "=" <expression> ";"
func (p *Parser) parseVarStatement() ast.Statement {
	stmnt := &ast.VarStatement{Token: p.curToken}

	if !p.expectPeek(token.IDENT) {
		return nil
	}

	stmnt.Name = &ast.Identifier{Token: p.curToken, Value: p.curToken.Literal}

	if !p.expectPeek(token.ASSIGN) {
		return nil
	}

	p.nextToken()

	stmnt.Value = p.parseExpression(LOWEST)

	if p.peekTokenIs(token.SEMICOLON) {
		p.nextToken()
	}

	return stmnt
}

// parses production of return statement --> "return" <expression> ";"
func (p *Parser) parseReturnStatement() *ast.ReturnStatement {
	stmnt := &ast.ReturnStatement{Token: p.curToken}

	p.nextToken()

	stmnt.ReturnValue = p.parseExpression(LOWEST)

	if p.peekTokenIs(token.SEMICOLON) {
		p.nextToken()
	}

	return stmnt
}

// Creates and returns ExpressionStatement from current token,
// it calls parseExpression to assign it to Expression property of the new ExpresisonStatement.
// Sets precedence to the lowest since it's the most outer expression in the whole statement.
func (p *Parser) parseExpressionStatement() *ast.ExpressionStatement {
	stmnt := &ast.ExpressionStatement{Token: p.curToken}

	stmnt.Expression = p.parseExpression(LOWEST)
	p.nextToken()

	if p.peekTokenIs(token.SEMICOLON) {
		p.nextToken()
	}

	return stmnt
}

// Checks precedence of current token,
// if not defined in the precedence map returns lowest precedence.
func (p *Parser) curPrecedence() int {
	if p, ok := precedences[p.curToken.Type]; ok {
		return p
	}
	return LOWEST
}

// Checks precedence of next token,
// if not defined in the precedence map returns lowest precedence.
func (p *Parser) peekPrecedence() int {
	if p, ok := precedences[p.peekToken.Type]; ok {
		return p
	}
	return LOWEST
}

func (p *Parser) parseGroupedExpression() ast.Expression {
	p.nextToken()

	exp := p.parseExpression(LOWEST)

	if !p.expectPeek(token.RPAREN) {
		return nil
	}

	return exp
}

func (p *Parser) parseIfExpression() ast.Expression {
	exp := &ast.IfExpression{Token: p.curToken}

	if !p.expectPeek(token.LPAREN) {
		return nil
	}

	p.nextToken()
	exp.Condition = p.parseExpression(LOWEST)

	if !p.expectPeek(token.RPAREN) {
		return nil
	}

	if !p.expectPeek(token.LBRACE) {
		return nil
	}

	exp.Consequence = p.parseBlockStatement()

	if p.peekTokenIs(token.ELSE) {
		p.nextToken()

		if !p.expectPeek(token.LBRACE) {
			return nil
		}

		exp.Alternative = p.parseBlockStatement()
	}

	return exp
}

func (p *Parser) parseBlockStatement() *ast.BlockStatement {
	block := &ast.BlockStatement{Token: p.curToken}
	block.Statements = []ast.Statement{}

	p.nextToken()

	for !p.curTokenIs(token.RBRACE) {
		stmnt := p.parseStatement()
		if stmnt != nil {
			block.Statements = append(block.Statements, stmnt)
		}
		p.nextToken()
	}

	return block
}

// Together with parsePrefixExpression and parseInfixExpression, parseExpression creates recursively the AST tree.
// parseExpression starts with creating left-hand side of the each node.
// Then, it continuously tries to parse the right side expression.
func (p *Parser) parseExpression(precedence int) ast.Expression {
	prefixFunc := p.prefixParseFuncs[p.curToken.Type]
	if prefixFunc == nil {
		p.noPrefixParseFuncError(p.curToken.Type)
		return nil
	}
	leftExp := prefixFunc()

	for !p.peekTokenIs(token.SEMICOLON) && precedence < p.peekPrecedence() {
		infixFunc := p.infixParseFuncs[p.peekToken.Type]
		if infixFunc == nil {
			return leftExp
		}

		p.nextToken()

		leftExp = infixFunc(leftExp)
	}

	return leftExp
}

// Returns a error message if wrong operator was used as prefix operator. e.g. in "*5;" statement.
func (p *Parser) noPrefixParseFuncError(t token.Type) {
	msg := fmt.Sprintf("no prefix parse function for %s found", t)
	p.errors = append(p.errors, msg)
}

// Creates a PrefixExpression with current token as prefix operator
// and expression as the right side of the PrefixExpression starting from next token. --> "!<expression>"
func (p *Parser) parsePrefixExpression() ast.Expression {
	expression := &ast.PrefixExpression{
		Token:    p.curToken,
		Operator: p.curToken.Literal,
	}

	p.nextToken()

	expression.Right = p.parseExpression(PREFIX)

	return expression
}

// It's given left side expression as an argument.
// It creates InfixExpression with given expression on the left and current token as the operator.
// Then it calls parseExpression with precedence of the current operator to assign it on it's right side.
func (p *Parser) parseInfixExpression(left ast.Expression) ast.Expression {
	expression := &ast.InfixExpression{
		Token:    p.curToken,
		Operator: p.curToken.Literal,
		Left:     left,
	}

	precedence := p.curPrecedence()
	p.nextToken()
	expression.Right = p.parseExpression(precedence)

	return expression
}

// Parses integer tokens into the IntegerLiterals AST nodes.
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

// Parses boolean tokens into the BooleanLiteral AST nodes.
func (p *Parser) parseBooleanLiteral() ast.Expression {
	return &ast.BooleanLiteral{Token: p.curToken, Value: p.curTokenIs(token.TRUE)}
}

func (p *Parser) parseFunctionLiteral() ast.Expression {
	fl := &ast.FunctionLiteral{Token: p.curToken}

	if !p.expectPeek(token.LPAREN) {
		return nil
	}

	fl.Parameters = p.parseFunctionParameters()

	if !p.expectPeek(token.LBRACE) {
		return nil
	}

	fl.Body = p.parseBlockStatement()

	return fl
}

func (p *Parser) parseFunctionParameters() []*ast.Identifier {
	identifiers := []*ast.Identifier{}

	if p.peekTokenIs(token.RPAREN) {
		p.nextToken()
		return identifiers
	}

	p.nextToken()

	ident := &ast.Identifier{Token: p.curToken, Value: p.curToken.Literal}
	identifiers = append(identifiers, ident)

	for p.peekTokenIs(token.COMMA) {
		p.nextToken()
		p.nextToken()

		ident := &ast.Identifier{Token: p.curToken, Value: p.curToken.Literal}
		identifiers = append(identifiers, ident)
	}

	if !p.expectPeek(token.RPAREN) {
		return nil
	}

	return identifiers
}

func (p *Parser) parseCallExpression(function ast.Expression) ast.Expression {
	exp := &ast.CallExpression{Token: p.curToken, Function: function}

	exp.Arguments = p.parseCallArguments()

	return exp
}

func (p *Parser) parseCallArguments() []ast.Expression {
	args := []ast.Expression{}

	if p.peekTokenIs(token.RPAREN) {
		p.nextToken()
		return args
	}

	p.nextToken()
	args = append(args, p.parseExpression(LOWEST))

	for p.peekTokenIs(token.COMMA) {
		p.nextToken()
		p.nextToken()
		args = append(args, p.parseExpression(LOWEST))
	}

	if !p.expectPeek(token.RPAREN) {
		return nil
	}

	return args
}
