package cmd

import (
	"fmt"
	"io/ioutil"
	"strings"

	"github.com/atotto/clipboard"
	"github.com/spf13/cobra"
)

func NewCopyCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "copy",
		Short: "Copy main program to clipboard",
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

func getProgName() (string, error) {
	var fname string

	fs, err := ioutil.ReadDir(".")
	if err != nil {
		return fname, err
	}

	for _, f := range fs {
		if !f.IsDir() && strings.HasPrefix(f.Name(), "main.") {
			return f.Name(), nil
		}
	}

	return fname, fmt.Errorf("not found a main program")
}
