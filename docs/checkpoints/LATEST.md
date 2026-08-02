# Latest Checkpoint

**[2026-08-02 — Unit 10, records reconciled and the route matrix built](2026-08-02-unit10-records-reconciled.md)**
(preceding: [DEMO-0001A complete, the demo runs](2026-08-02-demo0001a-shell-complete.md))

State: DEMO-0001 is in its third phase, **field/event content recovery**, on
branch `demo/whelk-content-parity`. The foundation (units 0–9) is merged to
`main` at `297ba88`, tagged `demo-0001-foundation-v0.1`. Worktree clean, no
emulator running, no resident instrumentation. SCN-0001 remains the evidence
program that DEMO-0001 consumes.

Units 8 and 9 landed in that merge but were never checkpointed, so the
checkpoint chain and `CURRENT_FOCUS.md` still said "next action: Unit 8" while
the readiness matrix and deviations register described a completed Unit 9. This
unit closed that gap and three more found alongside it.

**The route content matrix now puts a number on the critical path.**
`docs/demo/DEMO-0001-CONTENT-MATRIX.md` is keyed by the 19 SCN-0001 beats
rather than by subsystem, and counts how many beats each unresolved dependency
blocks. **Compression (X1) blocks 8 of 19** — B02, B03, B06, B07, B09, B11,
B14, B18 — more than any other single dependency, and it still has zero records
and zero code. Map headers block 6; field sprites and music sequences 5 each;
event dispatch 4; the dialogue corpus 3. The matrix also surfaced two
requirements nobody owned, now added to readiness: **B19 action animations** and
**X4 transition effects**.

**The finding that re-sequences the program: the evidence is already frozen.**
Every preserved Mesen savestate carries `ppu.vram` (64 KiB), `ppu.cgram`,
`ppu.oamRam`, `spc.ram` (64 KiB ARAM), `memoryManager.workRam` and the full
PPU/DMA register state. `internal/mesenstate` parses those files today but
exposes only WRAM and SRAM. The corpus covers field, battle and post-victory
states — so **a decompressed FF6 tileset is already in hand, hashed, and
reproducible with no emulator**, which is exactly the known-good output a
compression recovery needs on the far side of its falsifier. Maps, sprites,
backgrounds and the audio driver become analysis rather than sessions.
EXP-0048 does need a live session and is re-ordered behind that work — deferred,
not blocked.

**Three counting errors, corrected from source, all the same failure mode as
retired deviation D1** — two records of one fact, disagreeing, with nothing
comparing them. (1) The readiness summary claimed 53 requirements when the file
has carried 55 since `969b5dd`, with Unknown 33 vs 36 and Evidence Ready 7 vs
6; the wrong figure had propagated into four checkpoints and three dashboards.
Both columns are now recounted: **57 rows — 1 Validated, 14 Integrated, 6
Implemented, 29 Unknown**. (2) `MESEN_CAPABILITY_MATRIX.md` recorded VRAM,
CGRAM, OAM, ARAM and DSP access as `Unknown` although `bridge.lua:238` has had
all five wired since it was written and two experiments had already used them —
work was being deferred as "needs instrumentation" when the instrument existed.
(3) `STATISTICS.md` was stale on every count. Rows now distinguish
wired-and-exercised, wired-but-never-exercised (OAM — never once read), and
genuinely absent (live DMA register capture, which no probe or Go path
implements despite a command, a skill, an agent and a playbook for it).

No Go code changed. All gates green.

Exact next action: **Unit 11 — fix the palette defect that makes the demo's own
text invisible.** `BlitTile` adds `PaletteBase` to ink values 1–3, both scenes
pass `white = 3` / `gray = 2`, and `GrayPalette()` defines only indices 1–3, so
every "white" string resolves to black on black. `Sum256` hashes indices and is
deliberately palette-independent, so the frame goldens structurally cannot see
it. Fix the convention, add the missing assertion, regenerate the goldens, and
record the project-authored palette as a deviation retired by F2.
