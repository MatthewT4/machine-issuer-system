package errorlist

import "errors"

var (
	ErrInvalidCredentials = errors.New("invalid credentials")

	ErrUnexpectedSigningMethod = errors.New("unexpected signing method")
	ErrInvalidToken            = errors.New("invalid token")
	ErrInvalidTokenClaims      = errors.New("invalid token claims")
)
