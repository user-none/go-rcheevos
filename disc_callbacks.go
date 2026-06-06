package rcheevos

/*
#include <stdint.h>
#include <stddef.h>
*/
import "C"
import (
	"runtime/cgo"
	"unsafe"
)

// goCDOpenTrack resolves the disc from its cgo.Handle, selects the requested
// track, and returns a cgo.Handle to a per-track cdTrack (0 on failure).
//
//export goCDOpenTrack
func goCDOpenTrack(discHandle C.uintptr_t, track C.uint32_t) C.uintptr_t {
	disc, ok := cgo.Handle(uintptr(discHandle)).Value().(CDReader)
	if !ok {
		return 0
	}
	start := trackStartSector(disc, uint32(track))
	h := cgo.NewHandle(&cdTrack{disc: disc, firstSector: start})
	return C.uintptr_t(uintptr(h))
}

// goCDReadSector fills buffer with cooked sector data and returns bytes read.
//
//export goCDReadSector
func goCDReadSector(handle C.uintptr_t, sector C.uint32_t, buffer unsafe.Pointer, requestedBytes C.size_t) C.size_t {
	if buffer == nil || requestedBytes == 0 {
		return 0
	}
	t, ok := cgo.Handle(uintptr(handle)).Value().(*cdTrack)
	if !ok {
		return 0
	}
	buf := unsafe.Slice((*byte)(buffer), int(requestedBytes))
	return C.size_t(t.read(int(sector), buf))
}

// goCDCloseTrack releases the per-track cgo.Handle.
//
//export goCDCloseTrack
func goCDCloseTrack(handle C.uintptr_t) {
	cgo.Handle(uintptr(handle)).Delete()
}

// goCDFirstTrackSector returns the track's first data sector.
//
//export goCDFirstTrackSector
func goCDFirstTrackSector(handle C.uintptr_t) C.uint32_t {
	t, ok := cgo.Handle(uintptr(handle)).Value().(*cdTrack)
	if !ok {
		return 0
	}
	return C.uint32_t(t.firstSector)
}
