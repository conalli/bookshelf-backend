package routes

import "net/http"

func SetCmd(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("Hello from setcmd"))
}
