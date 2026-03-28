package resolvers

import (
	"gitlab.com/olaris/olaris-server/metadata/db"
)

type searchArgs struct {
	Name string
}

// Search searches for media content.
func (r *Resolver) Search(args *searchArgs) *[]*MovieOrSeriesResolver {
	var l []*MovieOrSeriesResolver

	for _, movie := range db.SearchMovieByTitle(args.Name) {
		l = append(l, &MovieOrSeriesResolver{r: &MovieResolver{r: movie}})
	}
	for _, serie := range db.SearchSeriesByTitle(args.Name) {
		l = append(l, &MovieOrSeriesResolver{r: &SeriesResolver{r: serie}})
	}

	return &l
}
