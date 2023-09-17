package cook

import (
	"github.com/spf13/cobra"
	"github.com/zcubbs/hotpot/pkg/recipe"
	"github.com/zcubbs/x/must"
	"github.com/zcubbs/x/progress"
	"github.com/zcubbs/x/style"
)

var (
	configPath string
)

// Cmd represents the cook command
var Cmd = &cobra.Command{
	Use:   "cook",
	Short: "cook commands",
	Long:  ``,
	Run: func(cmd *cobra.Command, args []string) {
		style.PrintColoredHeader("Cooking the cluster")
		must.Succeed(progress.RunTask(cook(), true))
	},
}

func cook() func() error {
	return func() error {
		return recipe.Cook(configPath)
	}
}

func init() {
	Cmd.Flags().StringVarP(&configPath, "config", "c", "./recipe.yaml", "yaml config file path (default is ./recipe.yaml)")

	_ = Cmd.MarkFlagRequired("config")
}
