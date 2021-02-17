package model

type InvalidArgumentError struct {
	error string
}

func NewInvalidArgument(error string) error {
	return &InvalidArgumentError{error: error}
}

func (e *InvalidArgumentError) Error() string {
	return e.error
}

type ForbiddenError struct {
	error string
}

func NewForbiddenError(error string) error {
	return &ForbiddenError{error: error}
}

func (e *ForbiddenError) Error() string {
	return e.error
}

type NotFoundError struct {
	error string
}

func NewNotFoundError(error string) error {
	return &NotFoundError{error: error}
}

func (e *NotFoundError) Error() string {
	return e.error
}

type ConflictError struct {
	error string
}

func NewConflictError(error string) error {
	return &ConflictError{error: error}
}

func (e *ConflictError) Error() string {
	return e.error
}
