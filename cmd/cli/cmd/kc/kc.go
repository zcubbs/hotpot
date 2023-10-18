package kc

import (
	"fmt"
	"github.com/spf13/cobra"
	"github.com/zcubbs/go-k8s/k3s"
	"github.com/zcubbs/x/must"
	"github.com/zcubbs/x/progress"
)

var url string

// Cmd represents the cook command
var Cmd = &cobra.Command{
	Use:   "kc",
	Short: "Print kubeconfig",
	Long: `Use -u or --url to override server url.
Example: hotpot kc -u https://localhost:6443`,
	Run: func(cmd *cobra.Command, args []string) {
		verbose := cmd.Flag("verbose").Value.String() == "true"
		must.Succeed(progress.RunTask(printKc(verbose), true))
	},
}

func printKc(verbose bool) func() error {
	return func() error {
		kc, err := getKubeConfig("", verbose)
		if err != nil {
			return err
		}

		if kc != "" {
			err := k3s.PrintKubeconfig(kc, url)
			if err != nil {
				return fmt.Errorf("failed to print kubeconfig \n %w", err)
			}
		} else {
			return fmt.Errorf("kubeconfig not found")
		}

		return nil
	}
}

func init() {
	Cmd.Flags().StringVarP(&url, "url", "u", "", "override kubeconfig url")
}
