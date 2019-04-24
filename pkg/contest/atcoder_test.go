package contest

import (
	"reflect"
	"testing"
)

const (
	testContestID = "abc123"
)

func TestSqrapeQuizzes(t *testing.T) {
	at := NewAtCoder(testContestID)

	got, err := at.SqrapeQuizzes()
	if err != nil {
		t.Error(err)
	}

	expect := []string{"abc123_a", "abc123_b", "abc123_c", "abc123_d"}
	if !reflect.DeepEqual(got, expect) {
		t.Errorf("Could not get valid quizzes\ngot: %v\nexpect: %v",
			got, expect)
	}
}

func TestSqrapeSample(t *testing.T) {
	at := NewAtCoder(testContestID)

	ss, err := at.SqrapeSample("abc123_a")
	if err != nil {
		t.Error(err)
	}

	if got, expect := len(ss), 2; got != expect {
		t.Errorf("Invalid number of samples\ngot: %v\nexpect: %v",
			got, expect)
	}
}

func TestGetCsrfToken(t *testing.T) {
	at := NewAtCoder(testContestID)

	token, err := at.getCsrfToken(atCoderLoginURL)
	if err != nil {
		t.Error(err)
	}
	if token == "" {
		t.Error("Failed to get CSRF token")
	}
}
