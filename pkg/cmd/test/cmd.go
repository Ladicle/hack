package test

import (
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"path/filepath"
	"time"

	"github.com/atotto/clipboard"
	"github.com/pkg/browser"
	"github.com/spf13/cobra"

	"github.com/Ladicle/hack/pkg/config"
	"github.com/Ladicle/hack/pkg/contest"
	"github.com/Ladicle/hack/pkg/lang"
	"github.com/Ladicle/hack/pkg/sample"
)

type Options struct {
	Timeout time.Duration
	Copy    bool
	Open    bool
}

func NewCommand(f *config.File, out io.Writer) *cobra.Command {
	var opts Options

	cmd := &cobra.Command{
		Use:          "test",
		Aliases:      []string{"tt", "t"},
		Short:        "Test your program",
		SilenceUsage: true,
		RunE: func(cmd *cobra.Command, args []string) error {
			return opts.Run(f, out)
		},
	}

	cmd.Flags().DurationVar(&opts.Timeout, "timeout", 2*time.Second, "set timeout duration.")
	cmd.Flags().BoolVar(&opts.Copy, "copy", true, "copy program to clipboard after passing all tests.")
	cmd.Flags().BoolVar(&opts.Open, "open", true, "open task page after passing all tests.")

	return cmd
}

func (o *Options) Run(f *config.File, out io.Writer) error {
	wd, err := os.Getwd()
	if err != nil {
		return err
	}

	prog, err := lang.FindProg(wd)
	if err != nil {
		return err
	}
	tester, err := lang.GetTester(prog, o.Timeout)
	if err != nil {
		return err
	}
	num, err := sample.CntInputs(wd)
	if err != nil {
		return err
	}

	var cntErr int
	for sampleID := 1; sampleID <= num; sampleID++ {
		err = tester.Run(sampleID)
		if err == nil {
			fmt.Fprintf(out, "[AC] Sample #%d\n", sampleID)
			continue
		}
		var langErr lang.Error
		if ok := errors.As(err, &langErr); !ok {
			return err
		}
		fmt.Fprintln(out, langErr)
		cntErr++
	}
	if cntErr > 0 {
		return fmt.Errorf("Fail %d samples", cntErr)
	}

	if o.Copy {
		data, err := ioutil.ReadFile(prog)
		if err != nil {
			return err
		}
		if err := clipboard.WriteAll(string(data)); err != nil {
			return err
		}
		fmt.Fprintf(out, "Copy %v!\n", prog)
	}

	if o.Open {
		taskURL := contest.GetTaskURL(getContestID(wd), getTaskID(wd))
		if err := browser.OpenURL(taskURL); err != nil {
			return err
		}
	}
	return nil
}

// getContestID returns the parent directory name as the contest ID.
func getContestID(dir string) string {
	curBase := filepath.Base(dir)
	return filepath.Base(dir[:len(dir)-len(curBase)])
}

// getTaskID returns the specified directory name as the task ID.
func getTaskID(dir string) string {
	return filepath.Base(dir)
}
