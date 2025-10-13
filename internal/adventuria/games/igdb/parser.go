package igdb

import (
	"adventuria/internal/adventuria"
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
			for i, platform := range platforms {
				res[i] = games.Platform{
					IdDb:     platform.GetId(),
					Name:     platform.GetName(),
					Checksum: platform.GetChecksum(),
				}
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
			for i, company := range companies {
				res[i] = games.Company{
					IdDb:     company.GetId(),
					Name:     company.GetName(),
					Checksum: company.GetChecksum(),
				}
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
			for i, game := range gamesIgdb {
				releaseDate, err := types.ParseDateTime(game.GetFirstReleaseDate().AsTime())
				if err != nil {
					return
				}

				platformIds := make([]uint64, len(game.GetPlatforms()))
				for i, platform := range game.GetPlatforms() {
					platformIds[i] = platform.Id
				}

				companyIds := make([]uint64, len(game.GetInvolvedCompanies()))
				for i, company := range game.GetInvolvedCompanies() {
					companyIds[i] = company.Id
				}

				res[i] = games.Game{
					IdDb:        game.GetId(),
					Name:        game.GetName(),
					ReleaseDate: releaseDate,
					Platforms: games.CollectionReference{
						Ids:        platformIds,
						Collection: adventuria.CollectionPlatforms,
					},
					Companies: games.CollectionReference{
						Ids:        companyIds,
						Collection: adventuria.CollectionCompanies,
					},
					Cover:    game.GetCover().GetUrl(),
					Checksum: game.GetChecksum(),
				}
			}

			offset += limit
			ch <- res
		}
	}()

	return ch, nil
}
