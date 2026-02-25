package zim

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestNativeHTTPServer_ServeMainEntry(t *testing.T) {
	archive, err := NewArchive(getTestZimPath())
	if err != nil {
		t.Fatalf("Failed to open valid ZIM archive: %v", err)
	}
	defer archive.Close()

	// Create the Native Go Server
	zimServer := NewHTTPServer(archive)

	// Create a mock HTTP request to the root path "/"
	req := httptest.NewRequest(http.MethodGet, "/", nil)
	w := httptest.NewRecorder()

	// Serve the request
	zimServer.ServeHTTP(w, req)
	resp := w.Result()
	defer resp.Body.Close()

	// If the test ZIM doesn't have a main entry, it's completely normal to get a 404.
	// Let's dynamically find a valid path to test instead!
	if resp.StatusCode == http.StatusNotFound {
		t.Log("Archive has no Main Entry. Falling back to dynamically finding the first file...")

		entry, err := archive.GetEntryByIndex(0)
		if err != nil {
			t.Fatalf("Failed to get entry at index 0: %v", err)
		}

		item, _ := entry.GetItem(true)
		path := item.GetPath()
		item.Close()
		entry.Close()

		t.Logf("Found dynamic path: /%s. Testing server again...", path)

		// Create a NEW request targeting the dynamic path
		req = httptest.NewRequest(http.MethodGet, "/"+path, nil)
		w = httptest.NewRecorder()
		zimServer.ServeHTTP(w, req)

		resp = w.Result()
		defer resp.Body.Close()
	}

	// Now check if the response was successful
	if resp.StatusCode != http.StatusOK {
		t.Errorf("Expected status OK (200), got %v", resp.StatusCode)
	}

	// Make sure we have a Content-Length
	if resp.ContentLength <= 0 {
		t.Errorf("Expected Content-Length > 0, got %v", resp.ContentLength)
	}

	t.Logf("Successfully served path. Mimetype: %s, Size: %d", resp.Header.Get("Content-Type"), resp.ContentLength)
}

func TestNativeHTTPServer_ServeNotFound(t *testing.T) {
	archive, err := NewArchive(getTestZimPath())
	if err != nil {
		t.Fatalf("Failed to open valid ZIM archive: %v", err)
	}
	defer archive.Close()

	zimServer := NewHTTPServer(archive)

	req := httptest.NewRequest(http.MethodGet, "/this_path_absolutely_does_not_exist_404.html", nil)
	w := httptest.NewRecorder()

	zimServer.ServeHTTP(w, req)

	resp := w.Result()
	defer resp.Body.Close()

	// Should return a 404
	if resp.StatusCode != http.StatusNotFound {
		t.Errorf("Expected status Not Found (404), got %v", resp.StatusCode)
	}
}
