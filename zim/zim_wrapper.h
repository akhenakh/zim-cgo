#ifndef ZIM_WRAPPER_H
#define ZIM_WRAPPER_H

#include <stdint.h>
#include <stdbool.h>

#ifdef __cplusplus
extern "C" {
#endif

typedef void* zim_archive_t;
typedef void* zim_entry_t;
typedef void* zim_item_t;
typedef void* zim_query_t;
typedef void* zim_searcher_t;
typedef void* zim_search_t;
typedef void* zim_search_result_set_t;
typedef void* zim_search_iterator_t;
typedef void* zim_suggestion_searcher_t;
typedef void* zim_suggestion_search_t;
typedef void* zim_suggestion_result_set_t;
typedef void* zim_suggestion_iterator_t;
typedef void* zim_creator_t;
typedef void* zim_writer_item_t;

// Archive
zim_archive_t zim_archive_new(const char* path);
void zim_archive_free(zim_archive_t archive);
uint64_t zim_archive_get_entry_count(zim_archive_t archive);
bool zim_archive_has_entry_by_path(zim_archive_t archive, const char* path);
zim_entry_t zim_archive_get_entry_by_path(zim_archive_t archive, const char* path);
zim_entry_t zim_archive_get_main_entry(zim_archive_t archive);
zim_entry_t zim_archive_get_entry_by_index(zim_archive_t archive, uint32_t idx);
bool zim_archive_has_fulltext_index(zim_archive_t archive);

// Entry
void zim_entry_free(zim_entry_t entry);
bool zim_entry_is_redirect(zim_entry_t entry);
zim_item_t zim_entry_get_item(zim_entry_t entry, bool follow);

// Item
void zim_item_free(zim_item_t item);
char* zim_item_get_path(zim_item_t item);     // Caller must free()
char* zim_item_get_title(zim_item_t item);    // Caller must free()
char* zim_item_get_mimetype(zim_item_t item); // Caller must free()
uint64_t zim_item_get_size(zim_item_t item);
char* zim_item_get_data(zim_item_t item, uint64_t* size); // Caller must free()

// --- Search API ---
zim_query_t zim_query_new(const char* query_str);
void zim_query_free(zim_query_t query);

zim_searcher_t zim_searcher_new(zim_archive_t archive);
void zim_searcher_free(zim_searcher_t searcher);

zim_search_t zim_searcher_search(zim_searcher_t searcher, zim_query_t query);
void zim_search_free(zim_search_t search);

int zim_search_get_estimated_matches(zim_search_t search);

zim_search_result_set_t zim_search_get_results(zim_search_t search, int start, int max_results);
void zim_search_result_set_free(zim_search_result_set_t set);
int zim_search_result_set_get_size(zim_search_result_set_t set);

// --- Search Iterator ---
zim_search_iterator_t zim_search_result_set_begin(zim_search_result_set_t set);
zim_search_iterator_t zim_search_result_set_end(zim_search_result_set_t set);
void zim_search_iterator_free(zim_search_iterator_t it);
bool zim_search_iterator_equal(zim_search_iterator_t a, zim_search_iterator_t b);
void zim_search_iterator_next(zim_search_iterator_t it);

char* zim_search_iterator_get_path(zim_search_iterator_t it);
char* zim_search_iterator_get_title(zim_search_iterator_t it);
char* zim_search_iterator_get_snippet(zim_search_iterator_t it);
int zim_search_iterator_get_score(zim_search_iterator_t it);
int zim_search_iterator_get_word_count(zim_search_iterator_t it);

// --- Suggestion API ---
zim_suggestion_searcher_t zim_suggestion_searcher_new(zim_archive_t archive);
void zim_suggestion_searcher_free(zim_suggestion_searcher_t searcher);
void zim_suggestion_searcher_set_verbose(zim_suggestion_searcher_t searcher, bool verbose);

zim_suggestion_search_t zim_suggestion_searcher_suggest(zim_suggestion_searcher_t searcher, const char* query);
void zim_suggestion_search_free(zim_suggestion_search_t search);
int zim_suggestion_search_get_estimated_matches(zim_suggestion_search_t search);

zim_suggestion_result_set_t zim_suggestion_search_get_results(zim_suggestion_search_t search, int start, int max_results);
void zim_suggestion_result_set_free(zim_suggestion_result_set_t set);
int zim_suggestion_result_set_get_size(zim_suggestion_result_set_t set);

zim_suggestion_iterator_t zim_suggestion_result_set_begin(zim_suggestion_result_set_t set);
zim_suggestion_iterator_t zim_suggestion_result_set_end(zim_suggestion_result_set_t set);
void zim_suggestion_iterator_free(zim_suggestion_iterator_t it);
bool zim_suggestion_iterator_equal(zim_suggestion_iterator_t a, zim_suggestion_iterator_t b);
void zim_suggestion_iterator_next(zim_suggestion_iterator_t it);

char* zim_suggestion_iterator_get_path(zim_suggestion_iterator_t it);
char* zim_suggestion_iterator_get_title(zim_suggestion_iterator_t it);
char* zim_suggestion_iterator_get_snippet(zim_suggestion_iterator_t it);
bool zim_suggestion_iterator_has_snippet(zim_suggestion_iterator_t it);

// --- Writer API ---
typedef void* zim_creator_t;
typedef void* zim_writer_item_t;

zim_creator_t zim_creator_new();
void zim_creator_free(zim_creator_t creator);

void zim_creator_config_verbose(zim_creator_t creator, bool verbose);
void zim_creator_config_compression(zim_creator_t creator, int compression); // 1 = None, 5 = Zstd

bool zim_creator_start_zim_creation(zim_creator_t creator, const char* filepath);
bool zim_creator_add_item(zim_creator_t creator, zim_writer_item_t item);
bool zim_creator_add_metadata(zim_creator_t creator, const char* name, const char* content);
bool zim_creator_add_illustration(zim_creator_t creator, unsigned int size, const char* content, uint64_t content_len);
bool zim_creator_set_main_path(zim_creator_t creator, const char* main_path);
bool zim_creator_finish_zim_creation(zim_creator_t creator);

zim_writer_item_t zim_writer_string_item_new(const char* path, const char* mimetype, const char* title, const char* content, uint64_t content_len, bool front_article);
zim_writer_item_t zim_writer_file_item_new(const char* path, const char* mimetype, const char* title, const char* filepath, bool front_article);
void zim_writer_item_free(zim_writer_item_t item);

#ifdef __cplusplus
}
#endif

#endif // ZIM_WRAPPER_H
