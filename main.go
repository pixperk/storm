package main

import (
	"log"

	"github.com/pixperk/storm/parser"
)

func main() {
	ast, err := parser.ParseDSL("examples/schema.storm")
	if err != nil {
		log.Fatalf("Failed to parse DSL: %v", err)
	}

	parser.DebugPrint(ast)
}
