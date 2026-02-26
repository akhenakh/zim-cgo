package main

import (
	"embed"
	"encoding/json"
	"flag"
	"fmt"
	"html/template"
	"log"
	"math/rand"
	"net/http"
	"strings"

	"github.com/akhenakh/kiwix-go/zim"
)

//go:embed templates/*
var templateFS embed.FS

type Server struct {
	archive      *zim.Archive
	searcher     *zim.Searcher
	suggSearcher *zim.SuggestionSearcher
	templates    *template.Template
	entryCount   uint64
}

func main() {
	zimPath := flag.String("z", "", "path to zim file")
	flag.Parse()

	if *zimPath == "" {
		log.Fatal("please provide a zim file path using -z flag")
	}

	archive, err := zim.NewArchive(*zimPath)
	if err != nil {
		log.Fatalf("failed to open zim archive: %v", err)
	}
	defer archive.Close()

	searcher, err := zim.NewSearcher(archive)
	if err != nil {
		log.Fatalf("failed to create searcher: %v", err)
	}
	defer searcher.Close()

	suggSearcher, err := zim.NewSuggestionSearcher(archive)
	if err != nil {
		log.Fatalf("failed to create suggestion searcher: %v", err)
	}
	defer suggSearcher.Close()

	tmpl, err := template.ParseFS(templateFS, "templates/*.html")
	if err != nil {
		log.Fatalf("failed to parse templates: %v", err)
	}

	s := &Server{
		archive:      archive,
		searcher:     searcher,
		suggSearcher: suggSearcher,
		templates:    tmpl,
		entryCount:   archive.GetEntryCount(),
	}

	http.HandleFunc("/content/", s.handleContent)
	http.HandleFunc("/search/suggestions", s.handleSuggestions)
	http.HandleFunc("/search/results", s.handleSearchResults)
	http.HandleFunc("/api/random", s.handleRandomAPI)
	http.HandleFunc("/api/main", s.handleMainEntry)
	http.HandleFunc("/random", s.handleRandom)
	http.HandleFunc("/", s.handleMain)

	log.Println("server starting on :8080")
	if err := http.ListenAndServe(":8080", nil); err != nil {
		log.Fatal(err)
	}
}

func (s *Server) handleMain(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path == "/" || r.URL.Path == "" {
		s.renderShell(w, "")
		return
	}
	s.renderShell(w, r.URL.Path)
}

func (s *Server) handleRandom(w http.ResponseWriter, r *http.Request) {
	idx := rand.Uint32() % uint32(s.entryCount)

	entry, err := s.archive.GetEntryByIndex(idx)
	if err != nil {
		http.Error(w, "failed to get random entry", http.StatusInternalServerError)
		return
	}
	defer entry.Close()

	item, err := entry.GetItem(true)
	if err != nil {
		http.Error(w, "failed to get random item", http.StatusInternalServerError)
		return
	}
	defer item.Close()

	path := item.GetPath()
	if path == "" {
		http.Error(w, "random entry has no path", http.StatusInternalServerError)
		return
	}

	http.Redirect(w, r, "/content/"+path, http.StatusFound)
}

func (s *Server) handleMainEntry(w http.ResponseWriter, r *http.Request) {
	entry, err := s.archive.GetMainEntry()
	if err != nil {
		http.Error(w, "failed to get main entry", http.StatusNotFound)
		return
	}
	defer entry.Close()

	item, err := entry.GetItem(true)
	if err != nil {
		http.Error(w, "failed to get main item", http.StatusNotFound)
		return
	}
	defer item.Close()

	path := item.GetPath()
	if path == "" {
		http.Error(w, "main entry has no path", http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{"path": path})
}

func (s *Server) handleRandomAPI(w http.ResponseWriter, r *http.Request) {
	var path string
	for range 100 {
		idx := rand.Uint32() % uint32(s.entryCount)
		entry, err := s.archive.GetEntryByIndex(idx)
		if err != nil {
			continue
		}
		defer entry.Close()

		item, err := entry.GetItem(true)
		if err != nil {
			continue
		}
		defer item.Close()

		mimetype := item.GetMimetype()
		if strings.HasPrefix(mimetype, "text/html") {
			path = item.GetPath()
			break
		}
	}

	if path == "" {
		http.Error(w, "no HTML page found", http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{"path": path})
}

func (s *Server) handleSuggestions(w http.ResponseWriter, r *http.Request) {
	query := r.URL.Query().Get("q")
	if len(query) < 2 {
		json.NewEncoder(w).Encode([]zim.SuggestionResult{})
		return
	}

	search, err := s.suggSearcher.Suggest(query)
	if err != nil {
		http.Error(w, "suggestion search failed", http.StatusInternalServerError)
		return
	}
	defer search.Close()

	results, err := search.GetResults(0, 10)
	if err != nil {
		http.Error(w, "failed to get suggestion results", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(results)
}

func (s *Server) handleSearchResults(w http.ResponseWriter, r *http.Request) {
	defer func() {
		if r := recover(); r != nil {
			http.Error(w, "search failed", http.StatusInternalServerError)
		}
	}()

	query := r.URL.Query().Get("q")
	if query == "" {
		s.renderShell(w, "")
		return
	}

	queryObj, err := zim.NewQuery(query)
	if err != nil {
		http.Error(w, "failed to create query", http.StatusInternalServerError)
		return
	}
	defer queryObj.Close()

	search, err := s.searcher.Search(queryObj)
	if err != nil {
		http.Error(w, "search failed", http.StatusInternalServerError)
		return
	}
	defer search.Close()

	results, err := search.GetResults(0, 50)
	if err != nil {
		http.Error(w, "failed to get search results", http.StatusInternalServerError)
		return
	}

	s.renderSearchResults(w, query, results)
}

func (s *Server) handleContent(w http.ResponseWriter, r *http.Request) {
	path := strings.TrimPrefix(r.URL.Path, "/content/")

	entry, err := s.archive.GetEntryByPath(path)
	if err != nil {
		http.NotFound(w, r)
		return
	}
	defer entry.Close()

	item, err := entry.GetItem(true)
	if err != nil {
		http.Error(w, "failed to get item", http.StatusInternalServerError)
		return
	}
	defer item.Close()

	data := item.GetData()
	mimetype := item.GetMimetype()

	if strings.HasPrefix(mimetype, "text/html") {
		baseURL := getBaseURL(path)
		baseTag := "<base href=\"" + baseURL + "\">"
		data = []byte(baseTag + string(data))
	}

	if mimetype != "" {
		w.Header().Set("Content-Type", mimetype)
	}
	w.Header().Set("Content-Length", fmt.Sprintf("%d", len(data)))
	w.WriteHeader(http.StatusOK)
	w.Write(data)
}

func (s *Server) renderShell(w http.ResponseWriter, iframeSrc string) {
	data := struct {
		IframeSrc string
	}{
		IframeSrc: iframeSrc,
	}

	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	if err := s.templates.ExecuteTemplate(w, "shell.html", data); err != nil {
		http.Error(w, "template execution failed", http.StatusInternalServerError)
	}
}

func getBaseURL(path string) string {
	dir := path[:strings.LastIndex(path, "/")+1]
	return "/content/" + dir
}

func (s *Server) renderSearchResults(w http.ResponseWriter, query string, results []zim.SearchResult) {
	type searchResult struct {
		Path      string
		Title     template.HTML
		Snippet   template.HTML
		Score     int
		WordCount int
	}
	sr := make([]searchResult, len(results))
	for i, r := range results {
		title := strings.ReplaceAll(r.Title, "&lt;i>", "<i>")
		title = strings.ReplaceAll(title, "&lt;/i>", "</i>")
		snippet := r.Snippet

		sr[i] = searchResult{
			Path:      r.Path,
			Title:     template.HTML(title),
			Snippet:   template.HTML(snippet),
			Score:     r.Score,
			WordCount: r.WordCount,
		}
	}

	data := struct {
		Query   string
		Results []searchResult
	}{
		Query:   query,
		Results: sr,
	}

	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	if err := s.templates.ExecuteTemplate(w, "search.html", data); err != nil {
		log.Printf("template error: %v", err)
		http.Error(w, "template execution failed", http.StatusInternalServerError)
	}
}
