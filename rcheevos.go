package rcheevos

/*
#cgo CFLAGS: -I${SRCDIR}/rcheevos/include -I${SRCDIR}/rcheevos/src -DRC_DISABLE_LUA
#cgo CFLAGS: -DRC_CLIENT_SUPPORTS_HASH
#cgo LDFLAGS: -lm

#include <stdlib.h>
#include <stdint.h>
#include "rc_client.h"
#include "rc_api_request.h"
#include "rc_hash.h"
#include "rc_consoles.h"
#include "rc_error.h"

// Forward declarations for Go callback bridges
extern uint32_t goReadMemoryCallback(uint32_t address, uint8_t* buffer, uint32_t num_bytes, rc_client_t* client);
extern void goServerCallCallback(const rc_api_request_t* request, rc_client_server_callback_t callback, void* callback_data, rc_client_t* client);
extern void goEventCallback(const rc_client_event_t* event, rc_client_t* client);
extern void goLogCallback(const char* message, const rc_client_t* client);
extern void goLoginCallback(int result, const char* error_message, rc_client_t* client, void* userdata);
extern void goLoadGameCallback(int result, const char* error_message, rc_client_t* client, void* userdata);
extern void goLeaderboardEntriesCallback(int result, const char* error_message, rc_client_leaderboard_entry_list_t* list, rc_client_t* client, void* userdata);
extern uint64_t goGetTimeMillisecsCallback(const rc_client_t* client);
extern void goChangeMediaCallback(int result, const char* error_message, rc_client_t* client, void* userdata);
extern void goFetchHashLibraryCallback(int result, const char* error_message, rc_client_hash_library_t* list, rc_client_t* client, void* userdata);
extern void goFetchGameTitlesCallback(int result, const char* error_message, rc_client_game_title_list_t* list, rc_client_t* client, void* userdata);
extern void goFetchAllUserProgressCallback(int result, const char* error_message, rc_client_all_user_progress_t* list, rc_client_t* client, void* userdata);

// C wrapper functions that call the Go callbacks
static uint32_t c_read_memory_callback(uint32_t address, uint8_t* buffer, uint32_t num_bytes, rc_client_t* client) {
    return goReadMemoryCallback(address, buffer, num_bytes, client);
}

static void c_server_call_callback(const rc_api_request_t* request, rc_client_server_callback_t callback, void* callback_data, rc_client_t* client) {
    goServerCallCallback(request, callback, callback_data, client);
}

static void c_event_callback(const rc_client_event_t* event, rc_client_t* client) {
    goEventCallback(event, client);
}

static void c_log_callback(const char* message, const rc_client_t* client) {
    goLogCallback(message, client);
}

static void c_login_callback(int result, const char* error_message, rc_client_t* client, void* userdata) {
    goLoginCallback(result, error_message, client, userdata);
}

static void c_load_game_callback(int result, const char* error_message, rc_client_t* client, void* userdata) {
    goLoadGameCallback(result, error_message, client, userdata);
}

static void c_leaderboard_entries_callback(int result, const char* error_message, rc_client_leaderboard_entry_list_t* list, rc_client_t* client, void* userdata) {
    goLeaderboardEntriesCallback(result, error_message, list, client, userdata);
}

// Helper to create client with our C wrapper callbacks
static rc_client_t* create_client_with_callbacks(void) {
    return rc_client_create(c_read_memory_callback, c_server_call_callback);
}

static void set_event_handler(rc_client_t* client) {
    rc_client_set_event_handler(client, c_event_callback);
}

static void enable_logging(rc_client_t* client, int level) {
    rc_client_enable_logging(client, level, c_log_callback);
}

// The begin_ helpers pass the request's handle from Go through rcheevos as
// the callback userdata, so each completion reaches its own request.

static rc_client_async_handle_t* begin_login_with_password(rc_client_t* client, const char* username, const char* password, uintptr_t userdata) {
    return rc_client_begin_login_with_password(client, username, password, c_login_callback, (void*)userdata);
}

static rc_client_async_handle_t* begin_login_with_token(rc_client_t* client, const char* username, const char* token, uintptr_t userdata) {
    return rc_client_begin_login_with_token(client, username, token, c_login_callback, (void*)userdata);
}

static rc_client_async_handle_t* begin_load_game(rc_client_t* client, const char* hash, uintptr_t userdata) {
    return rc_client_begin_load_game(client, hash, c_load_game_callback, (void*)userdata);
}

static rc_client_async_handle_t* begin_identify_and_load_game(rc_client_t* client, uint32_t console_id, const char* file_path, const uint8_t* data, size_t data_size, uintptr_t userdata) {
    return rc_client_begin_identify_and_load_game(client, console_id, file_path, data, data_size, c_load_game_callback, (void*)userdata);
}

static rc_client_async_handle_t* begin_fetch_leaderboard_entries(rc_client_t* client, uint32_t leaderboard_id, uint32_t first_entry, uint32_t count, uintptr_t userdata) {
    return rc_client_begin_fetch_leaderboard_entries(client, leaderboard_id, first_entry, count, c_leaderboard_entries_callback, (void*)userdata);
}

static rc_client_async_handle_t* begin_fetch_leaderboard_entries_around_user(rc_client_t* client, uint32_t leaderboard_id, uint32_t count, uintptr_t userdata) {
    return rc_client_begin_fetch_leaderboard_entries_around_user(client, leaderboard_id, count, c_leaderboard_entries_callback, (void*)userdata);
}

static uint64_t c_get_time_millisecs_callback(const rc_client_t* client) {
    return goGetTimeMillisecsCallback(client);
}

static void set_get_time_millisecs_function(rc_client_t* client) {
    rc_client_set_get_time_millisecs_function(client, c_get_time_millisecs_callback);
}

static void c_change_media_callback(int result, const char* error_message, rc_client_t* client, void* userdata) {
    goChangeMediaCallback(result, error_message, client, userdata);
}

static rc_client_async_handle_t* begin_change_media(rc_client_t* client, const char* hash, uintptr_t userdata) {
    return rc_client_begin_change_media(client, hash, c_change_media_callback, (void*)userdata);
}

static rc_client_async_handle_t* begin_identify_and_change_media(rc_client_t* client, const char* file_path, const uint8_t* data, size_t data_size, uintptr_t userdata) {
    return rc_client_begin_identify_and_change_media(client, file_path, data, data_size, c_change_media_callback, (void*)userdata);
}

static void c_fetch_hash_library_callback(int result, const char* error_message, rc_client_hash_library_t* list, rc_client_t* client, void* userdata) {
    goFetchHashLibraryCallback(result, error_message, list, client, userdata);
}

static void c_fetch_game_titles_callback(int result, const char* error_message, rc_client_game_title_list_t* list, rc_client_t* client, void* userdata) {
    goFetchGameTitlesCallback(result, error_message, list, client, userdata);
}

static void c_fetch_all_user_progress_callback(int result, const char* error_message, rc_client_all_user_progress_t* list, rc_client_t* client, void* userdata) {
    goFetchAllUserProgressCallback(result, error_message, list, client, userdata);
}

static rc_client_async_handle_t* begin_fetch_hash_library(rc_client_t* client, uint32_t console_id, uintptr_t userdata) {
    return rc_client_begin_fetch_hash_library(client, console_id, c_fetch_hash_library_callback, (void*)userdata);
}

static rc_client_async_handle_t* begin_fetch_game_titles(rc_client_t* client, const uint32_t* game_ids, uint32_t num_game_ids, uintptr_t userdata) {
    return rc_client_begin_fetch_game_titles(client, game_ids, num_game_ids, c_fetch_game_titles_callback, (void*)userdata);
}

static rc_client_async_handle_t* begin_fetch_all_user_progress(rc_client_t* client, uint32_t console_id, uintptr_t userdata) {
    return rc_client_begin_fetch_all_user_progress(client, console_id, c_fetch_all_user_progress_callback, (void*)userdata);
}

// Server response helper
static void server_response_callback(rc_client_server_callback_t callback, void* callback_data, const char* body, int body_len, int http_status) {
    rc_api_server_response_t response;
    response.body = body;
    response.body_length = body_len;
    response.http_status_code = http_status;
    callback(&response, callback_data);
}
*/
import "C"
import (
	"bytes"
	"math"
	"runtime/cgo"
	"sync"
	"sync/atomic"
	"unsafe"
)

// Console IDs
const (
	ConsoleUnknown      = C.RC_CONSOLE_UNKNOWN
	ConsoleMegaDrive    = C.RC_CONSOLE_MEGA_DRIVE
	ConsoleNintendo64   = C.RC_CONSOLE_NINTENDO_64
	ConsoleSNES         = C.RC_CONSOLE_SUPER_NINTENDO
	ConsoleGameBoy      = C.RC_CONSOLE_GAMEBOY
	ConsoleGBA          = C.RC_CONSOLE_GAMEBOY_ADVANCE
	ConsoleGBC          = C.RC_CONSOLE_GAMEBOY_COLOR
	ConsoleNES          = C.RC_CONSOLE_NINTENDO
	ConsolePCEngine     = C.RC_CONSOLE_PC_ENGINE
	ConsoleSegaCD       = C.RC_CONSOLE_SEGA_CD
	ConsoleSega32X      = C.RC_CONSOLE_SEGA_32X
	ConsoleMasterSystem = C.RC_CONSOLE_MASTER_SYSTEM
	ConsolePlayStation  = C.RC_CONSOLE_PLAYSTATION
	ConsoleNeoGeoPocket = C.RC_CONSOLE_NEOGEO_POCKET
	ConsoleGameGear     = C.RC_CONSOLE_GAME_GEAR
	ConsoleSG1000       = C.RC_CONSOLE_SG1000
	ConsoleSaturn       = C.RC_CONSOLE_SATURN
)

// Log levels
const (
	LogLevelNone    = C.RC_CLIENT_LOG_LEVEL_NONE
	LogLevelError   = C.RC_CLIENT_LOG_LEVEL_ERROR
	LogLevelWarn    = C.RC_CLIENT_LOG_LEVEL_WARN
	LogLevelInfo    = C.RC_CLIENT_LOG_LEVEL_INFO
	LogLevelVerbose = C.RC_CLIENT_LOG_LEVEL_VERBOSE
)

// Event types
const (
	EventAchievementTriggered             = C.RC_CLIENT_EVENT_ACHIEVEMENT_TRIGGERED
	EventLeaderboardStarted               = C.RC_CLIENT_EVENT_LEADERBOARD_STARTED
	EventLeaderboardFailed                = C.RC_CLIENT_EVENT_LEADERBOARD_FAILED
	EventLeaderboardSubmitted             = C.RC_CLIENT_EVENT_LEADERBOARD_SUBMITTED
	EventAchievementChallengeIndicatorShow = C.RC_CLIENT_EVENT_ACHIEVEMENT_CHALLENGE_INDICATOR_SHOW
	EventAchievementChallengeIndicatorHide = C.RC_CLIENT_EVENT_ACHIEVEMENT_CHALLENGE_INDICATOR_HIDE
	EventAchievementProgressIndicatorShow  = C.RC_CLIENT_EVENT_ACHIEVEMENT_PROGRESS_INDICATOR_SHOW
	EventAchievementProgressIndicatorHide  = C.RC_CLIENT_EVENT_ACHIEVEMENT_PROGRESS_INDICATOR_HIDE
	EventAchievementProgressIndicatorUpdate = C.RC_CLIENT_EVENT_ACHIEVEMENT_PROGRESS_INDICATOR_UPDATE
	EventLeaderboardTrackerShow           = C.RC_CLIENT_EVENT_LEADERBOARD_TRACKER_SHOW
	EventLeaderboardTrackerHide           = C.RC_CLIENT_EVENT_LEADERBOARD_TRACKER_HIDE
	EventLeaderboardTrackerUpdate         = C.RC_CLIENT_EVENT_LEADERBOARD_TRACKER_UPDATE
	EventReset                            = C.RC_CLIENT_EVENT_RESET
	EventGameCompleted                    = C.RC_CLIENT_EVENT_GAME_COMPLETED
	EventServerError                      = C.RC_CLIENT_EVENT_SERVER_ERROR
	EventDisconnected                     = C.RC_CLIENT_EVENT_DISCONNECTED
	EventReconnected                      = C.RC_CLIENT_EVENT_RECONNECTED
)

// Achievement states
const (
	AchievementStateInactive = C.RC_CLIENT_ACHIEVEMENT_STATE_INACTIVE
	AchievementStateActive   = C.RC_CLIENT_ACHIEVEMENT_STATE_ACTIVE
	AchievementStateUnlocked = C.RC_CLIENT_ACHIEVEMENT_STATE_UNLOCKED
	AchievementStateDisabled = C.RC_CLIENT_ACHIEVEMENT_STATE_DISABLED
)

// ReadMemoryFunc is called by rcheevos to read emulator memory
type ReadMemoryFunc func(address uint32, buffer []byte) uint32

// ServerCallFunc is called by rcheevos to make HTTP requests to the RetroAchievements API
type ServerCallFunc func(request *ServerRequest)

// EventFunc is called when achievement events occur
type EventFunc func(event *Event)

// LogFunc is called for logging
type LogFunc func(message string)

// LoginCallbackFunc is called when login completes
type LoginCallbackFunc func(result int, errorMessage string)

// LoadGameCallbackFunc is called when game loading completes
type LoadGameCallbackFunc func(result int, errorMessage string)

// ServerRequest represents an HTTP request to the RetroAchievements API
type ServerRequest struct {
	URL         string
	PostData    string
	ContentType string
	callback    C.rc_client_server_callback_t
	callbackData unsafe.Pointer
	client       *Client
	responded    atomic.Bool
}

// Respond sends the server response back to rcheevos. Only the first call
// has an effect: rcheevos frees the request's state when it handles the
// response, so a second response would use freed memory. A response for a
// client that has been destroyed is dropped, since rcheevos would reach the
// freed client. A body larger than math.MaxInt32 bytes, more than rcheevos
// takes, is reported as a client error with no body instead.
//
// Destroy waits for a Respond in progress, so a callback reached from
// Respond must not destroy the client.
func (r *ServerRequest) Respond(body []byte, httpStatus int) {
	if !r.responded.CompareAndSwap(false, true) {
		return
	}
	c := r.client
	c.respondMu.RLock()
	defer c.respondMu.RUnlock()
	if c.destroyed {
		return
	}
	if len(body) > math.MaxInt32 {
		body = nil
		httpStatus = C.RC_API_SERVER_RESPONSE_CLIENT_ERROR
	}

	var bodyPtr *C.char
	if len(body) > 0 {
		bodyPtr = (*C.char)(unsafe.Pointer(&body[0]))
	}
	C.server_response_callback(r.callback, r.callbackData, bodyPtr, C.int(len(body)), C.int(httpStatus))
}

// Event represents an achievement event
type Event struct {
	Type                 uint32
	Achievement          *Achievement
	Leaderboard          *Leaderboard
	LeaderboardTracker   *LeaderboardTracker
	LeaderboardScoreboard *LeaderboardScoreboard
	ServerError          *ServerError
	Subset               *EventSubset
}

// LeaderboardTracker represents a leaderboard tracker display
type LeaderboardTracker struct {
	ID      uint32
	Display string
}

// LeaderboardScoreboard represents scoreboard data after a leaderboard submission
type LeaderboardScoreboard struct {
	LeaderboardID  uint32
	SubmittedScore string
	BestScore      string
	NewRank        uint32
	NumEntries     uint32
	TopEntries     []*LeaderboardScoreboardEntry
}

// LeaderboardScoreboardEntry represents a single entry in the scoreboard
type LeaderboardScoreboardEntry struct {
	Username string
	Rank     uint32
	Score    string
}

// ServerError represents an API error
type ServerError struct {
	ErrorMessage string
	API          string
	Result       int
	RelatedID    uint32
}

// EventSubset represents subset info in an event
type EventSubset struct {
	ID    uint32
	Title string
}

// Achievement represents achievement information
type Achievement struct {
	ID               uint32
	Title            string
	Description      string
	BadgeName        string
	Points           uint32
	State            uint8
	Category         uint8
	Bucket           uint8
	Unlocked         uint8
	Type             uint8
	MeasuredProgress string
	MeasuredPercent  float32
	Rarity           float32
	RarityHardcore   float32
	UnlockTime       int64
}

// Leaderboard represents leaderboard information
type Leaderboard struct {
	ID            uint32
	Title         string
	Description   string
	TrackerValue  string
	State         uint8
	Format        uint8
	LowerIsBetter bool
}

// User represents user information
type User struct {
	DisplayName string
	Username    string
	Token       string
	Score       uint32
	ScoreSoftcore uint32
}

// Game represents game information
type Game struct {
	ID        uint32
	ConsoleID uint32
	Title     string
	Hash      string
	BadgeName string
}

// Client wraps the rc_client
type Client struct {
	ptr            *C.rc_client_t
	readMemoryFunc ReadMemoryFunc
	serverCallFunc ServerCallFunc
	eventFunc      EventFunc
	logFunc        LogFunc
	// pending holds the requests started and not yet finished, by the
	// handle passed to rcheevos as their userdata
	pending map[cgo.Handle]requestKind
	mu      sync.RWMutex

	// respondMu is held by Respond while it hands a response to rcheevos,
	// and by Destroy to wait for those to finish before freeing the client.
	// destroyed is set under it.
	respondMu sync.RWMutex
	destroyed bool
}

var (
	// clientMap maps C pointers to Go Client instances for callbacks
	clientMap   = make(map[*C.rc_client_t]*Client)
	clientMapMu sync.RWMutex
)

func registerClient(c *Client) {
	clientMapMu.Lock()
	clientMap[c.ptr] = c
	clientMapMu.Unlock()
}

func unregisterClient(c *Client) {
	clientMapMu.Lock()
	delete(clientMap, c.ptr)
	clientMapMu.Unlock()
}

func getClient(ptr *C.rc_client_t) *Client {
	clientMapMu.RLock()
	c := clientMap[ptr]
	clientMapMu.RUnlock()
	return c
}

// NewClient creates a new rcheevos client
func NewClient(readMemory ReadMemoryFunc, serverCall ServerCallFunc) *Client {
	c := &Client{
		readMemoryFunc: readMemory,
		serverCallFunc: serverCall,
		pending:        make(map[cgo.Handle]requestKind),
	}
	c.ptr = C.create_client_with_callbacks()
	registerClient(c)
	C.set_event_handler(c.ptr)
	return c
}

// Destroy releases resources associated with the client. It waits for any
// server response being handled, then drops responses that arrive later
// and the callbacks of requests still pending, which are never called.
func (c *Client) Destroy() {
	c.respondMu.Lock()
	if c.destroyed {
		c.respondMu.Unlock()
		return
	}
	c.destroyed = true
	c.respondMu.Unlock()

	// Callbacks rcheevos makes from here on find no client and are
	// ignored, and no new request can start
	unregisterClient(c)
	c.mu.Lock()
	ptr := c.ptr
	c.ptr = nil
	c.mu.Unlock()
	c.dropRequests()

	C.rc_client_destroy(ptr)
}

// SetEventHandler sets the event callback
func (c *Client) SetEventHandler(handler EventFunc) {
	c.mu.Lock()
	c.eventFunc = handler
	c.mu.Unlock()
}

// EnableLogging enables logging at the specified level
func (c *Client) EnableLogging(level int, logFunc LogFunc) {
	c.mu.Lock()
	c.logFunc = logFunc
	c.mu.Unlock()
	C.enable_logging(c.ptr, C.int(level))
}

// SetHardcoreEnabled enables or disables hardcore mode
func (c *Client) SetHardcoreEnabled(enabled bool) {
	val := C.int(0)
	if enabled {
		val = 1
	}
	C.rc_client_set_hardcore_enabled(c.ptr, val)
}

// GetHardcoreEnabled returns whether hardcore mode is enabled
func (c *Client) GetHardcoreEnabled() bool {
	return C.rc_client_get_hardcore_enabled(c.ptr) != 0
}

// LoginWithPassword starts a login with username and password. callback
// receives the result of this login only; it may be called before
// LoginWithPassword returns, for example when a login is already in
// progress.
func (c *Client) LoginWithPassword(username, password string, callback LoginCallbackFunc) {
	h, ok := c.beginRequest(requestOther, callback)
	if !ok {
		if callback != nil {
			callback(InvalidState, msgClientRequired)
		}
		return
	}

	cUsername := C.CString(username)
	cPassword := C.CString(password)
	defer C.free(unsafe.Pointer(cUsername))
	defer C.free(unsafe.Pointer(cPassword))

	C.begin_login_with_password(c.ptr, cUsername, cPassword, C.uintptr_t(h))
}

// LoginWithToken starts a login with username and token. callback receives
// the result of this login only; it may be called before LoginWithToken
// returns.
func (c *Client) LoginWithToken(username, token string, callback LoginCallbackFunc) {
	h, ok := c.beginRequest(requestOther, callback)
	if !ok {
		if callback != nil {
			callback(InvalidState, msgClientRequired)
		}
		return
	}

	cUsername := C.CString(username)
	cToken := C.CString(token)
	defer C.free(unsafe.Pointer(cUsername))
	defer C.free(unsafe.Pointer(cToken))

	C.begin_login_with_token(c.ptr, cUsername, cToken, C.uintptr_t(h))
}

// Logout logs out the current user. A login in progress is reported
// aborted to its callback, and so is a game load in progress.
func (c *Client) Logout() {
	c.retireLoads()
	C.rc_client_logout(c.ptr)
}

// GetUser returns the current user info, or nil if not logged in
func (c *Client) GetUser() *User {
	user := C.rc_client_get_user_info(c.ptr)
	if user == nil {
		return nil
	}
	return &User{
		DisplayName:   C.GoString(user.display_name),
		Username:      C.GoString(user.username),
		Token:         C.GoString(user.token),
		Score:         uint32(user.score),
		ScoreSoftcore: uint32(user.score_softcore),
	}
}

// LoadGame starts loading a game by hash. callback receives the result of
// this load only; it may be called before LoadGame returns. A load still in
// progress is abandoned and reported aborted to its own callback.
func (c *Client) LoadGame(hash string, callback LoadGameCallbackFunc) {
	// rcheevos rejects an empty hash before touching the load in progress
	if hash != "" {
		c.retireLoads()
	}
	h, ok := c.beginRequest(requestLoad, callback)
	if !ok {
		if callback != nil {
			callback(InvalidState, msgClientRequired)
		}
		return
	}

	cHash := C.CString(hash)
	defer C.free(unsafe.Pointer(cHash))

	C.begin_load_game(c.ptr, cHash, C.uintptr_t(h))
}

// IdentifyAndLoadGame identifies a game by ROM data and loads it. callback
// receives the result of this load only; it may be called before
// IdentifyAndLoadGame returns. A load still in progress is abandoned and
// reported aborted to its own callback.
func (c *Client) IdentifyAndLoadGame(consoleID uint32, filePath string, data []byte, callback LoadGameCallbackFunc) {
	c.retireLoads()
	h, ok := c.beginRequest(requestLoad, callback)
	if !ok {
		if callback != nil {
			callback(InvalidState, msgClientRequired)
		}
		return
	}

	var cFilePath *C.char
	if filePath != "" {
		cFilePath = C.CString(filePath)
		defer C.free(unsafe.Pointer(cFilePath))
	}

	var dataPtr *C.uint8_t
	if len(data) > 0 {
		dataPtr = (*C.uint8_t)(unsafe.Pointer(&data[0]))
	}

	C.begin_identify_and_load_game(c.ptr, C.uint32_t(consoleID), cFilePath, dataPtr, C.size_t(len(data)), C.uintptr_t(h))
}

// UnloadGame unloads the current game. A load in progress is abandoned and
// reported aborted to its callback.
func (c *Client) UnloadGame() {
	c.retireLoads()
	C.rc_client_unload_game(c.ptr)
}

// IsGameLoaded returns whether a game is currently loaded
func (c *Client) IsGameLoaded() bool {
	return C.rc_client_is_game_loaded(c.ptr) != 0
}

// GetGame returns the current game info, or nil if no game is loaded
func (c *Client) GetGame() *Game {
	game := C.rc_client_get_game_info(c.ptr)
	if game == nil {
		return nil
	}
	return &Game{
		ID:        uint32(game.id),
		ConsoleID: uint32(game.console_id),
		Title:     C.GoString(game.title),
		Hash:      C.GoString(game.hash),
		BadgeName: C.GoString(game.badge_name),
	}
}

// DoFrame processes achievements for the current frame
func (c *Client) DoFrame() {
	C.rc_client_do_frame(c.ptr)
}

// Idle processes the periodic queue (call when paused)
func (c *Client) Idle() {
	C.rc_client_idle(c.ptr)
}

// Reset resets the runtime state
func (c *Client) Reset() {
	C.rc_client_reset(c.ptr)
}

// GetRichPresenceMessage returns the current rich presence message, cut
// to 255 bytes when longer
func (c *Client) GetRichPresenceMessage() string {
	buffer := make([]byte, 256)
	length := C.rc_client_get_rich_presence_message(c.ptr, (*C.char)(unsafe.Pointer(&buffer[0])), C.size_t(len(buffer)))
	return bufferString(buffer, uint64(length))
}

// bufferString returns the text rcheevos wrote into buffer, given the
// length it returned. That length can be larger than the buffer: for a
// message that did not fit, rcheevos returns the length it needed, not the
// length it wrote. The text ends at the first NUL byte and never past the
// buffer.
func bufferString(buffer []byte, length uint64) string {
	if length < uint64(len(buffer)) {
		buffer = buffer[:length]
	}
	if i := bytes.IndexByte(buffer, 0); i >= 0 {
		buffer = buffer[:i]
	}
	return string(buffer)
}

// HasRichPresence returns whether the game supports rich presence
func (c *Client) HasRichPresence() bool {
	return C.rc_client_has_rich_presence(c.ptr) != 0
}

// HasAchievements returns whether the game has achievements
func (c *Client) HasAchievements() bool {
	return C.rc_client_has_achievements(c.ptr) != 0
}

// HasLeaderboards returns whether the game has leaderboards
func (c *Client) HasLeaderboards() bool {
	return C.rc_client_has_leaderboards(c.ptr) != 0
}

// GetAchievement returns achievement info by ID
func (c *Client) GetAchievement(id uint32) *Achievement {
	ach := C.rc_client_get_achievement_info(c.ptr, C.uint32_t(id))
	if ach == nil {
		return nil
	}
	return achievementFromC(ach)
}

func achievementFromC(ach *C.rc_client_achievement_t) *Achievement {
	return &Achievement{
		ID:               uint32(ach.id),
		Title:            C.GoString(ach.title),
		Description:      C.GoString(ach.description),
		BadgeName:        C.GoString(&ach.badge_name[0]),
		Points:           uint32(ach.points),
		State:            uint8(ach.state),
		Category:         uint8(ach.category),
		Bucket:           uint8(ach.bucket),
		Unlocked:         uint8(ach.unlocked),
		Type:             uint8(ach._type),
		MeasuredProgress: C.GoString(&ach.measured_progress[0]),
		MeasuredPercent:  float32(ach.measured_percent),
		Rarity:           float32(ach.rarity),
		RarityHardcore:   float32(ach.rarity_hardcore),
		UnlockTime:       int64(ach.unlock_time),
	}
}

// GetLeaderboard returns leaderboard info by ID
func (c *Client) GetLeaderboard(id uint32) *Leaderboard {
	lb := C.rc_client_get_leaderboard_info(c.ptr, C.uint32_t(id))
	if lb == nil {
		return nil
	}
	return leaderboardFromC(lb)
}

func leaderboardFromC(lb *C.rc_client_leaderboard_t) *Leaderboard {
	var trackerValue string
	if lb.tracker_value != nil {
		trackerValue = C.GoString(lb.tracker_value)
	}
	return &Leaderboard{
		ID:            uint32(lb.id),
		Title:         C.GoString(lb.title),
		Description:   C.GoString(lb.description),
		TrackerValue:  trackerValue,
		State:         uint8(lb.state),
		Format:        uint8(lb.format),
		LowerIsBetter: lb.lower_is_better != 0,
	}
}

// ConsoleName returns the name for a console ID
func ConsoleName(consoleID uint32) string {
	return C.GoString(C.rc_console_name(C.uint32_t(consoleID)))
}

// cBeginFetchLeaderboardEntries is an internal wrapper for C.begin_fetch_leaderboard_entries
func cBeginFetchLeaderboardEntries(client *C.rc_client_t, leaderboardID, firstEntry, count uint32, h cgo.Handle) {
	C.begin_fetch_leaderboard_entries(client, C.uint32_t(leaderboardID), C.uint32_t(firstEntry), C.uint32_t(count), C.uintptr_t(h))
}

// cBeginFetchLeaderboardEntriesAroundUser is an internal wrapper for C.begin_fetch_leaderboard_entries_around_user
func cBeginFetchLeaderboardEntriesAroundUser(client *C.rc_client_t, leaderboardID, count uint32, h cgo.Handle) {
	C.begin_fetch_leaderboard_entries_around_user(client, C.uint32_t(leaderboardID), C.uint32_t(count), C.uintptr_t(h))
}

// cSetGetTimeMillisecsFunction is an internal wrapper
func cSetGetTimeMillisecsFunction(client *C.rc_client_t) {
	C.set_get_time_millisecs_function(client)
}

// cBeginChangeMedia is an internal wrapper
func cBeginChangeMedia(client *C.rc_client_t, hash *C.char, h cgo.Handle) {
	C.begin_change_media(client, hash, C.uintptr_t(h))
}

// cBeginIdentifyAndChangeMedia is an internal wrapper
func cBeginIdentifyAndChangeMedia(client *C.rc_client_t, filePath *C.char, data *C.uint8_t, dataSize C.size_t, h cgo.Handle) {
	C.begin_identify_and_change_media(client, filePath, data, dataSize, C.uintptr_t(h))
}

// cBeginFetchHashLibrary is an internal wrapper
func cBeginFetchHashLibrary(client *C.rc_client_t, consoleID uint32, h cgo.Handle) {
	C.begin_fetch_hash_library(client, C.uint32_t(consoleID), C.uintptr_t(h))
}

// cBeginFetchGameTitles is an internal wrapper
func cBeginFetchGameTitles(client *C.rc_client_t, gameIDs *C.uint32_t, numGameIDs uint32, h cgo.Handle) {
	C.begin_fetch_game_titles(client, gameIDs, C.uint32_t(numGameIDs), C.uintptr_t(h))
}

// cBeginFetchAllUserProgress is an internal wrapper
func cBeginFetchAllUserProgress(client *C.rc_client_t, consoleID uint32, h cgo.Handle) {
	C.begin_fetch_all_user_progress(client, C.uint32_t(consoleID), C.uintptr_t(h))
}
