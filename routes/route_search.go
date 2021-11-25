package routes

import (
	"log"
	"net/http"

	"github.com/conalli/bookshelf-backend/controllers"
	"github.com/gorilla/mux"
)

// Search takes the apiKey and cmd route variables and redirects the user to the url
// associated with the cmd or to a google search of the cmd if no url can be found.
func Search(w http.ResponseWriter, r *http.Request) {
	log.Println("Search endpoint hit")
	vars := mux.Vars(r)
	apiKey := vars["apiKey"]
	cmd := vars["cmd"]
	url, err := controllers.GetURL(apiKey, cmd)
	if err != nil {
		log.Printf("%d: %s", err.Status(), err.Error())
	}
	http.Redirect(w, r, url, http.StatusSeeOther)
}
