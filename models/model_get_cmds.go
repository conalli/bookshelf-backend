package models

type GetCmdsReq struct {
	Name     string `json:"name"`
	Password string `json:"password"`
}
