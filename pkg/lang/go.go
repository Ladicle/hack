package lang

import (
	"context"
	"fmt"
	"os/exec"
)

var _ Tester = GoTester{}

type GoTester struct{ Options }

func (t GoTester) Run(sampleID int) error {
	cmd := exec.Command("go", "build", "-o", defaultBinaryName, t.Program)
	if out, err := cmd.CombinedOutput(); err != nil {
		return Error{
			ID:    sampleID,
			Type:  RuntimeErr,
			Extra: fmt.Sprintf("%v: %s", err, out)}
	}
	ctx, cancel := context.WithTimeout(context.Background(), t.Timeout)
	defer cancel()
	return runProgram(ctx, sampleID, defaultBinaryName)
}
