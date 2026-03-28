package resolvers

// MovieOrEpisodeResolver resolves the MovieOrEpisode union (Movie | Episode).
type MovieOrEpisodeResolver struct {
	r interface{}
}

// ToMovie tries to convert media to Movie
func (r *MovieOrEpisodeResolver) ToMovie() (*MovieResolver, bool) {
	res, ok := r.r.(*MovieResolver)
	return res, ok
}

// ToEpisode tries to convert media to Episode
func (r *MovieOrEpisodeResolver) ToEpisode() (*EpisodeResolver, bool) {
	res, ok := r.r.(*EpisodeResolver)
	return res, ok
}
