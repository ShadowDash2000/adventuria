package steam

import (
	"adventuria/internal/adventuria/games/github"
	"context"
	"encoding/json"
	"net/http"

	steamstore "github.com/ShadowDash2000/steam-store-go"
)

type Parser struct{}

type Release struct {
	Assets []Asset `json:"assets"`
}

type Asset struct {
	Name               string `json:"name"`
	BrowserDownloadURL string `json:"browser_download_url"`
}

func NewParser() *Parser {
	return &Parser{}
}

func (p *Parser) FetchLatestRelease(ctx context.Context) ([]steamstore.SteamSpyAppDetailsResponse, error) {
	downloadUrl, err := github.FetchLatestReleaseDownloadUrl(ctx, "ShadowDash2000", "steam-spy-scraper")
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
		Data []steamstore.SteamSpyAppDetailsResponse `json:"data"`
	}
	if err = json.NewDecoder(fileRes.Body).Decode(&games); err != nil {
		return nil, err
	}

	return games.Data, nil
}
