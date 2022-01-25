package lang

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"time"

	"github.com/Ladicle/hack/pkg/sample"
	"github.com/fatih/color"
)

const (
	ProgName   = "main"
	BinaryName = "./main"

	LangCpp    = "cpp"
	LangGo     = "go"
	LangPython = "py"
)

// Tester is a interface for testing programs.
type Tester interface {
	// Run executes program which is passed sample input with the specified ID.
	Run(sampleID int) error
}

type Options struct {
	Program string
	Timeout time.Duration
}

type ErrorType string

const (
	RuntimeErr  ErrorType = "RE"
	TimeoutErr  ErrorType = "TLE"
	WrongAnswer ErrorType = "WA"
)

type Error struct {
	ID    int
	Type  ErrorType
	Extra string

	Output       []byte
	Expect       []byte
	SampleInFile string
}

func (e Error) Error() string {
	var (
		buf     bytes.Buffer
		errType string
	)
	switch e.Type {
	case RuntimeErr, WrongAnswer:
		errType = color.RedString("%v", e.Type)
	case TimeoutErr:
		errType = color.YellowString("%v", e.Type)
	}
	buf.WriteString(fmt.Sprintf("[%s] Sample #%d", errType, e.ID))
	if e.Extra != "" {
		buf.WriteString("\n")
		buf.WriteString(e.Extra)
	}
	return buf.String()
}

func GetTester(program string, timeout time.Duration) (Tester, error) {
	opts := Options{
		Program: program,
		Timeout: timeout,
	}
	lang := filepath.Ext(program)[1:]
	switch lang {
	case LangCpp:
		return &CppTester{Options: opts}, nil
	case LangGo:
		return &GoTester{Options: opts}, nil
	case LangPython:
		return &PythonTester{Options: opts}, nil
	}
	return nil, fmt.Errorf("%q is unsupported program", lang)
}

func runProgram(ctx context.Context, sampleID int, args ...string) error {
	var (
		out    bytes.Buffer
		errout bytes.Buffer
	)
	c := exec.CommandContext(ctx, args[0], args[1:]...)
	c.Stdout = &out
	c.Stderr = &errout

	// Pass sample input file to Standard Input
	sampleInput, err := os.Open(sample.Name(sampleID, sample.ExtSampleIn))
	if err != nil {
		return err
	}
	defer sampleInput.Close()
	c.Stdin = sampleInput

	start := time.Now()
	if err := c.Run(); err != nil {
		// TLE
		if err == context.DeadlineExceeded || err.Error() == "signal: killed" {
			return Error{
				ID:   sampleID,
				Type: TimeoutErr,
				Extra: fmt.Sprintf(
					"Past %v:\nStdout:\n%s\nIf you want to change timeout, use `--timeout(-t)` flag.",
					time.Now().Sub(start), out.String()),
			}
		}
		// RE
		if len(errout.Bytes()) > 0 {
			return Error{
				ID:   sampleID,
				Type: RuntimeErr,
				Extra: fmt.Sprintf("%v:\nStdout:\n%sStderr:\n%s",
					err, out.String(), errout.String()),
			}
		}
		return err
	}

	expect, err := ioutil.ReadFile(sample.Name(sampleID, sample.ExtSampleOut))
	if err != nil {
		return err
	}
	// AC
	if bytes.Compare(out.Bytes(), expect) == 0 {
		return nil
	}
	// WA
	return Error{
		ID:           sampleID,
		Type:         WrongAnswer,
		Output:       out.Bytes(),
		Expect:       expect,
		SampleInFile: sampleInput.Name(),
	}
}

// FindProg finds a file that has 'main' as a name prefix in the current directory.
func FindProg(dir string) (string, error) {
	entries, err := os.ReadDir(dir)
	if err != nil {
		return "", err
	}
	for _, entry := range entries {
		if !entry.IsDir() && strings.HasPrefix(filepath.Base(entry.Name()), ProgName) {
			return entry.Name(), nil
		}
	}
	return "", errors.New("not found 'main' program in the current directory")
}
