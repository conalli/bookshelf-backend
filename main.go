package main

import (
	"fmt"
	"net/http"

	"github.com/gorilla/mux"
)

func main() {
	router := mux.NewRouter()
	
	router.HandleFunc("/", HealthCheck).Methods("GET")
	router.HandleFunc("/hello", HelloWorld).Methods("GET")

	http.Handle("/", router)
	http.ListenAndServe(":8080", router)
}

func HealthCheck(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)

		fmt.Fprintf(w, "Server is up and running on port 8080")
	}

	func HelloWorld(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)

		w.Write([]byte("Hello World"))
	}