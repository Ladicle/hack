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
		path := filepath.Join(opt.WorkDir, fmt.Sprintf("%v.in", s))
		out, err := runGoFile(path)
		if err != nil {
			return err
		}

		path = filepath.Join(opt.WorkDir, fmt.Sprintf("%v.out", s))
		sampleOut, err := ioutil.ReadFile(path)
		if err != nil {
			fmt.Fprintln(c.IO, sampleOut)
			return err
		}

		got := strings.TrimSuffix(string(out), "\n")
		want := strings.TrimSuffix(string(sampleOut), "\n")
		if got == want {
			fmt.Fprintf(c.IO, "[%v] input #%v\n", aurora.Green("AC"), s)
		} else {
			fmt.Fprintf(c.IO, "[%v] input #%v\n", aurora.Red("WA"), s)
			gotL := strings.Split(got, "\n")
			wantL := strings.Split(want, "\n")
			for i, g := range gotL {
				if g == wantL[i] {
					fmt.Fprintln(c.IO, g)
				} else {
					fmt.Fprintf(c.IO, "%v\n%v\n", aurora.Red(g), aurora.Green(wantL[i]))
				}
			}
		}
	}
	return nil
}

func runGoFile(path string) ([]byte, error) {
	f, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	cmd := exec.Command("go", "run", "main.go")
	cmd.Stdin = f
	return cmd.Output()
}
