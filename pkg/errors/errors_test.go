package errors_test

import (
	"testing"

	"github.com/conalli/bookshelf-backend/pkg/errors"
	"go.mongodb.org/mongo-driver/mongo"
)

func TestNewApiError(t *testing.T) {
	t.Parallel()
	err := errors.NewApiError(500, "internal server error", "Something went wrong...")
	status := 500
	errStatus := err.Status()
	if status != errStatus {
		t.Errorf("Wanted error with status: %d, got status %d", status, errStatus)
	}
	msg := "Something went wrong... -- internal server error"
	errMsg := err.Error()
	if msg != errMsg {
		t.Errorf("Wanted error with status: %s, got status %s", msg, errMsg)
	}
}

func TestNewBadRequestError(t *testing.T) {
	t.Parallel()
	err := errors.NewBadRequestError("This is a bad request")
	status := 400
	errStatus := err.Status()
	if status != errStatus {
		t.Errorf("Wanted error with status: %d, got status %d", status, errStatus)
	}
	msg := "This is a bad request -- bad request"
	errMsg := err.Error()
	if msg != errMsg {
		t.Errorf("Wanted error with status: %s, got status %s", msg, errMsg)
	}
}

func TestNewWrongCredentialsError(t *testing.T) {
	t.Parallel()
	err := errors.NewWrongCredentialsError("You have the wrong credentials")
	status := 401
	errStatus := err.Status()
	if status != errStatus {
		t.Errorf("Wanted error with status: %d, got status %d", status, errStatus)
	}
	msg := "You have the wrong credentials -- wrong credentials"
	errMsg := err.Error()
	if msg != errMsg {
		t.Errorf("Wanted error with status: %s, got status %s", msg, errMsg)
	}
}

func TestNewInternalServerError(t *testing.T) {
	t.Parallel()
	err := errors.NewInternalServerError()
	status := 500
	errStatus := err.Status()
	if status != errStatus {
		t.Errorf("Wanted error with status: %d, got status %d", status, errStatus)
	}
	msg := " -- internal server error"
	errMsg := err.Error()
	if msg != errMsg {
		t.Errorf("Wanted error with status: %s, got status %s", msg, errMsg)
	}
}

func TestNewJWTTokenError(t *testing.T) {
	t.Parallel()
	err := errors.NewJWTTokenError("Access token invalid")
	status := 500
	errStatus := err.Status()
	if status != errStatus {
		t.Errorf("Wanted error with status: %d, got status %d", status, errStatus)
	}
	msg := "Access token invalid -- invalid JWT token"
	errMsg := err.Error()
	if msg != errMsg {
		t.Errorf("Wanted error with status: %s, got status %s", msg, errMsg)
	}
}

func TestNewJWTClaimsError(t *testing.T) {
	t.Parallel()
	err := errors.NewJWTClaimsError("Access token invalid")
	status := 500
	errStatus := err.Status()
	if status != errStatus {
		t.Errorf("Wanted error with status: %d, got status %d", status, errStatus)
	}
	msg := "Access token invalid -- invalid JWT claims"
	errMsg := err.Error()
	if msg != errMsg {
		t.Errorf("Wanted error with status: %s, got status %s", msg, errMsg)
	}
}

func TestParseGetUserError(t *testing.T) {
	t.Parallel()
	err := mongo.ErrNoDocuments
	res := errors.ParseGetUserError("This is a no documents error", err)
	status := 400
	errStatus := res.Status()
	if status != errStatus {
		t.Errorf("Wanted error with status: %d, got status %d", status, errStatus)
	}
	msg := "mongo: no documents in result"
	errMsg := err.Error()
	if msg != errMsg {
		t.Errorf("Wanted error with status: %s, got status %s", msg, errMsg)
	}

	err = errors.NewInternalServerError()
	res = errors.ParseGetUserError("This is not a no documents error", err)
	status = 500
	errStatus = res.Status()
	if status != errStatus {
		t.Errorf("Wanted error with status: %d, got status %d", status, errStatus)
	}
	msg = " -- internal server error"
	errMsg = err.Error()
	if msg != errMsg {
		t.Errorf("Wanted error with status: %s, got status %s", msg, errMsg)
	}
}
