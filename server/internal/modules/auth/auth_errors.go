package auth

import "errors"

var (
	ErrAccountNotFound = errors.New("account not found")
	ErrEmailExists     = errors.New("email already exists")
)
