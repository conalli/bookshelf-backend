package testutils

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"

	"github.com/conalli/bookshelf-backend/pkg/http/request"
	"github.com/conalli/bookshelf-backend/pkg/jwtauth"
	"github.com/conalli/bookshelf-backend/pkg/logs"
)

// RequestWithCookie provides a helper for testing handlers that require jwt cookies.
func RequestWithCookie(method, url string, body io.Reader, APIKey string, log logs.Logger) (*http.Response, error) {
	client := &http.Client{}
	req, err := http.NewRequest(method, url, body)
	if err != nil {
		return nil, err
	}
	jwt, err := jwtauth.NewTokens(APIKey, log)
	if err != nil {
		return nil, err
	}
	req.AddCookie(&http.Cookie{Name: "bookshelfjwt", Value: jwt["access_token"]})
	return client.Do(req)
}

// MakeRequestBody takes in a struct and attempts to marshal it and turn it into a new buffer.
func MakeRequestBody[T request.APIRequest](data T) (*bytes.Buffer, error) {
	body, err := json.Marshal(data)
	if err != nil {
		return nil, err
	}
	return bytes.NewBuffer(body), nil
}
