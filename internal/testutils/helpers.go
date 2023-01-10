package testutils

import (
	"bytes"
	"crypto/rand"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/conalli/bookshelf-backend/pkg/http/request"
	"github.com/conalli/bookshelf-backend/pkg/logs"
	"github.com/conalli/bookshelf-backend/pkg/services/accounts"
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
	req.AddCookie(code)
	req.AddCookie(access)
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
