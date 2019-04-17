package cmd

import (
	"fmt"
	"io/ioutil"

	"github.com/atotto/clipboard"
	"github.com/spf13/cobra"
)

func NewSwitchCmd() *cobra.Command {
	return &cobra.Command{
		Use:     "switch [contest_id]",
		Aliases: []string{"sw"},
		Short:   "Copy main program to clipboard",
		RunE: func(cmd *cobra.Command, args []string) error {
			fname, err := getProgName()
			if err != nil {
				return err
			}

			code, err := ioutil.ReadFile(fname)
			if err != nil {
				return err
			}

			if err := clipboard.WriteAll(string(code)); err != nil {
				return err
			}
			fmt.Printf("Copy %v to the clipboard\n", fname)
			return nil
		},
		SilenceUsage: true,
	}
}
