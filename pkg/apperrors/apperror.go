package apperrors

import (
	"net/http"
)

type Error struct {
	Status      int    `json:"status,omitempty"`
	Code        string `json:"code,omitempty"`
	Message     string `json:"message,omitempty"`
	Description string `json:"description,omitempty"`
}

type AppErrorOption func(appError *Error)

func New(opts ...AppErrorOption) Error {
	appErr := Error{
		Status: http.StatusInternalServerError,
	}
	for _, opt := range opts {
		opt(&appErr)
	}
	return appErr
}

func NewError(status int, code, message string) Error {
	return Error{Status: status, Code: code, Message: message}
}

func (e Error) Error() string {
	if e.Description != "" {
		return e.Description
	}

	return e.Message
}

func ErrorWithDescription(err Error, description string) error {
	return Error{
		Status:      err.Status,
		Code:        err.Code,
		Message:     err.Message,
		Description: description,
	}
}

func (e Error) WithDescription(description string) Error {
	e.Description = description
	return e
}

func WithStatus(status int) AppErrorOption {
	return func(appError *Error) {
		appError.Status = status
	}
}

func WithCode(code string) AppErrorOption {
	return func(appError *Error) {
		appError.Code = code
	}
}

func WithMessage(message string) AppErrorOption {
	return func(appError *Error) {
		appError.Message = message
	}
}

func WithDescription(description string) AppErrorOption {
	return func(appError *Error) {
		appError.Description = description
	}
}

func FromError(err error) error {
	if IsNotFoundError(err) {
		return Error{
			Status:  http.StatusNotFound,
			Code:    "NOT_FOUND",
			Message: "not found resources",
		}
	}
	if IsConstraintViolationError(err) {
		return Error{
			Status:  http.StatusBadRequest,
			Code:    "CONSTRAINT_VIOLATION",
			Message: "constraint violation",
		}
	}
	if IsObjectNotInPrerequisiteStateError(err) {
		return Error{
			Status:  http.StatusBadRequest,
			Code:    "OBJECT_NOT_IN_PREREQUISITE_STATE",
			Message: "object not in prerequisite state",
		}
	}
	return Error{
		Status:  http.StatusInternalServerError,
		Code:    "SERVER_ERROR",
		Message: "server error",
	}
}
