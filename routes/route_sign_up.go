package routes

import "net/http"

func SignUp(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("Hello from signup"))
}
