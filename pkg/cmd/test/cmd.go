package test

import (
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"strconv"
	"time"

	"github.com/atotto/clipboard"
	"github.com/fatih/color"
	"github.com/pkg/browser"
	"github.com/spf13/cobra"

	"github.com/Ladicle/hack/pkg/config"
	"github.com/Ladicle/hack/pkg/contest"
	"github.com/Ladicle/hack/pkg/lang"
	"github.com/Ladicle/hack/pkg/sample"
)

type Options struct {
	SampleID int

	Timeout time.Duration
	Copy    bool
	Open    bool
	Color   bool

	workingDir string
	programID  string
	tester     lang.Tester
}

func NewCommand(f *config.File, out io.Writer) *cobra.Command {
	var opts Options

	cmd := &cobra.Command{
		Use:          "test [<SAMPLE>]",
		Aliases:      []string{"tt", "t"},
		Short:        "Test your program",
		SilenceUsage: true,
		RunE: func(cmd *cobra.Command, args []string) error {
			if err := opts.Validate(args); err != nil {
				return err
			}
			if err := opts.Complete(); err != nil {
				return err
			}
			return opts.Run(f, out)
		},
	}

	cmd.Flags().DurationVar(&opts.Timeout, "timeout", 2*time.Second, "set timeout duration.")
	cmd.Flags().BoolVar(&opts.Copy, "copy", true, "copy program to clipboard after passing all tests.")
	cmd.Flags().BoolVar(&opts.Open, "open", true, "open task page after passing all tests.")
	cmd.Flags().BoolVarP(&opts.Color, "color", "C", false, "enable color output even if not in tty.")

	return cmd
}

func (o *Options) Validate(args []string) error {
	if len(args) == 1 {
		id, err := strconv.Atoi(args[0])
		if err != nil {
			return err
		}
		o.SampleID = id
	}
	return nil
}

func (o *Options) Complete() error {
	if o.Color {
		color.NoColor = false
	}

	wd, err := os.Getwd()
	if err != nil {
		return err
	}
	o.workingDir = wd

	prog, err := lang.FindProg(wd)
	if err != nil {
		return err
	}
	o.programID = prog

	tester, err := lang.GetTester(prog, o.Timeout)
	if err != nil {
		return err
	}
	o.tester = tester
	return nil
}

func (o *Options) Run(f *config.File, out io.Writer) error {
	// Test only the specified sample
	if o.SampleID != 0 {
		if _, err := testProgram(o.tester, o.SampleID, out); err != nil {
			return err
		}
		return nil
	}

	// Test all samples
	num, err := sample.CntInputs(o.workingDir)
	if err != nil {
		return err
	}

	var cntErr int
	for sampleID := 1; sampleID <= num; sampleID++ {
		if pass, err := testProgram(o.tester, o.SampleID, out); err != nil {
			return err
		} else if !pass {
			cntErr++
		}
	}
	if cntErr > 0 {
		return fmt.Errorf("Fail %v, Pass %v samples",
			color.RedString("%d", cntErr), color.GreenString("%d", num-cntErr))
	}

	if o.Copy {
		data, err := ioutil.ReadFile(o.programID)
		if err != nil {
			return err
		}
		if err := clipboard.WriteAll(string(data)); err != nil {
			return err
		}
		fmt.Fprintf(out, "Copy %v!\n", o.programID)
	}

	if o.Open {
		var (
			contestID = contest.GetContestID(o.workingDir)
			taskID    = contest.GetTaskID(o.workingDir)
		)
		taskURL := contest.GetTaskURL(contestID, taskID)
		if err := browser.OpenURL(taskURL); err != nil {
			return err
		}
	}
	return nil
}

// testProgram tests program with the specified sample ID, then returns pass flag of the
// sample test case and error.
func testProgram(tester lang.Tester, sampleID int, out io.Writer) (pass bool, err error) {
	err = tester.Run(sampleID)
	if err == nil {
		fmt.Fprintf(out, "[%v] Sample #%d\n", color.GreenString("AC"), sampleID)
		return true, nil
	}
	var langErr lang.Error
	if ok := errors.As(err, &langErr); !ok {
		return false, err
	}
	fmt.Fprintln(out, langErr)
	return false, nil
}
