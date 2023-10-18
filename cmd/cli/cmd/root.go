package cmd

import (
	"fmt"
	"github.com/spf13/cobra"
	"github.com/zcubbs/hotpot/cmd/cli/cmd/cook"
	"github.com/zcubbs/hotpot/cmd/cli/cmd/eightysix"
	"github.com/zcubbs/hotpot/cmd/cli/cmd/kc"
	"os"
)

var (
	Version string
	Commit  string
	Date    string
)

var (
	rootCmd = &cobra.Command{
		Use:   "",
		Short: "",
		Long:  "",
	}

	versionCmd = &cobra.Command{
		Use:   "version",
		Short: "Print the version number of hotpot",
		Long:  "",
		Run: func(cmd *cobra.Command, args []string) {
			fmt.Println(getVersion())
		},
	}

	aboutCmd = &cobra.Command{
		Use:   "about",
		Short: "Print information about hotpot",
		Long:  "",
		Run: func(cmd *cobra.Command, args []string) {
			About()
		},
	}
)

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func init() {
	rootCmd.PersistentFlags().BoolP("verbose", "v", false, "verbose output")
	rootCmd.DisableSuggestions = true

	rootCmd.AddCommand(versionCmd)
	rootCmd.AddCommand(aboutCmd)
	rootCmd.AddCommand(cook.Cmd)
	rootCmd.AddCommand(kc.Cmd)
	rootCmd.AddCommand(eightysix.Cmd)
}

func About() {
	fmt.Println("HotPot: Cooking Your Cluster to Perfection 🍲")
	fmt.Println(getFullVersion())
	fmt.Println(getDescription())
	fmt.Println("Author: zakaria.elbouwab")
	fmt.Println("License: MIT")
	fmt.Println("Repository: https://github.com/zcubbs/hotpot")
}

func getVersion() string {
	return fmt.Sprintf("v%s", Version)
}

func getFullVersion() string {
	return fmt.Sprintf(`
Version: v%s
Commit: %s
Date: %s
`, Version, Commit, Date)
}

func getDescription() string {
	return `
HotPot is your go-to CLI utility that marries the simplicity of cooking
with the robustness of Kubernetes deployments. Drawing inspiration from
crafting and culinary arts, HotPot serves up k3s clusters based on your
specific recipe (configuration). It aims to provide a reproducible,
consistent, and reliable way to deploy your clusters and applications.
It can also work with none k3s clusters by disabling the k3s feature.
`
}
