package validator

import (
	"errors"
	"fmt"
	"regexp"
	"strconv"
	"strings"

	"github.com/hashicorp/go-multierror"
	"github.com/pixperk/storm/internal/transform/ir"
	"github.com/pixperk/storm/internal/types/directive"
	fld "github.com/pixperk/storm/internal/types/field"
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

	// Validate field type
	if err := validateFieldType(field); err != nil {
		errList = multierror.Append(errList, err)
	}

	// Validate directives
	if err := validateDirectives(field); err != nil {
		errList = multierror.Append(errList, err)
	}

	// Validate relation fields
	relatedModelName := field.Type.Kind.String()
	if modelNames[relatedModelName] && relatedModelName != model.Name {
		// This is a relation field pointing to another model
		if !field.IsArray && !hasDirective(field, directive.DirBelongsTo) {
			errList = multierror.Append(errList, errors.New("relation field must have @belongsTo directive"))
		}
		if field.IsArray && !hasDirective(field, directive.DirHasMany) {
			errList = multierror.Append(errList, errors.New("array relation field must have @hasMany directive"))
		}
	}

	// Validate directive combinations
	if hasDirective(field, directive.DirID) {
		// ID field specific validations
		if field.Type.Kind != fld.KindInt {
			errList = multierror.Append(errList, errors.New("@id field must be of type Int"))
		}
		if field.IsArray {
			errList = multierror.Append(errList, errors.New("@id field cannot be an array"))
		}
		if hasDirective(field, directive.DirHasMany) {
			errList = multierror.Append(errList, errors.New("@id field cannot be combined with @hasMany"))
		}
		if hasDirective(field, directive.DirBelongsTo) {
			errList = multierror.Append(errList, errors.New("@id field cannot be combined with @belongsTo"))
		}
	}

	if hasDirective(field, directive.DirAuto) && !hasDirective(field, directive.DirID) {
		errList = multierror.Append(errList, errors.New("@auto can only be used with @id fields"))
	}

	return errList.ErrorOrNil()
}

// validateDatabaseConfig validates the database configuration
func validateDatabaseConfig(irData *ir.IR) error {
	errList := new(multierror.Error)

	// Check if database driver is provided
	if irData.DatabaseDriver == "" {
		errList = multierror.Append(errList, errors.New("database driver is required"))
	} else {
		// Validate supported database drivers
		switch irData.DatabaseDriver {
		case "mysql", "postgres", "postgresql", "sqlite", "sqlite3":
			// Valid database drivers
		default:
			errList = multierror.Append(errList, fmt.Errorf("unsupported database driver: %s", irData.DatabaseDriver))
		}
	}

	// Check if database URL is provided
	if irData.DatabaseURL == "" {
		errList = multierror.Append(errList, errors.New("database URL is required"))
	} else {
		// Basic URL validation (could be more sophisticated)
		if !strings.Contains(irData.DatabaseURL, "://") {
			errList = multierror.Append(errList, fmt.Errorf("invalid database URL format: %s", irData.DatabaseURL))
		}
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
				}

				// The field type should be a valid model
				relatedModelName := field.Type.Kind.String()
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
						relatedField.Type.Kind.String() == model.Name {
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
				}

				// The field type should be a valid model
				relatedModelName := field.Type.Kind.String()
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
						relatedField.Type.Kind.String() == model.Name {
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

// validateFieldType checks if the field type is valid
func validateFieldType(field ir.IRField) error {
	errList := new(multierror.Error)

	// Check if the field type is a known type
	fieldKind := field.Type.Kind
	switch fieldKind {
	case fld.KindInt, fld.KindFloat, fld.KindDecimal, fld.KindBigInt,
		fld.KindString, fld.KindText, fld.KindChar,
		fld.KindBoolean,
		fld.KindDateTime, fld.KindDate, fld.KindTime, fld.KindTimestamp,
		fld.KindBinary,
		fld.KindJSON, fld.KindUUID, fld.KindCUID, fld.KindPoint:
		// Valid built-in types
	case fld.KindCustom:
		errList = multierror.Append(errList, fmt.Errorf("custom type requires additional specification"))
	default:
		// Check if it's a model type (for relations)
		if fieldKind.String() == "" {
			errList = multierror.Append(errList, fmt.Errorf("unknown field type: %v", fieldKind))
		}
	}

	return errList.ErrorOrNil()
}

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
		if err := validateDirectiveArgs(field, dir); err != nil {
			errList = multierror.Append(errList, err)
		}

		// Validate directive compatibility with field type
		if err := validateDirectiveTypeCompatibility(field, dir); err != nil {
			errList = multierror.Append(errList, err)
		}
	}

	return errList.ErrorOrNil()
}

// validateDirectiveArgs checks if the directive arguments are valid
func validateDirectiveArgs(field ir.IRField, dir directive.Directive) error {
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
