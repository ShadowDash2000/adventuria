package v2_to_v3

import (
	"github.com/pocketbase/pocketbase/core"
	"github.com/spf13/cobra"
)

func NewV2ToV3MigratorCommand(pb core.App) *cobra.Command {
	command := &cobra.Command{
		Use: "v2-to-v3",
	}

	command.AddCommand(migrateActionsCommand(pb))

	return command
}
