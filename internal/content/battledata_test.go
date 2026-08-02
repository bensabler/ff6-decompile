package content_test

import (
	"testing"

	"github.com/bensabler/ff6-decompile/internal/content"
	"github.com/bensabler/ff6-decompile/internal/game/battledata"
)

// openTables loads the real archive, or skips. The archive is generated from
// the operator's own ROM and is never committed, so these are opt-in.
func openTables(t *testing.T) *content.BattleTables {
	t.Helper()
	a, err := content.Open(repoRoot(t))
	if err != nil {
		t.Skipf("no generated archive: %v", err)
	}
	tb, err := content.LoadBattleTables(a)
	if err != nil {
		t.Fatalf("loading battle tables: %v", err)
	}
	return tb
}

// TestArchivedTablesMatchTheRecordedFacts is the differential for this unit:
// it replays every formation and monster fact the experiments established and
// checks the decoders reproduce them from the archive.
//
// A decoder that agrees with itself proves nothing. These numbers come from
// live battle state, not from re-reading the same bytes.
func TestArchivedTablesMatchTheRecordedFacts(t *testing.T) {
	tb := openTables(t)

	t.Run("table sizes", func(t *testing.T) {
		if got := tb.Formations.Count(); got != 576 {
			t.Errorf("formations = %d, want 576 (the archived span)", got)
		}
		if got := tb.Monsters.Count(); got != 384 {
			t.Errorf("monsters = %d, want 384", got)
		}
		if got := tb.FormationFlags.Count(); got != 576 {
			t.Errorf("formation flags = %d, want 576", got)
		}
	})

	// Monster stats measured live and recorded in CURRENT_FOCUS / EXP-0028.
	t.Run("monster stats", func(t *testing.T) {
		for _, tt := range []struct {
			id     int
			power  uint8
			hp, mp uint16
			note   string
		}{
			{0, 16, 40, 15, "opening guard, EXP-0032/0033"},
			{19, 13, 24, 0, "mines encounter; power 13 is EXP-0018's live pair"},
			{25, 20, 27, 5, "scripted battle, EXP-0034"},
			{27, 110, 115, 30, "scripted battle 5, EXP-0039"},
			{77, 19, 35, 0, "mines encounter; power 19 is EXP-0018's live pair"},
		} {
			m, err := tb.Monster(tt.id)
			if err != nil {
				t.Fatalf("monster %d: %v", tt.id, err)
			}
			if m.HP() != tt.hp || m.MP() != tt.mp {
				t.Errorf("monster %d (%s): HP=%d MP=%d, want %d/%d", tt.id, tt.note, m.HP(), m.MP(), tt.hp, tt.mp)
			}
			if m.Power() != tt.power {
				t.Errorf("monster %d (%s): power=%d, want %d", tt.id, tt.note, m.Power(), tt.power)
			}
		}
	})

	// Formation compositions established by EXP-0030/0034/0038/0039.
	t.Run("formation compositions", func(t *testing.T) {
		for _, tt := range []struct {
			id   int
			ids  []byte
			note string
		}{
			{2, []byte{0, 0}, "opening scripted battles 1 and 3"},
			{1, []byte{25, 25}, "opening scripted battle 2"},
			{41, []byte{25, 0, 0}, "opening scripted battle 4"},
			{84, []byte{27, 27, 0, 0}, "scripted battle 5, en route to the shaft"},
			{44, []byte{19, 77}, "mines encounter, EXP-0030's verified record"},
			{14, []byte{19, 19, 19}, "mines encounter, EXP-0038"},
		} {
			f, err := tb.Formation(tt.id)
			if err != nil {
				t.Fatalf("formation %d: %v", tt.id, err)
			}
			got := f.OccupiedIDs()
			if len(got) != len(tt.ids) {
				t.Errorf("formation %d (%s): ids %v, want %v", tt.id, tt.note, got, tt.ids)
				continue
			}
			for i := range got {
				if got[i] != tt.ids[i] {
					t.Errorf("formation %d (%s): ids %v, want %v", tt.id, tt.note, got, tt.ids)
					break
				}
			}
		}
	})

	// The staged Whelk record, byte-for-byte from EXP-0039/B18.
	t.Run("whelk record bytes", func(t *testing.T) {
		f, err := tb.Formation(432)
		if err != nil {
			t.Fatal(err)
		}
		want := []byte{0x80, 0x03, 0x00, 0x34, 0xFF, 0xFF, 0xFF, 0xFF,
			0x48, 0xAB, 0x00, 0x00, 0x00, 0x00, 0x3F}
		for i := range want {
			if f[i] != want[i] {
				t.Fatalf("formation 432 = % 02X, want % 02X (the record SCN-0001 B18 staged)", f[:], want)
			}
		}
	})
}

// TestLeadingWordIsNotTheMonsterIDExtension records a negative result.
//
// Whelk is formation 432, whose id bytes read as records 0 and 52. That cannot
// be right: record 0 is the 40 HP opening guard, while EXP-0040 measured
// Whelk's two entities at 50000 and 1600 HP. Something outside the id bytes
// selects a different record, and the leading word at +$00/+$01 was the
// obvious candidate — it is the only other field in the record's first half.
//
// It is refuted. Formation 1 carries the *same* leading word, $0380, and its
// monsters are Confirmed to be record 25 (27 HP / 5 MP, matched against live
// battle state in EXP-0034). If the word added a high bit to the ids,
// formation 1 would resolve to record 281 and the opening battle would not
// match.
//
// This narrows the search rather than closing it: the extension, if it is one,
// lives in the trailing bytes +$08..+$0E, in the parallel flags table at
// ROMCPU:$CF5900, or in the loader's X-computation (CEN-MONSTER-0001 lists
// "X-computation (formation->monster-id mapping)" as unknown).
func TestLeadingWordIsNotTheMonsterIDExtension(t *testing.T) {
	tb := openTables(t)

	whelk, err := tb.Formation(432)
	if err != nil {
		t.Fatal(err)
	}
	opening, err := tb.Formation(1)
	if err != nil {
		t.Fatal(err)
	}

	if whelk.LeadingWord() != opening.LeadingWord() {
		t.Fatalf("this refutation rests on formations 1 and 432 sharing a leading word; "+
			"they now read $%04X and $%04X, so re-examine the hypothesis",
			opening.LeadingWord(), whelk.LeadingWord())
	}

	// Formation 1's composition is Confirmed against live state, so the word
	// demonstrably does not shift its ids.
	got := opening.OccupiedIDs()
	if len(got) != 2 || got[0] != 25 || got[1] != 25 {
		t.Fatalf("formation 1 ids = %v, want [25 25]; EXP-0034 verified monster 25 live", got)
	}
	m, err := tb.Monster(int(got[0]))
	if err != nil {
		t.Fatal(err)
	}
	if m.HP() != 27 {
		t.Errorf("monster %d HP = %d, want 27; the low-byte reading is what EXP-0034 confirmed", got[0], m.HP())
	}
}

// TestFormationVerifiedExtentIsHonest guards the distinction the package
// draws between archived and verified. Whelk sits past the verified span, and
// a caller that forgets is the one who will report an unverified id as fact.
func TestFormationVerifiedExtentIsHonest(t *testing.T) {
	tb := openTables(t)
	if battledata.FormationVerified(432) {
		t.Error("formation 432 is archived but no experiment has verified its composition")
	}
	if tb.Formations.Count() <= battledata.FormationsVerified {
		t.Error("the archive should hold more formations than are verified; that gap is the point")
	}
}
