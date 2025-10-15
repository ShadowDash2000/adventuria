package igdb

import (
	"adventuria/internal/adventuria"
	"adventuria/internal/adventuria/games"
	"fmt"
	"net/url"
	"strconv"
	"strings"

	"github.com/bestnite/go-igdb"
	"github.com/bestnite/go-igdb/endpoint"
	pb "github.com/bestnite/go-igdb/proto"
	"github.com/pocketbase/pocketbase/tools/types"
	"google.golang.org/protobuf/proto"
)

type Parser struct {
	filter string
	client *igdb.Client
}

func NewParser(clientID, clientSecret, filter string) *Parser {
	p := &Parser{
		client: igdb.New(clientID, clientSecret),
		filter: filter,
	}

	return p
}

func (p *Parser) ParsePlatforms() (chan []games.Platform, error) {
	count, err := p.client.Platforms.Count()
	if err != nil {
		return nil, err
	}

	ch := make(chan []games.Platform)

	go func() {
		defer close(ch)

		const limit = uint64(500)
		for offset := uint64(0); offset < count; offset += limit {
			platforms, err := p.client.Platforms.Paginated(0, limit)
			if err != nil {
				return
			}

			res := make([]games.Platform, len(platforms))
			for i, platform := range platforms {
				res[i] = games.Platform{
					IdDb:     platform.GetId(),
					Name:     platform.GetName(),
					Checksum: platform.GetChecksum(),
				}
			}

			ch <- res
		}
	}()

	return ch, nil
}

func (p *Parser) ParseCompanies() (chan []games.Company, error) {
	count, err := p.client.Companies.Count()
	if err != nil {
		return nil, err
	}

	ch := make(chan []games.Company)

	go func() {
		defer close(ch)

		const limit = uint64(500)
		for offset := uint64(0); offset < count; offset += limit {
			companies, err := p.client.Companies.Paginated(offset, limit)
			if err != nil {
				return
			}

			res := make([]games.Company, len(companies))
			for i, company := range companies {
				res[i] = games.Company{
					IdDb:     company.GetId(),
					Name:     company.GetName(),
					Checksum: company.GetChecksum(),
				}
			}

			ch <- res
		}
	}()

	return ch, nil
}

func (p *Parser) ParseGames() (chan []games.Game, error) {
	count, err := p.gamesCount()
	if err != nil {
		return nil, err
	}

	ch := make(chan []games.Game)

	go func() {
		defer close(ch)

		websiteTypes, err := p.fetchWebsiteTypes()
		if err != nil {
			return
		}

		const limit = uint64(500)
		for offset := uint64(0); offset < count; offset += limit {
			gamesIgdb, err := p.client.Games.Paginated(offset, limit)
			if err != nil {
				return
			}

			websites, err := p.fetchWebsites(gamesIgdb)
			if err != nil {
				return
			}

			res := make([]games.Game, len(gamesIgdb))
			for i, game := range gamesIgdb {
				releaseDate, err := types.ParseDateTime(game.GetFirstReleaseDate().AsTime())
				if err != nil {
					return
				}

				platformIds := make([]uint64, len(game.GetPlatforms()))
				for i, platform := range game.GetPlatforms() {
					platformIds[i] = platform.Id
				}

				companyIds := make([]uint64, len(game.GetInvolvedCompanies()))
				for i, company := range game.GetInvolvedCompanies() {
					companyIds[i] = company.Id
				}

				steamAppId := uint64(0)
				for _, website := range game.GetWebsites() {
					websiteFull, ok := websites[website.GetId()]
					if !ok {
						continue
					}

					websiteType, ok := websiteTypes[websiteFull.GetType().GetId()]
					if !ok {
						continue
					}

					if websiteType.GetType() != "Steam" {
						continue
					}

					steamAppId = p.parseSteamAppId(websiteFull.GetUrl())
				}

				res[i] = games.Game{
					IdDb:        game.GetId(),
					Name:        game.GetName(),
					ReleaseDate: releaseDate,
					Platforms: games.CollectionReference{
						Ids:        platformIds,
						Collection: adventuria.CollectionPlatforms,
					},
					Companies: games.CollectionReference{
						Ids:        companyIds,
						Collection: adventuria.CollectionCompanies,
					},
					SteamAppId: steamAppId,
					Cover:      game.GetCover().GetUrl(),
					Checksum:   game.GetChecksum(),
				}
			}

			ch <- res
		}
	}()

	return ch, nil
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

func (p *Parser) fetchWebsites(games []*pb.Game) (map[uint64]*pb.Website, error) {
	var websiteIds []uint64
	for _, game := range games {
		for _, website := range game.GetWebsites() {
			websiteIds = append(websiteIds, website.GetId())
		}
	}

	websites, err := p.client.Websites.GetByIDs(websiteIds)
	if err != nil {
		return nil, err
	}

	res := make(map[uint64]*pb.Website)
	for _, website := range websites {
		res[website.GetId()] = website
	}

	return res, err
}

func (p *Parser) fetchWebsiteTypes() (map[uint64]*pb.WebsiteType, error) {
	count, err := p.client.WebsiteTypes.Count()
	if err != nil {
		return nil, err
	}

	websiteTypes, err := p.client.WebsiteTypes.Paginated(0, count)
	if err != nil {
		return nil, err
	}

	res := make(map[uint64]*pb.WebsiteType)
	for _, websiteType := range websiteTypes {
		res[websiteType.GetId()] = websiteType
	}

	return res, err
}

func (p *Parser) parseSteamAppId(rawUrl string) uint64 {
	u, err := url.Parse(rawUrl)
	if err != nil {
		return 0
	}

	parts := strings.Split(strings.Trim(u.Path, "/"), "/")
	if len(parts) < 2 || parts[0] != "app" {
		return 0
	}

	idStr := parts[1]
	appId, err := strconv.ParseUint(idStr, 10, 64)
	if err != nil {
		return 0
	}

	return appId
}
