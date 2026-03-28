package version

import (
	"fmt"

	"github.com/spf13/cobra"

	"gitlab.com/olaris/olaris-server/helpers"
	"gitlab.com/olaris/olaris-server/pkg/cmd"
)

func NewVersionCommand() *cmd.CobraCommand {
	c := &cobra.Command{
		Use:   "version",
		Short: "Displays the current olaris-server version",
		Run: func(cmd *cobra.Command, args []string) {
			fmt.Println(helpers.Version)
		},
	}

	return &cmd.CobraCommand{Command: c}
}
