package common

import "strings"

func EncodeHTML(value string) string {
	replacer := strings.NewReplacer(
		"&", "&#38;",
		"<", "&#60;",
		">", "&#62;",
		"\u00a0", "&#160;",
	)
	return replacer.Replace(value)
}
