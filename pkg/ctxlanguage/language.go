package ctxlanguage

import (
	"context"
	"errors"
)

type ctxKey struct{}

var defaultLang = "en"

func SetDefaultLang(lang string) error {
	if lang == "" {
		return errors.New("language: default language is empty")
	}

	defaultLang = lang

	return nil
}

func GetCtxLang(ctx context.Context) string {
	lang, ok := ctx.Value(ctxKey{}).(string)
	if ok && lang != "" {
		return lang
	}

	return defaultLang
}

func SetCtxLang(ctx context.Context, lang string) context.Context {
	return context.WithValue(ctx, ctxKey{}, lang)
}
