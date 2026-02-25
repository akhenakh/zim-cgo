#include "kiwix_wrapper.h"
#include <kiwix/library.h>
#include <kiwix/manager.h>
#include <kiwix/book.h>
#include <kiwix/server.h>
#include <cstring>
#include <cstdlib>
#include <filesystem> 

using namespace kiwix;

extern "C" {

// Helper to pass strings back to Go
static char* copy_string(const std::string& str) {
    char* copy = (char*)malloc(str.length() + 1);
    if (copy) strcpy(copy, str.c_str());
    return copy;
}

// --- Library ---

kiwix_library_t kiwix_library_new() {
    try {
        auto shared_lib = Library::create();
        // Allocate a pointer to the shared_ptr so we don't lose the reference block
        return new std::shared_ptr<Library>(shared_lib);
    } catch(...) { return nullptr; }
}

void kiwix_library_free(kiwix_library_t lib) {
    if (lib) {
        delete static_cast<std::shared_ptr<Library>*>(lib);
    }
}

unsigned int kiwix_library_get_book_count(kiwix_library_t lib, bool local, bool remote) {
    if (!lib) return 0;
    auto s_lib = *static_cast<std::shared_ptr<Library>*>(lib);
    return s_lib->getBookCount(local, remote);
}

kiwix_book_t kiwix_library_get_book_by_id(kiwix_library_t lib, const char* id) {
    try {
        auto s_lib = *static_cast<std::shared_ptr<Library>*>(lib);
        // Use the thread-safe version to return a copy of the Book object
        Book book = s_lib->getBookByIdThreadSafe(id);
        return new Book(book);
    } catch(...) { return nullptr; }
}

// --- Manager ---

kiwix_manager_t kiwix_manager_new(kiwix_library_t lib) {
    try {
        auto s_lib = *static_cast<std::shared_ptr<Library>*>(lib);
        return new Manager(s_lib);
    } catch(...) { return nullptr; }
}

void kiwix_manager_free(kiwix_manager_t mgr) {
    delete static_cast<Manager*>(mgr);
}

bool kiwix_manager_add_book_from_path(kiwix_manager_t mgr, const char* path) {
    try {
        return static_cast<Manager*>(mgr)->addBookFromPath(path);
    } catch(...) { return false; }
}

void kiwix_manager_add_books_from_directory(kiwix_manager_t mgr, const char* dir_path) {
    try {
        auto manager = static_cast<Manager*>(mgr);
        
        // Manual directory traversal compatible with all libkiwix versions
        for (const auto& entry : std::filesystem::recursive_directory_iterator(dir_path)) {
            if (entry.is_regular_file() && entry.path().extension() == ".zim") {
                manager->addBookFromPath(entry.path().string());
            }
        }
    } catch(...) { 
        // Catch filesystem errors (e.g., directory does not exist) or kiwix errors
    }
}

// --- Book ---

void kiwix_book_free(kiwix_book_t book) {
    delete static_cast<Book*>(book);
}

char* kiwix_book_get_id(kiwix_book_t book) {
    try { return copy_string(static_cast<Book*>(book)->getId()); } catch(...) { return nullptr; }
}

char* kiwix_book_get_title(kiwix_book_t book) {
    try { return copy_string(static_cast<Book*>(book)->getTitle()); } catch(...) { return nullptr; }
}

char* kiwix_book_get_description(kiwix_book_t book) {
    try { return copy_string(static_cast<Book*>(book)->getDescription()); } catch(...) { return nullptr; }
}

char* kiwix_book_get_path(kiwix_book_t book) {
    try { return copy_string(static_cast<Book*>(book)->getPath()); } catch(...) { return nullptr; }
}

uint64_t kiwix_book_get_article_count(kiwix_book_t book) {
    try { return static_cast<Book*>(book)->getArticleCount(); } catch(...) { return 0; }
}

uint64_t kiwix_book_get_size(kiwix_book_t book) {
    try { return static_cast<Book*>(book)->getSize(); } catch(...) { return 0; }
}

// --- Server ---

kiwix_server_t kiwix_server_new(kiwix_library_t lib) {
    try {
        auto s_lib = *static_cast<std::shared_ptr<Library>*>(lib);
        return new Server(s_lib);
    } catch(...) { return nullptr; }
}

void kiwix_server_free(kiwix_server_t server) {
    delete static_cast<Server*>(server);
}

void kiwix_server_set_port(kiwix_server_t server, int port) {
    try { static_cast<Server*>(server)->setPort(port); } catch(...) { }
}

void kiwix_server_set_block_external_links(kiwix_server_t server, bool block) {
    try { static_cast<Server*>(server)->setBlockExternalLinks(block); } catch(...) { }
}

bool kiwix_server_start(kiwix_server_t server) {
    try { return static_cast<Server*>(server)->start(); } catch(...) { return false; }
}

void kiwix_server_stop(kiwix_server_t server) {
    try { static_cast<Server*>(server)->stop(); } catch(...) { }
}

}
