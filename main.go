package main

import (
	"log"

	"github.com/pixperk/storm/parser"
	"github.com/pixperk/storm/transform"
)

func main() {
	ast, err := parser.ParseDSL("examples/schema.storm")
	if err != nil {
		log.Fatalf("Failed to parse DSL: %v", err)
	}

	parser.DebugPrint(ast)
	ir := transform.ToIR(ast)
	transform.PrintIR(ir)
}
