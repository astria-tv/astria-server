package db

import (
	"github.com/jinzhu/gorm"
	log "github.com/sirupsen/logrus"
)

// MaxCastRoles is the maximum number of cast roles stored per movie or series.
const MaxCastRoles = 20

// Person stores a unique individual from TMDB. The TmdbID (TMDB person ID) is
// the canonical identifier — the same person appearing across many movies and
// series is stored exactly once.
type Person struct {
	gorm.Model
	TmdbID       int `gorm:"unique_index"`
	Name         string
	ProfilePath  string
	Biography    string `gorm:"type:text"`
	Birthday     string
	Deathday     string
	PlaceOfBirth string
}

// CastRole links a Person to a Movie or Series via a polymorphic owner.
type CastRole struct {
	gorm.Model
	PersonID  uint
	Person    Person
	Character string
	// Display order as returned by TMDB (lower = more prominent).
	Order     int
	OwnerID   uint
	OwnerType string
}

// FindCastForMovie returns cast roles for the given movie, ordered by prominence.
func FindCastForMovie(movieID uint) []CastRole {
	return findCastRoles(movieID, "movies")
}

// FindCastForSeries returns cast roles for the given series, ordered by prominence.
func FindCastForSeries(seriesID uint) []CastRole {
	return findCastRoles(seriesID, "series")
}

func findCastRoles(ownerID uint, ownerType string) []CastRole {
	var roles []CastRole
	db.Preload("Person").
		Where("owner_id = ? AND owner_type = ?", ownerID, ownerType).
		Order("`order` ASC").
		Find(&roles)
	return roles
}

// FindPersonByTmdbID looks up a person by their TMDB ID.
func FindPersonByTmdbID(tmdbID int) (*Person, error) {
	var person Person
	if err := db.Where("tmdb_id = ?", tmdbID).Take(&person).Error; err != nil {
		return nil, err
	}
	return &person, nil
}

// SavePerson inserts or updates a person record.
func SavePerson(person *Person) error {
	return db.Save(person).Error
}

// GetOrCreatePersonByTmdbID returns the existing person or creates a new one.
func GetOrCreatePersonByTmdbID(tmdbID int, name string, profilePath string) (*Person, error) {
	person, err := FindPersonByTmdbID(tmdbID)
	if err == nil {
		return person, nil
	}
	person = &Person{
		TmdbID:      tmdbID,
		Name:        name,
		ProfilePath: profilePath,
	}
	if err := db.Create(person).Error; err != nil {
		// Handle race condition: another goroutine may have created it
		person, err = FindPersonByTmdbID(tmdbID)
		if err != nil {
			return nil, err
		}
	}
	return person, nil
}

// UpdatePersonDetails updates biography and other details for an existing person.
func UpdatePersonDetails(person *Person) error {
	return db.Save(person).Error
}

// UpdateCastForMovie replaces all cast roles for the given movie.
func UpdateCastForMovie(movieID uint, roles []CastRole) error {
	return updateCastRoles(movieID, "movies", roles)
}

// UpdateCastForSeries replaces all cast roles for the given series.
func UpdateCastForSeries(seriesID uint, roles []CastRole) error {
	return updateCastRoles(seriesID, "series", roles)
}

func updateCastRoles(ownerID uint, ownerType string, roles []CastRole) error {
	tx := db.Begin()

	if err := tx.Unscoped().
		Where("owner_id = ? AND owner_type = ?", ownerID, ownerType).
		Delete(&CastRole{}).Error; err != nil {
		tx.Rollback()
		return err
	}

	for i := range roles {
		roles[i].OwnerID = ownerID
		roles[i].OwnerType = ownerType
		if err := tx.Create(&roles[i]).Error; err != nil {
			tx.Rollback()
			return err
		}
	}

	return tx.Commit().Error
}

// DeleteCastForMovie removes all cast role entries for a movie.
func DeleteCastForMovie(movieID uint) error {
	return db.Unscoped().
		Where("owner_id = ? AND owner_type = ?", movieID, "movies").
		Delete(&CastRole{}).Error
}

// DeleteCastForSeries removes all cast role entries for a series.
func DeleteCastForSeries(seriesID uint) error {
	return db.Unscoped().
		Where("owner_id = ? AND owner_type = ?", seriesID, "series").
		Delete(&CastRole{}).Error
}

// GarbageCollectOrphanedPeople removes Person records that are no longer
// referenced by any CastRole. This is intentionally lenient — person records
// are cheap and useful for caching, so we only clean up truly orphaned ones.
func GarbageCollectOrphanedPeople() {
	result := db.Exec(
		"DELETE FROM people WHERE id NOT IN (SELECT DISTINCT person_id FROM cast_roles)")
	if result.Error != nil {
		log.Warnln("Failed to garbage collect orphaned people:", result.Error)
	}
}

// FindMoviesWithoutCast returns all movies that have no CastRole rows.
func FindMoviesWithoutCast() []Movie {
	var movies []Movie
	db.Where("id NOT IN (SELECT DISTINCT owner_id FROM cast_roles WHERE owner_type = ?)", "movies").
		Find(&movies)
	return movies
}

// FindSeriesWithoutCast returns all series that have no CastRole rows.
func FindSeriesWithoutCast() []Series {
	var series []Series
	db.Where("id NOT IN (SELECT DISTINCT owner_id FROM cast_roles WHERE owner_type = ?)", "series").
		Find(&series)
	return series
}
