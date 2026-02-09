package rcheevos

/*
#include <stdlib.h>
#include <string.h>
#include "rc_hash.h"
*/
import "C"
import "unsafe"

// HashIterator is used to generate hashes for ROM identification
type HashIterator struct {
	iter   C.rc_hash_iterator_t
	path   *C.char
	buffer unsafe.Pointer // C-allocated buffer copy
}

// NewHashIterator creates a new hash iterator for a file path
func NewHashIterator(path string) *HashIterator {
	h := &HashIterator{}
	h.path = C.CString(path)
	C.rc_hash_initialize_iterator(&h.iter, h.path, nil, 0)
	return h
}

// NewHashIteratorFromBuffer creates a new hash iterator from memory
func NewHashIteratorFromBuffer(path string, data []byte) *HashIterator {
	h := &HashIterator{}
	h.path = C.CString(path)

	var dataPtr *C.uint8_t
	if len(data) > 0 {
		// Copy data to C memory since the iterator stores the pointer
		h.buffer = C.malloc(C.size_t(len(data)))
		C.memcpy(h.buffer, unsafe.Pointer(&data[0]), C.size_t(len(data)))
		dataPtr = (*C.uint8_t)(h.buffer)
	}

	C.rc_hash_initialize_iterator(&h.iter, h.path, dataPtr, C.size_t(len(data)))
	return h
}

// Destroy releases resources associated with the iterator
func (h *HashIterator) Destroy() {
	C.rc_hash_destroy_iterator(&h.iter)
	if h.path != nil {
		C.free(unsafe.Pointer(h.path))
		h.path = nil
	}
	if h.buffer != nil {
		C.free(h.buffer)
		h.buffer = nil
	}
}

// Next generates the next hash. Returns empty string if no more hashes.
func (h *HashIterator) Next() string {
	var hash [33]C.char
	if C.rc_hash_iterate(&hash[0], &h.iter) == 0 {
		return ""
	}
	return C.GoString(&hash[0])
}

// GenerateHash generates a hash for the specified console from the iterator
func (h *HashIterator) GenerateHash(consoleID uint32) string {
	var hash [33]C.char
	if C.rc_hash_generate(&hash[0], C.uint32_t(consoleID), &h.iter) == 0 {
		return ""
	}
	return C.GoString(&hash[0])
}

// HashFromBuffer generates a hash from ROM data for a specific console
func HashFromBuffer(consoleID uint32, data []byte) string {
	if len(data) == 0 {
		return ""
	}

	var hash [33]C.char
	dataPtr := (*C.uint8_t)(unsafe.Pointer(&data[0]))

	if C.rc_hash_generate_from_buffer(&hash[0], C.uint32_t(consoleID), dataPtr, C.size_t(len(data))) == 0 {
		return ""
	}
	return C.GoString(&hash[0])
}

// HashFromFile generates a hash from a ROM file for a specific console
func HashFromFile(consoleID uint32, path string) string {
	cPath := C.CString(path)
	defer C.free(unsafe.Pointer(cPath))

	var hash [33]C.char
	if C.rc_hash_generate_from_file(&hash[0], C.uint32_t(consoleID), cPath) == 0 {
		return ""
	}
	return C.GoString(&hash[0])
}
