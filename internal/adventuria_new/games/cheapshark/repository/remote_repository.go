package repository

import (
	"adventuria/internal/adventuria/games/github"
	"adventuria/internal/adventuria_new/games/cheapshark"
	"context"
	"encoding/json"
	"net/http"
	"slices"
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

func (r *RemoteRepository) FetchLatestRelease(ctx context.Context) ([]*cheapshark.CheapSharkResponse, error) {
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
		Data []*cheapshark.CheapSharkResponse `json:"data"`
	}
	if err = json.NewDecoder(fileRes.Body).Decode(&games); err != nil {
		return nil, err
	}

	games.Data = slices.DeleteFunc(games.Data, func(v *cheapshark.CheapSharkResponse) bool {
		return v == nil
	})

	return games.Data, nil
}
