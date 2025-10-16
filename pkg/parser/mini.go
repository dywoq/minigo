package parser

import (
	"fmt"
	"slices"

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

func parseDeclaration(c context) (ast.Node, error) {
	t := c.current()
	if t == nil {
		return nil, fmt.Errorf("unexpected EOF")
	}

	switch t.Kind {
	case token.Identifier:
		next := c.peek(1)
		if next != nil && next.Literal == ":=" {
			name := t.Literal
			c.advance(1)
			return parseVariable(name, c)
		}
		return nil, fmt.Errorf("unexpected identifier without := at %v", t.Position)

	case token.Keyword:
		if t.Literal == "func" {
			return parseFunction(c)
		}
		return nil, fmt.Errorf("unexpected keyword %q at %v", t.Literal, t.Position)
	}

	return nil, fmt.Errorf("unexpected declaration start: %v", t)
}

func parseVariable(name string, c context) (ast.Node, error) {
	_, err := c.expectLiteral(":=")
	if err != nil {
		return nil, err
	}
	val, kind, err := parseExpression(c, 0)
	if err != nil {
		return nil, err
	}
	return ast.Variable{
		Name:  name,
		Type:  ast.TypeFromKind(kind),
		Value: val,
	}, nil
}

func parseValue(c context) (ast.Node, token.Kind, error) {
	t := c.current()
	if t == nil {
		return nil, token.Illegal, fmt.Errorf("unexpected EOF in value")
	}

	if t.Literal == "(" {
		c.advance(1)
		expr, kind, err := parseExpression(c, 0)
		if err != nil {
			return nil, token.Illegal, err
		}
		_, err = c.expectLiteral(")")
		if err != nil {
			return nil, token.Illegal, err
		}

		if bin, ok := expr.(ast.BinaryExpression); ok {
			bin.HasParens = true
			return bin, kind, nil
		}

		return expr, kind, nil
	}

	switch t.Kind {
	case token.Integer, token.Float, token.String:
		c.advance(1)
		return ast.Value{Value: t.Literal}, t.Kind, nil

	case token.Type, token.Identifier:
		next := c.peek(1)
		if next != nil && next.Literal == "(" {
			to, err := c.expectKinds(token.Type, token.Identifier)
			if err != nil {
				return nil, token.Illegal, err
			}
			_, err = c.expectLiteral("(")
			if err != nil {
				return nil, token.Illegal, err
			}
			val, _, err := parseValue(c)
			if err != nil {
				return nil, token.Illegal, err
			}
			_, err = c.expectLiteral(")")
			if err != nil {
				return nil, token.Illegal, err
			}
			return ast.TypeConversion{
				To:    to.Literal,
				Value: val,
			}, t.Kind, nil
		}
		c.advance(1)
		return ast.Value{Value: t.Literal}, t.Kind, nil
	}

	return nil, token.Illegal, fmt.Errorf("unknown value: %v", t)
}

var precedence = map[string]int{
	"+": 1,
	"-": 1,
	"*": 2,
	"/": 2,
}

func parseExpression(c context, minPrec int) (ast.Node, token.Kind, error) {
	left, kind, err := parseValue(c)
	if err != nil {
		return nil, token.Illegal, err
	}

	for {
		t := c.current()
		if t == nil || !slices.Contains(token.BinaryOperators, t.Literal) {
			break
		}

		opPrec := precedence[t.Literal]
		if opPrec < minPrec {
			break
		}

		op := t.Literal
		c.advance(1)

		right, rKind, err := parseExpression(c, opPrec+1)
		if err != nil {
			return nil, token.Illegal, err
		}

		left = ast.BinaryExpression{
			Left:     left,
			Operator: op,
			Right:    right,
		}
		kind = rKind
	}

	return left, kind, nil
}

func parseFunction(c context) (ast.Node, error) {
	_, err := c.expectLiteral("func")
	if err != nil {
		return nil, err
	}

	nameToken, err := c.expectKind(token.Identifier)
	if err != nil {
		return nil, err
	}
	name := nameToken.Literal

	_, err = c.expectLiteral("(")
	if err != nil {
		return nil, err
	}

	var args []ast.FunctionArgument
	var variadicSeen bool

	for !c.eof() && c.current().Literal != ")" {
		argNameToken, err := c.expectKind(token.Identifier)
		if err != nil {
			return nil, err
		}
		isVariadic := false
		if c.current().Literal == "..." {
			if variadicSeen {
				return nil, fmt.Errorf("only one variadic parameter allowed")
			}
			isVariadic = true
			variadicSeen = true
			c.advance(1)
		}
		argTypeToken, err := c.expectKinds(token.Type, token.Identifier)
		if err != nil {
			return nil, err
		}

		arg := ast.FunctionArgument{
			Identifier: argNameToken.Literal,
			Type:       argTypeToken.Kind,
			Variadic:   isVariadic,
		}
		args = append(args, arg)
		if variadicSeen && c.current().Literal == "," {
			return nil, fmt.Errorf("variadic parameter must be last")
		}

		if c.current().Literal == "," {
			c.advance(1)
		}
	}

	_, err = c.expectLiteral(")")
	if err != nil {
		return nil, err
	}

	retType := ""
	if !c.eof() && c.current().Literal != "{" {
		tok := c.current()
		if tok.Kind == token.Type || tok.Kind == token.Identifier {
			retType = tok.Literal
			c.advance(1)
		}
	}

	_, err = c.expectLiteral("{")
	if err != nil {
		return nil, err
	}

	var body []ast.Node
	for !c.eof() && c.current().Literal != "}" {
		stmt, err := parseDeclaration(c)
		if err != nil {
			stmt, _, err = parseExpression(c, 0)
			if err != nil {
				return nil, err
			}
		}
		body = append(body, stmt)
	}

	_, err = c.expectLiteral("}")
	if err != nil {
		return nil, err
	}

	return ast.Function{
		Name:       name,
		ReturnType: retType,
		Arguments:  args,
		Body:       body,
	}, nil
}
