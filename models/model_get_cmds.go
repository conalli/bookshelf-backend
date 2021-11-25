package models

// GetCmdsReq represents the fields required in the request in order to return
// all of the users cmds.
type GetCmdsReq struct {
	Name     string `json:"name"`
	Password string `json:"password"`
}
