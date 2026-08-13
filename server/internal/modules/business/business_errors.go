package business

import "errors"

var (
	ErrBusinessAlreadyExists = errors.New("Business profile already exists.")
	ErrBusinessNotFound      = errors.New("Business profile doens't exists.")
)
