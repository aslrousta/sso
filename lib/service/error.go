package service

import "fmt"

// Error reports a service error.
type Error struct {
	Service string
	Message string
	Cause   error
}

func (e *Error) Error() string {
	if e.Cause == nil {
		return fmt.Sprintf("%s: %s", e.Service, e.Message)
	}
	return fmt.Sprintf("%s: %s: %v", e.Service, e.Message, e.Cause)
}

func (e *Error) Unwrap() error {
	return e.Cause
}

// NotFoundError reports a missing entity error.
type NotFoundError struct {
	Entity string
	Key    string
}

func (e *NotFoundError) Error() string {
	return fmt.Sprintf("%s not found (key = '%s')", e.Entity, e.Key)
}
