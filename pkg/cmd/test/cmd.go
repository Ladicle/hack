package test

import (
	"bytes"
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
	"github.com/sergi/go-diff/diffmatchpatch"
	"github.com/spf13/cobra"

	"github.com/Ladicle/hack/pkg/config"
	"github.com/Ladicle/hack/pkg/contest"
	"github.com/Ladicle/hack/pkg/lang"
	"github.com/Ladicle/hack/pkg/sample"
)

type Options struct {
	SampleID int

	Timeout time.Duration
	Submit  bool
	Copy    bool
	Open    bool
	Color   bool

	HideInput  bool
	HideOutput bool
	HideExpect bool
	HideDiff   bool

	workingDir string
	program    string
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

	cmd.Flags().DurationVarP(&opts.Timeout, "timeout", "t", 2*time.Second, "set timeout duration.")
	cmd.Flags().BoolVar(&opts.Copy, "copy", false, "copy program to clipboard after passing all tests.")
	cmd.Flags().BoolVar(&opts.Submit, "submit", true, "submit program after passing all tests.")
	cmd.Flags().BoolVar(&opts.Open, "open", true, "open task page after passing all tests.")
	cmd.Flags().BoolVarP(&opts.Color, "color", "C", false, "enable color output even if not in tty.")

	cmd.Flags().BoolVarP(&opts.HideInput, "hide-input", "I", false, "do not print input value when the test fails.")
	cmd.Flags().BoolVarP(&opts.HideOutput, "hide-output", "O", false, "do not print output value when the test fails.")
	cmd.Flags().BoolVarP(&opts.HideExpect, "hide-expect", "E", false, "do not print expected value when the test fails.")
	cmd.Flags().BoolVarP(&opts.HideDiff, "hide-diff", "D", false, "do not print diff value when the test fails.")

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
	o.program = prog

	tester, err := lang.GetTester(prog, o.Timeout)
	if err != nil {
		return err
	}
	o.tester = tester
	return nil
}

func (o Options) Run(f *config.File, out io.Writer) error {
	// Test only the specified sample
	if o.SampleID > 0 {
		if _, err := o.testProgram(o.SampleID, out); err != nil {
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
		if pass, err := o.testProgram(sampleID, out); err != nil {
			return err
		} else if !pass {
			cntErr++
		}
	}
	if cntErr > 0 {
		return fmt.Errorf("Fail %v, Pass %v samples",
			color.RedString("%d", cntErr), color.GreenString("%d", num-cntErr))
	}

	var (
		contestID = contest.GetContestID(o.workingDir)
		taskID    = contest.GetTaskID(o.workingDir)

		openURL string
	)

	if o.Copy {
		data, err := ioutil.ReadFile(o.program)
		if err != nil {
			return err
		}
		if err := clipboard.WriteAll(string(data)); err != nil {
			return err
		}
		openURL = contest.GetTaskURL(contestID, taskID)
		fmt.Fprintf(out, "Copy %v!\n", o.program)
	}

	if o.Submit {
		at, err := contest.NewAtCoder(contestID)
		if err != nil {
			return err
		}
		if err := at.Login(f.AtCoder.User, f.AtCoder.Pass); err != nil {
			return err
		}
		if err := at.SubmitCode(taskID, o.program); err != nil {
			return err
		}
		openURL = contest.GetSubmitMeURL(contestID)
		fmt.Fprintf(out, "Submit %v!\n", o.program)
	}

	if o.Open {
		if err := browser.OpenURL(openURL); err != nil {
			return err
		}
	}
	return nil
}

// testProgram tests program with the specified sample ID, then returns pass flag of the
// sample test case and error.
func (o Options) testProgram(sampleID int, out io.Writer) (pass bool, err error) {
	err = o.tester.Run(sampleID)
	if err == nil {
		fmt.Fprintf(out, "[%v] Sample #%d\n", color.GreenString("AC"), sampleID)
		return true, nil
	}
	var langErr lang.Error
	if ok := errors.As(err, &langErr); !ok {
		return false, err
	}
	fmt.Fprintln(out, langErr)
	if langErr.Type != lang.WrongAnswer {
		return false, nil
	}

	// Print debug information
	var input []byte
	if !o.HideInput {
		data, err := ioutil.ReadFile(langErr.SampleInFile)
		if err != nil {
			return false, err
		}
		input = data
	}

	var diff string
	if !o.HideDiff {
		dmp := diffmatchpatch.New()
		diff = dmp.DiffPrettyText(dmp.DiffMain(string(langErr.Output), string(langErr.Expect), false))
	}

	var buf bytes.Buffer
	var deco = color.New(color.FgYellow).SprintlnFunc()
	if !o.HideOutput {
		buf.WriteString(deco("Got:"))
		buf.Write(langErr.Output)
	}
	if !o.HideExpect {
		buf.WriteString(deco("\nExpect:"))
		buf.Write(langErr.Expect)
	}
	if !color.NoColor && !o.HideDiff {
		buf.WriteString(deco("\nDiff:"))
		buf.WriteString(diff)
	}
	if !o.HideInput {
		buf.WriteString(deco("\nInput:"))
		buf.Write(input)
	}
	fmt.Fprintln(out, buf.String())
	return false, nil
}
