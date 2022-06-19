package user

import (
	"errors"
	"gitlab.com/olaris/olaris-server/cmd/user/user_create"

	"github.com/goava/di"
	"github.com/spf13/cobra"
)

func Options() di.Option {
	return di.Options(
		di.Provide(NewUserCommand, di.Tags{"type": "user"}),
		di.Invoke(RegisterUserCommand),
		user_create.Options(),
	)
}

func RegisterUserCommand(deps struct {
	di.Inject
	RootCommand    *cobra.Command `di:"type=root"`
	NewUserCommand *cobra.Command `di:"type=user"`
}) {
	deps.RootCommand.AddCommand(deps.NewUserCommand)
}

func NewUserCommand() *cobra.Command {
	c := &cobra.Command{
		Use:   "user",
		Short: "Manage users",
		RunE: func(cmd *cobra.Command, args []string) error {
			return errors.New("Subcommand required")
		},
	}

	return c
}
