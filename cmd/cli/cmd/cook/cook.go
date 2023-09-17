package cook

import (
	"github.com/spf13/cobra"
	"github.com/zcubbs/hotpot/pkg/recipe"
	"github.com/zcubbs/x/must"
	"github.com/zcubbs/x/progress"
	"github.com/zcubbs/x/style"
)

var (
	recipePath string
)

// Cmd represents the cook command
var Cmd = &cobra.Command{
	Use:   "cook",
	Short: "cook commands",
	Long:  ``,
	Run: func(cmd *cobra.Command, args []string) {
		must.Succeed(progress.RunTask(cook(), true))
	},
}

func cook() func() error {
	return func() error {
		return recipe.Cook(recipePath,
			recipe.Hooks{
				Pre: func() error {
					style.PrintColoredHeader("Cooking the cluster")
					return nil
				},
				Post: func() error {
					return nil
				},
			},
		)
	}
}

func init() {
	Cmd.Flags().StringVarP(&recipePath, "recipe", "r", "./recipe.yaml", "yaml config file path (default is ./recipe.yaml)")

	_ = Cmd.MarkFlagRequired("recipe")
}
