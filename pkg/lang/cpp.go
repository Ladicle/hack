package lang

import (
	"context"
	"fmt"
	"os/exec"
)

var _ Tester = CppTester{}

type CppTester struct{ Options }

func (t CppTester) Run(sampleID int) error {
	cmd := exec.Command("g++", t.Program, "-std=c++14", "-pthread", "-o", BinaryName)
	if out, err := cmd.CombinedOutput(); err != nil {
		return Error{
			ID:    sampleID,
			Type:  RuntimeErr,
			Extra: fmt.Sprintf("%v: %s", err, out)}
	}
	ctx, cancel := context.WithTimeout(context.Background(), t.Timeout)
	defer cancel()
	return runProgram(ctx, sampleID, BinaryName)
}
