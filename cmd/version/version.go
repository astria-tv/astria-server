package version

import (
	"fmt"
	"github.com/goava/di"
	"github.com/spf13/cobra"
	"gitlab.com/olaris/olaris-server/helpers"
)

func Options() di.Option {
	return di.Options(
		di.Provide(NewVersionCommand, di.Tags{"type": "version"}),
		di.Invoke(RegisterVersionCommand),
	)
}

func RegisterVersionCommand(deps struct {
	di.Inject
	RootCommand    *cobra.Command `di:"type=root"`
	VersionCommand *cobra.Command `di:"type=version"`
}) {
	deps.RootCommand.AddCommand(deps.VersionCommand)
}

func NewVersionCommand() *cobra.Command {
	c := &cobra.Command{
		Use:   "version",
		Short: "Displays the current olaris-server version",
		Run: func(cmd *cobra.Command, args []string) {
			fmt.Println(helpers.Version)
		},
	}

	return c
}
