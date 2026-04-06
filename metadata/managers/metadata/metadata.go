package metadata

import (
	"runtime"
	"sync"

	log "github.com/sirupsen/logrus"
	"github.com/astria-tv/astria-server/metadata/agents"
	"github.com/astria-tv/astria-server/metadata/db"
	mhelpers "github.com/astria-tv/astria-server/metadata/helpers"
)

// MetadataManager manages the metadata repository that is referenced by the files in the various
// libraries.
type MetadataManager struct {
	seriesCreationMutex sync.Mutex
	moviesCreationMutex sync.Mutex

	// Read/write lock for episode manipulation
	// TODO(Leon Handreke): Use proper locking,
	//  everywhere. This is a quickfix for the garbage collection routine.
	episodeLock sync.Map
	seasonLock  sync.Map
	seriesLock  sync.Map

	agent agents.MetadataRetrievalAgent

	eventBroker *metadataEventBroker
}

// NewMetadataManager creates a new MetadataManager
func NewMetadataManager(agent agents.MetadataRetrievalAgent) *MetadataManager {
	return &MetadataManager{
		agent:       agent,
		eventBroker: newMetadataEventBroker(),
	}
}

// RefreshAgentMetadataWithMissingArt loops over all series/episodes/seasons and movies with missing art (posters/backdrop) and tries to retrieve them.
func (m *MetadataManager) RefreshAgentMetadataWithMissingArt() {
	log.Debugln("Checking and updating media items for missing art.")
	wg := &sync.WaitGroup{}

	uuids := make(chan string, runtime.NumCPU())

	missingItems := db.ItemsWithMissingMetadata()
	log.Debugln(len(missingItems), " items appear to be missing art.")
	for _, UUID := range missingItems {
		wg.Add(1)
		go func(UUID string) {
			defer wg.Done()
			uuids <- UUID
			m.RefreshAgentMetadataForUUID(UUID)
			<-uuids
		}(UUID)
	}

	wg.Wait()
}

// RefreshCastForItemsWithMissingCast finds all movies and series that have no
// cast data and fetches it from the agent. This is a self-healing backfill:
// once every item has cast, the queries return nothing and this is a no-op.
func (m *MetadataManager) RefreshCastForItemsWithMissingCast() {
	movies := db.FindMoviesWithoutCast()
	series := db.FindSeriesWithoutCast()

	if len(movies) == 0 && len(series) == 0 {
		return
	}

	log.Infof("Backfilling cast for %d movies and %d series.", len(movies), len(series))

	for i := range movies {
		movie := &movies[i]
		mhelpers.WithLock(func() {
			m.refreshMovieCast(movie)
		}, movie.UUID)
	}

	for i := range series {
		s := &series[i]
		mhelpers.WithLock(func() {
			m.refreshSeriesCast(s)
		}, s.UUID)
	}
}

// RefreshAgentMetadataForUUID takes an UUID of a mediaitem and refreshes all metadata
func (m *MetadataManager) RefreshAgentMetadataForUUID(UUID string) bool {

	log.WithFields(log.Fields{"uuid": UUID}).
		Debugln("Looking to refresh metadata agent data.")
	movie, err := db.FindMovieByUUID(UUID)
	if err == nil {
		mhelpers.WithLock(func() {
			m.RefreshMovieMetadata(movie)
			db.SaveMovie(movie)
		}, movie.UUID)
		return true
	}

	series, err := db.FindSeriesByUUID(UUID)
	if err == nil {
		mhelpers.WithLock(func() {
			m.refreshSeriesMetadataFromAgent(series)
			db.SaveSeries(series)
		}, series.UUID)
		return true
	}

	season, err := db.FindSeasonByUUID(UUID)
	if err == nil {
		mhelpers.WithLock(func() {
			m.refreshSeasonMetadataFromAgent(season, season.GetSeries().TmdbID)
			db.SaveSeason(season)
		}, season.UUID)
		return true
	}

	episode, err := db.FindEpisodeByUUID(UUID)
	if err == nil {
		mhelpers.WithLock(func() {
			m.refreshEpisodeMetadataFromAgent(episode, episode.SeasonNum, episode.GetSeries().TmdbID)
			db.SaveEpisode(episode)
		}, episode.UUID)
		return true
	}

	return false
}
