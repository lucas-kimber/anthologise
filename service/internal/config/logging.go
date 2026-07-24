package config

import (
	"log/slog"
	"os"

	"github.com/spf13/viper"
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

// ConfigureSlog configures the default slog logger using Viper settings.
//
// The log-level setting controls the minimum enabled logging level.
// This can be set to debug, info, warn, or error. Invalid values default to info.
//
// When json-logging is true, logs use slog's JSON handler. Otherwise, they use
// the text handler.
func ConfigureSlog() {

	jsonLogging := viper.GetBool("json-logging")
	lvl := viper.GetString("log-level")

	logLevel, lvlOk := parseLogLevel(lvl)

	opts := &slog.HandlerOptions{
		Level: logLevel,
	}
	logger := buildLogger(jsonLogging, opts)

	slog.SetDefault(logger)

	if !lvlOk {
		slog.Warn("No valid logging level was passed, defaulted to INFO")
	}

	slog.Info("Logger Successfully Initialised")
}
