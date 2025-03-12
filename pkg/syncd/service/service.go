package service

import (
	"fmt"
	"github.com/zcubbs/hotpot/pkg/x/bash"
	"os"
)

const (
	serviceName     = "hotpot-syncd"
	serviceTemplate = `[Unit]
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
)

// Install creates and installs the systemd service
func Install() error {
	servicePath := fmt.Sprintf("/etc/systemd/system/%s.service", serviceName)

	// Write service file
	if err := os.WriteFile(servicePath, []byte(serviceTemplate), 0600); err != nil {
		return fmt.Errorf("failed to write service file: %w", err)
	}

	// Reload systemd
	if err := reloadDaemon(); err != nil {
		return fmt.Errorf("failed to reload systemd: %w", err)
	}

	return nil
}

// Enable enables and starts the systemd service
func Enable() error {
	// Enable service
	if err := bash.ExecuteCmd("systemctl", false, "enable", serviceName); err != nil {
		return fmt.Errorf("failed to enable service: %w", err)
	}

	// Start service
	if err := bash.ExecuteCmd("systemctl", false, "start", serviceName); err != nil {
		return fmt.Errorf("failed to start service: %w", err)
	}

	return nil
}

// Disable stops and disables the systemd service
func Disable() error {
	// Stop service
	if err := bash.ExecuteCmd("systemctl", false, "stop", serviceName); err != nil {
		return fmt.Errorf("failed to stop service: %w", err)
	}

	// Disable service
	if err := bash.ExecuteCmd("systemctl", false, "disable", serviceName); err != nil {
		return fmt.Errorf("failed to disable service: %w", err)
	}

	return nil
}

// Status returns the current status of the service
func Status() (string, error) {
	output, err := bash.ExecuteCmdWithOutput("systemctl", "false", "status", serviceName)
	if err != nil {
		return "", fmt.Errorf("failed to get service status: %w", err)
	}
	return output, nil
}

// reloadDaemon reloads the systemd daemon
func reloadDaemon() error {
	return bash.ExecuteCmd("systemctl", false, "daemon-reload")
}

// Uninstall removes the systemd service
func Uninstall() error {
	// Disable and stop service first
	_ = Disable()

	// Remove service file
	servicePath := fmt.Sprintf("/etc/systemd/system/%s.service", serviceName)
	if err := os.Remove(servicePath); err != nil && !os.IsNotExist(err) {
		return fmt.Errorf("failed to remove service file: %w", err)
	}

	// Reload systemd
	if err := reloadDaemon(); err != nil {
		return fmt.Errorf("failed to reload systemd: %w", err)
	}

	return nil
}
