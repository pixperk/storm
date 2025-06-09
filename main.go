package main

import (
	"log"

	"github.com/pixperk/storm/internal/parser"
	"github.com/pixperk/storm/internal/transform/ir"
)

func main() {
	ast, err := parser.ParseDSL("examples/schema.storm")
	if err != nil {
		log.Fatalf("Failed to parse DSL: %v", err)
	}

	parser.DebugPrint(ast)
	irVar := ir.ToIR(ast)
	ir.PrintIR(irVar)
}
