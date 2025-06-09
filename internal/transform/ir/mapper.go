package ir

import (
	"fmt"
	"strings"
)

func MapFieldType(s string) FieldType {
	switch s {
	case "Int":
		return TypeInt
	case "String":
		return TypeString
	case "Boolean":
		return TypeBoolean
	case "Float":
		return TypeFloat
	default:
		return TypeRelation // Fallback for custom types
	}
}

func MapDirectives(raws []string) ([]Directive, error) {
	var result []Directive
	for _, r := range raws {
		r = strings.TrimSpace(r)
		switch r {
		case "id":
			result = append(result, DirID)
		case "auto":
			result = append(result, DirAuto)
		case "unique":
			result = append(result, DirUnique)
		case "hasMany":
			result = append(result, DirHasMany)
		case "belongsTo":
			result = append(result, DirBelongsTo)
		default:
			return nil, fmt.Errorf("unknown directive: %s", r)
		}
	}
	return result, nil
}
