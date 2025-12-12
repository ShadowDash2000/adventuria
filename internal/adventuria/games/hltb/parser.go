package hltb

import (
	"context"
	"errors"

	"github.com/ShadowDash2000/howlongtobeat"
)

type Parser struct {
	client *howlongtobeat.Client
}

func NewParser() (*Parser, error) {
	c, err := howlongtobeat.New()
	if err != nil {
		return nil, err
	}

	return &Parser{client: c}, nil
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
