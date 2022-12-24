package root

import (
	"github.com/goava/di"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"gitlab.com/olaris/olaris-server/helpers"
	"gitlab.com/olaris/olaris-server/pkg/config"
	"strings"
)

func Options() di.Option {
	return di.Options(
		di.Provide(NewRootCommand, di.Tags{"type": "root"}),
	)
}

func NewRootCommand() *cobra.Command {
	c := &cobra.Command{
		Use: "olaris",
		PersistentPreRun: func(cmd *cobra.Command, args []string) {
			_ = viper.BindPFlag("config.dir", cmd.Flags().Lookup("config-dir"))
			_ = viper.BindPFlag("server.cacheDir", cmd.Flags().Lookup("cache-dir"))
			_ = viper.BindPFlag("server.verbose", cmd.Flags().Lookup("verbose"))
			_ = viper.BindPFlag("server.directFileAccess", cmd.Flags().Lookup("allow-direct-file-access"))
			_ = viper.BindPFlag("debug.streamingPages", cmd.Flags().Lookup("enable-streaming-debug-pages"))
			_ = viper.BindPFlag("debug.transcoderLog", cmd.Flags().Lookup("write-transcoder-log"))
			_ = viper.BindPFlag("rclone.configFile", cmd.Flags().Lookup("rclone-config"))

			// Configure Viper
			viper.SetConfigName("olaris")
			viper.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))
			viper.SetEnvPrefix("olaris")
			viper.AutomaticEnv()
			viper.AddConfigPath(viper.GetString("config.dir"))
		},
	}

	c.PersistentFlags().Bool("allow-direct-file-access", false, "Whether accessing files directly by path (without a valid JWT) is allowed")
	c.PersistentFlags().Bool("enable-streaming-debug-pages", false, "Whether to enable debug pages in the streaming server")
	c.PersistentFlags().Bool("write-transcoder-log", true, "Whether to write transcoder output to logfile")
	c.PersistentFlags().BoolP("verbose", "v", true, "verbose logging")
	c.PersistentFlags().String("config-dir", config.GetDefaultConfigDir(), "Default configuration directory for config files")
	c.PersistentFlags().String("rclone-config", helpers.GetDefaultRcloneConfigPath(), "Default rclone configuration file")
	c.PersistentFlags().String("cache-dir", helpers.GetDefaultCacheDir(), "Cache directory for transcoding an other temporarily files")

	return c
}
