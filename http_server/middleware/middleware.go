package middleware

import (
	"context"
	"github.com/google/uuid"
	"log/slog"
	"net/http"
	"slices"
	"time"
)

type Middleware func(handler http.Handler) http.Handler

func CreateStack(middlewares ...Middleware) Middleware {
	return func(next http.Handler) http.Handler {
		for _, middleware := range slices.Backward(middlewares) {
			next = middleware(next)
		}
		return next
	}
}

type wrappedWriter struct {
	http.ResponseWriter
	statusCode int
}

func (w *wrappedWriter) WriteHeader(statusCode int) {
	w.ResponseWriter.WriteHeader(statusCode)
	w.statusCode = statusCode
}

const requestIdKey = "x-request_id"

func AddRequestId(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ctx := context.WithValue(r.Context(), requestIdKey, uuid.New().String())
		req := r.WithContext(ctx)
		next.ServeHTTP(w, req)
	})
}

func Log(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()

		wrappedW := &wrappedWriter{
			ResponseWriter: w,
			statusCode:     http.StatusOK,
		}

		next.ServeHTTP(wrappedW, r)

		requestId, ok := r.Context().Value(requestIdKey).(string)
		if !ok {
			requestId = "_missing_"
		}

		slog.Info("Execution time for the request was: "+time.Since(start).String(),
			"REQUEST_ID", requestId,
			"METHOD", r.Method,
			"PATH", r.URL.Path,
			"STATUS", wrappedW.statusCode)
	})
}
