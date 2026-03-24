package db_test

import (
	"testing"

	"gitlab.com/olaris/olaris-server/metadata/db"
)

func TestAllPlayState(t *testing.T) {
	defer setupTest(t)()

	createData()

	pss := db.LatestPlayStates(1, 1)
	if len(pss) != 1 {
		t.Error("Expected one PlayState to return got", len(pss), "instead")
	}

	pss = db.LatestPlayStates(2, 1)
	if len(pss) != 2 {
		t.Error("Expected two PlayStates to return got", len(pss), "instead")
	}
}

func TestContinueMovie(t *testing.T) {
	defer setupTest(t)()
	createMovieData()

	movies := db.UpNextMovies(1)
	if movies[0].Title != "Test" {
		t.Error("Got the wrong movie expected Test but got:", movies[0].Title)
	}
}
func TestContinuePlayResume(t *testing.T) {
	defer setupTest(t)()

	createSeries1()
	createSeries2()
	createData()

	episodes := db.UpNextEpisodes(1)
	if len(episodes) != 3 {
		t.Errorf("exepected %v episodes got %v instead", 3, len(episodes))
	} else {

		if episodes[0].Name != "NFY - Episode 3" {
			t.Errorf("Expected the first Episode to be resumed to be %s got %s instead\n", "NFY - Episode 3", episodes[0].Name)
		}

		if episodes[1].Name != "NS - Episode S02E01" {
			t.Errorf("Expected the second Episode to be resumed to be NS - Episode S02E01 got %s instead\n", episodes[1].Name)
		}

		if episodes[2].Name != "AECW - Episode 3" {
			t.Errorf("Expected the second Episode to be resumed to be Episode 3 got %s instead\n", episodes[2].Name)
		}

		count := db.UnwatchedEpisodesInSeasonCount(1, 1)
		if count != 2 {
			t.Errorf("Expected the amount of unwatched episodes in the season to be 2 got %v instead\n", count)
		}

	}
}
