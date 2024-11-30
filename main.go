package main

import (
	"github.com/donseba/go-htmx"
	"github.com/joho/godotenv"
	"log/slog"
	"net/http"
	"shorvath/nutrition-tracker/http_server/middleware"
	"shorvath/nutrition-tracker/http_server/routes"
	"shorvath/nutrition-tracker/http_server/routes/api"
	"shorvath/nutrition-tracker/util"
)

type App struct {
	HTMX   *htmx.HTMX
	Router *http.ServeMux
}

func main() {
	if err := godotenv.Load(); err != nil {
		slog.Error("[main] Failed to load .env file!", "ERROR", err)
		panic(1)
	}

	app := &App{
		HTMX:   htmx.New(),
		Router: http.NewServeMux(),
	}

	routes.ServeRoute(app.Router, api.Prefix, api.Routes())

	middlewareStack := middleware.CreateStack(
		middleware.AddRequestId,
		middleware.Log,
	)
	server := http.Server{
		Addr:    ":" + util.SafeGetEnv("PORT"),
		Handler: middlewareStack(app.Router),
	}
	slog.Info("[main] Starting server on address " + server.Addr)
	if err := server.ListenAndServe(); err != nil {
		slog.Error("[main] Failed to serve address "+server.Addr, "ERROR", err)
	}
}
