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

	// Partition UUIDs by media type
	var movieUUIDs, seriesUUIDs []string
	for _, item := range items {
		switch item.MediaType {
		case "movie":
			movieUUIDs = append(movieUUIDs, item.MediaUUID)
		case "series":
			seriesUUIDs = append(seriesUUIDs, item.MediaUUID)
		}
	}

	// Batch-load all movies and series in two queries
	moviesMap := db.FindMoviesByUUIDs(movieUUIDs)
	seriesMap := db.FindSeriesByUUIDs(seriesUUIDs)

	// Batch-load play states for movies
	allUUIDs := append(movieUUIDs, seriesUUIDs...)
	playStates := db.FindPlayStatesByUUIDs(allUUIDs, userID)

	// Batch-load unwatched episode counts for series
	var seriesIDs []uint
	for _, s := range seriesMap {
		seriesIDs = append(seriesIDs, s.ID)
	}
	unwatchedCounts := db.BatchUnwatchedEpisodesInSeriesCounts(seriesIDs, userID)

	// Assemble resolvers preserving watchlist order
	onWL := true
	var resolvers []*MovieOrSeriesResolver
	for _, item := range items {
		switch item.MediaType {
		case "movie":
			movie, ok := moviesMap[item.MediaUUID]
			if !ok {
				continue
			}
			resolvers = append(resolvers, &MovieOrSeriesResolver{
				r: &MovieResolver{r: movie, playState: playStates[movie.UUID], onWatchlist: &onWL},
			})
		case "series":
			series, ok := seriesMap[item.MediaUUID]
			if !ok {
				continue
			}
			count := int32(unwatchedCounts[series.ID])
			resolvers = append(resolvers, &MovieOrSeriesResolver{
				r: &SeriesResolver{r: series, unwatchedEpisodesCount: &count, onWatchlist: &onWL},
			})
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
