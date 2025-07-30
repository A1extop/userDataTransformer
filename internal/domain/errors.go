package domain

import "errors"

var (
	ErrInvalidRequest = errors.New("invalid request")
	ErrNoAuthToken    = errors.New("no auth token")
	ErrNoData         = errors.New("xmlUsers len = 0")
	ErrFailedAdd      = errors.New("failed to add items to local store")
)
