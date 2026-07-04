package repository

import (
	"adventuria/internal/adventuria_new/games/how_long_to_beat"
	"context"
	"encoding/json"
	"net/http"
	"slices"

	"github.com/ShadowDash2000/howlongtobeat"
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

func (r *RemoteRepository) FetchLatestRelease(ctx context.Context) ([]*how_long_to_beat.HowLongToBeatResponse, error) {
	downloadUrl, err := r.github.FetchLatestReleaseDownloadUrl(ctx, "ShadowDash2000", "hltb-scraper")
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
		Data []howlongtobeat.SearchGameData `json:"data"`
	}
	if err = json.NewDecoder(fileRes.Body).Decode(&games); err != nil {
		return nil, err
	}

	result := make([]*how_long_to_beat.HowLongToBeatResponse, 0, len(games.Data))
	for _, g := range games.Data {
		result = append(result, &how_long_to_beat.HowLongToBeatResponse{
			ID:           g.GameID,
			Name:         g.GameName,
			ReleaseWorld: g.ReleaseWorld,
			CompMain:     g.CompMain,
		})
	}

	result = slices.DeleteFunc(result, func(v *how_long_to_beat.HowLongToBeatResponse) bool {
		return v == nil
	})

	return result, nil
}
