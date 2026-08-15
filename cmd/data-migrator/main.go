package main

import (
	"adventuria/internal/adventuria"
	v1 "adventuria/internal/migrator/v1"
	v2 "adventuria/internal/migrator/v2"
	"log"

	"github.com/pocketbase/pocketbase"
	"github.com/spf13/cobra"
)

func main() {
	pb := pocketbase.New()

	command := &cobra.Command{
		Use: "migrate-data",
	}

	registry := adventuria.NewRegistry(pb, pb.Logger())
	command.AddCommand(v1.NewV1MigratorCommand(pb, registry.Activities(), registry.Cells()))
	command.AddCommand(v2.NewV2MigratorCommand(pb))

	pb.RootCmd.AddCommand(command)

	err := pb.Start()
	if err != nil {
		log.Fatal(err)
	}
}
