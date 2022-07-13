package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/conalli/bookshelf-backend/pkg/db/mongodb"
	"github.com/conalli/bookshelf-backend/pkg/http/rest"
	"github.com/go-playground/validator/v10"
	"github.com/joho/godotenv"
	"go.uber.org/zap"
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
	loadEnv("development")
	logger, err := zap.NewProduction()
	if err != nil {
		log.Fatalf("Couldn't make a new logger, %v", err)
	}
	r := rest.NewRouter(logger.Sugar(), validator.New(), mongodb.New()).Walk().HandlerWithCORS()

	port := os.Getenv("PORT")
	log.Println("Server up and running on port: " + port)
	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%s", port), r))
}
