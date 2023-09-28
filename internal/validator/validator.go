package validator

import (
	"fmt"
	"log/slog"
	"regexp"
	"strings"

	"github.com/go-playground/validator/v10"
)

type ErrorResponse struct {
	Message string
}

// Error implements error.
func (e *ErrorResponse) Error() string {
	return e.Message
}

func Validate(form interface{}) (response *ErrorResponse) {
	validate := validator.New()

	// registering our custom validator
	ok := validate.RegisterValidation("validateAlphaWithSpecialChars", validateAlphaWithSpecialChars, true)
	if ok != nil {
		slog.Error("Validator: could not register validateAlphaWithSpecialChars", "Err", ok)
	}

	ok = validate.RegisterValidation("validateNotBlank", validateNotBlank, true)
	if ok != nil {
		slog.Error("Validator: could not register validateNotBlank", "Err", ok)
	}

	err := validate.Struct(form)
	if err != nil {
		validationErrors, ok := err.(validator.ValidationErrors)
		if !ok {
			return &ErrorResponse{
				Message: "Could not validate the form! try again.",
			}
		}

		// here we should consider all the validation tags
		for _, err := range validationErrors {
			switch err.Tag() {
			case "required":
				response = &ErrorResponse{
					Message: "The " + err.Field() + " is required!",
				}

			case "email":
				response = &ErrorResponse{
					Message: fmt.Sprintf("%v is not a valid email!", err.Value()),
				}

			case "min":
				response = &ErrorResponse{
					Message: fmt.Sprintf("%v should be atleast %v characters long!", err.Field(), err.Param()),
				}

			case "max":
				response = &ErrorResponse{
					Message: fmt.Sprintf("%v should be at most %v characters long!", err.Field(), err.Param()),
				}

			case "validateAlphaWithSpecialChars":
				response = &ErrorResponse{
					Message: fmt.Sprintf("%v should only be letters, number or _ only!", err.Field()),
				}

			case "validateNotBlank":
				response = &ErrorResponse{
					Message: fmt.Sprintf("%v should not be empty", err.Field()),
				}

			default:
				response = &ErrorResponse{
					Message: "Something is wrong with " + err.Field(),
				}
			}
		}

		return response
	}

	return nil
}

// Custom validation function to check for letters, spaces, underscores, and numbers
func validateAlphaWithSpecialChars(fl validator.FieldLevel) bool {
	str := fl.Field().String()

	// Define a regular expression pattern to match the allowed characters
	pattern := "^[a-zA-Z0-9_ ]*$"

	return regexp.MustCompile(pattern).MatchString(str)
}

// Custom non blank validation
func validateNotBlank(fl validator.FieldLevel) bool {
	str := fl.Field().String()

	return strings.TrimSpace(str) != ""
}
