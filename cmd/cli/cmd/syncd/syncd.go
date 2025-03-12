package syncd

import (
	"github.com/spf13/cobra"
)

// Cmd represents the syncd command
var Cmd = &cobra.Command{
	Use:   "syncd",
	Short: "Recipe synchronization daemon commands",
	Long:  `Manage the recipe synchronization daemon that keeps your recipe files in sync with a git repository.`,
}

func init() {
	Cmd.AddCommand(configCmd)
	Cmd.AddCommand(enableCmd)
	Cmd.AddCommand(disableCmd)
}
