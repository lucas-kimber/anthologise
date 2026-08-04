package config

import (
	"log/slog"
	"os"
)

func parseLogLevel(strLevel string) (slog.Level, bool) {

	var logLevel slog.Level
	ok := true

	if err := logLevel.UnmarshalText([]byte(strLevel)); err != nil {
		logLevel = slog.LevelInfo
		ok = false
	}
	return logLevel, ok
}

func buildLogger(jsonLogging bool, opts *slog.HandlerOptions) *slog.Logger {

	var handler slog.Handler = slog.NewTextHandler(os.Stdout, opts)

	if jsonLogging {
		handler = slog.NewJSONHandler(os.Stdout, opts)
	}

	return slog.New(handler)
}

// ConfigureSlog takes a configuration struct and returns a constructed slog logger.
//
// The log-level setting controls the minimum enabled logging level.
// This can be set to debug, info, warn, or error. Invalid values default to info.
//
// When json-logging is true, logs use slog's JSON handler. Otherwise, they use
// the text handler.
func ConfigureSlog(cfg LogConfig) *slog.Logger {

	logLevel, lvlOk := parseLogLevel(cfg.Level)

	lvl := &slog.HandlerOptions{
		Level: logLevel,
	}
	logger := buildLogger(cfg.FormatJSON, lvl)

	if !lvlOk {
		logger.Warn("No valid logging level was passed, defaulted to INFO")
	}

	logger.Info("Logger successfully initialised")

	return logger
}
