package repository

import "errors"

var (
	// ErrResourceNotFound is returned when a specified item was not found from repository.
	ErrResourceNotFound = errors.New("requested resource not found")
)
