# Latest Checkpoint

**[2026-08-02 — Units 10-14, the compression premise tested and refuted](2026-08-02-units10-14-content-parity-session.md)**
(preceding: [Unit 10, records reconciled](2026-08-02-unit10-records-reconciled.md))

State: branch `demo/whelk-content-parity`, six commits ahead of `main` at
`297ba88`, worktree clean, not pushed. No emulator running, no resident
instrumentation, no background processes.

**The session's result is not the code it wrote. It is that the program's
stated keystone was wrong, and the evidence to show that was already on disk.**

Readiness X1 said FF6's compression format gates maps, field sprites, battle
backgrounds, enemy graphics and party sprites. **EXP-0050 tested it.** A
verbatim search against preserved VRAM — no emulator, no new capture — finds
the BG tile data for milestones 02, 04 and 05 present in the ROM
**uncompressed**, in three contiguous spans totalling 20 KB:
`ROMFILE:0x208460`, `0x223000`, `0x224F00`, the first two **shared between the
Narshe exterior and the mines interior**. Verbatim coverage of the whole 64 KB
image: field 47-52 %, battle 38 %, Mode 7 opening 0 %.

Compression was **withdrawn** from the route matrix's pressure table rather
than re-ranked: 48-62 % is unmatched, but a verbatim search is defeated by
bit-plane reordering, runtime composition and WRAM-built tilemaps as readily as
by compression, so ranking it on the unmatched fraction would repeat the
original error inverted. **Map headers now lead at 6 beats**, and F1 has
concrete anchors for the first time.

That was possible because Unit 12 exposed what `internal/mesenstate` had been
parsing and hiding since it was written: `ppu.vram`, `ppu.cgram`, `ppu.oamRam`,
`spc.ram` and the full PPU/DMA register state, in every preserved savestate.
`ff6lab state` reads all of them now, plus `state ppu` and `state origin`.

**A parity-blocking defect was found and fixed.** `BlitTile` adds `PaletteBase`
to ink values 1-3; both scenes passed `white = 3` / `gray = 2` as if it were a
brightness level; `GrayPalette` defined only entries 0-3. Measured on the real
font: **32 % of drawn ink resolved to a visible colour** — two thirds of what
the demo drew was invisible, since Unit 4. The frame goldens structurally could
not see it, because `Sum256` hashes indices and is deliberately
palette-independent. `content.SubPalette` makes the units explicit, and the
same mistake now fails a test. Visible ink 32 % → 63 %, ink mask byte-identical.

**Three records disagreed with themselves**, all D1's failure mode: the
readiness summary vs its own tables (53 claimed, 55 actual, since `969b5dd`,
propagated into four checkpoints and three dashboards);
`MESEN_CAPABILITY_MATRIX.md` recording VRAM/CGRAM/OAM/ARAM as `Unknown` while
`bridge.lua` had them wired and two experiments had used them; and a stale
`STATISTICS.md`. All recounted from source, and the readiness comparison is now
`audit.CheckReadinessSummary`, which caught a bad status token this same
session had introduced.

Gates: gofmt clean; build/vet/test on both variants; `ff6lab audit` clean
(eleven checks); census clean; archive verify 8/8; restricted scan clean.

Exact next action: **Unit 15 — find the map header record that selects
`ROMFILE:0x208460` and `0x223000`.** Two field maps share those two blocks and
differ in a third, which is a strong discriminator for the selecting record.
Falsifier: no candidate reproduces both maps' block sets. Cheaper fallback if
it stalls: **read OAM from milestone 04** — every savestate carries it,
`ff6lab state oam` reaches it, and no experiment has ever looked.
