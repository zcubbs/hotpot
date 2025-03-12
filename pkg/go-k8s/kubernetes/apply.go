package kubernetes

import (
	"fmt"
	"github.com/zcubbs/hotpot/pkg/x/bash"
	"github.com/zcubbs/hotpot/pkg/x/yaml"
	"os"
	"time"
)

func ApplyManifest(manifestTmpl string, data interface{}, debug bool) error {
	return ApplyManifestWithKc(manifestTmpl, data, "", debug)
}

func ApplyManifestWithKc(manifestTmpl string, data interface{}, kubeconfig string, debug bool) error {
	b, err := yaml.ApplyTmpl(manifestTmpl, data, debug)
	if err != nil {
		return fmt.Errorf("failed to apply template \n %w", err)
	}

	// generate tmp file name
	tmpDir := os.TempDir()
	fn := fmt.Sprintf("%s%stmpManifest_%s.yaml",
		tmpDir,
		string(os.PathSeparator),
		time.Unix(time.Now().Unix(), 0).Format("20060102150405"),
	)
	if debug {
		fmt.Printf("tmp manifest file: %s\n", fn)
	}

	// write tmp manifest
	err = os.WriteFile(fn, b, 0600)
	if err != nil {
		return fmt.Errorf("failed to write tmp manifest \n %w", err)
	}

	fmt.Printf("Executing command: kubectl --kubeconfig %s apply -f %s\n", kubeconfig, fn)
	if kubeconfig != "" {
		err = bash.ExecuteCmd("kubectl", debug, "--kubeconfig", kubeconfig, "apply", "-f", fn)
		if err != nil {
			return fmt.Errorf("failed to apply manifest \n %w", err)
		}
	} else {
		err = bash.ExecuteCmd("kubectl", debug, "apply", "-f", fn)
		if err != nil {
			return fmt.Errorf("failed to apply manifest \n %w", err)
		}
	}

	// delete tmp manifest
	err = os.Remove(fn)
	if err != nil {
		return fmt.Errorf("failed to delete tmp manifest \n %w", err)
	}
	return nil
}
