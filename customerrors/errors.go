package customerrors

import "errors"

var (
	ErrNotFound = errors.New("operation with given id not found")
	ErrNotMatch = errors.New("code does not match")
)
