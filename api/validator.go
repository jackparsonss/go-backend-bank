package api

import (
	"go-backend/util"

	"github.com/go-playground/validator/v10"
)

// This code defines a custom validation function called `validCurrency` using the `validator` package
// in Go. The function takes a `fieldLevel` parameter of type `validator.FieldLevel` and returns a
// boolean value.
var validCurrency validator.Func = func(fieldLevel validator.FieldLevel) bool {
	if currency, ok := fieldLevel.Field().Interface().(string); ok {
		// check currency is valid
		return util.IsSupportedCurrency(currency)
	}
	return false
}
