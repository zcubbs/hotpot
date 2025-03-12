package syncd

import (
	"fmt"
	"github.com/spf13/cobra"
	"github.com/zcubbs/hotpot/pkg/syncd/service"
)

var disableCmd = &cobra.Command{
	Use:   "disable",
	Short: "Disable and stop the sync daemon",
	Long:  `Disable and stop the hotpot-syncd systemd service.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		fmt.Println("🔌 Disabling hotpot-syncd service...")

		if err := service.Disable(); err != nil {
			return fmt.Errorf("failed to disable service: %w", err)
		}

		fmt.Println("✅ Service disabled successfully")
		return nil
	},
}
