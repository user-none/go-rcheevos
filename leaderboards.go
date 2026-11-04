package rcheevos

/*
#include <stdlib.h>
#include "rc_client.h"
*/
import "C"
import (
	"time"
	"unsafe"
)

// Leaderboard states
const (
	LeaderboardStateInactive = C.RC_CLIENT_LEADERBOARD_STATE_INACTIVE
	LeaderboardStateActive   = C.RC_CLIENT_LEADERBOARD_STATE_ACTIVE
	LeaderboardStateTracking = C.RC_CLIENT_LEADERBOARD_STATE_TRACKING
	LeaderboardStateDisabled = C.RC_CLIENT_LEADERBOARD_STATE_DISABLED
)

// Leaderboard formats
const (
	LeaderboardFormatTime  = C.RC_CLIENT_LEADERBOARD_FORMAT_TIME
	LeaderboardFormatScore = C.RC_CLIENT_LEADERBOARD_FORMAT_SCORE
	LeaderboardFormatValue = C.RC_CLIENT_LEADERBOARD_FORMAT_VALUE
)

// Leaderboard list grouping options
const (
	LeaderboardListGroupingNone     = C.RC_CLIENT_LEADERBOARD_LIST_GROUPING_NONE
	LeaderboardListGroupingTracking = C.RC_CLIENT_LEADERBOARD_LIST_GROUPING_TRACKING
)

// Leaderboard bucket types
const (
	LeaderboardBucketUnknown     = C.RC_CLIENT_LEADERBOARD_BUCKET_UNKNOWN
	LeaderboardBucketInactive    = C.RC_CLIENT_LEADERBOARD_BUCKET_INACTIVE
	LeaderboardBucketActive      = C.RC_CLIENT_LEADERBOARD_BUCKET_ACTIVE
	LeaderboardBucketUnsupported = C.RC_CLIENT_LEADERBOARD_BUCKET_UNSUPPORTED
	LeaderboardBucketAll         = C.RC_CLIENT_LEADERBOARD_BUCKET_ALL
)

// LeaderboardBucket represents a group of leaderboards
type LeaderboardBucket struct {
	Leaderboards []*Leaderboard
	Label        string
	SubsetID     uint32
	BucketType   uint8
}

// LeaderboardList represents a list of leaderboard buckets
type LeaderboardList struct {
	Buckets []*LeaderboardBucket
	ptr     *C.rc_client_leaderboard_list_t
}

// CreateLeaderboardList creates a list of leaderboards matching the specified grouping.
// Grouping should be LeaderboardListGroupingNone or LeaderboardListGroupingTracking.
// The returned list must be destroyed by calling Destroy() when no longer needed.
func (c *Client) CreateLeaderboardList(grouping int) *LeaderboardList {
	ptr := C.rc_client_create_leaderboard_list(c.ptr, C.int(grouping))
	if ptr == nil {
		return nil
	}

	list := &LeaderboardList{ptr: ptr}

	// Convert C buckets to Go
	if ptr.num_buckets > 0 {
		list.Buckets = make([]*LeaderboardBucket, ptr.num_buckets)
		buckets := unsafe.Slice(ptr.buckets, ptr.num_buckets)

		for i := range buckets {
			bucket := &LeaderboardBucket{
				Label:      C.GoString(buckets[i].label),
				SubsetID:   uint32(buckets[i].subset_id),
				BucketType: uint8(buckets[i].bucket_type),
			}

			if buckets[i].num_leaderboards > 0 {
				bucket.Leaderboards = make([]*Leaderboard, buckets[i].num_leaderboards)
				leaderboards := unsafe.Slice(buckets[i].leaderboards, buckets[i].num_leaderboards)

				for j := range leaderboards {
					bucket.Leaderboards[j] = leaderboardFromC(leaderboards[j])
				}
			}

			list.Buckets[i] = bucket
		}
	}

	return list
}

// Destroy releases resources associated with the leaderboard list.
func (l *LeaderboardList) Destroy() {
	if l.ptr != nil {
		C.rc_client_destroy_leaderboard_list(l.ptr)
		l.ptr = nil
	}
}

// GetAllLeaderboards returns a flat slice of all leaderboards across all buckets.
func (l *LeaderboardList) GetAllLeaderboards() []*Leaderboard {
	var all []*Leaderboard
	for _, bucket := range l.Buckets {
		all = append(all, bucket.Leaderboards...)
	}
	return all
}

// LeaderboardEntry represents a single entry in a leaderboard
type LeaderboardEntry struct {
	User      string
	Display   string
	Submitted time.Time
	Rank      uint32
	Index     uint32
}

// LeaderboardEntryList represents a list of leaderboard entries
type LeaderboardEntryList struct {
	Entries      []*LeaderboardEntry
	TotalEntries uint32
	UserIndex    int32 // -1 if user not in list
	ptr          *C.rc_client_leaderboard_entry_list_t
}

// Destroy releases resources associated with the leaderboard entry list.
func (l *LeaderboardEntryList) Destroy() {
	if l.ptr != nil {
		C.rc_client_destroy_leaderboard_entry_list(l.ptr)
		l.ptr = nil
	}
}

// GetUserImageURL returns the URL for a leaderboard entry user's avatar.
func (e *LeaderboardEntry) GetUserImageURL() string {
	// This requires the original C entry pointer which we don't keep.
	// For simplicity, we'll construct the URL directly.
	// RetroAchievements avatar URL format: https://media.retroachievements.org/UserPic/{user}.png
	return "https://media.retroachievements.org/UserPic/" + e.User + ".png"
}

// LeaderboardEntriesCallback is called when leaderboard entries are fetched
type LeaderboardEntriesCallback func(result int, errorMessage string, entries *LeaderboardEntryList)

//export goLeaderboardEntriesCallback
func goLeaderboardEntriesCallback(result C.int, errorMessage *C.char, list *C.rc_client_leaderboard_entry_list_t, client *C.rc_client_t, userdata unsafe.Pointer) {
	var callback LeaderboardEntriesCallback
	if cb, ok := finishRequest(client, userdata); ok {
		callback = cb.(LeaderboardEntriesCallback)
	}
	if callback == nil {
		if list != nil {
			C.rc_client_destroy_leaderboard_entry_list(list)
		}
		return
	}

	errMsg := goErrorMessage(errorMessage)

	var entryList *LeaderboardEntryList
	if list != nil {
		entryList = &LeaderboardEntryList{
			TotalEntries: uint32(list.total_entries),
			UserIndex:    int32(list.user_index),
			ptr:          list,
		}

		if list.num_entries > 0 {
			entryList.Entries = make([]*LeaderboardEntry, list.num_entries)
			entries := unsafe.Slice(list.entries, list.num_entries)

			for i := range entries {
				entryList.Entries[i] = &LeaderboardEntry{
					User:      C.GoString(entries[i].user),
					Display:   C.GoString(&entries[i].display[0]),
					Submitted: time.Unix(int64(entries[i].submitted), 0),
					Rank:      uint32(entries[i].rank),
					Index:     uint32(entries[i].index),
				}
			}
		}
	}

	callback(int(result), errMsg, entryList)
}

// FetchLeaderboardEntries fetches leaderboard entries from the server.
// firstEntry is 1-based (1 = first entry).
// count is the number of entries to fetch.
// The callback receives an allocated list that must be destroyed by calling Destroy().
// It receives the result of this fetch only, and may be called before
// FetchLeaderboardEntries returns.
func (c *Client) FetchLeaderboardEntries(leaderboardID uint32, firstEntry, count uint32, callback LeaderboardEntriesCallback) {
	h, ok := c.beginRequest(requestOther, callback)
	if !ok {
		if callback != nil {
			callback(InvalidState, msgClientRequired, nil)
		}
		return
	}

	cBeginFetchLeaderboardEntries(c.ptr, leaderboardID, firstEntry, count, h)
}

// FetchLeaderboardEntriesAroundUser fetches leaderboard entries surrounding the current user.
// count is the total number of entries to fetch (user will be approximately in the middle).
// The callback receives an allocated list that must be destroyed by calling Destroy().
// It receives the result of this fetch only, and may be called before
// FetchLeaderboardEntriesAroundUser returns.
func (c *Client) FetchLeaderboardEntriesAroundUser(leaderboardID uint32, count uint32, callback LeaderboardEntriesCallback) {
	h, ok := c.beginRequest(requestOther, callback)
	if !ok {
		if callback != nil {
			callback(InvalidState, msgClientRequired, nil)
		}
		return
	}

	cBeginFetchLeaderboardEntriesAroundUser(c.ptr, leaderboardID, count, h)
}
