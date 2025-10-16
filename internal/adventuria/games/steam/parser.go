package steam

import (
	"adventuria/internal/adventuria/games"
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

func (p *Parser) ParsePrices(ctx context.Context, games []games.GameRecord) error {
	for _, game := range games {
		appDetail, err := p.client.GetSteamSpyAppDetails(ctx, uint(game.SteamAppId()))
		if err != nil {
			return err
		}

		game.SetSteamAppPrice(int(appDetail.Price))
	}

	return nil
}
