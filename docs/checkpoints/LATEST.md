# Latest Checkpoint

**[2026-08-02 — DEMO-0001A complete, the demo runs](2026-08-02-demo0001a-shell-complete.md)**
(preceding: [CORR-0001 pointer advance at ROMCPU:$C09B5C](2026-08-01-corr0001-c09b5c-pointer-advance.md))

State: the project has shifted into **playable-vertical-slice production**.
`DEMO-0001 — New Game to Whelk Victory` is the active program, on branch
`demo/new-game-to-whelk` (8 commits, worktree clean, not pushed). SCN-0001
remains the evidence program that DEMO-0001 consumes.

**The demo runs.** `go run -tags gui ./cmd/ff6demo` opens a window;
`go run ./cmd/ff6demo -headless -frames 120 -capture-last` needs no display and
is the authoritative mode for any frame-exact claim. It renders real FF6 text
from the battle HUD font extracted from the operator's own ROM. The binary
**cannot read a ROM**: `internal/rom` is barred transitively from `cmd/ff6demo`,
`internal/content` and `internal/engine`, and `ff6lab audit` walks the import
graph to enforce it.

Two findings are worth carrying forward.

**A parity-blocking defect, found and fixed.** The `hud-font` extractor read
`ROMFILE:0x046FC0` as a block start. That address is only the back-projected
anchor for VRAM tile `$000`; the real block starts at `0x047FB0`. **255 of the
257 tiles in the shipped font sheet were attack-table bytes rendered as tiles**,
and had been since 2026-07-30. `manifests/rom-regions.json` ROM-0016 recorded
the correct span the whole time — the extractor and the ledger were two
independent records of one fact that disagreed, and nothing compared them. Now
`TestHUDFontMatchesROMLedger` does, and a skip-guarded archive-vs-ROM
differential covers the general case across all 256 glyphs.

**EXP-0049 closed a mapping the census carried as Unknown**, from tracked
evidence with no emulator run: `VRAM tile = $100 + encoded byte`. The
discriminating trial decoded the whole EXP-0023 HUD tilemap through the
relation and got coherent game text across all 37 referenced tiles —
"Were-Rat", "Repo Man", "WEDGE", "VICKS", "MagiTek", "Item", "Row", "Tek", HP
digits. That also promotes `textenc` itself from "derived from a menu tilemap"
to "verified against rendered output". Secondary: `$BF` = `?`; five HP/ATB
gauge tiles identified (CEN-GFX-0005); the HUD layout registered
(CEN-BATTLE-0014); **47 tiles recorded as deliberately unidentified**, with a
test asserting the classifier never guesses one.

Readiness: **7 of 53 requirements Integrated, 1 Validated, 10 Implemented** —
from 0 Integrated at program start. Every Integrated row is engine or text
plumbing. No Field or Audio row has moved, and the map system, dialogue corpus,
event opcode table, sprite formats, audio sequences and FF6's compression
format all remain at zero records. The demo has a spine, not yet a game, and
`docs/demo/DEMO-0001-READINESS.md` says so per requirement.

ADR-0001 adopts Ebitengine, confined to `internal/engine/ebitenhost` (enforced
as a direct-import rule) and behind a `gui` build tag. Measured: v2.9.9 needs
cgo on **both** macOS and Linux, so the tag is what keeps every default gate
free of system libraries. The headless host and the whole test surface landed
one commit *before* ebiten entered `go.mod`, which is what makes the dependency
provably optional.

Nothing in flight: no emulator, no resident instrumentation, worktree clean.
All gates green on both build variants, including `go mod verify`, and the
restricted-file scan now also rejects tracked rendered images.

Exact next action: **Unit 8 — decode formations and monster records into Go**
(`internal/content/battledata`) from the already-archived `formations.bin` and
`monsters.bin`. Both formats are Confirmed (EXP-0028/0029/0030). Then Unit 9, a
battle HUD scene consuming the font, the ATB layer and those tables.
