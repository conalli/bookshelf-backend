package auth

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"os"

	"github.com/conalli/bookshelf-backend/pkg/logs"
	"github.com/go-playground/validator/v10"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/endpoints"
)

type Service interface {
	OAuthFlow(ctx context.Context, code string) (string, error)
}

type service struct {
	l logs.Logger
	v *validator.Validate
}

func NewService(l logs.Logger, v *validator.Validate) *service {
	return &service{l, v}
}

func (s *service) OAuthFlow(ctx context.Context, code string) (string, error) {
	url := fmt.Sprintf("%s?client_id=%s&client_secret=%s&code=%s&grant_type=authorization_code&redirect_uri=http://localhost:8080/api/oauth/redirect", endpoints.Google.TokenURL, os.Getenv("GOOGLE_CLIENT_ID"), os.Getenv("GOOGLE_CLIENT_SECRET"), code)
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, url, nil)
	if err != nil {
		s.l.Error(err)
		return "", err
	}
	req.Header.Set("accept", "application/json")
	res, err := http.DefaultClient.Do(req)
	if err != nil {
		s.l.Error(err)
		return "", err
	}
	defer res.Body.Close()
	var t oauth2.Token
	json.NewDecoder(res.Body).Decode(&t)

	s.l.Infof("%+v", t)
	return "", nil
}
