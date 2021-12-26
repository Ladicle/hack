package lang

import (
	"bytes"
	"context"
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"path/filepath"
	"time"

	"github.com/sergi/go-diff/diffmatchpatch"
)

const (
	defaultBinaryName = "./main"

	LangCpp    = ".cpp"
	LangGo     = ".go"
	LangPython = ".py"
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
}

func (e Error) Error() string {
	var buf bytes.Buffer
	buf.WriteString(fmt.Sprintf("[%v] Sample #%d", e.Type, e.ID))
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
	ext := filepath.Ext(program)
	switch ext {
	case LangCpp:
		return &CppTester{Options: opts}, nil
	case LangGo:
		return &GoTester{Options: opts}, nil
	case LangPython:
		return &PythonTester{Options: opts}, nil
	}
	return nil, fmt.Errorf("%q is unsupported program", ext)
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
	sampleInput, err := os.Open(fmt.Sprintf("%d.in", sampleID))
	if err != nil {
		return err
	}
	c.Stdin = sampleInput

	if err := c.Run(); err != nil {
		if err == context.DeadlineExceeded {
			return Error{ID: sampleID, Type: TimeoutErr}
		}
		if len(errout.Bytes()) > 0 {
			return Error{
				ID:    sampleID,
				Type:  RuntimeErr,
				Extra: fmt.Sprintf("%v: %s", err, errout.String()),
			}
		}
		return err
	}

	want, err := ioutil.ReadFile(fmt.Sprintf("%d.out", sampleID))
	if err != nil {
		return err
	}
	dmp := diffmatchpatch.New()
	diffs := dmp.DiffMain(out.String(), string(want), false)
	if len(diffs) == 1 && diffs[0].Type == diffmatchpatch.DiffEqual {
		return nil
	}
	return Error{
		ID:    sampleID,
		Type:  WrongAnswer,
		Extra: dmp.DiffPrettyText(diffs),
	}
}
