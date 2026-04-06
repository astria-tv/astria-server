package resolvers

import (
	"context"
	"sort"

	"github.com/astria-tv/astria-server/metadata/auth"
	"github.com/astria-tv/astria-server/metadata/db"
)

type sortable interface {
	TimeStamp() int64
	UpdatedAtTimeStamp() int64
}

// ByCreationDate is a sortable type to sort by creation date.
type ByCreationDate []sortable

func (a ByCreationDate) Len() int           { return len(a) }
func (a ByCreationDate) Swap(i, j int)      { a[i], a[j] = a[j], a[i] }
func (a ByCreationDate) Less(i, j int) bool { return a[i].TimeStamp() > a[j].TimeStamp() }

// ByUpdatedAt is a sortable type to sort by updated_at date.
type ByUpdatedAt []sortable

func (a ByUpdatedAt) Len() int      { return len(a) }
func (a ByUpdatedAt) Swap(i, j int) { a[i], a[j] = a[j], a[i] }
func (a ByUpdatedAt) Less(i, j int) bool {
	return a[i].UpdatedAtTimeStamp() > a[j].UpdatedAtTimeStamp()
}

// RecentlyAdded returns recently added media content.
func (r *Resolver) RecentlyAdded(ctx context.Context) *[]*MovieOrEpisodeResolver {
	userID, _ := auth.UserID(ctx)
	sortables := []sortable{}

	movies := db.RecentlyAddedMovies(userID)
	episodes := db.RecentlyAddedEpisodes(userID)

	for _, movie := range movies {
		sortables = append(sortables, movie)
	}
	for _, ep := range episodes {
		sortables = append(sortables, ep)
	}
	sort.Sort(ByCreationDate(sortables))

	// Batch-load play states for all items
	uuids := make([]string, 0, len(movies)+len(episodes))
	for _, m := range movies {
		uuids = append(uuids, m.UUID)
	}
	for _, e := range episodes {
		uuids = append(uuids, e.UUID)
	}
	playStates := db.FindPlayStatesByUUIDs(uuids, userID)

	l := []*MovieOrEpisodeResolver{}
	for _, item := range sortables {
		if res, ok := item.(*db.Episode); ok {
			l = append(l, &MovieOrEpisodeResolver{r: &EpisodeResolver{r: *res, playState: playStates[res.UUID]}})
		}
		if res, ok := item.(*db.Movie); ok {
			l = append(l, &MovieOrEpisodeResolver{r: &MovieResolver{r: *res, playState: playStates[res.UUID]}})
		}
	}

	return &l
}
