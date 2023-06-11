package library

import (
	"github.com/goava/di"
	"github.com/pkg/errors"
	"github.com/spf13/cobra"
	"gitlab.com/olaris/olaris-server/cmd/library/library_create"
)

func Options() di.Option {
	return di.Options(
		di.Provide(NewLibraryCommand, di.Tags{"type": "library"}),
		di.Invoke(RegisterLibraryCommand),
		library_create.Options(),
	)
}

func RegisterLibraryCommand(deps struct {
	di.Inject
	RootCommand    *cobra.Command `di:"type=root"`
	LibraryCommand *cobra.Command `di:"type=library"`
}) {
	deps.RootCommand.AddCommand(deps.LibraryCommand)
}

func NewLibraryCommand() *cobra.Command {
	c := &cobra.Command{
		Use:   "library",
		Short: "Manage libraries",
		RunE: func(cmd *cobra.Command, args []string) error {
			return errors.New("Subcommand required")
		},
	}

	return c
}
