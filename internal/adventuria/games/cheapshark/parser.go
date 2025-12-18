package cheapshark

import (
	"adventuria/internal/adventuria/games/github"
	"context"
	"encoding/json"
	"net/http"
)

type Parser struct{}

func NewParser() *Parser {
	return &Parser{}
}

func (p *Parser) FetchLatestRelease(ctx context.Context) ([]CheapSharkResponse, error) {
	downloadUrl, err := github.FetchLatestReleaseDownloadUrl(ctx, "ShadowDash2000", "cheapshark-scraper")
	if err != nil {
		return nil, err
	}

	client := http.DefaultClient
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, downloadUrl, nil)
	if err != nil {
		return nil, err
	}
	fileRes, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer fileRes.Body.Close()

	var games struct {
		Data []CheapSharkResponse `json:"data"`
	}
	if err = json.NewDecoder(fileRes.Body).Decode(&games); err != nil {
		return nil, err
	}

	return games.Data, nil
}
