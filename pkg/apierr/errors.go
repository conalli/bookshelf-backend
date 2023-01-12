package apierr

import (
	"encoding/json"
	"errors"
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
	status int
	err    error
	detail string
}

// APIErr represents the methods needed to return an APIErr.
type Error interface {
	Status() int
	Error() string
	Detail() string
}

// ResError represents an error response.
type ResError struct {
	Status int    `json:"status,omitempty"`
	Title  string `json:"title,omitempty"`
	Detail string `json:"detail,omitempty"`
}

// Status returns the status code of an APIError.
func (e APIError) Status() int {
	return e.status
}

func (e APIError) Error() string {
	return e.err.Error()
}

func (e APIError) Detail() string {
	return e.detail
}

// NewAPIError returns a new APIError with given arguments.
func NewAPIError(status int, value error, detail string) APIError {
	return APIError{
		status: status,
		err:    value,
		detail: detail,
	}
}

// NewBadRequestError returns a bad request APIError with given arguments.
func NewBadRequestError(detail string) APIError {
	return APIError{
		status: http.StatusBadRequest,
		err:    ErrBadRequest,
		detail: detail,
	}
}

// NewUnauthorizedError returns a wrong credentials APIError with given arguments.
func NewUnauthorizedError(detail string) APIError {
	return APIError{
		status: http.StatusUnauthorized,
		err:    ErrUnauthorized,
		detail: detail,
	}
}

// NewWrongCredentialsError returns a wrong credentials APIError with given arguments.
func NewWrongCredentialsError(detail string) APIError {
	return APIError{
		status: http.StatusUnauthorized,
		err:    ErrWrongCredentials,
		detail: detail,
	}
}

// NewInternalServerError returns an internal server error APIError.
func NewInternalServerError() APIError {
	return APIError{
		status: http.StatusInternalServerError,
		err:    ErrInternalServerError,
	}
}

// NewJWTTokenError returns a wrong credentials APIError with given arguments.
func NewJWTTokenError(detail string) APIError {
	return APIError{
		status: http.StatusUnauthorized,
		err:    ErrUnauthorized,
		detail: detail,
	}
}

// NewJWTClaimsError returns a wrong credentials APIError with given arguments.
func NewJWTClaimsError(detail string) APIError {
	return APIError{
		status: http.StatusInternalServerError,
		err:    ErrInvalidJWTClaims,
		detail: detail,
	}
}

// APIErrorResponse encodes the response with an APIErr.
func APIErrorResponse(w http.ResponseWriter, err Error) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(err.Status())
	res := ResError{
		Status: err.Status(),
		Title:  err.Error(),
		Detail: err.Detail(),
	}
	json.NewEncoder(w).Encode(res)
}
