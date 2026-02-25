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

type Compression int

const (
	CompressionNone Compression = 1
	CompressionZstd Compression = 5
)

// Creator represents the engine that builds a new ZIM archive
type Creator struct {
	ptr C.zim_creator_t
}

func NewCreator() (*Creator, error) {
	ptr := C.zim_creator_new()
	if ptr == nil {
		return nil, errors.New("failed to initialize ZIM creator")
	}

	c := &Creator{ptr: ptr}
	runtime.SetFinalizer(c, (*Creator).Close)
	return c, nil
}

func (c *Creator) Close() {
	if c.ptr != nil {
		C.zim_creator_free(c.ptr)
		c.ptr = nil
	}
}

func (c *Creator) ConfigVerbose(verbose bool) {
	C.zim_creator_config_verbose(c.ptr, C.bool(verbose))
}

func (c *Creator) ConfigCompression(comp Compression) {
	C.zim_creator_config_compression(c.ptr, C.int(comp))
}

func (c *Creator) StartZimCreation(filepath string) error {
	cPath := C.CString(filepath)
	defer C.free(unsafe.Pointer(cPath))

	if !bool(C.zim_creator_start_zim_creation(c.ptr, cPath)) {
		return errors.New("failed to start ZIM creation (is path writable?)")
	}
	return nil
}

func (c *Creator) SetMainPath(mainPath string) error {
	cPath := C.CString(mainPath)
	defer C.free(unsafe.Pointer(cPath))

	if !bool(C.zim_creator_set_main_path(c.ptr, cPath)) {
		return errors.New("failed to set main path")
	}
	return nil
}

func (c *Creator) AddItem(item *WriterItem) error {
	if !bool(C.zim_creator_add_item(c.ptr, item.ptr)) {
		return errors.New("failed to add item to archive (duplicate path?)")
	}
	return nil
}

func (c *Creator) AddMetadata(name, content string) error {
	cName := C.CString(name)
	cContent := C.CString(content)
	defer C.free(unsafe.Pointer(cName))
	defer C.free(unsafe.Pointer(cContent))

	if !bool(C.zim_creator_add_metadata(c.ptr, cName, cContent)) {
		return errors.New("failed to add metadata")
	}
	return nil
}

func (c *Creator) AddIllustration(size uint, content []byte) error {
	if len(content) == 0 {
		return errors.New("illustration content cannot be empty")
	}

	cContent := (*C.char)(unsafe.Pointer(&content[0]))
	if !bool(C.zim_creator_add_illustration(c.ptr, C.uint(size), cContent, C.uint64_t(len(content)))) {
		return errors.New("failed to add illustration")
	}
	return nil
}

func (c *Creator) FinishZimCreation() error {
	if !bool(C.zim_creator_finish_zim_creation(c.ptr)) {
		return errors.New("failed to finalize and finish ZIM creation")
	}
	return nil
}

// WriterItem represents an entry pending insertion into a ZIM archive
type WriterItem struct {
	ptr C.zim_writer_item_t
}

func NewStringItem(path, mimetype, title string, content []byte, isFrontArticle bool) (*WriterItem, error) {
	cPath := C.CString(path)
	cMime := C.CString(mimetype)
	cTitle := C.CString(title)
	defer C.free(unsafe.Pointer(cPath))
	defer C.free(unsafe.Pointer(cMime))
	defer C.free(unsafe.Pointer(cTitle))

	var cContent *C.char
	if len(content) > 0 {
		cContent = (*C.char)(unsafe.Pointer(&content[0]))
	}

	ptr := C.zim_writer_string_item_new(cPath, cMime, cTitle, cContent, C.uint64_t(len(content)), C.bool(isFrontArticle))
	if ptr == nil {
		return nil, errors.New("failed to create string item")
	}

	item := &WriterItem{ptr: ptr}
	runtime.SetFinalizer(item, (*WriterItem).Close)
	return item, nil
}

func NewFileItem(path, mimetype, title, filepath string, isFrontArticle bool) (*WriterItem, error) {
	cPath := C.CString(path)
	cMime := C.CString(mimetype)
	cTitle := C.CString(title)
	cFilepath := C.CString(filepath)
	defer C.free(unsafe.Pointer(cPath))
	defer C.free(unsafe.Pointer(cMime))
	defer C.free(unsafe.Pointer(cTitle))
	defer C.free(unsafe.Pointer(cFilepath))

	ptr := C.zim_writer_file_item_new(cPath, cMime, cTitle, cFilepath, C.bool(isFrontArticle))
	if ptr == nil {
		return nil, errors.New("failed to create file item")
	}

	item := &WriterItem{ptr: ptr}
	runtime.SetFinalizer(item, (*WriterItem).Close)
	return item, nil
}

func (i *WriterItem) Close() {
	if i.ptr != nil {
		C.zim_writer_item_free(i.ptr)
		i.ptr = nil
	}
}
