package zim

/*
#cgo CXXFLAGS: -std=c++11
#cgo LDFLAGS: -lzim
#include <stdlib.h>
#include "zim_wrapper.h"
*/
import "C"
import (
	"errors"
	"runtime"
	"unsafe"
)

// Archive represents a readable ZIM archive
type Archive struct {
	ptr C.zim_archive_t
}

// NewArchive opens a ZIM archive from the given file path
func NewArchive(path string) (*Archive, error) {
	cPath := C.CString(path)
	defer C.free(unsafe.Pointer(cPath))

	ptr := C.zim_archive_new(cPath)
	if ptr == nil {
		return nil, errors.New("failed to open archive or invalid format")
	}

	arch := &Archive{ptr: ptr}
	runtime.SetFinalizer(arch, (*Archive).Close)
	return arch, nil
}

// Close frees the underlying C++ archive resources
func (a *Archive) Close() {
	if a.ptr != nil {
		C.zim_archive_free(a.ptr)
		a.ptr = nil
	}
}

// GetEntryCount returns the number of user entries
func (a *Archive) GetEntryCount() uint64 {
	return uint64(C.zim_archive_get_entry_count(a.ptr))
}

type Entry struct {
	ptr C.zim_entry_t
}

// GetEntryByPath retrieves an entry by its URL path inside the archive
func (a *Archive) GetEntryByPath(path string) (*Entry, error) {
	cPath := C.CString(path)
	defer C.free(unsafe.Pointer(cPath))

	ptr := C.zim_archive_get_entry_by_path(a.ptr, cPath)
	if ptr == nil {
		return nil, errors.New("entry not found")
	}

	entry := &Entry{ptr: ptr}
	runtime.SetFinalizer(entry, (*Entry).Close)
	return entry, nil
}

// GetMainEntry retrieves the default index/home page of the ZIM archive
func (a *Archive) GetMainEntry() (*Entry, error) {
	ptr := C.zim_archive_get_main_entry(a.ptr)
	if ptr == nil {
		return nil, errors.New("main entry not found or archive has no main entry")
	}

	entry := &Entry{ptr: ptr}
	runtime.SetFinalizer(entry, (*Entry).Close)
	return entry, nil
}

// GetEntryByIndex retrieves an entry by its numerical index (sorted by path)
func (a *Archive) GetEntryByIndex(idx uint32) (*Entry, error) {
	ptr := C.zim_archive_get_entry_by_index(a.ptr, C.uint32_t(idx))
	if ptr == nil {
		return nil, errors.New("entry not found at the specified index")
	}

	entry := &Entry{ptr: ptr}
	runtime.SetFinalizer(entry, (*Entry).Close)
	return entry, nil
}

func (e *Entry) Close() {
	if e.ptr != nil {
		C.zim_entry_free(e.ptr)
		e.ptr = nil
	}
}

type Item struct {
	ptr C.zim_item_t
}

// GetItem resolves the entry payload. If 'follow' is true, it automatically resolves redirects.
func (e *Entry) GetItem(follow bool) (*Item, error) {
	ptr := C.zim_entry_get_item(e.ptr, C.bool(follow))
	if ptr == nil {
		return nil, errors.New("failed to retrieve item")
	}

	item := &Item{ptr: ptr}
	runtime.SetFinalizer(item, (*Item).Close)
	return item, nil
}

func (i *Item) Close() {
	if i.ptr != nil {
		C.zim_item_free(i.ptr)
		i.ptr = nil
	}
}

func (i *Item) GetPath() string {
	cStr := C.zim_item_get_path(i.ptr)
	if cStr == nil {
		return ""
	}
	defer C.free(unsafe.Pointer(cStr))
	return C.GoString(cStr)
}

func (i *Item) GetTitle() string {
	cStr := C.zim_item_get_title(i.ptr)
	if cStr == nil {
		return ""
	}
	defer C.free(unsafe.Pointer(cStr))
	return C.GoString(cStr)
}

func (i *Item) GetMimetype() string {
	cStr := C.zim_item_get_mimetype(i.ptr)
	if cStr == nil {
		return ""
	}
	defer C.free(unsafe.Pointer(cStr))
	return C.GoString(cStr)
}

func (i *Item) GetSize() uint64 {
	return uint64(C.zim_item_get_size(i.ptr))
}

// GetData reads the blob buffer for this item and returns it as a byte slice.
func (i *Item) GetData() []byte {
	var size C.uint64_t
	cData := C.zim_item_get_data(i.ptr, &size)
	if cData == nil {
		return nil
	}
	defer C.free(unsafe.Pointer(cData))

	// Convert C allocation directly into Go slice
	return C.GoBytes(unsafe.Pointer(cData), C.int(size))
}
