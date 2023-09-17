package cook

import (
	"fmt"
	"github.com/spf13/cobra"
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
		err := cook()
		if err != nil {
			fmt.Println(err)
		}
	},
}

func cook() error {

	return nil
}

func init() {
	Cmd.Flags().StringVarP(&configPath, "config", "c", "", "yaml config file path (default is ./config.yaml)")

	_ = Cmd.MarkFlagRequired("config")
}
