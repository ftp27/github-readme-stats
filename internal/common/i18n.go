package common

type I18n struct {
	Locale       string
	Translations map[string]map[string]string
}

const fallbackLocale = "en"

func NewI18n(locale string, translations map[string]map[string]string) *I18n {
	if locale == "" {
		locale = fallbackLocale
	}
	return &I18n{Locale: locale, Translations: translations}
}

func (i *I18n) T(key string) string {
	translations, ok := i.Translations[key]
	if !ok {
		return ""
	}
	if value, ok := translations[i.Locale]; ok {
		return value
	}
	return ""
}
