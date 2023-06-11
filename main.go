package main

import (
	"github.com/spf13/cobra"
	"gitlab.com/olaris/olaris-server/cmd"
	"gitlab.com/olaris/olaris-server/helpers"
	"os"

	"github.com/goava/di"
	"github.com/sirupsen/logrus"
)

// Workaround for setting a default command.
// https://github.com/spf13/cobra/issues/823
func resolveDefaultSubcommand(c *cobra.Command, defaultCommand string) {
	if len(os.Args) < 2 {
		os.Args = append(os.Args, defaultCommand)
		return
	}

	subcommand := os.Args[1]

	exclusions := []string{
		"--help",
		"-h",
		"help",
		"completion",
	}

	if helpers.ElementExists(exclusions, subcommand) {
		return
	}

	for _, v := range c.Commands() {
		if v.Name() == subcommand {
			return
		}
	}

	os.Args = append([]string{os.Args[0], defaultCommand}, os.Args[1:]...)
}

func main() {
	commandContainer, err := cmd.NewContainer()
	if err != nil {
		logrus.WithError(err).Fatal("failed to create command container")
	}

	var rootCommand *cobra.Command
	err = commandContainer.Resolve(&rootCommand, di.Tags{"type": "root"})
	if err != nil {
		logrus.WithError(err).Fatal("failed to resolve root command")
	}

	resolveDefaultSubcommand(rootCommand, "serve")

	err = rootCommand.Execute()
	if err != nil {
		logrus.WithError(err).Fatal("command execution failed")
	}

	os.Exit(0)
}
