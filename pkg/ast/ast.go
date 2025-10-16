package ast

import "github.com/dywoq/minigo/pkg/token"

// Node represents the node of the AST tree.
type Node interface {
	node()
}

type Type string

// Value presentation in code:
//
//  var x = "Hi!"
//  //       ^
//  //       |
//  //    The value
type Value struct {
	Value string `json:"value"`
}

// Variable presentation in code:
//
//  var x = "Hi!"
//  var y string = "Bye!"
//
// or:
//
//  z := "Goodbye!"
type Variable struct {
	Name  string `json:"string"`
	Value Node   `json:"value"`
	Type  Type   `json:"type"`
}

// Function presentation in code:
//
//  func greet(name string, additional ...string) {
//      // ...
//  }
type Function struct {
	Name       string             `json:"name"`
	ReturnType string             `json:"return_type"`
	Exported   bool               `json:"exported"`
	Arguments  []FunctionArgument `json:"arguments"`
	Body       []Node             `json:"body"`
}

// FunctionArgument presentation in code:
//
//  func greet(name string) {
//      //      ^
//      //      |
//      // The function argument
//  }
type FunctionArgument struct {
	Identifier string     `json:"identifier"`
	Type       token.Kind `json:"type"`
	Variadic   bool       `json:"variadic"`
}

// Call presentation in code:
//
//  print("Hi!", 10, 23)
type Call struct {
	Identifier string         `json:"identifier"`
	Arguments  []CallArgument `json:"arguments"`
}

// CallArgument presentation in code:
//
//  print("Hi!")
//  //     ^
//  //     |
//  //   The call argument
type CallArgument struct {
	Type  token.Kind `json:"type"`
	Value string     `json:"value"`
}

// FunctionValue presentation in code:
//
//  var greet = func(name string) {
//     // ...
//  }
type FunctionValue struct {
	ReturnType string             `json:"return_type"`
	Arguments  []FunctionArgument `json:"arguments"`
	Body       []Node             `json:"node"`
}

// TypeConversion presentation in code:
//
//  x := string("Hi!")
type TypeConversion struct {
	To    string `json:"to"`
	Value Node   `json:"value"`
}

// BinaryExpression presentation in code:
//
//  			2 + 3
//  //  Left -- ^ ^ ^ -- Right
//  //            |
//  //     Binary expression
type BinaryExpression struct {
	Left     Node   `json:"left"`
	Operator string `json:"operator"`
	Right    Node   `json:"right"`
}

// File represents the whole parsed file with node statements.
type File struct {
	Statements []Node `json:"statements"`
}

const (
	TypeString       Type = "string"
	TypeInteger      Type = "integer"
	TypeFloat        Type = "float"
	TypeOfIdentifier Type = "of-identifier"
	TypeUnknown      Type = "unknown"
)

func TypeFromKind(t token.Kind) Type {
	switch t {
	case token.Integer:
		return TypeInteger
	case token.Float:
		return TypeFloat
	case token.String:
		return TypeString
	case token.Identifier:
		return TypeOfIdentifier
	}
	return TypeUnknown
}

func (Value) node()            {}
func (Variable) node()         {}
func (Function) node()         {}
func (FunctionArgument) node() {}
func (Call) node()             {}
func (CallArgument) node()     {}
func (FunctionValue) node()    {}
func (File) node()             {}
func (TypeConversion) node()   {}
