package config

import (
	"strings"

	"github.com/spf13/viper"
)

func LoadConfig(cfgPath string, cfgName string) error {
	viper.AddConfigPath(cfgPath)
	viper.SetConfigType("yaml")
	viper.SetConfigName(cfgName)

	viper.AutomaticEnv()
	viper.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))

	return viper.ReadInConfig()
}