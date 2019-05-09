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
	"github.com/logrusorgru/aurora"
)

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
		fmt.Printf("[%v] %v\n", aurora.Red("CE").Bold(), fname)
		return err
	}

	fs, err := ioutil.ReadDir(".")
	if err != nil {
		return err
	}
	for _, f := range fs {
		if f.Name() == "testing_tool.py" {
			cmd := exec.Command("python", "testing_tool.py", "./solution")
			out, err := cmd.CombinedOutput()
			if err != nil {
				fmt.Printf("[%v] testing_tool.py\n", aurora.Red("WA").Bold())
				fmt.Printf("\nDetail:\n%v", string(out))
				return err
			}
			fmt.Printf("[%v] testing_tool.py\n",
				aurora.Green("AC").Bold())
			return nil
		}
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
		// Run program
		outf, err := tester.Run(id, timeout)
		if err != nil {
			if err == context.DeadlineExceeded {
				fmt.Printf("[%v] Sample #%v\n",
					aurora.Yellow("TLE").Bold(), id)
				continue
			}
			fmt.Printf("[%v] Sample #%v\n",
				aurora.Red("RTE").Bold(), id)
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
			fmt.Printf("[%v] Sample #%v\n",
				aurora.Green("AC").Bold(), id)
			continue
		}
		fmt.Printf("[%v] Sample #%v\n", aurora.Red("WA").Bold(), id)
		was = append(was, TestResult{
			SampleID: id,
			Attempt:  outf,
			Got:      got,
			Want:     want,
		})
	}

	// Show detail
	for _, wa := range was {
		fmt.Printf("\nCompare %v.out and %v\n", wa.SampleID, wa.Attempt)
		showOutputDiff(wa.Got, wa.Want)
	}
	return nil
}

func showOutputDiff(got, want string) {
	gotL := strings.Split(got, "\n")
	wantL := strings.Split(want, "\n")
	gn := len(gotL)

	for i, w := range wantL {
		if i >= gn {
			prettyPrintDiff("<empty>", w)
			continue
		}
		if w == gotL[i] {
			fmt.Println(w)
		} else {
			prettyPrintDiff(gotL[i], w)
		}
	}
	if gn > len(wantL) {
		for i := len(wantL); i < gn; i++ {
			fmt.Printf("%v\n", aurora.Red(gotL[i]))
		}
	}
}

func prettyPrintDiff(got, want string) {
	fmt.Printf("%v\n%v\n", aurora.Red(got), aurora.Green(want))
}
