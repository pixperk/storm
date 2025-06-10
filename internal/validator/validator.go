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
