package consumer

import "errors"

var (
	ErrConsumerAlreadyExists = errors.New("consumer already exists")
	ErrConsumerNotFound      = errors.New("consumer not found")
)
