package cmd

import (
	"fmt"

	"github.com/Ladicle/hack/pkg/config"
	"github.com/spf13/cobra"
)

func NewSwitchCmd() *cobra.Command {
	return &cobra.Command{
		Use:     "switch [host]/[contest_id]",
		Aliases: []string{"sw"},
		Short:   "Copy main program to clipboard",
		RunE: func(cmd *cobra.Command, args []string) error {
			if len(args) != 1 {
				return fmt.Errorf("required argument \"[host]/[contest_id]\" not set")
			}
			config.SetCurrent(args[0])
			fmt.Printf("Switch to %v\n", args[0])
			return nil
		},
		SilenceUsage: true,
	}
}
