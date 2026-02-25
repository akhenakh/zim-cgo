package kiwix

/*
// Libkiwix uses modern C++ features, typically requiring C++14 or C++17
#cgo CXXFLAGS: -std=c++17
#cgo LDFLAGS: -lkiwix -lzim -lmicrohttpd
#include <stdlib.h>
#include "kiwix_wrapper.h"
*/
import "C"
import (
	"errors"
	"runtime"
	"unsafe"
)

// Library represents a collection of ZIM files (books)
type Library struct {
	ptr C.kiwix_library_t
}

func NewLibrary() *Library {
	ptr := C.kiwix_library_new()
	if ptr == nil {
		return nil
	}
	lib := &Library{ptr: ptr}
	runtime.SetFinalizer(lib, (*Library).Close)
	return lib
}

func (l *Library) Close() {
	if l.ptr != nil {
		C.kiwix_library_free(l.ptr)
		l.ptr = nil
	}
}

func (l *Library) GetBookCount(local, remote bool) uint {
	return uint(C.kiwix_library_get_book_count(l.ptr, C.bool(local), C.bool(remote)))
}

// Manager allows editing and updating the Library
type Manager struct {
	ptr C.kiwix_manager_t
}

func NewManager(lib *Library) *Manager {
	ptr := C.kiwix_manager_new(lib.ptr)
	if ptr == nil {
		return nil
	}
	mgr := &Manager{ptr: ptr}
	runtime.SetFinalizer(mgr, (*Manager).Close)
	return mgr
}

func (m *Manager) Close() {
	if m.ptr != nil {
		C.kiwix_manager_free(m.ptr)
		m.ptr = nil
	}
}

func (m *Manager) AddBookFromPath(path string) bool {
	cPath := C.CString(path)
	defer C.free(unsafe.Pointer(cPath))
	return bool(C.kiwix_manager_add_book_from_path(m.ptr, cPath))
}

func (m *Manager) AddBooksFromDirectory(dirPath string) {
	cPath := C.CString(dirPath)
	defer C.free(unsafe.Pointer(cPath))
	C.kiwix_manager_add_books_from_directory(m.ptr, cPath)
}

// CPPServer serves the content of the library via the native C++ libmicrohttpd server
type CPPServer struct {
	ptr C.kiwix_server_t
}

func NewCPPServer(lib *Library) *CPPServer {
	ptr := C.kiwix_server_new(lib.ptr)
	if ptr == nil {
		return nil
	}
	srv := &CPPServer{ptr: ptr}
	runtime.SetFinalizer(srv, (*CPPServer).Close)
	return srv
}

func (s *CPPServer) Close() {
	if s.ptr != nil {
		C.kiwix_server_free(s.ptr)
		s.ptr = nil
	}
}

func (s *CPPServer) SetPort(port int) {
	C.kiwix_server_set_port(s.ptr, C.int(port))
}

func (s *CPPServer) SetBlockExternalLinks(block bool) {
	C.kiwix_server_set_block_external_links(s.ptr, C.bool(block))
}

func (s *CPPServer) Start() error {
	success := C.kiwix_server_start(s.ptr)
	if !success {
		return errors.New("failed to start kiwix C++ server")
	}
	return nil
}

func (s *CPPServer) Stop() {
	C.kiwix_server_stop(s.ptr)
}
