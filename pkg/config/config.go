package config

import (
	"os"
	"path"
	"strings"

	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"

	"github.com/astria-tv/astria-server/helpers"
)

var ConfigDir string

func GetDefaultConfigDir() string {
	defaultConfigDir := path.Join(helpers.GetHome(), ".config", "astria")
	if configDirEnv := os.Getenv("ASTRIA_CONFIG_DIR"); configDirEnv != "" {
		defaultConfigDir = configDirEnv
	}

	return defaultConfigDir
}

func InitViper() {
	viper.SetConfigName("astria")
	viper.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))
	viper.SetEnvPrefix("astria")
	viper.AutomaticEnv()

	viper.AddConfigPath(ConfigDir)
	viper.Set("configdir", ConfigDir)
	if err := viper.ReadInConfig(); err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); ok {
			// the user has no config file
		} else {
			logrus.WithError(err).WithField("configFile", viper.ConfigFileUsed()).Warnln("An error occurred while reading config file, contents are being ignored.")
		}
	}

	config := &Config{}
	if err := viper.Unmarshal(config); err != nil {
		logrus.Debugf("error applying configuration: %s\n", err.Error())
	}
}
