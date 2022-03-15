package goo

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"

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
	path := contest.GetAtCoderDir(f.BaseDir)
	if o.contestID != "" {
		path = filepath.Join(path, o.contestID)
	}
	if o.taskID != "" {
		if _, err := os.Stat(filepath.Join(path, o.taskID)); os.IsExist(err) {
			path = filepath.Join(path, o.taskID)
		} else {
			// Complement task ID from suffix. (e.g. c -> foo_contest_c)
			dirs, err := os.ReadDir(path)
			if err != nil {
				return err
			}
			for _, dir := range dirs {
				if !dir.IsDir() {
					continue
				}
				if strings.HasSuffix(dir.Name(), o.taskID) {
					path = filepath.Join(path, dir.Name())
					break
				}
			}
		}
	}
	_, err := fmt.Fprint(out, path)
	return err
}
