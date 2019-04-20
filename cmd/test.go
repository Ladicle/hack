package cmd

import (
	"context"
	"fmt"
	"io/ioutil"
	"strings"
	"time"

	"github.com/Ladicle/hack/pkg/lang"
	"github.com/Ladicle/hack/pkg/util"
	"github.com/logrusorgru/aurora"
	"github.com/spf13/cobra"
)

func NewTestCmd() *cobra.Command {
	var timeout time.Duration

	c := cobra.Command{
		Use:     "test",
		Aliases: []string{"t"},
		Short:   "Test main program",
		RunE: func(cmd *cobra.Command, args []string) error {
			return runTest(timeout)
		},
		SilenceUsage: true,
	}

	c.Flags().DurationVarP(&timeout, "timeout", "t", 2*time.Second,
		"Set execution time-limit")
	return &c
}

type Result struct {
	Attempt  string
	SampleID string
	Want     string
	Got      string
}

func runTest(timeout time.Duration) error {
	fname, err := util.GetProgName()
	if err != nil {
		return err
	}

	// Compile the main program
	tester := lang.GetTester(fname)
	if err := tester.Compile(); err != nil {
		fmt.Printf("[%v] %v", aurora.Red("CE").Bold(), fname)
		return err
	}

	sids, err := util.SampleIDs(".")
	if err != nil {
		return err
	}

	var was []Result
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
		got, want := string(bgot), string(bwant)
		if got == want {
			fmt.Printf("[%v] Sample #%v\n",
				aurora.Green("AC").Bold(), id)
			continue
		}
		fmt.Printf("[%v] Sample #%v\n", aurora.Red("WA").Bold(), id)
		was = append(was, Result{
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
}

func prettyPrintDiff(got, want string) {
	fmt.Printf("%v\n%v\n", aurora.Red(got), aurora.Green(want))
}
