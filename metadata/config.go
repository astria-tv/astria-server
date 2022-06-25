package metadata

import (
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

// FlagFuncs returns two functions: a function to register the flags on the
// command, and a function to register those flags with Viper. These two things
// need to happen in two different places.
func FlagFuncs() (flagFunc func(*cobra.Command) error, viperFunc func(*cobra.Command) error) {
	flagFunc = func(cmd *cobra.Command) error {
		cmd.Flags().String("db-conn", "", "sets the database connection string")
		cmd.Flags().String("sqlite-dir", "", "Path where the database is stored if using SQLite. (defaults to <config_dir>/metadb)")
		cmd.Flags().Bool("scan-hidden", false, "sets whether to scan hidden directories (directories starting with a .)")
		return nil
	}

	viperFunc = func(cmd *cobra.Command) error {
		_ = viper.BindPFlag("database.connection", cmd.Flags().Lookup("db-conn"))
		_ = viper.BindPFlag("server.sqliteDir", cmd.Flags().Lookup("sqlite-dir"))
		_ = viper.BindPFlag("metadata.scanHidden", cmd.Flags().Lookup("scan-hidden"))
		return nil
	}

	return
}
