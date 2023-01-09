package handlers

import (
	"net/http"

	"github.com/conalli/bookshelf-backend/pkg/errors"
	"github.com/conalli/bookshelf-backend/pkg/http/request"
	"github.com/conalli/bookshelf-backend/pkg/logs"
	"github.com/conalli/bookshelf-backend/pkg/services/auth"
)

func Refresh(a auth.Service, log logs.Logger) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		codeCookie := request.FilterCookies(r.Cookies(), auth.BookshelfTokenCode)
		accessCookie := request.FilterCookies(r.Cookies(), auth.BookshelfAccessToken)
		if codeCookie == nil || accessCookie == nil {
			log.Error("incorrect cookies in request")
			errors.APIErrorResponse(w, errors.NewBadRequestError("incorrect information in request"))
			return
		}
		access := accessCookie.Value
		code := codeCookie.Value
		tokens, apierr := a.RefreshTokens(r.Context(), access, code)
		if apierr != nil {
			log.Error("could not refresh tokens")
			errors.APIErrorResponse(w, apierr)
			return
		}
		cookies := tokens.NewTokenCookies(log)
		auth.AddCookiesToResponse(w, cookies)
		w.WriteHeader(http.StatusOK)
	}
}
