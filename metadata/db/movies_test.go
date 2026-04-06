package db_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/astria-tv/astria-server/metadata/db"
)

func TestUUIDable(t *testing.T) {
	defer setupTest(t)()

	createMovieData()

	if movie.UUID == "" || movie.MovieFiles[0].UUID == "" {
		t.Errorf("Movie/File was created without a UUID\n")
	} else {
		t.Log("Movie UUID:", movie.UUID)
	}
}

func TestSearchMovieByTitle(t *testing.T) {
	defer setupTest(t)()
	createMovieData()
	var movies []db.Movie
	movies = db.SearchMovieByTitle("max")
	if len(movies) == 0 {
		t.Error("Did not get any movies while searching")
		return
	}

	if movies[0].OriginalTitle != "Mad Max: Road Fury" {
		t.Error("Did not get the correct movie while searching")
	}
}

func TestCollectMovie(t *testing.T) {
	defer setupTest(t)()

	createMovieData()

	mov := db.FirstMovie()

	if len(mov.MovieFiles) != 0 {
		t.Error("Expected no movie files but got any still")
	}

	db.CollectMovieInfo(&mov)

	assert.Len(t, mov.MovieFiles, 1)
	assert.Len(t, mov.MovieFiles[0].Streams, 1)
}
