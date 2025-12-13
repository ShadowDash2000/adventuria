package steam

import (
	"context"
	"time"

	steamstore "github.com/ShadowDash2000/steam-store-go"
)

type Parser struct {
	client *steamstore.Client
}

func NewParser() *Parser {
	return &Parser{
		client: steamstore.New(
			steamstore.WithRateLimit(60*time.Second),
			steamstore.WithBurst(1),
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

type ParseAllAppsMessage struct {
	Apps map[string]steamstore.SteamSpyAppDetailsResponse
	Page uint
	Err  error
}

func (p *Parser) ParseAllApps(ctx context.Context, startPage uint) <-chan ParseAllAppsMessage {
	ch := make(chan ParseAllAppsMessage)

	go func() {
		defer close(ch)

		page := startPage
		for {
			select {
			case <-ctx.Done():
				return
			default:
				res, err := p.client.GetSteamSpyAppsPaginated(ctx, page)
				if err != nil {
					ch <- ParseAllAppsMessage{Page: page, Err: err}
					return
				}

				ch <- ParseAllAppsMessage{Apps: res, Page: page}

				page++
			}
		}
	}()

	return ch
}
