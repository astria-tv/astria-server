package resolvers

import (
	"fmt"

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
	return fmt.Sprintf("/olaris/m/images/tmdb/w185%s", r.r.ProfilePath)
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
