package syncd

import (
	"fmt"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/spf13/cobra"
	"github.com/zcubbs/hotpot/pkg/syncd/ui"
)

var configCmd = &cobra.Command{
	Use:   "config",
	Short: "Configure the sync daemon",
	Long:  `Interactive configuration for the recipe sync daemon.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		fmt.Println("🔧 Configuring hotpot-syncd...")

		p := tea.NewProgram(ui.InitialConfigModel())
		model, err := p.Run()
		if err != nil {
			return fmt.Errorf("failed to run configuration UI: %w", err)
		}

		if m, ok := model.(*ui.ConfigModel); ok && m.Done() {
			fmt.Println("✅ Configuration saved successfully")
		}

		return nil
	},
}
