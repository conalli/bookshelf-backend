package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/conalli/bookshelf-backend/auth/jwtauth"
	"github.com/conalli/bookshelf-backend/controllers"
	"github.com/conalli/bookshelf-backend/middleware"
	"github.com/gorilla/mux"
	"github.com/joho/godotenv"
)

func main() {
	loadEnv("development")
	router := mux.NewRouter()

	router.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("Hello"))
	}).Methods("GET")
	router.HandleFunc("/signup", controllers.SignUp).Methods("POST")
	router.HandleFunc("/login", controllers.LogIn).Methods("POST")
	router.HandleFunc("/getcmds/{apiKey}", jwtauth.Authorized(controllers.GetCmds)).Methods("GET")
	router.HandleFunc("/setcmd/{apiKey}", jwtauth.Authorized(controllers.SetCmd)).Methods("PATCH")
	router.HandleFunc("/delcmd/{apiKey}", jwtauth.Authorized(controllers.DelCmd)).Methods("PATCH")
	router.HandleFunc("/delacc/{apiKey}", jwtauth.Authorized(controllers.DelUser)).Methods("DELETE")
	router.HandleFunc("/search/{apiKey}/{cmd}", controllers.Search).Methods("GET")

	router.HandleFunc("/team/{apiKey}", jwtauth.Authorized(controllers.NewTeam)).Methods("POST")
	router.HandleFunc("/team/addmember/{apiKey}", jwtauth.Authorized(controllers.AddMember)).Methods("PUT")

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
