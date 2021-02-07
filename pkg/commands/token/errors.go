package token

import "errors"

var (
	errValidation      = errors.New("validation error")
	errTokenNamesMatch = errors.New("token names are equal")
)
