package validator

import (
	"fmt"

	"github.com/hashicorp/go-multierror"
	"github.com/pixperk/storm/internal/transform/ir"
	"github.com/pixperk/storm/internal/types/directive"
)

// validateDirectives checks if the field's directives are valid for its type
func validateDirectives(field ir.IRField) error {
	errList := new(multierror.Error)

	for _, dir := range field.Type.Directives {
		// First validate the directive is known
		switch dir.Kind {
		case directive.DirID, directive.DirAuto, directive.DirDefault, directive.DirUnique,
			directive.DirNullable, directive.DirHasMany, directive.DirBelongsTo, directive.DirHasOne,
			directive.DirIndex, directive.DirEnum, directive.DirUpdatedAt, directive.DirCreatedAt,
			directive.DirLength, directive.DirMin, directive.DirMax, directive.DirPrecision,
			directive.DirDefaultNow, directive.DirMap, directive.DirRelation:
			// Valid directive kind
		default:
			errList = multierror.Append(errList, fmt.Errorf("unknown directive: %s", dir.Kind.String()))
			continue
		}

		// Then validate directive arguments
		if err := validateDirectiveArgs(dir); err != nil {
			errList = multierror.Append(errList, err)
		}

		// Validate directive compatibility with field type
		if err := validateDirectiveTypeCompatibility(field, dir); err != nil {
			errList = multierror.Append(errList, err)
		}
	}

	return errList.ErrorOrNil()
}
