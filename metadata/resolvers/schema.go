package resolvers

import (
	"github.com/graph-gophers/graphql-go"
	"gitlab.com/bytesized/bytesized-streaming/metadata/db"
)

var SchemaTxt = `
	schema {
		query: Query
		mutation: Mutation
	}
	# The query type, represents all of the entry points into our object graph
	type Query {
		movies(): [Movie]!
		libraries(): [Library]!
		tvseries(): [TvSeries]!
		users(): [User]!
	}

	type Mutation {
		# Add a library to scan
		createLibrary(name: String!, file_path: String!, kind: Int!): LibRes!
		createUser(login: String!, password: String!, admin: Boolean!): CreateUserResponse!
	}

	interface LibRes {
		library: Library!
		error: Error
	}

	interface CreateUserResponse {
		user: User!
		error: Error
	}

	interface Error {
		message: String!
		hasError: Boolean!
	}

	interface User {
		login: String!
		admin: Boolean!
	}

	# A media library
	interface Library {
		# Library Type (0 - movies)
		kind: Int!
		# Human readable name of the Library
		name: String!
		# Path that this library manages
		file_path: String!
		movies: [Movie]!
		episodes: [Episode]!
	}

	interface TvSeries {
		name: String!
		overview: String!
		first_air_date: String!
		status: String!
		seasons: [Season]!
		backdrop_path: String!
		poster_path: String!
		tmdb_id: Int!
		type: String!
		uuid: String!
	}

	interface Season {
		name: String!
		overview: String!
		season_number: Int!
		air_date: String!
		poster_path: String!
		tmdb_id: Int!
		episodes: [Episode]!
		uuid: String!
	}

	interface Episode {
		name: String!
		overview: String!
		still_path: String!
		air_date: String!
		episode_number: String!
		tmdb_id: Int!
		# Filename
		file_name: String!
		# Absolute path to the filesystem
		file_path: String!
		uuid: String!
	}

	# A movie file
	interface Movie {
		# Title of the movie
		title: String!
		# Official Title
		original_title: String!
		# Filename
		file_name: String!
		# Absolute path to the filesystem
		file_path: String!
		# Release year
		year: String!
		# Library ID
		library_id: Int!
		# Short description of the movie
		overview: String!
		# IMDB ID
		imdb_id: String!
		# TMDB ID
		tmdb_id: Int!
		# ID to retrieve backdrop
		backdrop_path: String!
		# ID to retrieve poster
		poster_path: String!
		uuid: String!
	}
`

func InitSchema(ctx *db.MetadataContext) *graphql.Schema {
	Schema := graphql.MustParseSchema(SchemaTxt, &Resolver{ctx: ctx})
	return Schema
}
