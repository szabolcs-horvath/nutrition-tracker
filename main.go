package main

import (
	"context"
	"github.com/donseba/go-htmx"
	"github.com/joho/godotenv"
	"github.com/szabolcs-horvath/nutrition-tracker/http_server/middleware"
	"github.com/szabolcs-horvath/nutrition-tracker/http_server/routes"
	"github.com/szabolcs-horvath/nutrition-tracker/http_server/routes/api"
	"github.com/szabolcs-horvath/nutrition-tracker/util"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
)

type App struct {
	HTMX   *htmx.HTMX
	Router *http.ServeMux
}

func main() {
	envFile := getEnvFile()
	if err := godotenv.Load(envFile); err != nil {
		slog.Error("[main] Failed to load .env file!", "FILE", envFile, "ERROR", err)
		panic(1)
	}

	app := &App{
		HTMX:   htmx.New(),
		Router: http.NewServeMux(),
	}

	routes.ServeRoute(app.Router, api.Prefix, api.Routes())

	app.Router.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		h := app.HTMX.NewHandler(w, r)

		page := htmx.NewComponent("templates/index.html")

		_, err := h.Render(r.Context(), page)
		if err != nil {
			slog.ErrorContext(r.Context(), "Error rendering index.html")
		}
	})

	middlewareStack := middleware.CreateStack(
		middleware.AddRequestId,
		middleware.LogIncomingRequest,
		//middleware.Authenticate(util.SafeGetEnv("AUTH0_AUDIENCE"), util.SafeGetEnv("AUTH0_DOMAIN")),
		middleware.LogCompletedRequest,
	)
	server := http.Server{
		Addr:    ":" + util.SafeGetEnv("PORT"),
		Handler: middlewareStack(app.Router),
	}

	go func() {
		slog.Info("[main] Starting server on address " + server.Addr)
		if err := server.ListenAndServe(); err != nil {
			slog.Error("[main] Failed to serve address "+server.Addr, "ERROR", err)
		}
		slog.Info("[main] Stopped serving new connections on address " + server.Addr)
	}()

	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)
	<-sigChan
	shutdownCtx, shutdownRelease := context.WithTimeout(context.Background(), 5*time.Second)
	defer shutdownRelease()
	if err := server.Shutdown(shutdownCtx); err != nil {
		slog.Error("[main] HTTP shutdown error!", "ERROR", err)
	}
	slog.Info("[main] Graceful shutdown complete.")
}

func getEnvFile() string {
	envFile := ".env"
	if len(os.Args[1:]) > 0 {
		if os.Args[1] != "" {
			envFile = os.Args[1]
		}
	}
	slog.Info("[getEnvFile] Using .env file: " + envFile)
	return envFile
}
