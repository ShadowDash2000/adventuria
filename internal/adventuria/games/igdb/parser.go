package igdb

import (
	"adventuria/internal/adventuria"
	"adventuria/internal/adventuria/games"
	"context"
	"fmt"
	"strconv"

	"github.com/bestnite/go-igdb"
	"github.com/bestnite/go-igdb/endpoint"
	pb "github.com/bestnite/go-igdb/proto"
	"github.com/pocketbase/pocketbase/tools/types"
	"google.golang.org/protobuf/proto"
)

type Parser struct {
	filter string
	client *igdb.Client

	extGameSources map[uint64]*pb.ExternalGameSource
}

func NewParser(clientID, clientSecret, filter string) *Parser {
	p := &Parser{
		client: igdb.New(clientID, clientSecret),
		filter: filter,
	}

	return p
}

type ParsePlatformsMessage struct {
	Platforms []games.Platform
	Err       error
}

func (p *Parser) ParsePlatformsAll(ctx context.Context, limit uint64) (chan ParsePlatformsMessage, error) {
	count, err := p.client.Platforms.Count()
	if err != nil {
		return nil, err
	}

	return p.ParsePlatforms(ctx, count, limit)
}

func (p *Parser) ParsePlatforms(ctx context.Context, count, limit uint64) (chan ParsePlatformsMessage, error) {
	ch := make(chan ParsePlatformsMessage)

	go func() {
		defer close(ch)

		for offset := uint64(0); offset < count; offset += limit {
			select {
			case <-ctx.Done():
				return
			default:
				platforms, err := p.client.Platforms.Paginated(0, limit)
				if err != nil {
					ch <- ParsePlatformsMessage{Err: err}
					return
				}

				res := ParsePlatformsMessage{Platforms: make([]games.Platform, len(platforms))}
				for i, platform := range platforms {
					res.Platforms[i] = games.Platform{
						IdDb:     platform.GetId(),
						Name:     platform.GetName(),
						Checksum: strconv.FormatUint(platform.GetId(), 10),
					}
				}

				ch <- res
			}
		}
	}()

	return ch, nil
}

type ParseCompaniesMessage struct {
	Companies []games.Company
	Err       error
}

func (p *Parser) ParseCompaniesAll(ctx context.Context, limit uint64) (chan ParseCompaniesMessage, error) {
	count, err := p.client.Companies.Count()
	if err != nil {
		return nil, err
	}

	return p.ParseCompanies(ctx, count, limit)
}

func (p *Parser) ParseCompanies(ctx context.Context, count, limit uint64) (chan ParseCompaniesMessage, error) {
	ch := make(chan ParseCompaniesMessage)

	go func() {
		defer close(ch)

		for offset := uint64(0); offset < count; offset += limit {
			select {
			case <-ctx.Done():
				return
			default:
				companies, err := p.client.Companies.Paginated(offset, limit)
				if err != nil {
					ch <- ParseCompaniesMessage{Err: err}
					return
				}

				res := ParseCompaniesMessage{Companies: make([]games.Company, len(companies))}
				for i, company := range companies {
					res.Companies[i] = games.Company{
						IdDb:     company.GetId(),
						Name:     company.GetName(),
						Checksum: strconv.FormatUint(company.GetId(), 10),
					}
				}

				ch <- res
			}
		}
	}()

	return ch, nil
}

type ParseGamesMessage struct {
	Games []games.Game
	Err   error
}

func (p *Parser) ParseGamesAll(ctx context.Context, limit uint64) (chan ParseGamesMessage, error) {
	count, err := p.gamesCount()
	if err != nil {
		return nil, err
	}

	return p.ParseGames(ctx, count, limit)
}

func (p *Parser) ParseGames(ctx context.Context, count, limit uint64) (chan ParseGamesMessage, error) {
	if limit > count {
		limit = count
	}

	ch := make(chan ParseGamesMessage)

	go func() {
		defer close(ch)

		for offset := uint64(0); offset < count; offset += limit {
			select {
			case <-ctx.Done():
				return
			default:
				gamesIgdb, err := p.gamesPaginated(offset, limit)
				if err != nil {
					ch <- ParseGamesMessage{Err: err}
					return
				}
				steamAppIds, err := p.getSteamAppIds(gamesIgdb)
				if err != nil {
					ch <- ParseGamesMessage{Err: err}
					return
				}
				covers, err := p.fetchCovers(gamesIgdb)
				if err != nil {
					ch <- ParseGamesMessage{Err: err}
					return
				}
				companies, err := p.fetchInvolvedCompanies(gamesIgdb)
				if err != nil {
					ch <- ParseGamesMessage{Err: err}
					return
				}

				res := ParseGamesMessage{Games: make([]games.Game, len(gamesIgdb))}
				for i, game := range gamesIgdb {
					releaseDate, err := types.ParseDateTime(game.GetFirstReleaseDate().AsTime())
					if err != nil {
						ch <- ParseGamesMessage{Err: err}
						return
					}

					platformIds := make([]uint64, len(game.GetPlatforms()))
					for i, platform := range game.GetPlatforms() {
						platformIds[i] = platform.GetId()
					}

					var developersIds []uint64
					var publishersIds []uint64
					for _, involvedCompany := range game.GetInvolvedCompanies() {
						company, ok := companies[involvedCompany.GetId()]
						if !ok {
							continue
						}

						if company.GetDeveloper() {
							developersIds = append(developersIds, company.GetId())
							continue
						}
						if company.GetPublisher() {
							publishersIds = append(publishersIds, company.GetId())
							continue
						}
					}

					genreIds := make([]uint64, len(game.GetGenres()))
					for i, genre := range game.GetGenres() {
						genreIds[i] = genre.GetId()
					}

					res.Games[i] = games.Game{
						IdDb:        game.GetId(),
						Name:        game.GetName(),
						ReleaseDate: releaseDate,
						Platforms: games.CollectionReference{
							Ids:        platformIds,
							Collection: adventuria.CollectionPlatforms,
						},
						Developers: games.CollectionReference{
							Ids:        developersIds,
							Collection: adventuria.CollectionCompanies,
						},
						Publishers: games.CollectionReference{
							Ids:        publishersIds,
							Collection: adventuria.CollectionCompanies,
						},
						Genres: games.CollectionReference{
							Ids:        genreIds,
							Collection: adventuria.CollectionGenres,
						},
						SteamAppId: steamAppIds[game.GetId()],
						Cover:      covers[game.GetId()],
						Checksum:   strconv.FormatUint(game.GetId(), 10),
					}
				}

				ch <- res
			}
		}
	}()

	return ch, nil
}

func (p *Parser) gamesPaginated(offset, limit uint64) ([]*pb.Game, error) {
	return p.client.Games.Query(
		fmt.Sprintf("where %s; offset %d; limit %d; fields *; sort id asc;", p.filter, offset, limit),
	)
}

func (p *Parser) gamesCount() (uint64, error) {
	resp, err := p.client.Request(
		"POST",
		fmt.Sprintf("https://api.igdb.com/v4/%s/count.pb", endpoint.EPGames),
		fmt.Sprintf("where %s;", p.filter),
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

const (
	extGameSourceSteam = "Steam"
)

func (p *Parser) getSteamAppIds(games []*pb.Game) (map[uint64]uint64, error) {
	var extGamesIds []uint64
	for _, game := range games {
		for _, extGame := range game.GetExternalGames() {
			extGamesIds = append(extGamesIds, extGame.GetId())
		}
	}

	extGames, err := p.client.ExternalGames.GetByIDs(extGamesIds)
	if err != nil {
		return nil, err
	}

	extGameSources, err := p.externalGameSources()
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

		return nil, fmt.Errorf("getSteamAppId(): Can't parse Steam App Id = %v", extGame.GetUid())
	}

	return res, nil
}

func (p *Parser) externalGameSources() (map[uint64]*pb.ExternalGameSource, error) {
	if p.extGameSources != nil {
		return p.extGameSources, nil
	}

	var err error
	p.extGameSources, err = p.fetchExternalGameSources()
	if err != nil {
		return nil, err
	}

	return p.extGameSources, nil
}

func (p *Parser) fetchExternalGameSources() (map[uint64]*pb.ExternalGameSource, error) {
	count, err := p.client.ExternalGameSources.Count()
	if err != nil {
		return nil, err
	}

	sources, err := p.client.ExternalGameSources.Paginated(0, count)
	if err != nil {
		return nil, err
	}

	res := make(map[uint64]*pb.ExternalGameSource)
	for _, source := range sources {
		res[source.GetId()] = source
	}

	return res, nil
}

func (p *Parser) fetchCovers(games []*pb.Game) (map[uint64]string, error) {
	gameIds := make([]uint64, len(games))
	for i, game := range games {
		gameIds[i] = game.GetCover().GetId()
	}

	covers, err := p.client.Covers.GetByIDs(gameIds)
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

func (p *Parser) fetchInvolvedCompanies(games []*pb.Game) (map[uint64]*pb.InvolvedCompany, error) {
	var companyIds []uint64
	for _, game := range games {
		for _, company := range game.GetInvolvedCompanies() {
			companyIds = append(companyIds, company.GetId())
		}
	}

	companies, err := p.client.InvolvedCompanies.GetByIDs(companyIds)
	if err != nil {
		return nil, err
	}

	res := make(map[uint64]*pb.InvolvedCompany)
	for _, company := range companies {
		res[company.GetId()] = company
	}

	return res, nil
}

type ParseGenresMessage struct {
	Genres []games.Genre
	Err    error
}

func (p *Parser) ParseGenresAll(ctx context.Context, limit uint64) (chan ParseGenresMessage, error) {
	count, err := p.client.Genres.Count()
	if err != nil {
		return nil, err
	}

	return p.ParseGenres(ctx, count, limit)
}

func (p *Parser) ParseGenres(ctx context.Context, count, limit uint64) (chan ParseGenresMessage, error) {
	ch := make(chan ParseGenresMessage)

	go func() {
		defer close(ch)

		for offset := uint64(0); offset < count; offset += limit {
			select {
			case <-ctx.Done():
				return
			default:
				genres, err := p.client.Genres.Paginated(offset, limit)
				if err != nil {
					ch <- ParseGenresMessage{Err: err}
					return
				}

				res := ParseGenresMessage{Genres: make([]games.Genre, len(genres))}
				for i, genre := range genres {
					res.Genres[i] = games.Genre{
						IdDb:     genre.GetId(),
						Name:     genre.GetName(),
						Checksum: strconv.FormatUint(genre.GetId(), 10),
					}
				}

				ch <- res
			}
		}
	}()

	return ch, nil
}
