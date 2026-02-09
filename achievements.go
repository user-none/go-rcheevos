package rcheevos

/*
#include "rc_client.h"
*/
import "C"
import "unsafe"

// Achievement list grouping options
const (
	AchievementListGroupingLockState = C.RC_CLIENT_ACHIEVEMENT_LIST_GROUPING_LOCK_STATE
	AchievementListGroupingProgress  = C.RC_CLIENT_ACHIEVEMENT_LIST_GROUPING_PROGRESS
)

// Achievement categories
const (
	AchievementCategoryNone              = C.RC_CLIENT_ACHIEVEMENT_CATEGORY_NONE
	AchievementCategoryCore              = C.RC_CLIENT_ACHIEVEMENT_CATEGORY_CORE
	AchievementCategoryUnofficial        = C.RC_CLIENT_ACHIEVEMENT_CATEGORY_UNOFFICIAL
	AchievementCategoryCoreAndUnofficial = C.RC_CLIENT_ACHIEVEMENT_CATEGORY_CORE_AND_UNOFFICIAL
)

// Achievement bucket types
const (
	AchievementBucketUnknown          = C.RC_CLIENT_ACHIEVEMENT_BUCKET_UNKNOWN
	AchievementBucketLocked           = C.RC_CLIENT_ACHIEVEMENT_BUCKET_LOCKED
	AchievementBucketUnlocked         = C.RC_CLIENT_ACHIEVEMENT_BUCKET_UNLOCKED
	AchievementBucketUnsupported      = C.RC_CLIENT_ACHIEVEMENT_BUCKET_UNSUPPORTED
	AchievementBucketUnofficial       = C.RC_CLIENT_ACHIEVEMENT_BUCKET_UNOFFICIAL
	AchievementBucketRecentlyUnlocked = C.RC_CLIENT_ACHIEVEMENT_BUCKET_RECENTLY_UNLOCKED
	AchievementBucketActiveChallenge  = C.RC_CLIENT_ACHIEVEMENT_BUCKET_ACTIVE_CHALLENGE
	AchievementBucketAlmostThere      = C.RC_CLIENT_ACHIEVEMENT_BUCKET_ALMOST_THERE
	AchievementBucketUnsynced         = C.RC_CLIENT_ACHIEVEMENT_BUCKET_UNSYNCED
)

// Achievement types
const (
	AchievementTypeStandard    = C.RC_CLIENT_ACHIEVEMENT_TYPE_STANDARD
	AchievementTypeMissable    = C.RC_CLIENT_ACHIEVEMENT_TYPE_MISSABLE
	AchievementTypeProgression = C.RC_CLIENT_ACHIEVEMENT_TYPE_PROGRESSION
	AchievementTypeWin         = C.RC_CLIENT_ACHIEVEMENT_TYPE_WIN
)

// Achievement unlock states
const (
	AchievementUnlockedNone     = C.RC_CLIENT_ACHIEVEMENT_UNLOCKED_NONE
	AchievementUnlockedSoftcore = C.RC_CLIENT_ACHIEVEMENT_UNLOCKED_SOFTCORE
	AchievementUnlockedHardcore = C.RC_CLIENT_ACHIEVEMENT_UNLOCKED_HARDCORE
	AchievementUnlockedBoth     = C.RC_CLIENT_ACHIEVEMENT_UNLOCKED_BOTH
)

// AchievementBucket represents a group of achievements
type AchievementBucket struct {
	Achievements []*Achievement
	Label        string
	SubsetID     uint32
	BucketType   uint8
}

// AchievementList represents a list of achievement buckets
type AchievementList struct {
	Buckets []*AchievementBucket
	ptr     *C.rc_client_achievement_list_t
}

// CreateAchievementList creates a list of achievements matching the specified category and grouping.
// Category should be one of AchievementCategoryCore, AchievementCategoryUnofficial, or AchievementCategoryCoreAndUnofficial.
// Grouping should be AchievementListGroupingLockState or AchievementListGroupingProgress.
// The returned list must be destroyed by calling Destroy() when no longer needed.
func (c *Client) CreateAchievementList(category, grouping int) *AchievementList {
	ptr := C.rc_client_create_achievement_list(c.ptr, C.int(category), C.int(grouping))
	if ptr == nil {
		return nil
	}

	list := &AchievementList{ptr: ptr}

	// Convert C buckets to Go
	if ptr.num_buckets > 0 {
		list.Buckets = make([]*AchievementBucket, ptr.num_buckets)
		buckets := unsafe.Slice(ptr.buckets, ptr.num_buckets)

		for i := range buckets {
			bucket := &AchievementBucket{
				Label:      C.GoString(buckets[i].label),
				SubsetID:   uint32(buckets[i].subset_id),
				BucketType: uint8(buckets[i].bucket_type),
			}

			if buckets[i].num_achievements > 0 {
				bucket.Achievements = make([]*Achievement, buckets[i].num_achievements)
				achievements := unsafe.Slice(buckets[i].achievements, buckets[i].num_achievements)

				for j := range achievements {
					bucket.Achievements[j] = achievementFromC(achievements[j])
				}
			}

			list.Buckets[i] = bucket
		}
	}

	return list
}

// Destroy releases resources associated with the achievement list.
func (l *AchievementList) Destroy() {
	if l.ptr != nil {
		C.rc_client_destroy_achievement_list(l.ptr)
		l.ptr = nil
	}
}

// GetAllAchievements returns a flat slice of all achievements across all buckets.
func (l *AchievementList) GetAllAchievements() []*Achievement {
	var all []*Achievement
	for _, bucket := range l.Buckets {
		all = append(all, bucket.Achievements...)
	}
	return all
}

// GetAchievementImageURL returns the URL for an achievement's badge image.
// Pass AchievementStateUnlocked for the unlocked badge, or AchievementStateLocked for the locked/grayscale badge.
func (c *Client) GetAchievementImageURL(achievement *Achievement, state int) string {
	ach := C.rc_client_get_achievement_info(c.ptr, C.uint32_t(achievement.ID))
	if ach == nil {
		return ""
	}

	buffer := make([]byte, 256)
	result := C.rc_client_achievement_get_image_url(ach, C.int(state), (*C.char)(unsafe.Pointer(&buffer[0])), C.size_t(len(buffer)))
	if result != C.RC_OK {
		return ""
	}

	// Find null terminator
	for i, b := range buffer {
		if b == 0 {
			return string(buffer[:i])
		}
	}
	return string(buffer)
}

// UserGameSummary contains summary information about the user's progress in the current game.
type UserGameSummary struct {
	NumCoreAchievements        uint32
	NumUnofficialAchievements  uint32
	NumUnlockedAchievements    uint32
	NumUnsupportedAchievements uint32
	PointsCore                 uint32
	PointsUnlocked             uint32
}

// GetUserGameSummary returns a summary of the user's progress in the current game.
// This is used for displaying "You have unlocked X of Y achievements" messages.
func (c *Client) GetUserGameSummary() *UserGameSummary {
	var summary C.rc_client_user_game_summary_t
	C.rc_client_get_user_game_summary(c.ptr, &summary)

	return &UserGameSummary{
		NumCoreAchievements:        uint32(summary.num_core_achievements),
		NumUnofficialAchievements:  uint32(summary.num_unofficial_achievements),
		NumUnlockedAchievements:    uint32(summary.num_unlocked_achievements),
		NumUnsupportedAchievements: uint32(summary.num_unsupported_achievements),
		PointsCore:                 uint32(summary.points_core),
		PointsUnlocked:             uint32(summary.points_unlocked),
	}
}
