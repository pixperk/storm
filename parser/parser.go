package parser

import (
	"fmt"
	"os"
)

func ParseDSL(path string) (*DSLFile, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}
	return Parser.ParseString("", string(data))
}

func DebugPrint(ast *DSLFile) {
	for _, model := range ast.Models {
		fmt.Printf("Model: %s\n", model.Name)
		for _, field := range model.Fields {
			typeStr := field.Type.Name
			if field.Type.IsArray {
				typeStr += "[]"
			}
			fmt.Printf("  - %s: %s", field.Name, typeStr)

			if len(field.Directives) > 0 {
				fmt.Printf(" (directives: %v)", field.Directives)
			}
			fmt.Println()
		}
	}
}
