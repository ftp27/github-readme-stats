package handlers

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/ftp27/github-readme-stats/internal/github"
	"github.com/gin-gonic/gin"
	"github.com/spf13/viper"
)

const rateLimitSeconds = 60 * 5

type shieldsResponse struct {
	SchemaVersion int    `json:"schemaVersion"`
	Label         string `json:"label"`
	Message       string `json:"message"`
	Color         string `json:"color"`
	IsError       bool   `json:"isError"`
}

type patError struct {
	Type    string `json:"type"`
	Message string `json:"message"`
}

type patDetail struct {
	Status    string    `json:"status"`
	Remaining *int      `json:"remaining,omitempty"`
	ResetIn   *string   `json:"resetIn,omitempty"`
	Error     *patError `json:"error,omitempty"`
}

type patInfo struct {
	ValidPATs     []string             `json:"validPATs"`
	ExpiredPATs   []string             `json:"expiredPATs"`
	ExhaustedPATs []string             `json:"exhaustedPATs"`
	SuspendedPATs []string             `json:"suspendedPATs"`
	ErrorPATs     []string             `json:"errorPATs"`
	Details       map[string]patDetail `json:"details"`
}

func StatusUp(c *gin.Context) {
	typ := strings.ToLower(c.DefaultQuery("type", "boolean"))
	c.Header("Content-Type", "application/json")

	patsValid := true
	if _, err := github.RetryRateLimit(); err != nil {
		patsValid = false
	}

	if patsValid {
		c.Header("Cache-Control", fmt.Sprintf("max-age=0, s-maxage=%d", rateLimitSeconds))
	} else {
		c.Header("Cache-Control", "no-store")
	}

	switch typ {
	case "shields":
		c.JSON(http.StatusOK, shieldsUptimeBadge(patsValid))
	case "json":
		c.JSON(http.StatusOK, gin.H{"up": patsValid})
	default:
		c.Data(http.StatusOK, "application/json", []byte(strconv.FormatBool(patsValid)))
	}
}

func PATInfo(c *gin.Context) {
	c.Header("Content-Type", "application/json")

	info, err := collectPATInfo()
	if err != nil {
		log.Printf("pat-info error: %v", err)
		c.Header("Cache-Control", "no-store")
		c.String(http.StatusOK, "Something went wrong: "+err.Error())
		return
	}

	if info != nil {
		c.Header("Cache-Control", fmt.Sprintf("max-age=0, s-maxage=%d", rateLimitSeconds))
	}

	payload, err := json.MarshalIndent(info, "", "  ")
	if err != nil {
		c.Header("Cache-Control", "no-store")
		c.String(http.StatusOK, "Something went wrong: "+err.Error())
		return
	}

	c.Data(http.StatusOK, "application/json", payload)
}

func shieldsUptimeBadge(up bool) shieldsResponse {
	label := "Public Instance"
	message := "down"
	color := "red"
	if up {
		message = "up"
		color = "brightgreen"
	}

	return shieldsResponse{
		SchemaVersion: 1,
		Label:         label,
		Message:       message,
		Color:         color,
		IsError:       true,
	}
}

func collectPATInfo() (*patInfo, error) {
	keys := github.GetPATKeys()
	details := map[string]patDetail{}

	for _, key := range keys {
		token := viper.GetString(key)
		data, errors, restErr, err := github.QueryRateLimit(token)
		if err != nil {
			return nil, err
		}

		if restErr != nil {
			msg := strings.ToLower(restErr.Message)
			switch msg {
			case "bad credentials":
				details[key] = patDetail{Status: "expired"}
				continue
			case "sorry. your account was suspended.":
				details[key] = patDetail{Status: "suspended"}
				continue
			default:
				return nil, fmt.Errorf("%s", restErr.Message)
			}
		}

		if len(errors) > 0 && errors[0].Type != "RATE_LIMITED" {
			details[key] = patDetail{
				Status: "error",
				Error: &patError{
					Type:    errors[0].Type,
					Message: errors[0].Message,
				},
			}
			continue
		}

		isRateLimited := false
		if len(errors) > 0 && errors[0].Type == "RATE_LIMITED" {
			isRateLimited = true
		}
		if data.RateLimit.Remaining == 0 {
			isRateLimited = true
		}

		if isRateLimited {
			resetIn := "unknown"
			if data.RateLimit.ResetAt != "" {
				if t, err := time.Parse(time.RFC3339, data.RateLimit.ResetAt); err == nil {
					minutes := int(t.Sub(time.Now()).Minutes() + 0.5)
					resetIn = strconv.Itoa(minutes) + " minutes"
				}
			}
			remaining := 0
			details[key] = patDetail{
				Status:    "exhausted",
				Remaining: &remaining,
				ResetIn:   &resetIn,
			}
			continue
		}

		remaining := data.RateLimit.Remaining
		details[key] = patDetail{
			Status:    "valid",
			Remaining: &remaining,
		}
	}

	info := &patInfo{
		ValidPATs:     filterPATs(details, keys, "valid"),
		ExpiredPATs:   filterPATs(details, keys, "expired"),
		ExhaustedPATs: filterPATs(details, keys, "exhausted"),
		SuspendedPATs: filterPATs(details, keys, "suspended"),
		ErrorPATs:     filterPATs(details, keys, "error"),
		Details:       sortDetails(details, keys),
	}

	return info, nil
}

func filterPATs(details map[string]patDetail, keys []string, status string) []string {
	pats := []string{}
	for _, key := range keys {
		detail, ok := details[key]
		if !ok {
			continue
		}
		if detail.Status == status {
			pats = append(pats, key)
		}
	}
	return pats
}

func sortDetails(details map[string]patDetail, keys []string) map[string]patDetail {
	sorted := map[string]patDetail{}
	for _, key := range keys {
		if detail, ok := details[key]; ok {
			sorted[key] = detail
		}
	}
	return sorted
}
