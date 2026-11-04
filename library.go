package rcheevos

/*
#include <stdlib.h>
#include "rc_client.h"
*/
import "C"
import "unsafe"

// HashLibraryEntry represents a mapping from hash to game ID
type HashLibraryEntry struct {
	Hash   string
	GameID uint32
}

// HashLibrary represents a collection of hash-to-game mappings for a console
type HashLibrary struct {
	Entries []*HashLibraryEntry
	ptr     *C.rc_client_hash_library_t
}

// Destroy releases resources associated with the hash library.
func (l *HashLibrary) Destroy() {
	if l.ptr != nil {
		C.rc_client_destroy_hash_library(l.ptr)
		l.ptr = nil
	}
}

// FetchHashLibraryCallback is called when the hash library fetch completes
type FetchHashLibraryCallback func(result int, errorMessage string, library *HashLibrary)

//export goFetchHashLibraryCallback
func goFetchHashLibraryCallback(result C.int, errorMessage *C.char, list *C.rc_client_hash_library_t, client *C.rc_client_t, userdata unsafe.Pointer) {
	var callback FetchHashLibraryCallback
	if cb, ok := finishRequest(client, userdata); ok {
		callback = cb.(FetchHashLibraryCallback)
	}
	if callback == nil {
		if list != nil {
			C.rc_client_destroy_hash_library(list)
		}
		return
	}

	errMsg := goErrorMessage(errorMessage)

	var library *HashLibrary
	if list != nil {
		library = &HashLibrary{ptr: list}

		if list.num_entries > 0 {
			library.Entries = make([]*HashLibraryEntry, list.num_entries)
			entries := unsafe.Slice(list.entries, list.num_entries)

			for i := range entries {
				library.Entries[i] = &HashLibraryEntry{
					Hash:   C.GoString(&entries[i].hash[0]),
					GameID: uint32(entries[i].game_id),
				}
			}
		}
	}

	callback(int(result), errMsg, library)
}

// FetchHashLibrary fetches all known hashes for a console.
// This can be used to build a game database or identify ROMs offline.
// The callback receives an allocated library that must be destroyed by calling Destroy().
// It receives the result of this fetch only, and may be called before
// FetchHashLibrary returns.
func (c *Client) FetchHashLibrary(consoleID uint32, callback FetchHashLibraryCallback) {
	h, ok := c.beginRequest(requestOther, callback)
	if !ok {
		if callback != nil {
			callback(InvalidState, msgClientRequired, nil)
		}
		return
	}

	cBeginFetchHashLibrary(c.ptr, consoleID, h)
}

// GameTitleEntry represents a game's title and badge information
type GameTitleEntry struct {
	GameID    uint32
	Title     string
	BadgeName string
	BadgeURL  string
}

// GameTitleList represents a collection of game titles
type GameTitleList struct {
	Entries []*GameTitleEntry
	ptr     *C.rc_client_game_title_list_t
}

// Destroy releases resources associated with the game title list.
func (l *GameTitleList) Destroy() {
	if l.ptr != nil {
		C.rc_client_destroy_game_title_list(l.ptr)
		l.ptr = nil
	}
}

// FetchGameTitlesCallback is called when the game titles fetch completes
type FetchGameTitlesCallback func(result int, errorMessage string, titles *GameTitleList)

//export goFetchGameTitlesCallback
func goFetchGameTitlesCallback(result C.int, errorMessage *C.char, list *C.rc_client_game_title_list_t, client *C.rc_client_t, userdata unsafe.Pointer) {
	var callback FetchGameTitlesCallback
	if cb, ok := finishRequest(client, userdata); ok {
		callback = cb.(FetchGameTitlesCallback)
	}
	if callback == nil {
		if list != nil {
			C.rc_client_destroy_game_title_list(list)
		}
		return
	}

	errMsg := goErrorMessage(errorMessage)

	var titles *GameTitleList
	if list != nil {
		titles = &GameTitleList{ptr: list}

		if list.num_entries > 0 {
			titles.Entries = make([]*GameTitleEntry, list.num_entries)
			entries := unsafe.Slice(list.entries, list.num_entries)

			for i := range entries {
				var badgeURL string
				if entries[i].badge_url != nil {
					badgeURL = C.GoString(entries[i].badge_url)
				}
				titles.Entries[i] = &GameTitleEntry{
					GameID:    uint32(entries[i].game_id),
					Title:     C.GoString(entries[i].title),
					BadgeName: C.GoString(&entries[i].badge_name[0]),
					BadgeURL:  badgeURL,
				}
			}
		}
	}

	callback(int(result), errMsg, titles)
}

// FetchGameTitles fetches titles and badges for the specified game IDs.
// This is useful for displaying a library with game names/icons.
// The callback receives an allocated list that must be destroyed by calling Destroy().
// It receives the result of this fetch only, and may be called before
// FetchGameTitles returns.
func (c *Client) FetchGameTitles(gameIDs []uint32, callback FetchGameTitlesCallback) {
	if len(gameIDs) == 0 {
		if callback != nil {
			callback(OK, "", nil)
		}
		return
	}

	h, ok := c.beginRequest(requestOther, callback)
	if !ok {
		if callback != nil {
			callback(InvalidState, msgClientRequired, nil)
		}
		return
	}

	cBeginFetchGameTitles(c.ptr, (*C.uint32_t)(unsafe.Pointer(&gameIDs[0])), uint32(len(gameIDs)), h)
}

// UserProgressEntry represents a user's progress in a single game
type UserProgressEntry struct {
	GameID                       uint32
	NumAchievements              uint32
	NumUnlockedAchievements      uint32
	NumUnlockedAchievementsHardcore uint32
}

// AllUserProgress represents a user's progress across all games for a console
type AllUserProgress struct {
	Entries []*UserProgressEntry
	ptr     *C.rc_client_all_user_progress_t
}

// Destroy releases resources associated with the user progress list.
func (l *AllUserProgress) Destroy() {
	if l.ptr != nil {
		C.rc_client_destroy_all_user_progress(l.ptr)
		l.ptr = nil
	}
}

// FetchAllUserProgressCallback is called when the user progress fetch completes
type FetchAllUserProgressCallback func(result int, errorMessage string, progress *AllUserProgress)

//export goFetchAllUserProgressCallback
func goFetchAllUserProgressCallback(result C.int, errorMessage *C.char, list *C.rc_client_all_user_progress_t, client *C.rc_client_t, userdata unsafe.Pointer) {
	var callback FetchAllUserProgressCallback
	if cb, ok := finishRequest(client, userdata); ok {
		callback = cb.(FetchAllUserProgressCallback)
	}
	if callback == nil {
		if list != nil {
			C.rc_client_destroy_all_user_progress(list)
		}
		return
	}

	errMsg := goErrorMessage(errorMessage)

	var progress *AllUserProgress
	if list != nil {
		progress = &AllUserProgress{ptr: list}

		if list.num_entries > 0 {
			progress.Entries = make([]*UserProgressEntry, list.num_entries)
			entries := unsafe.Slice(list.entries, list.num_entries)

			for i := range entries {
				progress.Entries[i] = &UserProgressEntry{
					GameID:                          uint32(entries[i].game_id),
					NumAchievements:                 uint32(entries[i].num_achievements),
					NumUnlockedAchievements:         uint32(entries[i].num_unlocked_achievements),
					NumUnlockedAchievementsHardcore: uint32(entries[i].num_unlocked_achievements_hardcore),
				}
			}
		}
	}

	callback(int(result), errMsg, progress)
}

// FetchAllUserProgress fetches the user's progress for all games on a console.
// This is useful for displaying a game browser with completion percentages.
// The callback receives an allocated list that must be destroyed by calling Destroy().
// It receives the result of this fetch only, and may be called before
// FetchAllUserProgress returns.
func (c *Client) FetchAllUserProgress(consoleID uint32, callback FetchAllUserProgressCallback) {
	h, ok := c.beginRequest(requestOther, callback)
	if !ok {
		if callback != nil {
			callback(InvalidState, msgClientRequired, nil)
		}
		return
	}

	cBeginFetchAllUserProgress(c.ptr, consoleID, h)
}
