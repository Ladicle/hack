package contest

import "fmt"

// Contest manages programming contest.
type Contest struct {
	Name  string
	Set   func(output string, arg []string) error
	Usage string
}

var contests []Contest

// LoadContest loads all contests.
func LoadContest() {
	addContest(NewAtCorderContest())
	addContest(NewCodeJamContest())
	addContest(NewFreeContest())
}

func addContest(c Contest) {
	contests = append(contests, c)
}

// GetContest returns matched contest instance
func GetContest(name string) (*Contest, error) {
	if name == "" {
		return nil, fmt.Errorf("contest name is a required argument")
	}
	for _, c := range contests {
		if c.Name == name {
			return &c, nil
		}
	}
	return nil, fmt.Errorf("%s is an invalid contest name", name)
}

// ListContestName returns contest name list.
func ListContestName() []string {
	var list []string
	for _, c := range contests {
		list = append(list, c.Name)
	}
	return list
}
