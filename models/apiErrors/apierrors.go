package apiErrors

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"

	"go.mongodb.org/mongo-driver/mongo"
)

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

type ApiError struct {
	ErrStatus  int
	ErrValue   string
	ErrDetails string
}

type ApiErr interface {
	Status() int
	Error() string
}

type ResError struct {
	Status int    `json:"status,omitempty"`
	Error  string `json:"error,omitempty"`
}

func (e ApiError) Status() int {
	return e.ErrStatus
}

func (e ApiError) Error() string {
	return fmt.Sprintf("%s -- %s", e.ErrDetails, e.ErrValue)
}

func NewApiError(status int, value string, details string) ApiError {
	return ApiError{
		ErrStatus:  status,
		ErrValue:   value,
		ErrDetails: details,
	}
}

func NewBadRequestError(details string) ApiError {
	return ApiError{
		ErrStatus:  http.StatusBadRequest,
		ErrValue:   ErrBadRequest.Error(),
		ErrDetails: details,
	}
}

func NewWrongCredentialsError(details string) ApiError {
	return ApiError{
		ErrStatus:  http.StatusUnauthorized,
		ErrValue:   ErrWrongCredentials.Error(),
		ErrDetails: details,
	}
}

func NewInternalServerError() ApiError {
	return ApiError{
		ErrStatus: http.StatusInternalServerError,
		ErrValue:  ErrInternalServerError.Error(),
	}
}

func NewJWTTokenError(details string) ApiError {
	return ApiError{
		ErrStatus:  http.StatusInternalServerError,
		ErrValue:   ErrInvalidJWTToken.Error(),
		ErrDetails: details,
	}
}

func NewJWTClaimsError(details string) ApiError {
	return ApiError{
		ErrStatus:  http.StatusInternalServerError,
		ErrValue:   ErrInvalidJWTClaims.Error(),
		ErrDetails: details,
	}
}

func ParseGetUserError(value string, err error) ApiError {
	if err == mongo.ErrNoDocuments {
		return NewBadRequestError("error: could not find user with value " + value)
	}
	return NewInternalServerError()
}

func APIErrorResponse(w http.ResponseWriter, err ApiErr) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(err.Status())
	res := ResError{
		Status: err.Status(),
		Error:  err.Error(),
	}
	json.NewEncoder(w).Encode(res)
}
