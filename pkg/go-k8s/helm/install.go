package helm

import (
	"fmt"
	"github.com/zcubbs/hotpot/pkg/x/bash"
)

// Install installs helm
func Install(debug bool) error {
	fmt.Printf("🔨 Installing helm...\n")

	cmd := "curl -fsSL -o get_helm.sh https://raw.githubusercontent.com/helm/helm/master/scripts/get-helm-3 && " +
		"chmod 700 get_helm.sh && " +
		"./get_helm.sh && " +
		"rm get_helm.sh"

	err := bash.ExecuteCmd(cmd, debug)
	if err != nil {
		return fmt.Errorf("failed to install helm: %w", err)
	}

	return nil
}

// Uninstall uninstalls helm
func Uninstall(debug bool) error {
	fmt.Printf("🔨 Uninstalling helm...\n")

	cmd := "rm -f /usr/local/bin/helm"

	err := bash.ExecuteCmd(cmd, debug)
	if err != nil {
		return fmt.Errorf("failed to uninstall helm: %w", err)
	}

	return nil
}
