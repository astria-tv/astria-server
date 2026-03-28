package root

import (
	"github.com/spf13/cobra"

	"gitlab.com/olaris/olaris-server/pkg/cmd"
)

func NewRootCommand() *cmd.CobraCommand {
	c := &cobra.Command{
		Use: "olaris",
	}

	return &cmd.CobraCommand{Command: c}
}
