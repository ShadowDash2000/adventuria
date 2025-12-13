package steam

import (
	"adventuria/internal/adventuria"
	"adventuria/internal/adventuria/games"
	"context"
	"database/sql"
	"errors"
	"hash/crc32"

	steamstore "github.com/ShadowDash2000/steam-store-go"
	"github.com/pocketbase/dbx"
	"github.com/pocketbase/pocketbase/core"
)

type ParserController struct {
	parser *Parser
}

func New() (*ParserController, error) {
	return &ParserController{
		parser: NewParser(),
	}, nil
}

func (p *ParserController) Parse(ctx context.Context) {
	if err := p.parsePricesAndTags(ctx); err != nil {
		adventuria.PocketBase.Logger().Error("Failed to parse prices and tags", "error", err)
		return
	}
}

func (p *ParserController) parsePricesAndTags(ctx context.Context) error {
	tags, err := p.getTags()
	if err != nil {
		return err
	}

	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	ok, err := p.isWhereGamesWithoutPrice(ctx)
	if err != nil || !ok {
		return err
	}

	ch := p.parser.ParseAllApps(ctx, adventuria.GameSettings.SteamSpyLastPage())

	for msg := range ch {
		if msg.Err != nil {
			if errors.Is(msg.Err, steamstore.ErrNoApiKey) {
				adventuria.GameSettings.SetSteamSpyLastPage(0)
				_ = adventuria.PocketBase.Save(adventuria.GameSettings.ProxyRecord())
				break
			}
			adventuria.GameSettings.SetSteamSpyLastPage(msg.Page)
			_ = adventuria.PocketBase.Save(adventuria.GameSettings.ProxyRecord())
			return msg.Err
		}

		adventuria.GameSettings.SetSteamSpyLastPage(msg.Page)
		_ = adventuria.PocketBase.Save(adventuria.GameSettings.ProxyRecord())

		gameRecords, err := p.getGamesWithoutPriceBySteamAppDetails(ctx, msg.Apps)
		if err != nil {
			if errors.Is(err, sql.ErrNoRows) {
				continue
			}
			return err
		}

		appIdsToTags := make(map[uint]steamstore.SteamSpyTags)
		for _, appDetails := range msg.Apps {
			gameRecord, ok := gameRecords[appDetails.AppId]
			if ok {
				gameRecord.SetSteamAppPrice(int(appDetails.Price))
			} else {
				// delete, 'cause no changes to update
				delete(gameRecords, appDetails.AppId)
				adventuria.PocketBase.Logger().Debug(
					"steam.parsePricesAndTags(): Game not found",
					"appId", appDetails.AppId,
					"app", appDetails,
				)
			}

			appIdsToTags[appDetails.AppId] = appDetails.Tags
			for tag := range appDetails.Tags {
				if _, ok = tags[tag]; ok {
					continue
				}

				record := core.NewRecord(adventuria.GameCollections.Get(adventuria.CollectionTags))
				tagRecord := games.NewTagFromRecord(record)
				tagRecord.SetIdDb(uint64(crc32.ChecksumIEEE([]byte(tag))))
				tagRecord.SetName(tag)
				tagRecord.SetChecksum(tag)

				tags[tag] = tagRecord
			}
		}

		err = p.updateUnknowTags(tags)
		if err != nil {
			return err
		}

		for _, gameRecord := range gameRecords {
			for tag := range appIdsToTags[uint(gameRecord.SteamAppId())] {
				tagRecord := tags[tag]
				gameRecord.SetTags(append(gameRecord.Tags(), tagRecord.ID()))
			}

			err = adventuria.PocketBase.Save(gameRecord.ProxyRecord())
			if err != nil {
				return err
			}
		}
	}

	return nil
}

func (p *ParserController) updateUnknowTags(tags map[string]games.TagRecord) error {
	for _, tag := range tags {
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

func (p *ParserController) isWhereGamesWithoutPrice(ctx context.Context) (bool, error) {
	count := -1
	err := adventuria.PocketBase.RecordQuery(adventuria.GameCollections.Get(adventuria.CollectionGames)).
		WithContext(ctx).
		Where(
			dbx.And(
				dbx.Not(dbx.HashExp{"steam_app_id": 0}),
				dbx.HashExp{"steam_app_price": -1},
			),
		).
		Select("COUNT(DISTINCT [[id]])").
		OrderBy( /* reset */ ).
		Row(&count)
	if err != nil {
		return false, err
	}
	adventuria.PocketBase.Logger().Debug("steam.isWhereGamesWithoutPrice(): count", "count", count)
	return count > 0, nil
}

func (p *ParserController) getGamesWithoutPriceBySteamAppDetails(ctx context.Context, apps map[string]steamstore.SteamSpyAppDetailsResponse) (map[uint]games.GameRecord, error) {
	appIds := make([]any, 0, len(apps))
	for _, app := range apps {
		appIds = append(appIds, app.AppId)
	}

	var records []*core.Record
	err := adventuria.PocketBase.RecordQuery(adventuria.GameCollections.Get(adventuria.CollectionGames)).
		WithContext(ctx).
		Where(dbx.NewExp("steam_app_price = -1")).
		AndWhere(dbx.In("steam_app_id", appIds...)).
		All(&records)
	if err != nil {
		return nil, err
	}

	if len(records) == 0 {
		return nil, sql.ErrNoRows
	}

	res := make(map[uint]games.GameRecord, len(records))
	for _, record := range records {
		res[uint(record.GetInt("steam_app_id"))] = games.NewGameFromRecord(record)
	}

	return res, nil
}
