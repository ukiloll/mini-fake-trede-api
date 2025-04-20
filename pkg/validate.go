package pkg

import (
	"fmt"

	"github.com/go-playground/validator/v10"
)

type CustomValidationError struct {
	Field   string `json:"field"`
	Message string `json:"message"`
}

func validateStruct(s interface{}) []CustomValidationError {
	validate := validator.New()
	err := validate.Struct(s)
	if err == nil {
		return nil // No validation errors
	}

	var errors []CustomValidationError
	for _, fieldError := range err.(validator.ValidationErrors) {
		// Custom error message mapping
		message := getCustomErrorMessage(fieldError)
		errors = append(errors, CustomValidationError{
			Field:   fieldError.Field(),
			Message: message,
		})
	}

	return errors
}

// getCustomErrorMessage maps validation rules to human-readable messages
func getCustomErrorMessage(err validator.FieldError) string {
	switch err.Tag() {
	case "required":
		return fmt.Sprintf("%s is required", err.Field())
	case "min":
		return fmt.Sprintf("%s must be at least %s", err.Field(), err.Param())
	default:
		return fmt.Sprintf("%s is invalid", err.Field())
	}
}

// errors := validateStruct(bodyAsset)
// if len(errors) > 0 {
// 	w.Header().Set("Content-Type", "application/json")
// 	w.WriteHeader(http.StatusBadRequest)
// 	json.NewEncoder(w).Encode(map[string]interface{}{
// 		"errors": errors,
// 	})
// 	return
// }