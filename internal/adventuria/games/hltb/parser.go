package hltb

import (
	"context"
	"errors"
	"time"

	"github.com/ShadowDash2000/hltb-crashdummy-go"
	"github.com/ShadowDash2000/howlongtobeat"
)

type Parser struct {
	client       *howlongtobeat.Client
	cachedClient *hltb.Client
}

func NewParser(r time.Duration, b int) (*Parser, error) {
	c, err := howlongtobeat.New(
		howlongtobeat.WithRateLimit(r, b),
	)
	if err != nil {
		return nil, err
	}

	return &Parser{
		client: c,
		cachedClient: hltb.New(
			hltb.WithRateLimit(r, b),
		),
	}, nil
}

type WalkthroughTime struct {
	GameID   int
	Campaign float64
}

var ErrGameNotFound = errors.New("hltb: game not found")

func (p *Parser) ParseTime(ctx context.Context, search string) (*WalkthroughTime, error) {
	res, err := p.client.SearchSimple(ctx, search, howlongtobeat.SearchModifierHideDLC)
	if err != nil {
		return nil, err
	}

	if len(res) == 0 {
		return nil, ErrGameNotFound
	}

	return &WalkthroughTime{
		GameID:   res[0].GameID,
		Campaign: res[0].CompMain,
	}, nil
}

func (p *Parser) ParseBySteamAppId(ctx context.Context, appId uint64) (*WalkthroughTime, error) {
	res, err := p.cachedClient.GetBySteamAppId(ctx, appId)
	if err != nil {
		if errors.Is(err, hltb.ErrNotFound) {
			return nil, ErrGameNotFound
		}

		return nil, err
	}

	return &WalkthroughTime{
		GameID:   int(res.HltbId),
		Campaign: res.MainStory,
	}, nil
}

func (p *Parser) RefreshToken(ctx context.Context, search bool) error {
	return p.client.RefreshToken(ctx, search)
}
