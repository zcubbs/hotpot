package sops

import (
	"fmt"
	"runtime"
)

// GetSopsInstallCommand returns the command to install sops based on the OS
func GetSopsInstallCommand() (string, error) {
	switch runtime.GOOS {
	case "windows":
		// Windows installation command (using Chocolatey as an example)
		return "choco install sops", nil
	case "darwin":
		// macOS installation command (using Homebrew)
		return "brew install sops", nil
	case "linux":
		// Linux installation command (assuming use of a Debian-based system)
		return "wget https://github.com/mozilla/sops/releases/download/v3.7.1/sops-v3.7.1.linux -O /usr/local/bin/sops && chmod +x /usr/local/bin/sops", nil
	default:
		return "", fmt.Errorf("unsupported operating system")
	}
}
