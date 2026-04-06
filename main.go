package main

import (
	"os"

	"github.com/sirupsen/logrus"
	"github.com/spf13/pflag"
	"github.com/spf13/viper"

	"github.com/astria-tv/astria-server/cmd"
	"github.com/astria-tv/astria-server/helpers"
	"github.com/astria-tv/astria-server/pkg/config"
)

func main() {
	config.RegisterFlags(registerGlobalFlags)
	config.InitViper()

	err := cmd.New().GetCobraCommand().Execute()
	if err != nil {
		logrus.Fatal(err)
	}

	os.Exit(0)
}

func registerGlobalFlags(fs *pflag.FlagSet) {
	fs.Bool("allow_direct_file_access", false, "Whether accessing files directly by path (without a valid JWT) is allowed")
	fs.Bool("enable_streaming_debug_pages", false, "Whether to enable debug pages in the streaming server")
	fs.Bool("write_transcoder_log", true, "Whether to write transcoder output to logfile")

	fs.StringVar(&config.ConfigDir, "config_dir", config.GetDefaultConfigDir(), "Default configuration directory for config files")
	fs.String("rclone_config", helpers.GetDefaultRcloneConfigPath(), "Default rclone configuration file")
	fs.String("cache_dir", helpers.GetDefaultCacheDir(), "Cache directory for transcoding an other temporarily files")

	viper.BindPFlag("server.cacheDir", fs.Lookup("cache_dir"))
	viper.BindPFlag("server.directFileAccess", fs.Lookup("allow_direct_file_access"))
	viper.BindPFlag("debug.streamingPages", fs.Lookup("enable_streaming_debug_pages"))
	viper.BindPFlag("debug.transcoderLog", fs.Lookup("write_transcoder_log"))
	viper.BindPFlag("rclone.configFile", fs.Lookup("rclone_config"))
}
