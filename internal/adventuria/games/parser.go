package games

import (
	"adventuria/pkg/collections"
	"fmt"
	"github.com/bestnite/go-igdb"
	"github.com/bestnite/go-igdb/endpoint"
	pb "github.com/bestnite/go-igdb/proto"
	"github.com/pocketbase/dbx"
	"github.com/pocketbase/pocketbase/core"
	"github.com/pocketbase/pocketbase/tools/types"
	"google.golang.org/protobuf/proto"
)

type Parser struct {
	app    core.App
	cols   *collections.Collections
	filter string
	client *igdb.Client
}

func NewParser(clientID, clientSecret, filter string, cols *collections.Collections, app core.App) *Parser {
	p := &Parser{
		client: igdb.New(clientID, clientSecret),
		filter: filter,
		cols:   cols,
		app:    app,
	}

	app.Cron().MustAdd("games_parser", "0 3 * * 0", p.Parse)

	return p
}

func (p *Parser) Parse() {
	/*	err := p.ParsePlatforms()
		if err != nil {
			p.app.Logger().Error(err.Error())
		}

		err = p.ParseCompanies()
		if err != nil {
			p.app.Logger().Error(err.Error())
		}*/

	err := p.ParseGames()
	if err != nil {
		p.app.Logger().Error(err.Error())
	}
}

func (p *Parser) ParsePlatforms() error {
	col, err := p.cols.Get(CollectionPlatforms)
	if err != nil {
		return err
	}

	count, err := p.client.Platforms.Count()
	if err != nil {
		return err
	}

	platforms, err := p.client.Platforms.Paginated(0, count)
	if err != nil {
		return err
	}

	for _, platform := range platforms {
		igdbPlatform := &Platform{}

		record, err := p.app.FindFirstRecordByFilter(
			CollectionPlatforms,
			"igdb_id = {:id}",
			dbx.Params{"id": platform.Id},
		)
		if err != nil {
			igdbPlatform.SetProxyRecord(core.NewRecord(col))
		} else {
			igdbPlatform.SetProxyRecord(record)
		}

		if igdbPlatform.Checksum() == platform.Checksum {
			continue
		}

		igdbPlatform.SetIGDBID(int(platform.Id))
		igdbPlatform.SetName(platform.Name)
		igdbPlatform.SetChecksum(platform.Checksum)
		igdbPlatform.SetData(platform)
		p.app.Save(igdbPlatform)
	}

	return nil
}

func (p *Parser) ParseCompanies() error {
	col, err := p.cols.Get(CollectionCompanies)
	if err != nil {
		return err
	}

	count, err := p.client.Companies.Count()
	if err != nil {
		return err
	}
	offset := uint64(0)
	limit := uint64(500)

	for offset < count {
		companies, err := p.client.Companies.Paginated(offset, limit)
		if err != nil {
			return err
		}

		for _, company := range companies {
			igdbCompany := &Company{}

			record, err := p.app.FindFirstRecordByFilter(
				CollectionCompanies,
				"igdb_id = {:id}",
				dbx.Params{"id": company.Id},
			)
			if err != nil {
				igdbCompany.SetProxyRecord(core.NewRecord(col))
			} else {
				igdbCompany.SetProxyRecord(record)
			}

			if igdbCompany.Checksum() == company.Checksum {
				continue
			}

			igdbCompany.SetIGDBID(int(company.Id))
			igdbCompany.SetName(company.Name)
			igdbCompany.SetChecksum(company.Checksum)
			igdbCompany.SetData(company)
			p.app.Save(igdbCompany)
		}

		offset += limit
	}

	return nil
}

func (p *Parser) ParseGames() error {
	col, err := p.cols.Get(CollectionGames)
	if err != nil {
		return err
	}

	resp, err := p.client.Request(
		"POST",
		fmt.Sprintf("https://api.igdb.com/v4/%s/count.pb", endpoint.EPGames),
		fmt.Sprintf("where %s;", p.filter),
	)
	if err != nil {
		return err
	}

	var res pb.Count
	if err = proto.Unmarshal(resp.Body(), &res); err != nil {
		return fmt.Errorf("failed to unmarshal: %w", err)
	}

	count := uint64(res.Count)
	offset := uint64(0)
	limit := uint64(500)

	for offset < count {
		games, err := p.client.Games.Paginated(offset, limit)
		if err != nil {
			return err
		}

		for _, game := range games {
			igdbGame := &Game{}

			record, err := p.app.FindFirstRecordByFilter(
				CollectionGames,
				"igdb_id = {:id}",
				dbx.Params{"id": game.Id},
			)
			if err != nil {
				igdbGame.SetProxyRecord(core.NewRecord(col))
			} else {
				igdbGame.SetProxyRecord(record)
			}

			if igdbGame.Checksum() == game.Checksum {
				continue
			}

			platformIds, err := p.platformsToIds(game.Platforms)
			if err != nil {
				return err
			}

			igdbGame.SetIGDBID(int(game.Id))
			igdbGame.SetName(game.Name)
			date, err := types.ParseDateTime(game.FirstReleaseDate.AsTime())
			if err == nil {
				igdbGame.SetFirstReleaseDate(date)
			}
			igdbGame.SetPlatforms(platformIds)
			igdbGame.SetChecksum(game.Checksum)
			igdbGame.SetData(game)
			err = p.app.Save(igdbGame)
			if err != nil {
				return err
			}
		}

		offset += limit
	}

	return nil
}

func (p *Parser) platformsToIds(platforms []*pb.Platform) ([]string, error) {
	col, err := p.cols.Get(CollectionPlatforms)
	if err != nil {
		return nil, err
	}

	platformIGDBIds := make([]any, len(platforms))
	for i, platform := range platforms {
		platformIGDBIds[i] = platform.Id
	}

	records := make([]*core.Record, len(platformIGDBIds))
	err = p.app.
		RecordQuery(col).
		Where(
			dbx.Or(
				dbx.HashExp{"igdb_id": platformIGDBIds},
			),
		).
		All(&records)
	if err != nil {
		return nil, err
	}

	platformIds := make([]string, len(records))
	for i, record := range records {
		platformIds[i] = record.Id
	}

	return platformIds, nil
}
