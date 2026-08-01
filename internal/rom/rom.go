// Package rom loads and identity-checks the user-supplied Final Fantasy III
// (USA) ROM image that every extractor reads from.
//
// The ROM itself is never committed; the repository tracks only its SHA-256.
// Extraction refuses to run against any other revision: all tracked offsets
// and every piece of recorded evidence are specific to this image
// (docs/research/ROM_IDENTITY.md).
package rom

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"os"
)

// SupportedSHA256 is the only ROM revision the project's evidence covers:
// Final Fantasy III (USA), 3,145,728 bytes, no copier header.
const SupportedSHA256 = "0f51b4fca41b7fd509e4b8f9d543151f68efa5e97b08493e4b2a0c06f5d8d5e2"

// Size is the exact byte length of the supported image.
const Size = 3145728

// EnvVar names the environment variable consulted when no explicit path is
// given (the same convention the Mesen bridge uses for FF6_OUT).
const EnvVar = "FF6_ROM"

// ROM is a verified ROM image. The zero value is not usable; obtain one via
// Load.
type ROM struct {
	data   []byte
	sha256 string
}

// ResolvePath returns the explicit path when non-empty, else the FF6_ROM
// environment variable, else an error telling the operator what to set.
func ResolvePath(explicit string) (string, error) {
	if explicit != "" {
		return explicit, nil
	}
	if p := os.Getenv(EnvVar); p != "" {
		return p, nil
	}
	return "", fmt.Errorf("no ROM path: pass -rom <path> or set %s", EnvVar)
}

// Load reads and verifies a ROM image. It refuses wrong sizes and wrong
// hashes outright rather than extracting garbage from an unsupported
// revision.
func Load(path string) (*ROM, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("reading ROM: %w", err)
	}
	if len(data) != Size {
		return nil, fmt.Errorf("unsupported ROM: %d bytes, want %d (headerless FF3us)", len(data), Size)
	}
	sum := sha256.Sum256(data)
	got := hex.EncodeToString(sum[:])
	if got != SupportedSHA256 {
		return nil, fmt.Errorf("unsupported ROM revision: sha256 %s, supported %s", got, SupportedSHA256)
	}
	return &ROM{data: data, sha256: got}, nil
}

// SHA256 returns the verified image hash.
func (r *ROM) SHA256() string { return r.sha256 }

// Slice returns length bytes at the given file offset, bounds-checked.
// Callers use ROMFILE offsets (docs/research/ADDRESS_NOTATION.md).
func (r *ROM) Slice(off, length int) ([]byte, error) {
	if off < 0 || length < 0 || off+length > len(r.data) {
		return nil, fmt.Errorf("ROM slice out of range: off=0x%X len=0x%X size=0x%X", off, length, len(r.data))
	}
	out := make([]byte, length)
	copy(out, r.data[off:off+length])
	return out, nil
}
