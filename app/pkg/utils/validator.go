package utils

import (
	"encoding/json"
	"errors"
	"fmt"

	"github.com/gin-gonic/gin/binding"
	"github.com/go-playground/locales/en"
	ut "github.com/go-playground/universal-translator"
	"github.com/go-playground/validator/v10"

	en_translation "github.com/go-playground/validator/v10/translations/en"
)

var (
	Validate = validator.New()
	transl   ut.Translator
)

type Errors struct {
	Field   string `json:"field" exemple:"name"`
	Message string `json:"message" exemple:"name is required"`
}

func init() {
	if value, ok := binding.Validator.Engine().(*validator.Validate); ok {
		en := en.New()
		unt := ut.New(en, en)
		transl, _ = unt.GetTranslator("en")
		en_translation.RegisterDefaultTranslations(value, transl)
	}
}

func ValidatorError(validateError error) []Errors {

	var jsonValidatorError validator.ValidationErrors
	var jsonError *json.UnmarshalFieldError

	errorsCauses := []Errors{}
	if errors.As(validateError, &jsonError) {
		e := validateError.(*json.UnmarshalTypeError)
		cause := Errors{
			Field:   toSnakeCase(e.Field),
			Message: fmt.Sprintf("it can not be %s", e.Value),
		}

		errorsCauses = append(errorsCauses, cause)

	} else if errors.As(validateError, &jsonValidatorError) {
		for _, e := range validateError.(validator.ValidationErrors) {
			cause := Errors{
				Field:   toSnakeCase(e.Namespace()),
				Message: toSnakeCase(e.Translate(transl)),
			}

			errorsCauses = append(errorsCauses, cause)
		}
	} else {
		cause := Errors{
			Field:   "payload",
			Message: fmt.Sprintf("please check your payload: %s", toSnakeCase(validateError.Error())),
		}

		errorsCauses = append(errorsCauses, cause)
	}

	return errorsCauses
}

func toSnakeCase(str string) string {
	var snakeCase []rune
	for i, r := range str {
		if r >= 'A' && r <= 'Z' {
			if i > 0 {
				snakeCase = append(snakeCase, '_')
			}
			snakeCase = append(snakeCase, r+'a'-'A')
		} else {
			snakeCase = append(snakeCase, r)
		}
	}
	return string(snakeCase)
}
