package http

import "golang.org/x/text/language"

var languageMatcher = language.NewMatcher([]language.Tag{
	language.Russian,
	language.English,
})

var languageKeys = [...]string{
	"ru",
	"en",
}

func resolveLanguage(header string) string {
	tags, _, err := language.ParseAcceptLanguage(header)
	if err != nil || len(tags) == 0 {
		return ""
	}

	_, index, confidence := languageMatcher.Match(tags...)
	if confidence == language.No {
		return ""
	}

	return languageKeys[index]
}
