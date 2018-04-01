package contest

import "testing"

func TestGetContest(t *testing.T) {
	testCases := []struct {
		name string
		want *Contest
		err  bool
	}{
		{
			name: "atcoder",
		},
		{
			name: "codejam",
		},
		{
			name: "hoge",
			err:  true,
		},
		{
			name: "",
			err:  true,
		},
	}

	LoadContest()

	for i, tc := range testCases {
		c, err := GetContest(tc.name)
		if err != nil {
			if tc.err {
				continue
			}
			t.Errorf("#%v: GetContest(%q): %v", i, tc.name, err)
			continue
		}
		if tc.err {
			t.Errorf("#%v: GetContest(%q) expect err", i, tc.name)
			continue
		}
		if got, want := c.Name, tc.name; got != want {
			t.Errorf("#%v: GetContest(%q) = %q, want %q", i, tc.name, got, want)
		}
	}
}
