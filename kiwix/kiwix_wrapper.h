#ifndef KIWIX_WRAPPER_H
#define KIWIX_WRAPPER_H

#include <stdint.h>
#include <stdbool.h>

#ifdef __cplusplus
extern "C" {
#endif

// Opaque types
typedef void* kiwix_library_t;
typedef void* kiwix_manager_t;
typedef void* kiwix_book_t;
typedef void* kiwix_server_t;

// --- Library ---
kiwix_library_t kiwix_library_new();
void kiwix_library_free(kiwix_library_t lib);
unsigned int kiwix_library_get_book_count(kiwix_library_t lib, bool local, bool remote);
kiwix_book_t kiwix_library_get_book_by_id(kiwix_library_t lib, const char* id);

// --- Manager ---
kiwix_manager_t kiwix_manager_new(kiwix_library_t lib);
void kiwix_manager_free(kiwix_manager_t mgr);
bool kiwix_manager_add_book_from_path(kiwix_manager_t mgr, const char* path);
void kiwix_manager_add_books_from_directory(kiwix_manager_t mgr, const char* dir_path);

// --- Book ---
void kiwix_book_free(kiwix_book_t book);
char* kiwix_book_get_id(kiwix_book_t book);
char* kiwix_book_get_title(kiwix_book_t book);
char* kiwix_book_get_description(kiwix_book_t book);
char* kiwix_book_get_path(kiwix_book_t book);
uint64_t kiwix_book_get_article_count(kiwix_book_t book);
uint64_t kiwix_book_get_size(kiwix_book_t book);

// --- Server ---
kiwix_server_t kiwix_server_new(kiwix_library_t lib);
void kiwix_server_free(kiwix_server_t server);
void kiwix_server_set_port(kiwix_server_t server, int port);
void kiwix_server_set_block_external_links(kiwix_server_t server, bool block);
bool kiwix_server_start(kiwix_server_t server);
void kiwix_server_stop(kiwix_server_t server);


#ifdef __cplusplus
}
#endif

#endif // KIWIX_WRAPPER_H
