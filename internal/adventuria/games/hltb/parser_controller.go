package hltb

import (
	"adventuria/internal/adventuria"
	"adventuria/internal/adventuria/games"
	"context"
	"errors"
)

type ParserController struct {
	parser *Parser
	ctx    context.Context
}

func New(ctx context.Context) (*ParserController, error) {
	p, err := NewParser()
	if err != nil {
		return nil, err
	}

	return &ParserController{
		parser: p,
		ctx:    ctx,
	}, nil
}

func (p *ParserController) Parse(limit int) {
	if err := p.parseTime(p.ctx, limit); err != nil {
		adventuria.PocketBase.Logger().Error("Failed to parse time", "error", err)
		return
	}
}

func (p *ParserController) parseTime(ctx context.Context, limit int) error {
	for {
		gameRecords, err := p.getGamesWithoutTime(limit)
		if err != nil {
			return err
		}

		if len(gameRecords) == 0 {
			break
		}

		for _, game := range gameRecords {
			gameTime, err := p.parser.ParseTime(ctx, game.Name())
			if err != nil {
				if errors.Is(err, ErrGameNotFound) {
					adventuria.PocketBase.Logger().Debug("parseTime(): Game not found", "game", game.Name())
					game.SetCampaignTime(0)
				} else {
					return err
				}
			} else {
				game.SetHltbID(gameTime.GameID)
				game.SetCampaignTime(gameTime.Campaign)
			}

			err = adventuria.PocketBase.Save(game.ProxyRecord())
			if err != nil {
				return err
			}
		}
	}
	return nil
}

func (p *ParserController) getGamesWithoutTime(limit int) ([]games.GameRecord, error) {
	records, err := adventuria.PocketBase.FindRecordsByFilter(
		adventuria.GameCollections.Get(adventuria.CollectionGames),
		"campaign_time = -1",
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
