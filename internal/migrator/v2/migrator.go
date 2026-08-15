package v2

import (
	"github.com/pocketbase/pocketbase/core"
	"github.com/spf13/cobra"
)

func NewV2MigratorCommand(pb core.App) *cobra.Command {
	command := &cobra.Command{
		Use: "v2",
	}

	command.AddCommand(migrateActionsCommand(pb))

	return command
}
