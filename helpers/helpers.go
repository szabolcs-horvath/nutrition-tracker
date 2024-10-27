package helpers

import (
	"log/slog"
	"os"
)

func SafeGetEnv(key string) string {
	if os.Getenv(key) == "" {
		slog.Error("The environment variable '%s' is not set.", key)
		panic(1)
	}
	return os.Getenv(key)
}
