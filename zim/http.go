package zim

import (
	"fmt"
	"net/http"
	"strings"
)

// HTTPServer is a native Go HTTP handler that serves files directly from a ZIM archive
type HTTPServer struct {
	Archive *Archive
}

// NewHTTPServer creates a new native Go HTTP server for a ZIM archive
func NewHTTPServer(archive *Archive) *HTTPServer {
	return &HTTPServer{
		Archive: archive,
	}
}

// ServeHTTP implements the standard Go http.Handler interface
func (s *HTTPServer) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	// Remove the leading slash to match ZIM internal paths (e.g., "index.html" or "A/index.html")
	path := strings.TrimPrefix(r.URL.Path, "/")

	var entry *Entry
	var err error

	if path == "" {
		// Serve the main/home page
		entry, err = s.Archive.GetMainEntry()
		if err != nil {
			// Fallback if no main entry is explicitly set
			entry, err = s.Archive.GetEntryByPath("index.html")
		}
	} else {
		// Serve a specific file
		entry, err = s.Archive.GetEntryByPath(path)
	}

	if err != nil {
		http.NotFound(w, r)
		return
	}
	defer entry.Close()

	// Follow redirects (e.g. if the entry is just a pointer to another entry)
	item, err := entry.GetItem(true)
	if err != nil {
		http.Error(w, "Internal Server Error: Failed to load item", http.StatusInternalServerError)
		return
	}
	defer item.Close()

	// Get data and metadata
	data := item.GetData()
	mimetype := item.GetMimetype()
	size := item.GetSize()

	// Set required HTTP headers
	if mimetype != "" {
		w.Header().Set("Content-Type", mimetype)
	}
	w.Header().Set("Content-Length", fmt.Sprintf("%d", size))

	// Ensure caching headers are set since ZIM content is static
	w.Header().Set("Cache-Control", "public, max-age=86400")

	// Write the payload to the HTTP response
	w.WriteHeader(http.StatusOK)
	w.Write(data)
}
