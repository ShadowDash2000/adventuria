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
	ParseGames(ctx context.Context, filter string, count, offset, limit uint64) (<-chan ParseGamesMessage, error)
	GamesCount(ctx context.Context, filter string) (uint64, error)
	GetCompaniesByIDs(ctx context.Context, ids []uint64) ([]*Company, error)
	GetTagsByIDs(ctx context.Context, ids []uint64) ([]*Tag, error)
	GetThemesByIDs(ctx context.Context, ids []uint64) ([]*Theme, error)
	ParsePlatformsAll(ctx context.Context, limit uint64) (<-chan ParsePlatformsMessage, error)
	ParseGenresAll(ctx context.Context, limit uint64) (<-chan ParseGenresMessage, error)
	ParseGameTypesAll(ctx context.Context, limit uint64) (<-chan ParseGameTypesMessage, error)
}

type activities interface {
	GetOrCreate(ctx context.Context, data model.ActivityCreate) (*model.Activity, error)
	Save(ctx context.Context, activity *model.Activity) (*model.Activity, error)
}

type platforms interface {
	GetOrCreate(ctx context.Context, data model.PlatformCreate) (*model.Platform, error)
	Save(ctx context.Context, platform *model.Platform) (*model.Platform, error)
}

type companies interface {
	GetOrCreate(ctx context.Context, data model.CompanyCreate) (*model.Company, error)
	Save(ctx context.Context, company *model.Company) (*model.Company, error)
}

type tags interface {
	GetOrCreate(ctx context.Context, data model.TagCreate) (*model.Tag, error)
	Save(ctx context.Context, tag *model.Tag) (*model.Tag, error)
}

type themes interface {
	GetOrCreate(ctx context.Context, data model.ThemeCreate) (*model.Theme, error)
	Save(ctx context.Context, theme *model.Theme) (*model.Theme, error)
}

type genres interface {
	GetOrCreate(ctx context.Context, data model.GenreCreate) (*model.Genre, error)
	Save(ctx context.Context, genre *model.Genre) (*model.Genre, error)
}

type gameTypes interface {
	GetOrCreate(ctx context.Context, data model.GameTypeCreate) (*model.GameType, error)
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
	GetFirst(ctx context.Context) (*model.Settings, error)
	ChangeIGDBGamesParsedByID(ctx context.Context, id string, amount int) error
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

func (i *IGDB) ParseGames(ctx context.Context, filter string, limit uint64) error {
	gamesCount, err := i.remoteRepository.GamesCount(ctx, filter)
	if err != nil {
		return err
	}

	settings, err := i.settings.GetFirst(ctx)
	if err != nil {
		return err
	}

	if settings.IgdbGamesParsed() > uint(gamesCount) {
		return nil
	}

	offset := uint64(settings.IgdbGamesParsed())
	ch, err := i.remoteRepository.ParseGames(ctx, filter, gamesCount, offset, limit)
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

		res := make([]*model.Activity, 0, len(msg.Games))
		for _, game := range msg.Games {
			activity, err := i.activities.GetOrCreate(ctx, model.ActivityCreate{
				IdDb:     game.Id,
				Type:     model.ActivityTypeGame,
				Name:     game.Name,
				Checksum: game.Checksum,
			})
			if err != nil {
				return err
			}

			if !activity.IsNew() && activity.Checksum() == game.Checksum {
				continue
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

			res = append(res, activity)
		}

		for _, activity := range res {
			_, err = i.activities.Save(ctx, activity)
			if err != nil {
				return err
			}
		}

		err = i.settings.ChangeIGDBGamesParsedByID(ctx, settings.ID(), len(msg.Games))
		if err != nil {
			return err
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

	res := make([]*model.Company, 0, len(companies))
	for _, company := range companies {
		c, err := i.companies.GetOrCreate(ctx, model.CompanyCreate{
			IdDb:     company.Id,
			Name:     company.Name,
			Checksum: company.Checksum,
		})
		if err != nil {
			return err
		}

		if !c.IsNew() && c.Checksum() == company.Checksum {
			continue
		}

		c.SetName(company.Name)
		c.SetChecksum(company.Checksum)

		res = append(res, c)
	}

	for _, company := range res {
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

	res := make([]*model.Tag, 0, len(tags))
	for _, tag := range tags {
		t, err := i.tags.GetOrCreate(ctx, model.TagCreate{
			IdDb:     tag.Id,
			Name:     tag.Name,
			Checksum: tag.Checksum,
		})
		if err != nil {
			return err
		}

		if !t.IsNew() && t.Checksum() == tag.Checksum {
			continue
		}

		t.SetName(tag.Name)
		t.SetChecksum(tag.Checksum)

		res = append(res, t)
	}

	for _, tag := range res {
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

	res := make([]*model.Theme, 0, len(themes))
	for _, theme := range themes {
		t, err := i.themes.GetOrCreate(ctx, model.ThemeCreate{
			IdDb:     theme.Id,
			Name:     theme.Name,
			Checksum: theme.Checksum,
		})
		if err != nil {
			return err
		}

		if !t.IsNew() && t.Checksum() == theme.Checksum {
			continue
		}

		t.SetName(theme.Name)
		t.SetChecksum(theme.Checksum)

		res = append(res, t)
	}

	for _, theme := range res {
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

		res := make([]*model.Platform, 0, len(msg.Platforms))
		for _, platform := range msg.Platforms {
			p, err := i.platforms.GetOrCreate(ctx, model.PlatformCreate{
				IdDb:     platform.Id,
				Name:     platform.Name,
				Checksum: platform.Checksum,
			})
			if err != nil {
				return err
			}

			if !p.IsNew() && p.Checksum() == platform.Checksum {
				continue
			}

			p.SetName(platform.Name)
			p.SetChecksum(platform.Checksum)

			res = append(res, p)
		}

		for _, platform := range res {
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

		res := make([]*model.Genre, 0, len(msg.Genres))
		for _, genre := range msg.Genres {
			g, err := i.genres.GetOrCreate(ctx, model.GenreCreate{
				IdDb:     genre.Id,
				Name:     genre.Name,
				Checksum: genre.Checksum,
			})
			if err != nil {
				return err
			}

			if !g.IsNew() && g.Checksum() == genre.Checksum {
				continue
			}

			g.SetName(genre.Name)
			g.SetChecksum(genre.Checksum)

			res = append(res, g)
		}

		for _, genre := range res {
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

		res := make([]*model.GameType, 0, len(msg.GameTypes))
		for _, gameType := range msg.GameTypes {
			t, err := i.gameTypes.GetOrCreate(ctx, model.GameTypeCreate{
				IdDb:     gameType.Id,
				Name:     gameType.Name,
				Checksum: gameType.Checksum,
			})
			if err != nil {
				return err
			}

			if !t.IsNew() && t.Checksum() == gameType.Checksum {
				continue
			}

			t.SetName(gameType.Name)
			t.SetChecksum(gameType.Checksum)

			res = append(res, t)
		}

		for _, gameType := range res {
			_, err = i.gameTypes.Save(ctx, gameType)
			if err != nil {
				return err
			}
		}
	}

	return nil
}
