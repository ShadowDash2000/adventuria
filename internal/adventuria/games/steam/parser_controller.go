package steam

import (
	"adventuria/internal/adventuria"
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
	if err := p.parseSteamSpy(ctx); err != nil {
		adventuria.PocketBase.Logger().Error("Failed to parse steam spy", "error", err)
		return
	}
}

func (p *ParserController) parseSteamSpy(ctx context.Context) error {
	apps, err := p.parser.FetchLatestRelease(ctx)
	if err != nil {
		return err
	}

	for _, app := range apps {
		select {
		case <-ctx.Done():
			return ctx.Err()
		default:
		}

		if ok := p.isAppExist(ctx, app.AppId); ok {
			continue
		}

		appRecord := NewSteamSpyRecordFromRecord(
			core.NewRecord(
				adventuria.GameCollections.Get(adventuria.CollectionSteamSpy),
			),
		)
		appRecord.SetIdDb(app.AppId)
		appRecord.SetName(app.Name)
		appRecord.SetPrice(uint(app.Price))
		if err = adventuria.PocketBase.Save(appRecord.ProxyRecord()); err != nil {
			return err
		}
	}

	return nil
}

func (p *ParserController) isAppExist(ctx context.Context, id uint) bool {
	err := adventuria.PocketBase.
		RecordQuery(adventuria.GameCollections.Get(adventuria.CollectionSteamSpy)).
		WithContext(ctx).
		Where(dbx.HashExp{"id_db": id}).
		One(&core.Record{})
	return err == nil
}
