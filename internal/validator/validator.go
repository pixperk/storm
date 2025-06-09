package validator

import (
	"errors"
	"fmt"
	"regexp"

	"github.com/hashicorp/go-multierror"
	"github.com/pixperk/storm/internal/transform/ir"
)

var validIdentifierRegex = regexp.MustCompile("^[a-zA-Z][a-zA-Z0-9_]*$")
var reservedKeywords = map[string]bool{
	"type":    true,
	"model":   true,
	"select":  true,
	"package": true,
}

// Helper function to check if a field has a specific directive
func hasDirective(field ir.IRField, target ir.Directive) bool {
	for _, d := range field.Directives {
		if d == target {
			return true
		}
	}
	return false
}

// ValidateIR validates an entire DSL file containing multiple models
func ValidateIR(file []ir.IRModel) error {
	modelNames := make(map[string]bool)
	errList := new(multierror.Error)

	for _, model := range file {
		if _, exists := modelNames[model.Name]; exists {
			errList = multierror.Append(errList, fmt.Errorf("duplicate model name: %s", model.Name))
		} else {
			modelNames[model.Name] = true
		}
	}

	for _, model := range file {
		if err := ValidateModel(model, modelNames); err != nil {
			errList = multierror.Append(errList, fmt.Errorf("model %s: %w", model.Name, err))
		}
	}

	return errList.ErrorOrNil()
}

func ValidateModel(model ir.IRModel, modelNames map[string]bool) error {
	errList := new(multierror.Error)

	if model.Name == "" {
		errList = multierror.Append(errList, errors.New("model name cannot be empty"))
	}
	if !validIdentifierRegex.MatchString(model.Name) {
		errList = multierror.Append(errList, fmt.Errorf("invalid model name: %s", model.Name))
	}
	if len(model.Fields) == 0 {
		errList = multierror.Append(errList, errors.New("model must have at least one field"))
	}

	fieldNames := make(map[string]bool)
	idCount := 0

	for _, field := range model.Fields {
		if fieldNames[field.Name] {
			errList = multierror.Append(errList, fmt.Errorf("duplicate field name: %s", field.Name))
		} else {
			fieldNames[field.Name] = true
		}

		if hasDirective(field, ir.DirID) {
			idCount++
		}

		if err := ValidateField(field, model, modelNames); err != nil {
			errList = multierror.Append(errList, fmt.Errorf("field %s: %w", field.Name, err))
		}
	}

	if idCount > 1 {
		errList = multierror.Append(errList, errors.New("model must have less than one @id directive field"))
	}

	return errList.ErrorOrNil()
}

func ValidateField(field ir.IRField, model ir.IRModel, modelNames map[string]bool) error {
	errList := new(multierror.Error)

	if field.Name == "" {
		errList = multierror.Append(errList, errors.New("field name cannot be empty"))
	}
	if !validIdentifierRegex.MatchString(field.Name) {
		errList = multierror.Append(errList, fmt.Errorf("invalid field name: %s", field.Name))
	}
	if reservedKeywords[field.Name] {
		errList = multierror.Append(errList, fmt.Errorf("field name '%s' is reserved", field.Name))
	}

	// Validate directives
	for _, dir := range field.Directives {
		switch dir {
		case ir.DirID, ir.DirAuto, ir.DirUnique, ir.DirHasMany, ir.DirBelongsTo:
			// valid
		default:
			errList = multierror.Append(errList, fmt.Errorf("unknown directive: %s", dir.String()))
		}
	}

	if hasDirective(field, ir.DirID) && field.Type != ir.TypeInt {
		errList = multierror.Append(errList, errors.New("@id field must be of type Int"))
	}
	if hasDirective(field, ir.DirID) && field.IsArray {
		errList = multierror.Append(errList, errors.New("@id field cannot be an array"))
	}
	if hasDirective(field, ir.DirAuto) && !hasDirective(field, ir.DirID) {
		errList = multierror.Append(errList, errors.New("@auto can only be used with @id fields"))
	}
	if hasDirective(field, ir.DirHasMany) && !field.IsArray {
		errList = multierror.Append(errList, errors.New("@hasMany directive can only be used with array fields"))
	}
	if hasDirective(field, ir.DirBelongsTo) && field.IsArray {
		errList = multierror.Append(errList, errors.New("@belongsTo directive cannot be used with array fields"))
	}

	if field.Type == ir.TypeRelation {
		relation := field.Relation
		if !modelNames[relation] {
			errList = multierror.Append(errList, fmt.Errorf("relation field refers to non-existent model: %s", relation))
		}
		if relation == model.Name {
			errList = multierror.Append(errList, fmt.Errorf("self-referential relation is not allowed: %s", field.Name))
		}
		if !field.IsArray && modelNames[relation] && !hasDirective(field, ir.DirBelongsTo) {
			errList = multierror.Append(errList, errors.New("relation field must have @belongsTo directive"))
		}
		if field.IsArray && !hasDirective(field, ir.DirHasMany) {
			errList = multierror.Append(errList, errors.New("array relation field must have @hasMany directive"))
		}
	}

	if hasDirective(field, ir.DirID) && hasDirective(field, ir.DirHasMany) {
		errList = multierror.Append(errList, errors.New("@id field cannot be combined with @hasMany"))
	}
	if hasDirective(field, ir.DirID) && hasDirective(field, ir.DirBelongsTo) {
		errList = multierror.Append(errList, errors.New("@id field cannot be combined with @belongsTo"))
	}

	return errList.ErrorOrNil()
}
