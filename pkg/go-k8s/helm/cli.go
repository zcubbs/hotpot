// Package helm.
/*
Copyright © 2023 zcubbs https://github.com/zcubbs
*/
package helm

import (
	"fmt"
	"github.com/zcubbs/hotpot/pkg/x/bash"
)

const (
	TmpHelmScript    = "/tmp/get_helm.sh"
	InstallScriptUrl = "https://raw.githubusercontent.com/helm/helm/main/scripts/get-helm-3"
)

func InstallCli(debug bool) error {
	if debug {
		fmt.Printf("installing helm %s", "curl -fsSL -o "+TmpHelmScript+InstallScriptUrl)
	}

	err := bash.ExecuteCmd(
		"curl",
		debug,
		"-fsSL",
		"-o",
		TmpHelmScript,
		"https://raw.githubusercontent.com/helm/helm/main/scripts/get-helm-3",
	)
	if err != nil {
		return err
	}

	// chmod 700 get_helm.sh
	err = bash.Chmod("/tmp/get_helm.sh", 0700, debug)
	if err != nil {
		return err
	}

	// sh ./get_helm.sh
	_, err = bash.ExecuteScript(
		TmpHelmScript,
		debug,
		TmpHelmScript,
	)
	if err != nil {
		return err
	}

	return nil
}

func IsHelmInstalled() (bool, error) {
	// helm version
	err := bash.ExecuteCmd(
		"helm",
		false,
		"version",
	)
	if err != nil {
		return false, err
	}
	return true, nil
}
