// Package metadata implements metadata server features such as media indexing,
// media metadata lookup on external services and exposing this data via APIs.
package metadata

import (
	"net/http"

	"github.com/astria-tv/astria-server/helpers"
	"github.com/astria-tv/astria-server/metadata/app"
	"github.com/astria-tv/astria-server/metadata/resolvers"
	"github.com/gorilla/mux"
	"github.com/graph-gophers/graphql-transport-ws/graphqlws"

	"github.com/astria-tv/astria-server/metadata/auth"
)

// RegisterRoutes defines the handlers for metadata endpoints such as graphql and REST methods.
func RegisterRoutes(menv *app.MetadataContext, r *mux.Router) {
	imageManager := NewImageManager()

	schema, handler := resolvers.NewRelayHandler(menv)
	r.Handle("/query", auth.MiddleWare(graphqlws.NewHandlerFunc(schema, handler)))

	r.HandleFunc("/v1/auth", auth.UserHandler).Methods("POST")

	r.HandleFunc("/v1/version", versionHandler).Methods("GET")

	r.HandleFunc("/v1/user", auth.CreateUserHandler).Methods("POST")
	r.HandleFunc("/v1/user/setup", auth.ReadyForSetup)

	// Images need CORS so the Chromecast receiver can load artwork.
	imageRouter := r.PathPrefix("/images").Subrouter()
	imageRouter.Use(addImageCORSHeaders)
	imageRouter.HandleFunc("/{provider}/{size}/{id}", imageManager.HTTPHandler)
}

// addImageCORSHeaders adds CORS headers for the cast receiver on image responses.
func addImageCORSHeaders(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		origin := r.Header.Get("Origin")
		if origin == "https://cast.astria.tv" {
			w.Header().Set("Access-Control-Allow-Origin", "https://cast.astria.tv")
			w.Header().Set("Access-Control-Allow-Methods", "GET, HEAD, OPTIONS")
			w.Header().Set("Access-Control-Max-Age", "86400")

			if r.Method == http.MethodOptions {
				w.WriteHeader(http.StatusNoContent)
				return
			}
		}

		next.ServeHTTP(w, r)
	})
}

func versionHandler(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte(helpers.Version))
}
