package cards

import (
	"bytes"
	"fmt"
	"html/template"
	"math"
	"strings"

	"github.com/ftp27/github-readme-stats/internal/common"
	"github.com/ftp27/github-readme-stats/internal/fetchers"
)

const (
	cardMinWidth             = 287
	cardDefaultWidth         = 287
	rankCardMinWidth         = 420
	rankCardDefaultWidth     = 450
	rankOnlyCardMinWidth     = 290
	rankOnlyCardDefaultWidth = 290
)

var longLocales = map[string]bool{
	"az":      true,
	"bg":      true,
	"cs":      true,
	"de":      true,
	"el":      true,
	"es":      true,
	"fil":     true,
	"fi":      true,
	"fr":      true,
	"hu":      true,
	"id":      true,
	"ja":      true,
	"ml":      true,
	"my":      true,
	"nl":      true,
	"pl":      true,
	"pt-br":   true,
	"pt-pt":   true,
	"ru":      true,
	"sr":      true,
	"sr-latn": true,
	"sw":      true,
	"ta":      true,
	"uk-ua":   true,
	"uz":      true,
	"zh-tw":   true,
}

type StatCardOptions struct {
	Hide              []string
	ShowIcons         bool
	HideTitle         bool
	HideBorder        bool
	CardWidth         int
	HideRank          bool
	IncludeAllCommits bool
	CommitsYear       *int
	LineHeight        int
	TitleColor        string
	RingColor         string
	IconColor         string
	TextColor         string
	TextBold          bool
	BgColor           string
	Theme             string
	CustomTitle       string
	BorderRadius      float64
	BorderColor       string
	NumberFormat      string
	NumberPrecision   *int
	Locale            string
	DisableAnimations bool
	RankIcon          string
	Show              []string
}

type statTextItem struct {
	Icon          template.HTML
	Label         string
	Value         string
	ID            string
	UnitSymbol    string
	Index         int
	ShowIcons     bool
	ShiftValuePos float64
	Bold          bool
}

func RenderStatsCard(stats fetchers.StatsData, options StatCardOptions) (string, error) {
	lineHeight := options.LineHeight
	if lineHeight == 0 {
		lineHeight = 25
	}

	colors, err := common.GetCardColors(common.ColorsArgs{
		TitleColor:  options.TitleColor,
		TextColor:   options.TextColor,
		IconColor:   options.IconColor,
		BgColor:     options.BgColor,
		BorderColor: options.BorderColor,
		RingColor:   options.RingColor,
		Theme:       options.Theme,
	})
	if err != nil {
		return "", err
	}

	apostrophe := "'"
	if strings.HasSuffix(strings.ToLower(strings.TrimSpace(stats.Name)), "s") {
		apostrophe = ""
	}

	locale := strings.ToLower(options.Locale)
	translations := map[string]map[string]string{}
	for key, value := range common.StatCardLocales(stats.Name, apostrophe) {
		translations[key] = value
	}
	for key, value := range common.WakatimeCardLocales() {
		translations[key] = value
	}

	i18n := common.NewI18n(locale, translations)

	statMap := map[string]statTextItem{}
	statMap["stars"] = statTextItem{
		Icon:  template.HTML(common.Icons.Star),
		Label: i18n.T("statcard.totalstars"),
		Value: fmt.Sprintf("%d", stats.TotalStars),
		ID:    "stars",
	}

	commitsLabel := i18n.T("statcard.commits") + totalCommitsYearLabel(options.IncludeAllCommits, options.CommitsYear, i18n)
	statMap["commits"] = statTextItem{
		Icon:  template.HTML(common.Icons.Commits),
		Label: commitsLabel,
		Value: fmt.Sprintf("%d", stats.TotalCommits),
		ID:    "commits",
	}

	statMap["prs"] = statTextItem{
		Icon:  template.HTML(common.Icons.PRs),
		Label: i18n.T("statcard.prs"),
		Value: fmt.Sprintf("%d", stats.TotalPRs),
		ID:    "prs",
	}

	if contains(options.Show, "prs_merged") {
		statMap["prs_merged"] = statTextItem{
			Icon:  template.HTML(common.Icons.PRsMerged),
			Label: i18n.T("statcard.prs-merged"),
			Value: fmt.Sprintf("%d", stats.TotalPRsMerged),
			ID:    "prs_merged",
		}
	}

	if contains(options.Show, "prs_merged_percentage") {
		precision := clampPrecision(options.NumberPrecision)
		statMap["prs_merged_percentage"] = statTextItem{
			Icon:       template.HTML(common.Icons.PRsMergedPercentage),
			Label:      i18n.T("statcard.prs-merged-percentage"),
			Value:      fmt.Sprintf("%.*f", precision, stats.MergedPRsPercentage),
			ID:         "prs_merged_percentage",
			UnitSymbol: "%",
		}
	}

	if contains(options.Show, "reviews") {
		statMap["reviews"] = statTextItem{
			Icon:  template.HTML(common.Icons.Reviews),
			Label: i18n.T("statcard.reviews"),
			Value: fmt.Sprintf("%d", stats.TotalReviews),
			ID:    "reviews",
		}
	}

	statMap["issues"] = statTextItem{
		Icon:  template.HTML(common.Icons.Issues),
		Label: i18n.T("statcard.issues"),
		Value: fmt.Sprintf("%d", stats.TotalIssues),
		ID:    "issues",
	}

	if contains(options.Show, "discussions_started") {
		statMap["discussions_started"] = statTextItem{
			Icon:  template.HTML(common.Icons.DiscussionsStarted),
			Label: i18n.T("statcard.discussions-started"),
			Value: fmt.Sprintf("%d", stats.TotalDiscussionsStarted),
			ID:    "discussions_started",
		}
	}

	if contains(options.Show, "discussions_answered") {
		statMap["discussions_answered"] = statTextItem{
			Icon:  template.HTML(common.Icons.DiscussionsAnswered),
			Label: i18n.T("statcard.discussions-answered"),
			Value: fmt.Sprintf("%d", stats.TotalDiscussionsAnswered),
			ID:    "discussions_answered",
		}
	}

	statMap["contribs"] = statTextItem{
		Icon:  template.HTML(common.Icons.Contribs),
		Label: i18n.T("statcard.contribs"),
		Value: fmt.Sprintf("%d", stats.ContributedTo),
		ID:    "contribs",
	}

	statItems := []string{}
	index := 0
	for _, key := range orderStatKeys() {
		if contains(options.Hide, key) {
			continue
		}
		item, ok := statMap[key]
		if !ok {
			continue
		}
		item.Index = index
		item.ShowIcons = options.ShowIcons
		item.ShiftValuePos = 79.01
		if longLocales[locale] {
			item.ShiftValuePos += 50
		}
		item.Bold = options.TextBold
		item.Value = applyNumberFormat(item.ID, item.Value, options.NumberFormat, options.NumberPrecision)
		statItems = append(statItems, renderStatItem(item))
		index++
	}

	if len(statItems) == 0 && options.HideRank {
		return "", common.NewCustomError("Could not render stats card.", "Either stats or rank are required.")
	}

	baseHeight := int(math.Max(45+float64((len(statItems)+1)*lineHeight), 0))
	if !options.HideRank {
		if len(statItems) > 0 {
			baseHeight = int(math.Max(float64(baseHeight), 150))
		} else {
			baseHeight = int(math.Max(float64(baseHeight), 180))
		}
	}

	renderHeight := baseHeight
	if options.HideTitle {
		renderHeight -= 30
		if renderHeight < 0 {
			renderHeight = 0
		}
	}

	progress := 100 - stats.Rank.Percentile
	cssStyles := getStyles(colors.TextColor, colors.IconColor, colors.RingColor, options.ShowIcons, progress)

	titleText := i18n.T("statcard.title")
	if len(statItems) == 0 {
		titleText = i18n.T("statcard.ranktitle")
	}
	if options.CustomTitle != "" {
		titleText = options.CustomTitle
	}

	iconWidth := 0
	if options.ShowIcons && len(statItems) > 0 {
		iconWidth = 17
	}

	minCardWidth := 0
	if options.HideRank {
		minCardWidth = int(common.ClampValue(50+common.MeasureText(titleText, 10)*2, cardMinWidth, math.Inf(1))) + iconWidth
	} else if len(statItems) > 0 {
		minCardWidth = rankCardMinWidth + iconWidth
	} else {
		minCardWidth = rankOnlyCardMinWidth + iconWidth
	}

	defaultCardWidth := 0
	if options.HideRank {
		defaultCardWidth = cardDefaultWidth + iconWidth
	} else if len(statItems) > 0 {
		defaultCardWidth = rankCardDefaultWidth + iconWidth
	} else {
		defaultCardWidth = rankOnlyCardDefaultWidth + iconWidth
	}

	width := options.CardWidth
	if width == 0 {
		width = defaultCardWidth
	}
	if width < minCardWidth {
		width = minCardWidth
	}

	rankCircle := ""
	if !options.HideRank {
		rankCircle = fmt.Sprintf(`<g data-testid="rank-circle" transform="translate(%d, %d)">
      <circle class="rank-circle-rim" cx="-10" cy="8" r="40" />
      <circle class="rank-circle" cx="-10" cy="8" r="40" />
      <g class="rank-text">%s</g>
    </g>`, calculateRankXTranslation(width, minCardWidth, len(statItems) > 0, iconWidth), baseHeight/2-50, common.RankIcon(options.RankIcon, stats.Rank.Level, stats.Rank.Percentile))
	}

	labels := buildAccessibilityLabels(statMap, options.Hide, options.IncludeAllCommits, options.CommitsYear, i18n)

	gradientDefs := ""
	bgFill := colors.BgColor
	if len(colors.BgGradient) > 0 {
		gradientDefs = renderGradient(colors.BgGradient)
		bgFill = "url(#gradient)"
	}

	body := rankCircle + "<svg x=\"0\" y=\"0\">" + strings.Join(common.FlexLayout(statItems, lineHeight, "column", nil), "") + "</svg>"

	card := cardTemplateData{
		Width:        width,
		Height:       renderHeight,
		BorderRadius: options.BorderRadius,
		Title:        titleText,
		HideBorder:   options.HideBorder,
		HideTitle:    options.HideTitle,
		CSS:          template.CSS(cssStyles),
		BgFill:       bgFill,
		BorderColor:  colors.BorderColor,
		GradientDefs: template.HTML(gradientDefs),
		TitleColor:   colors.TitleColor,
		Animations:   !options.DisableAnimations,
		A11yTitle:    titleText + ", Rank: " + stats.Rank.Level,
		A11yDesc:     labels,
		Body:         template.HTML(body),
	}

	return renderCard(card)
}

type cardTemplateData struct {
	Width        int
	Height       int
	BorderRadius float64
	Title        string
	HideBorder   bool
	HideTitle    bool
	CSS          template.CSS
	BgFill       string
	BorderColor  string
	GradientDefs template.HTML
	TitleColor   string
	Animations   bool
	A11yTitle    string
	A11yDesc     string
	Body         template.HTML
}

func renderCard(data cardTemplateData) (string, error) {
	tmpl := template.Must(template.New("card").Parse(`
<svg width="{{.Width}}" height="{{.Height}}" viewBox="0 0 {{.Width}} {{.Height}}" fill="none" xmlns="http://www.w3.org/2000/svg" role="img" aria-labelledby="descId">
  <title id="titleId">{{.A11yTitle}}</title>
  <desc id="descId">{{.A11yDesc}}</desc>
  <style>
    .header {
      font: 600 18px 'Segoe UI', Ubuntu, Sans-Serif;
      fill: {{.TitleColor}};
      animation: fadeInAnimation 0.8s ease-in-out forwards;
    }
    @supports(-moz-appearance: auto) {
      .header { font-size: 15.5px; }
    }
    {{.CSS}}
    @keyframes scaleInAnimation {
      from { transform: translate(-5px, 5px) scale(0); }
      to { transform: translate(-5px, 5px) scale(1); }
    }
    @keyframes fadeInAnimation {
      from { opacity: 0; }
      to { opacity: 1; }
    }
    {{if not .Animations}}* { animation-duration: 0s !important; animation-delay: 0s !important; }{{end}}
  </style>
  {{.GradientDefs}}
  <rect data-testid="card-bg" x="0.5" y="0.5" rx="{{.BorderRadius}}" height="99%" stroke="{{.BorderColor}}" width="{{.WidthMinusOne}}" fill="{{.BgFill}}" stroke-opacity="{{.BorderOpacity}}" />
  {{if not .HideTitle}}
  <g data-testid="card-title" transform="translate(25, 35)">
    <text x="0" y="0" class="header" data-testid="header">{{.Title}}</text>
  </g>
  {{end}}
  <g data-testid="main-card-body" transform="translate(0, 0)">
    {{.Body}}
  </g>
</svg>`))

	payload := map[string]any{
		"Width":         data.Width,
		"WidthMinusOne": data.Width - 1,
		"Height":        data.Height,
		"BorderRadius":  data.BorderRadius,
		"Title":         data.Title,
		"HideTitle":     data.HideTitle,
		"CSS":           data.CSS,
		"GradientDefs":  data.GradientDefs,
		"BgFill":        data.BgFill,
		"BorderColor":   data.BorderColor,
		"BorderOpacity": borderOpacity(data.HideBorder),
		"TitleColor":    data.TitleColor,
		"Animations":    data.Animations,
		"A11yTitle":     data.A11yTitle,
		"A11yDesc":      data.A11yDesc,
		"Body":          data.Body,
	}

	var buf bytes.Buffer
	if err := tmpl.Execute(&buf, payload); err != nil {
		return "", err
	}
	return buf.String(), nil
}

func borderOpacity(hide bool) int {
	if hide {
		return 0
	}
	return 1
}

func renderStatItem(item statTextItem) string {
	tmpl := template.Must(template.New("stat-item").Parse(`
<g class="stagger" style="animation-delay: {{.Delay}}ms" transform="translate(25, 0)">
  {{if .ShowIcons}}
  <svg data-testid="icon" class="icon" viewBox="0 0 16 16" version="1.1" width="16" height="16">
    {{.Icon}}
  </svg>
  {{end}}
  <text class="stat {{if .Bold}}bold{{else}}not_bold{{end}}" {{if .ShowIcons}}x="25"{{end}} y="12.5">{{.Label}}:</text>
  <text class="stat {{if .Bold}}bold{{else}}not_bold{{end}}" x="{{.ValueX}}" y="12.5" data-testid="{{.ID}}">{{.Value}}{{if .UnitSymbol}} {{.UnitSymbol}}{{end}}</text>
</g>`))

	delay := (item.Index + 3) * 150
	valueX := 120 + item.ShiftValuePos
	if item.ShowIcons {
		valueX = 140 + item.ShiftValuePos
	}

	var buf bytes.Buffer
	_ = tmpl.Execute(&buf, map[string]any{
		"Delay":      delay,
		"ShowIcons":  item.ShowIcons,
		"Icon":       item.Icon,
		"Label":      item.Label,
		"Value":      item.Value,
		"UnitSymbol": item.UnitSymbol,
		"ID":         item.ID,
		"Bold":       item.Bold,
		"ValueX":     valueX,
	})

	return buf.String()
}

func getStyles(textColor string, iconColor string, ringColor string, showIcons bool, progress float64) string {
	display := "none"
	if showIcons {
		display = "block"
	}

	return fmt.Sprintf(`
    .stat { font: 600 14px 'Segoe UI', Ubuntu, "Helvetica Neue", Sans-Serif; fill: %s; }
    @supports(-moz-appearance: auto) { .stat { font-size:12px; } }
    .stagger { opacity: 0; animation: fadeInAnimation 0.3s ease-in-out forwards; }
    .rank-text { font: 800 24px 'Segoe UI', Ubuntu, Sans-Serif; fill: %s; animation: scaleInAnimation 0.3s ease-in-out forwards; }
    .rank-percentile-header { font-size: 14px; }
    .rank-percentile-text { font-size: 16px; }
    .not_bold { font-weight: 400 }
    .bold { font-weight: 700 }
    .icon { fill: %s; display: %s; }
    .rank-circle-rim { stroke: %s; fill: none; stroke-width: 6; opacity: 0.2; }
    .rank-circle { stroke: %s; stroke-dasharray: 250; fill: none; stroke-width: 6; stroke-linecap: round; opacity: 0.8; transform-origin: -10px 8px; transform: rotate(-90deg); animation: rankAnimation 1s forwards ease-in-out; }
    @keyframes rankAnimation { from { stroke-dashoffset: %f; } to { stroke-dashoffset: %f; } }
  `, textColor, textColor, iconColor, display, ringColor, ringColor, calculateCircleProgress(0), calculateCircleProgress(progress))
}

func calculateCircleProgress(value float64) float64 {
	radius := 40.0
	circumference := math.Pi * (radius * 2)
	if value < 0 {
		value = 0
	}
	if value > 100 {
		value = 100
	}
	return ((100 - value) / 100) * circumference
}

func calculateRankXTranslation(width int, minCardWidth int, hasStats bool, iconWidth int) int {
	if hasStats {
		minX := rankCardMinWidth + iconWidth - 70
		if width > rankCardDefaultWidth {
			xMaxExpansion := minX + (rankCardDefaultWidth-minCardWidth)/2
			return xMaxExpansion + (width - rankCardDefaultWidth)
		}
		return minX + (width-minCardWidth)/2
	}
	return width/2 + 20 - 10
}

func renderGradient(colors []string) string {
	if len(colors) <= 1 {
		return ""
	}
	angle := colors[0]
	stops := []string{}
	for idx, color := range colors[1:] {
		offset := 0
		if len(colors) > 2 {
			offset = (idx * 100) / (len(colors) - 2)
		}
		stops = append(stops, fmt.Sprintf("<stop offset=\"%d%%\" stop-color=\"#%s\" />", offset, color))
	}
	return fmt.Sprintf(`<defs><linearGradient id="gradient" gradientTransform="rotate(%s)" gradientUnits="userSpaceOnUse">%s</linearGradient></defs>`, angle, strings.Join(stops, ""))
}

func totalCommitsYearLabel(includeAllCommits bool, commitsYear *int, i18n *common.I18n) string {
	if includeAllCommits {
		return ""
	}
	if commitsYear != nil {
		return fmt.Sprintf(" (%d)", *commitsYear)
	}
	return " (" + i18n.T("wakatimecard.lastyear") + ")"
}

func contains(list []string, target string) bool {
	for _, item := range list {
		if item == target {
			return true
		}
	}
	return false
}

func orderStatKeys() []string {
	return []string{"stars", "commits", "prs", "prs_merged", "prs_merged_percentage", "reviews", "issues", "discussions_started", "discussions_answered", "contribs"}
}

func clampPrecision(value *int) int {
	if value == nil {
		return 2
	}
	precision := int(common.ClampValue(float64(*value), 0, 2))
	return precision
}

func applyNumberFormat(id string, raw string, numberFormat string, numberPrecision *int) string {
	if id == "prs_merged_percentage" {
		return raw
	}

	if numberFormat == "" {
		numberFormat = "short"
	}

	if strings.ToLower(numberFormat) == "long" {
		return raw
	}

	value := 0.0
	fmt.Sscanf(raw, "%f", &value)
	return common.KFormatter(value, numberPrecision)
}

func buildAccessibilityLabels(stats map[string]statTextItem, hide []string, includeAllCommits bool, commitsYear *int, i18n *common.I18n) string {
	items := []string{}
	for _, key := range orderStatKeys() {
		if contains(hide, key) {
			continue
		}
		stat, ok := stats[key]
		if !ok {
			continue
		}
		if key == "commits" {
			items = append(items, fmt.Sprintf("%s%s : %s", i18n.T("statcard.commits"), totalCommitsYearLabel(includeAllCommits, commitsYear, i18n), stat.Value))
			continue
		}
		items = append(items, fmt.Sprintf("%s: %s", stat.Label, stat.Value))
	}
	return strings.Join(items, ", ")
}
