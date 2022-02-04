package handlers

import (
	"encoding/json"
	"log"
	"net/http"
	"time"

	"github.com/conalli/bookshelf-backend/pkg/errors"
	"github.com/conalli/bookshelf-backend/pkg/jwtauth"
	"github.com/conalli/bookshelf-backend/pkg/user"
)

// SignUp is the handler for the signup endpoint. Checks db for username and if
// unique adds new user with given credentials.
func SignUp(u user.Service) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		log.Println("SignUp endpoint hit")
		var newUserReq user.SignUpRequest
		json.NewDecoder(r.Body).Decode(&newUserReq)
		newUser, err := u.NewUser(r.Context(), newUserReq)
		if err != nil {
			log.Printf("error returned while trying to create a new user: %v", err)
			errors.APIErrorResponse(w, err)
			return
		}
		log.Printf("successfully created a new user: %+v", newUser)
		var token string
		token, err = jwtauth.NewToken(newUser.APIKey)
		if err != nil {
			log.Printf("error returned while trying to create a new token: %v", err)
			errors.APIErrorResponse(w, err)
			return
		}
		// Use Secure during production.
		cookie := http.Cookie{
			Name:     "bookshelfjwt",
			Value:    token,
			Expires:  time.Now().Add(30 * time.Minute),
			HttpOnly: true,
			// Secure:   true,
			SameSite: http.SameSiteNoneMode,
		}
		log.Println("successfully returned token as cookie")
		http.SetCookie(w, &cookie)
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		res := newUser
		json.NewEncoder(w).Encode(res)
		return
	}
}
