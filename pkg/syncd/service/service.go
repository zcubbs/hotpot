package service

import (
	"fmt"
	"github.com/zcubbs/hotpot/pkg/x/bash"
	"os"
	"path/filepath"
	"runtime"
)

const (
	serviceName = "hotpot-syncd"
)

// getServiceConfig returns the appropriate service configuration based on OS
func getServiceConfig() (string, string) {
	switch runtime.GOOS {
	case "darwin":
		homeDir, err := os.UserHomeDir()
		if err != nil {
			homeDir = "/Users/" + os.Getenv("USER")
		}
		plistPath := filepath.Join(homeDir, "Library", "LaunchAgents", "com.zcubbs.hotpot.syncd.plist")
		plistTemplate := `<?xml version="1.0" encoding="UTF-8"?>
<!DOCTYPE plist PUBLIC "-//Apple//DTD PLIST 1.0//EN" "http://www.apple.com/DTDs/PropertyList-1.0.dtd">
<plist version="1.0">
<dict>
	<key>Label</key>
	<string>com.zcubbs.hotpot.syncd</string>
	<key>ProgramArguments</key>
	<array>
		<string>/usr/local/bin/hotpot</string>
		<string>syncd</string>
		<string>run</string>
	</array>
	<key>RunAtLoad</key>
	<true/>
	<key>KeepAlive</key>
	<true/>
	<key>StandardOutPath</key>
	<string>/tmp/hotpot-syncd.log</string>
	<key>StandardErrorPath</key>
	<string>/tmp/hotpot-syncd.err</string>
</dict>
</plist>`
		return plistPath, plistTemplate

	default: // Linux
		servicePath := fmt.Sprintf("/etc/systemd/system/%s.service", serviceName)
		systemdTemplate := `[Unit]
Description=Hotpot Recipe Sync Daemon
After=network.target

[Service]
Type=simple
ExecStart=/usr/local/bin/hotpot syncd run
Restart=always
RestartSec=10
User=root

[Install]
WantedBy=multi-user.target
`
		return servicePath, systemdTemplate
	}
}

// Install creates and installs the service
func Install() error {
	servicePath, serviceTemplate := getServiceConfig()

	// Create parent directory if it doesn't exist
	if err := os.MkdirAll(filepath.Dir(servicePath), 0750); err != nil {
		return fmt.Errorf("failed to create service directory: %w", err)
	}

	// Write service file
	if err := os.WriteFile(servicePath, []byte(serviceTemplate), 0600); err != nil {
		return fmt.Errorf("failed to write service file: %w", err)
	}

	// Load and start the service based on OS
	switch runtime.GOOS {
	case "darwin":
		// Load the service
		if err := bash.ExecuteCmd("launchctl", false, "load", servicePath); err != nil {
			return fmt.Errorf("failed to load service: %w", err)
		}

		// Start the service
		if err := bash.ExecuteCmd("launchctl", false, "start", "com.zcubbs.hotpot.syncd"); err != nil {
			return fmt.Errorf("failed to start service: %w", err)
		}

	default: // Linux
		// Reload systemd
		if err := reloadDaemon(); err != nil {
			return fmt.Errorf("failed to reload systemd: %w", err)
		}
	}

	return nil
}

// Enable enables and starts the service
func Enable() error {
	switch runtime.GOOS {
	case "darwin":
		servicePath, _ := getServiceConfig()
		// Load the service
		if err := bash.ExecuteCmd("launchctl", false, "load", servicePath); err != nil {
			return fmt.Errorf("failed to load service: %w", err)
		}

		// Start the service
		if err := bash.ExecuteCmd("launchctl", false, "start", "com.zcubbs.hotpot.syncd"); err != nil {
			return fmt.Errorf("failed to start service: %w", err)
		}

	default: // Linux
		// Enable service
		if err := bash.ExecuteCmd("systemctl", false, "enable", serviceName); err != nil {
			return fmt.Errorf("failed to enable service: %w", err)
		}

		// Start service
		if err := bash.ExecuteCmd("systemctl", false, "start", serviceName); err != nil {
			return fmt.Errorf("failed to start service: %w", err)
		}
	}

	return nil
}

// Disable stops and disables the service
func Disable() error {
	switch runtime.GOOS {
	case "darwin":
		servicePath, _ := getServiceConfig()
		// Stop the service
		if err := bash.ExecuteCmd("launchctl", false, "stop", "com.zcubbs.hotpot.syncd"); err != nil {
			return fmt.Errorf("failed to stop service: %w", err)
		}

		// Unload the service
		if err := bash.ExecuteCmd("launchctl", false, "unload", servicePath); err != nil {
			return fmt.Errorf("failed to unload service: %w", err)
		}

	default: // Linux
		// Stop service
		if err := bash.ExecuteCmd("systemctl", false, "stop", serviceName); err != nil {
			return fmt.Errorf("failed to stop service: %w", err)
		}

		// Disable service
		if err := bash.ExecuteCmd("systemctl", false, "disable", serviceName); err != nil {
			return fmt.Errorf("failed to disable service: %w", err)
		}
	}

	return nil
}

// Status returns the current status of the service
func Status() (string, error) {
	switch runtime.GOOS {
	case "darwin":
		output, err := bash.ExecuteCmdWithOutput("launchctl", "false", "list", "com.zcubbs.hotpot.syncd")
		if err != nil {
			return "", fmt.Errorf("failed to get service status: %w", err)
		}
		return output, nil

	default: // Linux
		output, err := bash.ExecuteCmdWithOutput("systemctl", "false", "status", serviceName)
		if err != nil {
			return "", fmt.Errorf("failed to get service status: %w", err)
		}
		return output, nil
	}
}

// reloadDaemon reloads the systemd daemon
func reloadDaemon() error {
	return bash.ExecuteCmd("systemctl", false, "daemon-reload")
}

// Uninstall removes the service
func Uninstall() error {
	// Disable and stop service first
	_ = Disable()

	// Remove service file
	servicePath, _ := getServiceConfig()
	if err := os.Remove(servicePath); err != nil && !os.IsNotExist(err) {
		return fmt.Errorf("failed to remove service file: %w", err)
	}

	// Reload daemon if on Linux
	if runtime.GOOS != "darwin" {
		if err := reloadDaemon(); err != nil {
			return fmt.Errorf("failed to reload systemd: %w", err)
		}
	}

	return nil
}
