package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"URLShortner/internal/database"
	api_http "URLShortner/internal/http"
)

func main() {
	db, err := database.NewRedisClient()
	if err != nil {
		log.Fatalf("Could not initialize Redis client: %s", err)
	}

	api := &api_http.API{
		Redis: db,
	}

	port := getEnv("APP_PORT", "8080")
	host := getEnv("APP_HOST", "0.0.0.0")

	server := &http.Server{
		Addr:    fmt.Sprintf("%s:%s", host, port),
		Handler: api.Routes(),
	}

	log.Printf("Starting server on %s:%s", host, port)
	if err := server.ListenAndServe(); err != nil {
		log.Fatalf("Server failed to start: %s", err)
	}
}

func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}