package utils

import (
	"encoding/json"
	"errors"
	"fmt"
	"strconv"
	"strings"

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

var cardNumber validator.Func = func(fl validator.FieldLevel) bool {
	value, ok := fl.Field().Interface().(string)
	if !ok {
		return false
	}

	if len(value) < 13 || len(value) > 19 {
		return false
	}

	for _, r := range value {
		if r < '0' || r > '9' {
			return false
		}
	}

	return true
}

var cardExpirate validator.Func = func(fl validator.FieldLevel) bool {
	value, ok := fl.Field().Interface().(string)
	if !ok {
		return false
	}

	if len(value) != 5 {
		return false
	}

	expiry := strings.Split(value, "/")
	if len(expiry) != 2 {
		return false
	}

	month := expiry[0]
	year := expiry[1]

	if len(month) != 2 || len(year) != 2 {
		return false
	}

	monthInt, err := strconv.Atoi(month)
	if err != nil || monthInt < 1 || monthInt > 12 {
		return false
	}

	yearInt, err := strconv.Atoi(year)
	if err != nil || yearInt < 0 || yearInt > 99 {
		return false
	}

	return true
}

// init initializes the validator engine with English translations.
// It sets up the translation system using the "en" locale and registers
// the default translations for the validator.
func init() {
	if value, ok := binding.Validator.Engine().(*validator.Validate); ok {
		en := en.New()
		unt := ut.New(en, en)
		transl, _ = unt.GetTranslator("en")

		en_translation.RegisterDefaultTranslations(value, transl)

		value.RegisterValidation("cnumber", cardNumber)
		value.RegisterTranslation("cnumber", transl, func(ut ut.Translator) error {
			return ut.Add("cnumber", "{0} must be a valid card number", true)
		}, func(ut ut.Translator, fe validator.FieldError) string {
			t, _ := ut.T("cnumber", fe.Field())
			return t
		})

		value.RegisterValidation("cexpirate", cardExpirate)
		value.RegisterTranslation("cexpirate", transl, func(ut ut.Translator) error {
			return ut.Add("cexpirate", "{0} must be in the format MM/YY", true)
		}, func(ut ut.Translator, fe validator.FieldError) string {
			t, _ := ut.T("cexpirate", fe.Field())
			return t
		})
	}
}

// ValidatorError processes validation errors and returns a slice of Errors.
// It handles different types of validation errors including JSON unmarshal errors
// and validator.ValidationErrors. The function converts these errors into a
// standardized Errors format with field names in snake_case and appropriate messages.
//
// Parameters:
//   - validateError: error - The validation error to be processed.
//
// Returns:
//   - []Errors: A slice of Errors containing the processed error details.
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
				Message: e.Translate(transl),
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
		if r == '.' {
			snakeCase = append(snakeCase, '.')
		} else if r >= 'A' && r <= 'Z' {
			if i > 0 && str[i-1] != '.' {
				snakeCase = append(snakeCase, '_')
			}
			snakeCase = append(snakeCase, r+'a'-'A')
		} else {
			snakeCase = append(snakeCase, r)
		}
	}
	// Remove leading and trailing underscores
	result := strings.Trim(string(snakeCase), "_")
	return result
}
