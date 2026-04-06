package cmd

import (
	"github.com/astria-tv/astria-server/cmd/dumpdebug"
	"github.com/astria-tv/astria-server/cmd/identify"
	"github.com/astria-tv/astria-server/cmd/identify_movie"
	"github.com/astria-tv/astria-server/cmd/library"
	"github.com/astria-tv/astria-server/cmd/library_create"
	"github.com/astria-tv/astria-server/cmd/root"
	"github.com/astria-tv/astria-server/cmd/serve"
	"github.com/astria-tv/astria-server/cmd/user"
	"github.com/astria-tv/astria-server/cmd/user_create"
	"github.com/astria-tv/astria-server/cmd/version"
	"github.com/astria-tv/astria-server/pkg/cmd"
	"github.com/astria-tv/astria-server/streaming"
)

func New() *cmd.CobraCommand {
	rootCmd := root.NewRootCommand()

	serveCmd := serve.NewServeCommand(streaming.NewStreamingController())
	rootCmd.GetCobraCommand().AddCommand(serveCmd.GetCobraCommand())
	rootCmd.GetCobraCommand().Flags().AddFlagSet(serveCmd.GetCobraCommand().Flags())
	rootCmd.GetCobraCommand().Run = serveCmd.GetCobraCommand().Run

	userCmd := user.NewUserCommand()
	rootCmd.GetCobraCommand().AddCommand(userCmd.GetCobraCommand())
	userCmd.GetCobraCommand().AddCommand(user_create.NewUserCreateCommand().GetCobraCommand())

	identifyCmd := identify.NewIdentifyCommand()
	rootCmd.GetCobraCommand().AddCommand(identifyCmd.GetCobraCommand())
	identifyCmd.GetCobraCommand().AddCommand(identify_movie.NewIdentifyMovieCommand().GetCobraCommand())

	libraryCmd := library.NewLibraryCommand()
	rootCmd.GetCobraCommand().AddCommand(libraryCmd.GetCobraCommand())
	libraryCmd.GetCobraCommand().AddCommand(library_create.NewLibraryCreateCommand().GetCobraCommand())

	rootCmd.GetCobraCommand().AddCommand(dumpdebug.NewDumpDebugCommand().GetCobraCommand())
	rootCmd.GetCobraCommand().AddCommand(version.NewVersionCommand().GetCobraCommand())

	return rootCmd
}
