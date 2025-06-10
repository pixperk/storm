package validator

import (
	"fmt"
	"strconv"

	"github.com/hashicorp/go-multierror"
	"github.com/pixperk/storm/internal/types/directive"
)

// validateDirectiveArgs checks if the directive arguments are valid
func validateDirectiveArgs(dir directive.Directive) error {
	errList := new(multierror.Error)

	switch dir.Kind {
	case directive.DirLength:
		// @length requires exactly one integer argument
		if len(dir.Args) != 1 {
			errList = multierror.Append(errList, fmt.Errorf("@length directive requires exactly one integer argument"))
		} else if _, err := strconv.Atoi(dir.Args[0]); err != nil {
			errList = multierror.Append(errList, fmt.Errorf("@length directive argument must be an integer"))
		}

	case directive.DirPrecision:
		// @precision requires exactly two integer arguments
		if len(dir.Args) != 2 {
			errList = multierror.Append(errList, fmt.Errorf("@precision directive requires exactly two integer arguments"))
		} else {
			p, err1 := strconv.Atoi(dir.Args[0])
			s, err2 := strconv.Atoi(dir.Args[1])

			if err1 != nil || err2 != nil {
				errList = multierror.Append(errList, fmt.Errorf("@precision directive arguments must be integers"))
			} else if p <= 0 {
				errList = multierror.Append(errList, fmt.Errorf("@precision first argument (precision) must be positive"))
			} else if s < 0 || s > p {
				errList = multierror.Append(errList, fmt.Errorf("@precision second argument (scale) must be between 0 and precision"))
			}
		}

	case directive.DirMin, directive.DirMax:
		// @min and @max require exactly one numeric argument
		if len(dir.Args) != 1 {
			errList = multierror.Append(errList, fmt.Errorf("@%s directive requires exactly one numeric argument", dir.Kind.String()))
		} else {
			_, err := strconv.ParseFloat(dir.Args[0], 64)
			if err != nil {
				errList = multierror.Append(errList, fmt.Errorf("@%s directive argument must be a number", dir.Kind.String()))
			}
		}

	case directive.DirDefault:
		// @default requires exactly one argument
		if len(dir.Args) != 1 {
			errList = multierror.Append(errList, fmt.Errorf("@default directive requires exactly one argument"))
		}

	case directive.DirHasMany, directive.DirBelongsTo, directive.DirID, directive.DirAuto,
		directive.DirUnique, directive.DirIndex, directive.DirUpdatedAt, directive.DirCreatedAt:
		// These directives don't require arguments
		if len(dir.Args) > 0 {
			errList = multierror.Append(errList, fmt.Errorf("@%s directive does not accept arguments", dir.Kind.String()))
		}
	}

	return errList.ErrorOrNil()
}
