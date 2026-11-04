package rcheevos

import (
	"sync"
	"testing"
)

// requestRecorder is a server call function that keeps the requests it is
// given instead of answering them
type requestRecorder struct {
	mu   sync.Mutex
	reqs []*ServerRequest
}

func (r *requestRecorder) call(req *ServerRequest) {
	r.mu.Lock()
	r.reqs = append(r.reqs, req)
	r.mu.Unlock()
}

func (r *requestRecorder) get(t *testing.T, i int) *ServerRequest {
	t.Helper()
	r.mu.Lock()
	defer r.mu.Unlock()
	if i >= len(r.reqs) {
		t.Fatalf("server request %d not made; %d made", i, len(r.reqs))
	}
	return r.reqs[i]
}

// result records the calls one request's callback received
type result struct {
	calls []int
	msgs  []string
}

func (r *result) callback(code int, msg string) {
	r.calls = append(r.calls, code)
	r.msgs = append(r.msgs, msg)
}

func newRecordedClient(t *testing.T) (*Client, *requestRecorder) {
	t.Helper()
	rec := &requestRecorder{}
	c := NewClient(func(uint32, []byte) uint32 { return 0 }, rec.call)
	t.Cleanup(c.Destroy)
	return c, rec
}

// logIn logs c in through rec, whose next request is the login
func logIn(t *testing.T, c *Client, rec *requestRecorder, next int) {
	t.Helper()
	var login result
	c.LoginWithPassword("user", "pass", login.callback)
	rec.get(t, next).Respond([]byte(`{"Success":true,"User":"user","Token":"token"}`), 200)
	if len(login.calls) != 1 || login.calls[0] != OK {
		t.Fatalf("login calls = %v %v, want one OK", login.calls, login.msgs)
	}
}

func pendingCount(c *Client) int {
	c.mu.RLock()
	defer c.mu.RUnlock()
	return len(c.pending)
}

// TestOverlappingLogins tests that a second login while one is in progress
// is rejected to its own callback, and the first login's result still
// reaches the first callback only.
func TestOverlappingLogins(t *testing.T) {
	c, rec := newRecordedClient(t)
	var first, second result

	c.LoginWithPassword("user", "pass", first.callback)
	c.LoginWithPassword("user", "pass", second.callback)
	if len(second.calls) != 1 || second.calls[0] != InvalidState {
		t.Fatalf("second login calls = %v %v, want one InvalidState", second.calls, second.msgs)
	}

	rec.get(t, 0).Respond([]byte(`{"Success":true,"User":"user","Token":"token"}`), 200)
	if len(first.calls) != 1 || first.calls[0] != OK {
		t.Fatalf("first login calls = %v %v, want one OK", first.calls, first.msgs)
	}
	if len(second.calls) != 1 {
		t.Fatalf("second login got the first login's result: %v", second.calls)
	}
	if n := pendingCount(c); n != 0 {
		t.Fatalf("%d requests still pending", n)
	}
}

// TestSupersededLoad tests that a load replaced by a newer one is reported
// aborted to its own callback, once, even after its server response
// arrives, and the newer load's callback is left waiting for its own
// result.
func TestSupersededLoad(t *testing.T) {
	c, rec := newRecordedClient(t)
	logIn(t, c, rec, 0)
	var older, newer result

	c.LoadGame("0123456789abcdef0123456789abcdef", older.callback)
	olderReq := rec.get(t, 1)
	c.LoadGame("fedcba9876543210fedcba9876543210", newer.callback)

	if len(older.calls) != 1 || older.calls[0] != Aborted {
		t.Fatalf("older load calls = %v %v, want one Aborted", older.calls, older.msgs)
	}
	olderReq.Respond([]byte(`{"Success":true,"GameID":0}`), 200)
	if len(older.calls) != 1 {
		t.Fatalf("older load called again after its response: %v %v", older.calls, older.msgs)
	}
	if len(newer.calls) != 0 {
		t.Fatalf("newer load calls = %v %v, want none yet", newer.calls, newer.msgs)
	}
	if n := pendingCount(c); n != 1 {
		t.Fatalf("%d requests pending, want 1", n)
	}
}

// TestLogoutAbortsLoadOnce tests that logging out reports a load waiting
// for the login aborted exactly once, although rcheevos also reports it.
func TestLogoutAbortsLoadOnce(t *testing.T) {
	c, _ := newRecordedClient(t)
	var login, load result

	c.LoginWithPassword("user", "pass", login.callback)
	c.LoadGame("0123456789abcdef0123456789abcdef", load.callback)
	c.Logout()

	if len(load.calls) != 1 || load.calls[0] != Aborted {
		t.Fatalf("load calls = %v %v, want one Aborted", load.calls, load.msgs)
	}
	if n := pendingCount(c); n != 1 {
		t.Fatalf("%d requests pending, want 1 (the login)", n)
	}
}

// TestUnloadAbortsLoad tests that unloading reports a pending load aborted
func TestUnloadAbortsLoad(t *testing.T) {
	c, rec := newRecordedClient(t)
	logIn(t, c, rec, 0)
	var load result

	c.LoadGame("0123456789abcdef0123456789abcdef", load.callback)
	c.UnloadGame()

	if len(load.calls) != 1 || load.calls[0] != Aborted {
		t.Fatalf("load calls = %v %v, want one Aborted", load.calls, load.msgs)
	}
	if n := pendingCount(c); n != 0 {
		t.Fatalf("%d requests still pending", n)
	}
}

// TestEmptyHashKeepsLoad tests that a load with no hash, which rcheevos
// rejects before touching the load in progress, does not abort it.
func TestEmptyHashKeepsLoad(t *testing.T) {
	c, rec := newRecordedClient(t)
	logIn(t, c, rec, 0)
	var load, empty result

	c.LoadGame("0123456789abcdef0123456789abcdef", load.callback)
	c.LoadGame("", empty.callback)

	if len(empty.calls) != 1 || empty.calls[0] != InvalidState {
		t.Fatalf("empty hash calls = %v %v, want one InvalidState", empty.calls, empty.msgs)
	}
	if len(load.calls) != 0 {
		t.Fatalf("load in progress calls = %v, want none", load.calls)
	}
	if n := pendingCount(c); n != 1 {
		t.Fatalf("%d requests pending, want 1", n)
	}
}

// TestOverlappingFetches tests that two fetches of the same kind each get
// their own result, answered in the opposite order.
func TestOverlappingFetches(t *testing.T) {
	c, rec := newRecordedClient(t)
	got := map[string]string{}
	fetch := func(name string) FetchHashLibraryCallback {
		return func(code int, _ string, lib *HashLibrary) {
			if code != OK || lib == nil || len(lib.Entries) != 1 {
				got[name] = "bad result"
				return
			}
			got[name] = lib.Entries[0].Hash
			lib.Destroy()
		}
	}

	c.FetchHashLibrary(ConsoleSaturn, fetch("first"))
	c.FetchHashLibrary(ConsoleMegaDrive, fetch("second"))

	rec.get(t, 1).Respond([]byte(`{"Success":true,"MD5List":{"bbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbb":2}}`), 200)
	rec.get(t, 0).Respond([]byte(`{"Success":true,"MD5List":{"aaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaa":1}}`), 200)

	if got["first"] != "aaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaa" || got["second"] != "bbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbb" {
		t.Fatalf("results = %v", got)
	}
	if n := pendingCount(c); n != 0 {
		t.Fatalf("%d requests still pending", n)
	}
}

// TestDestroyDuringResponses tests that responses handled on several
// goroutines while the client is destroyed are either handled whole before
// Destroy or dropped after it. Run with -race.
func TestDestroyDuringResponses(t *testing.T) {
	const fetches = 16
	c, rec := newRecordedClient(t)
	var mu sync.Mutex
	done := 0
	for range fetches {
		c.FetchHashLibrary(ConsoleSaturn, func(_ int, _ string, lib *HashLibrary) {
			if lib != nil {
				lib.Destroy()
			}
			mu.Lock()
			done++
			mu.Unlock()
		})
	}

	var wg sync.WaitGroup
	for i := range fetches {
		req := rec.get(t, i)
		wg.Add(1)
		go func() {
			defer wg.Done()
			req.Respond([]byte(`{"Success":true,"MD5List":{}}`), 200)
		}()
	}
	c.Destroy()
	wg.Wait()

	mu.Lock()
	defer mu.Unlock()
	if done > fetches {
		t.Fatalf("%d callbacks for %d requests", done, fetches)
	}
	if n := pendingCount(c); n != 0 {
		t.Fatalf("%d requests still pending", n)
	}
}

// TestDestroyWithPendingRequests tests that destroying a client drops its
// pending requests without calling them, that a response arriving later is
// dropped, and that requests made afterwards fail at once.
func TestDestroyWithPendingRequests(t *testing.T) {
	c, rec := newRecordedClient(t)
	called := false
	c.FetchHashLibrary(ConsoleSaturn, func(int, string, *HashLibrary) { called = true })
	req := rec.get(t, 0)

	c.Destroy()
	if n := pendingCount(c); n != 0 {
		t.Fatalf("%d requests still pending after Destroy", n)
	}

	req.Respond([]byte(`{"Success":true,"MD5List":{}}`), 200)
	if called {
		t.Fatal("callback called for a response after Destroy")
	}

	var login result
	c.LoginWithPassword("user", "pass", login.callback)
	if len(login.calls) != 1 || login.calls[0] != InvalidState {
		t.Fatalf("login after Destroy calls = %v, want one InvalidState", login.calls)
	}
	c.Destroy()
}
