package model

import "fmt"

type ErrorType string

const (
	ErrorTypeValidation = "validation"
	ErrorTypeNotFound   = "not_found"
	ErrorTypeInternal   = "internal"
)

type Error struct {
	Type    ErrorType
	Message string
	Err     error
}

func (e *Error) Error() string {
	if e.Err != nil {
		return fmt.Sprintf("%s: %s: %v", e.Type, e.Message, e.Err)
	}

	return fmt.Sprintf("%s: %s", e.Type, e.Message)
}

func NewValidationError(message string) *Error {
	return &Error{
		Type:    ErrorTypeValidation,
		Message: message,
	}
}

func NewNotFoundError(message string) *Error {
	return &Error{
		Type:    ErrorTypeNotFound,
		Message: message,
	}
}

func NewInternalError(err error) *Error {
	return &Error{
		Type:    ErrorTypeInternal,
		Message: "internal error occurred",
		Err:     err,
	}
}
