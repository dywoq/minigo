package ast

import "github.com/dywoq/minigo/pkg/token"

// Node represents the node of the AST tree.
type Node interface {
	node()
}

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
	Literal Node       `json:"literal"`
	Type    token.Kind `json:"type"`
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

func (Value) node()            {}
func (Variable) node()         {}
func (Function) node()         {}
func (FunctionArgument) node() {}
func (Call) node()             {}
func (CallArgument) node()     {}
func (FunctionValue) node()    {}
