package resolvers

import (
	"context"
	"fmt"
	"github.com/astria-tv/astria-server/metadata/auth"
)

// CreateNoAuthorisationError returns a standard error for unauthorised requests.
func CreateNoAuthorisationError() error {
	return fmt.Errorf("you are not authorised for this action")
}

func ifAdmin(ctx context.Context) error {
	admin, ok := auth.UserAdmin(ctx)
	if ok && admin {
		return nil
	}
	return CreateNoAuthorisationError()
}
