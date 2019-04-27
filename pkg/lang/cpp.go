package lang

import (
	"fmt"
	"os/exec"
	"time"
)

type CppTester struct {
	ProgName string
}

func (t *CppTester) Compile() error {
	c := exec.Command("g++", t.ProgName, "-std=c++14",
		"-pthread", "-o", binaryName)
	if out, err := c.CombinedOutput(); err != nil {
		return fmt.Errorf("%v - %v", err.Error(), string(out))
	}
	return nil
}

func (t *CppTester) Run(sampleID string, timeout time.Duration) (string, error) {
	return runBinary(sampleID, timeout)
}
