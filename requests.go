package rcheevos

/*
#include "rc_client.h"
*/
import "C"
import (
	"runtime/cgo"
	"unsafe"
)

// requestKind tells pending requests apart where rcheevos treats them
// differently
type requestKind int

const (
	// requestOther is a request rcheevos always completes with its
	// callback, unless the client is destroyed
	requestOther requestKind = iota
	// requestLoad is a game load. rcheevos abandons a load in progress
	// without calling its callback when another load starts, the game is
	// unloaded, or the user logs out.
	requestLoad
)

// Messages for requests that end without reaching rcheevos
const (
	msgClientRequired = "client is required"
	msgLoadAborted    = "The requested game is no longer active"
)

// beginRequest registers callback as a pending request of c and returns
// the handle passed to rcheevos as the request's userdata. ok is false when
// the client has been destroyed; the request must then not be started.
func (c *Client) beginRequest(kind requestKind, callback any) (h cgo.Handle, ok bool) {
	c.mu.Lock()
	defer c.mu.Unlock()
	if c.ptr == nil {
		return 0, false
	}
	h = cgo.NewHandle(callback)
	c.pending[h] = kind
	return h, true
}

// finishRequest ends the pending request whose handle rcheevos passed back
// as userdata and returns its callback. ok is false when the request is no
// longer pending: it was already finished, retired as aborted, or its
// client destroyed. The callback must then not be called.
func finishRequest(client *C.rc_client_t, userdata unsafe.Pointer) (callback any, ok bool) {
	c := getClient(client)
	if c == nil {
		return nil, false
	}
	h := cgo.Handle(uintptr(userdata))

	c.mu.Lock()
	_, ok = c.pending[h]
	delete(c.pending, h)
	c.mu.Unlock()

	if !ok {
		return nil, false
	}
	callback = h.Value()
	h.Delete()
	return callback, true
}

// retireLoads ends every pending game load, reporting it aborted. It is
// called before a call that makes rcheevos abandon a load in progress,
// which rcheevos does without calling the load's callback. Should
// rcheevos still call one back, finishRequest no longer finds it.
func (c *Client) retireLoads() {
	c.mu.Lock()
	var loads []cgo.Handle
	for h, kind := range c.pending {
		if kind == requestLoad {
			loads = append(loads, h)
			delete(c.pending, h)
		}
	}
	c.mu.Unlock()

	for _, h := range loads {
		callback := h.Value().(LoadGameCallbackFunc)
		h.Delete()
		if callback != nil {
			callback(Aborted, msgLoadAborted)
		}
	}
}

// dropRequests ends every pending request without calling its callback.
// The client is being destroyed, and rcheevos frees the requests without
// calling back.
func (c *Client) dropRequests() {
	c.mu.Lock()
	for h := range c.pending {
		h.Delete()
	}
	clear(c.pending)
	c.mu.Unlock()
}
