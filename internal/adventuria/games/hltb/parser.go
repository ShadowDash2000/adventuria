package hltb

import (
	"context"
	"errors"

	"github.com/ShadowDash2000/hltb-crashdummy-go"
	"github.com/ShadowDash2000/howlongtobeat"
)

type Parser struct {
	client       *howlongtobeat.Client
	cachedClient *hltb.Client
}

func NewParser() (*Parser, error) {
	c, err := howlongtobeat.New()
	if err != nil {
		return nil, err
	}

	return &Parser{client: c, cachedClient: hltb.New()}, nil
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
