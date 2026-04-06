package resolvers

import (
	"context"
	"fmt"

	"gitlab.com/olaris/olaris-server/metadata/auth"
	"gitlab.com/olaris/olaris-server/metadata/db"
)

// PersonResolver resolves a Person.
type PersonResolver struct {
	r db.Person
}

// Name returns the person's name.
func (r *PersonResolver) Name() string {
	return r.r.Name
}

// TmdbID returns the TMDB person ID.
func (r *PersonResolver) TmdbID() int32 {
	return int32(r.r.TmdbID)
}

// ProfilePath returns a URL for the profile image, served through the image cache.
func (r *PersonResolver) ProfilePath() string {
	if r.r.ProfilePath == "" {
		return ""
	}
	return fmt.Sprintf("/astria/m/images/tmdb/w185%s", r.r.ProfilePath)
}

// Biography returns the person's biography.
func (r *PersonResolver) Biography() string {
	return r.r.Biography
}

// Birthday returns the person's date of birth.
func (r *PersonResolver) Birthday() string {
	return r.r.Birthday
}

// Deathday returns the person's date of death.
func (r *PersonResolver) Deathday() string {
	return r.r.Deathday
}

// PlaceOfBirth returns the person's place of birth.
func (r *PersonResolver) PlaceOfBirth() string {
	return r.r.PlaceOfBirth
}

// CastRoles returns all cast roles this person has across movies and series.
func (r *PersonResolver) CastRoles(ctx context.Context) []*PersonCastRoleResolver {
	userID, _ := auth.UserID(ctx)
	roles := db.FindCastRolesForPerson(r.r.ID)
	var resolvers []*PersonCastRoleResolver
	for _, role := range roles {
		switch role.OwnerType {
		case "movies":
			movie, err := db.FindMovieByID(role.OwnerID)
			if err != nil {
				continue
			}
			ps, _ := db.FindPlayState(movie.UUID, userID)
			wlMap := db.IsOnWatchlistByUUIDs([]string{movie.UUID}, userID)
			onWL := wlMap[movie.UUID]
			resolvers = append(resolvers, &PersonCastRoleResolver{
				character: role.Character,
				media:     &MovieResolver{r: *movie, playState: ps, onWatchlist: &onWL},
			})
		case "series":
			series, err := db.FindSeries(role.OwnerID)
			if err != nil {
				continue
			}
			wlMap := db.IsOnWatchlistByUUIDs([]string{series.UUID}, userID)
			onWL := wlMap[series.UUID]
			resolvers = append(resolvers, &PersonCastRoleResolver{
				character: role.Character,
				media:     &SeriesResolver{r: *series, onWatchlist: &onWL},
			})
		}
	}
	return resolvers
}

// CastRoleResolver resolves a single cast role.
type CastRoleResolver struct {
	r db.CastRole
}

// Person returns the person who played this role.
func (r *CastRoleResolver) Person() *PersonResolver {
	return &PersonResolver{r: r.r.Person}
}

// Character returns the character name played.
func (r *CastRoleResolver) Character() string {
	return r.r.Character
}

// PersonCastRoleResolver resolves a cast role from the person's perspective.
type PersonCastRoleResolver struct {
	character string
	media     interface{}
}

// Character returns the character name played.
func (r *PersonCastRoleResolver) Character() string {
	return r.character
}

// Media returns the media union resolver.
func (r *PersonCastRoleResolver) Media() *MovieOrSeriesResolver {
	return &MovieOrSeriesResolver{r: r.media}
}

type personArgs struct {
	TmdbID int32
}

// Person looks up a person by TMDB ID.
func (r *Resolver) Person(ctx context.Context, args *personArgs) *PersonResolver {
	person, err := db.FindPersonByTmdbID(int(args.TmdbID))
	if err != nil {
		return nil
	}

	// Lazily fetch full profile details on first access.
	if person.Biography == "" {
		if err := r.env.MetadataRetrievalAgent.UpdatePersonMD(person, person.TmdbID); err == nil {
			db.UpdatePersonDetails(person)
		}
	}

	return &PersonResolver{r: *person}
}
