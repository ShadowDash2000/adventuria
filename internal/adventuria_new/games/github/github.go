package github

import "context"

type repository interface {
	FetchLatestReleaseDownloadUrl(ctx context.Context, owner, repo string) (string, error)
}

type Github struct {
	repository repository
}

func NewGithub(repository repository) *Github {
	return &Github{repository: repository}
}

func (g *Github) FetchLatestReleaseDownloadUrl(ctx context.Context, owner, repo string) (string, error) {
	return g.repository.FetchLatestReleaseDownloadUrl(ctx, owner, repo)
}
