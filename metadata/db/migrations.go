package db

import (
	"github.com/jinzhu/gorm"
	"gopkg.in/gormigrate.v1"
)

// allMigrations returns the ordered list of all database migrations.
// Each migration should use the tx parameter (not the global db) to ensure
// transactional safety.
func allMigrations() []*gormigrate.Migration {
	return []*gormigrate.Migration{
		{
			ID:      "2019-06-19-new-filepaths",
			Migrate: migrateNewFilepaths,
		},
		{
			ID:      "2019-08-03-remove-unidentified-episodes",
			Migrate: migrateRemoveUnidentifiedEpisodes,
		},
		{
			ID:      "2019-08-03-remove-unidentified-movies",
			Migrate: migrateRemoveUnidentifiedMovies,
		},
	}
}

// migrateNewFilepaths prefixes all file paths with "local#" to support
// the rclone-based library types introduced alongside file locators.
func migrateNewFilepaths(tx *gorm.DB) error {
	type MovieFile struct {
		gorm.Model
		FilePath string
	}
	var movieFiles []MovieFile
	tx.Find(&movieFiles)
	for _, f := range movieFiles {
		f.FilePath = "local#" + f.FilePath
		tx.Save(f)
	}

	type EpisodeFile struct {
		gorm.Model
		FilePath string
	}
	var episodeFiles []EpisodeFile
	tx.Find(&episodeFiles)
	for _, f := range episodeFiles {
		f.FilePath = "local#" + f.FilePath
		tx.Save(f)
	}

	return nil
}

// migrateRemoveUnidentifiedEpisodes deletes episodes that were never
// identified (tmdb_id = 0), leaving only their EpisodeFiles behind.
func migrateRemoveUnidentifiedEpisodes(tx *gorm.DB) error {
	return tx.Exec("DELETE FROM episodes WHERE tmdb_id = 0;").Error
}

// migrateRemoveUnidentifiedMovies deletes movies that were never
// identified (tmdb_id = 0).
func migrateRemoveUnidentifiedMovies(tx *gorm.DB) error {
	return tx.Exec("DELETE FROM movies WHERE tmdb_id = 0;").Error
}
