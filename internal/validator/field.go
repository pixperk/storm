package validator

import (
	"errors"
	"fmt"

	"github.com/hashicorp/go-multierror"
	"github.com/pixperk/storm/internal/transform/ir"
	"github.com/pixperk/storm/internal/types/directive"
	fld "github.com/pixperk/storm/internal/types/field"
)

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
	relatedModelName := field.Type.String()
	if modelNames[relatedModelName] && relatedModelName != model.Name {
		// This is a relation field pointing to another model
		if !field.IsArray && !(hasDirective(field, directive.DirBelongsTo) || hasDirective(field, directive.DirHasOne)) {
			errList = multierror.Append(errList, errors.New("relation field must have @belongsTo or @hasOne directive"))
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
		fld.KindJSON, fld.KindUUID, fld.KindCUID, fld.KindPoint, fld.KindCustom:

	default:
		// Check if it's a model type (for relations)
		if fieldKind.String() == "" {
			errList = multierror.Append(errList, fmt.Errorf("unknown field type: %v", fieldKind))
		}
	}

	return errList.ErrorOrNil()
}
