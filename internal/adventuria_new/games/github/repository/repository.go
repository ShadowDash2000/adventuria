package repository

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
)

type Repository struct{}

func NewRepository() *Repository {
	return &Repository{}
}

type release struct {
	Assets []asset `json:"assets"`
}

type asset struct {
	Name               string `json:"name"`
	BrowserDownloadURL string `json:"browser_download_url"`
}

func (r *Repository) FetchLatestReleaseDownloadUrl(ctx context.Context, owner, repo string) (string, error) {
	client := http.DefaultClient

	req, err := http.NewRequestWithContext(
		ctx,
		http.MethodGet,
		fmt.Sprintf("https://api.github.com/repos/%s/%s/releases/latest", owner, repo),
		nil,
	)
	if err != nil {
		return "", err
	}
	res, err := client.Do(req)
	if err != nil {
		return "", err
	}
	defer res.Body.Close()

	var release release
	if err = json.NewDecoder(res.Body).Decode(&release); err != nil {
		panic(err)
	}

	if release.Assets == nil {
		return "", errors.New("no assets found")
	}

	return release.Assets[0].BrowserDownloadURL, nil
}
