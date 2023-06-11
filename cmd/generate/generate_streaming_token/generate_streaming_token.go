package generate_streaming_token

import (
	"fmt"
	"github.com/goava/di"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"gitlab.com/olaris/olaris-server/metadata/auth"
	"os"
)

func Options() di.Option {
	return di.Options(
		di.Provide(NewGenerateStreamingTokenCommand, di.Tags{"type": "generate_streaming_token"}),
		di.Invoke(RegisterGenerateStreamingTokenCommand),
	)
}

func RegisterGenerateStreamingTokenCommand(deps struct {
	di.Inject
	GenerateCommand               *cobra.Command `di:"type=generate"`
	GenerateStreamingTokenCommand *cobra.Command `di:"type=generate_streaming_token"`
}) {
	deps.GenerateCommand.AddCommand(deps.GenerateStreamingTokenCommand)
}

func NewGenerateStreamingTokenCommand() *cobra.Command {
	var filepath string

	c := &cobra.Command{
		Use:   "streamingtoken",
		Short: "Generates a streaming token",
		Long:  "Generates a streaming token for the given file path.",
		Run: func(cmd *cobra.Command, args []string) {
			token, err := auth.CreateStreamingJWT(0, filepath)
			if err != nil {
				log.WithError(err).WithField("filepath", filepath).Fatal("Failed to create streaming token")
				os.Exit(1)
			}
			fmt.Println(token)
		},
	}

	c.Flags().StringVar(&filepath, "filepath", "", "Filepath to generate a token for")
	_ = c.MarkFlagRequired("filepath")

	return c
}
