package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/conalli/bookshelf-backend/pkg/errors"
	"github.com/conalli/bookshelf-backend/pkg/http/request"
	"github.com/conalli/bookshelf-backend/pkg/logs"
	"github.com/conalli/bookshelf-backend/pkg/services/accounts"
)

// AddCmdResponse represents the data returned upon successfully adding a cmd.
type AddCmdResponse struct {
	NumAdded int    `json:"num_added"`
	Cmd      string `json:"cmd"`
	URL      string `json:"url"`
}

// AddCmd is the handler for the setcmd endpoint. Checks credentials + JWT and if
// authorized sets new cmd.
func AddCmd(u accounts.UserService, log logs.Logger) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		APIKey, ok := request.GetAPIKeyFromContext(r)
		if len(APIKey) < 1 || !ok {
			log.Error("could not get APIKey from context")
			errors.APIErrorResponse(w, errors.NewInternalServerError())
			return
		}
		setCmdReq, parseErr := request.DecodeJSONRequest[request.AddCmd](r.Body)
		if parseErr != nil {
			errRes := errors.NewBadRequestError("could not parse request body")
			errors.APIErrorResponse(w, errRes)
		}
		numUpdated, err := u.AddCmd(r.Context(), setCmdReq, APIKey)
		if err != nil {
			log.Errorf("error returned while trying to add a new cmd: %v", err)
			errors.APIErrorResponse(w, err)
			return
		}
		if numUpdated == 0 {
			log.Errorf("could not update cmds... maybe %s:%s already exists?", setCmdReq.Cmd, setCmdReq.URL)
			err := errors.NewBadRequestError("error: could not update cmds")
			errors.APIErrorResponse(w, err)
			return
		}
		log.Infof("successfully set cmd: %s, url: %s", setCmdReq.Cmd, setCmdReq.URL)
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		res := AddCmdResponse{
			NumAdded: numUpdated,
			Cmd:      setCmdReq.Cmd,
			URL:      setCmdReq.URL,
		}
		json.NewEncoder(w).Encode(res)
	}
}
