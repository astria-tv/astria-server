package user_create

import (
	"github.com/goava/di"
	"github.com/spf13/cobra"

	"gitlab.com/olaris/olaris-server/metadata/app"
	"gitlab.com/olaris/olaris-server/metadata/db"
)

func Options() di.Option {
	return di.Options(
		di.Provide(NewUserCreateCommand, di.Tags{"type": "user_create"}),
		di.Invoke(RegisterUserCreateCommand),
	)
}

func RegisterUserCreateCommand(deps struct {
	di.Inject
	UserCommand       *cobra.Command `di:"type=user"`
	UserCreateCommand *cobra.Command `di:"type=user_create"`
}) {
	deps.UserCommand.AddCommand(deps.UserCreateCommand)
}

func NewUserCreateCommand() *cobra.Command {
	var username string
	var password string
	var admin bool

	c := &cobra.Command{
		Use:   "create",
		Short: "Create a new user",
		RunE: func(cmd *cobra.Command, args []string) error {
			mctx := app.NewDefaultMDContext()
			defer mctx.Db.Close()

			_, err := db.CreateUser(username, password, admin)

			return err
		},
	}

	c.Flags().StringVar(&username, "username", "", "")
	c.MarkFlagRequired("username")

	c.Flags().StringVar(&password, "password", "", "")
	c.MarkFlagRequired("password")

	c.Flags().BoolVar(&admin, "admin", false, "Whether the new user should be an admin")

	return c
}
