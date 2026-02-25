package zim

/*
#include <stdlib.h>
#include "zim_wrapper.h"
*/
import "C"
import (
	"errors"
	"runtime"
	"unsafe"
)

// HasFulltextIndex returns true if the archive has a built-in search index
func (a *Archive) HasFulltextIndex() bool {
	return bool(C.zim_archive_has_fulltext_index(a.ptr))
}

// Query represents a search query string
type Query struct {
	ptr C.zim_query_t
}

func NewQuery(queryStr string) (*Query, error) {
	cQuery := C.CString(queryStr)
	defer C.free(unsafe.Pointer(cQuery))

	ptr := C.zim_query_new(cQuery)
	if ptr == nil {
		return nil, errors.New("failed to create search query")
	}

	q := &Query{ptr: ptr}
	runtime.SetFinalizer(q, (*Query).Close)
	return q, nil
}

func (q *Query) Close() {
	if q.ptr != nil {
		C.zim_query_free(q.ptr)
		q.ptr = nil
	}
}

// Searcher performs fulltext search over a ZIM Archive
type Searcher struct {
	ptr C.zim_searcher_t
}

func NewSearcher(archive *Archive) (*Searcher, error) {
	ptr := C.zim_searcher_new(archive.ptr)
	if ptr == nil {
		return nil, errors.New("failed to initialize searcher (does the archive have a fulltext index?)")
	}

	s := &Searcher{ptr: ptr}
	runtime.SetFinalizer(s, (*Searcher).Close)
	return s, nil
}

func (s *Searcher) Close() {
	if s.ptr != nil {
		C.zim_searcher_free(s.ptr)
		s.ptr = nil
	}
}

// Search executes a query and holds the results
type Search struct {
	ptr C.zim_search_t
}

func (s *Searcher) Search(query *Query) (*Search, error) {
	ptr := C.zim_searcher_search(s.ptr, query.ptr)
	if ptr == nil {
		return nil, errors.New("search execution failed")
	}

	search := &Search{ptr: ptr}
	runtime.SetFinalizer(search, (*Search).Close)
	return search, nil
}

func (s *Search) Close() {
	if s.ptr != nil {
		C.zim_search_free(s.ptr)
		s.ptr = nil
	}
}

func (s *Search) GetEstimatedMatches() int {
	return int(C.zim_search_get_estimated_matches(s.ptr))
}

// SearchResult holds metadata for a single matched entry
type SearchResult struct {
	Path      string
	Title     string
	Snippet   string
	Score     int
	WordCount int
}

// GetResults fetches a slice of results, handling the C++ iterator safely in the background
func (s *Search) GetResults(start, maxResults int) ([]SearchResult, error) {
	setPtr := C.zim_search_get_results(s.ptr, C.int(start), C.int(maxResults))
	if setPtr == nil {
		return nil, errors.New("failed to retrieve search results")
	}
	defer C.zim_search_result_set_free(setPtr)

	beginIt := C.zim_search_result_set_begin(setPtr)
	endIt := C.zim_search_result_set_end(setPtr)
	defer C.zim_search_iterator_free(beginIt)
	defer C.zim_search_iterator_free(endIt)

	var results []SearchResult

	// Iterate through results
	for !bool(C.zim_search_iterator_equal(beginIt, endIt)) {
		cPath := C.zim_search_iterator_get_path(beginIt)
		cTitle := C.zim_search_iterator_get_title(beginIt)
		cSnippet := C.zim_search_iterator_get_snippet(beginIt)

		res := SearchResult{
			Path:      C.GoString(cPath),
			Title:     C.GoString(cTitle),
			Snippet:   C.GoString(cSnippet),
			Score:     int(C.zim_search_iterator_get_score(beginIt)),
			WordCount: int(C.zim_search_iterator_get_word_count(beginIt)),
		}

		C.free(unsafe.Pointer(cPath))
		C.free(unsafe.Pointer(cTitle))
		C.free(unsafe.Pointer(cSnippet))

		results = append(results, res)
		C.zim_search_iterator_next(beginIt)
	}

	return results, nil
}

// --- Suggestion API ---

// SuggestionSearcher provides suggestion search over titles in a ZIM Archive
type SuggestionSearcher struct {
	ptr C.zim_suggestion_searcher_t
}

func NewSuggestionSearcher(archive *Archive) (*SuggestionSearcher, error) {
	ptr := C.zim_suggestion_searcher_new(archive.ptr)
	if ptr == nil {
		return nil, errors.New("failed to create suggestion searcher")
	}

	s := &SuggestionSearcher{ptr: ptr}
	runtime.SetFinalizer(s, (*SuggestionSearcher).Close)
	return s, nil
}

func (s *SuggestionSearcher) Close() {
	if s.ptr != nil {
		C.zim_suggestion_searcher_free(s.ptr)
		s.ptr = nil
	}
}

func (s *SuggestionSearcher) SetVerbose(verbose bool) {
	C.zim_suggestion_searcher_set_verbose(s.ptr, C.bool(verbose))
}

// SuggestionSearch represents a suggestion search query
type SuggestionSearch struct {
	ptr C.zim_suggestion_search_t
}

func (s *SuggestionSearcher) Suggest(query string) (*SuggestionSearch, error) {
	cQuery := C.CString(query)
	defer C.free(unsafe.Pointer(cQuery))

	ptr := C.zim_suggestion_searcher_suggest(s.ptr, cQuery)
	if ptr == nil {
		return nil, errors.New("suggestion search failed")
	}

	search := &SuggestionSearch{ptr: ptr}
	runtime.SetFinalizer(search, (*SuggestionSearch).Close)
	return search, nil
}

func (s *SuggestionSearch) Close() {
	if s.ptr != nil {
		C.zim_suggestion_search_free(s.ptr)
		s.ptr = nil
	}
}

func (s *SuggestionSearch) GetEstimatedMatches() int {
	return int(C.zim_suggestion_search_get_estimated_matches(s.ptr))
}

// SuggestionResult holds metadata for a single suggestion
type SuggestionResult struct {
	Path    string
	Title   string
	Snippet string
}

// GetResults fetches a slice of suggestion results
func (s *SuggestionSearch) GetResults(start, maxResults int) ([]SuggestionResult, error) {
	setPtr := C.zim_suggestion_search_get_results(s.ptr, C.int(start), C.int(maxResults))
	if setPtr == nil {
		return nil, errors.New("failed to retrieve suggestion results")
	}
	defer C.zim_suggestion_result_set_free(setPtr)

	beginIt := C.zim_suggestion_result_set_begin(setPtr)
	endIt := C.zim_suggestion_result_set_end(setPtr)
	defer C.zim_suggestion_iterator_free(beginIt)
	defer C.zim_suggestion_iterator_free(endIt)

	var results []SuggestionResult

	for !bool(C.zim_suggestion_iterator_equal(beginIt, endIt)) {
		cPath := C.zim_suggestion_iterator_get_path(beginIt)
		cTitle := C.zim_suggestion_iterator_get_title(beginIt)
		cSnippet := C.zim_suggestion_iterator_get_snippet(beginIt)

		var path, title, snippet string
		if cPath != nil {
			path = C.GoString(cPath)
			C.free(unsafe.Pointer(cPath))
		}
		if cTitle != nil {
			title = C.GoString(cTitle)
			C.free(unsafe.Pointer(cTitle))
		}
		if cSnippet != nil {
			snippet = C.GoString(cSnippet)
			C.free(unsafe.Pointer(cSnippet))
		}

		results = append(results, SuggestionResult{
			Path:    path,
			Title:   title,
			Snippet: snippet,
		})

		C.zim_suggestion_iterator_next(beginIt)
	}

	return results, nil
}
