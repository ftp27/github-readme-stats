package github

import (
	"encoding/json"
	"errors"
	"net/http"
	"time"
)

var restClient = &http.Client{Timeout: 10 * time.Second}

type RestResponse struct {
	TotalCount int `json:"total_count"`
}

func SearchCommits(username string, token string) (RestResponse, *RestError, int, error) {
	req, err := http.NewRequest("GET", "https://api.github.com/search/commits?q=author:"+username, nil)
	if err != nil {
		return RestResponse{}, nil, 0, err
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Accept", "application/vnd.github.cloak-preview")
	if token != "" {
		req.Header.Set("Authorization", "token "+token)
	}

	resp, err := restClient.Do(req)
	if err != nil {
		return RestResponse{}, nil, 0, err
	}
	defer resp.Body.Close()

	if resp.StatusCode >= 400 {
		var restErr RestError
		if err := json.NewDecoder(resp.Body).Decode(&restErr); err != nil {
			return RestResponse{}, nil, resp.StatusCode, errors.New("github api error")
		}
		return RestResponse{}, &restErr, resp.StatusCode, nil
	}

	var data RestResponse
	if err := json.NewDecoder(resp.Body).Decode(&data); err != nil {
		return RestResponse{}, nil, resp.StatusCode, err
	}

	return data, nil, resp.StatusCode, nil
}
