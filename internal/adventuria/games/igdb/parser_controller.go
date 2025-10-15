package igdb

import (
	"adventuria/internal/adventuria"
	"adventuria/internal/adventuria/games"
	"errors"
	"os"

	"github.com/pocketbase/dbx"
	"github.com/pocketbase/pocketbase/core"
)

type ParserController struct {
	parser games.Parser
}

func NewParserController() (*ParserController, error) {
	twitchClientId, ok := os.LookupEnv("TWITCH_CLIENT_ID")
	if !ok {
		return nil, errors.New("IGDB: TWITCH_CLIENT_ID not found")
	}
	twitchClientSecret, ok := os.LookupEnv("TWITCH_CLIENT_SECRET")
	if !ok {
		return nil, errors.New("IGDB: TWITCH_CLIENT_SECRET not found")
	}
	igdbParseFilter, ok := os.LookupEnv("IGDB_PARSE_FILTER")
	if !ok {
		return nil, errors.New("IGDB: IGDB_PARSE_FILTER not found")
	}

	p := &ParserController{
		parser: NewParser(twitchClientId, twitchClientSecret, igdbParseFilter),
	}

	return p, nil
}

func (p *ParserController) Parse() {
	if err := p.parseCompanies(); err != nil {
		adventuria.PocketBase.Logger().Error("Failed to parse companies", "error", err)
		return
	}
	if err := p.parsePlatforms(); err != nil {
		adventuria.PocketBase.Logger().Error("Failed to parse platforms", "error", err)
		return
	}
	if err := p.parseGames(); err != nil {
		adventuria.PocketBase.Logger().Error("Failed to parse games", "error", err)
	}
}

func (p *ParserController) parseGames() error {
	ch, err := p.parser.ParseGames()
	if err != nil {
		return err
	}

	collection, err := adventuria.GameCollections.Get(adventuria.CollectionGames)
	if err != nil {
		return err
	}

	for gamesIGDB := range ch {
		records := make([]games.UpdatableRecord, len(gamesIGDB))
		for i, game := range gamesIGDB {
			record := core.NewRecord(collection)

			gameRecord := games.NewGameFromRecord(record)
			gameRecord.SetIdDb(game.IdDb)
			gameRecord.SetName(game.Name)
			gameRecord.SetReleaseDate(game.ReleaseDate)
			gameRecord.SetSteamAppId(game.SteamAppId)
			gameRecord.SetChecksum(game.Checksum)

			platformIds, err := p.collectionReferenceToIds(game.Platforms)
			if err != nil {
				return err
			}
			gameRecord.SetPlatforms(platformIds)

			companyIds, err := p.collectionReferenceToIds(game.Companies)
			if err != nil {
				return err
			}
			gameRecord.SetCompanies(companyIds)

			records[i] = gameRecord
		}

		err = p.batchUpdate(records)
		if err != nil {
			return err
		}
	}

	return nil
}

func (p *ParserController) parsePlatforms() error {
	ch, err := p.parser.ParsePlatforms()
	if err != nil {
		return err
	}

	collection, err := adventuria.GameCollections.Get(adventuria.CollectionPlatforms)
	if err != nil {
		return err
	}

	for platforms := range ch {
		records := make([]games.UpdatableRecord, len(platforms))
		for i, platform := range platforms {
			record := core.NewRecord(collection)

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

func (p *ParserController) parseCompanies() error {
	ch, err := p.parser.ParseCompanies()
	if err != nil {
		return err
	}

	collection, err := adventuria.GameCollections.Get(adventuria.CollectionCompanies)
	if err != nil {
		return err
	}

	for companies := range ch {
		records := make([]games.UpdatableRecord, len(companies))
		for i, company := range companies {
			record := core.NewRecord(collection)

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

func (p *ParserController) batchUpdate(records []games.UpdatableRecord) error {
	checksums, err := p.obtainChecksums(records)
	if err != nil {
		return err
	}

	for _, record := range records {
		if checksum, ok := checksums[int(record.IdDb())]; ok {
			if checksum == record.Checksum() {
				continue
			}
		}

		err = adventuria.PocketBase.Save(record.ProxyRecord())
		if err != nil {
			return err
		}
	}

	return nil
}

func (p *ParserController) obtainChecksums(updatables []games.UpdatableRecord) (map[int]string, error) {
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

	checksums := make(map[int]string, len(records))
	for _, record := range records {
		checksums[record.GetInt("id_db")] = record.GetString("checksum")
	}

	return checksums, nil
}

func (p *ParserController) collectionReferenceToIds(reference games.CollectionReference) ([]string, error) {
	col, err := adventuria.GameCollections.Get(reference.Collection)
	if err != nil {
		return nil, err
	}

	idsDb := make([]any, len(reference.Ids))
	for i, id := range reference.Ids {
		idsDb[i] = int(id)
	}

	records := make([]*core.Record, len(idsDb))
	err = adventuria.PocketBase.
		RecordQuery(col).
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
