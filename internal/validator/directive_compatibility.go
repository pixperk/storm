package validator

import (
	"fmt"

	"github.com/hashicorp/go-multierror"
	"github.com/pixperk/storm/internal/transform/ir"
	"github.com/pixperk/storm/internal/types/directive"
	fld "github.com/pixperk/storm/internal/types/field"
)

// validateDirectiveTypeCompatibility checks if the directive is compatible with the field type
func validateDirectiveTypeCompatibility(field ir.IRField, dir directive.Directive) error {
	errList := new(multierror.Error)
	fieldKind := field.Type.Kind

	switch dir.Kind {
	case directive.DirLength:
		// @length can only be used with string types
		if fieldKind != fld.KindString && fieldKind != fld.KindChar && fieldKind != fld.KindText {
			errList = multierror.Append(errList, fmt.Errorf("@length directive can only be used with string types"))
		}

	case directive.DirPrecision:
		// @precision can only be used with decimal or float types
		if fieldKind != fld.KindDecimal && fieldKind != fld.KindFloat {
			errList = multierror.Append(errList, fmt.Errorf("@precision directive can only be used with decimal or float types"))
		}

	case directive.DirMin, directive.DirMax:
		// @min and @max can only be used with numeric types
		if fieldKind != fld.KindInt && fieldKind != fld.KindFloat &&
			fieldKind != fld.KindDecimal && fieldKind != fld.KindBigInt {
			errList = multierror.Append(errList, fmt.Errorf("@%s directive can only be used with numeric types", dir.Kind.String()))
		}

	case directive.DirHasMany:
		// @hasMany can only be used with array fields
		if !field.IsArray {
			errList = multierror.Append(errList, fmt.Errorf("@hasMany directive can only be used with array fields"))
		}

	case directive.DirBelongsTo:
		// @belongsTo cannot be used with array fields
		if field.IsArray {
			errList = multierror.Append(errList, fmt.Errorf("@belongsTo directive cannot be used with array fields"))
		}

	case directive.DirID:
		// @id is typically used with Int fields
		if fieldKind != fld.KindInt {
			errList = multierror.Append(errList, fmt.Errorf("@id directive is recommended to be used with Int fields"))
		}

	case directive.DirDefaultNow:
		// @defaultNow can only be used with date/time types
		if fieldKind != fld.KindDateTime && fieldKind != fld.KindDate &&
			fieldKind != fld.KindTime && fieldKind != fld.KindTimestamp {
			errList = multierror.Append(errList, fmt.Errorf("@defaultNow directive can only be used with date/time types"))
		}

	case directive.DirUpdatedAt, directive.DirCreatedAt:
		// @updatedAt and @createdAt can only be used with date/time types
		if fieldKind != fld.KindDateTime && fieldKind != fld.KindTimestamp {
			errList = multierror.Append(errList, fmt.Errorf("@%s directive can only be used with DateTime or Timestamp types", dir.Kind.String()))
		}
	}

	return errList.ErrorOrNil()
}
