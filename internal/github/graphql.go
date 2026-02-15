package github

import (
	"bytes"
	"encoding/json"
	"errors"
	"net/http"
	"time"
)

var graphqlClient = &http.Client{Timeout: 10 * time.Second}

func DoGraphQL(query string, variables map[string]any, token string) (GraphQLResponse, *RestError, error) {
	body, err := json.Marshal(map[string]any{
		"query":     query,
		"variables": variables,
	})
	if err != nil {
		return GraphQLResponse{}, nil, err
	}

	req, err := http.NewRequest("POST", githubGraphQLEndpoint, bytes.NewReader(body))
	if err != nil {
		return GraphQLResponse{}, nil, err
	}

	req.Header.Set("Authorization", "bearer "+token)
	req.Header.Set("Content-Type", "application/json")

	resp, err := graphqlClient.Do(req)
	if err != nil {
		return GraphQLResponse{}, nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode >= 400 {
		var restErr RestError
		if err := json.NewDecoder(resp.Body).Decode(&restErr); err != nil {
			return GraphQLResponse{}, nil, errors.New("github api error")
		}
		return GraphQLResponse{}, &restErr, nil
	}

	var gqlResp GraphQLResponse
	if err := json.NewDecoder(resp.Body).Decode(&gqlResp); err != nil {
		return GraphQLResponse{}, nil, err
	}

	return gqlResp, nil, nil
}
