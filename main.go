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
		slog.Error("[main] Failed to load .env file!", "ERROR", err)
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
	slog.Info("[main] Starting server on address " + server.Addr)
	if err := server.ListenAndServe(); err != nil {
		slog.Error("[main] Failed to serve address "+server.Addr, "ERROR", err)
	}
}
