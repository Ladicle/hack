package cmd

import (
	"fmt"

	"github.com/Ladicle/hack/pkg/contest"
	"github.com/spf13/cobra"
)

func NewInitCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:     "init",
		Aliases: []string{"i"},
		Short:   "Initialize workspace for the contest",
	}

	atcoder := &cobra.Command{
		Use:     "atcoder [contest_id]",
		Aliases: []string{"at", "a"},
		Example: "at abc123",
		Short:   "Initialize workspace for AtCoder",
		RunE: func(cmd *cobra.Command, args []string) error {
			if len(args) != 1 {
				return fmt.Errorf("required argument \"contest id\" not set (eg. abc123)")
			}
			id := args[0]
			if err := contest.NewAtCoder(id).SetupQuizDir(); err != nil {
				return err
			}
			fmt.Printf("Completed initialization for AtCoder %v\n", id)
			return nil
		},
		SilenceUsage: true,
	}

	cmd.AddCommand(atcoder)

	return cmd
}
