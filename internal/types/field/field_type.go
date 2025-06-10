package field

import "github.com/pixperk/storm/internal/types/directive"

// FieldType represents a complete field type with options
type FieldType struct {
	Kind       FieldKind
	Directives []directive.Directive
}

// NewFieldType creates a new FieldType with the given kind
func NewFieldType(kind FieldKind, dirs []directive.Directive) *FieldType {
	return &FieldType{
		Kind:       kind,
		Directives: dirs,
	}
}
