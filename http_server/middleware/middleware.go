package middleware

import (
	"bytes"
	"context"
	"github.com/google/uuid"
	"github.com/szabolcs-horvath/nutrition-tracker/util"
	"io"
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
	statusCode   int
	responseBody []byte
}

func (w *wrappedWriter) WriteHeader(statusCode int) {
	w.ResponseWriter.WriteHeader(statusCode)
	w.statusCode = statusCode
}

func (w *wrappedWriter) Write(data []byte) (int, error) {
	w.responseBody = data
	return w.ResponseWriter.Write(data)
}

func findRequestId(r *http.Request) string {
	requestId, ok := r.Context().Value(requestIdKey).(string)
	if !ok {
		requestId = "_missing_"
	}
	return requestId
}

const requestIdKey = "x-request_id"

func AddRequestId(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ctx := context.WithValue(r.Context(), requestIdKey, uuid.New().String())
		req := r.WithContext(ctx)
		next.ServeHTTP(w, req)
	})
}

func LogIncomingRequest(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		keyValues := make([]any, 0)
		keyValues = append(keyValues, "METHOD", r.Method)
		keyValues = append(keyValues, "PATH", r.URL.Path)
		keyValues = append(keyValues, "REQUEST_ID", findRequestId(r))

		body, _ := io.ReadAll(r.Body)
		r.Body = io.NopCloser(bytes.NewReader(body))
		if !bytes.Equal(body, []byte("")) {
			keyValues = append(keyValues, "BODY", body)
		}

		slog.InfoContext(r.Context(), "[LogIncomingRequest]", keyValues...)
		next.ServeHTTP(w, r)
	})
}

func IsAuthenticated(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		session, err := util.CookieStoreInstance.Get(r, "auth-session")
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		if session.Values["profile"] == nil {
			if r.URL.Path == "/auth/login" || r.URL.Path == "/auth/callback" {
				slog.Warn("[IsAuthenticated] session.Values[\"profile\"] is nil", "PATH", r.URL.Path)
			} else {
				slog.Info("[IsAuthenticated] redirecting to the login path...", "PATH", r.URL.Path)
				http.Redirect(w, r, "/auth/login", http.StatusSeeOther)
				return
			}
		}
		next.ServeHTTP(w, r)
	})
}

func LogCompletedRequest(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		keyValues := make([]any, 0)
		keyValues = append(keyValues, "METHOD", r.Method)
		keyValues = append(keyValues, "PATH", r.URL.Path)

		wrappedW := &wrappedWriter{
			ResponseWriter: w,
			statusCode:     http.StatusOK,
		}
		start := time.Now()
		next.ServeHTTP(wrappedW, r)
		duration := time.Since(start)

		keyValues = append(keyValues, "STATUS", wrappedW.statusCode)
		keyValues = append(keyValues, "DURATION", duration.String())
		keyValues = append(keyValues, "REQUEST_ID", findRequestId(r))
		keyValues = append(keyValues, "BODY", string(wrappedW.responseBody))

		slog.InfoContext(r.Context(), "[LogCompletedRequest]", keyValues...)
	})
}
