package rcheevos

import (
	"testing"
	"time"
)

func TestLeaderboardConstants(t *testing.T) {
	// States
	if LeaderboardStateInactive != 0 {
		t.Errorf("LeaderboardStateInactive = %d, want 0", LeaderboardStateInactive)
	}
	if LeaderboardStateActive != 1 {
		t.Errorf("LeaderboardStateActive = %d, want 1", LeaderboardStateActive)
	}
	if LeaderboardStateTracking != 2 {
		t.Errorf("LeaderboardStateTracking = %d, want 2", LeaderboardStateTracking)
	}
	if LeaderboardStateDisabled != 3 {
		t.Errorf("LeaderboardStateDisabled = %d, want 3", LeaderboardStateDisabled)
	}

	// Formats
	if LeaderboardFormatTime != 0 {
		t.Errorf("LeaderboardFormatTime = %d, want 0", LeaderboardFormatTime)
	}
	if LeaderboardFormatScore != 1 {
		t.Errorf("LeaderboardFormatScore = %d, want 1", LeaderboardFormatScore)
	}
	if LeaderboardFormatValue != 2 {
		t.Errorf("LeaderboardFormatValue = %d, want 2", LeaderboardFormatValue)
	}

	// Grouping
	if LeaderboardListGroupingNone != 0 {
		t.Errorf("LeaderboardListGroupingNone = %d, want 0", LeaderboardListGroupingNone)
	}
	if LeaderboardListGroupingTracking != 1 {
		t.Errorf("LeaderboardListGroupingTracking = %d, want 1", LeaderboardListGroupingTracking)
	}

	// Bucket types
	if LeaderboardBucketUnknown != 0 {
		t.Errorf("LeaderboardBucketUnknown = %d, want 0", LeaderboardBucketUnknown)
	}
	if LeaderboardBucketInactive != 1 {
		t.Errorf("LeaderboardBucketInactive = %d, want 1", LeaderboardBucketInactive)
	}
	if LeaderboardBucketActive != 2 {
		t.Errorf("LeaderboardBucketActive = %d, want 2", LeaderboardBucketActive)
	}
	if LeaderboardBucketUnsupported != 3 {
		t.Errorf("LeaderboardBucketUnsupported = %d, want 3", LeaderboardBucketUnsupported)
	}
	if LeaderboardBucketAll != 4 {
		t.Errorf("LeaderboardBucketAll = %d, want 4", LeaderboardBucketAll)
	}
}

func TestLeaderboardFields(t *testing.T) {
	lb := &Leaderboard{
		ID:            12345,
		Title:         "Best Time",
		Description:   "Complete the level as fast as possible",
		TrackerValue:  "1:23.45",
		State:         LeaderboardStateTracking,
		Format:        LeaderboardFormatTime,
		LowerIsBetter: true,
	}

	if lb.ID != 12345 {
		t.Errorf("ID = %d, want 12345", lb.ID)
	}
	if lb.Title != "Best Time" {
		t.Errorf("Title = %q, want %q", lb.Title, "Best Time")
	}
	if lb.TrackerValue != "1:23.45" {
		t.Errorf("TrackerValue = %q, want %q", lb.TrackerValue, "1:23.45")
	}
	if lb.State != LeaderboardStateTracking {
		t.Errorf("State = %d, want %d", lb.State, LeaderboardStateTracking)
	}
	if lb.Format != LeaderboardFormatTime {
		t.Errorf("Format = %d, want %d", lb.Format, LeaderboardFormatTime)
	}
	if !lb.LowerIsBetter {
		t.Error("LowerIsBetter should be true")
	}
}

func TestLeaderboardListDestroy(t *testing.T) {
	// Test that Destroy is safe to call on nil
	var list *LeaderboardList
	if list != nil {
		list.Destroy()
	}

	// Test empty list
	list = &LeaderboardList{}
	list.Destroy() // Should not panic
	list.Destroy() // Should be safe to call twice
}

func TestLeaderboardListGetAllLeaderboards(t *testing.T) {
	// Test with nil buckets
	list := &LeaderboardList{
		Buckets: nil,
	}
	all := list.GetAllLeaderboards()
	if all != nil {
		t.Error("GetAllLeaderboards should return nil for nil buckets")
	}

	// Test with empty buckets
	list = &LeaderboardList{
		Buckets: []*LeaderboardBucket{},
	}
	all = list.GetAllLeaderboards()
	if all != nil {
		t.Error("GetAllLeaderboards should return nil for empty buckets")
	}

	// Test with buckets containing leaderboards
	list = &LeaderboardList{
		Buckets: []*LeaderboardBucket{
			{
				Label: "Active",
				Leaderboards: []*Leaderboard{
					{ID: 1, Title: "First"},
					{ID: 2, Title: "Second"},
				},
			},
			{
				Label: "Inactive",
				Leaderboards: []*Leaderboard{
					{ID: 3, Title: "Third"},
				},
			},
		},
	}

	all = list.GetAllLeaderboards()
	if len(all) != 3 {
		t.Errorf("GetAllLeaderboards returned %d leaderboards, want 3", len(all))
	}
}

func TestLeaderboardEntryFields(t *testing.T) {
	now := time.Now()
	entry := &LeaderboardEntry{
		User:      "TestUser",
		Display:   "1:23.45",
		Submitted: now,
		Rank:      1,
		Index:     0,
	}

	if entry.User != "TestUser" {
		t.Errorf("User = %q, want %q", entry.User, "TestUser")
	}
	if entry.Display != "1:23.45" {
		t.Errorf("Display = %q, want %q", entry.Display, "1:23.45")
	}
	if entry.Rank != 1 {
		t.Errorf("Rank = %d, want 1", entry.Rank)
	}
	if entry.Index != 0 {
		t.Errorf("Index = %d, want 0", entry.Index)
	}
}

func TestLeaderboardEntryGetUserImageURL(t *testing.T) {
	entry := &LeaderboardEntry{
		User: "TestUser",
	}

	url := entry.GetUserImageURL()
	expected := "https://media.retroachievements.org/UserPic/TestUser.png"
	if url != expected {
		t.Errorf("GetUserImageURL() = %q, want %q", url, expected)
	}
}

func TestLeaderboardEntryListDestroy(t *testing.T) {
	// Test that Destroy is safe to call on nil
	var list *LeaderboardEntryList
	if list != nil {
		list.Destroy()
	}

	// Test empty list
	list = &LeaderboardEntryList{}
	list.Destroy() // Should not panic
	list.Destroy() // Should be safe to call twice
}

func TestLeaderboardEntryListFields(t *testing.T) {
	list := &LeaderboardEntryList{
		Entries: []*LeaderboardEntry{
			{User: "User1", Rank: 1},
			{User: "User2", Rank: 2},
			{User: "User3", Rank: 3},
		},
		TotalEntries: 100,
		UserIndex:    1,
	}

	if len(list.Entries) != 3 {
		t.Errorf("len(Entries) = %d, want 3", len(list.Entries))
	}
	if list.TotalEntries != 100 {
		t.Errorf("TotalEntries = %d, want 100", list.TotalEntries)
	}
	if list.UserIndex != 1 {
		t.Errorf("UserIndex = %d, want 1", list.UserIndex)
	}
}

func TestLeaderboardTrackerFields(t *testing.T) {
	tracker := &LeaderboardTracker{
		ID:      123,
		Display: "1:23.45",
	}

	if tracker.ID != 123 {
		t.Errorf("ID = %d, want 123", tracker.ID)
	}
	if tracker.Display != "1:23.45" {
		t.Errorf("Display = %q, want %q", tracker.Display, "1:23.45")
	}
}

func TestLeaderboardScoreboardFields(t *testing.T) {
	scoreboard := &LeaderboardScoreboard{
		LeaderboardID:  123,
		SubmittedScore: "1:23.45",
		BestScore:      "1:20.00",
		NewRank:        5,
		NumEntries:     100,
		TopEntries: []*LeaderboardScoreboardEntry{
			{Username: "Top1", Rank: 1, Score: "1:00.00"},
			{Username: "Top2", Rank: 2, Score: "1:05.00"},
		},
	}

	if scoreboard.LeaderboardID != 123 {
		t.Errorf("LeaderboardID = %d, want 123", scoreboard.LeaderboardID)
	}
	if scoreboard.SubmittedScore != "1:23.45" {
		t.Errorf("SubmittedScore = %q, want %q", scoreboard.SubmittedScore, "1:23.45")
	}
	if scoreboard.NewRank != 5 {
		t.Errorf("NewRank = %d, want 5", scoreboard.NewRank)
	}
	if len(scoreboard.TopEntries) != 2 {
		t.Errorf("len(TopEntries) = %d, want 2", len(scoreboard.TopEntries))
	}
}
