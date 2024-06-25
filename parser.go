package runevm

import (
	"fmt"
	"strings"
)

type Parser struct {
	input *TokenStream
}

func newParser(input *TokenStream) *Parser {
	return &Parser{input: input}
}

var FALSE = &expression{
	Type:  boolExpr,
	Value: false,
}

var precedence = map[string]int{
	"=":  1,
	"||": 2,
	"&&": 3,
	"<":  7, ">": 7, "<=": 7, ">=": 7, "==": 7, "!=": 7,
	"+": 10, "-": 10,
	"*": 20, "/": 20, "%": 20,
}

func (p *Parser) isPunc(ch string) *Token {
	tok := p.input.peek()
	if tok != nil && tok.Type == "punc" && (ch == "" || tok.Value == ch) {
		return tok
	}
	return nil
}

func (p *Parser) isKw(kw string) *Token {
	tok := p.input.peek()
	if tok != nil && tok.Type == "kw" && (kw == "" || tok.Value == kw) {
		return tok
	}
	return nil
}

func (p *Parser) isOp(op string) *Token {
	tok := p.input.peek()
	if tok != nil && tok.Type == "op" && (op == "" || tok.Value == op) {
		return tok
	}
	return nil
}

func (p *Parser) skipPunc(ch string) {
	if p.isPunc(ch) != nil {
		p.input.next()
	} else {
		p.input.error(p.input.current, fmt.Sprintf("Expecting punctuation: \"%s\"", ch))
	}
}

func (p *Parser) skipKw(kw string) {
	if p.isKw(kw) != nil {
		p.input.next()
	} else {
		p.input.error(p.input.current, fmt.Sprintf("Expecting keyword: \"%s\"", kw))
	}
}

func (p *Parser) unexpected(tok *Token) {
	if tok != nil {
		p.input.error(tok, fmt.Sprintf("Unexpected token: \"%s\"", tok.Value))
	}
}

func (p *Parser) parseBinaryExpression(left *expression, prec int) *expression {
	for {
		tok := p.isOp("")
		if tok == nil {
			break
		}
		opPrec := precedence[tok.Value]
		if opPrec <= prec {
			break
		}

		p.input.next()

		right := p.parseAtom()

		nextTok := p.isOp("")

		if nextTok != nil && precedence[nextTok.Value] > opPrec {
			right = p.parseBinaryExpression(right, opPrec)
		}

		var exprType exprType
		if tok.Value == "=" {
			exprType = assignExpr
		} else {
			exprType = binaryExpr
		}

		left = &expression{
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

func (p *Parser) parseDelimited(start, stop, separator string, parser func() *expression) []*expression {
	var a []*expression
	first := true
	p.skipPunc(start)
	for !p.input.eof() {
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

func (p *Parser) parseEnclosed(start string, stop string, parser func() *expression) []*expression {
	var a []*expression
	p.skipPunc(start)
	for !p.input.eof() {
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

func (p *Parser) parseFunctionCall(funcExpr *expression) *expression {
	tok := p.input.peek()
	return &expression{
		Type: callExpr,
		Func: funcExpr,
		Args: p.parseDelimited("(", ")", ",", p.parseExpression),
		File: tok.File,
		Line: tok.Line,
		Col:  tok.Col,
	}
}

func (p *Parser) parseVarname() string {
	name := p.input.next()
	if name.Type != "var" {
		p.input.error(name, fmt.Sprintf("Expecting variable name, but got: '%s'", name.Value))
	}
	return name.Value
}

func (p *Parser) parseIf() *expression {
	tok := p.input.peek()
	p.skipKw("if")
	cond := p.parseExpression()
	if p.isPunc("{") == nil {
		p.skipKw("then")
	}
	then := p.parseExpression()
	ret := &expression{
		Type: ifExpr,
		Cond: cond,
		Then: then,
		File: tok.File,
		Line: tok.Line,
		Col:  tok.Col,
	}

	var elifBlocks []*expression
	hasElif := false
	for p.isKw("elif") != nil {
		hasElif = true
		tok = p.input.peek()
		p.input.next()
		elifCond := p.parseExpression()
		if p.isPunc("{") == nil {
			p.skipKw("then")
		}
		elifThen := p.parseExpression()
		elifBlocks = append(elifBlocks, &expression{
			Type: ifExpr,
			Cond: elifCond,
			Then: elifThen,
			File: tok.File,
			Line: tok.Line,
			Col:  tok.Col,
		})
	}

	if hasElif {
		if p.isKw("else") != nil {
			p.input.next()
			elseBlock := p.parseExpression()
			elifBlocks = append(elifBlocks, elseBlock)
		} else {
			p.input.error(tok, "Expecting 'else' after 'elif'")
		}
		ret.Else = &expression{
			Type:  blockExpr,
			Block: elifBlocks,
		}
	} else if p.isKw("else") != nil {
		p.input.next()
		ret.Else = p.parseExpression()
	}

	return ret
}

func (p *Parser) parseWhileExpr() *expression {
	tok := p.input.peek()
	p.skipKw("while")
	cond := p.parseExpression()
	if p.isPunc("{") == nil {
		p.input.error(p.input.current, fmt.Sprintf("Expecting token '{', bot got: \"%s\"", p.input.current.Value))
	}
	body := p.parseBlock()
	return &expression{
		Type: whileExpr,
		Cond: cond,
		Body: body,
		File: tok.File,
		Line: tok.Line,
		Col:  tok.Col,
	}
}

func (p *Parser) parseFunctionDecl() *expression {
	tok := p.input.peek()
	paramExprs := p.parseDelimited("(", ")", ",", func() *expression {
		return &expression{
			Type:  varExpr,
			Value: p.parseVarname(),
		}
	})
	var params []string
	for _, expr := range paramExprs {
		params = append(params, expr.Value.(string))
	}
	return &expression{
		Type:   funExpr,
		Params: params,
		Body:   p.parseExpression(),
		File:   tok.File,
		Line:   tok.Line,
		Col:    tok.Col,
	}
}

func (p *Parser) parseBoolExpr() *expression {
	tok := p.input.next()
	return &expression{
		Type:   boolExpr,
		Value:  tok.Value == "true",
		File:   tok.File,
		Line:   tok.Line,
		Col:    tok.Col,
		Length: tok.Length,
	}
}

func (p *Parser) parseArrayDecl() *expression {
	tok := p.input.peek()
	p.skipKw("array")
	values := p.parseDelimited("{", "}", ",", p.parseExpression)
	return &expression{
		Type:  arrayExpr,
		Block: values,
		File:  tok.File,
		Line:  tok.Line,
		Col:   tok.Col,
	}
}

func (p *Parser) parseTableDecl() *expression {
	tok := p.input.peek()
	p.skipKw("table")
	pairs := p.parseDelimited("{", "}", ",", p.parsePairDecl)
	return &expression{
		Type:  tableExpr,
		Block: pairs,
		File:  tok.File,
		Line:  tok.Line,
		Col:   tok.Col,
	}
}

func (p *Parser) parsePairDecl() *expression {
	key := p.parseExpression()
	_, ok := key.Value.(string)
	if !ok {
		Error(key, "key must be of type string, but got: '%v'", key.Value)
	}
	// remove any occurences of whitespaces including space, tabs and newlines
	key.Value = strings.Join(strings.Fields(key.Value.(string)), "")
	p.skipPunc(":")
	value := p.parseExpression()
	tok := p.input.peek()
	return &expression{
		Type:  pairExpr,
		Left:  key,
		Right: value,
		File:  tok.File,
		Line:  tok.Line,
		Col:   tok.Col,
	}
}

func (p *Parser) parseAccessOrCall(expr *expression) *expression {
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

func (p *Parser) parseFieldAccessExpr(expr *expression) *expression {
	tok := p.input.peek()
	p.skipPunc(".")
	fieldName := p.parseVarname()
	return &expression{
		Type:  varExpr,    // We use Var type to represent field access
		Value: expr.Value, // Variable name that should be stored in the environment e.g. in "person.name" the string "person"
		Left:  expr,
		Index: &expression{
			Type:  strExpr,   // Field name as a string
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

func (p *Parser) parseIndexExpr(expr *expression) *expression {
	tok := p.input.peek()
	p.skipPunc("[")
	indexExpr := p.parseExpression()
	switch indexExpr.Value.(type) {
	case string, int:
		// Valid type, proceed without doing anything
	default:
		p.input.error(p.input.current, fmt.Sprintf("index expression must be of type string for tables or type int for arrays, but got '%T'", indexExpr.Value))
		return nil
	}
	p.skipPunc("]")

	return &expression{
		Type:  varExpr, // We use Var type to represent variable or array access
		Value: expr.Value,
		Left:  expr,
		Index: indexExpr,
		File:  tok.File,
		Line:  tok.Line,
		Col:   tok.Col,
	}
}

func (p *Parser) parseNotExpr() *expression {
	tok := p.input.peek()
	p.skipKw("not")
	expr := p.parseExpression()

	return &expression{
		Type:     unaryExpr,
		Operator: "not",
		Right:    expr,
		File:     tok.File,
		Line:     tok.Line,
		Col:      tok.Col,
	}
}

func (p *Parser) parseReturnExpr() *expression {
	tok := p.input.peek()
	p.skipKw("return")
	var expr *expression
	if p.isOp("=") != nil {
		p.input.next()
		expr = p.parseExpression()
	} else {
		// Inject a false expression if the return has no argument
		expr = FALSE
	}
	return &expression{
		Type:  returnExpr,
		Right: expr,
		File:  tok.File,
		Line:  tok.Line,
		Col:   tok.Col,
	}
}

func (p *Parser) parseBreakExpr() *expression {
	tok := p.input.peek()
	p.skipKw("break")
	return &expression{
		Type:  breakExpr,
		Right: FALSE,
		File:  tok.File,
		Line:  tok.Line,
		Col:   tok.Col,
	}
}

func (p *Parser) parseContinueExpr() *expression {
	tok := p.input.peek()
	p.skipKw("continue")
	return &expression{
		Type:  continueExpr,
		Right: FALSE,
		File:  tok.File,
		Line:  tok.Line,
		Col:   tok.Col,
	}
}

func (p *Parser) parseAtom() *expression {
	var expr *expression
	if p.isPunc("(") != nil {
		p.input.next()
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
		p.input.next()
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
	} else if p.isKw("break") != nil {
		expr = p.parseBreakExpr()
	} else if p.isKw("continue") != nil {
		expr = p.parseContinueExpr()

	} else {
		tok := p.input.next()
		if tok.Type == "var" || tok.Type == "num" || tok.Type == "str" {
			expr = &expression{
				Type:   exprType(tok.Type),
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

func (p *Parser) parseProgram() *expression {
	var prog []*expression
	for !p.input.eof() {
		prog = append(prog, p.parseExpression())
	}
	return &expression{
		Type:  blockExpr,
		Block: prog,
	}
}

func (p *Parser) parseImport() *expression {
	tok := p.input.next()
	path := p.parseExpression()
	return &expression{
		Type: importExpr,
		Left: path,
		File: tok.File,
		Line: tok.Line,
		Col:  tok.Col,
	}
}

func (p *Parser) parseBlock() *expression {
	block := p.parseEnclosed("{", "}", p.parseExpression)
	if len(block) == 0 {
		return FALSE
	}
	if len(block) == 1 {
		return block[0]
	}
	return &expression{
		Type:  blockExpr,
		Block: block,
	}
}

func (p *Parser) parseExpression() *expression {
	left := p.parseAtom()
	left = p.parseBinaryExpression(left, 0)
	return p.parseAccessOrCall(left)
}
