package scanner

import (
	"errors"
	"fmt"
	"io"
	"time"
	"unicode"

	"github.com/dywoq/minigo/cmd/minigo/token"
)

// Scanner is responsible for scanning the code
// and tokenizing it.
type Scanner struct {
	input      []byte // used to prevent multiple using io.ReadAll
	r          io.Reader
	p          *token.Position
	d          debug
	scanning   bool
	tokenizers []tokenizer
}

type debug struct {
	s  *Scanner
	w  io.Writer
	on bool
}

var (
	ErrWorking = errors.New("scanner is working")
)

// New returns a pointer to Scanner with the given io.Reader instance.
// If something fails, the function returns nil and an error.
func New(r io.Reader) (*Scanner, error) {
	bytes, err := io.ReadAll(r)
	if err != nil {
		return nil, err
	}
	return &Scanner{
		bytes,
		r,
		token.NewPosition(1, 1, 0),
		debug{},
		false,
		[]tokenizer{tokenizeNumber},
	}, nil
}

// NewDebug returns a pointer to Scanner,
// the only difference from New is that the debugger is automatically on.
// But you must pass io.Writer for the debugger to write messages.
func NewDebug(r io.Reader, w io.Writer) (*Scanner, error) {
	s := &Scanner{
		r:          r,
		p:          token.NewPosition(1, 1, 0),
		scanning:   false,
		tokenizers: []tokenizer{tokenizeNumber},
	}
	s.d = debug{s, w, true}
	bytes, err := io.ReadAll(r)
	if err != nil {
		return nil, err
	}
	s.input = bytes
	return s, nil
}

// Debug returns the methods that allows you
// to manage the debugger.
func (s *Scanner) Debug() debug {
	return s.d
}

// SetReader sets the reader and updates the underlying input.
// Returns an error if scanner is already working,
// or something went wrong when trying to get the content with io.ReadAll.
func (s *Scanner) SetReader(r io.Reader) error {
	if s.scanning {
		return ErrWorking
	}
	bytes, err := io.ReadAll(r)
	if err != nil {
		return err
	}
	s.input = bytes
	return nil
}

// Scan scans the given input, and tokenizes it.
// If the current character doesn't satisfy the requirements of one of tokenizers,
// Scan tries other tokenizer.
func (s *Scanner) Scan() ([]*token.Token, error) {
	result := []*token.Token{}
	for !s.eof() {
		for _, tokenizer := range s.tokenizers {
			s.skipWhitespace()
			tok, err := tokenizer(s)
			if err != nil {
				if err == errNoMatch {
					s.debug("got errNoMatch, trying other tokenizer")
					continue
				}
				return nil, err
			}
			result = append(result, tok)
		}
	}
	result = append(result, token.NewToken("", token.Eof, s.p))
	return result, nil
}

func (s *Scanner) skipWhitespace() {
	for {
		if r, _ := s.current(); !unicode.IsSpace(r) {
			break
		}
		s.advance(1)
	}
}

func (s *Scanner) slice(start, end int) (string, error) {
	switch {
	case start < 0:
		return "", errors.New("start is negative")
	case end > len(s.input):
		return "", errors.New("end is higher than the input")
	case start > end:
		return "", errors.New("start is higher than the end")
	}
	return string(s.input[start:end]), nil
}

func (s *Scanner) new(literal string, kind token.Kind) *token.Token {
	return token.NewToken(literal, kind, s.p)
}

func (s *Scanner) advance(n int) error {
	if s.eof() {
		return io.EOF
	}
	for range n {
		r, err := s.current()
		if err != nil {
			return err
		}
		s.p.Position++
		if r == '\n' || s.eof() {
			s.p.Column = 1
			s.p.Line++
		} else {
			s.p.Column++
		}
	}
	return nil
}

func (s *Scanner) current() (rune, error) {
	if s.eof() {
		return 0, io.EOF
	}
	return rune(s.input[s.p.Position]), nil
}

func (s *Scanner) eof() bool {
	return s.p.Position >= len(s.input)
}

func (s *Scanner) debug(v ...any) error {
	if !s.d.on {
		return nil
	}
	_, err := io.WriteString(s.d.w, fmt.Sprintf("%s %v\n", time.Now().String(), v))
	if err != nil {
		return err
	}
	return err
}

func (s *Scanner) debugf(format string, v ...any) error {
	res := fmt.Sprintf(format, v...)
	return s.debug(res)
}

func (s *Scanner) position() *token.Position {
	return s.p
}

// Set turns on the debugging mode.
func (d *debug) Set(b bool) {
	d.on = b
}

// On returns true if the debugging mode is on.
func (d *debug) On() bool {
	return d.on
}

// SetWriter sets a instance that implements io.Writer interface.
func (d *debug) SetWriter(w io.Writer) {
	d.w = w
}
