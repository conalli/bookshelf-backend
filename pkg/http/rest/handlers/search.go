package handlers

import (
	"log"
	"net/http"

	"github.com/conalli/bookshelf-backend/pkg/search"
	"github.com/gorilla/mux"
)

// Search takes the APIKey and cmd route variables and redirects the user to the url
// associated with the cmd or to a google search of the cmd if no url can be found.
func Search(s search.Service) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		log.Println("Search endpoint hit")
		vars := mux.Vars(r)
		APIKey := vars["APIKey"]
		cmd := vars["cmd"]
		url, err := s.Search(r.Context(), APIKey, cmd)
		if err != nil {
			log.Printf("%v", err)
			// log.Printf("%d: %s", err.Status(), err.Error())
		}
		http.Redirect(w, r, url, http.StatusSeeOther)
		return
	}
}
