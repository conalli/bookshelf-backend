package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/conalli/bookshelf-backend/pkg/http/rest"
	"github.com/joho/godotenv"
)

func loadEnv(env string) {
	if env == "production" {
		return
	}
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Could not load .env file")
	}
}

func main() {
	loadEnv("production")

	router := rest.RouterWithCORS(true)

	port := os.Getenv("PORT")
	log.Println("Server up and running on port: " + port)
	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%s", port), router))
}
