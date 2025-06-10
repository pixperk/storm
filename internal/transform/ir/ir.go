package ir

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/pixperk/storm/internal/parser"
	"github.com/pixperk/storm/internal/types/directive"
	"github.com/pixperk/storm/internal/types/field"
)

type IRModel struct {
	Name   string
	Fields []IRField
}

type IRField struct {
	Name    string
	Type    field.FieldType
	IsArray bool
}

type IR struct {
	DatabaseDriver string
	DatabaseURL    string
	Models         []IRModel
}

func ToIR(ast *parser.DSLFile) (*IR, error) {
	ir := &IR{
		DatabaseDriver: ast.DatabaseDriver,
		DatabaseURL:    ast.DatabaseURL,
		Models:         make([]IRModel, 0, len(ast.Models)),
	}

	for _, m := range ast.Models {
		model := IRModel{
			Name:   m.Name,
			Fields: make([]IRField, 0, len(m.Fields)),
		}

		for _, f := range m.Fields {
			fieldKind := MapFieldType(f.Type.Name)

			directives := make([]directive.Directive, 0, len(f.Directives))
			for _, rawDir := range f.Directives {
				args := make([]string, 0, len(rawDir.Args))
				for _, arg := range rawDir.Args {
					switch {
					case arg.String != nil:
						args = append(args, *arg.String)
					case arg.Ident != nil:
						args = append(args, *arg.Ident)
					case arg.Int != nil:
						args = append(args, strconv.Itoa(*arg.Int))
					case arg.Float != nil:
						args = append(args, fmt.Sprintf("%f", *arg.Float))
					}
				}

				if dirObj := MapDirective(rawDir.Name, args); dirObj != nil {
					directives = append(directives, *dirObj)
				}
			}

			field := field.NewFieldType(fieldKind, directives)
			model.Fields = append(model.Fields, IRField{
				Name:    f.Name,
				Type:    *field,
				IsArray: f.Type.IsArray,
			})
		}

		ir.Models = append(ir.Models, model)
	}

	return ir, nil
}

// PrintIR prints the IR representation of the DSL file with database type information.
func PrintIR(ir *IR) {
	fmt.Println("=============== IR Models ===============")
	fmt.Printf("Database Driver: %s\n", ir.DatabaseDriver)
	fmt.Printf("Database URL: %s\n", ir.DatabaseURL)

	for _, m := range ir.Models {
		fmt.Printf("\n┌─── Model: %s ───┐\n", m.Name)

		if len(m.Fields) == 0 {
			fmt.Println("│  No fields defined                │")
		} else {
			for _, f := range m.Fields {
				// Display field name and type
				typeStr := f.Type.Kind.String()
				if f.IsArray {
					typeStr += "[]"
				}

				fmt.Printf("│  %-15s : %-10s │\n", f.Name, typeStr)

				// Display database types
				fmt.Printf("│    ├─ MySQL    : %-15s │\n", f.Type.MySQLType())
				fmt.Printf("│    ├─ Postgres : %-15s │\n", f.Type.PostgresType())
				fmt.Printf("│    └─ SQLite   : %-15s │\n", f.Type.SQLiteType())

				// Display directives if any
				if len(f.Type.Directives) > 0 {
					fmt.Println("│    ┌─ Directives:            │")
					for i, dir := range f.Type.Directives {
						prefix := "│    ├─"
						if i == len(f.Type.Directives)-1 {
							prefix = "│    └─"
						}

						if len(dir.Args) > 0 {
							fmt.Printf("%s @%-10s(%s) │\n", prefix, dir.Kind.String(), strings.Join(dir.Args, ", "))
						} else {
							fmt.Printf("%s @%-10s       │\n", prefix, dir.Kind.String())
						}
					}
				}

				// Add separator between fields
				fmt.Println("│                               │")
			}
		}

		fmt.Println("└───────────────────────────────┘")
	}

	fmt.Println("\n=========== End of IR Models ===========")
}
