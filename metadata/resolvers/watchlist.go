package resolvers

import (
	"context"

	"github.com/astria-tv/astria-server/metadata/auth"
	"github.com/astria-tv/astria-server/metadata/db"
)

type addToWatchlistArgs struct {
	MediaUUID string
	MediaType string
}

type removeFromWatchlistArgs struct {
	MediaUUID string
}

// Watchlist returns all movies and series on the current user's watchlist.
func (r *Resolver) Watchlist(ctx context.Context) []*MovieOrSeriesResolver {
	userID, _ := auth.UserID(ctx)
	items, _ := db.GetWatchlistItems(userID)

	var resolvers []*MovieOrSeriesResolver
	for _, item := range items {
		switch item.MediaType {
		case "movie":
			movie, err := db.FindMovieByUUID(item.MediaUUID)
			if err != nil {
				continue
			}
			resolvers = append(resolvers, &MovieOrSeriesResolver{r: &MovieResolver{r: *movie}})
		case "series":
			series, err := db.FindSeriesByUUID(item.MediaUUID)
			if err != nil {
				continue
			}
			resolvers = append(resolvers, &MovieOrSeriesResolver{r: &SeriesResolver{r: *series}})
		}
	}
	return resolvers
}

// AddToWatchlist adds a movie or series to the user's watchlist.
func (r *Resolver) AddToWatchlist(ctx context.Context, args *addToWatchlistArgs) *BoolResponseResolver {
	userID, _ := auth.UserID(ctx)
	item := db.WatchlistItem{
		UserID:    userID,
		MediaUUID: args.MediaUUID,
		MediaType: args.MediaType,
	}
	err := db.AddWatchlistItem(&item)
	return &BoolResponseResolver{success: err == nil}
}

// RemoveFromWatchlist removes a movie or series from the user's watchlist.
func (r *Resolver) RemoveFromWatchlist(ctx context.Context, args *removeFromWatchlistArgs) *BoolResponseResolver {
	userID, _ := auth.UserID(ctx)
	err := db.DeleteWatchlistItem(args.MediaUUID, userID)
	return &BoolResponseResolver{success: err == nil}
}
