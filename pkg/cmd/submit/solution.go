package submit

import (
	"context"
	"fmt"
	"io/ioutil"
	"os/exec"
	"time"

	"github.com/Ladicle/hack/pkg/config"
	"github.com/Ladicle/hack/pkg/contest"
	"github.com/Ladicle/hack/pkg/format"
	"github.com/Ladicle/hack/pkg/lang"
	"github.com/Ladicle/hack/pkg/util"
	"github.com/atotto/clipboard"
	"github.com/golang/glog"
	"github.com/skratchdot/open-golang/open"
)

const gcjTestTool = "testing_tool.py"

type TestResult struct {
	Attempt  string
	SampleID string
	Want     string
	Got      string
}

func NewTestResult(attempt, sampleID, want, got string) *TestResult {
	return &TestResult{
		Attempt:  attempt,
		SampleID: sampleID,
		Want:     want,
		Got:      got,
	}
}

type Solution struct {
	QuizID   string
	ProgFile string

	*format.HackRobot
}

func NewSolution(quizID, progFile string, hr *format.HackRobot) *Solution {
	return &Solution{
		QuizID:   quizID,
		ProgFile: progFile,

		HackRobot: hr,
	}
}

func (s *Solution) Copy() error {
	code, err := ioutil.ReadFile(s.ProgFile)
	if err != nil {
		return err
	}

	if err := clipboard.WriteAll(string(code)); err != nil {
		return err
	}
	s.Info("Copy %q to the clipboard.", s.ProgFile)
	return nil
}

func (s *Solution) Open() error {
	at, err := contest.NewAtCoder(config.CurrentContestID())
	if err != nil {
		return err
	}
	return open.Run(at.SubmittedURL())
}

func (s *Solution) Submit() error {
	at, err := contest.NewAtCoder(config.CurrentContestID())
	if err != nil {
		return err
	}
	if err := at.SubmitCode(s.QuizID, s.ProgFile); err != nil {
		return err
	}
	s.Info("Success to submit the %v code!", s.QuizID)
	return nil
}

func (s *Solution) Test(limit time.Duration) error {
	// compile source code
	tester := lang.GetTester(s.ProgFile)
	if err := tester.Compile(); err != nil {
		s.State(format.StateCompileError, s.ProgFile)
		return err
	}

	// run interactive test tool if it exists
	if exist, err := util.InDir(gcjTestTool); err != nil {
		return err
	} else if exist {
		cmd := exec.Command("python", gcjTestTool, lang.ExeBinary)
		out, err := cmd.CombinedOutput()
		if err != nil {
			s.State(format.StateWrongAnswer, gcjTestTool)
			glog.Error(string(out))
			return err
		}
		s.State(format.StateAnswerIsCorrect, gcjTestTool)
	}

	// check source code using sample code
	sids, err := util.SampleIDs(".")
	if err != nil {
		return err
	}
	if len(sids) == 0 {
		s.Fatal("There is no samples.")
	}
	var was []*TestResult
	for _, id := range sids {
		name := fmt.Sprintf("Sample #%v", id)

		// Run program
		outf, err := tester.Run(id, limit)
		if err != nil {
			if err == context.DeadlineExceeded {
				s.State(format.StateTimeLimitExceeded, name)
				continue
			}
			s.State(format.StateRuntimeError, name)
			return err
		}

		// Check the answer
		got, err := util.CleanRead(outf)
		if err != nil {
			return err
		}
		want, err := util.CleanRead(fmt.Sprintf("%v.out", id))
		if err != nil {
			return err
		}
		if got == want {
			s.State(format.StateAnswerIsCorrect, name)
			continue
		}
		s.State(format.StateWrongAnswer, name)
		was = append(was, NewTestResult(id, outf, got, want))
	}

	// Show error detail
	for _, wa := range was {
		s.Printfln("\nCompare %v.out and %v", wa.SampleID, wa.Attempt)
		s.PrettyDiff(wa.Got, wa.Want)
	}
	if len(was) > 0 {
		return fmt.Errorf("test is failed")
	}
	return nil
}
