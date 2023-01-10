package handlers

import (
	"net/http"

	"github.com/conalli/bookshelf-backend/pkg/logs"
	"github.com/conalli/bookshelf-backend/pkg/services/auth"
)

func LogOut(a auth.Service, log logs.Logger) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		// TODO: delete refresh token
		auth.RemoveBookshelfCookies(w)
		w.WriteHeader(http.StatusOK)
	}
}
