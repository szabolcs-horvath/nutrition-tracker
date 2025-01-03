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

const (
	Auth0DisabledFlag = "AUTH0_DISABLED"
	TLSDisabledFlag   = "TLS_DISABLED"

	PortKey        = "PORT"
	TLSCertFileKey = "TLS_CERT_FILE"
	TLSKeyFile     = "TLS_KEY_FILE"

	DefaultHttpPort  = "80"
	DefaultHttpsPort = "443"
)

var (
	DefaultCertFile = util.GetPwdSafe() + "/../certs/cert.pem"
	DefaultKeyFile  = util.GetPwdSafe() + "/../certs/key.pem"

	auth0Enabled bool
	tlsEnabled   bool
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

	auth0Enabled = !util.GetEnvFlag(Auth0DisabledFlag)
	tlsEnabled = !util.GetEnvFlag(TLSDisabledFlag)

	if auth0Enabled {
		authenticator, err := util.NewAuthenticator()
		if err != nil {
			panic(fmt.Errorf("couldn't initialize the Authenticator instance: %v", err.Error()))
		}
		util.AuthenticatorInstance = authenticator

		gob.Register(map[string]interface{}{})
		util.CookieStoreInstance = sessions.NewCookieStore([]byte(util.GetEnvSafe("COOKIE_STORE_AUTH_KEY")))
		util.CookieStoreInstance.Options = &sessions.Options{
			Path:     "/",
			Secure:   auth0Enabled,
			HttpOnly: false,
			SameSite: http.SameSiteLaxMode,
		}
	}
}

func getHttpServer(router *http.ServeMux, middlewares []middleware.Middleware) http.Server {
	var port string
	if tlsEnabled {
		port = util.GetEnvOrElse(PortKey, DefaultHttpsPort)
	} else {
		port = util.GetEnvOrElse(PortKey, DefaultHttpPort)
	}

	middlewareStack := middleware.CreateStack(middlewares...)

	return http.Server{
		Addr:    ":" + port,
		Handler: middlewareStack(router),
	}
}

func runServer(server *http.Server) {
	slog.Info("[runServer] Starting server on address " + server.Addr)
	if tlsEnabled {
		certFile := util.GetEnvOrElse(TLSCertFileKey, DefaultCertFile)
		keyFile := util.GetEnvOrElse(TLSKeyFile, DefaultKeyFile)
		if err := server.ListenAndServeTLS(certFile, keyFile); err != nil {
			slog.Info("[runServer] Stopped serving address "+server.Addr, "err", err)
		}
	} else {
		if err := server.ListenAndServe(); err != nil {
			slog.Info("[runServer] Stopped serving address "+server.Addr, "err", err)
		}
	}
	slog.Info("[runServer] Stopped serving new connections on address " + server.Addr)
}

func gracefulShutdown(server *http.Server) {
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)
	<-sigChan

	shutdownCtx, shutdownRelease := context.WithTimeout(context.Background(), 5*time.Second)
	defer shutdownRelease()
	if err := server.Shutdown(shutdownCtx); err != nil {
		slog.Error("[gracefulShutdown] HTTP shutdown error!", "ERROR", err)
	}
	slog.Info("[gracefulShutdown] Graceful shutdown complete.")
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
	if auth0Enabled {
		middlewares = append(middlewares, middleware.IsAuthenticated)
	}
	middlewares = append(middlewares, middleware.LogCompletedRequest)

	server := getHttpServer(router, middlewares)

	go runServer(&server)
	gracefulShutdown(&server)
}
