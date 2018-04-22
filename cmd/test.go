package cmd

import (
	"context"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"time"

	"github.com/logrusorgru/aurora"
)

const (
	defaultTimeout = 2 * time.Second
	timeoutErrMsg  = "signal: killed"

	statusWA = "WA"
	statusLT = "LT"
	statusAC = "AC"
)

// NewTestCmd tests contest testrmation.
func NewTestCmd(io io.Writer) Command {
	s := testCmd{IO: io}
	return Command{
		Name:        "test",
		Short:       "test [NUMBER]",
		Description: "Test tests your code with all samples if you don't specified the number",
		Run:         s.run,
	}
}

type testCmd struct {
	IO io.Writer
}

func (c *testCmd) run(args []string, opt Option) error {
	if len(args) > 1 {
		return c.test(opt.WorkDir, args[0])
	}

	fs, err := ioutil.ReadDir(opt.WorkDir)
	if err != nil {
		return err
	}

	var samples []string
	for _, f := range fs {
		n := strings.Split(f.Name(), ".")
		if len(n) < 2 {
			continue
		}
		if n[1] == "out" {
			samples = append(samples, n[0])
		}
	}

	for _, s := range samples {
		if err := c.test(opt.WorkDir, s); err != nil {
			return err
		}
	}
	return nil
}

func (c *testCmd) test(workDir, s string) error {
	inPath := filepath.Join(workDir, fmt.Sprintf("%v.in", s))
	out, err := runGoFile(inPath, false, defaultTimeout)
	if err != nil {
		if err.Error() == timeoutErrMsg {
			c.printResutl(statusLT, s)
			return nil
		}
		fmt.Fprintln(c.IO, string(out))
		return err
	}

	outPath := filepath.Join(workDir, fmt.Sprintf("%v.out", s))
	sampleOut, err := ioutil.ReadFile(outPath)
	if err != nil {
		return err
	}

	got := strings.TrimSuffix(string(out), "\n")
	want := strings.TrimSuffix(string(sampleOut), "\n")
	if got == want {
		c.printResutl(statusAC, s)
		return nil
	}

	c.printResutl(statusWA, s)
	c.showOutputDiff(got, want)

	if out, err := runGoFile(inPath, true, defaultTimeout); err != nil {
		if err.Error() == timeoutErrMsg {
			c.printResutl(statusLT, s)
			return nil
		}
		return err
	} else {
		fmt.Fprintln(c.IO, "# debug")
		fmt.Fprintln(c.IO, string(out))
	}
	return nil
}

func (c *testCmd) showOutputDiff(got, want string) {
	gotL := strings.Split(got, "\n")
	wantL := strings.Split(want, "\n")
	gn := len(gotL)

	for i, w := range wantL {
		if i >= gn {
			c.prettyPrintDiff("<empty>", w)
			continue
		}
		if w == gotL[i] {
			fmt.Fprintln(c.IO, w)
		} else {
			c.prettyPrintDiff(gotL[i], w)
		}
	}
}

func (c *testCmd) prettyPrintDiff(got, want string) {
	fmt.Fprintf(c.IO, "%v\n%v\n", aurora.Red(got), aurora.Green(want))
}

func (c *testCmd) printResutl(status, id string) {
	var state aurora.Value
	switch status {
	case statusAC:
		state = aurora.Green(statusAC)
	case statusWA:
		state = aurora.Red(statusWA)
	case statusLT:
		state = aurora.Brown(statusLT)
	}
	fmt.Fprintf(c.IO, "[%v] input #%v\n", state, id)
}

func runGoFile(path string, debug bool, timeout time.Duration) ([]byte, error) {
	f, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()

	var cmd *exec.Cmd
	if debug {
		cmd = exec.CommandContext(ctx, "go", "run", "main.go", "--debug")
	} else {
		cmd = exec.CommandContext(ctx, "go", "run", "main.go")
	}
	cmd.Stdin = f
	return cmd.CombinedOutput()
}