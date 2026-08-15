package main

import (
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

	command.AddCommand(v2.NewV2MigratorCommand(pb))

	pb.RootCmd.AddCommand(command)

	err := pb.Start()
	if err != nil {
		log.Fatal(err)
	}
}
