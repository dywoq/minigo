package parser

import (
	"fmt"

	"github.com/dywoq/minigo/pkg/ast"
	"github.com/dywoq/minigo/pkg/token"
)

type context interface {
	advance(n int) error
	current() *token.Token
	peek(n int) *token.Token
	eof() bool
	debug(v ...any) error
	debugf(format string, v ...any) error
	expectLiteral(literal string) (*token.Token, error)
	expectLiterals(literals ...string) (*token.Token, error)
	expectKind(kind token.Kind) (*token.Token, error)
	expectKinds(kinds ...token.Kind) (*token.Token, error)
}

type mini func(context) (ast.Node, error)

type expectError struct {
	got      *token.Token
	literals []string
	kinds    []token.Kind
}

func (e *expectError) Error() string {
	if len(e.literals) > 0 {
		return fmt.Sprintf("expected one of %v, got %q at %v", e.literals, e.got.Literal, e.got.Position)
	}
	if len(e.kinds) > 0 {
		return fmt.Sprintf("expected one of %v, got %v at %v", e.kinds, e.got.Kind, e.got.Position)
	}
	return fmt.Sprintf("unexpected token %v at %v", e.got, e.got.Position)
}

// Handles: var x = 2, var x int = 3, x := 2
func parseDeclaration(c context) (ast.Node, error) {
	t := c.current()
	if t == nil {
		return nil, fmt.Errorf("unexpected EOF")
	}

	switch {
	// short variable: x := 2
	case t.Kind == token.Identifier:
		next := c.peek(1)
		if next != nil && next.Literal == ":=" {
			name := t.Literal
			c.advance(1) // consume identifier
			return parseShortVariable(name, c)
		}
		return nil, fmt.Errorf("unexpected identifier without := at %v", t.Position)

	// normal variable: var x = 2 or var x int = 3
	case t.Literal == "var":
		c.advance(1)
		return parseVariable(c)
	}

	return nil, fmt.Errorf("unexpected declaration start: %v", t)
}

func parseVariable(c context) (ast.Node, error) {
	ident, err := c.expectKind(token.Identifier)
	if err != nil {
		return nil, err
	}

	var typeName string
	t := c.current()
	if t != nil && t.Kind == token.Type {
		typeName = t.Literal
		c.advance(1)
	}

	var val ast.Node
	if t := c.current(); t != nil && t.Literal == "=" {
		c.advance(1)
		v, _, err := parseValue(c)
		if err != nil {
			return nil, err
		}
		val = v
	}

	return ast.Variable{
		Name:  ident.Literal,
		Type:  typeName,
		Value: val,
	}, nil
}

func parseShortVariable(name string, c context) (ast.Node, error) {
	_, err := c.expectLiteral(":=")
	if err != nil {
		return nil, err
	}
	val, kind, err := parseValue(c)
	if err != nil {
		return nil, err
	}
	return ast.Variable{
		Name:  name,
		Type:  string(kind),
		Value: val,
	}, nil
}

func parseValue(c context) (ast.Node, token.Kind, error) {
	t := c.current()
	if t == nil {
		return nil, token.Illegal, fmt.Errorf("unexpected EOF in value")
	}

	switch t.Kind {
	case token.Integer, token.Float, token.String:
		c.advance(1)
		return ast.Value{Value: t.Literal}, t.Kind, nil

	case token.Type, token.Identifier:
		to, err := c.expectKinds(token.Type, token.Identifier)
		if err != nil {
			return nil, token.Illegal, err
		}
		if _, err := c.expectLiteral("("); err != nil {
			return nil, token.Illegal, err
		}
		val, _, err := parseValue(c)
		if err != nil {
			return nil, token.Illegal, err
		}
		if _, err := c.expectLiteral(")"); err != nil {
			return nil, token.Illegal, err
		}
		return ast.TypeConversion{
			To:    to.Literal,
			Value: val,
		}, t.Kind, nil
	}

	return nil, token.Illegal, fmt.Errorf("unknown value: %v", t)
}
