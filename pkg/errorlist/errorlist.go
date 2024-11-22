package errorlist

import "errors"

var (
	ErrInvalidCredentials = errors.New("invalid credentials")

	ErrUnexpectedSigningMethod = errors.New("unexpected signing method")
	ErrInvalidToken            = errors.New("invalid token")
	ErrInvalidTokenClaims      = errors.New("invalid token claims")
	ErrInvalidRole             = errors.New("invalid role")
	ErrHandlerNotAllowed       = errors.New("handler not allowed")
)
