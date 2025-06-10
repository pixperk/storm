package validator

import (
	"fmt"
	"regexp"

	"github.com/hashicorp/go-multierror"
	"github.com/pixperk/storm/internal/transform/ir"
	"github.com/pixperk/storm/internal/types/directive"
)

var validIdentifierRegex = regexp.MustCompile("^[a-zA-Z][a-zA-Z0-9_]*$")
var reservedKeywords = map[string]bool{
	"type":    true,
	"model":   true,
	"select":  true,
	"package": true,
}

// Helper function to check if a field has a specific directive
func hasDirective(field ir.IRField, targetKind directive.DirectiveKind) bool {
	for _, dir := range field.Type.Directives {
		if dir.Kind == targetKind {
			return true
		}
	}
	return false
}

// ValidateIR validates an entire IR file
func ValidateIR(irData *ir.IR) error {
	errList := new(multierror.Error)

	// Validate database configuration
	if err := validateDatabaseConfig(irData); err != nil {
		errList = multierror.Append(errList, err)
	}

	// Validate models
	modelNames := make(map[string]bool)
	for _, model := range irData.Models {
		if _, exists := modelNames[model.Name]; exists {
			errList = multierror.Append(errList, fmt.Errorf("duplicate model name: %s", model.Name))
		} else {
			modelNames[model.Name] = true
		}
	}

	for _, model := range irData.Models {
		if err := ValidateModel(model, modelNames); err != nil {
			errList = multierror.Append(errList, fmt.Errorf("model %s: %w", model.Name, err))
		}
	}

	// Validate relational consistency
	if err := validateRelationalConsistency(irData.Models); err != nil {
		errList = multierror.Append(errList, err)
	}

	return errList.ErrorOrNil()
}

// validateRelationalConsistency checks that all relations are properly defined
func validateRelationalConsistency(models []ir.IRModel) error {
	errList := new(multierror.Error)
	modelMap := make(map[string]ir.IRModel)

	// Create a map of model names to models for easy lookup
	for _, model := range models {
		modelMap[model.Name] = model
	}

	// Check each model's relations
	for _, model := range models {
		for _, field := range model.Fields {
			// Check for @hasMany directive
			if hasDirective(field, directive.DirHasMany) {
				// The field must be an array
				if !field.IsArray {
					errList = multierror.Append(errList, fmt.Errorf("model %s: field %s with @hasMany must be an array", model.Name, field.Name))
					continue
				} // The field type should be a valid model
				relatedModelName := field.Type.String()
				relatedModel, exists := modelMap[relatedModelName]
				if !exists {
					errList = multierror.Append(errList, fmt.Errorf("model %s: field %s references non-existent model %s", model.Name, field.Name, relatedModelName))
					continue
				}

				// Find a corresponding field in the related model with @belongsTo pointing back to this model
				foundBackReference := false
				for _, relatedField := range relatedModel.Fields {
					if hasDirective(relatedField, directive.DirBelongsTo) &&
						!relatedField.IsArray &&
						relatedField.Type.String() == model.Name {
						foundBackReference = true
						break
					}
				}

				if !foundBackReference {
					errList = multierror.Append(errList, fmt.Errorf("model %s: field %s has @hasMany but no corresponding @belongsTo in model %s", model.Name, field.Name, relatedModelName))
				}
			}

			// Check for @belongsTo directive
			if hasDirective(field, directive.DirBelongsTo) {
				// The field must not be an array
				if field.IsArray {
					errList = multierror.Append(errList, fmt.Errorf("model %s: field %s with @belongsTo cannot be an array", model.Name, field.Name))
					continue
				} // The field type should be a valid model
				relatedModelName := field.Type.String()
				relatedModel, exists := modelMap[relatedModelName]
				if !exists {
					errList = multierror.Append(errList, fmt.Errorf("model %s: field %s references non-existent model %s", model.Name, field.Name, relatedModelName))
					continue
				}

				// Find a corresponding field in the related model with @hasMany pointing back to this model
				foundBackReference := false
				for _, relatedField := range relatedModel.Fields {
					if hasDirective(relatedField, directive.DirHasMany) &&
						relatedField.IsArray &&
						relatedField.Type.String() == model.Name {
						foundBackReference = true
						break
					}
				}

				if !foundBackReference {
					errList = multierror.Append(errList, fmt.Errorf("model %s: field %s has @belongsTo but no corresponding @hasMany in model %s", model.Name, field.Name, relatedModelName))
				}
			}
		}
	}

	return errList.ErrorOrNil()
}
