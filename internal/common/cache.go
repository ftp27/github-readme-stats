package common

import "github.com/spf13/viper"

const (
	minSeconds  = 60
	hourSeconds = 60 * minSeconds
	daySeconds  = 24 * hourSeconds
)

type CacheDurations struct {
	OneMinute   int
	FiveMinutes int
	TenMinutes  int
	Fifteen     int
	Thirty      int
	TwoHours    int
	FourHours   int
	SixHours    int
	EightHours  int
	TwelveHours int
	OneDay      int
	TwoDay      int
	SixDay      int
	TenDay      int
}

type CacheTTL struct {
	Default int
	Min     int
	Max     int
}

type CacheTTLGroup struct {
	StatsCard CacheTTL
	TopLangs  CacheTTL
	PinCard   CacheTTL
	GistCard  CacheTTL
	Wakatime  CacheTTL
	Error     int
	Durations CacheDurations
}

var CacheConfig = CacheTTLGroup{
	StatsCard: CacheTTL{Default: daySeconds, Min: 12 * hourSeconds, Max: 2 * daySeconds},
	TopLangs:  CacheTTL{Default: 6 * daySeconds, Min: 2 * daySeconds, Max: 10 * daySeconds},
	PinCard:   CacheTTL{Default: 10 * daySeconds, Min: daySeconds, Max: 10 * daySeconds},
	GistCard:  CacheTTL{Default: 2 * daySeconds, Min: daySeconds, Max: 10 * daySeconds},
	Wakatime:  CacheTTL{Default: daySeconds, Min: 12 * hourSeconds, Max: 2 * daySeconds},
	Error:     10 * minSeconds,
	Durations: CacheDurations{
		OneMinute:   minSeconds,
		FiveMinutes: 5 * minSeconds,
		TenMinutes:  10 * minSeconds,
		Fifteen:     15 * minSeconds,
		Thirty:      30 * minSeconds,
		TwoHours:    2 * hourSeconds,
		FourHours:   4 * hourSeconds,
		SixHours:    6 * hourSeconds,
		EightHours:  8 * hourSeconds,
		TwelveHours: 12 * hourSeconds,
		OneDay:      daySeconds,
		TwoDay:      2 * daySeconds,
		SixDay:      6 * daySeconds,
		TenDay:      10 * daySeconds,
	},
}

func ResolveCacheSeconds(requested int, def int, min int, max int) int {
	cacheSeconds := int(ClampValue(float64(requested), float64(min), float64(max)))
	if requested == 0 {
		cacheSeconds = def
	}

	if env := viper.GetString("CACHE_SECONDS"); env != "" {
		if parsed, ok := ParseInt(env); ok {
			cacheSeconds = parsed
		}
	}

	return cacheSeconds
}
