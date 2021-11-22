package apiErrors

import (
	"errors"
	"fmt"
	"net/http"
)

type ApiError struct {
	ErrStatus int    `json:"errStatus,omitempty"`
	ErrValue  string `json:"err,omitempty"`
}

var (
	ErrBadRequest          = errors.New("bad request")
	ErrWrongCredentials    = errors.New("wrong credentials")
	ErrNotFound            = errors.New("not found")
	ErrUnauthorized        = errors.New("unauthorized")
	ErrForbidden           = errors.New("forbidden")
	ErrPermissionDenied    = errors.New("permission denied")
	ErrBadQueryParams      = errors.New("invalid query params")
	ErrInternalServerError = errors.New("internal server error")
	ErrRequestTimeoutError = errors.New("request timeout")
	ErrInvalidJWTToken     = errors.New("invalid JWT token")
	ErrInvalidJWTClaims    = errors.New("invalid JWT claims")
)

func (e ApiError) Error() string {
	return fmt.Sprintf("error: status-%d, value %s", e.ErrStatus, e.ErrValue)
}

func (e ApiError) Status() int {
	return e.ErrStatus
}

func NewApiError(status int, value string) ApiError {
	return ApiError{
		ErrStatus: status,
		ErrValue:  value,
	}
}

func NewBadRequestError() ApiError {
	return ApiError{
		ErrStatus: http.StatusBadRequest,
		ErrValue:  ErrBadRequest.Error(),
	}
}

func NewInternalServerError() ApiError {
	return ApiError{
		ErrStatus: http.StatusInternalServerError,
		ErrValue:  ErrInternalServerError.Error(),
	}
}
