package submit

import (
	"context"
	"errors"
	"fmt"
	"io/ioutil"
	"os/exec"
	"strings"
	"time"

	"github.com/Ladicle/hack/pkg/format"
	"github.com/Ladicle/hack/pkg/lang"
	"github.com/Ladicle/hack/pkg/util"
	"github.com/golang/glog"
)

const testTool = "testing_tool.py"

type TestResult struct {
	Attempt  string
	SampleID string
	Want     string
	Got      string
}

func runTest(timeout time.Duration, hr *format.HackRobot) error {
	fname, err := util.GetProgName()
	if err != nil {
		return err
	}
	if fname == "" {
		hr.Fatal("Not found the %v program in this directory.", "main.[go|cpp]")
	}

	// Compile the main program
	tester := lang.GetTester(fname)
	if err := tester.Compile(); err != nil {
		hr.State(format.StateCompileError, fname)
		return err
	}

	fs, err := ioutil.ReadDir(".")
	if err != nil {
		return err
	}
	for _, f := range fs {
		if f.Name() != testTool {
			continue
		}

		cmd := exec.Command("python", "testing_tool.py", "./solution")
		out, err := cmd.CombinedOutput()
		if err != nil {
			hr.State(format.StateWrongAnswer, testTool)
			glog.Error(string(out))
			return err
		}
		hr.State(format.StateAnswerIsCorrect, testTool)
		return nil
	}

	sids, err := util.SampleIDs(".")
	if err != nil {
		return err
	}
	if len(sids) == 0 {
		return errors.New("There is no sample inputs.")
	}

	var was []TestResult
	for _, id := range sids {
		sampleName := fmt.Sprintf("Sample #%v", id)

		// Run program
		outf, err := tester.Run(id, timeout)
		if err != nil {
			if err == context.DeadlineExceeded {
				hr.State(format.StateTimeLimitExceeded, sampleName)
				continue
			}
			hr.State(format.StateAnswerIsCorrect, sampleName)
			return err
		}

		// Check the answer
		bgot, err := ioutil.ReadFile(outf)
		if err != nil {
			return err
		}

		bwant, err := ioutil.ReadFile(fmt.Sprintf("%v.out", id))
		got, want := strings.TrimSuffix(string(bgot), "\n"),
			strings.TrimSuffix(string(bwant), "\n")

		if got == want {
			hr.State(format.StateAnswerIsCorrect, sampleName)
			continue
		}
		hr.State(format.StateWrongAnswer, sampleName)

		was = append(was, TestResult{
			SampleID: id,
			Attempt:  outf,
			Got:      got,
			Want:     want,
		})
	}

	// Show detail
	for _, wa := range was {
		hr.Printfln("\nCompare %v.out and %v", wa.SampleID, wa.Attempt)
		hr.PrettyDiff(wa.Got, wa.Want)
	}
	return nil
}
