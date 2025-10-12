package main

import (
	"fmt"
	"os"

	"github.com/dywoq/minigo/scanner"
)

func main() {
	f, err := os.Open("main.dl")
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
		if token != nil {
			fmt.Printf("%s %s\n", token.Literal, token.Kind)
		}
	}
}
