package rcheevos

import (
	"testing"
)

func TestAchievementConstants(t *testing.T) {
	// Categories
	if AchievementCategoryNone != 0 {
		t.Errorf("AchievementCategoryNone = %d, want 0", AchievementCategoryNone)
	}
	if AchievementCategoryCore != 1 {
		t.Errorf("AchievementCategoryCore = %d, want 1", AchievementCategoryCore)
	}
	if AchievementCategoryUnofficial != 2 {
		t.Errorf("AchievementCategoryUnofficial = %d, want 2", AchievementCategoryUnofficial)
	}
	if AchievementCategoryCoreAndUnofficial != 3 {
		t.Errorf("AchievementCategoryCoreAndUnofficial = %d, want 3", AchievementCategoryCoreAndUnofficial)
	}

	// Grouping
	if AchievementListGroupingLockState != 0 {
		t.Errorf("AchievementListGroupingLockState = %d, want 0", AchievementListGroupingLockState)
	}
	if AchievementListGroupingProgress != 1 {
		t.Errorf("AchievementListGroupingProgress = %d, want 1", AchievementListGroupingProgress)
	}

	// Bucket types
	if AchievementBucketUnknown != 0 {
		t.Errorf("AchievementBucketUnknown = %d, want 0", AchievementBucketUnknown)
	}
	if AchievementBucketLocked != 1 {
		t.Errorf("AchievementBucketLocked = %d, want 1", AchievementBucketLocked)
	}
	if AchievementBucketUnlocked != 2 {
		t.Errorf("AchievementBucketUnlocked = %d, want 2", AchievementBucketUnlocked)
	}

	// Types
	if AchievementTypeStandard != 0 {
		t.Errorf("AchievementTypeStandard = %d, want 0", AchievementTypeStandard)
	}
	if AchievementTypeMissable != 1 {
		t.Errorf("AchievementTypeMissable = %d, want 1", AchievementTypeMissable)
	}
	if AchievementTypeProgression != 2 {
		t.Errorf("AchievementTypeProgression = %d, want 2", AchievementTypeProgression)
	}
	if AchievementTypeWin != 3 {
		t.Errorf("AchievementTypeWin = %d, want 3", AchievementTypeWin)
	}

	// Unlock states
	if AchievementUnlockedNone != 0 {
		t.Errorf("AchievementUnlockedNone = %d, want 0", AchievementUnlockedNone)
	}
	if AchievementUnlockedSoftcore != 1 {
		t.Errorf("AchievementUnlockedSoftcore = %d, want 1", AchievementUnlockedSoftcore)
	}
	if AchievementUnlockedHardcore != 2 {
		t.Errorf("AchievementUnlockedHardcore = %d, want 2", AchievementUnlockedHardcore)
	}
	if AchievementUnlockedBoth != 3 {
		t.Errorf("AchievementUnlockedBoth = %d, want 3", AchievementUnlockedBoth)
	}
}

func TestAchievementListDestroy(t *testing.T) {
	// Test that Destroy is safe to call on nil
	var list *AchievementList
	// Should not panic
	if list != nil {
		list.Destroy()
	}

	// Test empty list
	list = &AchievementList{}
	list.Destroy() // Should not panic
	list.Destroy() // Should be safe to call twice
}

func TestAchievementListGetAllAchievements(t *testing.T) {
	// Test with nil buckets
	list := &AchievementList{
		Buckets: nil,
	}
	all := list.GetAllAchievements()
	if all != nil {
		t.Error("GetAllAchievements should return nil for nil buckets")
	}

	// Test with empty buckets
	list = &AchievementList{
		Buckets: []*AchievementBucket{},
	}
	all = list.GetAllAchievements()
	if all != nil {
		t.Error("GetAllAchievements should return nil for empty buckets")
	}

	// Test with buckets containing achievements
	list = &AchievementList{
		Buckets: []*AchievementBucket{
			{
				Label: "Locked",
				Achievements: []*Achievement{
					{ID: 1, Title: "First"},
					{ID: 2, Title: "Second"},
				},
			},
			{
				Label: "Unlocked",
				Achievements: []*Achievement{
					{ID: 3, Title: "Third"},
				},
			},
		},
	}

	all = list.GetAllAchievements()
	if len(all) != 3 {
		t.Errorf("GetAllAchievements returned %d achievements, want 3", len(all))
	}

	// Verify order
	if all[0].ID != 1 || all[1].ID != 2 || all[2].ID != 3 {
		t.Error("GetAllAchievements returned achievements in wrong order")
	}
}

func TestUserGameSummaryFields(t *testing.T) {
	summary := &UserGameSummary{
		NumCoreAchievements:        10,
		NumUnofficialAchievements:  5,
		NumUnlockedAchievements:    7,
		NumUnsupportedAchievements: 1,
		PointsCore:                 100,
		PointsUnlocked:             70,
	}

	if summary.NumCoreAchievements != 10 {
		t.Errorf("NumCoreAchievements = %d, want 10", summary.NumCoreAchievements)
	}
	if summary.NumUnofficialAchievements != 5 {
		t.Errorf("NumUnofficialAchievements = %d, want 5", summary.NumUnofficialAchievements)
	}
	if summary.NumUnlockedAchievements != 7 {
		t.Errorf("NumUnlockedAchievements = %d, want 7", summary.NumUnlockedAchievements)
	}
	if summary.NumUnsupportedAchievements != 1 {
		t.Errorf("NumUnsupportedAchievements = %d, want 1", summary.NumUnsupportedAchievements)
	}
	if summary.PointsCore != 100 {
		t.Errorf("PointsCore = %d, want 100", summary.PointsCore)
	}
	if summary.PointsUnlocked != 70 {
		t.Errorf("PointsUnlocked = %d, want 70", summary.PointsUnlocked)
	}
}

func TestAchievementFields(t *testing.T) {
	ach := &Achievement{
		ID:               12345,
		Title:            "Test Achievement",
		Description:      "Test Description",
		BadgeName:        "12345",
		Points:           10,
		State:            AchievementStateActive,
		Category:         AchievementCategoryCore,
		Bucket:           AchievementBucketLocked,
		Unlocked:         AchievementUnlockedNone,
		Type:             AchievementTypeStandard,
		MeasuredProgress: "5/10",
		MeasuredPercent:  50.0,
		Rarity:           25.5,
		RarityHardcore:   10.2,
		UnlockTime:       0,
	}

	if ach.ID != 12345 {
		t.Errorf("ID = %d, want 12345", ach.ID)
	}
	if ach.Title != "Test Achievement" {
		t.Errorf("Title = %q, want %q", ach.Title, "Test Achievement")
	}
	if ach.Points != 10 {
		t.Errorf("Points = %d, want 10", ach.Points)
	}
	if ach.MeasuredPercent != 50.0 {
		t.Errorf("MeasuredPercent = %f, want 50.0", ach.MeasuredPercent)
	}
	if ach.Rarity != 25.5 {
		t.Errorf("Rarity = %f, want 25.5", ach.Rarity)
	}
}
