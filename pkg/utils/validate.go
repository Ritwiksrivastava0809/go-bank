package utils

import (
	"github.com/go-playground/validator/v10"
)

// Validator is a global instance of the validator
var Validator *validator.Validate

// InitValidator initializes the validator instance
func InitValidator() {
	Validator = validator.New()
}

// ValidateStruct validates a struct based on struct tags
func ValidateStruct(s interface{}) error {
	return Validator.Struct(s)
}

var ValidCurrency validator.Func = func(feildLevel validator.FieldLevel) bool {
	if currency, ok := feildLevel.Field().Interface().(string); ok {
		// check if the currency is valid
		flag := ISSupportedCurrency(currency)
		return flag

	}
	return false
}
