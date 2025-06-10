package field

import (
	"strconv"

	"github.com/pixperk/storm/internal/types/directive"
)

// FieldType represents a complete field type with options
type FieldType struct {
	Kind       FieldKind
	ModelName  string // Used for custom/relation types to store the model name
	Directives []directive.Directive
}

// NewFieldType creates a new FieldType with the given kind
func NewFieldType(kind FieldKind, modelName string, dirs []directive.Directive) *FieldType {
	return &FieldType{
		Kind:       kind,
		ModelName:  modelName,
		Directives: dirs,
	}
}

func (ft *FieldType) GetDirective(kind directive.DirectiveKind) []string {
	for _, dir := range ft.Directives {
		if dir.Kind == kind {
			return dir.Args
		}
	}
	return nil
}

// GetLength returns the length of the field type, defaulting to 255 if not specified
func returnLength(ft FieldType) string {
	if lengthArr := ft.GetDirective(directive.DirLength); len(lengthArr) > 0 {
		if l, err := strconv.Atoi(lengthArr[0]); err == nil && l > 0 {
			return lengthArr[0]
		}
	}
	return "255"
}

// GetPrecisionScale returns the precision and scale values from the precision directive
func (ft FieldType) GetPrecisionScale() (string, string) {
	if precisionArr := ft.GetDirective(directive.DirPrecision); len(precisionArr) >= 2 {
		precision := precisionArr[0]
		scale := precisionArr[1]
		if p, err := strconv.Atoi(precision); err == nil && p > 0 {
			if s, err := strconv.Atoi(scale); err == nil && s >= 0 && s <= p {
				return precision, scale
			}
		}
	}
	return "", ""
}

// String returns a string representation of the field type
func (ft FieldType) String() string {
	if ft.Kind == KindCustom && ft.ModelName != "" {
		return ft.ModelName
	}
	return ft.Kind.String()
}
