package ir

import (
	"fmt"

	"github.com/pixperk/storm/internal/parser"
)

type IRModel struct {
	Name   string
	Fields []IRField
}

type IRField struct {
	Name       string
	Type       FieldType
	Relation   string
	IsArray    bool
	Directives []Directive
}

func ToIR(ast *parser.DSLFile) ([]IRModel, error) {
	var ir []IRModel
	relation := ""

	for _, m := range ast.Models {
		model := IRModel{
			Name:   m.Name,
			Fields: []IRField{},
		}

		for _, f := range m.Fields {
			fieldType := MapFieldType(f.Type.Name)
			if fieldType == TypeRelation {
				relation = f.Type.Name // Store the related model name
			} else {
				relation = ""
			}
			directives, err := MapDirectives(f.Directives)
			if err != nil {
				return nil, err
			}

			model.Fields = append(model.Fields, IRField{
				Name:       f.Name,
				Type:       fieldType,
				IsArray:    f.Type.IsArray,
				Directives: directives,
				Relation:   relation,
			})
		}

		ir = append(ir, model)
	}
	return ir, nil
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
