package common

import (
	"strconv"
	"strings"
)

func ParseBoolean(value string) *bool {
	if value == "" {
		return nil
	}
	lower := strings.ToLower(value)
	if lower == "true" {
		v := true
		return &v
	}
	if lower == "false" {
		v := false
		return &v
	}
	return nil
}

func ParseArray(value string) []string {
	if value == "" {
		return []string{}
	}
	parts := strings.Split(value, ",")
	out := make([]string, 0, len(parts))
	for _, part := range parts {
		if part == "" {
			continue
		}
		out = append(out, part)
	}
	return out
}

func ClampValue(number, min, max float64) float64 {
	if number != number {
		return min
	}
	if number < min {
		return min
	}
	if number > max {
		return max
	}
	return number
}

func ParseInt(value string) (int, bool) {
	if value == "" {
		return 0, false
	}
	parsed, err := strconv.Atoi(value)
	if err != nil {
		return 0, false
	}
	return parsed, true
}

func ParseFloat(value string) (float64, bool) {
	if value == "" {
		return 0, false
	}
	parsed, err := strconv.ParseFloat(value, 64)
	if err != nil {
		return 0, false
	}
	return parsed, true
}

func LowercaseTrim(value string) string {
	return strings.ToLower(strings.TrimSpace(value))
}
