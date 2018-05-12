package cmd

import "testing"

func TestHandleCmd(t *testing.T) {
	addCmd(Command{
		Name: "test",
		Run: func(args []string, opt Option) error {
			return nil
		},
	})

	cases := []struct {
		Name   string
		ExcErr bool
	}{
		{Name: "test"},
		{Name: "unknown", ExcErr: true},
	}

	for _, tc := range cases {
		if err := HandleCmd(tc.Name, []string{}, Option{}); err != nil && !tc.ExcErr {
			t.Fatalf("Could not handle the %v command.: %v", tc.Name, err)
		}
	}
}
