package main

import (
	v2tov3 "adventuria/internal/migrator/v2-to-v3"
	"log"

	"github.com/pocketbase/pocketbase"
	"github.com/spf13/cobra"
)

func main() {
	pb := pocketbase.New()

	command := &cobra.Command{
		Use: "migrate-data",
	}

	command.AddCommand(v2tov3.NewV2ToV3MigratorCommand(pb))

	pb.RootCmd.AddCommand(command)

	err := pb.Start()
	if err != nil {
		log.Fatal(err)
	}
}
