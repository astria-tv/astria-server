package agents_test

import (
	"github.com/stretchr/testify/assert"
	"github.com/astria-tv/astria-server/metadata/agents"
	"github.com/astria-tv/astria-server/metadata/db"
	"testing"
)

func TestSeasonLookup(t *testing.T) {
	const testSeasonTmdbID = 2426
	var season db.Season
	a := agents.NewTmdbAgent()

	a.UpdateSeasonMD(&season, testSeasonTmdbID, 1)

	assert.EqualValues(t, 7625, season.TmdbID)
}

func TestTmdbMovieLookup(t *testing.T) {
	movie := db.Movie{}
	a := agents.NewTmdbAgent()

	err := a.UpdateMovieMD(&movie, 76341)
	assert.NoError(t, err)

	assert.Equal(t, "Mad Max: Fury Road", movie.OriginalTitle)
	assert.EqualValues(t, 76341, movie.TmdbID)
}
