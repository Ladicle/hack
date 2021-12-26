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
)

const (
	dirPerm = 0755

	indentLv1 = " "
	indentLv2 = "   "

	example = `  # Initialize directories for AtCoder Beginner Contest 228.
  hack init abc228`
)

type Options struct {
	// ID is a identifier for the contest.
	ID string
}

func NewCommand(f *config.File, out io.Writer) *cobra.Command {
	var opts Options

	return &cobra.Command{
		Use:          "init <CONTEST>",
		Aliases:      []string{"ini", "i"},
		Short:        "create directories and download sample test cases for AtCoder.",
		Example:      example,
		SilenceUsage: true,
		RunE: func(cmd *cobra.Command, args []string) error {
			if err := opts.Validate(args); err != nil {
				return err
			}
			return opts.Run(f, out)
		},
	}
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

		fmt.Fprintf(out, "%v✓ Scraping task %v\n", indentLv1, task)
		for i, sample := range samples {
			id := i + 1 // convert to 1-index
			if samples == nil {
				continue
			}
			fmt.Fprintf(out, "%v✓ Scraping sample #%d\n", indentLv2, id)
			if err := sample.Write(dir, id); err != nil {
				return err
			}
		}
	}
	return nil
}
