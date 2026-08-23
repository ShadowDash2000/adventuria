package v2

import (
	"adventuria/internal/adventuria/model"
	"adventuria/internal/adventuria/reviews"
	"context"

	"github.com/pocketbase/pocketbase/core"
	"github.com/spf13/cobra"
)

type reviewsService interface {
	Create(ctx context.Context, input reviews.CreateInput) (*model.Review, error)
}

type itemsService interface {
	GetByName(ctx context.Context, name string) (*model.Item, error)
}

func NewV2MigratorCommand(pb core.App, reviews reviewsService, items itemsService) *cobra.Command {
	command := &cobra.Command{
		Use: "v2",
	}

	command.AddCommand(migrateActionsCommand(pb, reviews, items))

	return command
}
