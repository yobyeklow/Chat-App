package utils

import (
	"net/http"
)

type ErrorCode string

const (
	ErrCodeBadRequest     ErrorCode = "BAD_REQUEST"
	ErrCodeNotFound       ErrorCode = "NOT_FOUND"
	ErrCodeConflict       ErrorCode = "CONFLICT"
	ErrCodeInternal       ErrorCode = "INTERNAL_SERVER_CODE"
	ErrCodeUnauthorized   ErrorCode = "UNAUTHORIZED"
	ErrCodeTooManyRequest ErrorCode = "TOO_MANY_REQUEST"
)

type AppError struct {
	Code    ErrorCode
	Message string
	Err     error
}

func (appError *AppError) Error() string {
	return ""
}
func NewError(message string, code ErrorCode) error {
	return &AppError{
		Message: message,
		Code:    code,
	}
}
func WrapError(message string, code ErrorCode, err error) error {
	return &AppError{
		Message: message,
		Code:    code,
		Err:     err,
	}
}
func httpStatusFromCode(code ErrorCode) int {
	switch code {
	case ErrCodeBadRequest:
		return http.StatusBadRequest
	case ErrCodeNotFound:
		return http.StatusNotFound
	case ErrCodeInternal:
		return http.StatusInternalServerError
	case ErrCodeUnauthorized:
		return http.StatusUnauthorized
	case ErrCodeTooManyRequest:
		return http.StatusTooManyRequests
	default:
		return http.StatusConflict
	}
}
