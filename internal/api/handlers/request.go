package handlers

import (
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
)

// bindJSON binds the request body to the given object
//   - if error validation returns error response
func bindJSON(c *gin.Context, obj interface{}) error {
	if err := c.BindJSON(obj); err != nil {
		BadRequest(c, "Invalid inputs, please check and try again", ExtractValidationErrors(err))
		return err
	}
	return nil
}

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
