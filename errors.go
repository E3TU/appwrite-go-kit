package base

import "errors"

var (
	// ErrInvalidConfig reports a Config that failed validation.
	ErrInvalidConfig = errors.New("appwrite-base: invalid config")

	// ErrDisabled reports use of a capability that is not configured.
	ErrDisabled = errors.New("appwrite-base: capability disabled")
)
