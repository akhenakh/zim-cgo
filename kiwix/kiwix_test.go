package kiwix

import (
	"path/filepath"
	"testing"
	"time"
)

func getTestDataDir() string {
	return filepath.Join("..", "testdata")
}

func getTestZimPath() string {
	return filepath.Join(getTestDataDir(), "test.zim")
}

func TestKiwixLibraryAndManager_SingleFile(t *testing.T) {
	lib := NewLibrary()
	if lib == nil {
		t.Fatalf("Failed to create new Library")
	}
	defer lib.Close()

	mgr := NewManager(lib)
	if mgr == nil {
		t.Fatalf("Failed to create new Manager")
	}
	defer mgr.Close()

	// Test adding a single book directly
	success := mgr.AddBookFromPath(getTestZimPath())
	if !success {
		t.Errorf("Expected AddBookFromPath to return true for a valid ZIM")
	}

	// Verify the book was added to the local library
	count := lib.GetBookCount(true, false)
	if count != 1 {
		t.Errorf("Expected library to have 1 local book, got %d", count)
	}
}

func TestKiwixLibraryAndManager_DirectoryScan(t *testing.T) {
	lib := NewLibrary()
	defer lib.Close()

	mgr := NewManager(lib)
	defer mgr.Close()

	// Test adding from a directory
	mgr.AddBooksFromDirectory(getTestDataDir())

	// Verify the directory scan found the test.zim file
	count := lib.GetBookCount(true, false)
	if count < 1 {
		t.Errorf("Expected library to find at least 1 book in testdata directory, got %d", count)
	}
}

func TestKiwixCPPServer_Lifecycle(t *testing.T) {
	lib := NewLibrary()
	defer lib.Close()

	mgr := NewManager(lib)
	defer mgr.Close()
	mgr.AddBookFromPath(getTestZimPath())

	server := NewCPPServer(lib)
	if server == nil {
		t.Fatalf("Failed to create new CPPServer")
	}
	defer server.Close()

	// Configure server
	server.SetPort(18080)
	server.SetBlockExternalLinks(true)

	errCh := make(chan error, 1)
	go func() {
		errCh <- server.Start()
	}()

	time.Sleep(200 * time.Millisecond)

	select {
	case err := <-errCh:
		if err != nil {
			t.Fatalf("CPPServer failed to start: %v", err)
		}
	default:
		t.Log("CPPServer started successfully. Shutting down...")
		server.Stop()
	}
}
