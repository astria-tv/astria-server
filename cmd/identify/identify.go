package identify

import (
	"github.com/pkg/errors"
	"github.com/spf13/cobra"

	"gitlab.com/olaris/olaris-server/pkg/cmd"
)

func NewIdentifyCommand() *cmd.CobraCommand {
	c := &cobra.Command{
		Use:   "identify",
		Short: "Identify media files",
		RunE: func(cmd *cobra.Command, args []string) error {
			return errors.New("Subcommand required")
		},
	}

	return &cmd.CobraCommand{Command: c}
}
