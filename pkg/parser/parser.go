package parser

import (
	"errors"
	"io"

	"github.com/dywoq/minigo/pkg/ast"
	"github.com/dywoq/minigo/pkg/token"
)

type Parser struct {
	tokens  []*token.Token
	pos     int
	parsing bool
	d       debug
	mini    []mini
}

var ErrWorking = errors.New("parser: parser is working")

var miniParsers = []mini{
	parseVariable,
}

type debug struct {
	p  *Parser
	on bool
	w  io.Writer
}

// Debug returns the methods that allow you
// to manage the debugger.
func (p *Parser) Debug() debug {
	return p.d
}

// SetTokens sets the tokens.
// If the parser is currently working, the function returns ErrWorking.
// If t slice is nil - an error is returned.
func (p *Parser) SetTokens(t []*token.Token) error {
	if p.parsing {
		return ErrWorking
	}
	if t == nil {
		return errors.New("parser: given tokens slice is nil")
	}
	p.tokens = t
	return nil
}

// Parse parses the given tokens.
// Returns ast.Program node with the set statements.
func (p *Parser) Parse() (ast.Node, error) {
	panic("implement me!")
}

func (p *Parser) advance(n int) error {
	if p.pos+n >= len(p.tokens) {
		return errors.New("p.pos+n is overflow")
	}
	for range n {
		p.pos++
	}
	return nil
}

func (p *Parser) eof() bool {
	t, _ := p.current()
	return p.pos >= len(p.tokens) || t.Kind == token.Eof
}

func (p *Parser) current() (*token.Token, error) {
	if p.eof() {
		return nil, io.EOF
	}
	return p.tokens[p.pos], nil
}

// On reports whether the debug mode is on.
func (d *debug) On() bool {
	return d.on
}

// Set sets the debug mode to b.
// If the parser is currently working, the function returns ErrWorking.
func (d *debug) Set(b bool) error {
	if d.p.parsing {
		return ErrWorking
	}
	d.on = b
	return nil
}

// SetWriter sets the writer.
// If w is nil or the parser is currently working, an error is returned.
func (d *debug) SetWriter(w io.Writer) error {
	if d.p.parsing {
		return ErrWorking
	}
	if w == nil {
		return errors.New("debug: given io.Writer is nil")
	}
	d.w = w
	return nil
}
