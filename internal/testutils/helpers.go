package testutils

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"

	"github.com/conalli/bookshelf-backend/pkg/http/request"
	"github.com/conalli/bookshelf-backend/pkg/logs"
	"github.com/conalli/bookshelf-backend/pkg/services/auth"
)

// RequestWithCookie provides a helper for testing handlers that require jwt cookies.
func RequestWithCookie(method, url string, body io.Reader, APIKey string, log logs.Logger) (*http.Response, error) {
	client := &http.Client{}
	req, err := http.NewRequest(method, url, body)
	if err != nil {
		return nil, err
	}
	tokens, err := auth.NewTokens(log, APIKey)
	if err != nil {
		return nil, err
	}
	cookies := tokens.NewTokenCookies(log)
	code := request.FilterCookies(cookies, auth.BookshelfTokenCode)
	access := request.FilterCookies(cookies, auth.BookshelfAccessToken)
	refresh := request.FilterCookies(cookies, auth.BookshelfRefreshToken)
	req.AddCookie(code)
	req.AddCookie(access)
	req.AddCookie(refresh)
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
