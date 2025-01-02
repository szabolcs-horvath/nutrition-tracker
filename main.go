package main

import (
	"context"
	"encoding/gob"
	"fmt"
	"github.com/gorilla/sessions"
	"github.com/joho/godotenv"
	"github.com/szabolcs-horvath/nutrition-tracker/http_server/middleware"
	"github.com/szabolcs-horvath/nutrition-tracker/http_server/routes"
	"github.com/szabolcs-horvath/nutrition-tracker/http_server/routes/api"
	"github.com/szabolcs-horvath/nutrition-tracker/http_server/routes/auth"
	"github.com/szabolcs-horvath/nutrition-tracker/http_server/routes/htmx"
	"github.com/szabolcs-horvath/nutrition-tracker/util"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
)

func init() {
	envFile := ".env"
	if len(os.Args[1:]) > 0 {
		if os.Args[1] != "" {
			envFile = os.Args[1]
		}
	}
	slog.Info("[getEnvFile] Using .env file: " + envFile)

	if err := godotenv.Load(envFile); err != nil {
		slog.Error("[main] Failed to load .env file!", "FILE", envFile, "ERROR", err)
		panic(1)
	}

	if os.Getenv("AUTH0_DISABLED") != "true" {
		authenticator, err := util.NewAuthenticator()
		if err != nil {
			panic(fmt.Errorf("couldn't initialize the Authenticator instance: %v", err.Error()))
		}
		util.AuthenticatorInstance = authenticator

		gob.Register(map[string]interface{}{})
		util.CookieStoreInstance = sessions.NewCookieStore([]byte(util.SafeGetEnv("SESSION_KEY")))
		util.CookieStoreInstance.Options = &sessions.Options{
			Path:     "/",
			Secure:   false,
			HttpOnly: false,
			SameSite: http.SameSiteLaxMode,
		}
	}
}

func main() {
	router := http.NewServeMux()

	router.HandleFunc("/{$}", routes.RootHandler)
	routes.ServeRoute(router, api.Prefix, api.Routes())
	routes.ServeRouteHandlers(router, auth.Prefix, auth.Routes())
	routes.ServeRouteHandlers(router, htmx.Prefix, htmx.Routes())
	routes.ServeFS(router, "/static", "web/static/vendor")

	middlewares := make([]middleware.Middleware, 0)
	middlewares = append(middlewares, middleware.AddRequestId)
	middlewares = append(middlewares, middleware.LogIncomingRequest)
	if os.Getenv("AUTH0_DISABLED") != "true" {
		middlewares = append(middlewares, middleware.IsAuthenticated)
	}
	middlewares = append(middlewares, middleware.LogCompletedRequest)
	middlewareStack := middleware.CreateStack(middlewares...)

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
