package rcheevos

/*
#include <stdlib.h>
#include "rc_client.h"
*/
import "C"
import (
	"sync"
	"unsafe"
)

// ChangeMediaCallback is called when media change completes
type ChangeMediaCallback func(result int, errorMessage string)

var (
	changeMediaCallbacks   = make(map[*C.rc_client_t]ChangeMediaCallback)
	changeMediaCallbacksMu sync.RWMutex
)

//export goChangeMediaCallback
func goChangeMediaCallback(result C.int, errorMessage *C.char, client *C.rc_client_t, userdata unsafe.Pointer) {
	changeMediaCallbacksMu.RLock()
	callback := changeMediaCallbacks[client]
	changeMediaCallbacksMu.RUnlock()

	if callback != nil {
		var errMsg string
		if errorMessage != nil {
			errMsg = C.GoString(errorMessage)
		}
		callback(int(result), errMsg)
	}
}

// ChangeMedia changes the active disc in a multi-disc game using a known hash.
// This allows achievements to span multiple discs.
func (c *Client) ChangeMedia(hash string, callback ChangeMediaCallback) {
	changeMediaCallbacksMu.Lock()
	changeMediaCallbacks[c.ptr] = callback
	changeMediaCallbacksMu.Unlock()

	cHash := C.CString(hash)
	defer C.free(unsafe.Pointer(cHash))

	cBeginChangeMedia(c.ptr, cHash)
}

// IdentifyAndChangeMedia identifies and changes the active disc in a multi-disc game.
// Either filePath or data should be provided (not both).
func (c *Client) IdentifyAndChangeMedia(filePath string, data []byte, callback ChangeMediaCallback) {
	changeMediaCallbacksMu.Lock()
	changeMediaCallbacks[c.ptr] = callback
	changeMediaCallbacksMu.Unlock()

	var cFilePath *C.char
	if filePath != "" {
		cFilePath = C.CString(filePath)
		defer C.free(unsafe.Pointer(cFilePath))
	}

	var dataPtr *C.uint8_t
	if len(data) > 0 {
		dataPtr = (*C.uint8_t)(unsafe.Pointer(&data[0]))
	}

	cBeginIdentifyAndChangeMedia(c.ptr, cFilePath, dataPtr, C.size_t(len(data)))
}

// GetGameImageURL returns the URL for the game's badge/box art image.
func (c *Client) GetGameImageURL() string {
	game := C.rc_client_get_game_info(c.ptr)
	if game == nil {
		return ""
	}

	buffer := make([]byte, 256)
	result := C.rc_client_game_get_image_url(game, (*C.char)(unsafe.Pointer(&buffer[0])), C.size_t(len(buffer)))
	if result != C.RC_OK {
		return ""
	}

	for i, b := range buffer {
		if b == 0 {
			return string(buffer[:i])
		}
	}
	return string(buffer)
}

// GetUserImageURL returns the URL for the current user's avatar.
func (c *Client) GetUserImageURL() string {
	user := C.rc_client_get_user_info(c.ptr)
	if user == nil {
		return ""
	}

	buffer := make([]byte, 256)
	result := C.rc_client_user_get_image_url(user, (*C.char)(unsafe.Pointer(&buffer[0])), C.size_t(len(buffer)))
	if result != C.RC_OK {
		return ""
	}

	for i, b := range buffer {
		if b == 0 {
			return string(buffer[:i])
		}
	}
	return string(buffer)
}
