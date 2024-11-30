package helpers

import (
	"encoding/json"
	"log/slog"
	"net/http"
	"os"
)

func SafeGetEnv(key string) string {
	if os.Getenv(key) == "" {
		slog.Error("[SafeGetEnv] The environment variable '" + key + "' is not set.")
		panic(1)
	}
	return os.Getenv(key)
}

func WriteJson(w http.ResponseWriter, status int, data interface{}) error {
	js, err := json.Marshal(data)
	if err != nil {
		return err
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	if _, err = w.Write(js); err != nil {
		return err
	}
	return nil
}
