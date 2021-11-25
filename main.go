package main

import (
	"log"
	"net/http"
	"os"

	"github.com/conalli/bookshelf-backend/routes"
	"github.com/gorilla/mux"
	"github.com/joho/godotenv"
)

func main() {
	loadEnv()
	router := mux.NewRouter()

	router.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("Hello"))
	}).Methods("GET")
	router.HandleFunc("/signup", routes.SignUp).Methods("POST")
	router.HandleFunc("/login", routes.LogIn).Methods("POST")
	router.HandleFunc("/setcmd", routes.SetCmd).Methods("PUT")
	router.HandleFunc("/getcmds", routes.GetCmds).Methods("GET")
	router.HandleFunc("/search/{apiKey}/{cmd}", routes.Search).Methods("GET")

	http.Handle("/", router)
	port := os.Getenv("PORT")
	log.Println("Server up and running on port" + port)
	http.ListenAndServe(port, router)
}

func loadEnv() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Could not load .env file")
	}
}
