package handlers

import (
	"net/http"
	"os"

	"github.com/conalli/bookshelf-backend/pkg/errors"
	"github.com/conalli/bookshelf-backend/pkg/http/request"
	"github.com/conalli/bookshelf-backend/pkg/logs"
	"github.com/conalli/bookshelf-backend/pkg/services/search"
	"github.com/gorilla/mux"
)

// Search takes the APIKey and cmd route variables and redirects the user to the url
// associated with the cmd or to a google search of the cmd if no url can be found.
func Search(s search.Service, log logs.Logger) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		APIKey, ok := request.GetAPIKeyFromContext(r)
		if len(APIKey) < 1 || !ok {
			log.Error("could not get APIKey from context")
			errors.APIErrorResponse(w, errors.NewInternalServerError())
			return
		}
		vars := mux.Vars(r)
		args := vars["args"]
		log.Info(args)
		url, err := s.Search(r.Context(), APIKey, args)
		if err != nil {
			log.Errorf("could not find cmd: %v", err)
			errURL := os.Getenv("ALLOWED_URL_BASE") + "/webcli/error"
			http.Redirect(w, r, errURL, http.StatusSeeOther)
		}
		http.Redirect(w, r, url, http.StatusSeeOther)
	}
}
