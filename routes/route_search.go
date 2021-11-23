package routes

import (
	"log"
	"net/http"

	"github.com/conalli/bookshelf-backend/controllers"
	"github.com/gorilla/mux"
)

func Search(w http.ResponseWriter, r *http.Request) {
	log.Println("Search endpoint hit")
	vars := mux.Vars(r)
	apiKey := vars["apiKey"]
	cmd := vars["cmd"]
	url, err := controllers.GetURL(apiKey, cmd)
	if err != nil {
		log.Printf("%d: %s", err.Status(), err.Error())
	}
	http.Redirect(w, r, url, http.StatusTemporaryRedirect)
}
