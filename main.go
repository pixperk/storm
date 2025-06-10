package main

import (
	"log"

	"github.com/pixperk/storm/internal/parser"
	"github.com/pixperk/storm/internal/transform/ir"
	"github.com/pixperk/storm/internal/validator"
)

func main() {
	ast, err := parser.ParseDSL("examples/schema.storm")
	if err != nil {
		log.Fatalf("Failed to parse DSL: %v", err)
	}

	irVar, err := ir.ToIR(ast)
	if err != nil {
		log.Fatalf("Failed to transform to IR: %v", err)
	}
	err = validator.ValidateIR(irVar)
	if err != nil {
		log.Fatalf("Validation error: %v", err)
	}
	ir.PrintIR(irVar)
}
