package model

type BaseError struct {
	Message string
}

func (e *BaseError) Error() string {
	return e.Message
}

type ErrNotFound struct {
	BaseError
}

type ErrInternal struct {
	BaseError
}

type ErrBadRequest struct {
	BaseError
}
