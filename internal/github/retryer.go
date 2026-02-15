package github

import (
	"errors"
	"log"
	"strings"
)

type Fetcher func(token string) (GraphQLResponse, *RestError, error)

type RetryConfig struct {
	IsRateLimited func(response GraphQLResponse) bool
}

func RetryGraphQL(fetcher Fetcher, config RetryConfig) (GraphQLResponse, error) {
	keys := GetPATKeys()
	if len(keys) == 0 {
		return GraphQLResponse{}, errors.New("no github api tokens found")
	}

	for _, key := range keys {
		token := getToken(key)
		resp, restErr, err := fetcher(token)
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

		if config.IsRateLimited != nil && config.IsRateLimited(resp) {
			log.Printf("%s Failed", key)
			continue
		}

		return resp, nil
	}

	return GraphQLResponse{}, errors.New("all github api tokens failed")
}

func getToken(key string) string {
	return strings.TrimSpace(getEnv(key))
}
