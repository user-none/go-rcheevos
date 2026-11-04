package rcheevos

/*
#include <stdlib.h>
#include <string.h>
#include <stdint.h>
#include "rc_hash.h"
#include "rc_client.h"

// Go callback bridges implemented (//export) in disc_callbacks.go.
extern uintptr_t goCDOpenTrack(uintptr_t disc_handle, uint32_t track);
extern size_t goCDReadSector(uintptr_t handle, uint32_t sector, void* buffer, size_t requested_bytes);
extern void goCDCloseTrack(uintptr_t handle);
extern uint32_t goCDFirstTrackSector(uintptr_t handle);

// goLoadGameCallback is the shared load-game bridge (//export in callbacks.go).
extern void goLoadGameCallback(int result, const char* error_message, rc_client_t* client, void* userdata);

// go_track tags an opened track handle as either backed by a Go CDReader (via a
// runtime/cgo handle) or delegated to the default cdreader, so the shared
// read/close/first-sector callbacks route to the right place. Discs whose path
// is a "godisc:<handle>" token are read through Go; any other path falls back
// to the built-in cdreader so existing file-based hashing is unchanged.
typedef struct {
    int is_go;
    uintptr_t go_handle; // cgo.Handle to a cdTrack, when is_go
    void* def_handle;    // default cdreader track handle, otherwise
} go_track;

static struct rc_hash_cdreader g_default_cdreader;
static struct rc_hash_cdreader g_go_cdreader;

static void* go_open_track_iterator(const char* path, uint32_t track, const struct rc_hash_iterator* iterator) {
    go_track* t = (go_track*)malloc(sizeof(go_track));
    if (!t)
        return NULL;

    if (path && strncmp(path, "godisc:", 7) == 0) {
        uintptr_t disc_handle = (uintptr_t)strtoull(path + 7, NULL, 10);
        uintptr_t h = goCDOpenTrack(disc_handle, track);
        if (h == 0) {
            free(t);
            return NULL;
        }
        t->is_go = 1;
        t->go_handle = h;
        t->def_handle = NULL;
        return t;
    }

    t->is_go = 0;
    t->go_handle = 0;
    if (g_default_cdreader.open_track_iterator)
        t->def_handle = g_default_cdreader.open_track_iterator(path, track, iterator);
    else if (g_default_cdreader.open_track)
        t->def_handle = g_default_cdreader.open_track(path, track);
    else
        t->def_handle = NULL;
    if (!t->def_handle) {
        free(t);
        return NULL;
    }
    return t;
}

static size_t go_read_sector(void* track_handle, uint32_t sector, void* buffer, size_t requested_bytes) {
    go_track* t = (go_track*)track_handle;
    if (!t)
        return 0;
    if (t->is_go)
        return goCDReadSector(t->go_handle, sector, buffer, requested_bytes);
    if (g_default_cdreader.read_sector)
        return g_default_cdreader.read_sector(t->def_handle, sector, buffer, requested_bytes);
    return 0;
}

static void go_close_track(void* track_handle) {
    go_track* t = (go_track*)track_handle;
    if (!t)
        return;
    if (t->is_go)
        goCDCloseTrack(t->go_handle);
    else if (g_default_cdreader.close_track)
        g_default_cdreader.close_track(t->def_handle);
    free(t);
}

static uint32_t go_first_track_sector(void* track_handle) {
    go_track* t = (go_track*)track_handle;
    if (!t)
        return 0;
    if (t->is_go)
        return goCDFirstTrackSector(t->go_handle);
    if (g_default_cdreader.first_track_sector)
        return g_default_cdreader.first_track_sector(t->def_handle);
    return 0;
}

// init_go_cdreader captures the built-in cdreader (used for non-Go paths) and
// installs the Go-routing cdreader globally, so both the standalone iterator
// path and the rc_client identify path read through it.
static void init_go_cdreader(void) {
    rc_hash_get_default_cdreader(&g_default_cdreader);
    g_go_cdreader.open_track = NULL;
    g_go_cdreader.read_sector = go_read_sector;
    g_go_cdreader.close_track = go_close_track;
    g_go_cdreader.first_track_sector = go_first_track_sector;
    g_go_cdreader.open_track_iterator = go_open_track_iterator;
    rc_hash_init_custom_cdreader(&g_go_cdreader);
}

static void c_load_game_callback_disc(int result, const char* error_message, rc_client_t* client, void* userdata) {
    goLoadGameCallback(result, error_message, client, userdata);
}

static rc_client_async_handle_t* begin_identify_and_load_game_disc(rc_client_t* client, uint32_t console_id, const char* path, uintptr_t userdata) {
    return rc_client_begin_identify_and_load_game(client, console_id, path, NULL, 0, c_load_game_callback_disc, (void*)userdata);
}
*/
import "C"
import (
	"runtime/cgo"
	"strconv"
	"strings"
	"unsafe"
)

// CDReader supplies disc sectors and track layout for disc-based hashing. A
// caller passes any value implementing it (its signatures match a streaming
// disc reader exactly, so a concrete reader satisfies it structurally without
// importing this package). rcheevos reads sectors through it instead of opening
// a file itself, so any image format the caller supports works. This package
// never closes the reader; the caller owns its lifetime.
type CDReader interface {
	// ReadSectorData returns the cooked user data of the sector at the given
	// absolute LBA (header stripped for data sectors; audio sectors whole).
	ReadSectorData(lba int) ([]byte, error)
	// NumTracks returns the number of tracks on the disc.
	NumTracks() int
	// Track returns the TOC fields for track index i in [0, NumTracks).
	Track(i int) (number int, typ string, frames int, pregap int, startLBA int, control uint8)
}

func init() {
	C.init_go_cdreader()
}

// Track sentinels from rc_hash.h (RC_HASH_CDTRACK_*), as the uint32 the
// cdreader receives.
const (
	cdTrackFirstData = 0xFFFFFFFF // (uint32)-1
	cdTrackLast      = 0xFFFFFFFE // (uint32)-2
	cdTrackLargest   = 0xFFFFFFFD // (uint32)-3
)

// cdTrack is the per-open-track state referenced by a cgo.Handle. read fills a
// buffer from successive cooked sectors, mirroring the C cdreader contract
// (requested bytes may span more than one sector).
type cdTrack struct {
	disc        CDReader
	firstSector int
}

func (t *cdTrack) read(sector int, buf []byte) int {
	n := 0
	for n < len(buf) {
		data, err := t.disc.ReadSectorData(sector)
		if err != nil {
			break
		}
		copied := copy(buf[n:], data)
		n += copied
		if copied < len(data) {
			break // buffer filled mid-sector
		}
		sector++
	}
	return n
}

// trackStartSector returns the absolute sector where the requested track's data
// begins (start LBA past any pregap). track is either a 1-based track number or
// an RC_HASH_CDTRACK_ sentinel.
func trackStartSector(disc CDReader, track uint32) int {
	n := disc.NumTracks()
	switch track {
	case cdTrackFirstData:
		for i := 0; i < n; i++ {
			_, typ, _, pregap, startLBA, _ := disc.Track(i)
			if !strings.HasPrefix(typ, "AUDIO") {
				return startLBA + pregap
			}
		}
	case cdTrackLargest:
		best, start := -1, 0
		for i := 0; i < n; i++ {
			_, _, frames, pregap, startLBA, _ := disc.Track(i)
			if frames > best {
				best = frames
				start = startLBA + pregap
			}
		}
		if best >= 0 {
			return start
		}
	case cdTrackLast:
		if n > 0 {
			_, _, _, pregap, startLBA, _ := disc.Track(n - 1)
			return startLBA + pregap
		}
	default: // literal 1-based track number
		for i := 0; i < n; i++ {
			num, _, _, pregap, startLBA, _ := disc.Track(i)
			if num == int(track) {
				return startLBA + pregap
			}
		}
	}
	return 0
}

// discToken encodes a disc's cgo.Handle as the synthetic path rcheevos passes
// to the cdreader's open_track; the C side recognizes the "godisc:" prefix and
// routes back into Go.
func discToken(h cgo.Handle) string {
	return "godisc:" + strconv.FormatUint(uint64(uintptr(h)), 10)
}

// HashFromDisc computes the hash for a disc-based console by reading sectors
// through disc. Returns "" on failure.
func HashFromDisc(consoleID uint32, disc CDReader) string {
	h := cgo.NewHandle(disc)
	defer h.Delete()

	it := NewHashIterator(discToken(h))
	defer it.Destroy()
	return it.GenerateHash(consoleID)
}

// IdentifyAndLoadGameFromDisc identifies a disc-based game by reading sectors
// through disc and loads it for achievement tracking. consoleID selects the
// hashing algorithm (e.g. 39 for Saturn). The disc is read synchronously while
// the hash is computed during this call, so it need not outlive the call; the
// load callback fires later from the async server round-trip. callback
// receives the result of this load only; it may be called before
// IdentifyAndLoadGameFromDisc returns. A load still in progress is abandoned
// and reported aborted to its own callback.
func (c *Client) IdentifyAndLoadGameFromDisc(consoleID uint32, disc CDReader, callback LoadGameCallbackFunc) {
	c.retireLoads()
	req, ok := c.beginRequest(requestLoad, callback)
	if !ok {
		if callback != nil {
			callback(InvalidState, msgClientRequired)
		}
		return
	}

	h := cgo.NewHandle(disc)
	defer h.Delete()

	token := C.CString(discToken(h))
	defer C.free(unsafe.Pointer(token))

	C.begin_identify_and_load_game_disc(c.ptr, C.uint32_t(consoleID), token, C.uintptr_t(req))
}
