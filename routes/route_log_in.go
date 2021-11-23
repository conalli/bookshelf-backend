package routes

import "net/http"

func LogIn(w http.ResponseWriter, r *http.Request) {
	// TODO: add log in with JWT
	w.Write([]byte("Hello from login"))
}
