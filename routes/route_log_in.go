package routes

import "net/http"

func LogIn(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("Hello from login"))
}
