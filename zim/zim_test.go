package zim

import (
	"path/filepath"
	"testing"
)

func getTestZimPath() string {
	// Adjusts the path to point to ../testdata/test.zim
	return filepath.Join("..", "testdata", "test.zim")
}

func getTestSearchZimPath() string {
	return filepath.Join("..", "testdata", "devdocs_en_markdown_2026-01.zim")
}

func TestZIMArchive_OpenValid(t *testing.T) {
	archive, err := NewArchive(getTestZimPath())
	if err != nil {
		t.Fatalf("Failed to open valid ZIM archive: %v", err)
	}
	defer archive.Close()

	count := archive.GetEntryCount()
	t.Logf("Successfully opened ZIM archive. Entry count: %d", count)

	if count == 0 {
		t.Log("Warning: Archive opened successfully but has 0 entries (might be a tiny/empty test ZIM).")
	}
}

func TestZIMArchive_OpenInvalid(t *testing.T) {
	invalidPath := filepath.Join("..", "testdata", "does_not_exist.zim")
	_, err := NewArchive(invalidPath)
	if err == nil {
		t.Errorf("Expected an error when opening a non-existent ZIM archive, got nil")
	}
}

func TestZIMArchive_GetNonExistentEntry(t *testing.T) {
	archive, err := NewArchive(getTestZimPath())
	if err != nil {
		t.Fatalf("Failed to open valid ZIM archive: %v", err)
	}
	defer archive.Close()

	// Requesting a path that is virtually guaranteed not to exist
	_, err = archive.GetEntryByPath("this_path_does_not_exist_123456789.html")
	if err == nil {
		t.Errorf("Expected an error when fetching a non-existent entry, got nil")
	}
}
