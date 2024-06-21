package runevm

import (
	"fmt"
	"strings"
)

type Parser struct {
	input *TokenStream
}

func NewParser(input *TokenStream) *Parser {
	return &Parser{input: input}
}

var FALSE = &Expr{
	Type:  Bool,
	Value: false,
}

var PRECEDENCE = map[string]int{
	"=":  1,
	"||": 2,
	"&&": 3,
	"<":  7, ">": 7, "<=": 7, ">=": 7, "==": 7, "!=": 7,
	"+": 10, "-": 10,
	"*": 20, "/": 20, "%": 20,
}

func (p *Parser) isPunc(ch string) *Token {
	tok := p.input.Peek()
	if tok != nil && tok.Type == "punc" && (ch == "" || tok.Value == ch) {
		return tok
	}
	return nil
}

func (p *Parser) isKw(kw string) *Token {
	tok := p.input.Peek()
	if tok != nil && tok.Type == "kw" && (kw == "" || tok.Value == kw) {
		return tok
	}
	return nil
}

func (p *Parser) isOp(op string) *Token {
	tok := p.input.Peek()
	if tok != nil && tok.Type == "op" && (op == "" || tok.Value == op) {
		return tok
	}
	return nil
}

func (p *Parser) skipPunc(ch string) {
	if p.isPunc(ch) != nil {
		p.input.Next()
	} else {
		p.input.Error(p.input.current, fmt.Sprintf("Expecting punctuation: \"%s\"", ch))
	}
}

func (p *Parser) skipKw(kw string) {
	if p.isKw(kw) != nil {
		p.input.Next()
	} else {
		p.input.Error(p.input.current, fmt.Sprintf("Expecting keyword: \"%s\"", kw))
	}
}

func (p *Parser) unexpected(tok *Token) {
	if tok != nil {
		p.input.Error(tok, fmt.Sprintf("Unexpected token: \"%s\"", tok.Value))
	}
}

func (p *Parser) parseBinaryExpression(left *Expr, prec int) *Expr {
	for {
		tok := p.isOp("")
		if tok == nil {
			break
		}
		opPrec := PRECEDENCE[tok.Value]
		if opPrec <= prec {
			break
		}

		p.input.Next()

		right := p.parseAtom()

		nextTok := p.isOp("")

		if nextTok != nil && PRECEDENCE[nextTok.Value] > opPrec {
			right = p.parseBinaryExpression(right, opPrec)
		}

		var exprType ExprType
		if tok.Value == "=" {
			exprType = Assign
		} else {
			exprType = Binary
		}

		left = &Expr{
			Type:     exprType,
			Operator: tok.Value,
			Left:     left,
			Right:    right,
			File:     tok.File,
			Line:     tok.Line,
			Col:      tok.Col,
		}
	}
	return left
}

func (p *Parser) parseDelimited(start, stop, separator string, parser func() *Expr) []*Expr {
	var a []*Expr
	first := true
	p.skipPunc(start)
	for !p.input.Eof() {
		if p.isPunc(stop) != nil {
			break
		}
		if !first {
			p.skipPunc(separator)
		}
		if p.isPunc(stop) != nil {
			break
		}
		a = append(a, parser())
		first = false
	}
	p.skipPunc(stop)
	return a
}

func (p *Parser) parseEnclosed(start string, stop string, parser func() *Expr) []*Expr {
	var a []*Expr
	p.skipPunc(start)
	for !p.input.Eof() {
		if p.isPunc(stop) != nil {
			break
		}
		if p.isPunc(stop) != nil {
			break
		}
		a = append(a, parser())
	}
	p.skipPunc(stop)
	return a
}

func (p *Parser) parseFunctionCall(funcExpr *Expr) *Expr {
	tok := p.input.Peek()
	return &Expr{
		Type: Call,
		Func: funcExpr,
		Args: p.parseDelimited("(", ")", ",", p.parseExpression),
		File: tok.File,
		Line: tok.Line,
		Col:  tok.Col,
	}
}

func (p *Parser) parseVarname() string {
	name := p.input.Next()
	if name.Type != "var" {
		p.input.Error(name, fmt.Sprintf("Expecting variable name, but got: '%s'", name.Value))
	}
	return name.Value
}

func (p *Parser) parseIf() *Expr {
	tok := p.input.Peek()
	p.skipKw("if")
	cond := p.parseExpression()
	if p.isPunc("{") == nil {
		p.skipKw("then")
	}
	then := p.parseExpression()
	ret := &Expr{
		Type: If,
		Cond: cond,
		Then: then,
		File: tok.File,
		Line: tok.Line,
		Col:  tok.Col,
	}

	var elifBlocks []*Expr
	hasElif := false
	for p.isKw("elif") != nil {
		hasElif = true
		tok = p.input.Peek()
		p.input.Next()
		elifCond := p.parseExpression()
		if p.isPunc("{") == nil {
			p.skipKw("then")
		}
		elifThen := p.parseExpression()
		elifBlocks = append(elifBlocks, &Expr{
			Type: If,
			Cond: elifCond,
			Then: elifThen,
			File: tok.File,
			Line: tok.Line,
			Col:  tok.Col,
		})
	}

	if hasElif {
		if p.isKw("else") != nil {
			p.input.Next()
			elseBlock := p.parseExpression()
			elifBlocks = append(elifBlocks, elseBlock)
		} else {
			p.input.Error(tok, "Expecting 'else' after 'elif'")
		}
		ret.Else = &Expr{
			Type: Prog,
			Prog: elifBlocks,
		}
	} else if p.isKw("else") != nil {
		p.input.Next()
		ret.Else = p.parseExpression()
	}

	return ret
}

func (p *Parser) parseWhileExpr() *Expr {
	tok := p.input.Peek()
	p.skipKw("while")
	cond := p.parseExpression()
	if p.isPunc("{") == nil {
		p.input.Error(p.input.current, fmt.Sprintf("Expecting token '{', bot got: \"%s\"", p.input.current.Value))
	}
	body := p.parseBlock()
	return &Expr{
		Type: While,
		Cond: cond,
		Body: body,
		File: tok.File,
		Line: tok.Line,
		Col:  tok.Col,
	}
}

func (p *Parser) parseFunctionDecl() *Expr {
	tok := p.input.Peek()
	paramExprs := p.parseDelimited("(", ")", ",", func() *Expr {
		return &Expr{
			Type:  Var,
			Value: p.parseVarname(),
		}
	})
	var params []string
	for _, expr := range paramExprs {
		params = append(params, expr.Value.(string))
	}
	return &Expr{
		Type:   Fun,
		Params: params,
		Body:   p.parseExpression(),
		File:   tok.File,
		Line:   tok.Line,
		Col:    tok.Col,
	}
}

func (p *Parser) parseBoolExpr() *Expr {
	tok := p.input.Next()
	return &Expr{
		Type:   Bool,
		Value:  tok.Value == "true",
		File:   tok.File,
		Line:   tok.Line,
		Col:    tok.Col,
		Length: tok.Length,
	}
}

func (p *Parser) parseArrayDecl() *Expr {
	tok := p.input.Peek()
	p.skipKw("array")
	values := p.parseDelimited("{", "}", ",", p.parseExpression)
	return &Expr{
		Type: Array,
		Prog: values,
		File: tok.File,
		Line: tok.Line,
		Col:  tok.Col,
	}
}

func (p *Parser) parseTableDecl() *Expr {
	tok := p.input.Peek()
	p.skipKw("table")
	pairs := p.parseDelimited("{", "}", ",", p.parsePairDecl)
	return &Expr{
		Type: Table,
		Prog: pairs,
		File: tok.File,
		Line: tok.Line,
		Col:  tok.Col,
	}
}

func (p *Parser) parsePairDecl() *Expr {
	key := p.parseExpression()
	_, ok := key.Value.(string)
	if !ok {
		Error(key, "key must be of type string, but got: '%v'", key.Value)
	}
	// remove any occurences of whitespaces including space, tabs and newlines
	key.Value = strings.Join(strings.Fields(key.Value.(string)), "")
	p.skipPunc(":")
	value := p.parseExpression()
	tok := p.input.Peek()
	return &Expr{
		Type:  Pair,
		Left:  key,
		Right: value,
		File:  tok.File,
		Line:  tok.Line,
		Col:   tok.Col,
	}
}

func (p *Parser) parseAccessOrCall(expr *Expr) *Expr {
	// Function call
	if p.isPunc("(") != nil {
		return p.parseFunctionCall(expr)
	}
	// Array/table access
	if p.isPunc("[") != nil {
		return p.parseIndexExpr(expr)
	}
	// Field access
	if p.isPunc(".") != nil {
		return p.parseFieldAccessExpr(expr)
	}
	return expr
}

func (p *Parser) parseFieldAccessExpr(expr *Expr) *Expr {
	tok := p.input.Peek()
	p.skipPunc(".")
	fieldName := p.parseVarname()
	return &Expr{
		Type:  Var,        // We use Var type to represent field access
		Value: expr.Value, // Variable name that should be stored in the environment e.g. in "person.name" the string "person"
		Left:  expr,
		Index: &Expr{
			Type:  Str,       // Field name as a string
			Value: fieldName, // Field name e.g. in "person.name" the string "name"
			File:  tok.File,
			Line:  tok.Line,
			Col:   tok.Col,
		},
		File: tok.File,
		Line: tok.Line,
		Col:  tok.Col,
	}
}

func (p *Parser) parseIndexExpr(expr *Expr) *Expr {
	tok := p.input.Peek()
	p.skipPunc("[")
	indexExpr := p.parseExpression()
	switch indexExpr.Value.(type) {
	case string, int:
		// Valid type, proceed without doing anything
	default:
		p.input.Error(p.input.current, fmt.Sprintf("index expression must be of type string for tables or type int for arrays, but got '%T'", indexExpr.Value))
		return nil
	}
	p.skipPunc("]")

	return &Expr{
		Type:  Var, // We use Var type to represent variable or array access
		Value: expr.Value,
		Left:  expr,
		Index: indexExpr,
		File:  tok.File,
		Line:  tok.Line,
		Col:   tok.Col,
	}
}

func (p *Parser) parseNotExpr() *Expr {
	tok := p.input.Peek()
	p.skipKw("not")
	expr := p.parseExpression()

	return &Expr{
		Type:     Unary,
		Operator: "not",
		Right:    expr,
		File:     tok.File,
		Line:     tok.Line,
		Col:      tok.Col,
	}
}

func (p *Parser) parseReturnExpr() *Expr {
	tok := p.input.Peek()
	p.skipKw("return")
	var expr *Expr
	if p.isOp("<") != nil {
		p.input.Next()
		expr = p.parseExpression()
	} else {
		// Inject a false expression if the return has no argument
		expr = FALSE
	}
	return &Expr{
		Type:  Return,
		Right: expr,
		File:  tok.File,
		Line:  tok.Line,
		Col:   tok.Col,
	}
}

func (p *Parser) parseAtom() *Expr {
	var expr *Expr
	if p.isPunc("(") != nil {
		p.input.Next()
		expr = p.parseExpression()
		p.skipPunc(")")
	} else if p.isPunc("{") != nil {
		expr = p.parseBlock()
	} else if p.isKw("if") != nil {
		expr = p.parseIf()
	} else if p.isKw("while") != nil {
		expr = p.parseWhileExpr()
	} else if p.isKw("true") != nil || p.isKw("false") != nil {
		expr = p.parseBoolExpr()
	} else if p.isKw("fun") != nil {
		p.input.Next()
		expr = p.parseFunctionDecl()
	} else if p.isKw("array") != nil {
		expr = p.parseArrayDecl()
	} else if p.isKw("table") != nil {
		expr = p.parseTableDecl()
	} else if p.isKw("import") != nil {
		expr = p.parseImport()
	} else if p.isKw("not") != nil {
		expr = p.parseNotExpr()
	} else if p.isKw("return") != nil {
		expr = p.parseReturnExpr()

	} else {
		tok := p.input.Next()
		if tok.Type == "var" || tok.Type == "num" || tok.Type == "str" {
			expr = &Expr{
				Type:   ExprType(tok.Type),
				Value:  tok.Value,
				File:   tok.File,
				Line:   tok.Line,
				Col:    tok.Col,
				Length: tok.Length,
			}
		} else {
			p.unexpected(tok)
		}
	}

	return p.parseAccessOrCall(expr)
}

func (p *Parser) parseProgram() *Expr {
	var prog []*Expr
	for !p.input.Eof() {
		prog = append(prog, p.parseExpression())
	}
	return &Expr{
		Type: Prog,
		Prog: prog,
	}
}

func (p *Parser) parseImport() *Expr {
	tok := p.input.Next()
	path := p.parseExpression()
	return &Expr{
		Type: Import,
		Left: path,
		File: tok.File,
		Line: tok.Line,
		Col:  tok.Col,
	}
}

func (p *Parser) parseBlock() *Expr {
	block := p.parseEnclosed("{", "}", p.parseExpression)
	if len(block) == 0 {
		return FALSE
	}
	if len(block) == 1 {
		return block[0]
	}
	return &Expr{
		Type: Prog,
		Prog: block,
	}
}

func (p *Parser) parseExpression() *Expr {
	left := p.parseAtom()
	left = p.parseBinaryExpression(left, 0)
	return p.parseAccessOrCall(left)
}
