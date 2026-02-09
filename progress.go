package rcheevos

/*
#include "rc_client.h"
*/
import "C"
import "unsafe"

// ProgressSize returns the number of bytes needed to serialize the runtime state.
// Use this to allocate a buffer for SerializeProgress.
func (c *Client) ProgressSize() int {
	return int(C.rc_client_progress_size(c.ptr))
}

// SerializeProgress serializes the runtime state into a buffer.
// This should be called when creating a savestate to preserve achievement progress.
// Returns the serialized data, or nil and an error code on failure.
func (c *Client) SerializeProgress() ([]byte, int) {
	size := c.ProgressSize()
	if size == 0 {
		return nil, OK
	}

	buffer := make([]byte, size)
	result := C.rc_client_serialize_progress_sized(c.ptr, (*C.uint8_t)(unsafe.Pointer(&buffer[0])), C.size_t(size))
	if result != C.RC_OK {
		return nil, int(result)
	}

	return buffer, OK
}

// DeserializeProgress restores the runtime state from a buffer.
// This should be called when loading a savestate to restore achievement progress.
// Returns OK on success, or an error code on failure.
func (c *Client) DeserializeProgress(data []byte) int {
	if len(data) == 0 {
		return OK
	}

	result := C.rc_client_deserialize_progress_sized(c.ptr, (*C.uint8_t)(unsafe.Pointer(&data[0])), C.size_t(len(data)))
	return int(result)
}
