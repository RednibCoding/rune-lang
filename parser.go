package runevm

import (
	"fmt"
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

func (p *Parser) maybeBinary(left *Expr, myPrec int) *Expr {
	for tok := p.isOp(""); tok != nil; tok = p.isOp("") {
		hisPrec := PRECEDENCE[tok.Value]
		if hisPrec <= myPrec {
			break
		}
		p.input.Next()
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
			Right:    p.maybeBinary(p.parseAtom(), hisPrec),
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

func (p *Parser) parseCall(funcExpr *Expr) *Expr {
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
		p.input.Error(p.input.current, "Expecting variable name")
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
			elifBlocks = append(elifBlocks, p.parseExpression())
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

func (p *Parser) parseWhile() *Expr {
	tok := p.input.Peek()
	p.skipKw("while")
	cond := p.parseExpression()
	if p.isPunc("{") == nil {
		p.input.Error(p.input.current, fmt.Sprintf("Expecting token '{', bot got: \"%s\"", p.input.current.Value))
	}
	body := p.parseProg()
	return &Expr{
		Type: While,
		Cond: cond,
		Body: body,
		File: tok.File,
		Line: tok.Line,
		Col:  tok.Col,
	}
}

func (p *Parser) parseFun() *Expr {
	tok := p.input.Peek()
	varNames := p.parseDelimited("(", ")", ",", func() *Expr {
		return &Expr{
			Type:  Var,
			Value: p.parseVarname(),
		}
	})
	var vars []string
	for _, expr := range varNames {
		vars = append(vars, expr.Value.(string))
	}
	return &Expr{
		Type: Fun,
		Vars: vars,
		Body: p.parseExpression(),
		File: tok.File,
		Line: tok.Line,
		Col:  tok.Col,
	}
}

func (p *Parser) parseBool() *Expr {
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

func (p *Parser) parseArray() *Expr {
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

func (p *Parser) maybeCall(expr func() *Expr) *Expr {
	exprNode := expr()
	// Function call
	if p.isPunc("(") != nil {
		return p.parseCall(exprNode)
	}
	// Array access
	if p.isPunc("[") != nil {
		return p.parseIndex(exprNode)
	}
	return exprNode
}

func (p *Parser) parseIndex(arrayExpr *Expr) *Expr {
	tok := p.input.Next() // skip the '['
	indexExpr := p.parseExpression()
	p.skipPunc("]")

	return &Expr{
		Type:  Var, // We use Var type to represent variable or array access
		Value: arrayExpr.Value,
		Index: indexExpr,
		File:  tok.File,
		Line:  tok.Line,
		Col:   tok.Col,
	}
}

func (p *Parser) parseAtom() *Expr {
	return p.maybeCall(func() *Expr {
		if p.isPunc("(") != nil {
			p.input.Next()
			exp := p.parseExpression()
			p.skipPunc(")")
			return exp
		}
		if p.isPunc("{") != nil {
			return p.parseProg()
		}
		if p.isKw("if") != nil {
			return p.parseIf()
		}
		if p.isKw("while") != nil {
			return p.parseWhile()
		}
		if p.isKw("true") != nil || p.isKw("false") != nil {
			return p.parseBool()
		}
		if p.isKw("fun") != nil {
			p.input.Next()
			return p.parseFun()
		}
		if p.isKw("array") != nil {
			return p.parseArray()
		}
		tok := p.input.Next()
		if tok.Type == "var" || tok.Type == "num" || tok.Type == "str" {
			return &Expr{
				Type:   ExprType(tok.Type),
				Value:  tok.Value,
				File:   tok.File,
				Line:   tok.Line,
				Col:    tok.Col,
				Length: tok.Length,
			}
		}

		p.unexpected(tok)
		return nil
	})
}

func (p *Parser) parseToplevel() *Expr {
	var prog []*Expr
	for !p.input.Eof() {
		prog = append(prog, p.parseExpression())
		// if !p.input.Eof() {
		// p.skipPunc(";")

		// }
	}
	return &Expr{
		Type: Prog,
		Prog: prog,
	}
}

func (p *Parser) parseProg() *Expr {
	prog := p.parseEnclosed("{", "}", p.parseExpression)
	if len(prog) == 0 {
		return FALSE
	}
	if len(prog) == 1 {
		return prog[0]
	}
	return &Expr{
		Type: Prog,
		Prog: prog,
	}
}

func (p *Parser) parseExpression() *Expr {
	return p.maybeCall(func() *Expr {
		return p.maybeBinary(p.parseAtom(), 0)
	})
}
