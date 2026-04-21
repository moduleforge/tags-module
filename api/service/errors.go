package service

import "errors"

// ErrNotFound is returned when a requested resource does not exist or is not
// visible to the caller (to avoid leaking existence).
var ErrNotFound = errors.New("not found")

// ErrForbidden is returned when the caller lacks permission for the operation
// but is known to have visibility of the resource (e.g., subject trying PUT).
var ErrForbidden = errors.New("forbidden")

// ErrInvalidInput is returned when the caller supplies invalid or missing input.
var ErrInvalidInput = errors.New("invalid input")

// ErrConflict is returned when a uniqueness constraint would be violated.
var ErrConflict = errors.New("conflict")
