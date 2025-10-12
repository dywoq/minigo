package scanner

import (
	"errors"
	"fmt"
	"slices"
	"unicode"

	"github.com/dywoq/minigo/pkg/token"
)

type context interface {
	advance(n int) error
	current() (rune, error)
	eof() bool
	debug(v ...any) error
	debugf(format string, v ...any) error
	position() *token.Position
	slice(start, end int) (string, error)
	new(literal string, kind token.Kind) *token.Token
}

type tokenizer func(c context) (*token.Token, error)

var errNoMatch = errors.New("no match")

func tokenizeNumber(c context) (*token.Token, error) {
	r, err := c.current()
	if err != nil {
		return nil, err
	}
	if !unicode.IsNumber(r) {
		return nil, errNoMatch
	}
	start := c.position().Position
	c.debug("tokenizing number")
	for {
		if err := c.advance(1); err != nil {
			break
		}

		r, err = c.current()
		if err != nil {
			break
		}

		if !unicode.IsNumber(r) {
			break
		}
	}
	if r == '.' {
		c.debug("detected dot, consuming fractional part")
		if err := c.advance(1); err != nil {
			return nil, fmt.Errorf("expected a number after dot at %d", c.position().Line)
		}

		r, _ = c.current()
		if !unicode.IsNumber(r) {
			return nil, fmt.Errorf("expected a number after dot at %d", c.position().Line)
		}

		for {
			if err := c.advance(1); err != nil {
				break
			}

			r, err = c.current()
			if err != nil {
				break
			}

			if !unicode.IsNumber(r) {
				break
			}
		}
		str, err := c.slice(start, c.position().Position)
		if err != nil {
			return nil, err
		}
		return c.new(str, token.Float), nil
	}

	str, err := c.slice(start, c.position().Position)
	if err != nil {
		return nil, err
	}
	return c.new(str, token.Integer), nil
}

func tokenizeKeyword(c context) (*token.Token, error) {
	str, err := selectWordAndCheck(c, token.Keywords)
	if err != nil {
		return nil, err
	}
	return c.new(str, token.Keyword), nil
}

func tokenizeType(c context) (*token.Token, error) {
	str, err := selectWordAndCheck(c, token.Types)
	if err != nil {
		return nil, err
	}
	return c.new(str, token.Type), nil
}

func selectWordAndCheck(c context, collection token.Collection) (string, error) {
	if r, _ := c.current(); !unicode.IsLetter(r) {
		return "", errNoMatch
	}
	start := c.position().Position
	for {
		r, _ := c.current()
		if c.eof() || !unicode.IsLetter(r) {
			break
		}
		c.advance(1)
	}
	str, err := c.slice(start, c.position().Position)
	if err != nil {
		return "", err
	}
	if !slices.Contains(collection, str) {
		c.position().Position = start
		return "", errNoMatch
	}
	return str, nil
}
