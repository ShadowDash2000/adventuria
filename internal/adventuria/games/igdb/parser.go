package igdb

import (
	"adventuria/internal/adventuria/games"
	"fmt"

	"github.com/bestnite/go-igdb"
	"github.com/bestnite/go-igdb/endpoint"
	pb "github.com/bestnite/go-igdb/proto"
	"github.com/pocketbase/pocketbase/tools/types"
	"google.golang.org/protobuf/proto"
)

type Parser struct {
	filter string
	client *igdb.Client
}

func NewParser(clientID, clientSecret, filter string) *Parser {
	p := &Parser{
		client: igdb.New(clientID, clientSecret),
		filter: filter,
	}

	return p
}

func (p *Parser) ParsePlatforms() (chan []games.Platform, error) {
	count, err := p.client.Platforms.Count()
	if err != nil {
		return nil, err
	}

	ch := make(chan []games.Platform)

	go func() {
		defer close(ch)

		offset := uint64(0)
		limit := uint64(500)

		for offset < count {
			platforms, err := p.client.Platforms.Paginated(0, limit)
			if err != nil {
				return
			}

			res := make([]games.Platform, len(platforms))
			for _, platform := range platforms {
				res = append(res, games.Platform{
					IdDb:     platform.Id,
					Name:     platform.Name,
					Checksum: platform.Checksum,
				})
			}

			offset += limit
			ch <- res
		}
	}()

	return ch, nil
}

func (p *Parser) ParseCompanies() (chan []games.Company, error) {
	count, err := p.client.Companies.Count()
	if err != nil {
		return nil, err
	}

	ch := make(chan []games.Company)

	go func() {
		defer close(ch)

		offset := uint64(0)
		limit := uint64(500)

		for offset < count {
			companies, err := p.client.Companies.Paginated(offset, limit)
			if err != nil {
				return
			}

			res := make([]games.Company, len(companies))
			for _, company := range companies {
				res = append(res, games.Company{
					IdDb:     company.Id,
					Name:     company.Name,
					Checksum: company.Checksum,
				})
			}

			offset += limit
			ch <- res
		}
	}()

	return ch, nil
}

func (p *Parser) ParseGames() (chan []games.Game, error) {
	resp, err := p.client.Request(
		"POST",
		fmt.Sprintf("https://api.igdb.com/v4/%s/count.pb", endpoint.EPGames),
		fmt.Sprintf("where %s;", p.filter),
	)
	if err != nil {
		return nil, err
	}

	var res pb.Count
	if err = proto.Unmarshal(resp.Body(), &res); err != nil {
		return nil, fmt.Errorf("failed to unmarshal: %w", err)
	}

	ch := make(chan []games.Game)

	go func() {
		defer close(ch)

		count := uint64(res.Count)
		offset := uint64(0)
		limit := uint64(500)

		for offset < count {
			gamesIgdb, err := p.client.Games.Paginated(offset, limit)
			if err != nil {
				return
			}

			res := make([]games.Game, len(gamesIgdb))
			for _, game := range gamesIgdb {
				releaseDate, err := types.ParseDateTime(game.FirstReleaseDate.AsTime())
				if err != nil {
					return
				}

				platforms := make([]uint64, len(game.Platforms))
				for _, platform := range game.Platforms {
					platforms = append(platforms, platform.Id)
				}

				res = append(res, games.Game{
					IdDb:        game.Id,
					Name:        game.Name,
					ReleaseDate: releaseDate,
					Platforms:   platforms,
					Checksum:    game.Checksum,
				})
			}

			offset += limit
			ch <- res
		}
	}()

	return ch, nil
}
