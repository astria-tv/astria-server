package version

import (
	"fmt"

	"github.com/spf13/cobra"

	"github.com/astria-tv/astria-server/helpers"
	"github.com/astria-tv/astria-server/pkg/cmd"
)

func NewVersionCommand() *cmd.CobraCommand {
	c := &cobra.Command{
		Use:   "version",
		Short: "Displays the current astria-server version",
		Run: func(cmd *cobra.Command, args []string) {
			fmt.Println(helpers.Version)
		},
	}

	return &cmd.CobraCommand{Command: c}
}
