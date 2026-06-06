package types

import "errors"

var (
	ErrNotFound          = errors.New("Not found")
	ErrInvalidSearchTerm = errors.New("Invalid search term")
	ErrNoBodyFound       = errors.New("No body found")
)
