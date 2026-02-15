package github

import (
	"encoding/json"
	"errors"
	"log"
	"os"
	"regexp"
	"sort"
	"strings"

	"github.com/spf13/viper"
)

var patKeyPattern = regexp.MustCompile(`^PAT_\d+$`)

func GetPATKeys() []string {
	keys := map[string]bool{}
	for _, entry := range os.Environ() {
		parts := strings.SplitN(entry, "=", 2)
		if len(parts) == 0 {
			continue
		}
		key := parts[0]
		if patKeyPattern.MatchString(key) {
			keys[key] = true
		}
	}

	for _, key := range viper.AllKeys() {
		upper := strings.ToUpper(key)
		if patKeyPattern.MatchString(upper) {
			keys[upper] = true
		}
	}

	result := make([]string, 0, len(keys))
	for key := range keys {
		result = append(result, key)
	}
	sort.Strings(result)
	return result
}

func isRateLimited(errors []GraphQLError) bool {
	if len(errors) == 0 {
		return false
	}
	err := errors[0]
	if err.Type == "RATE_LIMITED" {
		return true
	}
	return strings.Contains(strings.ToLower(err.Message), "rate limit")
}

func isBadCredential(msg string) bool {
	return strings.EqualFold(msg, "Bad credentials")
}

func isSuspended(msg string) bool {
	return strings.EqualFold(msg, "Sorry. Your account was suspended.")
}

func RetryRateLimit() (GraphQLResponse, error) {
	keys := GetPATKeys()
	if len(keys) == 0 {
		return GraphQLResponse{}, errors.New("no github api tokens found")
	}

	for _, key := range keys {
		token := getEnv(key)
		data, gqlErrors, restErr, err := QueryRateLimit(token)
		if err != nil {
			return GraphQLResponse{}, err
		}
		if restErr != nil {
			if isBadCredential(restErr.Message) || isSuspended(restErr.Message) {
				log.Printf("%s Failed", key)
				continue
			}
			return GraphQLResponse{}, errors.New(restErr.Message)
		}
		if isRateLimited(gqlErrors) {
			log.Printf("%s Failed", key)
			continue
		}
		payload, _ := json.Marshal(data)
		return GraphQLResponse{Data: payload, Errors: gqlErrors}, nil
	}

	return GraphQLResponse{}, errors.New("all github api tokens failed")
}

func getEnv(key string) string {
	return viper.GetString(key)
}
