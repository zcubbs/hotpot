package eightysix

import (
	"fmt"
	"github.com/spf13/cobra"
	"github.com/zcubbs/go-k8s/k3s"
	"github.com/zcubbs/x/progress"
	"os"
	"strings"
)

var silent bool
var purgeExtraDirs []string

// Cmd represents the cook command
var Cmd = &cobra.Command{
	Use:     "86",
	Aliases: []string{"wipe", "destroy", "clear"},
	Short:   "Clear cluster",
	Long:    ``,
	Run: func(cmd *cobra.Command, args []string) {
		verbose := cmd.Flag("verbose").Value.String() == "true"
		if err := clearCluster(verbose); err != nil {
			fmt.Println(err)
		}
	},
}

func clearCluster(verbose bool) error {
	// if not silent prompt for confirmation
	if !silent {
		fmt.Println("Are you sure you want to clear the cluster? (y/n)")
		var response string
		_, err := fmt.Scanln(&response)
		if err != nil {
			return fmt.Errorf("failed to read response \n %w", err)
		}
		if response != "y" {
			fmt.Println("Aborted.")
			return nil
		}
	}

	return progress.RunTask(func() error {
		fmt.Printf("Clearing cluster...\n")
		fmt.Printf("    ├─ uninstalling k3s... \n")
		err := k3s.Uninstall(verbose)
		if err != nil && !strings.Contains(err.Error(), "no such file or directory") { // ignore if k3s is not installed
			return err
		}

		// purge extra dirs
		for _, v := range purgeExtraDirs {
			fmt.Printf("    ├─ purging extra dir %s... \n", v)
			err := purgeDir(v)
			if err != nil {
				return err
			}
		}
		fmt.Printf("    └─ ok\n")
		return nil
	}, true)
}

func purgeDir(dir string) error {
	// delete dir
	err := os.RemoveAll(dir)
	if err != nil {
		fmt.Printf("failed to delete dir %s: %s\n", dir, err)
	}

	return nil
}

func init() {
	Cmd.Flags().BoolVarP(&silent, "silent", "s", false, "silent exec")
	Cmd.Flags().StringSliceVarP(&purgeExtraDirs, "purge", "p", []string{}, "extra dirs to purge")
}
