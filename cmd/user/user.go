package user

import (
	"errors"

	"github.com/spf13/cobra"

	"github.com/astria-tv/astria-server/pkg/cmd"
)

func NewUserCommand() *cmd.CobraCommand {
	c := &cobra.Command{
		Use:   "user",
		Short: "Manage users",
		RunE: func(cmd *cobra.Command, args []string) error {
			return errors.New("Subcommand required")
		},
	}

	return &cmd.CobraCommand{Command: c}
}
