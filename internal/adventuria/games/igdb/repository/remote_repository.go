package repository

import (
	igdbparser "adventuria/internal/adventuria/games/igdb"
	"context"
	"fmt"
	"strconv"
	"strings"

	"git.nite07.com/shadod/go-igdb"
	"git.nite07.com/shadod/go-igdb/endpoint"
	pb "git.nite07.com/shadod/go-igdb/proto"
	"google.golang.org/protobuf/proto"
)

type RemoteRepository struct {
	client *igdb.Client
	filter string

	extGameSources map[uint64]*pb.ExternalGameSource
}

func NewRemoteRepository(clientID, clientSecret, filter string) *RemoteRepository {
	return &RemoteRepository{
		client: igdb.New(clientID, clientSecret),
		filter: filter,
	}
}

func (r *RemoteRepository) ParseGamesAll(ctx context.Context, offset, limit uint64) (<-chan igdbparser.ParseGamesMessage, uint64, error) {
	count, err := r.GamesCount(ctx)
	if err != nil {
		return nil, 0, err
	}

	ch, err := r.ParseGames(ctx, count, offset, limit)
	if err != nil {
		return nil, 0, err
	}

	return ch, count, nil
}

func (r *RemoteRepository) ParseGames(ctx context.Context, count, offset, limit uint64) (<-chan igdbparser.ParseGamesMessage, error) {
	if limit > count {
		limit = count
	}

	ch := make(chan igdbparser.ParseGamesMessage)

	go func() {
		defer close(ch)

		for ; offset < count; offset += limit {
			select {
			case <-ctx.Done():
				return
			default:
				gamesIgdb, err := r.gamesPaginated(ctx, offset, limit)
				if err != nil {
					ch <- igdbparser.ParseGamesMessage{Err: err}
					return
				}
				steamAppIds, err := r.getSteamAppIds(ctx, gamesIgdb)
				if err != nil {
					ch <- igdbparser.ParseGamesMessage{Err: err}
					return
				}
				covers, err := r.fetchCovers(ctx, gamesIgdb)
				if err != nil {
					ch <- igdbparser.ParseGamesMessage{Err: err}
					return
				}
				companies, err := r.fetchInvolvedCompanies(ctx, gamesIgdb)
				if err != nil {
					ch <- igdbparser.ParseGamesMessage{Err: err}
					return
				}

				res := igdbparser.ParseGamesMessage{Games: make([]*igdbparser.Game, len(gamesIgdb))}
				for i, game := range gamesIgdb {
					platformIds := make([]uint64, len(game.GetPlatforms()))
					for i, platform := range game.GetPlatforms() {
						platformIds[i] = platform.GetId()
					}

					var developerIds []uint64
					var publisherIds []uint64
					for _, involvedCompany := range game.GetInvolvedCompanies() {
						company, ok := companies[involvedCompany.GetId()]
						if !ok {
							continue
						}

						if company.GetDeveloper() {
							developerIds = append(developerIds, company.GetCompany().GetId())
							continue
						}
						if company.GetPublisher() {
							publisherIds = append(publisherIds, company.GetCompany().GetId())
							continue
						}
					}

					genreIds := make([]uint64, len(game.GetGenres()))
					for i, genre := range game.GetGenres() {
						genreIds[i] = genre.GetId()
					}

					keywordIds := make([]uint64, len(game.GetKeywords()))
					for i, keyword := range game.GetKeywords() {
						keywordIds[i] = keyword.GetId()
					}

					themeIds := make([]uint64, len(game.GetThemes()))
					for i, theme := range game.GetThemes() {
						themeIds[i] = theme.GetId()
					}

					res.Games[i] = &igdbparser.Game{
						Id:          strconv.FormatUint(game.GetId(), 10),
						Name:        game.GetName(),
						Slug:        game.GetSlug(),
						ReleaseDate: game.GetFirstReleaseDate().AsTime(),
						Platforms:   platformIds,
						Developers:  developerIds,
						Publishers:  publisherIds,
						Genres:      genreIds,
						Tags:        keywordIds,
						Themes:      themeIds,
						GameType:    game.GetGameType().GetId(),
						SteamAppId:  steamAppIds[game.GetId()],
						Cover:       covers[game.GetId()],
						Checksum:    game.GetChecksum(),
					}
				}

				ch <- res
			}
		}
	}()

	return ch, nil
}

func (r *RemoteRepository) GamesCount(ctx context.Context) (uint64, error) {
	resp, err := r.client.Request(
		ctx,
		"POST",
		fmt.Sprintf("https://api.igdb.com/v4/%s/count.pb", endpoint.EPGames),
		fmt.Sprintf("where %s;", r.filter),
	)
	if err != nil {
		return 0, err
	}

	var res pb.Count
	if err = proto.Unmarshal(resp.Body(), &res); err != nil {
		return 0, fmt.Errorf("failed to unmarshal: %w", err)
	}

	return uint64(res.Count), nil
}

func (r *RemoteRepository) gamesPaginated(ctx context.Context, offset, limit uint64) ([]*pb.Game, error) {
	return r.client.Games.Query(ctx,
		fmt.Sprintf("where %s; offset %d; limit %d; fields *; sort id asc;", r.filter, offset, limit),
	)
}

const (
	extGameSourceSteam = "Steam"
)

func (r *RemoteRepository) getSteamAppIds(ctx context.Context, games []*pb.Game) (map[uint64]uint64, error) {
	var extGamesIds []uint64
	for _, game := range games {
		for _, extGame := range game.GetExternalGames() {
			extGamesIds = append(extGamesIds, extGame.GetId())
		}
	}

	extGames, err := r.client.ExternalGames.GetByIDs(ctx, extGamesIds)
	if err != nil {
		return nil, err
	}

	extGameSources, err := r.externalGameSources(ctx)
	if err != nil {
		return nil, err
	}

	res := make(map[uint64]uint64)
	for _, extGame := range extGames {
		extGameSource, ok := extGameSources[extGame.GetExternalGameSource().GetId()]
		if !ok {
			continue
		}
		if extGameSource.GetName() != extGameSourceSteam {
			continue
		}
		if uid, err := strconv.ParseUint(extGame.GetUid(), 10, 64); err == nil {
			res[extGame.GetGame().GetId()] = uid
			continue
		}
		if uids := strings.Split(extGame.GetUid(), ","); len(uids) > 0 {
			if uid, err := strconv.ParseUint(uids[0], 10, 64); err == nil {
				res[extGame.GetGame().GetId()] = uid
				continue
			}
		}

		return nil, fmt.Errorf("getSteamAppId(): Can't parse Steam App Id = %v", extGame.GetUid())
	}

	return res, nil
}

func (r *RemoteRepository) externalGameSources(ctx context.Context) (map[uint64]*pb.ExternalGameSource, error) {
	if r.extGameSources != nil {
		return r.extGameSources, nil
	}

	var err error
	r.extGameSources, err = r.fetchExternalGameSources(ctx)
	if err != nil {
		return nil, err
	}

	return r.extGameSources, nil
}

func (r *RemoteRepository) fetchExternalGameSources(ctx context.Context) (map[uint64]*pb.ExternalGameSource, error) {
	count, err := r.client.ExternalGameSources.Count(ctx)
	if err != nil {
		return nil, err
	}

	sources, err := r.client.ExternalGameSources.Paginated(ctx, 0, count)
	if err != nil {
		return nil, err
	}

	res := make(map[uint64]*pb.ExternalGameSource)
	for _, source := range sources {
		res[source.GetId()] = source
	}

	return res, nil
}

func (r *RemoteRepository) fetchCovers(ctx context.Context, games []*pb.Game) (map[uint64]string, error) {
	gameIds := make([]uint64, len(games))
	for i, game := range games {
		gameIds[i] = game.GetCover().GetId()
	}

	covers, err := r.client.Covers.GetByIDs(ctx, gameIds)
	if err != nil {
		return nil, err
	}

	res := make(map[uint64]string)
	for _, cover := range covers {
		res[cover.GetGame().GetId()] = fmt.Sprintf(
			"https://images.igdb.com/igdb/image/upload/t_cover_big/%s.webp",
			cover.GetImageId(),
		)
	}

	return res, nil
}

func (r *RemoteRepository) fetchInvolvedCompanies(ctx context.Context, games []*pb.Game) (map[uint64]*pb.InvolvedCompany, error) {
	var companyIds []uint64
	for _, game := range games {
		for _, company := range game.GetInvolvedCompanies() {
			companyIds = append(companyIds, company.GetId())
		}
	}

	companies, err := r.client.InvolvedCompanies.GetByIDs(ctx, companyIds)
	if err != nil {
		return nil, err
	}

	res := make(map[uint64]*pb.InvolvedCompany)
	for _, company := range companies {
		res[company.GetId()] = company
	}

	return res, nil
}

func (r *RemoteRepository) ParsePlatformsAll(ctx context.Context, limit uint64) (<-chan igdbparser.ParsePlatformsMessage, error) {
	count, err := r.client.Platforms.Count(ctx)
	if err != nil {
		return nil, err
	}

	return r.ParsePlatforms(ctx, count, limit)
}

func (r *RemoteRepository) ParsePlatforms(ctx context.Context, count, limit uint64) (<-chan igdbparser.ParsePlatformsMessage, error) {
	ch := make(chan igdbparser.ParsePlatformsMessage)

	go func() {
		defer close(ch)

		for offset := uint64(0); offset < count; offset += limit {
			select {
			case <-ctx.Done():
				return
			default:
				platforms, err := r.client.Platforms.Paginated(ctx, 0, limit)
				if err != nil {
					ch <- igdbparser.ParsePlatformsMessage{Err: err}
					return
				}

				res := igdbparser.ParsePlatformsMessage{Platforms: make([]*igdbparser.Platform, len(platforms))}
				for i, platform := range platforms {
					res.Platforms[i] = &igdbparser.Platform{
						Id:       strconv.FormatUint(platform.GetId(), 10),
						Name:     platform.GetName(),
						Checksum: platform.GetChecksum(),
					}
				}

				ch <- res
			}
		}
	}()

	return ch, nil
}

func (r *RemoteRepository) ParseGenresAll(ctx context.Context, limit uint64) (<-chan igdbparser.ParseGenresMessage, error) {
	count, err := r.client.Genres.Count(ctx)
	if err != nil {
		return nil, err
	}

	return r.ParseGenres(ctx, count, limit)
}

func (r *RemoteRepository) ParseGenres(ctx context.Context, count, limit uint64) (<-chan igdbparser.ParseGenresMessage, error) {
	ch := make(chan igdbparser.ParseGenresMessage)

	go func() {
		defer close(ch)

		for offset := uint64(0); offset < count; offset += limit {
			select {
			case <-ctx.Done():
				return
			default:
				genres, err := r.client.Genres.Paginated(ctx, offset, limit)
				if err != nil {
					ch <- igdbparser.ParseGenresMessage{Err: err}
					return
				}

				res := igdbparser.ParseGenresMessage{Genres: make([]*igdbparser.Genre, len(genres))}
				for i, genre := range genres {
					res.Genres[i] = &igdbparser.Genre{
						Id:       strconv.FormatUint(genre.GetId(), 10),
						Name:     genre.GetName(),
						Checksum: genre.GetChecksum(),
					}
				}

				ch <- res
			}
		}
	}()

	return ch, nil
}

func (r *RemoteRepository) ParseGameTypesAll(ctx context.Context, limit uint64) (<-chan igdbparser.ParseGameTypesMessage, error) {
	count, err := r.client.GameTypes.Count(ctx)
	if err != nil {
		return nil, err
	}

	return r.ParseGameTypes(ctx, count, limit)
}

func (r *RemoteRepository) ParseGameTypes(ctx context.Context, count, limit uint64) (<-chan igdbparser.ParseGameTypesMessage, error) {
	ch := make(chan igdbparser.ParseGameTypesMessage)

	go func() {
		defer close(ch)

		for offset := uint64(0); offset < count; offset += limit {
			select {
			case <-ctx.Done():
				return
			default:
				gameTypes, err := r.client.GameTypes.Paginated(ctx, offset, limit)
				if err != nil {
					ch <- igdbparser.ParseGameTypesMessage{Err: err}
					return
				}

				res := igdbparser.ParseGameTypesMessage{GameTypes: make([]*igdbparser.GameType, len(gameTypes))}
				for i, gameType := range gameTypes {
					res.GameTypes[i] = &igdbparser.GameType{
						Id:       strconv.FormatUint(gameType.GetId(), 10),
						Name:     gameType.GetType(),
						Checksum: gameType.GetChecksum(),
					}
				}

				ch <- res
			}
		}
	}()

	return ch, nil
}

func (r *RemoteRepository) GetCompaniesByIDs(ctx context.Context, ids []uint64) ([]*igdbparser.Company, error) {
	if len(ids) == 0 {
		return nil, nil
	}

	comps, err := r.client.Companies.GetByIDs(ctx, ids)
	if err != nil {
		return nil, err
	}

	res := make([]*igdbparser.Company, len(comps))
	for i, c := range comps {
		res[i] = &igdbparser.Company{
			Id:       strconv.FormatUint(c.GetId(), 10),
			Name:     c.GetName(),
			Checksum: c.GetChecksum(),
		}
	}

	return res, nil
}

func (r *RemoteRepository) GetTagsByIDs(ctx context.Context, ids []uint64) ([]*igdbparser.Tag, error) {
	if len(ids) == 0 {
		return nil, nil
	}

	tags, err := r.client.Keywords.GetByIDs(ctx, ids)
	if err != nil {
		return nil, err
	}

	res := make([]*igdbparser.Tag, len(tags))
	for i, tag := range tags {
		res[i] = &igdbparser.Tag{
			Id:       strconv.FormatUint(tag.GetId(), 10),
			Name:     tag.GetName(),
			Checksum: tag.GetChecksum(),
		}
	}

	return res, nil
}

func (r *RemoteRepository) GetThemesByIDs(ctx context.Context, ids []uint64) ([]*igdbparser.Theme, error) {
	if len(ids) == 0 {
		return nil, nil
	}

	themes, err := r.client.Themes.GetByIDs(ctx, ids)
	if err != nil {
		return nil, err
	}

	res := make([]*igdbparser.Theme, len(themes))
	for i, theme := range themes {
		res[i] = &igdbparser.Theme{
			Id:       strconv.FormatUint(theme.GetId(), 10),
			Name:     theme.GetName(),
			Checksum: theme.GetChecksum(),
		}
	}

	return res, nil
}
