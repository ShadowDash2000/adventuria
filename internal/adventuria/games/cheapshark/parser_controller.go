package cheapshark

import (
	"adventuria/internal/adventuria"
	"adventuria/internal/adventuria/schema"
	"context"

	"github.com/pocketbase/dbx"
	"github.com/pocketbase/pocketbase/core"
)

type ParserController struct {
	parser *Parser
}

func New() *ParserController {
	return &ParserController{
		parser: NewParser(),
	}
}

func (p *ParserController) Parse(ctx context.Context) {
	if err := p.parseCheapshark(ctx); err != nil {
		adventuria.PocketBase.Logger().Error("Failed to parse cheapshark", "error", err)
		return
	}
}

func (p *ParserController) parseCheapshark(ctx context.Context) error {
	deals, err := p.parser.FetchLatestRelease(ctx)
	if err != nil {
		return err
	}

	for _, deal := range deals {
		select {
		case <-ctx.Done():
			return ctx.Err()
		default:
		}

		if ok := p.isDealExist(ctx, deal.SteamAppID); ok {
			continue
		}

		dealRecord := NewCheapsharkRecordFromRecord(
			core.NewRecord(
				adventuria.GameCollections.Get(schema.CollectionCheapshark),
			),
		)
		dealRecord.SetIdDb(deal.SteamAppID)
		dealRecord.SetName(deal.Title)
		dealRecord.SetPrice(deal.NormalPrice)
		if err = adventuria.PocketBase.Save(dealRecord.ProxyRecord()); err != nil {
			return err
		}
	}

	return nil
}

func (p *ParserController) isDealExist(ctx context.Context, id uint) bool {
	err := adventuria.PocketBase.
		RecordQuery(adventuria.GameCollections.Get(schema.CollectionCheapshark)).
		WithContext(ctx).
		Where(dbx.HashExp{"id_db": id}).
		One(&core.Record{})
	return err == nil
}
