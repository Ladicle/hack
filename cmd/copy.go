package cmd

import (
	"fmt"
	"io/ioutil"

	"github.com/Ladicle/hack/pkg/util"
	"github.com/atotto/clipboard"
	"github.com/spf13/cobra"
)

func NewCopyCmd() *cobra.Command {
	return &cobra.Command{
		Use:     "copy",
		Aliases: []string{"c"},
		Short:   "Copy main program to clipboard",
		RunE: func(cmd *cobra.Command, args []string) error {
			fname, err := util.GetProgName()
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
