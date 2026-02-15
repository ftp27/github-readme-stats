package common

import (
	"bytes"
	"fmt"
	"html/template"
	"strings"
)

const ErrorCardLength = 576.5

var upstreamApiErrors = map[string]bool{
	TryAgainLater:            true,
	SecondaryErrors.MaxRetry: true,
}

func RenderError(message string, secondary string, options RenderOptions) (string, error) {
	colors, err := GetCardColors(ColorsArgs{
		TitleColor:  options.TitleColor,
		TextColor:   options.TextColor,
		BgColor:     options.BgColor,
		BorderColor: options.BorderColor,
		Theme:       options.Theme,
	})
	if err != nil {
		return "", err
	}

	showRepoLink := true
	if options.ShowRepoLink != nil {
		showRepoLink = *options.ShowRepoLink
	}

	issueLink := ""
	if showRepoLink && !upstreamApiErrors[secondary] {
		issueLink = " file an issue at https://tiny.one/readme-stats"
	}

	tmpl := template.Must(template.New("error").Parse(`
<svg width="{{.Width}}" height="120" viewBox="0 0 {{.Width}} 120" fill="{{.BgColor}}" xmlns="http://www.w3.org/2000/svg">
  <style>
    .text { font: 600 16px 'Segoe UI', Ubuntu, Sans-Serif; fill: {{.TitleColor}} }
    .small { font: 600 12px 'Segoe UI', Ubuntu, Sans-Serif; fill: {{.TextColor}} }
    .gray { fill: #858585 }
  </style>
  <rect x="0.5" y="0.5" width="{{.RectWidth}}" height="99%" rx="4.5" fill="{{.BgColor}}" stroke="{{.BorderColor}}"/>
  <text x="25" y="45" class="text">Something went wrong!{{.IssueLink}}</text>
  <text data-testid="message" x="25" y="55" class="text small">
    <tspan x="25" dy="18">{{.Message}}</tspan>
    <tspan x="25" dy="18" class="gray">{{.Secondary}}</tspan>
  </text>
</svg>`))

	var buf bytes.Buffer
	data := map[string]any{
		"Width":       ErrorCardLength,
		"RectWidth":   ErrorCardLength - 1,
		"BgColor":     colors.BgColor,
		"BorderColor": colors.BorderColor,
		"TitleColor":  colors.TitleColor,
		"TextColor":   colors.TextColor,
		"IssueLink":   issueLink,
		"Message":     message,
		"Secondary":   secondary,
	}
	if err := tmpl.Execute(&buf, data); err != nil {
		return "", err
	}
	return buf.String(), nil
}

func FlexLayout(items []string, gap int, direction string, sizes []int) []string {
	lastSize := 0
	result := []string{}
	for i, item := range items {
		if strings.TrimSpace(item) == "" {
			continue
		}
		size := 0
		if len(sizes) > i {
			size = sizes[i]
		}
		transform := fmt.Sprintf("translate(%d, 0)", lastSize)
		if direction == "column" {
			transform = fmt.Sprintf("translate(0, %d)", lastSize)
		}
		lastSize += size + gap
		result = append(result, fmt.Sprintf("<g transform=\"%s\">%s</g>", transform, item))
	}
	return result
}

func MeasureText(value string, fontSize float64) float64 {
	widths := []float64{
		0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0,
		0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0,
		0, 0, 0, 0, 0.2796875, 0.2765625,
		0.3546875, 0.5546875, 0.5546875, 0.8890625, 0.665625, 0.190625,
		0.3328125, 0.3328125, 0.3890625, 0.5828125, 0.2765625, 0.3328125,
		0.2765625, 0.3015625, 0.5546875, 0.5546875, 0.5546875, 0.5546875,
		0.5546875, 0.5546875, 0.5546875, 0.5546875, 0.5546875, 0.5546875,
		0.2765625, 0.2765625, 0.584375, 0.5828125, 0.584375, 0.5546875,
		1.0140625, 0.665625, 0.665625, 0.721875, 0.721875, 0.665625,
		0.609375, 0.7765625, 0.721875, 0.2765625, 0.5, 0.665625,
		0.5546875, 0.8328125, 0.721875, 0.7765625, 0.665625, 0.7765625,
		0.721875, 0.665625, 0.609375, 0.721875, 0.665625, 0.94375,
		0.665625, 0.665625, 0.609375, 0.2765625, 0.3546875, 0.2765625,
		0.4765625, 0.5546875, 0.3328125, 0.5546875, 0.5546875, 0.5,
		0.5546875, 0.5546875, 0.2765625, 0.5546875, 0.5546875, 0.221875,
		0.240625, 0.5, 0.221875, 0.8328125, 0.5546875, 0.5546875,
		0.5546875, 0.5546875, 0.3328125, 0.5, 0.2765625, 0.5546875,
		0.5, 0.721875, 0.5, 0.5, 0.5, 0.3546875, 0.259375, 0.353125, 0.5890625,
	}

	avg := 0.5279276315789471
	sum := 0.0
	for _, c := range value {
		idx := int(c)
		if idx >= 0 && idx < len(widths) {
			sum += widths[idx]
		} else {
			sum += avg
		}
	}

	return sum * fontSize
}

type RenderOptions struct {
	TitleColor   string
	TextColor    string
	BgColor      string
	BorderColor  string
	Theme        string
	ShowRepoLink *bool
}
