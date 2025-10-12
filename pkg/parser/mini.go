package parser

import "github.com/dywoq/minigo/pkg/token"

type context interface {
	advance(n int) error
	current() (*token.Token, error)
	eof() bool
}

type mini func(context) (*token.Token, error)

func parseVariable(c context) (*token.Token, error) {
	return nil, nil
}
