package routes

import (
	"fmt"
	"log"
	"net/http"

	"github.com/conalli/bookshelf-backend/controllers"
	"github.com/gorilla/mux"
)

func Search(w http.ResponseWriter, r *http.Request) {
	fmt.Println("Search endpoint hit")
	vars := mux.Vars(r)
	apiKey := vars["apiKey"]
	cmd := vars["cmd"]
	url, err := controllers.GetURL(apiKey, cmd)
	if err != nil {
		log.Println(err)
	}
	http.Redirect(w, r, url, http.StatusTemporaryRedirect)
}
