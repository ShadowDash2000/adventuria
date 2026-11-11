package main

import (
	"adventuria/internal/adventuria"
	migratorRepo "adventuria/internal/migrator/repository"
	v1 "adventuria/internal/migrator/v1"
	v2 "adventuria/internal/migrator/v2"
	"adventuria/pkg/pbtransaction"
	"context"
	"log"

	"ganigeorgiev.com/pocketbase/pocketbase"
	"ganigeorgiev.com/pocketbase/pocketbase/core"
	"github.com/spf13/cobra"
)

func main() {
	pb := pocketbase.New()

	command := &cobra.Command{
		Use: "migrate-data",
	}

	registry := adventuria.NewRegistry(pb, pb.Logger())
	playersRepo := migratorRepo.NewPlayers(pb)
	actionsRepo := migratorRepo.NewActions(pb)
	seasonsRepo := migratorRepo.NewSeasons(pb)
	cellsRepo := migratorRepo.NewCells(pb)

	transactionRunner := func(ctx context.Context, fn func(ctx context.Context) error) error {
		return pbtransaction.RunInTransaction(ctx, pb, func(ctx context.Context, txApp core.App) error {
			return fn(ctx)
		})
	}

	command.AddCommand(v1.NewV1MigratorCommand(
		transactionRunner,
		playersRepo,
		actionsRepo,
		seasonsRepo,
		registry.PlayerProgress(),
		registry.PlayerStats(),
		registry.Worlds(),
		cellsRepo,
		registry.Activities(),
		registry.Reviews(),
	))
	command.AddCommand(v2.NewV2MigratorCommand(
		transactionRunner,
		playersRepo,
		actionsRepo,
		seasonsRepo,
		registry.Reviews(),
		registry.Items(),
		registry.PlayerProgress(),
		registry.PlayerStats(),
		registry.Worlds(),
	))

	pb.RootCmd.AddCommand(command)

	err := pb.Start()
	if err != nil {
		log.Fatal(err)
	}
}
