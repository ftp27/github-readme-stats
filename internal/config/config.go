package config

import (
	"strings"

	"github.com/spf13/viper"
)

func Load() {
	viper.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))
	viper.AutomaticEnv()

	viper.SetDefault("PORT", "9000")
	viper.SetDefault("FETCH_MULTI_PAGE_STARS", "false")

	viper.AddConfigPath(".")
	viper.AddConfigPath("..")

	viper.SetConfigFile(".env")
	viper.SetConfigType("env")
	if err := viper.ReadInConfig(); err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); !ok {
			// Ignore missing config file.
		}
	}

	mergeConfig("config", "yaml")
	mergeConfig("config", "yml")
	mergeConfig("config", "json")
}

func mergeConfig(name string, configType string) {
	viper.SetConfigName(name)
	viper.SetConfigType(configType)
	if err := viper.MergeInConfig(); err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); !ok {
			// Ignore missing config file.
		}
	}
}
