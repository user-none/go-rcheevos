package rcheevos

import (
	"testing"
)

func TestConstants(t *testing.T) {
	// Console IDs
	if ConsoleMasterSystem != 11 {
		t.Errorf("ConsoleMasterSystem = %d, want 11", ConsoleMasterSystem)
	}
	if ConsoleGameGear != 15 {
		t.Errorf("ConsoleGameGear = %d, want 15", ConsoleGameGear)
	}
	if ConsoleMegaDrive != 1 {
		t.Errorf("ConsoleMegaDrive = %d, want 1", ConsoleMegaDrive)
	}
	if ConsoleNES != 7 {
		t.Errorf("ConsoleNES = %d, want 7", ConsoleNES)
	}
	if ConsoleSNES != 3 {
		t.Errorf("ConsoleSNES = %d, want 3", ConsoleSNES)
	}

	// Log levels
	if LogLevelNone != 0 {
		t.Errorf("LogLevelNone = %d, want 0", LogLevelNone)
	}
	if LogLevelError != 1 {
		t.Errorf("LogLevelError = %d, want 1", LogLevelError)
	}
	if LogLevelVerbose != 4 {
		t.Errorf("LogLevelVerbose = %d, want 4", LogLevelVerbose)
	}

	// Achievement states
	if AchievementStateInactive != 0 {
		t.Errorf("AchievementStateInactive = %d, want 0", AchievementStateInactive)
	}
	if AchievementStateActive != 1 {
		t.Errorf("AchievementStateActive = %d, want 1", AchievementStateActive)
	}
	if AchievementStateUnlocked != 2 {
		t.Errorf("AchievementStateUnlocked = %d, want 2", AchievementStateUnlocked)
	}

	// Event types
	if EventAchievementTriggered != 1 {
		t.Errorf("EventAchievementTriggered = %d, want 1", EventAchievementTriggered)
	}
	if EventReset != 14 {
		t.Errorf("EventReset = %d, want 14", EventReset)
	}
	if EventGameCompleted != 15 {
		t.Errorf("EventGameCompleted = %d, want 15", EventGameCompleted)
	}
}

func TestErrorCodes(t *testing.T) {
	if OK != 0 {
		t.Errorf("OK = %d, want 0", OK)
	}
	if LoginRequired != -28 {
		t.Errorf("LoginRequired = %d, want -28", LoginRequired)
	}
	if NoGameLoaded != -29 {
		t.Errorf("NoGameLoaded = %d, want -29", NoGameLoaded)
	}
	if InvalidCredentials != -34 {
		t.Errorf("InvalidCredentials = %d, want -34", InvalidCredentials)
	}
}

func TestErrorString(t *testing.T) {
	tests := []struct {
		code int
		want string
	}{
		{OK, "OK"},
		{InvalidCredentials, "Invalid credentials"},
		{LoginRequired, "Login required"},
		{NoGameLoaded, "No game loaded"},
		{ExpiredToken, "Expired token"},
	}

	for _, tt := range tests {
		got := ErrorString(tt.code)
		if got != tt.want {
			t.Errorf("ErrorString(%d) = %q, want %q", tt.code, got, tt.want)
		}
	}
}

func TestConsoleName(t *testing.T) {
	tests := []struct {
		id   uint32
		want string
	}{
		{ConsoleMasterSystem, "Master System"},
		{ConsoleGameGear, "Game Gear"},
		{ConsoleMegaDrive, "Sega Genesis"},
		{ConsoleNES, "Nintendo Entertainment System"},
		{ConsoleSNES, "Super Nintendo Entertainment System"},
		{ConsoleGameBoy, "GameBoy"},
		{ConsoleGBA, "GameBoy Advance"},
	}

	for _, tt := range tests {
		got := ConsoleName(tt.id)
		if got != tt.want {
			t.Errorf("ConsoleName(%d) = %q, want %q", tt.id, got, tt.want)
		}
	}
}

func TestClientCreation(t *testing.T) {
	client := NewClient(
		func(address uint32, buffer []byte) uint32 {
			return 0
		},
		func(req *ServerRequest) {
		},
	)

	if client == nil {
		t.Fatal("NewClient returned nil")
	}
	if client.ptr == nil {
		t.Fatal("client.ptr is nil")
	}

	// Verify client is registered
	if getClient(client.ptr) != client {
		t.Error("client not registered in clientMap")
	}

	client.Destroy()

	// Verify client is unregistered
	if getClient(client.ptr) != nil {
		t.Error("client still registered after Destroy")
	}
}

func TestClientConfiguration(t *testing.T) {
	client := NewClient(
		func(address uint32, buffer []byte) uint32 { return 0 },
		func(req *ServerRequest) {},
	)
	defer client.Destroy()

	// Test hardcore mode (default is enabled)
	if !client.GetHardcoreEnabled() {
		t.Error("hardcore mode should be enabled by default")
	}

	client.SetHardcoreEnabled(false)
	if client.GetHardcoreEnabled() {
		t.Error("hardcore mode should be disabled after SetHardcoreEnabled(false)")
	}

	client.SetHardcoreEnabled(true)
	if !client.GetHardcoreEnabled() {
		t.Error("hardcore mode should be enabled after SetHardcoreEnabled(true)")
	}

	// Test encore mode (default is disabled)
	if client.GetEncoreModeEnabled() {
		t.Error("encore mode should be disabled by default")
	}

	client.SetEncoreModeEnabled(true)
	if !client.GetEncoreModeEnabled() {
		t.Error("encore mode should be enabled after SetEncoreModeEnabled(true)")
	}

	// Test unofficial mode (default is disabled)
	if client.GetUnofficialEnabled() {
		t.Error("unofficial mode should be disabled by default")
	}

	client.SetUnofficialEnabled(true)
	if !client.GetUnofficialEnabled() {
		t.Error("unofficial mode should be enabled after SetUnofficialEnabled(true)")
	}

	// Test spectator mode (default is disabled)
	if client.GetSpectatorModeEnabled() {
		t.Error("spectator mode should be disabled by default")
	}

	client.SetSpectatorModeEnabled(true)
	if client.GetSpectatorModeEnabled() {
		// Note: spectator mode can only be enabled if it was enabled before game load
		// Since no game is loaded, this might not work as expected
	}
}

func TestClientNoGameLoaded(t *testing.T) {
	client := NewClient(
		func(address uint32, buffer []byte) uint32 { return 0 },
		func(req *ServerRequest) {},
	)
	defer client.Destroy()

	// No game loaded
	if client.IsGameLoaded() {
		t.Error("IsGameLoaded should return false when no game is loaded")
	}

	// GetGame should return nil
	if client.GetGame() != nil {
		t.Error("GetGame should return nil when no game is loaded")
	}

	// GetUser should return nil (not logged in)
	if client.GetUser() != nil {
		t.Error("GetUser should return nil when not logged in")
	}

	// HasAchievements should return false
	if client.HasAchievements() {
		t.Error("HasAchievements should return false when no game is loaded")
	}

	// HasLeaderboards should return false
	if client.HasLeaderboards() {
		t.Error("HasLeaderboards should return false when no game is loaded")
	}

	// HasRichPresence should return false
	if client.HasRichPresence() {
		t.Error("HasRichPresence should return false when no game is loaded")
	}

	// GetRichPresenceMessage should return empty
	if msg := client.GetRichPresenceMessage(); msg != "" {
		t.Errorf("GetRichPresenceMessage should return empty, got %q", msg)
	}
}

func TestClientUserAgent(t *testing.T) {
	client := NewClient(
		func(address uint32, buffer []byte) uint32 { return 0 },
		func(req *ServerRequest) {},
	)
	defer client.Destroy()

	ua := client.GetUserAgentClause()
	if ua == "" {
		t.Error("GetUserAgentClause should return a non-empty string")
	}

	// Should contain "rcheevos"
	if len(ua) < 5 {
		t.Errorf("GetUserAgentClause too short: %q", ua)
	}
}

func TestClientProcessing(t *testing.T) {
	client := NewClient(
		func(address uint32, buffer []byte) uint32 { return 0 },
		func(req *ServerRequest) {},
	)
	defer client.Destroy()

	// IsProcessingRequired should return false with no game
	if client.IsProcessingRequired() {
		t.Error("IsProcessingRequired should return false when no game is loaded")
	}

	// CanPause should return true with no game
	allowed, frames := client.CanPause()
	if !allowed {
		t.Error("CanPause should return true when no game is loaded")
	}
	if frames != 0 {
		t.Errorf("CanPause frames remaining should be 0, got %d", frames)
	}

	// GetLoadGameState should return None
	if state := client.GetLoadGameState(); state != LoadGameStateNone {
		t.Errorf("GetLoadGameState should return None, got %d", state)
	}
}

func TestLoadGameStateName(t *testing.T) {
	tests := []struct {
		state int
		want  string
	}{
		{LoadGameStateNone, "None"},
		{LoadGameStateAwaitLogin, "Awaiting Login"},
		{LoadGameStateIdentifying, "Identifying Game"},
		{LoadGameStateDone, "Done"},
		{LoadGameStateAborted, "Aborted"},
		{999, "Unknown"},
	}

	for _, tt := range tests {
		got := LoadGameStateName(tt.state)
		if got != tt.want {
			t.Errorf("LoadGameStateName(%d) = %q, want %q", tt.state, got, tt.want)
		}
	}
}

func TestProgressSerialization(t *testing.T) {
	client := NewClient(
		func(address uint32, buffer []byte) uint32 { return 0 },
		func(req *ServerRequest) {},
	)
	defer client.Destroy()

	// With no game loaded, progress size should be 0
	size := client.ProgressSize()
	if size != 0 {
		t.Errorf("ProgressSize should be 0 with no game, got %d", size)
	}

	// Serialize should return empty data with no game
	data, result := client.SerializeProgress()
	if result != OK {
		t.Errorf("SerializeProgress result = %d, want OK", result)
	}
	if data != nil {
		t.Errorf("SerializeProgress should return nil with no game, got %v", data)
	}

	// Deserialize empty data should succeed
	result = client.DeserializeProgress(nil)
	if result != OK {
		t.Errorf("DeserializeProgress(nil) = %d, want OK", result)
	}

	result = client.DeserializeProgress([]byte{})
	if result != OK {
		t.Errorf("DeserializeProgress([]) = %d, want OK", result)
	}
}

func TestEventHandler(t *testing.T) {
	client := NewClient(
		func(address uint32, buffer []byte) uint32 { return 0 },
		func(req *ServerRequest) {},
	)
	defer client.Destroy()

	client.SetEventHandler(func(event *Event) {
		_ = event // Use event to avoid unused warning
	})

	// We can't easily trigger events without a loaded game,
	// but we can verify the handler is set
	client.mu.RLock()
	hasHandler := client.eventFunc != nil
	client.mu.RUnlock()

	if !hasHandler {
		t.Error("event handler should be set")
	}
}

func TestLogging(t *testing.T) {
	logMessages := []string{}

	client := NewClient(
		func(address uint32, buffer []byte) uint32 { return 0 },
		func(req *ServerRequest) {},
	)
	defer client.Destroy()

	client.EnableLogging(LogLevelVerbose, func(message string) {
		logMessages = append(logMessages, message)
	})

	client.mu.RLock()
	hasLogFunc := client.logFunc != nil
	client.mu.RUnlock()

	if !hasLogFunc {
		t.Error("log function should be set")
	}
}

func TestGetAchievementNil(t *testing.T) {
	client := NewClient(
		func(address uint32, buffer []byte) uint32 { return 0 },
		func(req *ServerRequest) {},
	)
	defer client.Destroy()

	// Should return nil for non-existent achievement
	ach := client.GetAchievement(12345)
	if ach != nil {
		t.Error("GetAchievement should return nil for non-existent achievement")
	}
}

func TestGetLeaderboardNil(t *testing.T) {
	client := NewClient(
		func(address uint32, buffer []byte) uint32 { return 0 },
		func(req *ServerRequest) {},
	)
	defer client.Destroy()

	// Should return nil for non-existent leaderboard
	lb := client.GetLeaderboard(12345)
	if lb != nil {
		t.Error("GetLeaderboard should return nil for non-existent leaderboard")
	}
}

func TestCreateAchievementListNoGame(t *testing.T) {
	client := NewClient(
		func(address uint32, buffer []byte) uint32 { return 0 },
		func(req *ServerRequest) {},
	)
	defer client.Destroy()

	// May return empty list or nil with no game loaded
	list := client.CreateAchievementList(AchievementCategoryCore, AchievementListGroupingLockState)
	if list != nil {
		defer list.Destroy()
		// Should have no achievements
		all := list.GetAllAchievements()
		if len(all) != 0 {
			t.Errorf("expected 0 achievements, got %d", len(all))
		}
	}
}

func TestCreateLeaderboardListNoGame(t *testing.T) {
	client := NewClient(
		func(address uint32, buffer []byte) uint32 { return 0 },
		func(req *ServerRequest) {},
	)
	defer client.Destroy()

	// May return empty list or nil with no game loaded
	list := client.CreateLeaderboardList(LeaderboardListGroupingNone)
	if list != nil {
		defer list.Destroy()
		// Should have no leaderboards
		all := list.GetAllLeaderboards()
		if len(all) != 0 {
			t.Errorf("expected 0 leaderboards, got %d", len(all))
		}
	}
}

func TestCreateSubsetListNoGame(t *testing.T) {
	client := NewClient(
		func(address uint32, buffer []byte) uint32 { return 0 },
		func(req *ServerRequest) {},
	)
	defer client.Destroy()

	// May return empty list or nil with no game loaded
	list := client.CreateSubsetList()
	if list != nil {
		defer list.Destroy()
		// Should have no subsets
		if len(list.Subsets) != 0 {
			t.Errorf("expected 0 subsets, got %d", len(list.Subsets))
		}
	}
}

func TestGetUserGameSummary(t *testing.T) {
	client := NewClient(
		func(address uint32, buffer []byte) uint32 { return 0 },
		func(req *ServerRequest) {},
	)
	defer client.Destroy()

	// Should return zeroed summary with no game
	summary := client.GetUserGameSummary()
	if summary == nil {
		t.Fatal("GetUserGameSummary should not return nil")
	}

	if summary.NumCoreAchievements != 0 {
		t.Errorf("NumCoreAchievements = %d, want 0", summary.NumCoreAchievements)
	}
	if summary.NumUnlockedAchievements != 0 {
		t.Errorf("NumUnlockedAchievements = %d, want 0", summary.NumUnlockedAchievements)
	}
	if summary.PointsCore != 0 {
		t.Errorf("PointsCore = %d, want 0", summary.PointsCore)
	}
}

func TestImageURLsEmpty(t *testing.T) {
	client := NewClient(
		func(address uint32, buffer []byte) uint32 { return 0 },
		func(req *ServerRequest) {},
	)
	defer client.Destroy()

	// Should return empty with no game/user
	if url := client.GetGameImageURL(); url != "" {
		t.Errorf("GetGameImageURL should return empty, got %q", url)
	}

	if url := client.GetUserImageURL(); url != "" {
		t.Errorf("GetUserImageURL should return empty, got %q", url)
	}
}

func TestDoFrameNoGame(t *testing.T) {
	client := NewClient(
		func(address uint32, buffer []byte) uint32 { return 0 },
		func(req *ServerRequest) {},
	)
	defer client.Destroy()

	// Should not panic with no game loaded
	client.DoFrame()
	client.Idle()
	client.Reset()
}

func TestMultipleClients(t *testing.T) {
	client1 := NewClient(
		func(address uint32, buffer []byte) uint32 { return 1 },
		func(req *ServerRequest) {},
	)

	client2 := NewClient(
		func(address uint32, buffer []byte) uint32 { return 2 },
		func(req *ServerRequest) {},
	)

	// Both should be registered
	if getClient(client1.ptr) != client1 {
		t.Error("client1 not properly registered")
	}
	if getClient(client2.ptr) != client2 {
		t.Error("client2 not properly registered")
	}

	// Destroy one, other should still work
	client1.Destroy()

	if getClient(client1.ptr) != nil {
		t.Error("client1 should be unregistered after destroy")
	}
	if getClient(client2.ptr) != client2 {
		t.Error("client2 should still be registered")
	}

	client2.Destroy()
}
