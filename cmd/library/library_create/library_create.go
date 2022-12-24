package library_create

import (
	"github.com/spf13/viper"
	"gitlab.com/olaris/olaris-server/helpers"
	"path"
	"time"

	"github.com/goava/di"
	"github.com/spf13/cobra"

	"gitlab.com/olaris/olaris-server/metadata/app"
	"gitlab.com/olaris/olaris-server/metadata/db"
)

const defaultTimeOffset = -24 * time.Hour

func Options() di.Option {
	return di.Options(
		di.Provide(NewLibraryCreateCommand, di.Tags{"type": "library_create"}),
		di.Invoke(RegisterLibraryCreateCommand),
	)
}

func RegisterLibraryCreateCommand(deps struct {
	di.Inject
	LibraryCommand       *cobra.Command `di:"type=library"`
	LibraryCreateCommand *cobra.Command `di:"type=library_create"`
}) {
	deps.LibraryCommand.AddCommand(deps.LibraryCreateCommand)
}

func NewLibraryCreateCommand() *cobra.Command {
	var name string
	var filePath string
	var mediaType int
	var backendType int
	var rcloneName string

	c := &cobra.Command{
		Use:   "create",
		Short: "Create a new library",
		PreRun: func(cmd *cobra.Command, args []string) {
			_ = viper.BindPFlag("server.dbLog", cmd.Flags().Lookup("db-log"))
			_ = viper.BindPFlag("server.sqliteDir", cmd.Flags().Lookup("sqlite-dir"))
			_ = viper.BindPFlag("database.connection", cmd.Flags().Lookup("db-conn"))
			_ = viper.BindPFlag("metadata.scanHidden", cmd.Flags().Lookup("scan-hidden"))
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			mctx := app.NewDefaultMDContext()
			defer mctx.Db.Close()

			lib := &db.Library{Name: name, FilePath: filePath, Kind: db.MediaType(mediaType), Backend: backendType, RcloneName: rcloneName}

			// Make sure we don't initialize the library with zero time (issue with strict mode in MySQL)
			lib.RefreshStartedAt = time.Now().Add(defaultTimeOffset)
			lib.RefreshCompletedAt = time.Now().Add(defaultTimeOffset)

			err := db.AddLibrary(lib)
			return err
		},
	}

	c.Flags().StringVar(&name, "name", "", "A name for this library")
	c.MarkFlagRequired("name")

	c.Flags().StringVar(&filePath, "path", "", "Path for this library")
	c.MarkFlagRequired("path")

	c.Flags().IntVar(&mediaType, "media-type", 0, "Media type, 0 for Movies, 1 for Series")
	c.Flags().IntVar(&backendType, "backend-type", 0, "Backend type, 0 for Local, 1 for Rclone")
	c.Flags().StringVar(&rcloneName, "rclone-name", "", "Name for the Rclone remote")

	c.Flags().Bool("db-log", false, "sets whether the database should log queries")
	c.Flags().String("db-conn", "", "sets the database connection string")
	c.Flags().String("sqlite-dir", path.Join(helpers.BaseConfigDir(), "metadb"), "Path where the database is stored if using SQLite")

	return c
}
