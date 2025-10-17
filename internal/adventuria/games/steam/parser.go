package steam

import (
	"context"

	steamstore "github.com/ShadowDash2000/steam-store-go"
)

type Parser struct {
	client *steamstore.Client
}

func NewParser(apiKey string) *Parser {
	return &Parser{
		client: steamstore.New(
			steamstore.WithKey(apiKey),
		),
	}
}

type AppDetail struct {
	Price uint
	Tags  map[string]uint
}

func (p *Parser) ParseAppDetails(ctx context.Context, appId uint) (*AppDetail, error) {
	appDetail, err := p.client.GetSteamSpyAppDetails(ctx, appId)
	if err != nil {
		return nil, err
	}

	return &AppDetail{
		Price: uint(appDetail.Price),
		Tags:  appDetail.Tags,
	}, nil
}
