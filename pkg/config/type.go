package config

// Config manages contest data
type Config struct {
	Contest      Contest
	CurrentQuizz string
}

// Contest contains contest data
type Contest struct {
	Name    string
	URL     string
	Path    string
	Quizzes []string
}
