package repository

import "errors"

var (
	// ErrNotFound is returned when a specified item was not found from repository.
	ErrNotFound = errors.New("requested item was not found")
)
