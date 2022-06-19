package config

import "github.com/spf13/viper"

// Config is the base struct populated from the configuration file on disk by
// Viper.
// NOTE: This is not actually used anywhere yet. We are currently pulling
// directly from Viper when we need a setting.
type Config struct {
	Debug    Debug
	Server   Server
	Database Database
	RClone   RClone
	Metadata Metadata
}

// Debug contains settings to enable additional debug information in various
// places throughout the app.
type Debug struct {
	StreamingPages bool
	TranscoderLog  bool
}

// Server contains settings that affect the built-in web server.
type Server struct {
	Port             int
	Verbose          bool
	DBLog            bool
	SQLiteDir        string
	CacheDir         string
	ConfigDir        string
	DirectFileAccess bool
}

// Database contains settings that affect the database connection.
type Database struct {
	DBConn string
}

type RClone struct {
	ConfigFile string
}

type Metadata struct {
	ScanHidden bool
}

// FromViper creates a new Config from a Viper configuration.
func FromViper() (*Config, error) {
	var cfg Config
	err := viper.Unmarshal(&cfg)
	return &cfg, err
}
