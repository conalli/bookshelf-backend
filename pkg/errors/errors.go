package errors

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
)

var (
	// ErrBadRequest represents an HTTP bad request error.
	ErrBadRequest = errors.New("bad request")
	// ErrWrongCredentials represents an HTTP wrong credentials error.
	ErrWrongCredentials = errors.New("wrong credentials")
	// ErrNotFound represents an HTTP not found error.
	ErrNotFound = errors.New("not found")
	// ErrUnauthorized represents an HTTP unauthorized error.
	ErrUnauthorized = errors.New("unauthorized")
	// ErrForbidden represents an HTTP forbidden error.
	ErrForbidden = errors.New("forbidden")
	// ErrPermissionDenied represents an HTTP permission denied error.
	ErrPermissionDenied = errors.New("permission denied")
	// ErrBadQueryParams represents an HTTP bad query params error.
	ErrBadQueryParams = errors.New("invalid query params")
	// ErrInternalServerError represents an HTTP internal server error.
	ErrInternalServerError = errors.New("internal server error")
	// ErrRequestTimeoutError represents an HTTP request timeout error.
	ErrRequestTimeoutError = errors.New("request timeout")
	// ErrInvalidJWTToken represents an HTTP invalid jwt token error.
	ErrInvalidJWTToken = errors.New("invalid JWT token")
	// ErrInvalidJWTClaims represents an HTTP invalid jwt claims error.
	ErrInvalidJWTClaims = errors.New("invalid JWT claims")
)

// APIError represents an Api/server error.
type APIError struct {
	ErrStatus  int
	ErrValue   string
	ErrDetails string
}

// APIErr represents the methods needed to return an APIErr.
type APIErr interface {
	Status() int
	Error() string
}

// ResError represents an error response.
type ResError struct {
	Status int    `json:"status,omitempty"`
	Error  string `json:"error,omitempty"`
}

// Status returns the status code of an APIError.
func (e APIError) Status() int {
	return e.ErrStatus
}

func (e APIError) Error() string {
	return fmt.Sprintf("%s -- %s", e.ErrDetails, e.ErrValue)
}

// NewAPIError returns a new APIError with given arguments.
func NewAPIError(status int, value string, details string) APIError {
	return APIError{
		ErrStatus:  status,
		ErrValue:   value,
		ErrDetails: details,
	}
}

// NewBadRequestError returns a bad request APIError with given arguments.
func NewBadRequestError(details string) APIError {
	return APIError{
		ErrStatus:  http.StatusBadRequest,
		ErrValue:   ErrBadRequest.Error(),
		ErrDetails: details,
	}
}

// NewWrongCredentialsError returns a wrong credentials APIError with given arguments.
func NewWrongCredentialsError(details string) APIError {
	return APIError{
		ErrStatus:  http.StatusUnauthorized,
		ErrValue:   ErrWrongCredentials.Error(),
		ErrDetails: details,
	}
}

// NewInternalServerError returns an internal server error APIError.
func NewInternalServerError() APIError {
	return APIError{
		ErrStatus: http.StatusInternalServerError,
		ErrValue:  ErrInternalServerError.Error(),
	}
}

// NewJWTTokenError returns a wrong credentials APIError with given arguments.
func NewJWTTokenError(details string) APIError {
	return APIError{
		ErrStatus:  http.StatusInternalServerError,
		ErrValue:   ErrInvalidJWTToken.Error(),
		ErrDetails: details,
	}
}

// NewJWTClaimsError returns a wrong credentials APIError with given arguments.
func NewJWTClaimsError(details string) APIError {
	return APIError{
		ErrStatus:  http.StatusInternalServerError,
		ErrValue:   ErrInvalidJWTClaims.Error(),
		ErrDetails: details,
	}
}

// APIErrorResponse encodes the response with an APIErr.
func APIErrorResponse(w http.ResponseWriter, err APIErr) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(err.Status())
	res := ResError{
		Status: err.Status(),
		Error:  err.Error(),
	}
	json.NewEncoder(w).Encode(res)
}
