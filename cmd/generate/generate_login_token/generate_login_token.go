package generate_login_token

import (
	"fmt"
	"github.com/goava/di"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"gitlab.com/olaris/olaris-server/metadata/app"
	"gitlab.com/olaris/olaris-server/metadata/auth"
	"gitlab.com/olaris/olaris-server/metadata/db"
	"time"
)

func Options() di.Option {
	return di.Options(
		di.Provide(NewGenerateLoginTokenCommand, di.Tags{"type": "generate_login_token"}),
		di.Invoke(RegisterGenerateLoginTokenCommand),
	)
}

func RegisterGenerateLoginTokenCommand(deps struct {
	di.Inject
	GenerateCommand           *cobra.Command `di:"type=generate"`
	GenerateLoginTokenCommand *cobra.Command `di:"type=generate_login_token"`
}) {
	deps.GenerateCommand.AddCommand(deps.GenerateLoginTokenCommand)
}

func NewGenerateLoginTokenCommand() *cobra.Command {
	var username string

	c := &cobra.Command{
		Use:   "logintoken",
		Short: "Generates a login token",
		Long:  "Generates a login token that can be used to authenticate the given user.",
		Run: func(cmd *cobra.Command, args []string) {
			mctx := app.NewDefaultMDContext()
			defer mctx.Db.Close()

			user, err := db.FindUserByUsername(username)
			if err != nil {
				log.WithError(err).WithField("user", username).Fatal("Failed to find user")
			}

			// Create a token quasi-unlimited validity.
			jwt, err := auth.CreateMetadataJWT(user, 1000*24*time.Hour)
			if err != nil {
				log.WithError(err).Fatal("Failed to create login token")
			}

			fmt.Println("Bearer", jwt)
		},
	}

	c.Flags().StringVar(&username, "username", "", "User to generate a token for")
	_ = c.MarkFlagRequired("username")

	return c
}
