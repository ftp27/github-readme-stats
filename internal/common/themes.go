package common

type Theme struct {
	TitleColor  string
	IconColor   string
	TextColor   string
	BgColor     string
	BorderColor string
	RingColor   string
}

var Themes = map[string]Theme{
	"default": {
		TitleColor:  "2f80ed",
		IconColor:   "4c71f2",
		TextColor:   "434d58",
		BgColor:     "fffefe",
		BorderColor: "e4e2e2",
	},
}
