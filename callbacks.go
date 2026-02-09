package rcheevos

/*
#include "rc_client.h"
#include "rc_api_request.h"
*/
import "C"
import "unsafe"

//export goReadMemoryCallback
func goReadMemoryCallback(address C.uint32_t, buffer *C.uint8_t, numBytes C.uint32_t, client *C.rc_client_t) C.uint32_t {
	c := getClient(client)
	if c == nil || c.readMemoryFunc == nil {
		return 0
	}

	goBuffer := unsafe.Slice((*byte)(unsafe.Pointer(buffer)), numBytes)
	return C.uint32_t(c.readMemoryFunc(uint32(address), goBuffer))
}

//export goServerCallCallback
func goServerCallCallback(request *C.rc_api_request_t, callback C.rc_client_server_callback_t, callbackData unsafe.Pointer, client *C.rc_client_t) {
	c := getClient(client)
	if c == nil || c.serverCallFunc == nil {
		return
	}

	req := &ServerRequest{
		URL:          C.GoString(request.url),
		PostData:     C.GoString(request.post_data),
		ContentType:  C.GoString(request.content_type),
		callback:     callback,
		callbackData: callbackData,
	}

	c.serverCallFunc(req)
}

//export goEventCallback
func goEventCallback(event *C.rc_client_event_t, client *C.rc_client_t) {
	c := getClient(client)
	if c == nil {
		return
	}

	c.mu.RLock()
	handler := c.eventFunc
	c.mu.RUnlock()

	if handler == nil {
		return
	}

	ev := &Event{
		Type: uint32(event._type),
	}

	if event.achievement != nil {
		ev.Achievement = achievementFromC(event.achievement)
	}

	if event.leaderboard != nil {
		ev.Leaderboard = leaderboardFromC(event.leaderboard)
	}

	if event.leaderboard_tracker != nil {
		ev.LeaderboardTracker = &LeaderboardTracker{
			ID:      uint32(event.leaderboard_tracker.id),
			Display: C.GoString(&event.leaderboard_tracker.display[0]),
		}
	}

	if event.leaderboard_scoreboard != nil {
		sb := event.leaderboard_scoreboard
		scoreboard := &LeaderboardScoreboard{
			LeaderboardID:  uint32(sb.leaderboard_id),
			SubmittedScore: C.GoString(&sb.submitted_score[0]),
			BestScore:      C.GoString(&sb.best_score[0]),
			NewRank:        uint32(sb.new_rank),
			NumEntries:     uint32(sb.num_entries),
		}

		if sb.num_top_entries > 0 && sb.top_entries != nil {
			scoreboard.TopEntries = make([]*LeaderboardScoreboardEntry, sb.num_top_entries)
			topEntries := unsafe.Slice(sb.top_entries, sb.num_top_entries)
			for i := range topEntries {
				scoreboard.TopEntries[i] = &LeaderboardScoreboardEntry{
					Username: C.GoString(topEntries[i].username),
					Rank:     uint32(topEntries[i].rank),
					Score:    C.GoString(&topEntries[i].score[0]),
				}
			}
		}

		ev.LeaderboardScoreboard = scoreboard
	}

	if event.server_error != nil {
		ev.ServerError = &ServerError{
			ErrorMessage: C.GoString(event.server_error.error_message),
			API:          C.GoString(event.server_error.api),
			Result:       int(event.server_error.result),
			RelatedID:    uint32(event.server_error.related_id),
		}
	}

	if event.subset != nil {
		ev.Subset = &EventSubset{
			ID:    uint32(event.subset.id),
			Title: C.GoString(event.subset.title),
		}
	}

	handler(ev)
}

//export goLogCallback
func goLogCallback(message *C.char, client *C.rc_client_t) {
	c := getClient(client)
	if c == nil {
		return
	}

	c.mu.RLock()
	logFunc := c.logFunc
	c.mu.RUnlock()

	if logFunc != nil {
		logFunc(C.GoString(message))
	}
}

//export goLoginCallback
func goLoginCallback(result C.int, errorMessage *C.char, client *C.rc_client_t, userdata unsafe.Pointer) {
	c := getClient(client)
	if c == nil {
		return
	}

	c.mu.RLock()
	callback := c.loginCallback
	c.mu.RUnlock()

	if callback != nil {
		var errMsg string
		if errorMessage != nil {
			errMsg = C.GoString(errorMessage)
		}
		callback(int(result), errMsg)
	}
}

//export goLoadGameCallback
func goLoadGameCallback(result C.int, errorMessage *C.char, client *C.rc_client_t, userdata unsafe.Pointer) {
	c := getClient(client)
	if c == nil {
		return
	}

	c.mu.RLock()
	callback := c.loadGameCallback
	c.mu.RUnlock()

	if callback != nil {
		var errMsg string
		if errorMessage != nil {
			errMsg = C.GoString(errorMessage)
		}
		callback(int(result), errMsg)
	}
}
