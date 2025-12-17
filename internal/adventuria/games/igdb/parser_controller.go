package igdb

import (
	"adventuria/internal/adventuria"
	"adventuria/internal/adventuria/games"
	"adventuria/internal/adventuria/games/hltb"
	"adventuria/internal/adventuria/games/steam"
	"context"
	"database/sql"
	"errors"
	"math"
	"os"
	"regexp"
	"sort"
	"strings"
	"unicode"

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

		if err = p.saveKeywordsFromGames(ctx, msg.Games); err != nil {
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
			gameRecord.SetCover(game.Cover)
			gameRecord.SetChecksum(game.Checksum)

			if game.SteamAppId > 0 {
				steamSpyRecord, err := p.findSteamSpyByAppId(ctx, uint(game.SteamAppId))
				if err == nil {
					gameRecord.SetSteamSpy(steamSpyRecord.Id)
					gameRecord.SetSteamAppPrice(steamSpyRecord.Price())
				}
			}

			hltbRecord, err := p.findHltbByGameName(ctx, game.Name)
			if err == nil {
				gameRecord.SetHltbId(hltbRecord.IdDb())
				gameRecord.SetCampaign(hltbRecord.Campaign())
			}

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

			tagIds, err := p.collectionReferenceToIds(game.Tags)
			if err != nil {
				return err
			}
			gameRecord.SetTags(tagIds)

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

func (p *ParserController) saveKeywordsFromGames(ctx context.Context, gs []games.Game) error {
	uniq := make(map[uint64]struct{})
	for _, g := range gs {
		for _, id := range g.Tags.Ids {
			uniq[id] = struct{}{}
		}
	}

	if len(uniq) == 0 {
		return nil
	}

	ids := make([]uint64, 0, len(uniq))
	for id := range uniq {
		ids = append(ids, id)
	}

	tags, err := p.parser.FetchKeywordsByIDs(ctx, ids)
	if err != nil {
		return err
	}

	tagRecords := make([]games.UpdatableRecord, len(tags))
	for i, tag := range tags {
		record := core.NewRecord(adventuria.GameCollections.Get(adventuria.CollectionTags))

		tagRecord := games.NewTagFromRecord(record)
		tagRecord.SetIdDb(tag.IdDb)
		tagRecord.SetName(tag.Name)
		tagRecord.SetChecksum(tag.Checksum)

		tagRecords[i] = tagRecord
	}

	return p.batchUpdate(tagRecords)
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

func (p *ParserController) findHltbByGameName(ctx context.Context, gameName string) (*hltb.HowLongToBeatRecord, error) {
	parts := strings.Fields(normalizeTitle(gameName))
	if len(parts) == 0 {
		return nil, errors.New("game name is empty")
	}

	var records []*core.Record
	err := adventuria.PocketBase.
		RecordQuery(adventuria.GameCollections.Get(adventuria.CollectionHowLongToBeat)).
		WithContext(ctx).
		Where(dbx.Like("name", parts...)).
		All(&records)
	if err != nil {
		return nil, err
	}

	if len(records) == 0 {
		return nil, sql.ErrNoRows
	}

	type match struct {
		record   *core.Record
		distance int
		diffLen  int
	}

	matches := make([]match, len(records))
	targetName := strings.ToLower(gameName)

	for i, r := range records {
		dbName := r.GetString("name")
		dist := levenshteinDistance(targetName, dbName)
		matches[i] = match{
			record:   r,
			distance: dist,
			diffLen:  int(math.Abs(float64(len(targetName) - len(dbName)))),
		}
	}

	sort.Slice(matches, func(i, j int) bool {
		if matches[i].distance == matches[j].distance {
			return matches[i].diffLen < matches[j].diffLen
		}
		return matches[i].distance < matches[j].distance
	})

	return hltb.NewHowLongToBeatRecordFromRecord(matches[0].record), nil
}

var (
	regParens = regexp.MustCompile(`\s*[(\[{].*?[)\]}]\s*`)
	regSpaces = regexp.MustCompile(`\s+`)
)

func normalizeTitle(s string) string {
	s = strings.TrimSpace(s)
	if s == "" {
		return s
	}

	s = regParens.ReplaceAllString(s, " ")

	s = strings.ToLower(s)

	s = strings.Map(func(r rune) rune {
		switch {
		case unicode.IsLetter(r) || unicode.IsDigit(r):
			return r
		case r == '\'' || r == '.' || r == '/' || r == '\\':
			return r
		default:
			return ' '
		}
	}, s)

	s = regSpaces.ReplaceAllString(s, " ")
	return strings.TrimSpace(s)
}

func levenshteinDistance(s1, s2 string) int {
	r1 := []rune(strings.ToLower(s1))
	r2 := []rune(strings.ToLower(s2))
	n, m := len(r1), len(r2)

	if n == 0 {
		return m
	}
	if m == 0 {
		return n
	}

	matrix := make([][]int, n+1)
	for i := range matrix {
		matrix[i] = make([]int, m+1)
	}

	for i := 0; i <= n; i++ {
		matrix[i][0] = i
	}
	for j := 0; j <= m; j++ {
		matrix[0][j] = j
	}

	for i := 1; i <= n; i++ {
		for j := 1; j <= m; j++ {
			cost := 1
			if r1[i-1] == r2[j-1] {
				cost = 0
			}
			matrix[i][j] = min(matrix[i-1][j]+1, min(matrix[i][j-1]+1, matrix[i-1][j-1]+cost))
		}
	}
	return matrix[n][m]
}

func (p *ParserController) findSteamSpyByAppId(ctx context.Context, appId uint) (*steam.SteamSpyRecord, error) {
	var record core.Record
	err := adventuria.PocketBase.
		RecordQuery(adventuria.GameCollections.Get(adventuria.CollectionSteamSpy)).
		WithContext(ctx).
		Where(dbx.HashExp{"id_db": appId}).
		One(&record)
	if err != nil {
		return nil, err
	}

	return steam.NewSteamSpyRecordFromRecord(&record), nil
}
