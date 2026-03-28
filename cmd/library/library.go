package library

import (
	"github.com/pkg/errors"
	"github.com/spf13/cobra"

	"gitlab.com/olaris/olaris-server/pkg/cmd"
)

func NewLibraryCommand() *cmd.CobraCommand {
	c := &cobra.Command{
		Use:   "library",
		Short: "Manage libraries",
		RunE: func(cmd *cobra.Command, args []string) error {
			return errors.New("Subcommand required")
		},
	}

	return &cmd.CobraCommand{Command: c}
}
