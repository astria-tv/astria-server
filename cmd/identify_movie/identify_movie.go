package identify_movie

import (
	"fmt"
	"strings"

	"github.com/spf13/cobra"

	"github.com/astria-tv/astria-server/filesystem"
	"github.com/astria-tv/astria-server/metadata/agents"
	"github.com/astria-tv/astria-server/metadata/app"
	"github.com/astria-tv/astria-server/metadata/db"
	"github.com/astria-tv/astria-server/metadata/managers/metadata"
	"github.com/astria-tv/astria-server/pkg/cmd"
)

func NewIdentifyMovieCommand() *cmd.CobraCommand {
	var filePath string
	var agent string
	var id int
	var dbConn string
	var dbLog bool

	c := &cobra.Command{
		Use:   "movie",
		Short: "Identify a movie",
		Long:  "Identify a movie pased on it's file location path.\n* To identify a local file use local#/path/to/file.\n* To identify a Rclone file use rclone#/path/to/file.",
		RunE: func(cmd *cobra.Command, args []string) error {
			var a agents.MetadataRetrievalAgent
			switch strings.ToLower(agent) {
			case "tmdb":
				a = agents.NewTmdbAgent()
			default:
				return fmt.Errorf("unknown agent: %s", agent)
			}

			fl, err := filesystem.ParseFileLocator(filePath)
			if err != nil {
				return err
			}

			n, err := filesystem.GetNodeFromFileLocator(fl)
			if err != nil {
				return err
			}

			dbOptions := db.DatabaseOptions{
				Connection: dbConn,
				LogMode:    dbLog,
			}
			mctx := app.NewMDContext(dbOptions, a)
			defer mctx.Db.Close()

			f, err := db.FindMovieFileByPath(n)
			if err != nil {
				return err
			}

			mm := metadata.NewMetadataManager(a)
			movie, err := mm.GetOrCreateMovieByTmdbID(id)
			if err != nil {
				return err
			}

			f.Movie = *movie
			db.SaveMovieFile(f)
			return nil
		},
	}

	c.Flags().StringVar(&filePath, "path", "", "Path of the movie file")
	c.MarkFlagRequired("path")

	c.Flags().IntVar(&id, "id", 0, "ID of movie from agent")
	c.MarkFlagRequired("id")

	c.Flags().StringVar(&agent, "agent", "tmdb", "Agent, defaults to tmdb")
	c.Flags().StringVar(&dbConn, "db-conn", "", "sets the database connection string")
	c.Flags().BoolVar(&dbLog, "db-log", false, "sets whether the database should log queries")

	return &cmd.CobraCommand{Command: c}
}
