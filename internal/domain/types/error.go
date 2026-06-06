package types

import "errors"

var (
	ErrNotFound          = errors.New("Not found")
	ErrInvalidSearchTerm = errors.New("Invalid search term")
	ErrNoBodyFound       = errors.New("No body found")
	ErrTimeout           = errors.New("Timeout")
	ErrLooseConnection   = errors.New("Loose connection")
	ErrInvalidUserID     = errors.New("Invalid user ID")
	ErrForbidden         = errors.New("Forbidden")
)
