package rcheevos

/*
#include "rc_client.h"
*/
import "C"

// IsProcessingRequired returns whether there are any active achievements, leaderboards,
// or rich presence that need processing. Can be used to skip DoFrame calls when nothing
// needs to be evaluated.
func (c *Client) IsProcessingRequired() bool {
	return C.rc_client_is_processing_required(c.ptr) != 0
}

// CanPause determines if a sufficient number of frames have been processed since the
// last call to CanPause. This is used for hardcore anti-cheat to prevent pausing
// mid-leaderboard attempt.
//
// Returns true if pausing is allowed, false if not.
// If false, framesRemaining indicates how many more frames must be processed before
// pausing is allowed (can be converted to seconds for display).
//
// Only call this when the user is actually trying to pause.
func (c *Client) CanPause() (allowed bool, framesRemaining uint32) {
	var remaining C.uint32_t
	result := C.rc_client_can_pause(c.ptr, &remaining)
	return result != 0, uint32(remaining)
}

// AbortAsync marks an async process as aborted. The associated callback will not be called.
// Pass the async handle returned by login/load game functions.
func (c *Client) AbortAsync(handle *AsyncHandle) {
	if handle != nil && handle.ptr != nil {
		C.rc_client_abort_async(c.ptr, handle.ptr)
	}
}

// AsyncHandle represents an async operation that can be aborted
type AsyncHandle struct {
	ptr *C.rc_client_async_handle_t
}

// Load game state constants
const (
	LoadGameStateNone           = C.RC_CLIENT_LOAD_GAME_STATE_NONE
	LoadGameStateAwaitLogin     = C.RC_CLIENT_LOAD_GAME_STATE_AWAIT_LOGIN
	LoadGameStateIdentifying    = C.RC_CLIENT_LOAD_GAME_STATE_IDENTIFYING_GAME
	LoadGameStateFetchingData   = C.RC_CLIENT_LOAD_GAME_STATE_FETCHING_GAME_DATA
	LoadGameStateStartingSession = C.RC_CLIENT_LOAD_GAME_STATE_STARTING_SESSION
	LoadGameStateDone           = C.RC_CLIENT_LOAD_GAME_STATE_DONE
	LoadGameStateAborted        = C.RC_CLIENT_LOAD_GAME_STATE_ABORTED
)

// GetLoadGameState returns the current progress of the asynchronous load game process.
func (c *Client) GetLoadGameState() int {
	return int(C.rc_client_get_load_game_state(c.ptr))
}

// LoadGameStateName returns a human-readable name for a load game state
func LoadGameStateName(state int) string {
	switch state {
	case LoadGameStateNone:
		return "None"
	case LoadGameStateAwaitLogin:
		return "Awaiting Login"
	case LoadGameStateIdentifying:
		return "Identifying Game"
	case LoadGameStateFetchingData:
		return "Fetching Game Data"
	case LoadGameStateStartingSession:
		return "Starting Session"
	case LoadGameStateDone:
		return "Done"
	case LoadGameStateAborted:
		return "Aborted"
	default:
		return "Unknown"
	}
}
