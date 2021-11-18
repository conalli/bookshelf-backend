package routes

import "net/http"

func GetCmds(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("Hello from getcmds"))
}
