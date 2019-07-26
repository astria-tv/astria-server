package agents

import (
	"gitlab.com/olaris/olaris-server/metadata/db"
)

// MetadataRetrievalAgent can retrieve metadata for media items.
type MetadataRetrievalAgent interface {
	UpdateMovieMD(*db.Movie) error
	UpdateSeasonMD(*db.Season, *db.Series) error
	UpdateEpisodeMD(*db.Episode, *db.Season, *db.Series) error
	UpdateSeriesMD(*db.Series) error
}
