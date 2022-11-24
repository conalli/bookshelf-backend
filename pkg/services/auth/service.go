package auth

import (
	"context"

	"github.com/conalli/bookshelf-backend/pkg/logs"
	"github.com/go-playground/validator/v10"
)

type Service interface {
	OAuthRedirect(ctx context.Context, code string) (string, error)
}

type service struct {
	l logs.Logger
	v *validator.Validate
}

func NewService(l logs.Logger, v *validator.Validate) *service {
	return &service{l, v}
}

func (s *service) OAuthRedirect(ctx context.Context, code string) (string, error) {
	return "", nil
}
