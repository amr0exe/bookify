package service

import "errors"

var (
	ErrNoServiceFound  = errors.New("No records of service found.")
	ErrNoBusinessFound = errors.New("No business profile found for the given accoundId")
)
