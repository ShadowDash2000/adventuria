package igdb

import (
	"adventuria/internal/adventuria"
	"adventuria/internal/adventuria/games"
	"context"
	"errors"
	"os"

	"github.com/pocketbase/dbx"
	"github.com/pocketbase/pocketbase/core"
)

type ParserController struct {
	parser *Parser
}

func New() (*ParserController, error) {
	twitchClientId, ok := os.LookupEnv("TWITCH_CLIENT_ID")
	if !ok {
		return nil, errors.New("igdb: TWITCH_CLIENT_ID not found")
	}
	twitchClientSecret, ok := os.LookupEnv("TWITCH_CLIENT_SECRET")
	if !ok {
		return nil, errors.New("igdb: TWITCH_CLIENT_SECRET not found")
	}
	igdbParseFilter, ok := os.LookupEnv("IGDB_PARSE_FILTER")
	if !ok {
		return nil, errors.New("igdb: IGDB_PARSE_FILTER not found")
	}

	p := &ParserController{
		parser: NewParser(twitchClientId, twitchClientSecret, igdbParseFilter),
	}

	return p, nil
}

func (p *ParserController) Parse(ctx context.Context, limit uint64) {
	if err := p.parsePlatforms(ctx, limit); err != nil {
		adventuria.PocketBase.Logger().Error("Failed to parse platforms", "error", err)
		return
	}
	if err := p.parseGenres(ctx, limit); err != nil {
		adventuria.PocketBase.Logger().Error("Failed to parse genres", "error", err)
		return
	}
	if err := p.parseGames(ctx, limit); err != nil {
		adventuria.PocketBase.Logger().Error("Failed to parse games", "error", err)
	}
}

func (p *ParserController) parseGames(ctx context.Context, limit uint64) error {
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	gamesCount, err := p.parser.GamesCount(ctx)
	if err != nil {
		return err
	}

	adventuria.PocketBase.Logger().Info("igdb.parseGames", "games_count", gamesCount)

	gamesParsedPrev := adventuria.GameSettings.IGDBGamesParsed()
	if gamesParsedPrev >= gamesCount {
		adventuria.PocketBase.Logger().Info("IGDB: Nothing to parse", "games_parsed_prev", gamesParsedPrev, "games_count", gamesCount)
		cancel()
		return nil
	}

	ch, err := p.parser.ParseGames(ctx, gamesCount, gamesParsedPrev, limit)
	if err != nil {
		return err
	}

	count := gamesParsedPrev
	for msg := range ch {
		if msg.Err != nil {
			return msg.Err
		}

		if err = p.saveCompaniesFromGames(ctx, msg.Games); err != nil {
			return err
		}

		records := make([]games.UpdatableRecord, len(msg.Games))
		for i, game := range msg.Games {
			record := core.NewRecord(adventuria.GameCollections.Get(adventuria.CollectionGames))

			gameRecord := games.NewGameFromRecord(record)
			gameRecord.SetIdDb(game.IdDb)
			gameRecord.SetName(game.Name)
			gameRecord.SetSlug(game.Slug)
			gameRecord.SetReleaseDate(game.ReleaseDate)
			gameRecord.SetSteamAppId(game.SteamAppId)
			gameRecord.SetSteamAppPrice(-1)
			gameRecord.SetCampaignTime(-1)
			gameRecord.SetCover(game.Cover)
			gameRecord.SetChecksum(game.Checksum)

			platformIds, err := p.collectionReferenceToIds(game.Platforms)
			if err != nil {
				return err
			}
			gameRecord.SetPlatforms(platformIds)

			developerIds, err := p.collectionReferenceToIds(game.Developers)
			if err != nil {
				return err
			}
			gameRecord.SetDevelopers(developerIds)

			publisherIds, err := p.collectionReferenceToIds(game.Publishers)
			if err != nil {
				return err
			}
			gameRecord.SetPublishers(publisherIds)

			genreIds, err := p.collectionReferenceToIds(game.Genres)
			if err != nil {
				return err
			}
			gameRecord.SetGenres(genreIds)

			records[i] = gameRecord
		}

		err = p.batchUpdate(records)
		if err != nil {
			return err
		}

		count += uint64(len(records))

		adventuria.GameSettings.SetIGDBGamesParsed(count)
		if err = adventuria.PocketBase.Save(adventuria.GameSettings.ProxyRecord()); err != nil {
			adventuria.PocketBase.Logger().Error("igdb.parseGames: failed to save game settings", "error", err)
		}
	}

	return nil
}

func (p *ParserController) saveCompaniesFromGames(ctx context.Context, gs []games.Game) error {
	uniq := make(map[uint64]struct{})
	for _, g := range gs {
		for _, id := range g.Developers.Ids {
			uniq[id] = struct{}{}
		}
		for _, id := range g.Publishers.Ids {
			uniq[id] = struct{}{}
		}
	}

	if len(uniq) > 0 {
		ids := make([]uint64, 0, len(uniq))
		for id := range uniq {
			ids = append(ids, id)
		}

		companies, err := p.parser.FetchCompaniesByIDs(ctx, ids)
		if err != nil {
			return err
		}

		compRecords := make([]games.UpdatableRecord, len(companies))
		for i, company := range companies {
			record := core.NewRecord(adventuria.GameCollections.Get(adventuria.CollectionCompanies))

			companyRecord := games.NewCompanyFromRecord(record)
			companyRecord.SetIdDb(company.IdDb)
			companyRecord.SetName(company.Name)
			companyRecord.SetChecksum(company.Checksum)

			compRecords[i] = companyRecord
		}

		if err = p.batchUpdate(compRecords); err != nil {
			return err
		}
	}

	return nil
}

func (p *ParserController) parsePlatforms(ctx context.Context, limit uint64) error {
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	ch, err := p.parser.ParsePlatformsAll(ctx, limit)
	if err != nil {
		return err
	}

	for msg := range ch {
		if msg.Err != nil {
			return msg.Err
		}

		records := make([]games.UpdatableRecord, len(msg.Platforms))
		for i, platform := range msg.Platforms {
			record := core.NewRecord(adventuria.GameCollections.Get(adventuria.CollectionPlatforms))

			platformRecord := games.NewPlatformFromRecord(record)
			platformRecord.SetIdDb(platform.IdDb)
			platformRecord.SetName(platform.Name)
			platformRecord.SetChecksum(platform.Checksum)

			records[i] = platformRecord
		}

		err = p.batchUpdate(records)
		if err != nil {
			return err
		}
	}

	return nil
}

func (p *ParserController) parseCompanies(ctx context.Context, limit uint64) error {
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	ch, err := p.parser.ParseCompaniesAll(ctx, limit)
	if err != nil {
		return err
	}

	for msg := range ch {
		if msg.Err != nil {
			return msg.Err
		}

		records := make([]games.UpdatableRecord, len(msg.Companies))
		for i, company := range msg.Companies {
			record := core.NewRecord(adventuria.GameCollections.Get(adventuria.CollectionCompanies))

			companyRecord := games.NewCompanyFromRecord(record)
			companyRecord.SetIdDb(company.IdDb)
			companyRecord.SetName(company.Name)
			companyRecord.SetChecksum(company.Checksum)

			records[i] = companyRecord
		}

		err = p.batchUpdate(records)
		if err != nil {
			return err
		}
	}

	return nil
}

func (p *ParserController) parseGenres(ctx context.Context, limit uint64) error {
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	ch, err := p.parser.ParseGenresAll(ctx, limit)
	if err != nil {
		return err
	}

	for msg := range ch {
		if msg.Err != nil {
			return msg.Err
		}

		records := make([]games.UpdatableRecord, len(msg.Genres))
		for i, genre := range msg.Genres {
			record := core.NewRecord(adventuria.GameCollections.Get(adventuria.CollectionGenres))

			genreRecord := games.NewGenreFromRecord(record)
			genreRecord.SetIdDb(genre.IdDb)
			genreRecord.SetName(genre.Name)
			genreRecord.SetChecksum(genre.Checksum)

			records[i] = genreRecord
		}

		err = p.batchUpdate(records)
		if err != nil {
			return err
		}
	}

	return nil
}

func (p *ParserController) batchUpdate(records []games.UpdatableRecord) error {
	checksums, err := p.obtainChecksums(records)
	if err != nil {
		return err
	}

	for _, record := range records {
		if extRecord, ok := checksums[record.IdDb()]; ok {
			if extRecord.GetString("checksum") == record.Checksum() {
				continue
			}

			data := record.ProxyRecord().FieldsData()
			baseFields := []string{"id", "created", "updated"}
			for _, field := range baseFields {
				delete(data, field)
			}

			extRecord.Load(data)

			err = adventuria.PocketBase.Save(extRecord)
			if err != nil {
				return err
			}
			continue
		}

		err = adventuria.PocketBase.Save(record.ProxyRecord())
		if err != nil {
			return err
		}
	}

	return nil
}

func (p *ParserController) obtainChecksums(updatables []games.UpdatableRecord) (map[uint64]*core.Record, error) {
	if len(updatables) == 0 {
		return nil, nil
	}

	idsDb := make([]any, len(updatables))
	for i, updatable := range updatables {
		idsDb[i] = int(updatable.IdDb())
	}

	var records []*core.Record
	err := adventuria.PocketBase.
		RecordQuery(updatables[0].ProxyRecord().Collection()).
		Where(
			dbx.In("id_db", idsDb...),
		).
		All(&records)
	if err != nil {
		return nil, err
	}

	checksums := make(map[uint64]*core.Record, len(records))
	for _, record := range records {
		checksums[uint64(record.GetInt("id_db"))] = record
	}

	return checksums, nil
}

func (p *ParserController) collectionReferenceToIds(reference games.CollectionReference) ([]string, error) {
	idsDb := make([]any, len(reference.Ids))
	for i, id := range reference.Ids {
		idsDb[i] = int(id)
	}

	records := make([]*core.Record, len(idsDb))
	err := adventuria.PocketBase.
		RecordQuery(adventuria.GameCollections.Get(reference.Collection)).
		Where(
			dbx.In("id_db", idsDb...),
		).
		All(&records)
	if err != nil {
		return nil, err
	}

	ids := make([]string, len(records))
	for i, record := range records {
		ids[i] = record.Id
	}

	return ids, nil
}
