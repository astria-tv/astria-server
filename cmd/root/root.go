package root

import (
	"github.com/spf13/cobra"

	"github.com/astria-tv/astria-server/pkg/cmd"
)

func NewRootCommand() *cmd.CobraCommand {
	c := &cobra.Command{
		Use: "astria",
	}

	return &cmd.CobraCommand{Command: c}
}
