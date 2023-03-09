package testutils

import (
	"bytes"
	"crypto/rand"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"os"

	"github.com/conalli/bookshelf-backend/pkg/http/request"
	"github.com/conalli/bookshelf-backend/pkg/logs"
	"github.com/conalli/bookshelf-backend/pkg/services/accounts"
	"github.com/conalli/bookshelf-backend/pkg/services/auth"
	"github.com/conalli/bookshelf-backend/pkg/services/bookmarks"
)

type requestOptions struct {
	headers map[string]string
	body    io.Reader
	APIKey  string
	log     logs.Logger
}

func NewRequestOptions() *requestOptions {
	return &requestOptions{}
}

type RequestOption func(*requestOptions)

func WithHeaders(headers map[string]string) RequestOption {
	return func(ro *requestOptions) {
		ro.headers = headers
	}
}

func WithBody(body io.Reader) RequestOption {
	return func(ro *requestOptions) {
		ro.body = body
	}
}

func WithAPIKey(APIKey string) RequestOption {
	return func(ro *requestOptions) {
		ro.APIKey = APIKey
	}
}

func WithLogger(log logs.Logger) RequestOption {
	return func(ro *requestOptions) {
		ro.log = log
	}
}

// RequestWithCookie provides a helper for testing handlers that require jwt cookies.
func RequestWithCookie(method, url string, options ...RequestOption) (*http.Response, error) {
	client := &http.Client{}
	ro := NewRequestOptions()
	for _, opt := range options {
		opt(ro)
	}
	req, err := http.NewRequest(method, url, ro.body)
	if err != nil {
		return nil, err
	}
	for key, val := range ro.headers {
		req.Header.Add(key, val)
	}
	tokens, err := auth.NewTokens(ro.log, ro.APIKey)
	if err != nil {
		return nil, err
	}
	cookies := tokens.NewTokenCookies(ro.log, http.SameSiteStrictMode)
	code := request.FilterCookies(cookies, auth.BookshelfTokenCode)
	access := request.FilterCookies(cookies, auth.BookshelfAccessToken)
	req.AddCookie(code)
	req.AddCookie(access)
	return client.Do(req)
}

// MakeJSONRequestBody takes in a struct and attempts to marshal it and turn it into a new buffer.
func MakeJSONRequestBody[T request.APIRequest](data T) (*bytes.Buffer, error) {
	body, err := json.Marshal(data)
	if err != nil {
		return nil, err
	}
	return bytes.NewBuffer(body), nil
}

func MakeFileRequestBody(path, filename string) (*bytes.Buffer, string, error) {
	body := new(bytes.Buffer)
	writer := multipart.NewWriter(body)
	ff, err := writer.CreateFormFile(bookmarks.BookmarksFileKey, filename)
	if err != nil {
		return nil, "", err
	}
	file, err := os.Open(path)
	if err != nil {
		return nil, "", err
	}
	_, err = io.Copy(ff, file)
	if err != nil {
		return nil, "", err
	}
	writer.Close()
	return body, writer.FormDataContentType(), nil
}

func randomID(length int) (string, error) {
	bytes := make([]byte, length)
	if _, err := rand.Read(bytes); err != nil {
		return "", err
	}
	return hex.EncodeToString(bytes), nil
}

func IsSameUser(u1 accounts.User, u2 accounts.User) bool {
	id := u1.ID == u2.ID
	key := u1.APIKey == u2.APIKey
	name := u1.Name == u2.Name
	giv := u1.GivenName == u2.GivenName
	fam := u1.FamilyName == u2.FamilyName
	pic := u1.PictureURL == u2.PictureURL
	email := u1.Email == u2.Email
	verif := u1.EmailVerified == u2.EmailVerified
	loc := u1.Locale == u2.Locale
	prov := u1.Provider == u2.Provider
	cmd := fmt.Sprint(u1.Cmds) == fmt.Sprint(u2.Cmds)
	team := fmt.Sprint(u1.Teams) == fmt.Sprint(u2.Teams)
	return id && key && name && giv && fam && pic && email && verif && loc && prov && cmd && team
}
