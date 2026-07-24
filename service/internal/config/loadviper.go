package config

import (
	"strings"

	"github.com/spf13/viper"
)

func LoadViper() {

	viper.SetDefault("json-logging", false)

	viper.SetEnvKeyReplacer(strings.NewReplacer("-", "_"))
	viper.AutomaticEnv()
}
