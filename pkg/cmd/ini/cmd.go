package ini

import (
	"errors"
	"fmt"
	"io"
	"os"
	"path/filepath"

	"github.com/spf13/cobra"

	"github.com/Ladicle/hack/pkg/config"
	"github.com/Ladicle/hack/pkg/contest"
	"github.com/Ladicle/hack/pkg/lang"
)

const (
	dirPerm  = 0755
	filePerm = 0644

	indentLv1 = " "
	indentLv2 = "   "

	example = `  # Initialize directories for AtCoder Beginner Contest 228.
  hack init abc228`
)

type Options struct {
	// ID is a identifier for the contest.
	ID string
	// Lang is a name of programming language.
	Lang string
}

func NewCommand(f *config.File, out io.Writer) *cobra.Command {
	var opts Options

	cmd := &cobra.Command{
		Use:          "init <CONTEST>",
		Aliases:      []string{"ini", "i"},
		Short:        "Create directories and download samples",
		Example:      example,
		SilenceUsage: true,
		RunE: func(cmd *cobra.Command, args []string) error {
			if err := opts.Validate(args); err != nil {
				return err
			}
			return opts.Run(f, out)
		},
	}

	cmd.Flags().StringVarP(&opts.Lang, "lang", "l", "", "programming Language. (e.g. go)")

	return cmd
}

func (o *Options) Validate(args []string) error {
	if len(args) != 1 {
		return errors.New("only accepts <CONTEST> value")
	}
	o.ID = args[0]
	return nil
}

func (o *Options) Run(f *config.File, out io.Writer) error {
	at, err := contest.NewAtCoder(o.ID)
	if err != nil {
		return err
	}
	if err := at.Login(f.AtCoder.User, f.AtCoder.Pass); err != nil {
		return err
	}

	tasks, err := at.ScrapeTasks()
	if err != nil {
		return err
	}

	if err := os.Mkdir(o.ID, dirPerm); err != nil {
		return err
	}
	fmt.Fprintf(out, "Initialize directory for %v:\n", o.ID)

	for _, task := range tasks {
		samples, err := at.ScrapeTask(task)
		if err != nil {
			return err
		}

		dir := filepath.Join(o.ID, task)
		if err := os.Mkdir(dir, dirPerm); err != nil {
			return err
		}

		if o.Lang != "" {
			prog := filepath.Join(dir, fmt.Sprintf("%v.%v", lang.ProgName, o.Lang))
			f, err := os.Create(prog)
			if err != nil {
				return err
			}
			f.Close()
		}

		fmt.Fprintf(out, "%v✓ Scraping task %v\n", indentLv1, task)
		for i, sample := range samples {
			id := i + 1 // convert to 1-index
			if samples == nil {
				continue
			}
			fmt.Fprintf(out, "%v✓ Scraping sample #%d\n", indentLv2, id)
			if err := sample.Write(dir, id, filePerm); err != nil {
				return err
			}
		}
	}
	return nil
}
