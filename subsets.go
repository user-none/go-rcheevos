package rcheevos

/*
#include "rc_client.h"
*/
import "C"
import "unsafe"

// Subset represents a game subset (e.g., "Base Game", "Bonus", "Challenge Mode")
type Subset struct {
	ID              uint32
	Title           string
	BadgeName       string
	NumAchievements uint32
	NumLeaderboards uint32
}

// SubsetList represents a list of subsets for the current game
type SubsetList struct {
	Subsets []*Subset
	ptr     *C.rc_client_subset_list_t
}

// GetSubsetInfo returns information about a specific subset.
func (c *Client) GetSubsetInfo(subsetID uint32) *Subset {
	subset := C.rc_client_get_subset_info(c.ptr, C.uint32_t(subsetID))
	if subset == nil {
		return nil
	}

	return &Subset{
		ID:              uint32(subset.id),
		Title:           C.GoString(subset.title),
		BadgeName:       C.GoString(&subset.badge_name[0]),
		NumAchievements: uint32(subset.num_achievements),
		NumLeaderboards: uint32(subset.num_leaderboards),
	}
}

// CreateSubsetList creates a list of all subsets for the currently loaded game.
// The returned list must be destroyed by calling Destroy() when no longer needed.
func (c *Client) CreateSubsetList() *SubsetList {
	ptr := C.rc_client_create_subset_list(c.ptr)
	if ptr == nil {
		return nil
	}

	list := &SubsetList{ptr: ptr}

	if ptr.num_subsets > 0 {
		list.Subsets = make([]*Subset, ptr.num_subsets)
		subsets := unsafe.Slice(ptr.subsets, ptr.num_subsets)

		for i := range subsets {
			s := subsets[i]
			list.Subsets[i] = &Subset{
				ID:              uint32(s.id),
				Title:           C.GoString(s.title),
				BadgeName:       C.GoString(&s.badge_name[0]),
				NumAchievements: uint32(s.num_achievements),
				NumLeaderboards: uint32(s.num_leaderboards),
			}
		}
	}

	return list
}

// Destroy releases resources associated with the subset list.
func (l *SubsetList) Destroy() {
	if l.ptr != nil {
		C.rc_client_destroy_subset_list(l.ptr)
		l.ptr = nil
	}
}

// GetUserSubsetSummary returns a summary of the user's progress in a specific subset.
func (c *Client) GetUserSubsetSummary(subsetID uint32) *UserGameSummary {
	var summary C.rc_client_user_game_summary_t
	C.rc_client_get_user_subset_summary(c.ptr, C.uint32_t(subsetID), &summary)

	return &UserGameSummary{
		NumCoreAchievements:        uint32(summary.num_core_achievements),
		NumUnofficialAchievements:  uint32(summary.num_unofficial_achievements),
		NumUnlockedAchievements:    uint32(summary.num_unlocked_achievements),
		NumUnsupportedAchievements: uint32(summary.num_unsupported_achievements),
		PointsCore:                 uint32(summary.points_core),
		PointsUnlocked:             uint32(summary.points_unlocked),
	}
}
