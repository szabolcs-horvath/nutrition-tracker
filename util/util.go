package util

import (
	"encoding/json"
	"fmt"
	"html/template"
	"log/slog"
	"math"
	"net/http"
	"os"
	"strconv"
	"strings"
)

func Map[T, V any](slice []T, fn func(T) V) []V {
	result := make([]V, len(slice))
	for i, t := range slice {
		result[i] = fn(t)
	}
	return result
}

func SafeGetEnv(key string) string {
	if os.Getenv(key) == "" {
		slog.Error("[SafeGetEnv] The environment variable '" + key + "' is not set.")
		panic(1)
	}
	return os.Getenv(key)
}

func ReadJson(r *http.Request, target interface{}) error {
	return json.NewDecoder(r.Body).Decode(target)
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

func TemplateFuncs() template.FuncMap {
	return template.FuncMap{
		"formatFloat": func(value float64, precision int) string {
			p := math.Pow(10, float64(precision))
			rounded := math.Round(value*p) / p

			formatted := fmt.Sprintf("%."+strconv.Itoa(precision)+"f", rounded)

			if strings.Contains(formatted, ".") {
				formatted = strings.TrimRight(formatted, "0")
				if formatted[len(formatted)-1] == '.' {
					formatted = formatted[:len(formatted)-1]
				}
			}

			return formatted
		},
	}
}
