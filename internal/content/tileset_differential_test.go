package content_test

// The differential that makes ASSET-GFX-0002 evidence rather than a guess:
// the extracted tileset must equal the VRAM a preserved Mesen capture holds.
//
// It skips without the savestate corpus, so CI never needs one. Run it with
// the archive present:
//
//	go test ./internal/content/ -run TestNarsheTilesetMatchesCapturedVRAM

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/bensabler/ff6-decompile/internal/content"
	"github.com/bensabler/ff6-decompile/internal/graphics/tile4bpp"
	"github.com/bensabler/ff6-decompile/internal/mesenstate"
)

// capturedStates are the two SCN-0001 milestones EXP-0050 found loading this
// block at VRAM:$0000. Two independent captures, not one.
var capturedStates = []string{
	"local_artifacts/scenarios/SCN-0001/04-free-movement/run1-04.mss",
	"local_artifacts/scenarios/SCN-0001/02-narshe-entry/run1-02.mss",
}

func TestNarsheTilesetMatchesCapturedVRAM(t *testing.T) {
	root := repoRoot(t)
	a, err := content.Open(root)
	if err != nil {
		t.Skipf("no archive: %v", err)
	}
	ts, err := content.LoadNarsheTileset(a)
	if err != nil {
		t.Skipf("tileset not in the archive: %v", err)
	}

	for _, rel := range capturedStates {
		path := filepath.Join(root, rel)
		if _, err := os.Stat(path); err != nil {
			t.Skipf("no preserved savestate at %s", rel)
		}
		t.Run(filepath.Base(rel), func(t *testing.T) {
			st, err := mesenstate.Open(path)
			if err != nil {
				t.Fatal(err)
			}
			vram, err := st.VideoRAM()
			if err != nil {
				t.Fatal(err)
			}
			// The block occupies VRAM:$0000 upward, one 4bpp tile per 32 B.
			for n := 0; n < content.TilesetTiles; n++ {
				want, err := tile4bpp.Decode(vram[n*tile4bpp.EncodedSize : (n+1)*tile4bpp.EncodedSize])
				if err != nil {
					t.Fatalf("tile %d from VRAM: %v", n, err)
				}
				got, _ := ts.Tile(n)
				if *got != want {
					t.Fatalf("tile %d differs between the archive and captured VRAM\n got  %v\n want %v",
						n, got[0], want[0])
				}
			}
		})
	}
}
