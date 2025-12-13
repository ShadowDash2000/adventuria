package igdb

import (
	"context"
	"os"
	"testing"

	"github.com/joho/godotenv"
)

func init() {
	_ = godotenv.Load("../../../../.env")
}

func envOrSkip(t *testing.T, key string) string {
	t.Helper()
	v := os.Getenv(key)
	if v == "" {
		t.Skipf("env %s is not set; skipping", key)
	}
	return v
}

func newParser(t *testing.T) (*Parser, error) {
	twitchClientId := envOrSkip(t, "TWITCH_CLIENT_ID")
	twitchClientSecret := envOrSkip(t, "TWITCH_CLIENT_SECRET")
	igdbParseFilter := envOrSkip(t, "IGDB_PARSE_FILTER")

	return NewParser(twitchClientId, twitchClientSecret, igdbParseFilter), nil
}

func Test_ParsePlatforms(t *testing.T) {
	p, err := newParser(t)
	if err != nil {
		t.Fatal(err)
	}

	ch, err := p.ParsePlatforms(context.Background(), 50, 50)
	if err != nil {
		t.Fatal(err)
	}

	for msg := range ch {
		if msg.Err != nil {
			t.Error(msg)
			return
		}

		for _, game := range msg.Platforms {
			t.Logf("%+v", game)
		}
	}
}

func Test_ParseCompanies(t *testing.T) {
	p, err := newParser(t)
	if err != nil {
		t.Fatal(err)
	}

	ch, err := p.ParseCompanies(context.Background(), 50, 50)
	if err != nil {
		t.Fatal(err)
	}

	for msg := range ch {
		if msg.Err != nil {
			t.Error(msg)
			return
		}

		for _, game := range msg.Companies {
			t.Logf("%+v", game)
		}
	}
}

func Test_ParseGames(t *testing.T) {
	p, err := newParser(t)
	if err != nil {
		t.Fatal(err)
	}

	ch, err := p.ParseGames(context.Background(), 10, 0, 10)
	if err != nil {
		t.Fatal(err)
	}

	for msg := range ch {
		if msg.Err != nil {
			t.Error(msg)
			return
		}

		for _, game := range msg.Games {
			t.Logf("%+v", game)
		}
	}
}
