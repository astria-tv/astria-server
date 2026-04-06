package db_test

import (
	"testing"

	"github.com/astria-tv/astria-server/metadata/db"
)

var movie db.Movie

func setupTest(t *testing.T) func() {
	dbc := db.NewDb(db.DatabaseOptions{
		Connection: db.InMemory,
		LogMode:    false,
	})

	// Test teardown - return a closure for use by 'defer'
	return func() {
		// t is from the outer setupTest scope
		dbc.Close()
	}
}

func createMovieData() {
	mi := db.MediaItem{FilePath: "/tmp/test.mkv"}
	stream := db.Stream{CodecName: "test"}
	mf := db.MovieFile{MediaItem: mi, Streams: []db.Stream{stream}}

	movie = db.Movie{
		Title:         "Test",
		OriginalTitle: "Mad Max: Road Fury",
		MovieFiles:    []db.MovieFile{mf},
	}
	db.SaveMovie(&movie)
	ps := db.PlayState{
		MediaUUID: movie.UUID,
		Finished:  false, Playtime: 33, UserID: 1}
	db.SavePlayState(&ps)
}

func createData() {
	series := db.Series{Name: "All Episodes completely watched"}
	episode := &db.Episode{SeasonNum: 1, EpisodeNum: 1, Name: "AECW - Episode 1"}
	db.SaveEpisode(episode)
	db.SavePlayState(&db.PlayState{
		MediaUUID: episode.UUID,
		UserID:    1,
		Finished:  true, Playtime: 13,
	})

	episode2 := &db.Episode{SeasonNum: 1, EpisodeNum: 2, Name: "AECW - Episode 2"}
	db.SaveEpisode(episode2)
	db.SavePlayState(&db.PlayState{
		MediaUUID: episode2.UUID,
		UserID:    1,
		Finished:  true, Playtime: 14,
	})
	episode3 := &db.Episode{SeasonNum: 1, EpisodeNum: 3, Name: "AECW - Episode 3"}
	episode4 := &db.Episode{SeasonNum: 1, EpisodeNum: 4, Name: "AECW - Episode 4"}

	season := db.Season{Name: "Season 1", SeasonNumber: 1, Episodes: []*db.Episode{episode, episode2, episode3, episode4}}
	series.Seasons = []*db.Season{&season}
	db.CreateSeries(&series)

}

func createSeries1() {
	series2 := db.Series{Name: "Not finished watching an episode yet"}
	ep := &db.Episode{SeasonNum: 3, EpisodeNum: 3, Name: "NFY - Episode 3"}
	db.SaveEpisode(ep)
	db.SavePlayState(&db.PlayState{
		MediaUUID: ep.UUID,
		UserID:    1,
		Finished:  false, Playtime: 33,
	})

	ep2 := &db.Episode{SeasonNum: 3, EpisodeNum: 4, Name: "NFY - Episode 4"}
	s := db.Season{Name: "Season 3", SeasonNumber: 3, Episodes: []*db.Episode{ep, ep2}}
	series2.Seasons = []*db.Season{&s}
	db.CreateSeries(&series2)
}

func createSeries2() {
	series := db.Series{Name: "Next Season"}
	episode := &db.Episode{SeasonNum: 1, EpisodeNum: 1, Name: "NS - Episode 1"}
	db.SaveEpisode(episode)
	db.SavePlayState(&db.PlayState{
		MediaUUID: episode.UUID,
		UserID:    1,
		Finished:  true, Playtime: 13,
	})
	episode2 := &db.Episode{SeasonNum: 1, EpisodeNum: 2, Name: "NS - Episode 2"}
	db.SaveEpisode(episode2)
	db.SavePlayState(&db.PlayState{
		MediaUUID: episode2.UUID,
		UserID:    1,
		Finished:  true, Playtime: 14,
	})

	episode3 := &db.Episode{SeasonNum: 2, EpisodeNum: 1, Name: "NS - Episode S02E01"}
	episode4 := &db.Episode{SeasonNum: 2, EpisodeNum: 2, Name: "NS - Episode S02E02"}

	season := db.Season{Name: "Season 1", SeasonNumber: 1, Episodes: []*db.Episode{episode, episode2}}
	season2 := db.Season{Name: "Season 2", SeasonNumber: 2, Episodes: []*db.Episode{episode4, episode3}}
	series.Seasons = []*db.Season{&season, &season2}
	db.CreateSeries(&series)
}
