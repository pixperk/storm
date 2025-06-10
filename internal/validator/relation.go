package validator

import (
	"fmt"

	"github.com/hashicorp/go-multierror"
	"github.com/pixperk/storm/internal/transform/ir"
	"github.com/pixperk/storm/internal/types/directive"
)

func warn(message string, args ...interface{}) {
	// This function can be used to log warnings during validation
	// For now, we just print to stdout, but it can be replaced with a proper logging mechanism
	fmt.Printf("Warning: "+message+"\n", args...)
}

func validateRelationalConsistency(models []ir.IRModel) error {
	errList := new(multierror.Error)
	modelMap := make(map[string]ir.IRModel)

	// Build model lookup
	for _, model := range models {
		modelMap[model.Name] = model
	}

	// Validate each model's fields
	for _, model := range models {
		for _, field := range model.Fields {
			validateFieldRelations(model, field, modelMap, errList)
		}
	}

	return errList.ErrorOrNil()
}

func validateFieldRelations(model ir.IRModel, field ir.IRField, modelMap map[string]ir.IRModel, errList *multierror.Error) {
	switch {
	case hasDirective(field, directive.DirHasMany):
		validateHasMany(model, field, modelMap, errList)

	case hasDirective(field, directive.DirBelongsTo):
		validateBelongsTo(model, field, modelMap, errList)

	case hasDirective(field, directive.DirHasOne):
		validateHasOne(model, field, modelMap, errList)
	}
}

func validateHasMany(model ir.IRModel, field ir.IRField, modelMap map[string]ir.IRModel, errList *multierror.Error) {
	if !field.IsArray {
		// This is a real error, not just a warning
		_ = multierror.Append(errList, fmt.Errorf(
			"model %s: field %s with @hasMany must be an array", model.Name, field.Name))
		return
	}

	relatedModel, ok := modelMap[field.Type.String()]
	if !ok {
		// This is a real error, not just a warning
		_ = multierror.Append(errList, fmt.Errorf(
			"model %s: field %s references non-existent model %s", model.Name, field.Name, field.Type.String()))
		return
	}
	if !hasBackReference(relatedModel, model.Name, directive.DirBelongsTo, false) &&
		!hasBackReference(relatedModel, model.Name, directive.DirHasMany, true) {
		// Issue a warning for unidirectional @hasMany relationships instead of an error
		warn("model %s: field %s has @hasMany but no corresponding @belongsTo or @hasMany in model %s (unidirectional relation)",
			model.Name, field.Name, relatedModel.Name)
	}
}

func validateBelongsTo(model ir.IRModel, field ir.IRField, modelMap map[string]ir.IRModel, errList *multierror.Error) {
	if field.IsArray {
		_ = multierror.Append(errList, fmt.Errorf(
			"model %s: field %s with @belongsTo cannot be an array", model.Name, field.Name))
		return
	}

	relatedModel, ok := modelMap[field.Type.String()]
	if !ok {
		_ = multierror.Append(errList, fmt.Errorf(
			"model %s: field %s references non-existent model %s", model.Name, field.Name, field.Type.String()))
		return
	}

	// Check for circular @belongsTo relations (belongsTo in both directions)
	if hasBelongsToBackReference(relatedModel, model.Name) {
		_ = multierror.Append(errList, fmt.Errorf(
			"circular @belongsTo relation detected: both model %s and model %s have @belongsTo pointing to each other",
			model.Name, relatedModel.Name))
		return
	}

	if !hasBackReference(relatedModel, model.Name, directive.DirHasMany, true) &&
		!hasBackReference(relatedModel, model.Name, directive.DirHasOne, false) {
		// Issue a warning for unidirectional @belongsTo relationships instead of an error
		warn("model %s: field %s has @belongsTo but no corresponding @hasMany or @hasOne in model %s (unidirectional relation)",
			model.Name, field.Name, relatedModel.Name)
	}
}

func validateHasOne(model ir.IRModel, field ir.IRField, modelMap map[string]ir.IRModel, errList *multierror.Error) {
	if field.IsArray {
		_ = multierror.Append(errList, fmt.Errorf(
			"model %s: field %s with @hasOne cannot be an array", model.Name, field.Name))
		return
	}

	relatedModel, ok := modelMap[field.Type.String()]
	if !ok {
		_ = multierror.Append(errList, fmt.Errorf(
			"model %s: field %s references non-existent model %s", model.Name, field.Name, field.Type.String()))
		return
	}
	if !hasBackReference(relatedModel, model.Name, directive.DirBelongsTo, false) {
		// Issue a warning for unidirectional @hasOne relationships instead of an error
		warn("model %s: field %s has @hasOne but no corresponding @belongsTo in model %s (unidirectional relation)",
			model.Name, field.Name, relatedModel.Name)
	}
}

func hasBackReference(model ir.IRModel, targetModelName string, directiveType directive.DirectiveKind, shouldBeArray bool) bool {
	for _, field := range model.Fields {
		if hasDirective(field, directiveType) &&
			field.IsArray == shouldBeArray &&
			field.Type.String() == targetModelName {
			return true
		}
	}
	return false
}

// hasBelongsToBackReference checks if a model has a field with @belongsTo directive pointing to the target model
func hasBelongsToBackReference(model ir.IRModel, targetModelName string) bool {
	return hasBackReference(model, targetModelName, directive.DirBelongsTo, false)
}
