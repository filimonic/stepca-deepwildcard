package validator

import "deepwildcard/internal/deepwildcard/validator/errors"

type Result struct {
	Allowed bool
	Reason  *errors.ValidateError
}

func ResultAllowed() *Result {
	return &Result{
		Allowed: true,
		Reason:  nil,
	}
}

func ResultDenied(reasonCode string, reasonMessage string) *Result {
	return &Result{
		Allowed: false,
		Reason:  errors.NewValidateError(reasonCode, "%s", reasonMessage),
	}
}
