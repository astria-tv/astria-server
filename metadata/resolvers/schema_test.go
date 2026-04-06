package resolvers

import (
	"github.com/astria-tv/astria-server/metadata/app"
	"testing"
)

func TestInitSchema(t *testing.T) {
	env := app.NewTestingMDContext(nil)
	InitSchema(env)
}
