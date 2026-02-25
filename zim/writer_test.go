package zim

import (
	"os"
	"path/filepath"
	"testing"
)

func TestZIMCreator_RoundTrip(t *testing.T) {
	// 1. Setup a temporary directory for our generated ZIM file
	tempDir := t.TempDir()
	outPath := filepath.Join(tempDir, "output.zim")

	// 2. Initialize Creator
	creator, err := NewCreator()
	if err != nil {
		t.Fatalf("Failed to create ZIM Creator: %v", err)
	}
	defer creator.Close()

	// 3. Configure and Start Creation
	creator.ConfigCompression(CompressionNone) // Faster for tests
	if err := creator.StartZimCreation(outPath); err != nil {
		t.Fatalf("Failed to start creation: %v", err)
	}

	// 4. Add Metadata
	creator.AddMetadata("Title", "Go Roundtrip Test Archive")
	creator.AddMetadata("Description", "A ZIM file generated natively via Go bindings.")
	creator.AddMetadata("Language", "eng")
	creator.AddMetadata("Creator", "GoZim")
	creator.AddMetadata("Publisher", "GoZim")

	// 5. Build and Add Items
	mainPageContent := []byte("<html><body><h1>Welcome to GoZim!</h1></body></html>")
	mainItem, err := NewStringItem("index.html", "text/html", "Home Page", mainPageContent, true)
	if err != nil {
		t.Fatalf("Failed to create string item: %v", err)
	}
	defer mainItem.Close()

	if err := creator.AddItem(mainItem); err != nil {
		t.Fatalf("Failed to add main item: %v", err)
	}

	// Tell the creator which file is the landing page
	creator.SetMainPath("index.html")

	aboutPageContent := []byte("<html><body>About this generator...</body></html>")
	aboutItem, err := NewStringItem("about.html", "text/html", "About", aboutPageContent, true)
	if err != nil {
		t.Fatalf("Failed to create about item: %v", err)
	}
	defer aboutItem.Close()

	if err := creator.AddItem(aboutItem); err != nil {
		t.Fatalf("Failed to add about item: %v", err)
	}

	// 6. Finish Creation
	if err := creator.FinishZimCreation(); err != nil {
		t.Fatalf("Failed to finish ZIM creation: %v", err)
	}

	// Ensure file exists on disk
	info, err := os.Stat(outPath)
	if err != nil || info.Size() == 0 {
		t.Fatalf("Generated ZIM file does not exist or is empty")
	}
	t.Logf("Successfully generated ZIM archive at %s (Size: %d bytes)", outPath, info.Size())

	// ---------------------------------------------------------
	// 7. Verify the written archive using our Read API
	// ---------------------------------------------------------

	archive, err := NewArchive(outPath)
	if err != nil {
		t.Fatalf("Failed to open generated ZIM archive for reading: %v", err)
	}
	defer archive.Close()

	// Should have exactly 2 user entries (index.html, about.html)
	if archive.GetEntryCount() != 2 {
		t.Errorf("Expected 2 user entries, got %d", archive.GetEntryCount())
	}

	// Fetch the main entry
	entry, err := archive.GetMainEntry()
	if err != nil {
		t.Fatalf("Failed to retrieve Main Entry from generated ZIM: %v", err)
	}
	defer entry.Close()

	// Read its payload (MUST follow redirects because MainEntry is a redirect pointer!)
	item, err := entry.GetItem(true)
	if err != nil {
		t.Fatalf("Failed to get Main Entry item: %v", err)
	}
	defer item.Close()

	if item.GetPath() != "index.html" {
		t.Errorf("Expected Main Entry path to be 'index.html', got %q", item.GetPath())
	}

	readData := item.GetData()
	if string(readData) != string(mainPageContent) {
		t.Errorf("Content mismatch. Expected %q, got %q", string(mainPageContent), string(readData))
	}

	t.Log("Round-trip Creation -> Write -> Read was 100% successful!")
}
