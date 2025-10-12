package token

// Kind represents the token kind in string.
type Kind string

// Collection represents a slice of reserved words,
// which can't be usually used in the code.
//
// To check if a value exists in the collection, you simply can use slices.Contains:
//  c := token.Collection{"func", "if", "return"}
//  if !slices.Contains(c, "func") {
//     panic("func doesn't exist in the collection")
//  }
type Collection []string

// Position represents the token position.
// It should be used as a pointer to provide correct position information.
type Position struct {
	Line     int `json:"line"`
	Column   int `json:"column"`
	Position int `json:"position"`
}

// Token is a stream of characters,
// with the literal, kind and position.
type Token struct {
	Literal  string    `json:"literal"`
	Kind     Kind      `json:"kind"`
	Position *Position `json:"position"`
}

// NewTokens returns a pointer to Token struct..
func NewToken(literal string, kind Kind, position *Position) *Token {
	return &Token{literal, kind, position}
}

// NewPosition returns a pointer to Position struct.
func NewPosition(line int, column int, position int) *Position {
	return &Position{line, column, position}
}

// A token kind.
var (
	Identifier Kind = "identifier"
	Integer    Kind = "integer"
	Float      Kind = "float"
	Type       Kind = "type"
	Keyword    Kind = "keyword"
	Separator  Kind = "separator"
	Eof        Kind = "eof"
	Illegal    Kind = "illegal"
)

// A collection of reserved words.
var (
	Keywords Collection = Collection{
		"const",
		"func",
		"import",
		"package",
		"type",
		"var",
		"map",
		"break",
		"case",
		"continue",
		"default",
		"else",
		"for",
		"if",
		"range",
		"return",
		"switch",
	}

	Separators Collection = Collection{
		",",
		";",
		"*",
		"(",
		")",
		"[",
		"]",
		"{",
		"}",
		"...",
		".",
		"//",
	}

	Types Collection = Collection{
		"int",
		"string",
		"bool",
		"float",
		"rune",
	}
)
