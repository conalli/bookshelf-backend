package routes

import (
	"encoding/json"
	"log"
	"net/http"
	"time"

	"github.com/conalli/bookshelf-backend/controllers"
	"github.com/conalli/bookshelf-backend/models"
	"github.com/conalli/bookshelf-backend/utils/apiErrors"
	"github.com/conalli/bookshelf-backend/utils/auth/jwtauth"
)

func LogIn(w http.ResponseWriter, r *http.Request) {
	log.Println("LogIn endpoint hit")
	var logInReq models.Credentials
	json.NewDecoder(r.Body).Decode(&logInReq)

	name, err := controllers.CheckCredentials(logInReq)
	if err != nil {
		log.Printf("error returned while trying to get check credentials: %v", err)
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(err.Status())
		checkCredErr := apiErrors.ResError{
			Status: err.Status(),
			Error:  err.Error(),
		}
		json.NewEncoder(w).Encode(checkCredErr)
		return
	}
	var token string
	token, err = jwtauth.NewToken(name)
	if err != nil {
		log.Printf("error returned while trying to create a new token: %v", err)
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(err.Status())
		newTknErr := apiErrors.ResError{
			Status: err.Status(),
			Error:  err.Error(),
		}
		json.NewEncoder(w).Encode(newTknErr)
		return
	}
	cookie := http.Cookie{
		Name:     "jwt",
		Value:    token,
		Expires:  time.Now().Add(15 * time.Minute),
		HttpOnly: true,
	}
	log.Println("successfully returned token as cookie")
	http.SetCookie(w, &cookie)
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	res := struct{ status string }{
		status: "success",
	}
	json.NewEncoder(w).Encode(res)
	return
}
