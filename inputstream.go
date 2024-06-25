package runevm

import (
	"fmt"
	"os"
)

type InputStream struct {
	filepath string
	source   string
	Pos      int
	line     int
	Col      int
}

func newInputStream(source string, filepath string) *InputStream {
	p := &InputStream{
		filepath: filepath,
		source:   source,
		Pos:      0,
		line:     1,
		Col:      1,
	}
	return p
}

func (p *InputStream) next() byte {
	if p.Pos >= len(p.source) {
		return 0
	}
	ch := p.source[p.Pos]
	p.Pos++
	if ch == '\n' {
		p.line++
		p.Col = 1
	} else {
		p.Col++
	}
	return ch
}

func (p *InputStream) peek() byte {
	if p.Pos >= len(p.source) {
		return 0
	}
	return p.source[p.Pos]
}

func (p *InputStream) eof() bool {
	ch := p.peek()
	eof := ch == 0
	return eof
}

func (p *InputStream) error(tok *Token, msg string) {
	fmt.Printf("error [%s:%d:%d]: %s\n", tok.File, tok.Line, tok.Col, msg)
	os.Exit(0)
}
