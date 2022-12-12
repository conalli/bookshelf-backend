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
		APIKey, ok := request.GetAPIKeyFromContext(r)
		if len(APIKey) < 1 || !ok {
			log.Error("could not get APIKey from context")
			errors.APIErrorResponse(w, errors.NewInternalServerError())
			return
		}
		code := request.FilterCookies(r.Cookies(), auth.BookshelfTokenCode)
		if code == nil {
			log.Error("no code cookie in request")
			errors.APIErrorResponse(w, errors.NewBadRequestError("incorrect information in request"))
			return
		}
		tokens, apierr := a.RefreshTokens(r.Context(), APIKey, code.Value)
		if apierr != nil {
			log.Error("could not refresh tokens")
			errors.APIErrorResponse(w, apierr)
			return
		}
		cookies := tokens.NewTokenCookies(log)
		log.Info("successfully returned token as cookie")
		auth.AddCookiesToResponse(w, cookies)
		w.WriteHeader(http.StatusOK)
	}
}
