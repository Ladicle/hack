package cmd

import (
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/logrusorgru/aurora"
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
		inPath := filepath.Join(opt.WorkDir, fmt.Sprintf("%v.in", s))
		out, err := runGoFile(inPath, false)
		if err != nil {
			fmt.Fprintln(c.IO, string(out))
			return err
		}

		outPath := filepath.Join(opt.WorkDir, fmt.Sprintf("%v.out", s))
		sampleOut, err := ioutil.ReadFile(outPath)
		if err != nil {
			return err
		}

		got := strings.TrimSuffix(string(out), "\n")
		want := strings.TrimSuffix(string(sampleOut), "\n")
		if got == want {
			c.printResutl(true, s)
			continue
		}

		c.printResutl(false, s)
		c.showOutputDiff(got, want)

		if out, err := runGoFile(inPath, true); err != nil {
			return err
		} else {
			fmt.Fprintln(c.IO, "# debug")
			fmt.Fprintln(c.IO, string(out))
		}
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

func (c *testCmd) printResutl(result bool, id string) {
	state := aurora.Green("AC")
	if !result {
		state = aurora.Red("WA")
	}
	fmt.Fprintf(c.IO, "[%v] input #%v\n", state, id)
}

func runGoFile(path string, debug bool) ([]byte, error) {
	f, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	var cmd *exec.Cmd
	if debug {
		cmd = exec.Command("go", "run", "main.go", "--debug")
	} else {
		cmd = exec.Command("go", "run", "main.go")
	}
	cmd.Stdin = f
	return cmd.Output()
}
