// Package db handles database queries for the metadata server
package db

import (
	"fmt"
	"os"
	"path"
	"strings"

	"github.com/jinzhu/gorm"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	"gitlab.com/olaris/olaris-server/helpers"
	"gitlab.com/olaris/olaris-server/metadata/db/dialects/mysql"
	"gitlab.com/olaris/olaris-server/metadata/db/dialects/postgres"
	"gitlab.com/olaris/olaris-server/metadata/db/dialects/sqlite"
	"gopkg.in/gormigrate.v1"
)

var db *gorm.DB

const (
	SQLite      = "sqlite3"
	MySQL       = "mysql"
	PostgresSQL = "postgres"
	CockroachDB = "cockroachdb"
)

// DatabaseOptions holds information about how the database instance should be initialized
type DatabaseOptions struct {
	Connection string
	LogMode    bool
}

func getDefaultDbPath() (string, error) {
	dbDir := viper.GetString("server.sqliteDir")
	if err := helpers.EnsurePath(dbDir); err != nil {
		return "", err
	}

	return path.Join(dbDir, "metadata.db"), nil
}

func defaultDb(logMode bool) *gorm.DB {
	dbPath, err := getDefaultDbPath()
	if err != nil {
		panic(fmt.Sprintf("failed to get default database path: %s\n", err))
	}
	db, err = sqlite.NewSQLiteDatabase(dbPath, logMode)
	if err != nil {
		panic(fmt.Sprintf("failed to connect database: %s\n", err))
	}

	log.WithField("path", dbPath).Println("using default (sqlite3) database")
	return db
}

// NewDb initializes a new database instance.
func NewDb(options DatabaseOptions) *gorm.DB {
	var err error

	databaseTokens := strings.Split(options.Connection, "://")
	if len(databaseTokens) == 0 {
		db = defaultDb(options.LogMode)
	} else if len(databaseTokens) == 2 {
		engine := databaseTokens[0]
		connection := databaseTokens[1]
		switch engine {
		case SQLite:
			db, err = sqlite.NewSQLiteDatabase(connection, options.LogMode)
			if err != nil {
				log.Errorf("%s", err)
				os.Exit(1)
			}
			log.Println("using sqlite3 database driver")
		case MySQL:
			db, err = mysql.NewMySQLDatabase(connection, options.LogMode)
			if err != nil {
				log.Errorf("%s", err)
				os.Exit(1)
			}
			log.Println("using MySQL database driver")
		case CockroachDB, PostgresSQL:
			// CockroachDB uses the Postgres driver
			// https://www.cockroachlabs.com/docs/stable/build-a-go-app-with-cockroachdb-gorm.html
			db, err = postgres.NewPostgresDatabase(connection, options.LogMode)
			if err != nil {
				log.Errorf("%s", err)
				os.Exit(1)
			}
			log.Println("using postgres database driver")
		default:
			log.Errorf("unknown database engine: %s", engine)
			os.Exit(1)
		}
	} else {
		log.Debugf("unable to parse database connection string: %s, defaulting to sqlite3", options.Connection)
		db = defaultDb(options.LogMode)
	}

	err = migrateSchema(db)
	if err != nil {
		log.Fatalf("failed to migrate database: %s", err)
	}

	return db
}

var allModels = []interface{}{
	&Movie{}, &MovieFile{}, &Library{}, &Series{}, &Season{}, &Episode{},
	&EpisodeFile{}, &User{}, &Invite{}, &PlayState{}, &Stream{},
	&WatchlistItem{},
}

func initSchema(tx *gorm.DB) error {
	return tx.AutoMigrate(allModels...).Error
}

func migrateSchema(db *gorm.DB) error {
	m := gormigrate.New(db, gormigrate.DefaultOptions, allMigrations())
	m.InitSchema(initSchema)
	if err := m.Migrate(); err != nil {
		return err
	}

	// Apply any new column/table additions outside of versioned migrations.
	return db.AutoMigrate(allModels...).Error
}
