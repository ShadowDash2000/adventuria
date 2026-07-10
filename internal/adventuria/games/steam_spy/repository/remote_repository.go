package repository

import (
	"adventuria/internal/adventuria/games/steam_spy"
	"context"
	"encoding/json"
	"net/http"

	steamstore "github.com/ShadowDash2000/steam-store-go/apis/steam-spy"
)

type githubRepo interface {
	FetchLatestReleaseDownloadUrl(ctx context.Context, owner, repo string) (string, error)
}

type RemoteRepository struct {
	github githubRepo
}

func NewRemoteRepository(github githubRepo) *RemoteRepository {
	return &RemoteRepository{github: github}
}

func (r *RemoteRepository) FetchLatestRelease(ctx context.Context) ([]*steam_spy.SteamSpyResponse, error) {
	downloadUrl, err := r.github.FetchLatestReleaseDownloadUrl(ctx, "ShadowDash2000", "steam-spy-scraper")
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

	result := make([]*steam_spy.SteamSpyResponse, 0, len(games.Data))
	for _, g := range games.Data {
		result = append(result, &steam_spy.SteamSpyResponse{
			AppId: int(g.AppId),
			Name:  g.Name,
			Price: uint(g.Price),
		})
	}

	return result, nil
}
