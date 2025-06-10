package validator

import (
	"errors"
	"fmt"

	"github.com/hashicorp/go-multierror"
	"github.com/pixperk/storm/internal/transform/ir"
	"github.com/pixperk/storm/internal/types/directive"
)

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

		if hasDirective(field, directive.DirID) {
			idCount++
		}

		if err := ValidateField(field, model, modelNames); err != nil {
			errList = multierror.Append(errList, fmt.Errorf("field %s: %w", field.Name, err))
		}
	}

	if idCount > 1 {
		errList = multierror.Append(errList, errors.New("model must have at most one @id directive field"))
	}

	return errList.ErrorOrNil()
}
