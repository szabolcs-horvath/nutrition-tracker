package main

import (
	"context"
	"github.com/joho/godotenv"
	"github.com/szabolcs-horvath/nutrition-tracker/http_server/middleware"
	"github.com/szabolcs-horvath/nutrition-tracker/http_server/routes"
	"github.com/szabolcs-horvath/nutrition-tracker/http_server/routes/api"
	"github.com/szabolcs-horvath/nutrition-tracker/http_server/routes/htmx"
	"github.com/szabolcs-horvath/nutrition-tracker/util"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
)

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

func main() {
	envFile := getEnvFile()
	if err := godotenv.Load(envFile); err != nil {
		slog.Error("[main] Failed to load .env file!", "FILE", envFile, "ERROR", err)
		panic(1)
	}

	router := http.NewServeMux()
	routes.ServeRoute(router, api.Prefix, api.Routes())
	routes.ServeRouteHandlers(router, htmx.Prefix, htmx.Routes())
	routes.ServeFS(router, "/static", "web/static/vendor")

	middlewareStack := middleware.CreateStack(
		middleware.AddRequestId,
		middleware.LogIncomingRequest,
		//middleware.Authenticate(util.SafeGetEnv("AUTH0_AUDIENCE"), util.SafeGetEnv("AUTH0_DOMAIN")),
		middleware.LogCompletedRequest,
	)
	server := http.Server{
		Addr:    ":" + util.SafeGetEnv("PORT"),
		Handler: middlewareStack(router),
	}

	go func() {
		slog.Info("[main] Starting server on address " + server.Addr)
		if err := server.ListenAndServe(); err != nil {
			slog.Info("[main] Stopped serving address "+server.Addr, "err", err)
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
