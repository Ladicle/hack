package lang

import (
	"bytes"
	"context"
	"fmt"
	"os"
	"os/exec"
	"strings"
	"time"
)

const (
	TypeCpp = ".cpp"
	TypeGo  = ".go"

	ExeBinary = "solution"
)

// Tester is a interface for testing programs.
type Tester interface {
	Compile() error
	Run(sampleID string, timeout time.Duration) (string, error)
}

func GetTester(fileName string) Tester {
	switch {
	case strings.HasSuffix(fileName, TypeCpp):
		return &CppTester{ProgName: fileName}
	case strings.HasSuffix(fileName, TypeGo):
		return &GoTester{ProgName: fileName}
	}
	return nil
}

func runBinary(sampleID string, timeout time.Duration) (string, error) {
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()

	var errout bytes.Buffer
	c := exec.CommandContext(ctx, fmt.Sprintf("./%v", ExeBinary))
	c.Stderr = &errout

	inf, err := os.Open(fmt.Sprintf("%v.in", sampleID))
	if err != nil {
		return "", err
	}
	c.Stdin = inf

	outfName := attemptName(sampleID)
	outf, err := os.Create(outfName)
	if err != nil {
		return "", err
	}
	c.Stdout = outf

	if err := c.Run(); err != nil {
		if err == context.DeadlineExceeded {
			return outfName, err
		}
		return outfName, fmt.Errorf("%v - %v", err.Error(), errout.String())
	}
	return outfName, nil
}

func attemptName(id string) string {
	return fmt.Sprintf("attempt-sample%v-%v", id, time.Now().Format("20060102-150405"))
}
