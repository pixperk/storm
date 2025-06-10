package ir

import (
	"strings"

	"github.com/pixperk/storm/internal/types/directive"
	"github.com/pixperk/storm/internal/types/field"
)

// MapFieldType maps a string representation of a field type to a FieldType
func MapFieldType(s string) field.FieldKind {
	s = strings.ToLower(strings.TrimSpace(s))

	switch s {
	// Integer types
	case "int", "integer":
		return field.KindInt
	case "bigint":
		return field.KindBigInt // Floating point types
	case "float", "double", "real":
		return field.KindFloat

	// String types
	case "string", "varchar":
		return field.KindString

	case "text":
		return field.KindText
	case "char":
		return field.KindChar

	// Boolean type
	case "bool", "boolean":
		return field.KindBoolean

	// Date and time types
	case "datetime":
		return field.KindDateTime
	case "date":
		return field.KindDate
	case "time":
		return field.KindTime
	case "timestamp":
		return field.KindTimestamp

	// Binary data
	case "binary", "blob":
		return field.KindBinary

	// JSON
	case "json":
		return field.KindJSON
	case "point":
		return field.KindPoint

	default:
		// Handle as a custom type or return a default
		return field.KindCustom
	}
}

// MapDirective maps a string representation of a directive to a Directive
func MapDirective(name string, args []string) *directive.Directive {
	name = strings.ToLower(strings.TrimSpace(name))

	var kind directive.DirectiveKind

	switch name {
	case "id":
		kind = directive.DirID
	case "auto":
		kind = directive.DirAuto
	case "default":
		kind = directive.DirDefault
	case "unique":
		kind = directive.DirUnique
	case "nullable":
		kind = directive.DirNullable
	case "hasmany":
		kind = directive.DirHasMany
	case "belongsto":
		kind = directive.DirBelongsTo
	case "hasone":
		kind = directive.DirHasOne
	case "index":
		kind = directive.DirIndex
	case "enum":
		kind = directive.DirEnum
	case "updatedat":
		kind = directive.DirUpdatedAt
	case "createdat":
		kind = directive.DirCreatedAt
	// New directive types
	case "length":
		kind = directive.DirLength
	case "min":
		kind = directive.DirMin
	case "max":
		kind = directive.DirMax
	case "precision":
		kind = directive.DirPrecision
	case "defaultnow":
		kind = directive.DirDefaultNow
	case "map":
		kind = directive.DirMap
	case "relation":
		kind = directive.DirRelation
	default:
		// Handle unknown directive - could return error instead
		return nil
	}

	dir := directive.NewDirective(kind, args)
	// Add all arguments to the directive

	return dir
}
