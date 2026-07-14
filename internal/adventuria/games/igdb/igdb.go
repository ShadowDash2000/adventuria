package igdb

import (
	"adventuria/internal/adventuria/errs"
	"adventuria/internal/adventuria/model"
	"adventuria/internal/adventuria/schema"
	"context"
	"errors"
	"slices"
)

type repository interface {
	TableReferenceToID(ctx context.Context, reference TableReferenceSingle) (string, error)
	TableReferenceToIDs(ctx context.Context, reference TableReference) ([]string, error)
}

type remoteRepository interface {
	ParseGames(ctx context.Context, count, offset, limit uint64) (<-chan ParseGamesMessage, error)
	GamesCount(ctx context.Context) (uint64, error)
	GetCompaniesByIDs(ctx context.Context, ids []uint64) ([]*Company, error)
	GetTagsByIDs(ctx context.Context, ids []uint64) ([]*Tag, error)
	GetThemesByIDs(ctx context.Context, ids []uint64) ([]*Theme, error)
	ParsePlatformsAll(ctx context.Context, limit uint64) (<-chan ParsePlatformsMessage, error)
	ParseGenresAll(ctx context.Context, limit uint64) (<-chan ParseGenresMessage, error)
	ParseGameTypesAll(ctx context.Context, limit uint64) (<-chan ParseGameTypesMessage, error)
}

type activities interface {
	GetOrCreate(ctx context.Context, data model.ActivityCreate) (*model.Activity, error)
	GetChecksumsByIDs(ctx context.Context, ids []string) (map[string]string, error)
	Save(ctx context.Context, activity *model.Activity) (*model.Activity, error)
}

type platforms interface {
	GetOrCreate(ctx context.Context, data model.PlatformCreate) (*model.Platform, error)
	GetChecksumsByIDs(ctx context.Context, ids []string) (map[string]string, error)
	Save(ctx context.Context, platform *model.Platform) (*model.Platform, error)
}

type companies interface {
	GetOrCreate(ctx context.Context, data model.CompanyCreate) (*model.Company, error)
	GetChecksumsByIDs(ctx context.Context, ids []string) (map[string]string, error)
	Save(ctx context.Context, company *model.Company) (*model.Company, error)
}

type tags interface {
	GetOrCreate(ctx context.Context, data model.TagCreate) (*model.Tag, error)
	GetChecksumsByIDs(ctx context.Context, ids []string) (map[string]string, error)
	Save(ctx context.Context, tag *model.Tag) (*model.Tag, error)
}

type themes interface {
	GetOrCreate(ctx context.Context, data model.ThemeCreate) (*model.Theme, error)
	GetChecksumsByIDs(ctx context.Context, ids []string) (map[string]string, error)
	Save(ctx context.Context, theme *model.Theme) (*model.Theme, error)
}

type genres interface {
	GetOrCreate(ctx context.Context, data model.GenreCreate) (*model.Genre, error)
	GetChecksumsByIDs(ctx context.Context, ids []string) (map[string]string, error)
	Save(ctx context.Context, genre *model.Genre) (*model.Genre, error)
}

type gameTypes interface {
	GetOrCreate(ctx context.Context, data model.GameTypeCreate) (*model.GameType, error)
	GetChecksumsByIDs(ctx context.Context, ids []string) (map[string]string, error)
	Save(ctx context.Context, gameType *model.GameType) (*model.GameType, error)
}

type howLongToBeat interface {
	GetByNameAndYear(ctx context.Context, name string, year int) (*model.HowLongToBeat, error)
}

type steamSpy interface {
	GetByAppID(ctx context.Context, id int) (*model.SteamSpy, error)
}

type cheapShark interface {
	GetByAppID(ctx context.Context, id int) (*model.CheapShark, error)
}

type settings interface {
	GetFirstOrDefault(ctx context.Context) (*model.Settings, error)
	UpdateIGDBGamesParsedByID(ctx context.Context, id string, amount int) error
}

type IGDB struct {
	repository       repository
	remoteRepository remoteRepository
	activities       activities
	platforms        platforms
	companies        companies
	tags             tags
	themes           themes
	genres           genres
	gameTypes        gameTypes
	hltb             howLongToBeat
	steamSpy         steamSpy
	cheapShark       cheapShark
	settings         settings
}

func NewIGDB(
	repository repository,
	remoteRepository remoteRepository,
	activities activities,
	platforms platforms,
	companies companies,
	tags tags,
	themes themes,
	genres genres,
	gameTypes gameTypes,
	hltb howLongToBeat,
	steamSpy steamSpy,
	cheapShark cheapShark,
	settings settings,
) *IGDB {
	return &IGDB{
		repository:       repository,
		remoteRepository: remoteRepository,
		activities:       activities,
		platforms:        platforms,
		companies:        companies,
		tags:             tags,
		themes:           themes,
		genres:           genres,
		gameTypes:        gameTypes,
		hltb:             hltb,
		steamSpy:         steamSpy,
		cheapShark:       cheapShark,
		settings:         settings,
	}
}

func (i *IGDB) ParseGames(ctx context.Context, limit uint64) error {
	gamesCount, err := i.remoteRepository.GamesCount(ctx)
	if err != nil {
		return err
	}

	settings, err := i.settings.GetFirstOrDefault(ctx)
	if err != nil {
		return err
	}

	if settings.IgdbGamesParsed() > uint(gamesCount) {
		return nil
	}

	ch, err := i.remoteRepository.ParseGames(ctx, gamesCount, 0, limit)
	if err != nil {
		return err
	}

	for msg := range ch {
		if msg.Err != nil {
			return msg.Err
		}

		err = i.saveCompaniesFromGames(ctx, msg.Games)
		if err != nil {
			return err
		}

		err = i.saveTagsFromGames(ctx, msg.Games)
		if err != nil {
			return err
		}

		err = i.saveThemesFromGames(ctx, msg.Games)
		if err != nil {
			return err
		}

		res := make([]*model.Activity, len(msg.Games))
		idsToCheck := make([]string, 0, len(msg.Games))
		for j, game := range msg.Games {
			activity, err := i.activities.GetOrCreate(ctx, model.ActivityCreate{
				IdDb:     game.Id,
				Type:     model.ActivityTypeGame,
				Name:     game.Name,
				Checksum: game.Checksum,
			})
			if err != nil {
				return err
			}

			gameTypeLocalId, err := i.repository.TableReferenceToID(ctx, TableReferenceSingle{
				Id:         game.GameType,
				TableName:  schema.CollectionGameTypes,
				PrimaryKey: schema.GameTypeSchema.Id,
				SearchKey:  schema.GameTypeSchema.IdDb,
			})
			if err != nil {
				return err
			}

			platformLocalIds, err := i.repository.TableReferenceToIDs(ctx, TableReference{
				Ids:        game.Platforms,
				TableName:  schema.CollectionPlatforms,
				PrimaryKey: schema.PlatformSchema.Id,
				SearchKey:  schema.PlatformSchema.IdDb,
			})
			if err != nil {
				return err
			}

			developerLocalIds, err := i.repository.TableReferenceToIDs(ctx, TableReference{
				Ids:        game.Developers,
				TableName:  schema.CollectionCompanies,
				PrimaryKey: schema.CompanySchema.Id,
				SearchKey:  schema.CompanySchema.IdDb,
			})
			if err != nil {
				return err
			}

			publisherLocalIds, err := i.repository.TableReferenceToIDs(ctx, TableReference{
				Ids:        game.Publishers,
				TableName:  schema.CollectionCompanies,
				PrimaryKey: schema.CompanySchema.Id,
				SearchKey:  schema.CompanySchema.IdDb,
			})
			if err != nil {
				return err
			}

			genreLocalIds, err := i.repository.TableReferenceToIDs(ctx, TableReference{
				Ids:        game.Genres,
				TableName:  schema.CollectionGenres,
				PrimaryKey: schema.GenreSchema.Id,
				SearchKey:  schema.GenreSchema.IdDb,
			})
			if err != nil {
				return err
			}

			tagLocalIds, err := i.repository.TableReferenceToIDs(ctx, TableReference{
				Ids:        game.Tags,
				TableName:  schema.CollectionTags,
				PrimaryKey: schema.TagSchema.Id,
				SearchKey:  schema.TagSchema.IdDb,
			})
			if err != nil {
				return err
			}

			themeLocalIds, err := i.repository.TableReferenceToIDs(ctx, TableReference{
				Ids:        game.Themes,
				TableName:  schema.CollectionThemes,
				PrimaryKey: schema.ThemeSchema.Id,
				SearchKey:  schema.ThemeSchema.IdDb,
			})
			if err != nil {
				return err
			}

			hltb, err := i.hltb.GetByNameAndYear(ctx, game.Name, game.ReleaseDate.Year())
			if err != nil {
				if !errors.Is(err, errs.ErrHowLongToBeatNotFound) {
					return err
				}
			} else {
				activity.SetHltbId(uint(hltb.IdDb()))
				activity.SetHltbCampaignTime(hltb.Campaign())
			}

			if game.SteamAppId > 0 {
				steamAppPrice, err := i.getPriceBySteamAppID(ctx, int(game.SteamAppId))
				if err != nil {
					return err
				}

				activity.SetSteamAppPrice(steamAppPrice)
			}

			activity.SetName(game.Name)
			activity.SetSlug(game.Slug)
			activity.SetReleaseDate(game.ReleaseDate)
			activity.SetSteamAppId(uint(game.SteamAppId))
			activity.SetCover(game.Cover)
			activity.SetChecksum(game.Checksum)

			activity.SetGameType(gameTypeLocalId)
			activity.SetPlatforms(platformLocalIds)
			activity.SetDevelopers(developerLocalIds)
			activity.SetPublishers(publisherLocalIds)
			activity.SetGenres(genreLocalIds)
			activity.SetTags(tagLocalIds)
			activity.SetThemes(themeLocalIds)

			res[j] = activity

			if !activity.IsNew() {
				idsToCheck = append(idsToCheck, activity.ID())
			}
		}

		checksums, err := i.activities.GetChecksumsByIDs(ctx, idsToCheck)
		if err != nil {
			return err
		}

		for _, activity := range res {
			checksum, ok := checksums[activity.ID()]
			if ok && activity.Checksum() == checksum {
				continue
			}

			_, err = i.activities.Save(ctx, activity)
			if err != nil {
				return err
			}
		}
	}

	return nil
}

func (i *IGDB) getPriceBySteamAppID(ctx context.Context, id int) (uint, error) {
	steamSpy, err := i.steamSpy.GetByAppID(ctx, id)
	if err != nil {
		if !errors.Is(err, errs.ErrSteamSpyNotFound) {
			return 0, err
		}
	} else {
		return steamSpy.Price(), nil
	}

	cheapShark, err := i.cheapShark.GetByAppID(ctx, id)
	if err != nil {
		if !errors.Is(err, errs.ErrCheapSharkNotFound) {
			return 0, nil
		}
	} else {
		return uint(cheapShark.Price() * 100), nil
	}

	return 0, nil
}

func (i *IGDB) saveCompaniesFromGames(ctx context.Context, games []*Game) error {
	var uniqueIds []uint64
	for _, game := range games {
		uniqueIds = append(uniqueIds, game.Developers...)
		uniqueIds = append(uniqueIds, game.Publishers...)
	}

	if len(uniqueIds) == 0 {
		return nil
	}

	slices.Sort(uniqueIds)
	uniqueIds = slices.Compact(uniqueIds)

	companies, err := i.remoteRepository.GetCompaniesByIDs(ctx, uniqueIds)
	if err != nil {
		return err
	}

	res := make([]*model.Company, len(companies))
	idsToCheck := make([]string, 0, len(companies))
	for j, company := range companies {
		c, err := i.companies.GetOrCreate(ctx, model.CompanyCreate{
			IdDb:     company.Id,
			Name:     company.Name,
			Checksum: company.Checksum,
		})
		if err != nil {
			return err
		}

		c.SetName(company.Name)
		c.SetChecksum(company.Checksum)

		res[j] = c

		if !c.IsNew() {
			idsToCheck = append(idsToCheck, c.ID())
		}
	}

	checksums, err := i.companies.GetChecksumsByIDs(ctx, idsToCheck)
	if err != nil {
		return err
	}

	for _, company := range res {
		checksum, ok := checksums[company.ID()]
		if ok && company.Checksum() == checksum {
			continue
		}

		_, err = i.companies.Save(ctx, company)
		if err != nil {
			return err
		}
	}

	return nil
}

func (i *IGDB) saveTagsFromGames(ctx context.Context, games []*Game) error {
	var uniqueIds []uint64
	for _, game := range games {
		for _, id := range game.Tags {
			uniqueIds = append(uniqueIds, id)
		}
	}

	if len(uniqueIds) == 0 {
		return nil
	}

	slices.Sort(uniqueIds)
	uniqueIds = slices.Compact(uniqueIds)

	tags, err := i.remoteRepository.GetTagsByIDs(ctx, uniqueIds)
	if err != nil {
		return err
	}

	res := make([]*model.Tag, len(tags))
	idsToCheck := make([]string, 0, len(tags))
	for j, tag := range tags {
		t, err := i.tags.GetOrCreate(ctx, model.TagCreate{
			IdDb:     tag.Id,
			Name:     tag.Name,
			Checksum: tag.Checksum,
		})
		if err != nil {
			return err
		}

		t.SetName(tag.Name)
		t.SetChecksum(tag.Checksum)

		res[j] = t

		if !t.IsNew() {
			idsToCheck = append(idsToCheck, t.ID())
		}
	}

	checksums, err := i.tags.GetChecksumsByIDs(ctx, idsToCheck)
	if err != nil {
		return err
	}

	for _, tag := range res {
		checksum, ok := checksums[tag.ID()]
		if ok && tag.Checksum() == checksum {
			continue
		}

		_, err = i.tags.Save(ctx, tag)
		if err != nil {
			return err
		}
	}

	return nil
}

func (i *IGDB) saveThemesFromGames(ctx context.Context, games []*Game) error {
	var uniqueIds []uint64
	for _, game := range games {
		for _, id := range game.Themes {
			uniqueIds = append(uniqueIds, id)
		}
	}

	if len(uniqueIds) == 0 {
		return nil
	}

	slices.Sort(uniqueIds)
	uniqueIds = slices.Compact(uniqueIds)

	themes, err := i.remoteRepository.GetThemesByIDs(ctx, uniqueIds)
	if err != nil {
		return err
	}

	res := make([]*model.Theme, len(themes))
	idsToCheck := make([]string, 0, len(themes))
	for j, theme := range themes {
		t, err := i.themes.GetOrCreate(ctx, model.ThemeCreate{
			IdDb:     theme.Id,
			Name:     theme.Name,
			Checksum: theme.Checksum,
		})
		if err != nil {
			return err
		}

		t.SetName(theme.Name)
		t.SetChecksum(theme.Checksum)

		res[j] = t

		if !t.IsNew() {
			idsToCheck = append(idsToCheck, t.ID())
		}
	}

	checksums, err := i.themes.GetChecksumsByIDs(ctx, idsToCheck)
	if err != nil {
		return err
	}

	for _, theme := range res {
		checksum, ok := checksums[theme.ID()]
		if ok && theme.Checksum() == checksum {
			continue
		}

		_, err = i.themes.Save(ctx, theme)
		if err != nil {
			return err
		}
	}

	return nil
}

func (i *IGDB) ParsePlatforms(ctx context.Context, limit uint64) error {
	ch, err := i.remoteRepository.ParsePlatformsAll(ctx, limit)
	if err != nil {
		return err
	}

	for msg := range ch {
		if msg.Err != nil {
			return msg.Err
		}

		res := make([]*model.Platform, len(msg.Platforms))
		idsToCheck := make([]string, 0, len(msg.Platforms))
		for j, platform := range msg.Platforms {
			p, err := i.platforms.GetOrCreate(ctx, model.PlatformCreate{
				IdDb:     platform.Id,
				Name:     platform.Name,
				Checksum: platform.Checksum,
			})
			if err != nil {
				return err
			}

			p.SetName(platform.Name)
			p.SetChecksum(platform.Checksum)

			res[j] = p

			if !p.IsNew() {
				idsToCheck = append(idsToCheck, p.ID())
			}
		}

		checksums, err := i.platforms.GetChecksumsByIDs(ctx, idsToCheck)
		if err != nil {
			return err
		}

		for _, platform := range res {
			checksum, ok := checksums[platform.ID()]
			if ok && platform.Checksum() == checksum {
				continue
			}

			_, err = i.platforms.Save(ctx, platform)
			if err != nil {
				return err
			}
		}
	}

	return nil
}

func (i *IGDB) ParseGenres(ctx context.Context, limit uint64) error {
	ch, err := i.remoteRepository.ParseGenresAll(ctx, limit)
	if err != nil {
		return err
	}

	for msg := range ch {
		if msg.Err != nil {
			return msg.Err
		}

		res := make([]*model.Genre, len(msg.Genres))
		idsToCheck := make([]string, 0, len(msg.Genres))
		for j, genre := range msg.Genres {
			g, err := i.genres.GetOrCreate(ctx, model.GenreCreate{
				IdDb:     genre.Id,
				Name:     genre.Name,
				Checksum: genre.Checksum,
			})
			if err != nil {
				return err
			}

			g.SetName(genre.Name)
			g.SetChecksum(genre.Checksum)

			res[j] = g

			if !g.IsNew() {
				idsToCheck = append(idsToCheck, g.ID())
			}
		}

		checksums, err := i.genres.GetChecksumsByIDs(ctx, idsToCheck)
		if err != nil {
			return err
		}

		for _, genre := range res {
			checksum, ok := checksums[genre.ID()]
			if ok && genre.Checksum() == checksum {
				continue
			}

			_, err = i.genres.Save(ctx, genre)
			if err != nil {
				return err
			}
		}
	}

	return nil
}

func (i *IGDB) ParseGameTypes(ctx context.Context, limit uint64) error {
	ch, err := i.remoteRepository.ParseGameTypesAll(ctx, limit)
	if err != nil {
		return err
	}

	for msg := range ch {
		if msg.Err != nil {
			return msg.Err
		}

		res := make([]*model.GameType, len(msg.GameTypes))
		idsToCheck := make([]string, 0, len(msg.GameTypes))
		for j, gameType := range msg.GameTypes {
			t, err := i.gameTypes.GetOrCreate(ctx, model.GameTypeCreate{
				IdDb:     gameType.Id,
				Name:     gameType.Name,
				Checksum: gameType.Checksum,
			})
			if err != nil {
				return err
			}

			t.SetName(gameType.Name)
			t.SetChecksum(gameType.Checksum)

			res[j] = t

			if !t.IsNew() {
				idsToCheck = append(idsToCheck, t.ID())
			}
		}

		checksums, err := i.gameTypes.GetChecksumsByIDs(ctx, idsToCheck)
		if err != nil {
			return err
		}

		for _, gameType := range res {
			checksum, ok := checksums[gameType.ID()]
			if ok && gameType.Checksum() == checksum {
				continue
			}

			_, err = i.gameTypes.Save(ctx, gameType)
			if err != nil {
				return err
			}
		}
	}

	return nil
}
