package resolvers

import (
	"context"

	"github.com/astria-tv/astria-server/filesystem"
	"github.com/astria-tv/astria-server/metadata/db"
)

// MediaItemResolver is a resolver around mediaFile types.
type MediaFileItemResolver struct {
	r interface{}
}

// ToMovieFile tries to convert media to MovieFile
func (r *MediaFileItemResolver) ToMovieFile() (*MovieFileResolver, bool) {
	res, ok := r.r.(*MovieFileResolver)
	return res, ok
}

// ToEpisodeFile tries to convert media to EpisodeFile
func (r *MediaFileItemResolver) ToEpisodeFile() (*EpisodeFileResolver, bool) {
	res, ok := r.r.(*EpisodeFileResolver)
	return res, ok
}

type translateArgs struct {
	FileLocator string
}

// GetMediaFileFromFileLocator translates a fileLocator to an actual movie/episode file.
func (r *Resolver) GetMediaFileFromFileLocator(ctx context.Context, args *translateArgs) (mi *MediaFileItemResolver) {
	fl, err := filesystem.ParseFileLocator(args.FileLocator)

	if err != nil {
		return mi
	}
	l, err := filesystem.GetNodeFromFileLocator(fl)
	if err != nil {
		return mi
	}

	// Try to find an Episode with this path first
	e, err := db.FindEpisodeFileByPath(l)
	if err == nil {
		return &MediaFileItemResolver{r: &EpisodeFileResolver{r: *e}}
	}

	// Try to find a movie otherwise.
	m, err := db.FindMovieFileByPath(l)
	if err == nil {
		return &MediaFileItemResolver{r: &MovieFileResolver{r: *m}}
	}

	return mi
}
