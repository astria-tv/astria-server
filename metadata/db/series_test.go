package db_test

import (
	"testing"

	"github.com/astria-tv/astria-server/metadata/db"
)

func TestBatchUnwatchedEpisodesInSeriesCounts(t *testing.T) {
	defer setupTest(t)()

	createSeries1()
	createSeries2()
	createData()

	// Get all series
	series, err := db.FindAllSeries(nil)
	if err != nil {
		t.Fatalf("Failed to find series: %v", err)
	}

	if len(series) == 0 {
		t.Fatal("Expected at least one series")
	}

	seriesIDs := make([]uint, len(series))
	for i, s := range series {
		seriesIDs[i] = s.ID
	}

	batchCounts := db.BatchUnwatchedEpisodesInSeriesCounts(seriesIDs, 1)

	// Verify batch counts match individual counts
	for _, s := range series {
		individual := db.UnwatchedEpisodesInSeriesCount(s.ID, 1)
		batch := batchCounts[s.ID]
		if individual != batch {
			t.Errorf("Mismatch for series %d (%s): individual=%d, batch=%d", s.ID, s.Name, individual, batch)
		}
		if batch == 0 {
			t.Errorf("Expected non-zero unwatched count for series %d (%s)", s.ID, s.Name)
		}
	}
}
