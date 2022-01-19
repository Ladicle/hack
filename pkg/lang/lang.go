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
	"github.com/sergi/go-diff/diffmatchpatch"
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
	Input string
	Extra string
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

	want, err := ioutil.ReadFile(sample.Name(sampleID, sample.ExtSampleOut))
	if err != nil {
		return err
	}

	dmp := diffmatchpatch.New()
	diffs := dmp.DiffMain(out.String(), string(want), false)
	if len(diffs) == 1 && diffs[0].Type == diffmatchpatch.DiffEqual {
		return nil
	}

	input, err := ioutil.ReadFile(sampleInput.Name())
	if err != nil {
		return err
	}

	var buf bytes.Buffer
	var deco = color.New(color.FgYellow).SprintlnFunc()
	buf.WriteString(deco("Got:"))
	buf.Write(out.Bytes())
	buf.WriteString(deco("\nWant:"))
	buf.Write(want)
	if !color.NoColor {
		buf.WriteString(deco("\nDiff:"))
		buf.WriteString(dmp.DiffPrettyText(diffs))
	}
	buf.WriteString(deco("\nInput:"))
	buf.Write(input)

	return Error{
		ID:    sampleID,
		Type:  WrongAnswer,
		Extra: buf.String(),
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
