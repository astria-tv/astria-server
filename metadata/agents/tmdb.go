package agents

import (
	"time"

	"github.com/pkg/errors"
	"github.com/ryanbradynd05/go-tmdb"
	log "github.com/sirupsen/logrus"
	"github.com/astria-tv/astria-server/metadata/db"
)

const tmdbAPIKey = "0cdacd9ab172ac6ff69c8d84b2c938a8"

// TmdbAgent is a wrapper around themoviedb
type TmdbAgent struct {
	Tmdb *tmdb.TMDb
}

// NewTmdbAgent creates a new themoviedb agent.
func NewTmdbAgent() *TmdbAgent {
	return &TmdbAgent{tmdb.Init(tmdb.Config{
		APIKey:   tmdbAPIKey,
		Proxies:  nil,
		UseProxy: false,
	})}
}

// ParseTmdbDate parses a date string returned from the TMDB API
func ParseTmdbDate(tmdbDate string) (time.Time, error) {
	return time.Parse("2006-01-02", tmdbDate)
}

// UpdateEpisodeMD updates the metadata information for the given episode.
func (a *TmdbAgent) UpdateEpisodeMD(
	episode *db.Episode, seriesTmdbID int, seasonNum int, episodeNum int,
) error {
	fullEpisode, err := a.Tmdb.GetTvEpisodeInfo(
		seriesTmdbID, seasonNum, episodeNum, nil)
	if err != nil {
		return errors.Wrap(err, "Could not retrieve episode data from TMDB")
	}

	episode.AirDate = fullEpisode.AirDate
	episode.Name = fullEpisode.Name
	episode.TmdbID = fullEpisode.ID
	episode.Overview = fullEpisode.Overview
	episode.StillPath = fullEpisode.StillPath
	log.WithFields(log.Fields{"episodeName": episode.Name, "tmdbId": episode.TmdbID}).
		Debugln("found episode metadata")

	return nil

}

// UpdateSeasonMD updates the metadata information for the given season
func (a *TmdbAgent) UpdateSeasonMD(season *db.Season, seriesTmdbID int, seasonNum int) error {
	log.
		WithFields(log.Fields{
			"seasonNumber": seasonNum,
			"seriesTmdbID": seriesTmdbID}).
		Debugln("Looking for season metadata.")

	fullSeason, err := a.Tmdb.GetTvSeasonInfo(seriesTmdbID, seasonNum, nil)
	if err != nil {
		log.WithFields(log.Fields{"error": err}).Warnln("Could not grab season information.")
		return err
	}

	season.AirDate = fullSeason.AirDate
	season.Overview = fullSeason.Overview
	season.Name = fullSeason.Name
	season.TmdbID = fullSeason.ID
	season.PosterPath = fullSeason.PosterPath
	log.WithFields(log.Fields{"seasonName": season.Name, "tmdbId": season.TmdbID}).
		Debugln("Found season metadata.")
	return nil
}

// UpdateSeriesMD updates the metadata information for the given series.
func (a *TmdbAgent) UpdateSeriesMD(series *db.Series, tmdbID int) error {
	fullTv, err := a.Tmdb.GetTvInfo(tmdbID, nil)

	if err != nil {
		log.
			WithFields(log.Fields{
				"tmdbID": tmdbID,
				"error":  err}).
			Debugln("Could not grab full TV details.")
		return err
	}

	firstAirDate, _ := ParseTmdbDate(fullTv.FirstAirDate)

	series.Name = fullTv.Name
	series.OriginalName = fullTv.OriginalName
	series.FirstAirDate = fullTv.FirstAirDate
	series.FirstAirYear = uint64(firstAirDate.Year())
	series.Overview = fullTv.Overview
	series.Status = fullTv.Status
	series.Type = fullTv.Type
	series.BackdropPath = fullTv.BackdropPath
	series.PosterPath = fullTv.PosterPath
	return nil
}

// refreshAndSaveMovieMetadata updates
func (a *TmdbAgent) UpdateMovieMD(movie *db.Movie, tmdbID int) error {
	r, err := a.Tmdb.GetMovieInfo(tmdbID, nil)

	if err != nil {
		return errors.Wrap(err, "Failed to query TMDB for movie metadata")
	}

	releaseDate, _ := ParseTmdbDate(r.ReleaseDate)

	movie.TmdbID = r.ID
	movie.Title = r.Title
	movie.OriginalTitle = r.OriginalTitle
	movie.ReleaseDate = r.ReleaseDate
	movie.Year = uint64(releaseDate.Year())
	movie.Overview = r.Overview
	movie.BackdropPath = r.BackdropPath
	movie.PosterPath = r.PosterPath
	movie.ImdbID = r.ImdbID

	return nil
}

// TmdbSearchMovie directly exposes the TMDb search interface
func (a *TmdbAgent) TmdbSearchMovie(
	name string,
	options map[string]string,
) (*tmdb.MovieSearchResults, error) {
	return a.Tmdb.SearchMovie(name, options)
}

// TmdbSearchTv directly exposes the TMDb search interface
func (a *TmdbAgent) TmdbSearchTv(
	name string,
	options map[string]string,
) (*tmdb.TvSearchResults, error) {
	return a.Tmdb.SearchTv(name, options)
}

// GetMovieCast retrieves the cast for the given movie from TMDB.
// It returns CastRole structs with the Person populated (TmdbID, Name, ProfilePath).
// The caller is responsible for persisting Person records and setting PersonID.
func (a *TmdbAgent) GetMovieCast(tmdbID int) ([]db.CastRole, error) {
	credits, err := a.Tmdb.GetMovieCredits(tmdbID, nil)
	if err != nil {
		return nil, err
	}

	limit := len(credits.Cast)
	if limit > db.MaxCastRoles {
		limit = db.MaxCastRoles
	}

	roles := make([]db.CastRole, limit)
	for i := 0; i < limit; i++ {
		c := credits.Cast[i]
		roles[i] = db.CastRole{
			Character: c.Character,
			Order:     c.Order,
			Person: db.Person{
				TmdbID:      c.ID,
				Name:        c.Name,
				ProfilePath: c.ProfilePath,
			},
		}
	}
	return roles, nil
}

// GetSeriesCast retrieves the cast for the given TV series from TMDB.
func (a *TmdbAgent) GetSeriesCast(tmdbID int) ([]db.CastRole, error) {
	credits, err := a.Tmdb.GetTvCredits(tmdbID, nil)
	if err != nil {
		return nil, err
	}

	limit := len(credits.Cast)
	if limit > db.MaxCastRoles {
		limit = db.MaxCastRoles
	}

	roles := make([]db.CastRole, limit)
	for i := 0; i < limit; i++ {
		c := credits.Cast[i]
		roles[i] = db.CastRole{
			Character: c.Character,
			Order:     c.Order,
			Person: db.Person{
				TmdbID:      c.ID,
				Name:        c.Name,
				ProfilePath: c.ProfilePath,
			},
		}
	}
	return roles, nil
}

// UpdatePersonMD fetches full person details from TMDB and updates the struct.
func (a *TmdbAgent) UpdatePersonMD(person *db.Person, tmdbID int) error {
	p, err := a.Tmdb.GetPersonInfo(tmdbID, nil)
	if err != nil {
		return err
	}
	person.Name = p.Name
	person.ProfilePath = p.ProfilePath
	person.Biography = p.Biography
	person.Birthday = p.Birthday
	person.Deathday = p.Deathday
	person.PlaceOfBirth = p.PlaceOfBirth
	return nil
}
