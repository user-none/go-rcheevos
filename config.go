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

// GetTimeMillisecsFunc returns the current time in milliseconds.
// Used for accurate leaderboard timing.
type GetTimeMillisecsFunc func() uint64

var (
	timeMillisecsFuncs   = make(map[*C.rc_client_t]GetTimeMillisecsFunc)
	timeMillisecsFuncsMu sync.RWMutex
)

//export goGetTimeMillisecsCallback
func goGetTimeMillisecsCallback(client *C.rc_client_t) C.uint64_t {
	timeMillisecsFuncsMu.RLock()
	fn := timeMillisecsFuncs[client]
	timeMillisecsFuncsMu.RUnlock()

	if fn != nil {
		return C.uint64_t(fn())
	}
	return 0
}

// SetEncoreModeEnabled enables or disables encore mode.
// When enabled, achievements can be re-triggered even if already unlocked.
// Must be set before loading a game.
func (c *Client) SetEncoreModeEnabled(enabled bool) {
	val := C.int(0)
	if enabled {
		val = 1
	}
	C.rc_client_set_encore_mode_enabled(c.ptr, val)
}

// GetEncoreModeEnabled returns whether encore mode is enabled.
func (c *Client) GetEncoreModeEnabled() bool {
	return C.rc_client_get_encore_mode_enabled(c.ptr) != 0
}

// SetUnofficialEnabled enables or disables loading of unofficial achievements.
// Unofficial achievements are community-created and not part of the official set.
// Must be set before loading a game.
func (c *Client) SetUnofficialEnabled(enabled bool) {
	val := C.int(0)
	if enabled {
		val = 1
	}
	C.rc_client_set_unofficial_enabled(c.ptr, val)
}

// GetUnofficialEnabled returns whether unofficial achievements are enabled.
func (c *Client) GetUnofficialEnabled() bool {
	return C.rc_client_get_unofficial_enabled(c.ptr) != 0
}

// SetSpectatorModeEnabled enables or disables spectator mode.
// In spectator mode, events fire locally but nothing is submitted to the server.
// Can be toggled while a game is loaded.
func (c *Client) SetSpectatorModeEnabled(enabled bool) {
	val := C.int(0)
	if enabled {
		val = 1
	}
	C.rc_client_set_spectator_mode_enabled(c.ptr, val)
}

// GetSpectatorModeEnabled returns whether spectator mode is enabled.
func (c *Client) GetSpectatorModeEnabled() bool {
	return C.rc_client_get_spectator_mode_enabled(c.ptr) != 0
}

// SetHost sets the server hostname to use.
// Pass empty string to use the default RetroAchievements server.
func (c *Client) SetHost(hostname string) {
	if hostname == "" {
		C.rc_client_set_host(c.ptr, nil)
	} else {
		cHostname := C.CString(hostname)
		defer C.free(unsafe.Pointer(cHostname))
		C.rc_client_set_host(c.ptr, cHostname)
	}
}

// SetGetTimeMillisecsFunction sets a custom function to get the current time in milliseconds.
// This is used for accurate leaderboard timing.
// If not set, the default system time is used.
func (c *Client) SetGetTimeMillisecsFunction(fn GetTimeMillisecsFunc) {
	timeMillisecsFuncsMu.Lock()
	if fn != nil {
		timeMillisecsFuncs[c.ptr] = fn
	} else {
		delete(timeMillisecsFuncs, c.ptr)
	}
	timeMillisecsFuncsMu.Unlock()

	if fn != nil {
		cSetGetTimeMillisecsFunction(c.ptr)
	}
}

// GetUserAgentClause returns a string that can be added to the User-Agent header
// to identify the rcheevos version being used.
func (c *Client) GetUserAgentClause() string {
	buffer := make([]byte, 64)
	length := C.rc_client_get_user_agent_clause(c.ptr, (*C.char)(unsafe.Pointer(&buffer[0])), C.size_t(len(buffer)))
	if length == 0 {
		return ""
	}
	return string(buffer[:length])
}

// SetUserdata attaches client-specific data to the runtime.
// This can be retrieved later with GetUserdata.
func (c *Client) SetUserdata(userdata unsafe.Pointer) {
	C.rc_client_set_userdata(c.ptr, userdata)
}

// GetUserdata returns the client-specific data attached to the runtime.
func (c *Client) GetUserdata() unsafe.Pointer {
	return C.rc_client_get_userdata(c.ptr)
}

// SetReadMemoryFunction changes the memory read callback.
// This allows changing the memory source after the client is created.
func (c *Client) SetReadMemoryFunction(fn ReadMemoryFunc) {
	c.mu.Lock()
	c.readMemoryFunc = fn
	c.mu.Unlock()
}

// SetAllowBackgroundMemoryReads specifies whether rcheevos is allowed to read memory
// outside of DoFrame/Idle calls. This is used for async processing.
func (c *Client) SetAllowBackgroundMemoryReads(allowed bool) {
	val := C.int(0)
	if allowed {
		val = 1
	}
	C.rc_client_set_allow_background_memory_reads(c.ptr, val)
}
