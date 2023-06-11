package identify

import (
	"github.com/goava/di"
	"github.com/pkg/errors"
	"github.com/spf13/cobra"
	"gitlab.com/olaris/olaris-server/cmd/identify/identify_movie"
)

func Options() di.Option {
	return di.Options(
		di.Provide(NewIdentifyCommand, di.Tags{"type": "identify"}),
		di.Invoke(RegisterIdentifyCommand),
		identify_movie.Options(),
	)
}

func RegisterIdentifyCommand(deps struct {
	di.Inject
	RootCommand     *cobra.Command `di:"type=root"`
	IdentifyCommand *cobra.Command `di:"type=identify"`
}) {
	deps.RootCommand.AddCommand(deps.IdentifyCommand)
}

func NewIdentifyCommand() *cobra.Command {
	c := &cobra.Command{
		Use:   "identify",
		Short: "Identify media files",
		RunE: func(cmd *cobra.Command, args []string) error {
			return errors.New("Subcommand required")
		},
	}

	return c
}
