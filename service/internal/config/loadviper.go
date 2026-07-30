package config

import (
	"strings"

	"github.com/spf13/viper"
)

const (
	keyLogJSON  = "json-logging"
	keyLogLevel = "log-level"
)

// AppConfig contains all the app related options
type AppConfig struct{}

// LogConfig contains all the logging options
type LogConfig struct {
	// FormatJSON indicates if the logs should written as JSON; will be as text otherwise
	FormatJSON bool
	// LogLevel is the level that the project will log at, e.g: Error, Warn, etc.
	Level string
}

// DBConfig contains all the database options
type DBConfig struct{}

// Config contains all the runtime configuration for the project
type Config struct {
	App AppConfig
	Log LogConfig
	DB  DBConfig
}

func loadApp() AppConfig {
	return AppConfig{}
}

func loadLog() LogConfig {
	return LogConfig{
		viper.GetBool(keyLogJSON),
		viper.GetString(keyLogLevel),
	}
}

func loadDB() DBConfig {
	return DBConfig{}
}

func setDefaults() {

	viper.SetDefault(keyLogJSON, false)
	viper.SetDefault(keyLogLevel, "warn")
}

// LoadViper loads variables into viper and sets defaults for the service.
func LoadViper() Config {

	setDefaults()

	viper.SetEnvKeyReplacer(strings.NewReplacer("-", "_"))
	viper.AutomaticEnv()

	return Config{
		loadApp(),
		loadLog(),
		loadDB(),
	}
}
