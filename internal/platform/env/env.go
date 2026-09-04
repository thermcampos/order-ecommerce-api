package env

import (
	"log/slog"
	"os"
	"unicode/utf8"
)

func GetString(key, fallback string) string {
	slog.Debug("Fetching environment value", "key", key)
	if val := os.Getenv(key); val != "" {
		slog.Debug("Found environment value", "key", key, "length", utf8.RuneCountInString(val))
		return val
	}
	slog.Debug("No environment value. Returning fallback value", "key", key)
	return fallback
}
