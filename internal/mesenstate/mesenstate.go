// Package mesenstate reads Mesen 2 savestate (.mss) files without an
// emulator, so preserved states can be mined as evidence offline.
//
// A savestate is the ASCII signature "MSS", a short header, and one or
// more zlib streams. One of those streams is a flat sequence of named
// blocks, each encoded as a NUL-terminated name, a little-endian uint32
// length, and that many bytes of data — for example
// "memoryManager.workRam" carrying 131072 bytes of work RAM, or
// "cart.saveRam" carrying 8192 bytes of cartridge save RAM.
//
// Savestates are ROM-derived and never committed; this package only
// reads local artifacts.
package mesenstate

import (
	"bytes"
	"compress/zlib"
	"encoding/binary"
	"errors"
	"fmt"
	"io"
	"os"
	"sort"
)

// signature is the ASCII magic every Mesen savestate begins with.
const signature = "MSS"

// Block names the project reads.
const (
	WorkRAMBlock  = "memoryManager.workRam"
	SaveRAMBlock  = "cart.saveRam"
	VideoRAMBlock = "ppu.vram"
	CGRAMBlock    = "ppu.cgram"
	OAMBlock      = "ppu.oamRam"
	AudioRAMBlock = "spc.ram"
)

// Expected sizes of those blocks on SNES.
const (
	WorkRAMSize  = 128 * 1024
	SaveRAMSize  = 8 * 1024
	VideoRAMSize = 64 * 1024
	CGRAMSize    = 512 // 256 BGR555 entries
	// OAMSize is 512 bytes of four-byte entries plus the 32-byte high
	// table that carries each sprite's size bit and X sign bit.
	OAMSize      = 544
	AudioRAMSize = 64 * 1024
)

// maxSection bounds decompression so a corrupt or hostile file cannot
// exhaust memory.
const maxSection = 64 << 20

// State is the named-block section of a parsed savestate.
type State struct {
	blocks map[string][]byte
	names  []string
}

// Open reads and parses the savestate at path.
func Open(path string) (*State, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("mesenstate: read %s: %w", path, err)
	}
	st, err := Parse(data)
	if err != nil {
		return nil, fmt.Errorf("mesenstate: %s: %w", path, err)
	}
	return st, nil
}

// Parse locates the named-block section inside a savestate image.
//
// The header does not record where that section begins, so Parse walks
// candidate zlib streams and keeps the first one that decodes cleanly
// into named blocks and consumes its section exactly. A stream carrying
// anything else — the preview image, for instance — fails that check and
// is skipped.
func Parse(data []byte) (*State, error) {
	if len(data) < len(signature) || string(data[:len(signature)]) != signature {
		return nil, fmt.Errorf("not a Mesen savestate (missing %q signature)", signature)
	}
	for off := len(signature); off < len(data); off++ {
		if data[off] != 0x78 { // zlib CMF for the deflate window Mesen uses
			continue
		}
		section, consumed, err := inflateAt(data, off)
		if err != nil {
			continue
		}
		if st, err := parseBlocks(section); err == nil {
			return st, nil
		}
		if consumed > 1 {
			off += consumed - 1 // skip the stream we just decoded
		}
	}
	return nil, errors.New("no named-block section found")
}

// inflateAt decompresses the zlib stream starting at off, returning the
// section and how many input bytes it spanned. bytes.Reader is an
// io.ByteReader, so flate reads no further than the stream itself and
// the remaining length is exact.
func inflateAt(data []byte, off int) ([]byte, int, error) {
	r := bytes.NewReader(data[off:])
	zr, err := zlib.NewReader(r)
	if err != nil {
		return nil, 0, err
	}
	defer zr.Close()
	var buf bytes.Buffer
	if _, err := io.Copy(&buf, io.LimitReader(zr, maxSection)); err != nil {
		return nil, 0, err
	}
	return buf.Bytes(), len(data) - off - r.Len(), nil
}

// parseBlocks walks the flat name/length/data sequence. It requires the
// section to be consumed exactly and every name to be printable, which
// is what distinguishes the real block section from an arbitrary stream
// that merely happens to inflate.
func parseBlocks(section []byte) (*State, error) {
	st := &State{blocks: make(map[string][]byte)}
	dotted := false
	for pos := 0; pos < len(section); {
		end := bytes.IndexByte(section[pos:], 0)
		if end < 1 {
			return nil, fmt.Errorf("block name missing or empty at offset %d", pos)
		}
		name := section[pos : pos+end]
		if !printableASCII(name) {
			return nil, fmt.Errorf("block name is not printable ASCII at offset %d", pos)
		}
		if bytes.IndexByte(name, '.') >= 0 {
			dotted = true
		}
		lenAt := pos + end + 1
		if lenAt+4 > len(section) {
			return nil, fmt.Errorf("block %q: truncated length field", name)
		}
		size := uint64(binary.LittleEndian.Uint32(section[lenAt : lenAt+4]))
		body := lenAt + 4
		if uint64(body)+size > uint64(len(section)) {
			return nil, fmt.Errorf("block %q: declares %d bytes, past end of section", name, size)
		}
		key := string(name)
		if _, dup := st.blocks[key]; !dup {
			st.names = append(st.names, key)
		}
		st.blocks[key] = section[body : body+int(size)]
		pos = body + int(size)
	}
	if len(st.blocks) == 0 || !dotted {
		return nil, errors.New("no named blocks")
	}
	return st, nil
}

func printableASCII(b []byte) bool {
	for _, c := range b {
		if c < 0x20 || c > 0x7E {
			return false
		}
	}
	return true
}

// Names returns every block name, sorted.
func (s *State) Names() []string {
	out := append([]string(nil), s.names...)
	sort.Strings(out)
	return out
}

// Block returns a block's bytes. The slice aliases the decoded section;
// callers must not modify it.
func (s *State) Block(name string) ([]byte, bool) {
	b, ok := s.blocks[name]
	return b, ok
}

// WorkRAM returns the 128 KB SNES work RAM image, the domain the project
// writes as WRAM:+$xxxx.
func (s *State) WorkRAM() ([]byte, error) { return s.region(WorkRAMBlock, WorkRAMSize) }

// SaveRAM returns the cartridge save RAM image.
func (s *State) SaveRAM() ([]byte, error) { return s.region(SaveRAMBlock, SaveRAMSize) }

// VideoRAM returns the 64 KB VRAM image: the tile and tilemap data the PPU
// was actually drawing from when the state was written.
//
// For any graphics format this project has yet to decode, this is the
// **known-good output side of the comparison** — the bytes a ROM decompressor
// must reproduce. It needs no emulator, which is what makes a preserved
// savestate corpus usable as evidence long after the session ended.
func (s *State) VideoRAM() ([]byte, error) { return s.region(VideoRAMBlock, VideoRAMSize) }

// CGRAM returns the 512-byte palette image, 256 BGR555 entries. Decode
// entries with platform/bgr555.
func (s *State) CGRAM() ([]byte, error) { return s.region(CGRAMBlock, CGRAMSize) }

// OAM returns the 544-byte sprite attribute image: 128 four-byte entries
// followed by the 32-byte high table.
func (s *State) OAM() ([]byte, error) { return s.region(OAMBlock, OAMSize) }

// AudioRAM returns the 64 KB SPC700 address space: the loaded audio driver,
// its sequence data and its BRR samples, as they stood at capture time.
func (s *State) AudioRAM() ([]byte, error) { return s.region(AudioRAMBlock, AudioRAMSize) }

func (s *State) region(name string, want int) ([]byte, error) {
	b, ok := s.blocks[name]
	if !ok {
		return nil, fmt.Errorf("mesenstate: savestate has no %q block", name)
	}
	if len(b) != want {
		return nil, fmt.Errorf("mesenstate: %q is %d bytes, want %d", name, len(b), want)
	}
	return b, nil
}
