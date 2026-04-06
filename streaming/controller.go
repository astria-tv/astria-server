package streaming

import (
	"net/http"

	"gitlab.com/olaris/olaris-server/interfaces/web"

	"github.com/gorilla/mux"
	"github.com/graph-gophers/graphql-transport-ws/graphqlws"
	"gitlab.com/olaris/olaris-server/metadata/auth"
)

// AddCORSHeaders is a middleware that adds permissive CORS headers to all
// streaming responses. This is necessary for Chromecast and other cross-origin
// playback scenarios where the cast receiver is on a different domain.
func AddCORSHeaders(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		origin := r.Header.Get("Origin")
		if origin != "" {
			w.Header().Set("Access-Control-Allow-Origin", origin)
		} else {
			w.Header().Set("Access-Control-Allow-Origin", "*")
		}
		w.Header().Set("Access-Control-Allow-Methods", "GET, HEAD, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "*")
		w.Header().Set("Access-Control-Expose-Headers", "Content-Length, Content-Range, Content-Type, Accept-Ranges")
		w.Header().Set("Access-Control-Max-Age", "86400")

		if r.Method == http.MethodOptions {
			w.WriteHeader(http.StatusNoContent)
			return
		}

		next.ServeHTTP(w, r)
	})
}

// For Apple devices to handle HLS properly, the m3u8 playlists must be sent with the correct Content-Type
func AddM3U8Header(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/x-mpegURL")
		next.ServeHTTP(w, r)
	})
}

// Controller is the controller for the streaming package. It implements the
// web.Controller interface.
type Controller struct{}

// NewStreamingController creates and returns a new StreamingHandler.
func NewStreamingController() web.Controller {
	return &Controller{}
}

// RegisterRoutes registers the streaming handler's routes on the provided
// router.
func (sh *Controller) RegisterRoutes(router *mux.Router) {
	router.Use(AddCORSHeaders)

	router.Handle("/files/{fileLocator:.*}/{sessionID}/hls-transmuxing-manifest.m3u8", AddM3U8Header(http.HandlerFunc(serveHlsTransmuxingMasterPlaylist)))
	router.Handle("/files/{fileLocator:.*}/{sessionID}/hls-transcoding-manifest.m3u8", AddM3U8Header(http.HandlerFunc(serveHlsTranscodingMasterPlaylist)))
	router.HandleFunc("/files/{fileLocator:.*}/metadata.json", serveMetadata)
	router.Handle("/files/{fileLocator:.*}/{sessionID}/hls-manifest.m3u8", AddM3U8Header(http.HandlerFunc(serveHlsMasterPlaylist)))
	router.HandleFunc("/files/{fileLocator:.*}/{sessionID}/dash-manifest.mpd", serveDASHManifest)
	router.Handle("/files/{fileLocator:.*}/{sessionID}/{streamId}/{representationId}/media.m3u8", AddM3U8Header(http.HandlerFunc(serveHlsTranscodingMediaPlaylist)))
	router.HandleFunc("/files/{fileLocator:.*}/{sessionID}/{streamId}/{representationId}/{segmentId:[0-9]+}.m4s", buildMediaSegmentHandlerFunc())
	router.HandleFunc("/files/{fileLocator:.*}/{sessionID}/{streamId}/{representationId}/{segmentId:[0-9]+}.vtt", buildMediaSegmentHandlerFunc())
	router.HandleFunc("/files/{fileLocator:.*}/{sessionID}/{streamId}/{representationId}/init.mp4", buildInitHandlerFunc())

	// This handler just serves up the file for downloading. This is also used
	// internally by ffmpeg to access rclone files.
	router.HandleFunc("/files/{fileLocator:.*}", serveFile)

	router.HandleFunc("/debug/playbackSessions", servePlaybackSessionDebugPage)

	schema, handler := NewRelayHandler()
	router.Handle("/query", auth.MiddleWare(graphqlws.NewHandlerFunc(schema, handler)))
}
