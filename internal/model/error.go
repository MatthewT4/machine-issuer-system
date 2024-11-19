package model

type baseError struct {
	Message string
}

func (e *baseError) Error() string {
	return e.Message
}

type ErrNotFound struct {
	baseError
}

func (e *ErrNotFound) Error() string {
	return e.Message
}
