package rcheevos

import (
	"testing"
)

func TestHashFromBuffer(t *testing.T) {
	// Empty buffer should return empty hash
	hash := HashFromBuffer(ConsoleMasterSystem, nil)
	if hash != "" {
		t.Errorf("HashFromBuffer(nil) = %q, want empty", hash)
	}

	hash = HashFromBuffer(ConsoleMasterSystem, []byte{})
	if hash != "" {
		t.Errorf("HashFromBuffer([]) = %q, want empty", hash)
	}

	// Test with some data - SMS uses MD5 of the ROM
	// A simple 32KB ROM filled with 0xFF
	testROM := make([]byte, 32*1024)
	for i := range testROM {
		testROM[i] = 0xFF
	}

	hash = HashFromBuffer(ConsoleMasterSystem, testROM)
	if hash == "" {
		t.Error("HashFromBuffer should return a hash for valid ROM data")
	}

	// Hash should be 32 hex characters (MD5)
	if len(hash) != 32 {
		t.Errorf("hash length = %d, want 32", len(hash))
	}

	// Hash should be consistent
	hash2 := HashFromBuffer(ConsoleMasterSystem, testROM)
	if hash != hash2 {
		t.Errorf("hash not consistent: %q vs %q", hash, hash2)
	}
}

func TestHashFromBufferDifferentData(t *testing.T) {
	rom1 := make([]byte, 32*1024)
	rom2 := make([]byte, 32*1024)

	for i := range rom1 {
		rom1[i] = 0x00
	}
	for i := range rom2 {
		rom2[i] = 0xFF
	}

	hash1 := HashFromBuffer(ConsoleMasterSystem, rom1)
	hash2 := HashFromBuffer(ConsoleMasterSystem, rom2)

	if hash1 == hash2 {
		t.Error("different ROMs should produce different hashes")
	}
}

func TestHashFromBufferDifferentConsoles(t *testing.T) {
	// Same ROM data may produce different hashes for different consoles
	// because hashing algorithms vary by system
	testROM := make([]byte, 32*1024)
	for i := range testROM {
		testROM[i] = byte(i & 0xFF)
	}

	hashSMS := HashFromBuffer(ConsoleMasterSystem, testROM)
	hashGG := HashFromBuffer(ConsoleGameGear, testROM)

	// SMS and Game Gear use the same hashing, so they should match
	// But other consoles may differ
	if hashSMS == "" || hashGG == "" {
		t.Error("hashes should not be empty")
	}
}

func TestHashIteratorFromBuffer(t *testing.T) {
	testROM := make([]byte, 32*1024)
	for i := range testROM {
		testROM[i] = byte(i & 0xFF)
	}

	iter := NewHashIteratorFromBuffer("test.sms", testROM)
	if iter == nil {
		t.Fatal("NewHashIteratorFromBuffer returned nil")
	}
	defer iter.Destroy()

	// Should be able to generate at least one hash
	hash := iter.Next()
	if hash == "" {
		t.Error("iterator should produce at least one hash")
	}

	// Hash should be 32 characters
	if len(hash) != 32 {
		t.Errorf("hash length = %d, want 32", len(hash))
	}
}

func TestHashIteratorGenerateHash(t *testing.T) {
	testROM := make([]byte, 32*1024)
	for i := range testROM {
		testROM[i] = byte(i & 0xFF)
	}

	iter := NewHashIteratorFromBuffer("test.sms", testROM)
	if iter == nil {
		t.Fatal("NewHashIteratorFromBuffer returned nil")
	}
	defer iter.Destroy()

	// Generate hash for specific console
	hash := iter.GenerateHash(ConsoleMasterSystem)
	if hash == "" {
		t.Error("GenerateHash should produce a hash")
	}

	// Should match HashFromBuffer
	directHash := HashFromBuffer(ConsoleMasterSystem, testROM)
	if hash != directHash {
		t.Errorf("GenerateHash = %q, HashFromBuffer = %q", hash, directHash)
	}
}

func TestHashIteratorDestroy(t *testing.T) {
	testROM := make([]byte, 1024)
	iter := NewHashIteratorFromBuffer("test.sms", testROM)

	// Should not panic on destroy
	iter.Destroy()

	// Should be safe to call destroy again
	iter.Destroy()
}

func TestHashFromFileNonExistent(t *testing.T) {
	hash := HashFromFile(ConsoleMasterSystem, "/nonexistent/path/to/rom.sms")
	if hash != "" {
		t.Errorf("HashFromFile(nonexistent) = %q, want empty", hash)
	}
}

func TestNewHashIteratorNonExistent(t *testing.T) {
	iter := NewHashIterator("/nonexistent/path/to/rom.sms")
	if iter == nil {
		t.Fatal("NewHashIterator returned nil")
	}
	defer iter.Destroy()

	// Should return empty hash for non-existent file
	hash := iter.Next()
	if hash != "" {
		t.Errorf("iterator should return empty hash for non-existent file, got %q", hash)
	}
}
