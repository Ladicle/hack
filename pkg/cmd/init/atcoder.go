package init

import (
	"path/filepath"

	"github.com/Ladicle/hack/pkg/config"
	"github.com/Ladicle/hack/pkg/contest"
	"github.com/Ladicle/hack/pkg/format"
)

type atCoderInitializer struct {
	at  *contest.AtCoder
	cID string
	*format.HackRobot
}

func newAtCoderInitializer(hb *format.HackRobot, cid string) *atCoderInitializer {
	return &atCoderInitializer{
		HackRobot: hb,
		cID:       cid,
		at:        contest.NewAtCoder(cid),
	}
}

func (ai *atCoderInitializer) initAtCoderContest() error {
	ai.Info("Sure! I'll setup environment for %q contest.\n", ai.cID)
	ai.Start("Scraping %v quizzes :mag_right:", ai.cID)
	defer ai.End()
	qs, err := ai.at.ScrapeQuizzes()
	if err != nil {
		ai.Error()
		return err
	}

	ai.Start("Creating %v quiz directories :package:", len(qs))
	if err := contest.MkQuizDir(qs); err != nil {
		ai.Error()
		return err
	}
	for _, quiz := range qs {
		if err := ai.createSamples(quiz); err != nil {
			return err
		}
	}
	return nil
}

func (ai *atCoderInitializer) createSamples(quizID string) error {
	ai.Start("Scraping %v quizzes :inbox_tray:", quizID)
	defer ai.End()

	ss, err := ai.at.ScrapeSample(quizID)
	if err != nil {
		ai.Error()
		return err
	}

	qDir := filepath.Join(config.CurrentContestPath(), quizID)
	for _, sample := range ss {
		ai.Start("Scraping sample #%v :memo:", sample.ID)
		if err := contest.MkSample(qDir, sample); err != nil {
			ai.Error()
			return err
		}
	}
	return nil
}
