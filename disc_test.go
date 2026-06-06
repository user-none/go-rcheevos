package rcheevos

import (
	"crypto/md5"
	"encoding/hex"
	"testing"
)

// fakeCD is a minimal CDReader: one Mode 1 data track whose sector 0 holds a
// Saturn IP header. It exercises the cdreader bridge without a real image.
type fakeCD struct {
	sector0 []byte
}

func (f *fakeCD) ReadSectorData(lba int) ([]byte, error) {
	if lba == 0 {
		return f.sector0, nil
	}
	return make([]byte, 2048), nil
}

func (f *fakeCD) NumTracks() int { return 1 }

func (f *fakeCD) Track(i int) (number int, typ string, frames int, pregap int, startLBA int, control uint8) {
	return 1, "MODE1_RAW", 1000, 0, 0, 0x41
}

// TestHashFromDisc drives the full Go->C->Go cdreader path: rcheevos opens
// track 1, reads 512 bytes of sector 0 through the bridge, and MD5s them. The
// Saturn hash is the MD5 of those first 512 user-data bytes.
func TestHashFromDisc(t *testing.T) {
	user := make([]byte, 2048)
	copy(user, []byte("SEGA SEGASATURN "))
	for i := 16; i < len(user); i++ {
		user[i] = byte(i * 7)
	}

	sum := md5.Sum(user[:512])
	want := hex.EncodeToString(sum[:])

	got := HashFromDisc(ConsoleSaturn, &fakeCD{sector0: user})
	if got != want {
		t.Fatalf("HashFromDisc = %q, want %q", got, want)
	}
}

// TestHashFromDiscNotSaturn confirms a non-Saturn sector 0 fails identification
// (the signature check in rc_hash_sega_cd), yielding an empty hash.
func TestHashFromDiscNotSaturn(t *testing.T) {
	user := make([]byte, 2048) // no "SEGA SEGASATURN " signature
	if got := HashFromDisc(ConsoleSaturn, &fakeCD{sector0: user}); got != "" {
		t.Fatalf("HashFromDisc on non-Saturn disc = %q, want empty", got)
	}
}
