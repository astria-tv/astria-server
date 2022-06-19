package cmd

import (
	"github.com/goava/di"
	"gitlab.com/olaris/olaris-server/cmd/dumpdebug"
	"gitlab.com/olaris/olaris-server/cmd/generate"
	"gitlab.com/olaris/olaris-server/cmd/identify"
	"gitlab.com/olaris/olaris-server/cmd/library"
	"gitlab.com/olaris/olaris-server/cmd/root"
	"gitlab.com/olaris/olaris-server/cmd/serve"
	"gitlab.com/olaris/olaris-server/cmd/user"
	"gitlab.com/olaris/olaris-server/cmd/version"
	"gitlab.com/olaris/olaris-server/streaming"
)

// NewContainer returns a new dependency injection container for the
// command line.
func NewContainer() (*di.Container, error) {
	return di.New(
		streaming.Options(),

		// Commands
		root.Options(),
		user.Options(),
		serve.Options(),
		identify.Options(),
		library.Options(),
		dumpdebug.Options(),
		version.Options(),
		generate.Options(),
	)
}
