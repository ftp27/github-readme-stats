package handlers

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/ftp27/github-readme-stats/internal/cards"
	"github.com/ftp27/github-readme-stats/internal/common"
	"github.com/ftp27/github-readme-stats/internal/fetchers"
	"github.com/gin-gonic/gin"
	"github.com/spf13/viper"
)

func Stats(c *gin.Context) {
	query := c.Request.URL.Query()
	username := query.Get("username")
	locale := query.Get("locale")

	c.Header("Content-Type", "image/svg+xml")

	access := common.GuardAccess(username, "username", common.AccessColors{
		TitleColor:  query.Get("title_color"),
		TextColor:   query.Get("text_color"),
		BgColor:     query.Get("bg_color"),
		BorderColor: query.Get("border_color"),
		Theme:       query.Get("theme"),
	})
	if !access.IsPassed {
		c.Data(http.StatusOK, "image/svg+xml", []byte(access.Result))
		return
	}

	if locale != "" && !common.IsLocaleAvailable(locale) {
		svg, _ := common.RenderError(
			"Something went wrong",
			"Language not found",
			common.RenderOptions{
				TitleColor:  query.Get("title_color"),
				TextColor:   query.Get("text_color"),
				BgColor:     query.Get("bg_color"),
				BorderColor: query.Get("border_color"),
				Theme:       query.Get("theme"),
			},
		)
		c.Data(http.StatusOK, "image/svg+xml", []byte(svg))
		return
	}

	showStats := common.ParseArray(query.Get("show"))
	includeAllCommits := boolValue(common.ParseBoolean(query.Get("include_all_commits")))
	includeMerged := contains(showStats, "prs_merged") || contains(showStats, "prs_merged_percentage")
	includeDiscussions := contains(showStats, "discussions_started")
	includeDiscussionAnswers := contains(showStats, "discussions_answered")

	var commitsYear *int
	if parsed, ok := common.ParseInt(query.Get("commits_year")); ok {
		commitsYear = &parsed
	}

	stats, err := fetchers.FetchStats(
		username,
		includeAllCommits,
		common.ParseArray(query.Get("exclude_repo")),
		includeMerged,
		includeDiscussions,
		includeDiscussionAnswers,
		commitsYear,
	)
	if err != nil {
		setErrorCacheHeaders(c)
		secondary := common.RetrieveSecondaryMessage(err)
		showRepoLink := true
		if _, ok := err.(*common.MissingParamError); ok {
			showRepoLink = false
		}
		svg, _ := common.RenderError(
			err.Error(),
			secondary,
			common.RenderOptions{
				TitleColor:   query.Get("title_color"),
				TextColor:    query.Get("text_color"),
				BgColor:      query.Get("bg_color"),
				BorderColor:  query.Get("border_color"),
				Theme:        query.Get("theme"),
				ShowRepoLink: &showRepoLink,
			},
		)
		c.Data(http.StatusOK, "image/svg+xml", []byte(svg))
		return
	}

	cacheSeconds := common.ResolveCacheSeconds(
		parseIntDefault(query.Get("cache_seconds")),
		common.CacheConfig.StatsCard.Default,
		common.CacheConfig.StatsCard.Min,
		common.CacheConfig.StatsCard.Max,
	)
	setCacheHeaders(c, cacheSeconds)

	cardWidth := parseIntDefault(query.Get("card_width"))
	borderRadius := parseFloatDefault(query.Get("border_radius"), 4.5)
	lineHeight := parseIntDefault(query.Get("line_height"))
	showIcons := boolValue(common.ParseBoolean(query.Get("show_icons")))
	hideTitle := boolValue(common.ParseBoolean(query.Get("hide_title")))
	hideBorder := boolValue(common.ParseBoolean(query.Get("hide_border")))
	hideRank := boolValue(common.ParseBoolean(query.Get("hide_rank")))
	textBold := true
	if parsed := common.ParseBoolean(query.Get("text_bold")); parsed != nil {
		textBold = *parsed
	}
	disableAnimations := boolValue(common.ParseBoolean(query.Get("disable_animations")))

	var numberPrecision *int
	if parsed, ok := common.ParseInt(query.Get("number_precision")); ok {
		numberPrecision = &parsed
	}

	svg, err := cards.RenderStatsCard(stats, cards.StatCardOptions{
		Hide:              common.ParseArray(query.Get("hide")),
		ShowIcons:         showIcons,
		HideTitle:         hideTitle,
		HideBorder:        hideBorder,
		CardWidth:         cardWidth,
		HideRank:          hideRank,
		IncludeAllCommits: includeAllCommits,
		CommitsYear:       commitsYear,
		LineHeight:        lineHeight,
		TitleColor:        query.Get("title_color"),
		RingColor:         query.Get("ring_color"),
		IconColor:         query.Get("icon_color"),
		TextColor:         query.Get("text_color"),
		TextBold:          textBold,
		BgColor:           query.Get("bg_color"),
		Theme:             query.Get("theme"),
		CustomTitle:       query.Get("custom_title"),
		BorderRadius:      borderRadius,
		BorderColor:       query.Get("border_color"),
		NumberFormat:      query.Get("number_format"),
		NumberPrecision:   numberPrecision,
		Locale:            strings.ToLower(locale),
		DisableAnimations: disableAnimations,
		RankIcon:          query.Get("rank_icon"),
		Show:              showStats,
	})
	if err != nil {
		setErrorCacheHeaders(c)
		secondary := common.RetrieveSecondaryMessage(err)
		svg, _ := common.RenderError(
			err.Error(),
			secondary,
			common.RenderOptions{
				TitleColor:  query.Get("title_color"),
				TextColor:   query.Get("text_color"),
				BgColor:     query.Get("bg_color"),
				BorderColor: query.Get("border_color"),
				Theme:       query.Get("theme"),
			},
		)
		c.Data(http.StatusOK, "image/svg+xml", []byte(svg))
		return
	}

	c.Data(http.StatusOK, "image/svg+xml", []byte(svg))
}

func Pin(c *gin.Context) {
	sendNotImplementedSVG(c, "pin")
}

func TopLangs(c *gin.Context) {
	sendNotImplementedSVG(c, "top-langs")
}

func Wakatime(c *gin.Context) {
	sendNotImplementedSVG(c, "wakatime")
}

func Gist(c *gin.Context) {
	sendNotImplementedSVG(c, "gist")
}

func sendNotImplementedSVG(c *gin.Context, name string) {
	c.Header("Content-Type", "image/svg+xml")
	c.Header("Cache-Control", "no-store")

	svg := fmt.Sprintf(
		"<svg xmlns=\"http://www.w3.org/2000/svg\" width=\"495\" height=\"195\" role=\"img\" aria-label=\"%s card\">"+
			"<rect width=\"100%%\" height=\"100%%\" fill=\"#ffffff\" stroke=\"#e4e2e2\"/>"+
			"<text x=\"20\" y=\"40\" font-family=\"Verdana,Arial,Helvetica,sans-serif\" font-size=\"14\" fill=\"#333333\">"+
			"Something went wrong"+
			"</text>"+
			"<text x=\"20\" y=\"64\" font-family=\"Verdana,Arial,Helvetica,sans-serif\" font-size=\"12\" fill=\"#666666\">"+
			"Go server not implemented for %s"+
			"</text>"+
			"</svg>",
		name,
		name,
	)

	c.Data(http.StatusOK, "image/svg+xml", []byte(svg))
}

func boolValue(value *bool) bool {
	if value == nil {
		return false
	}
	return *value
}

func parseIntDefault(value string) int {
	parsed, ok := common.ParseInt(value)
	if !ok {
		return 0
	}
	return parsed
}

func parseFloatDefault(value string, def float64) float64 {
	parsed, ok := common.ParseFloat(value)
	if !ok {
		return def
	}
	return parsed
}

func contains(list []string, target string) bool {
	for _, item := range list {
		if item == target {
			return true
		}
	}
	return false
}

func setCacheHeaders(c *gin.Context, cacheSeconds int) {
	if cacheSeconds < 1 || strings.ToLower(viper.GetString("NODE_ENV")) == "development" {
		disableCaching(c)
		return
	}

	c.Header("Cache-Control", fmt.Sprintf("max-age=%d, s-maxage=%d, stale-while-revalidate=%d", cacheSeconds, cacheSeconds, common.CacheConfig.Durations.OneDay))
}

func setErrorCacheHeaders(c *gin.Context) {
	cacheSeconds, ok := common.ParseInt(strings.TrimSpace(viper.GetString("CACHE_SECONDS")))
	if ok && cacheSeconds < 1 {
		disableCaching(c)
		return
	}

	if strings.ToLower(viper.GetString("NODE_ENV")) == "development" {
		disableCaching(c)
		return
	}

	c.Header("Cache-Control", fmt.Sprintf("max-age=%d, s-maxage=%d, stale-while-revalidate=%d", common.CacheConfig.Error, common.CacheConfig.Error, common.CacheConfig.Durations.OneDay))
}

func disableCaching(c *gin.Context) {
	c.Header("Cache-Control", "no-cache, no-store, must-revalidate, max-age=0, s-maxage=0")
	c.Header("Pragma", "no-cache")
	c.Header("Expires", "0")
}
