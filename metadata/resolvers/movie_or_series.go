package resolvers

// MovieOrSeriesResolver resolves the MovieOrSeries union (Movie | Series).
type MovieOrSeriesResolver struct {
	r interface{}
}

// ToMovie tries to convert to Movie.
func (r *MovieOrSeriesResolver) ToMovie() (*MovieResolver, bool) {
	res, ok := r.r.(*MovieResolver)
	return res, ok
}

// ToSeries tries to convert to Series.
func (r *MovieOrSeriesResolver) ToSeries() (*SeriesResolver, bool) {
	res, ok := r.r.(*SeriesResolver)
	return res, ok
}
