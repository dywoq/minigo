package main

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/dywoq/minigo/pkg/parser"
	"github.com/dywoq/minigo/pkg/scanner"
)

func main() {
	f, err := os.Open("./main.dl")
	if err != nil {
		panic(err)
	}
	s, err := scanner.NewDebug(f, os.Stdout)
	if err != nil {
		panic(err)
	}

	tokens, err := s.Scan()
	if err != nil {
		panic(err)
	}

	for _, token := range tokens {
		fmt.Printf("%s %s %v\n", token.Literal, token.Kind, token.Position)
	}

	p, err := parser.NewDebug(tokens, os.Stdout)
	if err != nil {
		panic(err)
	}

	file, err := p.Parse()
	if err != nil {
		panic(err)
	}

	content, err := json.MarshalIndent(file, "| node |", "  ")
	if err != nil {
		panic(err)
	}
	fmt.Printf("ast: %v\n", string(content))
}
