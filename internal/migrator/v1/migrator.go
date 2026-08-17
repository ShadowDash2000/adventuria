package v1

import (
	"adventuria/internal/adventuria/model"
	"adventuria/internal/adventuria/reviews"
	"context"

	"github.com/pocketbase/pocketbase/core"
	"github.com/spf13/cobra"
)

type activities interface {
	GetByName(ctx context.Context, name string) (*model.Activity, error)
}

type cells interface {
	GetByName(ctx context.Context, name string, includeDisabled bool) (*model.CellInfo, error)
}

type reviewsService interface {
	Create(ctx context.Context, input reviews.CreateInput) (*model.Review, error)
}

func NewV1MigratorCommand(pb core.App, activities activities, cells cells, reviews reviewsService) *cobra.Command {
	command := &cobra.Command{Use: "v1"}
	command.AddCommand(migrateActionsCommand(pb, activities, cells, reviews))
	return command
}
