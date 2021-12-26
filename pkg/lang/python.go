package lang

import "context"

var _ Tester = PythonTester{}

type PythonTester struct {
	Options
}

func (t PythonTester) Run(sampleID int) error {
	ctx, cancel := context.WithTimeout(context.Background(), t.Timeout)
	defer cancel()
	return runProgram(ctx, sampleID, "python3", t.Program)
}
