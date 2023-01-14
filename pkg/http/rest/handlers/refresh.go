package handlers

import (
	"net/http"

	"github.com/conalli/bookshelf-backend/pkg/apierr"
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
			apierr.APIErrorResponse(w, apierr.NewBadRequestError("incorrect information in request"))
			return
		}
		access := accessCookie.Value
		code := codeCookie.Value
		tokens, apiErr := a.RefreshTokens(r.Context(), access, code)
		if apiErr != nil {
			log.Error("could not refresh tokens")
			apierr.APIErrorResponse(w, apiErr)
			return
		}
		cookies := tokens.NewTokenCookies(log, http.SameSiteLaxMode)
		auth.AddCookiesToResponse(w, cookies)
		w.WriteHeader(http.StatusOK)
	}
}
