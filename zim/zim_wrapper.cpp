#include "zim_wrapper.h"
#include <zim/archive.h>
#include <zim/entry.h>
#include <zim/item.h>
#include <zim/blob.h>
#include <zim/search.h>
#include <zim/suggestion.h>
#include <zim/writer/creator.h>
#include <zim/writer/item.h>
#include <cstring>
#include <cstdlib>

using namespace zim;

extern "C" {

zim_archive_t zim_archive_new(const char* path) {
    try {
        return new Archive(path);
    } catch(...) {
        return nullptr;
    }
}

void zim_archive_free(zim_archive_t archive) {
    delete static_cast<Archive*>(archive);
}

uint64_t zim_archive_get_entry_count(zim_archive_t archive) {
    if (!archive) return 0;
    return static_cast<Archive*>(archive)->getEntryCount();
}

bool zim_archive_has_entry_by_path(zim_archive_t archive, const char* path) {
    if (!archive) return false;
    return static_cast<Archive*>(archive)->hasEntryByPath(path);
}

zim_entry_t zim_archive_get_entry_by_path(zim_archive_t archive, const char* path) {
    try {
        Archive* arch = static_cast<Archive*>(archive);
        Entry entry = arch->getEntryByPath(path);
        return new Entry(entry);
    } catch(...) {
        return nullptr;
    }
}

zim_entry_t zim_archive_get_main_entry(zim_archive_t archive) {
    try {
        Archive* arch = static_cast<Archive*>(archive);
        if (!arch->hasMainEntry()) {
            return nullptr;
        }
        Entry entry = arch->getMainEntry();
        return new Entry(entry);
    } catch(...) {
        return nullptr;
    }
}

zim_entry_t zim_archive_get_entry_by_index(zim_archive_t archive, uint32_t idx) {
    try {
        Archive* arch = static_cast<Archive*>(archive);
        // getEntryByPath(idx) gets the idx'th entry sorted by path
        Entry entry = arch->getEntryByPath(idx); 
        return new Entry(entry);
    } catch(...) {
        return nullptr;
    }
}

bool zim_archive_has_fulltext_index(zim_archive_t archive) {
    if(!archive) return false;
    return static_cast<Archive*>(archive)->hasFulltextIndex();
}

void zim_entry_free(zim_entry_t entry) {
    delete static_cast<Entry*>(entry);
}

bool zim_entry_is_redirect(zim_entry_t entry) {
    if (!entry) return false;
    return static_cast<Entry*>(entry)->isRedirect();
}

zim_item_t zim_entry_get_item(zim_entry_t entry, bool follow) {
    try {
        Entry* e = static_cast<Entry*>(entry);
        Item item = e->getItem(follow);
        return new Item(item);
    } catch(...) {
        return nullptr;
    }
}

void zim_item_free(zim_item_t item) {
    delete static_cast<Item*>(item);
}

// Helper for strings
static char* copy_string(const std::string& str) {
    char* copy = (char*)malloc(str.length() + 1);
    if (copy) strcpy(copy, str.c_str());
    return copy;
}

char* zim_item_get_path(zim_item_t item) {
    try { return copy_string(static_cast<Item*>(item)->getPath()); } 
    catch(...) { return nullptr; }
}

char* zim_item_get_title(zim_item_t item) {
    try { return copy_string(static_cast<Item*>(item)->getTitle()); } 
    catch(...) { return nullptr; }
}

char* zim_item_get_mimetype(zim_item_t item) {
    try { return copy_string(static_cast<Item*>(item)->getMimetype()); } 
    catch(...) { return nullptr; }
}

uint64_t zim_item_get_size(zim_item_t item) {
    try { return static_cast<Item*>(item)->getSize(); } 
    catch(...) { return 0; }
}

char* zim_item_get_data(zim_item_t item, uint64_t* size) {
    try {
        Blob b = static_cast<Item*>(item)->getData();
        *size = b.size();
        char* buf = (char*)malloc(b.size());
        if (buf) {
            memcpy(buf, b.data(), b.size());
        }
        return buf;
    } catch(...) {
        *size = 0;
        return nullptr;
    }
}


// --- Search API ---

zim_query_t zim_query_new(const char* query_str) {
    try { return new Query(query_str); } catch(...) { return nullptr; }
}
void zim_query_free(zim_query_t query) { delete static_cast<Query*>(query); }

zim_searcher_t zim_searcher_new(zim_archive_t archive) {
    try {
        Archive* arch = static_cast<Archive*>(archive);
        return new Searcher(*arch);
    } catch(...) { return nullptr; }
}
void zim_searcher_free(zim_searcher_t searcher) { delete static_cast<Searcher*>(searcher); }

zim_search_t zim_searcher_search(zim_searcher_t searcher, zim_query_t query) {
    try {
        Searcher* s = static_cast<Searcher*>(searcher);
        Query* q = static_cast<Query*>(query);
        return new Search(s->search(*q));
    } catch(...) { return nullptr; }
}
void zim_search_free(zim_search_t search) { delete static_cast<Search*>(search); }

int zim_search_get_estimated_matches(zim_search_t search) {
    if(!search) return 0;
    return static_cast<Search*>(search)->getEstimatedMatches();
}

zim_search_result_set_t zim_search_get_results(zim_search_t search, int start, int max_results) {
    try {
        Search* s = static_cast<Search*>(search);
        return new SearchResultSet(s->getResults(start, max_results));
    } catch(...) { return nullptr; }
}
void zim_search_result_set_free(zim_search_result_set_t set) { delete static_cast<SearchResultSet*>(set); }

int zim_search_result_set_get_size(zim_search_result_set_t set) {
    if(!set) return 0;
    return static_cast<SearchResultSet*>(set)->size();
}

zim_search_iterator_t zim_search_result_set_begin(zim_search_result_set_t set) {
    try { return new SearchIterator(static_cast<SearchResultSet*>(set)->begin()); } catch(...) { return nullptr; }
}
zim_search_iterator_t zim_search_result_set_end(zim_search_result_set_t set) {
    try { return new SearchIterator(static_cast<SearchResultSet*>(set)->end()); } catch(...) { return nullptr; }
}
void zim_search_iterator_free(zim_search_iterator_t it) { delete static_cast<SearchIterator*>(it); }

bool zim_search_iterator_equal(zim_search_iterator_t a, zim_search_iterator_t b) {
    if(!a || !b) return false;
    return *static_cast<SearchIterator*>(a) == *static_cast<SearchIterator*>(b);
}
void zim_search_iterator_next(zim_search_iterator_t it) {
    if(it) {
        auto& iter = *static_cast<SearchIterator*>(it);
        ++iter;
    }
}

char* zim_search_iterator_get_path(zim_search_iterator_t it) {
    try { return copy_string(static_cast<SearchIterator*>(it)->getPath()); } catch(...) { return nullptr; }
}
char* zim_search_iterator_get_title(zim_search_iterator_t it) {
    try { return copy_string(static_cast<SearchIterator*>(it)->getTitle()); } catch(...) { return nullptr; }
}
char* zim_search_iterator_get_snippet(zim_search_iterator_t it) {
    try { return copy_string(static_cast<SearchIterator*>(it)->getSnippet()); } catch(...) { return nullptr; }
}
int zim_search_iterator_get_score(zim_search_iterator_t it) {
    try { return static_cast<SearchIterator*>(it)->getScore(); } catch(...) { return 0; }
}
int zim_search_iterator_get_word_count(zim_search_iterator_t it) {
    try { return static_cast<SearchIterator*>(it)->getWordCount(); } catch(...) { return 0; }
}

// --- Suggestion API ---

zim_suggestion_searcher_t zim_suggestion_searcher_new(zim_archive_t archive) {
    try {
        Archive* arch = static_cast<Archive*>(archive);
        return new SuggestionSearcher(*arch);
    } catch(...) { return nullptr; }
}
void zim_suggestion_searcher_free(zim_suggestion_searcher_t searcher) { 
    delete static_cast<SuggestionSearcher*>(searcher); 
}

void zim_suggestion_searcher_set_verbose(zim_suggestion_searcher_t searcher, bool verbose) {
    if (searcher) {
        static_cast<SuggestionSearcher*>(searcher)->setVerbose(verbose);
    }
}

zim_suggestion_search_t zim_suggestion_searcher_suggest(zim_suggestion_searcher_t searcher, const char* query) {
    try {
        SuggestionSearcher* s = static_cast<SuggestionSearcher*>(searcher);
        return new SuggestionSearch(s->suggest(query));
    } catch(...) { return nullptr; }
}
void zim_suggestion_search_free(zim_suggestion_search_t search) { 
    delete static_cast<SuggestionSearch*>(search); 
}

int zim_suggestion_search_get_estimated_matches(zim_suggestion_search_t search) {
    if(!search) return 0;
    return static_cast<SuggestionSearch*>(search)->getEstimatedMatches();
}

zim_suggestion_result_set_t zim_suggestion_search_get_results(zim_suggestion_search_t search, int start, int max_results) {
    try {
        SuggestionSearch* s = static_cast<SuggestionSearch*>(search);
        return new SuggestionResultSet(s->getResults(start, max_results));
    } catch(...) { return nullptr; }
}
void zim_suggestion_result_set_free(zim_suggestion_result_set_t set) { 
    delete static_cast<SuggestionResultSet*>(set); 
}

int zim_suggestion_result_set_get_size(zim_suggestion_result_set_t set) {
    if(!set) return 0;
    return static_cast<SuggestionResultSet*>(set)->size();
}

zim_suggestion_iterator_t zim_suggestion_result_set_begin(zim_suggestion_result_set_t set) {
    try { return new SuggestionIterator(static_cast<SuggestionResultSet*>(set)->begin()); } catch(...) { return nullptr; }
}
zim_suggestion_iterator_t zim_suggestion_result_set_end(zim_suggestion_result_set_t set) {
    try { return new SuggestionIterator(static_cast<SuggestionResultSet*>(set)->end()); } catch(...) { return nullptr; }
}
void zim_suggestion_iterator_free(zim_suggestion_iterator_t it) { 
    delete static_cast<SuggestionIterator*>(it); 
}

bool zim_suggestion_iterator_equal(zim_suggestion_iterator_t a, zim_suggestion_iterator_t b) {
    if(!a && !b) return true;
    if(!a || !b) return false;
    return *static_cast<SuggestionIterator*>(a) == *static_cast<SuggestionIterator*>(b);
}
void zim_suggestion_iterator_next(zim_suggestion_iterator_t it) {
    if(it) {
        auto& iter = *static_cast<SuggestionIterator*>(it);
        ++iter;
    }
}

char* zim_suggestion_iterator_get_path(zim_suggestion_iterator_t it) {
    try { 
        SuggestionIterator* iter = static_cast<SuggestionIterator*>(it);
        std::string result = iter->operator*().getPath();
        return copy_string(result); 
    } catch(...) { return nullptr; }
}
char* zim_suggestion_iterator_get_title(zim_suggestion_iterator_t it) {
    try { 
        SuggestionIterator* iter = static_cast<SuggestionIterator*>(it);
        std::string result = iter->operator*().getTitle();
        return copy_string(result); 
    } catch(...) { return nullptr; }
}
char* zim_suggestion_iterator_get_snippet(zim_suggestion_iterator_t it) {
    try { 
        SuggestionIterator* iter = static_cast<SuggestionIterator*>(it);
        std::string result = iter->operator*().getSnippet();
        return copy_string(result); 
    } catch(...) { return nullptr; }
}
bool zim_suggestion_iterator_has_snippet(zim_suggestion_iterator_t it) {
    try { 
        SuggestionIterator* iter = static_cast<SuggestionIterator*>(it);
        return iter->operator*().hasSnippet(); 
    } catch(...) { return false; }
}

// --- Writer API ---

zim_creator_t zim_creator_new() {
    try { return new zim::writer::Creator(); } catch(...) { return nullptr; }
}

void zim_creator_free(zim_creator_t creator) {
    delete static_cast<zim::writer::Creator*>(creator);
}

void zim_creator_config_verbose(zim_creator_t creator, bool verbose) {
    if (creator) static_cast<zim::writer::Creator*>(creator)->configVerbose(verbose);
}

void zim_creator_config_compression(zim_creator_t creator, int compression) {
    if (creator) static_cast<zim::writer::Creator*>(creator)->configCompression(static_cast<zim::Compression>(compression));
}

bool zim_creator_start_zim_creation(zim_creator_t creator, const char* filepath) {
    try {
        static_cast<zim::writer::Creator*>(creator)->startZimCreation(filepath);
        return true;
    } catch(...) { return false; }
}

bool zim_creator_add_item(zim_creator_t creator, zim_writer_item_t item) {
    try {
        auto shared_item = *static_cast<std::shared_ptr<zim::writer::Item>*>(item);
        static_cast<zim::writer::Creator*>(creator)->addItem(shared_item);
        return true;
    } catch(...) { return false; }
}

bool zim_creator_add_metadata(zim_creator_t creator, const char* name, const char* content) {
    try {
        static_cast<zim::writer::Creator*>(creator)->addMetadata(name, content);
        return true;
    } catch(...) { return false; }
}

bool zim_creator_add_illustration(zim_creator_t creator, unsigned int size, const char* content, uint64_t content_len) {
    try {
        std::string s_content(content, content_len);
        static_cast<zim::writer::Creator*>(creator)->addIllustration(size, s_content);
        return true;
    } catch(...) { return false; }
}

bool zim_creator_set_main_path(zim_creator_t creator, const char* main_path) {
    try {
        static_cast<zim::writer::Creator*>(creator)->setMainPath(main_path);
        return true;
    } catch(...) { return false; }
}

bool zim_creator_finish_zim_creation(zim_creator_t creator) {
    try {
        static_cast<zim::writer::Creator*>(creator)->finishZimCreation();
        return true;
    } catch(...) { return false; }
}

zim_writer_item_t zim_writer_string_item_new(const char* path, const char* mimetype, const char* title, const char* content, uint64_t content_len, bool front_article) {
    try {
        zim::writer::Hints hints;
        if (front_article) hints[zim::writer::FRONT_ARTICLE] = 1;
        std::string s_content(content, content_len);
        
        auto item = zim::writer::StringItem::create(path, mimetype, title, hints, s_content);
        return new std::shared_ptr<zim::writer::Item>(item);
    } catch(...) { return nullptr; }
}

zim_writer_item_t zim_writer_file_item_new(const char* path, const char* mimetype, const char* title, const char* filepath, bool front_article) {
    try {
        zim::writer::Hints hints;
        if (front_article) hints[zim::writer::FRONT_ARTICLE] = 1;
        
        auto item = std::make_shared<zim::writer::FileItem>(path, mimetype, title, hints, filepath);
        return new std::shared_ptr<zim::writer::Item>(item);
    } catch(...) { return nullptr; }
}

void zim_writer_item_free(zim_writer_item_t item) {
    if (item) delete static_cast<std::shared_ptr<zim::writer::Item>*>(item);
}



}
