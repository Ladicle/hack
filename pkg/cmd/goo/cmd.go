package goo

import (
	"fmt"
	"io"

	"github.com/Ladicle/hack/pkg/config"
	"github.com/Ladicle/hack/pkg/contest"
	"github.com/spf13/cobra"
)

type Options struct {
	contestID string
	taskID    string
}

func NewCommand(f *config.File, out io.Writer) *cobra.Command {
	return &cobra.Command{
		Use:          "go [<CONTEST>] [<TASK>]",
		Aliases:      []string{"g"},
		Short:        "Print path to the directory",
		SilenceUsage: true,
		RunE: func(cmd *cobra.Command, args []string) error {
			var opts Options
			if err := opts.Validate(args); err != nil {
				return err
			}
			return opts.Run(f, out)
		},
	}
}

func (o *Options) Validate(args []string) error {
	switch len(args) {
	case 0:
		// noop
	case 1:
		o.contestID = args[0]
	case 2:
		o.contestID = args[0]
		o.taskID = args[1]
	}
	return nil
}

func (o *Options) Run(f *config.File, out io.Writer) error {
	path, err := contest.GetDirFromID(f.BaseDir, o.contestID, o.taskID)
	if err != nil {
		return err
	}
	_, err = fmt.Fprint(out, path)
	return err
}
