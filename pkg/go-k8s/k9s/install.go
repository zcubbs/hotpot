package k9s

import (
	"fmt"
	"github.com/zcubbs/hotpot/pkg/x/bash"
)

// Install installs k9s
func Install(debug bool) error {
	fmt.Printf("🔨 Installing k9s...\n")

	cmd := "curl -fsSL -o k9s.tar.gz https://github.com/derailed/k9s/releases/latest/download/k9s_Linux_amd64.tar.gz && " +
		"tar -xzf k9s.tar.gz && " +
		"mv k9s /usr/local/bin/ && " +
		"rm k9s.tar.gz"

	err := bash.ExecuteCmd(cmd, debug)
	if err != nil {
		return fmt.Errorf("failed to install k9s: %w", err)
	}

	return nil
}

// Uninstall uninstalls k9s
func Uninstall(debug bool) error {
	fmt.Printf("🔨 Uninstalling k9s...\n")

	cmd := "rm -f /usr/local/bin/k9s"

	err := bash.ExecuteCmd(cmd, debug)
	if err != nil {
		return fmt.Errorf("failed to uninstall k9s: %w", err)
	}

	return nil
}
