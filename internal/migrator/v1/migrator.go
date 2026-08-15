package v1

import (
	"adventuria/internal/adventuria/model"
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

func NewV1MigratorCommand(pb core.App, activities activities, cells cells) *cobra.Command {
	command := &cobra.Command{Use: "v1"}
	command.AddCommand(migrateActionsCommand(pb, activities, cells))
	return command
}
