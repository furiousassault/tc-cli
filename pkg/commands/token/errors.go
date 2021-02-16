package token

import "errors"

var (
	errValidation      = errors.New("validation error")
	errTokenNameExists = errors.New("token name already exists")
)
