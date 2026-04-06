// Package react is a handler for the astria-react application.
package react

import (
	"embed"
	"fmt"
	"io/fs"
	"net/http"
)

//go:embed build/*
var embedded embed.FS

// spaHandler serves static files from the embedded filesystem, falling back to
// index.html for paths that don't match a file. This allows client-side routing
// to work on page refresh.
type spaHandler struct {
	fs         http.FileSystem
	fileServer http.Handler
}

func (h *spaHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	// Try to open the requested file.
	f, err := h.fs.Open(r.URL.Path)
	if err != nil {
		// File doesn't exist — serve index.html so the SPA router can handle it.
		r.URL.Path = "/"
		h.fileServer.ServeHTTP(w, r)
		return
	}
	f.Close()

	// File exists — serve it normally.
	h.fileServer.ServeHTTP(w, r)
}

// GetHandler implements a handler that serves up the compiled astria-react code.
func GetHandler() http.Handler {
	embeddedFS, err := fs.Sub(embedded, "build")

	if err != nil {
		panic(fmt.Sprintf("Failed to read embedded react files: %s", err.Error()))
	}

	fsys := http.FS(embeddedFS)
	handler := &spaHandler{
		fs:         fsys,
		fileServer: http.FileServer(fsys),
	}

	return handler
}
