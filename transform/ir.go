package transform

import (
	"fmt"

	"github.com/pixperk/storm/parser"
)

type IRModel struct {
	Name   string
	Fields []IRField
}

type IRField struct {
	Name       string
	Type       string //TODO : make an enum for types
	IsArray    bool
	Directives map[string]bool //TODO : make an enum for directives
}

func ToIR(ast *parser.DSLFile) []IRModel {
	var ir []IRModel

	for _, m := range ast.Models {
		model := IRModel{
			Name: m.Name,
		}

		for _, f := range m.Fields {
			directiveMap := map[string]bool{}
			for _, d := range f.Directives {
				directiveMap[d] = true
			}

			model.Fields = append(model.Fields, IRField{
				Name:       f.Name,
				Type:       f.Type.Name,
				IsArray:    f.Type.IsArray,
				Directives: directiveMap,
			})
		}

		ir = append(ir, model)
	}
	return ir
}

// DebugPrintIR prints the IR representation of the DSL file.
func PrintIR(models []IRModel) {
	for _, m := range models {
		fmt.Println("Model:", m.Name)
		for _, f := range m.Fields {
			arr := ""
			if f.IsArray {
				arr = "[]"
			}
			fmt.Printf("  - %s: %s%s", f.Name, f.Type, arr)
			if len(f.Directives) > 0 {
				fmt.Printf(" (directives: %v)", f.Directives)
			}
			fmt.Println()
		}
	}
}
