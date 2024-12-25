package main

import (
	"context"
	"github.com/joho/godotenv"
	"github.com/szabolcs-horvath/nutrition-tracker/http_server/middleware"
	"github.com/szabolcs-horvath/nutrition-tracker/http_server/routes"
	"github.com/szabolcs-horvath/nutrition-tracker/http_server/routes/api"
	"github.com/szabolcs-horvath/nutrition-tracker/repository"
	"github.com/szabolcs-horvath/nutrition-tracker/util"
	"html/template"
	"io"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
)

type Templates struct {
	templates *template.Template
}

type App struct {
	Router    *http.ServeMux
	Templates *Templates
}

func (app App) Render(w io.Writer, name string, data interface{}) error {
	return app.Templates.templates.ExecuteTemplate(w, name, data)
}

func newApp() *App {
	return &App{
		Router: http.NewServeMux(),
		Templates: &Templates{
			templates: template.Must(template.New("templates").Funcs(util.TemplateFuncs()).ParseGlob("templates/*.html")),
		},
	}
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

type IndexData struct {
	MealLogs []*repository.MealLog
	Items    []*repository.Item
}

func main() {
	envFile := getEnvFile()
	if err := godotenv.Load(envFile); err != nil {
		slog.Error("[main] Failed to load .env file!", "FILE", envFile, "ERROR", err)
		panic(1)
	}

	app := newApp()
	routes.ServeRoute(app.Router, api.Prefix, api.Routes())

	fs := http.FileServer(http.Dir("web/static/vendor"))
	app.Router.Handle("/static/", http.StripPrefix("/static/", fs))

	app.Router.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		mealLogs, err := repository.FindMealLogsForUserAndDate(r.Context(), 1)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		items, err := repository.ListItems(r.Context())
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		data := IndexData{
			MealLogs: mealLogs,
			Items:    items,
		}
		err = app.Render(w, "index", data)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
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
