package handlers

import (
	"net/http"

	"github.com/conalli/bookshelf-backend/pkg/logs"
	"github.com/conalli/bookshelf-backend/pkg/services/search"
	"github.com/gorilla/mux"
)

// Search takes the APIKey and cmd route variables and redirects the user to the url
// associated with the cmd or to a google search of the cmd if no url can be found.
func Search(s search.Service, log logs.Logger) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		log.Info("Search endpoint hit")
		vars := mux.Vars(r)
		APIKey := vars["APIKey"]
		cmd := vars["cmd"]
		url, err := s.Search(r.Context(), APIKey, cmd)
		if err != nil {
			log.Errorf("could not find cmd: %v", err)
		}
		http.Redirect(w, r, url, http.StatusSeeOther)
	}
}
