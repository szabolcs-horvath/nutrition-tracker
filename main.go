package main

import (
	"github.com/joho/godotenv"
	"log/slog"
	"net/http"
	"shorvath/nutrition-tracker/helpers"
	"shorvath/nutrition-tracker/http_server/middleware"
	"shorvath/nutrition-tracker/http_server/routes"
	"shorvath/nutrition-tracker/http_server/routes/api"
)

func main() {
	if err := godotenv.Load(); err != nil {
		slog.Error("Failed to load .env file!")
		panic(1)
	}

	router := http.NewServeMux()
	routes.ServeRoute(router, api.Prefix, api.Routes())

	middlewareStack := middleware.CreateStack(
		middleware.AddRequestId,
		middleware.Log,
	)

	server := http.Server{
		Addr:    ":" + helpers.SafeGetEnv("PORT"),
		Handler: middlewareStack(router),
	}

	slog.Info("Starting server on address " + server.Addr)
	server.ListenAndServe()
}
