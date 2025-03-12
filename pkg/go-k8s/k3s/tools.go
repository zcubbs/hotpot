// Package k3s.
/*
Copyright © 2023 zcubbs https://github.com/zcubbs
*/
package k3s

import (
	"fmt"
	"github.com/zcubbs/hotpot/pkg/x/bash"
	"golang.org/x/text/cases"
	"golang.org/x/text/language"
	"runtime"
)

func InstallK9s(debug bool) error {
	if debug {
		fmt.Printf("Arch: %s_%s\n", runtime.GOOS, runtime.GOARCH)
	}

	// wget -O /tmp/k9s.tar.gz https://github.com/derailed/k9s/releases/download/v0.27.3/k9s_Linux_amd64.tar.gz
	err := bash.ExecuteCmd(
		"wget",
		debug,
		"-O",
		"/tmp/k9s.tar.gz",
		fmt.Sprintf("https://github.com/derailed/k9s/releases/download/v0.27.3/k9s_%s_%s.tar.gz",
			cases.Title(language.Und).String(runtime.GOOS),
			runtime.GOARCH,
		),
	)
	if err != nil {
		return err
	}

	// tar -xvf k9s.tar.gz
	err = bash.ExtractTarGzWithFile("/tmp/k9s.tar.gz", "k9s", "/tmp", debug)
	if err != nil {
		return err
	}

	// mv /tmp/k9s /usr/local/bin/k9s
	err = bash.ExecuteCmd(
		"mv",
		debug,
		"/tmp/k9s",
		"/usr/local/bin/k9s",
	)
	if err != nil {
		return err
	}

	// chmod +x /usr/local/bin/k9s
	err = bash.Chmod("/usr/local/bin/k9s", 0700, debug)
	if err != nil {
		return fmt.Errorf("chmod /usr/local/bin/k9s: %w", err)
	}

	return nil
}
