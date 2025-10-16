package igdb

import (
	"adventuria/pkg/config"
	"context"
	"errors"
	"os"
	"testing"
)

func init() {
	config.LoadEnv("../../../../.env")
}

func newParser() (*Parser, error) {
	twitchClientId, ok := os.LookupEnv("TWITCH_CLIENT_ID")
	if !ok {
		return nil, errors.New("igdb: TWITCH_CLIENT_ID not found")
	}
	twitchClientSecret, ok := os.LookupEnv("TWITCH_CLIENT_SECRET")
	if !ok {
		return nil, errors.New("igdb: TWITCH_CLIENT_SECRET not found")
	}
	igdbParseFilter, ok := os.LookupEnv("IGDB_PARSE_FILTER")
	if !ok {
		return nil, errors.New("igdb: IGDB_PARSE_FILTER not found")
	}

	return NewParser(twitchClientId, twitchClientSecret, igdbParseFilter), nil
}

func Test_ParsePlatforms(t *testing.T) {
	p, err := newParser()
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
	p, err := newParser()
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
	p, err := newParser()
	if err != nil {
		t.Fatal(err)
	}

	ch, err := p.ParseGames(context.Background(), 10, 10)
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
