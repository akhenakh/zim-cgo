package main

import (
	"fmt"
	"log"

	"github.com/akhenakh/zim-cgo/zim"
)

func main() {
	// Open the archive
	archive, err := zim.NewArchive("../../testdata/devdocs_en_markdown_2026-01.zim")
	if err != nil {
		log.Fatalf("Error opening archive: %v", err)
	}
	defer archive.Close()

	fmt.Printf("Successfully loaded archive! Total entries: %d\n", archive.GetEntryCount())

	entry, err := archive.GetMainEntry()
	if err != nil {
		log.Fatalf("Error getting entry: %v", err)
	}
	defer entry.Close()

	// Resolve the entry payload into an Item (following redirects)
	item, err := entry.GetItem(true)
	if err != nil {
		log.Fatalf("Error getting item: %v", err)
	}
	defer item.Close()

	// Read metadata and data
	fmt.Printf("Item Title: %s\n", item.GetTitle())
	fmt.Printf("Item Mimetype: %s\n", item.GetMimetype())
	fmt.Printf("Item Size: %d bytes\n", item.GetSize())

	// Read content
	data := item.GetData()
	fmt.Printf("Preview: %s\n", string(data[:50])) // Print first 50 chars
}
