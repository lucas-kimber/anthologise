package config

import (
	"strings"

	"github.com/spf13/viper"
)

// LoadViper loads variables into viper and sets defaults for the service.
func LoadViper() {

	viper.SetDefault("json-logging", false)

	viper.SetEnvKeyReplacer(strings.NewReplacer("-", "_"))
	viper.AutomaticEnv()
}
