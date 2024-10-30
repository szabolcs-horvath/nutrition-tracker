package main

import (
	"log/slog"
	"net/http"
	"shorvath/nutrition-tracker/helpers"
	"shorvath/nutrition-tracker/http_server"

	"github.com/joho/godotenv"
)

func main() {
	if err := godotenv.Load(); err != nil {
		slog.Error("Failed to load .env file!")
		panic(1)
	}

	router := http.NewServeMux()
	api := http_server.ServeApiV1Routes(router)

	server := http.Server{
		Addr:    ":" + helpers.SafeGetEnv("PORT"),
		Handler: api,
	}

	slog.Info("Starting server on address " + server.Addr)
	server.ListenAndServe()
}
