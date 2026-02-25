package main

import (
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/akhenakh/kiwix-go/kiwix"
)

func main() {
	// Create a virtual Kiwix Library
	library := kiwix.NewLibrary()
	defer library.Close()

	// Wrap the Library in a Manager to add/edit ZIM books
	manager := kiwix.NewManager(library)
	defer manager.Close()

	// Scan a directory for ZIM files
	// (Ensure you have a local "zims" directory with some actual .zim files)
	zimPath := "./zims"
	fmt.Printf("Scanning directory '%s' for ZIM files...\n", zimPath)
	manager.AddBooksFromDirectory(zimPath)

	count := library.GetBookCount(true, false)
	if count == 0 {
		log.Println("No books found! Please place a .zim file in the 'zims' directory.")
	} else {
		fmt.Printf("Successfully loaded %d books into the library.\n", count)
	}

	// Create the Web Server
	server := kiwix.NewCPPServer(library)
	defer server.Close()

	port := 8080
	server.SetPort(port)
	server.SetBlockExternalLinks(true) // Great for offline-only compliance

	fmt.Printf("Starting Kiwix HTTP Server on port %d...\n", port)
	err := server.Start()
	if err != nil {
		log.Fatalf("Server error: %v", err)
	}

	fmt.Printf("Server is running! Navigate to http://localhost:%d in your browser.\n", port)

	// Block main goroutine until interrupted
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, os.Interrupt, syscall.SIGTERM)
	<-sigChan

	fmt.Println("\nShutting down server...")
	server.Stop()
}
