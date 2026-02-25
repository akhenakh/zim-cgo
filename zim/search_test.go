package zim

import (
	"testing"
)

func TestSearchAPI(t *testing.T) {
	archive, err := NewArchive(getTestSearchZimPath())
	if err != nil {
		t.Fatalf("Failed to open valid ZIM archive: %v", err)
	}
	defer archive.Close()

	if !archive.HasFulltextIndex() {
		t.Skip("Skipping Search API test because the test ZIM does not contain a fulltext index.")
	}

	searcher, err := NewSearcher(archive)
	if err != nil {
		t.Fatalf("Failed to create Searcher: %v", err)
	}
	defer searcher.Close()

	// Use the precise keyword based on the test data
	keyword := "markdown"
	query, err := NewQuery(keyword)
	if err != nil {
		t.Fatalf("Failed to create Query: %v", err)
	}
	defer query.Close()

	search, err := searcher.Search(query)
	if err != nil {
		t.Fatalf("Failed to execute Search: %v", err)
	}
	defer search.Close()

	matches := search.GetEstimatedMatches()
	t.Logf("Estimated matches for '%s': %d", keyword, matches)

	if matches == 0 {
		t.Fatalf("Expected at least 1 match for keyword '%s', got 0", keyword)
	}

	results, err := search.GetResults(0, 10)
	if err != nil {
		t.Fatalf("Failed to retrieve results: %v", err)
	}

	t.Logf("Retrieved %d actual search results.", len(results))

	if len(results) == 0 {
		t.Fatalf("Expected to retrieve at least 1 result, got 0")
	}

	// Validate the precision of the top result
	topResult := results[0]
	expectedPath := "index"
	expectedTitle := "Markdown Documentation"
	if topResult.Path != expectedPath {
		t.Errorf("Expected top result path to be %q, got %q", expectedPath, topResult.Path)
	}
	if topResult.Title != expectedTitle {
		t.Errorf("Expected top result title to be %q, got %q", expectedTitle, topResult.Title)
	}

	// Ensure the score is populated (usually out of 100 for the top result)
	if topResult.Score <= 0 {
		t.Errorf("Expected top result score to be > 0, got %d", topResult.Score)
	}

	for i, res := range results {
		t.Logf("Result %d: [%s] %s (Score: %d)", i, res.Path, res.Title, res.Score)
	}
}

func TestSuggestionAPI(t *testing.T) {
	archive, err := NewArchive(getTestSearchZimPath())
	if err != nil {
		t.Fatalf("Failed to open valid ZIM archive: %v", err)
	}
	defer archive.Close()

	searcher, err := NewSuggestionSearcher(archive)
	if err != nil {
		t.Fatalf("Failed to create SuggestionSearcher: %v", err)
	}
	defer searcher.Close()

	keyword := "markdown"
	search, err := searcher.Suggest(keyword)
	if err != nil {
		t.Fatalf("Failed to execute Suggestion: %v", err)
	}
	defer search.Close()

	matches := search.GetEstimatedMatches()
	t.Logf("Estimated matches for '%s': %d", keyword, matches)

	if matches == 0 {
		t.Fatalf("Expected at least 1 match for keyword '%s', got 0", keyword)
	}

	results, err := search.GetResults(0, 10)
	if err != nil {
		t.Fatalf("Failed to retrieve results: %v", err)
	}

	t.Logf("Retrieved %d actual suggestion results.", len(results))

	if len(results) == 0 {
		t.Fatalf("Expected to retrieve at least 1 result, got 0")
	}

	topResult := results[0]
	if topResult.Path == "" {
		t.Errorf("Expected top result path to be non-empty, got empty string")
	}
	if topResult.Title == "" {
		t.Errorf("Expected top result title to be non-empty, got empty string")
	}

	for i, res := range results {
		t.Logf("Result %d: [%s] %s (Snippet: %q)", i, res.Path, res.Title, res.Snippet)
	}
}
