package github

import (
	"bytes"
	"encoding/json"
	"errors"
	"net/http"
	"time"
)

type GraphQLError struct {
	Type    string `json:"type"`
	Message string `json:"message"`
}

type RateLimit struct {
	Remaining int    `json:"remaining"`
	ResetAt   string `json:"resetAt"`
}

type RateLimitData struct {
	RateLimit RateLimit `json:"rateLimit"`
}

type GraphQLResponse struct {
	Data   json.RawMessage `json:"data"`
	Errors []GraphQLError  `json:"errors"`
}

type RestError struct {
	Message string `json:"message"`
}

const githubGraphQLEndpoint = "https://api.github.com/graphql"

var httpClient = &http.Client{Timeout: 10 * time.Second}

func QueryRateLimit(token string) (RateLimitData, []GraphQLError, *RestError, error) {
	query := `query {
    rateLimit {
      remaining
      resetAt
    }
  }`

	body, err := json.Marshal(map[string]any{
		"query": query,
	})
	if err != nil {
		return RateLimitData{}, nil, nil, err
	}

	req, err := http.NewRequest("POST", githubGraphQLEndpoint, bytes.NewReader(body))
	if err != nil {
		return RateLimitData{}, nil, nil, err
	}

	req.Header.Set("Authorization", "bearer "+token)
	req.Header.Set("Content-Type", "application/json")

	resp, err := httpClient.Do(req)
	if err != nil {
		return RateLimitData{}, nil, nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode >= 400 {
		var restErr RestError
		if err := json.NewDecoder(resp.Body).Decode(&restErr); err != nil {
			return RateLimitData{}, nil, nil, errors.New("github api error")
		}
		return RateLimitData{}, nil, &restErr, nil
	}

	var gqlResp GraphQLResponse
	if err := json.NewDecoder(resp.Body).Decode(&gqlResp); err != nil {
		return RateLimitData{}, nil, nil, err
	}

	var data RateLimitData
	if err := json.Unmarshal(gqlResp.Data, &data); err != nil {
		return RateLimitData{}, gqlResp.Errors, nil, err
	}

	return data, gqlResp.Errors, nil, nil
}
