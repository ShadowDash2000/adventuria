package steam

import (
	"adventuria/internal/adventuria"
	"adventuria/internal/adventuria/games"
	"context"
	"errors"
	"hash/crc32"
	"os"

	"github.com/pocketbase/pocketbase/core"
	"golang.org/x/exp/maps"
)

type ParserController struct {
	parser *Parser

	tags map[string]games.TagRecord
}

func New() (*ParserController, error) {
	steamApiKey, ok := os.LookupEnv("STEAM_API_KEY")
	if !ok {
		return nil, errors.New("steam: STEAM_API_KEY not found")
	}

	return &ParserController{
		parser: NewParser(steamApiKey),
	}, nil
}

func (p *ParserController) Parse() {
	ctx := context.Background()

	if err := p.parsePricesAndTags(ctx, 100); err != nil {
		adventuria.PocketBase.Logger().Error("Failed to parse prices and tags", "error", err)
		return
	}
}

func (p *ParserController) parsePricesAndTags(ctx context.Context, limit int) error {
	tags, err := p.getTags()
	if err != nil {
		return err
	}
	p.tags = tags

	for {
		gameRecords, err := p.getSteamAppsWithoutPrice(limit)
		if err != nil {
			return err
		}

		if len(gameRecords) == 0 {
			break
		}

		tagsTmp := make(map[string]map[string]uint)
		for _, game := range gameRecords {
			appDetails, err := p.parser.ParseAppDetails(ctx, uint(game.SteamAppId()))
			if err != nil {
				return err
			}

			game.SetSteamAppPrice(int(appDetails.Price))

			tagsTmp[game.ID()] = appDetails.Tags
			for tag, _ := range appDetails.Tags {
				if _, ok := p.tags[tag]; ok {
					continue
				}

				record := core.NewRecord(adventuria.GameCollections.Get(adventuria.CollectionTags))
				tagRecord := games.NewTagFromRecord(record)
				tagRecord.SetIdDb(uint64(crc32.ChecksumIEEE([]byte(tag))))
				tagRecord.SetName(tag)
				tagRecord.SetChecksum(tag)

				p.tags[tag] = tagRecord
			}
		}

		err = p.updateUnknowTags()
		if err != nil {
			return err
		}

		for _, game := range gameRecords {
			for tag, _ := range tagsTmp[game.ID()] {
				tagRecord := p.tags[tag]
				game.SetTags(append(game.Tags(), tagRecord.ID()))
			}

			err = adventuria.PocketBase.Save(game.ProxyRecord())
			if err != nil {
				return err
			}
		}
	}

	maps.Clear(p.tags)

	return nil
}

func (p *ParserController) updateUnknowTags() error {
	for _, tag := range p.tags {
		if tag.ID() != "" {
			continue
		}

		err := adventuria.PocketBase.Save(tag.ProxyRecord())
		if err != nil {
			return err
		}
	}

	return nil
}

func (p *ParserController) getTags() (map[string]games.TagRecord, error) {
	records, err := adventuria.PocketBase.FindAllRecords(adventuria.GameCollections.Get(adventuria.CollectionTags))
	if err != nil {
		return nil, err
	}

	res := make(map[string]games.TagRecord, len(records))
	for _, record := range records {
		tagRecord := games.NewTagFromRecord(record)
		res[tagRecord.Name()] = tagRecord
	}

	return res, nil
}

func (p *ParserController) getSteamAppsWithoutPrice(limit int) ([]games.GameRecord, error) {
	records, err := adventuria.PocketBase.FindRecordsByFilter(
		adventuria.GameCollections.Get(adventuria.CollectionGames),
		"platforms.id_db ?= 6 && steam_app_id != 0 && steam_app_price = -1",
		"",
		limit,
		0,
		nil,
	)
	if err != nil {
		return nil, err
	}

	res := make([]games.GameRecord, len(records))
	for i, record := range records {
		res[i] = games.NewGameFromRecord(record)
	}

	return res, nil
}
