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

func Reduce[T, V any](slice []T, fn func(T, V), acc V) V {
	for _, t := range slice {
		fn(t, acc)
	}
	return acc
}

func GroupBy[T any, K comparable](slice []T, keyFn func(T) K) map[K][]T {
	return Reduce(slice, func(t T, acc map[K][]T) {
		key := keyFn(t)
		acc[key] = append(acc[key], t)
	}, map[K][]T{})
}

func GroupByKeys[T any, K comparable](slice []T, keys []K, keyFn func(T) K) map[K][]T {
	result := make(map[K][]T, len(keys))
	for _, k := range keys {
		result[k] = []T{}
	}
	return Reduce(slice, func(t T, acc map[K][]T) {
		key := keyFn(t)
		acc[key] = append(acc[key], t)
	}, result)
}

func FindFirst[T any](slice []T, predicate func(T) bool) (T, bool) {
	for _, t := range slice {
		if predicate(t) {
			return t, true
		}
	}
	var backup T
	return backup, false
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
		"percentage":          Percentage,
		"percentageRemaining": PercentageRemaining,
		"subtractInt64": func(a, b int64) int64 {
			return a - b
		},
		"subtractFloat64": func(a, b float64) int64 {
			return int64(a - b)
		},
		"mapOf": func(values ...interface{}) (map[string]interface{}, error) {
			if len(values)%2 != 0 {
				return nil, fmt.Errorf("dict expects an even number of arguments")
			}
			d := make(map[string]interface{})
			for i := 0; i < len(values); i += 2 {
				key, ok := values[i].(string)
				if !ok {
					return nil, fmt.Errorf("dict keys must be strings")
				}
				d[key] = values[i+1]
			}
			return d, nil
		},
		"formatFloat": func(value any, precision int) any {
			var floatValue float64
			switch value.(type) {
			case float64:
				floatValue = value.(float64)
			case *float64:
				floatValue = *value.(*float64)
			default:
				return value
			}
			p := math.Pow(10, float64(precision))
			rounded := math.Round(floatValue*p) / p

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

func Percentage(a, b float64) int64 {
	return int64(math.Min(a/b*100, 100))
}

func PercentageRemaining(a, b float64) int64 {
	return int64(math.Max(100-(a/b*100), 0))
}
