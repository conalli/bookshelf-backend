package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/conalli/bookshelf-backend/auth/jwtauth"
	"github.com/conalli/bookshelf-backend/middleware"
	"github.com/conalli/bookshelf-backend/routes"
	"github.com/gorilla/mux"
	"github.com/joho/godotenv"
)

func main() {
	loadEnv("production")
	router := mux.NewRouter()

	router.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("Hello"))
	}).Methods("GET")
	router.HandleFunc("/signup", routes.SignUp).Methods("POST")
	router.HandleFunc("/login", routes.LogIn).Methods("POST")
	router.HandleFunc("/getcmds/{apiKey}", jwtauth.Authorized(routes.GetCmds)).Methods("GET")
	router.HandleFunc("/setcmd/{apiKey}", jwtauth.Authorized(routes.SetCmd)).Methods("PATCH")
	router.HandleFunc("/delcmd/{apiKey}", jwtauth.Authorized(routes.DelCmd)).Methods("PATCH")
	router.HandleFunc("/delacc/{apiKey}", jwtauth.Authorized(routes.DelAcc)).Methods("DELETE")
	router.HandleFunc("/search/{apiKey}/{cmd}", routes.Search).Methods("GET")

	http.Handle("/", router)

	port := os.Getenv("PORT")
	log.Println("Server up and running on port: " + port)
	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%s", port), middleware.CORSMiddleware(router)))
}

func loadEnv(env string) {
	if env == "production" {
		return
	}
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Could not load .env file")
	}
}
