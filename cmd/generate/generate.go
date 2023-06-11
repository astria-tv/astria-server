package generate

import (
	"github.com/goava/di"
	"github.com/pkg/errors"
	"github.com/spf13/cobra"
	"gitlab.com/olaris/olaris-server/cmd/generate/generate_login_token"
	"gitlab.com/olaris/olaris-server/cmd/generate/generate_streaming_token"
)

func Options() di.Option {
	return di.Options(
		di.Provide(NewGenerateCommand, di.Tags{"type": "generate"}),
		di.Invoke(RegisterGenerateCommand),
		generate_login_token.Options(),
		generate_streaming_token.Options(),
	)
}

func RegisterGenerateCommand(deps struct {
	di.Inject
	RootCommand     *cobra.Command `di:"type=root"`
	GenerateCommand *cobra.Command `di:"type=generate"`
}) {
	deps.RootCommand.AddCommand(deps.GenerateCommand)
}

func NewGenerateCommand() *cobra.Command {
	c := &cobra.Command{
		Use:   "generate",
		Short: "Generate tokens",
		RunE: func(cmd *cobra.Command, args []string) error {
			return errors.New("Subcommand required")
		},
	}

	return c
}
