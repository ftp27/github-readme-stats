package common

import "strings"

type AccessGuardResult struct {
	IsPassed bool
	Result   string
}

type AccessColors struct {
	TitleColor  string
	TextColor   string
	BgColor     string
	BorderColor string
	Theme       string
}

func GuardAccess(id string, guardType string, colors AccessColors) AccessGuardResult {
	whitelist := Whitelist()
	gistWhitelist := GistWhitelist()

	var currentWhitelist []string
	notWhitelistedMsg := "This username is not whitelisted"
	if guardType == "gist" {
		currentWhitelist = gistWhitelist
		notWhitelistedMsg = "This gist ID is not whitelisted"
	} else {
		currentWhitelist = whitelist
	}

	if len(currentWhitelist) > 0 && !contains(currentWhitelist, id) {
		showRepo := false
		svg, _ := RenderError(
			notWhitelistedMsg,
			"Please deploy your own instance",
			RenderOptions{
				TitleColor:   colors.TitleColor,
				TextColor:    colors.TextColor,
				BgColor:      colors.BgColor,
				BorderColor:  colors.BorderColor,
				Theme:        colors.Theme,
				ShowRepoLink: &showRepo,
			},
		)
		return AccessGuardResult{IsPassed: false, Result: svg}
	}

	if guardType == "username" && len(currentWhitelist) == 0 && contains(Blacklist, id) {
		showRepo := false
		svg, _ := RenderError(
			"This username is blacklisted",
			"Please deploy your own instance",
			RenderOptions{
				TitleColor:   colors.TitleColor,
				TextColor:    colors.TextColor,
				BgColor:      colors.BgColor,
				BorderColor:  colors.BorderColor,
				Theme:        colors.Theme,
				ShowRepoLink: &showRepo,
			},
		)
		return AccessGuardResult{IsPassed: false, Result: svg}
	}

	return AccessGuardResult{IsPassed: true}
}

func contains(list []string, target string) bool {
	for _, item := range list {
		if strings.EqualFold(item, target) {
			return true
		}
	}
	return false
}
