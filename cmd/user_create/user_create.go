package user_create

import (
	"github.com/spf13/cobra"

	"github.com/astria-tv/astria-server/metadata/app"
	"github.com/astria-tv/astria-server/metadata/db"
	"github.com/astria-tv/astria-server/pkg/cmd"
)

func NewUserCreateCommand() *cmd.CobraCommand {
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

	return &cmd.CobraCommand{Command: c}
}
