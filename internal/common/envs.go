package common

import "github.com/spf13/viper"

func Whitelist() []string {
	if value := viper.GetString("WHITELIST"); value != "" {
		return SplitComma(value)
	}
	return nil
}

func GistWhitelist() []string {
	if value := viper.GetString("GIST_WHITELIST"); value != "" {
		return SplitComma(value)
	}
	return nil
}

func ExcludeRepositories() []string {
	if value := viper.GetString("EXCLUDE_REPO"); value != "" {
		return SplitComma(value)
	}
	return []string{}
}

func SplitComma(value string) []string {
	if value == "" {
		return []string{}
	}
	return ParseArray(value)
}
