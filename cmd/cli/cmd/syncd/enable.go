package syncd

import (
	"fmt"
	"github.com/spf13/cobra"
	"github.com/zcubbs/hotpot/pkg/syncd/service"
)

var enableCmd = &cobra.Command{
	Use:   "enable",
	Short: "Enable and start the sync daemon",
	Long:  `Enable and start the hotpot-syncd systemd service.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		fmt.Println("🔌 Enabling hotpot-syncd service...")

		if err := service.Enable(); err != nil {
			return fmt.Errorf("failed to enable service: %w", err)
		}

		status, err := service.Status()
		if err != nil {
			return fmt.Errorf("service enabled but failed to get status: %w", err)
		}

		fmt.Printf("✅ Service enabled successfully\n\nStatus:\n%s\n", status)
		return nil
	},
}
