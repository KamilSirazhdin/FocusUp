package utils

import (
	"fmt"
	"strings"

	"github.com/go-playground/validator/v10"
)

var validate = validator.New()

// ValidateStruct валидирует структуру с помощью validator
func ValidateStruct(s interface{}) error {
	err := validate.Struct(s)
	if err != nil {
		return formatValidationErrors(err)
	}
	return nil
}

// formatValidationErrors форматирует ошибки валидации в читаемый вид
func formatValidationErrors(err error) error {
	if err == nil {
		return nil
	}

	validationErrors, ok := err.(validator.ValidationErrors)
	if !ok {
		return err
	}

	errorMessages := []string{}
	for _, e := range validationErrors {
		field := e.Field()
		tag := e.Tag()
		param := e.Param()

		var message string
		switch tag {
		case "required":
			message = fmt.Sprintf("Поле '%s' обязательно", field)
		case "email":
			message = fmt.Sprintf("Поле '%s' должно быть email адресом", field)
		case "min":
			message = fmt.Sprintf("Поле '%s' должно содержать не менее %s символов", field, param)
		case "max":
			message = fmt.Sprintf("Поле '%s' должно содержать не более %s символов", field, param)
		default:
			message = fmt.Sprintf("Поле '%s' не прошло валидацию по правилу '%s'", field, tag)
		}
		errorMessages = append(errorMessages, message)
	}

	return fmt.Errorf(strings.Join(errorMessages, "; "))
}
