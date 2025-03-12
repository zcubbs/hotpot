// Package os provides a set of functions to interact with the operating system.
/*
Copyright © 2023 zcubbs https://github.com/zcubbs
*/
package os

import (
	"fmt"
	"github.com/zcubbs/hotpot/pkg/x/bash"
)

func RestartSystemdService(serviceName string, verbose bool) error {
	// Restart the service
	err := executeCmd("systemctl", verbose, "restart", serviceName)
	if err != nil {

		return fmt.Errorf("Failed to restart service: %w\n", err)
	}

	// Check the status of the service
	err = executeCmd("systemctl", verbose, "status", serviceName)
	if err != nil {
		return fmt.Errorf("Failed to get service status: %w\n", err)
	}

	return nil
}

func executeCmd(command string, verbose bool, args ...string) error {
	err := bash.ExecuteCmd(command, verbose, args...)
	if err != nil {
		cmd := fmt.Sprintf("%s %s", command, args)
		return fmt.Errorf("Failed to execute command: %s\n %w\n", cmd, err)
	}

	return nil
}
