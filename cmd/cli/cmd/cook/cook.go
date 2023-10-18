package cook

import (
	"fmt"
	"github.com/charmbracelet/lipgloss"
	"github.com/spf13/cobra"
	"github.com/zcubbs/hotpot/pkg/recipe"
	"github.com/zcubbs/x/must"
	"github.com/zcubbs/x/progress"
)

var (
	recipePath string
)

// Cmd represents the cook command
var Cmd = &cobra.Command{
	Use:   "cook",
	Short: "Cook commands",
	Long: `Cook cmd runs the recipe. Example: hotpot cook -r ./recipe.yaml.
Add -v or --verbose to enable verbose output.`,
	Run: func(cmd *cobra.Command, args []string) {
		verbose := cmd.Flag("verbose").Value.String() == "true"
		must.Succeed(progress.RunTask(cook(verbose), true))
	},
}

func cook(verbose bool) func() error {
	return func() error {
		return recipe.Cook(recipePath,
			recipe.Hooks{
				Pre: func(r *recipe.Recipe) error {
					style := lipgloss.NewStyle().Bold(true)
					r.Debug = verbose
					fmt.Println(style.Render("🍲 Cooking..."))
					return nil
				},
				Post: func(r *recipe.Recipe) error {
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
