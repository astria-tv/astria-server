// Package resolvers implements resolvers for the GraphQL interface.
package resolvers

import (
	"context"
	"sort"

	"gitlab.com/olaris/olaris-server/metadata/auth"
	"gitlab.com/olaris/olaris-server/metadata/db"
)

// UpNext returns episode/movie that could populate a dashboard.
func (r *Resolver) UpNext(ctx context.Context) *[]*MediaItemResolver {
	userID, _ := auth.UserID(ctx)
	sortables := []sortable{}

	movies := db.UpNextMovies(userID)
	episodes := db.UpNextEpisodes(userID)

	for _, movie := range movies {
		sortables = append(sortables, movie)
	}
	for _, ep := range episodes {
		sortables = append(sortables, ep)
	}
	sort.Sort(ByUpdatedAt(sortables))

	// Batch-load play states for all items
	uuids := make([]string, 0, len(movies)+len(episodes))
	for _, m := range movies {
		uuids = append(uuids, m.UUID)
	}
	for _, e := range episodes {
		uuids = append(uuids, e.UUID)
	}
	playStates := db.FindPlayStatesByUUIDs(uuids, userID)

	l := []*MediaItemResolver{}
	for _, item := range sortables {
		if res, ok := item.(*db.Episode); ok {
			l = append(l, &MediaItemResolver{r: &EpisodeResolver{r: *res, playState: playStates[res.UUID]}})
		}
		if res, ok := item.(*db.Movie); ok {
			l = append(l, &MediaItemResolver{r: &MovieResolver{r: *res, playState: playStates[res.UUID]}})
		}
	}

	return &l
}
