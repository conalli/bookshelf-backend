package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/conalli/bookshelf-backend/pkg/db/mongodb"
	"github.com/conalli/bookshelf-backend/pkg/db/redis"
	"github.com/conalli/bookshelf-backend/pkg/http/rest"
	"github.com/coreos/go-oidc/v3/oidc"
	"github.com/go-playground/validator/v10"
	"github.com/joho/godotenv"
	"go.uber.org/zap"
)

func loadEnv(env string) error {
	if env == "production" {
		return nil
	}
	return godotenv.Load()
}

func main() {
	logger, err := zap.NewProduction()
	if err != nil {
		log.Fatalf("Couldn't make a new logger, %v", err)
	}
	defer logger.Sync()
	if err = loadEnv("production"); err != nil {
		log.Fatal("Could not load .env file")
	}
	sugar := logger.Sugar()
	provider, err := oidc.NewProvider(context.Background(), "https://accounts.google.com")
	if err != nil {
		log.Fatal("Could not make oidc provider:", err)
	}
	r := rest.NewRouter(sugar, validator.New(), mongodb.New(sugar), redis.NewClient(sugar), provider).Walk().HandlerWithCORS()
	port := os.Getenv("PORT")
	log.Println("Server up and running on port: " + port)
	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%s", port), r))
}
