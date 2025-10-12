package scanner

import (
	"errors"
	"fmt"
	"unicode"

	"github.com/dywoq/minigo/token"
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

	str, err := c.slice(start, c.position().Position-1)
	if err != nil {
		return nil, err
	}
	return c.new(str, token.Integer), nil
}
