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
	Literal   Node               `json:"literal"`
	Arguments []FunctionArgument `json:"arguments"`
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

func (Value) node()            {}
func (Variable) node()         {}
func (Function) node()         {}
func (FunctionArgument) node() {}
