package middleware

import (
	"bytes"
	"context"
	"github.com/auth0/go-jwt-middleware/v2"
	"github.com/auth0/go-jwt-middleware/v2/jwks"
	"github.com/auth0/go-jwt-middleware/v2/validator"
	"github.com/google/uuid"
	"io"
	"log/slog"
	"net/http"
	"net/url"
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

func Authenticate(audience string, domain string) Middleware {
	return func(next http.Handler) http.Handler {
		issuerURL, _ := url.Parse("https://" + domain + "/")
		provider := jwks.NewCachingProvider(issuerURL, 5*time.Minute)

		jwtValidator, _ := validator.New(
			provider.KeyFunc,
			validator.RS256,
			issuerURL.String(),
			[]string{audience},
		)

		middleware := jwtmiddleware.New(jwtValidator.ValidateToken)
		return middleware.CheckJWT(next)
	}
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
