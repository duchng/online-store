package defined_errors

import (
	"net/http"

	"store-management/pkg/apperrors"
)

var (
	ErrIncorrectPassword = apperrors.Error{
		Status:      http.StatusUnauthorized,
		Code:        "001",
		Message:     "Password is incorrect",
		Description: "",
	}
	ErrPasswordMismath = apperrors.Error{
		Status:      http.StatusBadRequest,
		Code:        "002",
		Message:     "Password is incorrect",
		Description: "",
	}
)
