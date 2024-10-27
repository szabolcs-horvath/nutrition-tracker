package main

import (
	"fmt"
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
	port := helpers.SafeGetEnv("PORT")

	for k, v := range http_server.GetRoutes() {
		http.HandleFunc(k, v)
	}

	slog.Error("Exiting: %v", http.ListenAndServe(fmt.Sprintf(":%s", port), nil))
}
