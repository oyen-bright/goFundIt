package utils

import "github.com/go-playground/validator/v10"

func ExtractValidationErrors(err error) []map[string]interface{} {
	var errors []map[string]interface{}

	if validationErrors, ok := err.(validator.ValidationErrors); ok {
		for _, validationError := range validationErrors {
			errors = append(errors, map[string]interface{}{
				"field":   validationError.StructField(),
				"message": validationError.Error(),
				// "extras":  []interface{}{validationError.ActualTag(), validationError.Value()},
			})
		}
	} else {
		return []map[string]interface{}{
			{"message": err.Error()},
		}
	}

	return errors
}
