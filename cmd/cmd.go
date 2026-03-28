package cmd

import (
	"gitlab.com/olaris/olaris-server/cmd/dumpdebug"
	"gitlab.com/olaris/olaris-server/cmd/identify"
	"gitlab.com/olaris/olaris-server/cmd/identify_movie"
	"gitlab.com/olaris/olaris-server/cmd/library"
	"gitlab.com/olaris/olaris-server/cmd/library_create"
	"gitlab.com/olaris/olaris-server/cmd/root"
	"gitlab.com/olaris/olaris-server/cmd/serve"
	"gitlab.com/olaris/olaris-server/cmd/user"
	"gitlab.com/olaris/olaris-server/cmd/user_create"
	"gitlab.com/olaris/olaris-server/cmd/version"
	"gitlab.com/olaris/olaris-server/pkg/cmd"
	"gitlab.com/olaris/olaris-server/streaming"
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
