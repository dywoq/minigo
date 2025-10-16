package parser

import (
	"errors"
	"fmt"
	"io"
	"slices"
	"time"

	"github.com/dywoq/minigo/pkg/ast"
	"github.com/dywoq/minigo/pkg/token"
)

// Parser is responsible for parsing tokens given from the scanner.
type Parser struct {
	tokens  []*token.Token
	pos     int
	parsing bool
	d       debug
	mini    []mini
}

var ErrWorking = errors.New("parser: parser is working")

var miniParsers = []mini{
	parseDeclaration,
}

type debug struct {
	p  *Parser
	on bool
	w  io.Writer
}

// New returns a pointer to Parser with the debug automatically turned off.
// Notice that the debug writer is nil.
func New(tokens []*token.Token) (*Parser, error) {
	p := &Parser{
		tokens:  tokens,
		pos:     0,
		parsing: false,
		mini:    miniParsers,
	}
	p.d = debug{p: p, on: false, w: nil}
	return p, nil
}

// NewDebug returns a pointer to Parser with the debug automatically turned on.
// w writer is used by the debugger.
func NewDebug(tokens []*token.Token, w io.Writer) (*Parser, error) {
	p := &Parser{
		tokens:  tokens,
		pos:     0,
		parsing: false,
		mini:    miniParsers,
	}
	p.d = debug{p: p, on: true, w: w}
	return p, nil
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
// Returns ast.File node with the set statements.
func (p *Parser) Parse() (ast.File, error) {
	p.d.p.parsing = true
	p.debug("starting parsing")
	defer func() {
		p.d.p.parsing = false
		p.debug("ending parsing")
	}()

	f := ast.File{}
	for !p.eof() {
		t := p.current()
		if t.Kind == token.Eof {
			break
		}
		for _, mini := range p.mini {
			r, err := mini(p)
			if err != nil {
				return ast.File{}, err
			}
			f.Statements = append(f.Statements, r)
			p.debugf("parsed node at %v", p.pos)
		}
	}
	return f, nil
}

func (p *Parser) advance(n int) error {
	if p.pos+n >= len(p.tokens) {
		return errors.New("p.pos+n is overflow")
	}
	p.pos += n
	p.debugf("advancing by %d", n)
	return nil
}

func (p *Parser) eof() bool {
	return p.pos >= len(p.tokens)
}

func (p *Parser) current() *token.Token {
	if p.eof() {
		return nil
	}
	p.debug("getting current token")
	return p.tokens[p.pos]
}

func (p *Parser) debug(v ...any) error {
	if !p.d.on {
		return nil
	}
	_, err := io.WriteString(p.d.w, fmt.Sprintf("%s %v\n", time.Now().String(), v))
	if err != nil {
		return err
	}
	return nil
}

func (p *Parser) debugf(format string, v ...any) error {
	res := fmt.Sprintf(format, v...)
	return p.debug(res)
}

func (p *Parser) expectLiteral(literal string) (*token.Token, error) {
	t := p.current()
	p.debugf("expect literal \"%s\"...", literal)
	if t.Literal != literal {
		return nil, &expectError{
			got:      t,
			literals: []string{literal},
			kinds:    []token.Kind{},
		}
	}
	p.advance(1)
	return t, nil
}

func (p *Parser) expectLiterals(literals ...string) (*token.Token, error) {
	t := p.current()
	p.debugf("expect literals %v...", literals)
	if slices.Contains(literals, t.Literal) {
		p.advance(1)
		return t, nil
	}
	return nil, &expectError{
		got:      t,
		literals: literals,
		kinds:    []token.Kind{},
	}
}

func (p *Parser) expectKind(kind token.Kind) (*token.Token, error) {
	t := p.current()
	p.debugf("expect kind \"%s\"...", kind)
	if t.Kind != kind {
		return nil, &expectError{
			got:      t,
			literals: []string{},
			kinds:    []token.Kind{kind},
		}
	}
	p.advance(1)
	return t, nil
}

func (p *Parser) expectKinds(kinds ...token.Kind) (*token.Token, error) {
	t := p.current()
	p.debugf("expect kind %v...", kinds)
	if slices.Contains(kinds, t.Kind) {
		p.advance(1)
		return t, nil
	}
	return nil, &expectError{
		got:      t,
		literals: []string{},
		kinds:    kinds,
	}
}

func (p *Parser) peek(n int) *token.Token {
	if p.pos+n >= len(p.tokens) {
		return nil
	}
	return p.tokens[p.pos+n]
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
