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

func NewInputStream(source string, filepath string) *InputStream {
	// src := source + "\000"
	p := &InputStream{
		filepath: filepath,
		source:   source,
		Pos:      0,
		line:     1,
		Col:      0,
	}
	return p
}

func (p *InputStream) Next() byte {
	ch := p.source[p.Pos]
	p.Pos++
	if ch == '\n' {
		p.line++
		p.Col = 0
	} else {
		p.Col++
	}
	return ch
}

func (p *InputStream) Peek() byte {
	if p.Pos >= len(p.source) {
		return 0
	}
	return p.source[p.Pos]
}

func (p *InputStream) Eof() bool {
	ch := p.Peek()
	eof := ch == 0
	return eof
}

func (p *InputStream) Error(tok *Token, msg string) {
	fmt.Printf("error [%s:%d:%d]: %s\n", tok.File, tok.Line, tok.Col, msg)
	os.Exit(0)
}
