package main

import (
	"log/slog"
	"net/http"
	"shorvath/nutrition-tracker/helpers"
	"shorvath/nutrition-tracker/http_server"
	"shorvath/nutrition-tracker/http_server/middleware"

	"github.com/joho/godotenv"
)

func main() {
	if err := godotenv.Load(); err != nil {
		slog.Error("Failed to load .env file!")
		panic(1)
	}

	router := http.NewServeMux()
	api := routes.ServeApiV1Routes(router)
	middlewareStack := middleware.CreateStack(middleware.Log)

	server := http.Server{
		Addr:    ":" + helpers.SafeGetEnv("PORT"),
		Handler: middlewareStack(api),
	}

	slog.Info("Starting server on address " + server.Addr)
	server.ListenAndServe()
}
