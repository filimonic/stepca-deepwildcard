package errors

import (
	goErrors "errors"
	"fmt"
)

type ValidateError struct {
	Message error
	Code    string
}

func NewValidateError(code string, format string, a ...any) *ValidateError {
	return &ValidateError{
		Message: fmt.Errorf(format, a...),
		Code:    code,
	}
}

func (ve *ValidateError) String() string {
	return fmt.Sprintf("ValidateError<%s>: %s", ve.Code, ve.Message)
}

func (ve *ValidateError) Error() error {
	return goErrors.New(ve.String())
}
