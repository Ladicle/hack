package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

var (
	gitRepo string
	version string
)

var rootCmd = &cobra.Command{
	Use:   "hack",
	Short: "Hack assists your programming contest.",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Printf("%v -- %v\n", gitRepo, version)
	},
}

func Execute() error {
	return rootCmd.Execute()
}
