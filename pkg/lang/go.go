package lang

import (
	"fmt"
	"os/exec"
	"time"
)

type GoTester struct {
	ProgName string
}

func (t *GoTester) Compile() error {
	c := exec.Command("go", "build", "-o", ExeBinary, t.ProgName)
	if out, err := c.CombinedOutput(); err != nil {
		return fmt.Errorf("%v - %v", err.Error(), string(out))
	}
	return nil
}

func (t *GoTester) Run(sampleID string, timeout time.Duration) (string, error) {
	return runBinary(sampleID, timeout)
}
