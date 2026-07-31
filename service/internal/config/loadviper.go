package config

import (
	"strings"

	"github.com/spf13/viper"
)

const (
	envPrefix         = "ANTHOLOGISE"
	keyLogJSON        = "json-logging"
	keyLogLevel       = "log-level"
	keyAppStremioID   = "stremio-id"
	keyAppVersion     = "version-number"
	keyAppName        = "app-name"
	keyAppDescription = "manifest-description"
	keyAppLogoURL     = "logo-url"
)

// AppConfig contains all the app related options
type AppConfig struct {
	// StremioID is the ID for the add-on used by the manifest
	StremioID string
	// Version is the current app version
	Version string
	// Name is the name of the add-on used by the manifest
	Name string
	// Description is the description of the add-on used by the manifest
	Description string
	// LogoURL is the URL to the add-on's logo which will be used to display it in Stremio
	LogoURL string
}

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

func loadApp(v *viper.Viper) AppConfig {
	return AppConfig{
		v.GetString(keyAppStremioID),
		v.GetString(keyAppVersion),
		v.GetString(keyAppName),
		v.GetString(keyAppDescription),
		v.GetString(keyAppLogoURL),
	}
}

func loadLog(v *viper.Viper) LogConfig {
	return LogConfig{
		v.GetBool(keyLogJSON),
		v.GetString(keyLogLevel),
	}
}

func loadDB(v *viper.Viper) DBConfig {
	return DBConfig{}
}

func setDefaults(v *viper.Viper) {

	v.SetDefault(keyAppLogoURL, "https://raw.githubusercontent.com/Stremio/stremio-web/development/assets/images/stremio_symbol.png")

	v.SetDefault(keyLogJSON, false)
	v.SetDefault(keyLogLevel, "warn")
}

// LoadViper loads variables into viper and sets defaults for the service.
func LoadViper() Config {

	v := viper.New()

	setDefaults(v)

	v.SetEnvPrefix(envPrefix)
	v.SetEnvKeyReplacer(strings.NewReplacer("-", "_"))
	v.AutomaticEnv()

	return Config{
		loadApp(v),
		loadLog(v),
		loadDB(v),
	}
}
