package config

import (
	"gitlab.com/olaris/olaris-server/helpers"
	"path"
)

// GetDefaultConfigDir returns the default location for the configuration file.
func GetDefaultConfigDir() string {
	return path.Join(helpers.GetHome(), ".config", "olaris")
}
