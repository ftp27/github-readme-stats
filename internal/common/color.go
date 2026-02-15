package common

import (
	"regexp"
	"strings"
)

type CardColors struct {
	TitleColor  string
	IconColor   string
	TextColor   string
	BgColor     string
	BgGradient  []string
	BorderColor string
	RingColor   string
}

var hexColorPattern = regexp.MustCompile(`^([A-Fa-f0-9]{8}|[A-Fa-f0-9]{6}|[A-Fa-f0-9]{3}|[A-Fa-f0-9]{4})$`)

func IsValidHexColor(value string) bool {
	return hexColorPattern.MatchString(value)
}

func isValidGradient(colors []string) bool {
	if len(colors) <= 2 {
		return false
	}
	for _, color := range colors[1:] {
		if !IsValidHexColor(color) {
			return false
		}
	}
	return true
}

func fallbackColor(color string, fallback string) (string, []string) {
	var gradient []string
	colors := []string{}
	if color != "" {
		colors = strings.Split(color, ",")
	}

	if len(colors) > 1 && isValidGradient(colors) {
		gradient = colors
	}

	if gradient != nil {
		return "", gradient
	}
	if IsValidHexColor(color) {
		return "#" + color, nil
	}

	return fallback, nil
}

func GetCardColors(args ColorsArgs) (CardColors, error) {
	theme := Themes["default"]
	if args.Theme != "" {
		if selected, ok := Themes[args.Theme]; ok {
			theme = selected
		}
	}

	borderColor := theme.BorderColor
	if borderColor == "" {
		borderColor = Themes["default"].BorderColor
	}

	titleColor, titleGradient := fallbackColor(args.TitleColor, "#"+Themes["default"].TitleColor)
	if titleGradient != nil {
		return CardColors{}, ErrColor
	}

	ringColor, ringGradient := fallbackColor(args.RingColor, titleColor)
	if ringGradient != nil {
		return CardColors{}, ErrColor
	}

	iconColor, iconGradient := fallbackColor(args.IconColor, "#"+Themes["default"].IconColor)
	if iconGradient != nil {
		return CardColors{}, ErrColor
	}

	textColor, textGradient := fallbackColor(args.TextColor, "#"+Themes["default"].TextColor)
	if textGradient != nil {
		return CardColors{}, ErrColor
	}

	bgColor, bgGradient := fallbackColor(args.BgColor, "#"+Themes["default"].BgColor)
	if bgGradient == nil && bgColor == "" {
		bgColor = "#" + Themes["default"].BgColor
	}

	borderColorFallback, borderGradient := fallbackColor(args.BorderColor, "#"+borderColor)
	if borderGradient != nil {
		return CardColors{}, ErrColor
	}

	return CardColors{
		TitleColor:  titleColor,
		IconColor:   iconColor,
		TextColor:   textColor,
		BgColor:     bgColor,
		BgGradient:  bgGradient,
		BorderColor: borderColorFallback,
		RingColor:   ringColor,
	}, nil
}

type ColorsArgs struct {
	TitleColor  string
	TextColor   string
	IconColor   string
	BgColor     string
	BorderColor string
	RingColor   string
	Theme       string
}

var ErrColor = &ColorError{Message: "Unexpected behavior, all colors except background should be string."}

type ColorError struct {
	Message string
}

func (e *ColorError) Error() string {
	return e.Message
}
