package ini

import (
	"errors"
	"fmt"
	"io"
	"os"
	"path/filepath"

	"github.com/fatih/color"
	"github.com/spf13/cobra"

	"github.com/Ladicle/hack/pkg/config"
	"github.com/Ladicle/hack/pkg/contest"
	"github.com/Ladicle/hack/pkg/lang"
	"github.com/Ladicle/hack/pkg/sample"
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
	// Color is a flag of coloring.
	Color bool
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
			opts.Complete()
			return opts.Run(f, out)
		},
	}

	cmd.Flags().StringVarP(&opts.Lang, "lang", "l", "", "programming Language. (e.g. go)")
	cmd.Flags().BoolVarP(&opts.Color, "color", "C", false, "enable color output even if not in tty.")

	return cmd
}

func (o *Options) Validate(args []string) error {
	if len(args) != 1 {
		return errors.New("only accepts <CONTEST> value")
	}
	o.ID = args[0]
	return nil
}

func (o *Options) Complete() {
	if o.Color {
		color.NoColor = false
	}
}

func (o *Options) Run(f *config.File, out io.Writer) error {
	return o.initAtCoder(f, out)
}

func (o Options) initAtCoder(f *config.File, out io.Writer) error {
	at, err := contest.NewAtCoder(o.ID)
	if err != nil {
		return err
	}
	if err := at.Login(f.AtCoder.User, f.AtCoder.Pass); err != nil {
		return err
	}

	// Setup contest Directory
	contestDir := contest.GetContestDir(f.BaseDir, o.ID)
	if err := os.Mkdir(contestDir, dirPerm); err != nil && !os.IsExist(err) {
		return err
	}
	fmt.Fprintf(out, "Initialize directory for %v:\n", o.ID)

	tasks, err := at.ScrapeContest()
	if err != nil {
		return err
	}
	for _, task := range tasks {
		// Setup task directory
		taskDir := filepath.Join(contestDir, task)
		if err := os.Mkdir(taskDir, dirPerm); err != nil && !os.IsExist(err) {
			return err
		}
		if o.Lang != "" {
			prog := filepath.Join(taskDir, fmt.Sprintf("%v.%v", lang.ProgName, o.Lang))
			f, err := os.OpenFile(prog, os.O_RDONLY|os.O_CREATE, 0666)
			if err != nil {
				return err
			}
			f.Close()
		}

		// Setup sample directory
		sampleDir := filepath.Join(taskDir, sample.SampleDir)
		if err := os.Mkdir(sampleDir, dirPerm); err != nil && !os.IsExist(err) {
			return err
		}
		samples, err := at.ScrapeSamples(task)
		if err != nil {
			return err
		}
		green := color.New(color.FgGreen).SprintFunc()
		fmt.Fprintf(out, "%v%v Scraping task %v\n", indentLv1, green("✓"), task)
		for i, sample := range samples {
			id := i + 1 // convert to 1-index
			if samples == nil {
				continue
			}
			fmt.Fprintf(out, "%v%v Scraping sample #%d\n", indentLv2, green("✓"), id)
			if err := sample.Write(taskDir, id, filePerm); err != nil && !os.IsExist(err) {
				return err
			}
		}
	}
	return nil
}
