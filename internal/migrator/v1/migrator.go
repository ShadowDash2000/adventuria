package v1

import (
	"adventuria/internal/adventuria/model"
	"adventuria/internal/adventuria/reviews"
	"adventuria/internal/migrator"
	migratorModel "adventuria/internal/migrator/model"
	"context"

	"github.com/spf13/cobra"
)

type playersRepository interface {
	GetIDByName(ctx context.Context, name string) (string, error)
}

type actionsRepository interface {
	Create(ctx context.Context, actionCreate migratorModel.ActionImport) error
}

type seasonsRepository interface {
	GetBySlug(ctx context.Context, slug string) (*model.Season, error)
}

type activities interface {
	GetByName(ctx context.Context, name string) (*model.Activity, error)
}

type cells interface {
	GetIDByName(ctx context.Context, name string) (string, error)
}

type reviewsService interface {
	Create(ctx context.Context, input reviews.CreateInput) (*model.Review, error)
}

type playerProgressService interface {
	Save(ctx context.Context, progress *model.PlayerProgress) (*model.PlayerProgress, error)
}

type playerStatsService interface {
	Save(ctx context.Context, stats *model.PlayerStats) (*model.PlayerStats, error)
}

type worldsService interface {
	GetDefault(ctx context.Context) (*model.World, error)
}

func NewV1MigratorCommand(
	transactionRunner migrator.TransactionRunner,
	players playersRepository,
	actions actionsRepository,
	seasons seasonsRepository,
	playerProgress playerProgressService,
	playerStats playerStatsService,
	worlds worldsService,
	cells cells,
	activities activities,
	reviews reviewsService,
) *cobra.Command {
	command := &cobra.Command{
		Use: "v1",
	}

	command.AddCommand(migrateActionsCommand(
		transactionRunner,
		players,
		actions,
		cells,
		activities,
		reviews,
	))
	command.AddCommand(migrateProgressCommand(
		transactionRunner,
		players,
		seasons,
		playerProgress,
		playerStats,
		worlds,
	))

	return command
}
