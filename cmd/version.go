package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

func NewVersionCmd() *cobra.Command {
	return &cobra.Command{
		Use:     "version",
		Aliases: []string{"v"},
		Short:   "Show this command version",
		Run: func(cmd *cobra.Command, args []string) {
			fmt.Printf("%v -- %v\n", gitRepo, version)
		},
	}
}
