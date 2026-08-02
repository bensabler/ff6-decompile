# Checkpoint 2026-08-02 — Unit 10, records reconciled and the route matrix built

## Current question

None open. The next unit is chosen, unblocked, and needs no emulator: fix the
palette defect that renders the demo's own text invisible.

## State

DEMO-0001 is in its third phase, **field/event content recovery**, on branch
`demo/whelk-content-parity`. The foundation (units 0–9) is merged to `main` at
`297ba88` and tagged `demo-0001-foundation-v0.1`.

Units 8 and 9 had landed in the merge but were never checkpointed, so the
checkpoint chain, `CURRENT_FOCUS.md` and the milestone unit log all still
pointed at "next action: Unit 8" while the readiness matrix and the deviations
register described a completed Unit 9. This unit closed that gap and three
others found alongside it.

## Confirmed before this session

Everything in the merged foundation: the battle damage pipeline
(DISC-0001…0007), event flags (DISC-0008), the ATB program (EXP-0041…0047),
formations `ROMCPU:$CF6200`, monster records `ROMCPU:$CF0000`, the HUD font
block (GFX-0001, ROM-0016), the glyph relation `VRAM tile = $100 + encoded
byte` (EXP-0049), and CORR-0001's pointer advance at `ROMCPU:$C09B5C`.

ROM identity re-verified this session: SHA-256
`0f51b4fca41b7fd509e4b8f9d543151f68efa5e97b08493e4b2a0c06f5d8d5e2`, unchanged,
matching `docs/research/ROM_IDENTITY.md`.

## Work completed

**A route content dependency matrix**, `docs/demo/DEMO-0001-CONTENT-MATRIX.md`.
Keyed by the 19 SCN-0001 beats rather than by subsystem, because the demo is
built in route order while research is done in evidence order. It carries a
register of 45 distinct content dependencies (CD-CMP/MAP/EVT/TXT/SPR/BAT/AUD/INI),
each with its readiness row, its evidence anchor, its recovery-pipeline state
and its blocker; a beat → dependency map; and a **dependency-pressure table**
counting how many of the 19 beats each unresolved dependency blocks. It
introduces no facts — every state is sourced from an experiment, discovery,
census entry or readiness row.

### The number the matrix produced

**Compression (readiness X1) blocks 8 of the 19 route beats** — B02, B03, B06,
B07, B09, B11, B14, B18 — more than any other single dependency. Map headers
(F1) block 6, field sprites (F6) and music sequences (A3) block 5 each, event
dispatch (F8) blocks 4, the dialogue corpus (F10) blocks 3.

That agrees with readiness X1's own prose assessment ("highest-leverage
unscheduled research") and is now a count rather than a judgement. It is why
compression is scheduled as the next research unit.

### Two requirements nobody owned

Building the beat map surfaced content the route needs and no readiness row
covered. Both were added rather than noted:

- **B19 — action animations and effects.** The matrix covered battle
  backgrounds, enemy graphics, party sprites and the HUD, but nothing owned
  attack, spell or damage-effect animation. Acceptance step 14, "supports a
  functional Whelk battle", cannot pass without it.
- **X4 — scene transition effects** (fades, battle-entry wipes), named by the
  route, owned by no requirement.

A third item is an action rather than a gap and is queued: **`ppu.oamRam` has
never been read out of any capture**, although every preserved savestate
contains it and `bridge.lua` has had `read oam` wired since it was written.
Readiness F6 depends on it.

### Three counting errors, corrected from source

This is the same failure mode as retired deviation D1 — two independent records
of one fact, disagreeing, with nothing comparing them.

1. **The readiness summary disagreed with its own table.** Unit 0 reported "53
   requirements, 33 Unknown, 7 Evidence Ready". Counted from `969b5dd`, the
   file has always carried **55 rows, 36 Unknown, 6 Evidence Ready**. The wrong
   figure propagated into the milestone record, `MILESTONES.md`,
   `CURRENT_FOCUS.md` and four checkpoints. Both summary columns are now
   recounted from the tables; the current state is **57 rows — 1 Validated, 14
   Integrated, 6 Implemented, 1 Evidence Ready, 1 Extractor Ready, 2
   Researching, 2 Blocked, 1 Deferred, 29 Unknown**.
2. **`MESEN_CAPABILITY_MATRIX.md` understated the project's own instrument.**
   VRAM, CGRAM, OAM, ARAM and DSP access were all recorded `Unknown`, although
   `mesen/bridge.lua:238`'s `READ_MEMTYPES` has wired every one of them since
   the bridge was written, and EXP-0023 (VRAM, CGRAM) and EXP-0024 (ARAM, DSP)
   had already used them. That is a planning hazard, not a cosmetic error:
   graphics and audio work was being deferred as "needs new instrumentation"
   when the instrument existed. Rows now distinguish **wired-and-exercised**,
   **wired-but-never-exercised** (OAM), and **absent** (live DMA register
   capture, which genuinely does not exist in any probe or Go path).
3. **`STATISTICS.md` was a day stale** on every count — 36 experiments vs 48,
   57 census entries vs 68, 26 ROM regions vs 32, 0.47 % vs 0.49 % ownership.
   Recounted from the records rather than edited.

`DEMO-0001-ACCEPTANCE.md`'s header also read "3 of 6 gates" while its own table
showed four PASS. The table was right; no gate changed state.

### The finding that re-sequences the program

**The runtime evidence for almost every missing content family is already
captured and frozen.** Every preserved Mesen savestate under
`local_artifacts/experiments/` and `mesen/out/` carries `ppu.vram` (64 KiB),
`ppu.cgram` (512 B), `ppu.oamRam` (544 B), `spc.ram` (64 KiB ARAM),
`memoryManager.workRam` (128 KiB), the full PPU register set and per-channel
`dmaController.*` state. `internal/mesenstate` parses those files already —
`ff6lab state list` enumerates every block — but exposes only WRAM and SRAM.

The corpus covers field (`EXP-0035/recon-mines-inside.mss`,
`mesen/out/checkpoint3-mines.mss`), battle, and post-victory states. A
**decompressed FF6 tileset is therefore already in hand, hashed, and
reproducible without an emulator**, which is exactly the known-good output a
compression recovery needs on the other side of its falsifier.

This converts maps, sprites, backgrounds and the audio driver from work that
needs a live session into work that needs analysis. EXP-0048 (name the
execution-path invoker) does need a session, so it is re-ordered behind the
offline units — deferred, not blocked.

## Last raw observation

`mesen/bridge.lua:238-247` — `READ_MEMTYPES` wiring `vram`, `cgram`, `oam`,
`aram`, `dsp`, read directly rather than taken from a summary.

`./ff6lab state list local_artifacts/experiments/EXP-0035/recon-mines-inside.mss`
— block inventory including `ppu.vram` 65536, `ppu.cgram` 512, `ppu.oamRam`
544, `spc.ram` 65536, `memoryManager.workRam` 131072.

## Active emulator state

None. No emulator running, no resident instrumentation, no savestate loaded.

## Breakpoints/watchers

None.

## Evidence paths

No new evidence was produced. This unit read existing tracked records,
`mesen/bridge.lua`, and the preserved savestate corpus under
`local_artifacts/` (gitignored, unmodified).

## Files changed

- `docs/demo/DEMO-0001-CONTENT-MATRIX.md` — **new**
- `docs/demo/DEMO-0001-READINESS.md` — B19 and X4 added; summary recounted with
  the correction recorded; route-view cross-reference; tail prose corrected
- `docs/demo/DEMO-0001-new-game-to-whelk.md` — unit log completed through
  Unit 10; branch updated; "exact next action" rewritten
- `docs/demo/DEMO-0001-ACCEPTANCE.md` — gate header reconciled with its table
- `docs/research/MESEN_CAPABILITY_MATRIX.md` — five capability rows corrected
- `dashboards/CURRENT_FOCUS.md` — branch, readiness figures, next action
- `dashboards/STATISTICS.md` — recounted
- `dashboards/MILESTONES.md` — D1 row
- `dashboards/RESEARCH_QUEUE.md` — P0 re-sequenced behind the compression
  keystone and the offline-access unit
- `docs/checkpoints/LATEST.md`, this file

No Go code changed.

## Tests and quality gates

`gofmt -l .` clean; `go build ./...`, `go vet ./...`, `go test ./...` pass;
`ff6lab audit` clean; `census validate` clean. Recorded in the commit.

## Git status

Branch `demo/whelk-content-parity`, cut from `main` at `297ba88`. One commit
for this unit. Not pushed.

## Unresolved decisions

None.

## Blockers

None for the next unit. The standing content blockers are unchanged and now
ranked by beat count: compression (8 beats), map headers (6), field sprites and
music sequences (5 each), event dispatch (4), dialogue corpus (3).

## Exact next action

**Unit 11 — fix the palette defect that makes the demo's own text invisible.**

`framebuf.BlitTile` computes `row[px] = v + o.PaletteBase`
(`internal/graphics/framebuf/framebuf.go:128`). HUD font ink values are 1–3.
Both scenes pass `white = 3` / `gray = 2`
(`internal/game/scenes/battle.go:127`, `boot.go:86`). `GrayPalette()` defines
only indices 1–3 and leaves 4–255 `{0,0,0}`, and `cmd/ff6demo` passes `nil` for
the palette, so `GrayPalette` is what ships.

Therefore **every "white" string in the running demo resolves to black on
black**, and "gray" strings draw one of their three ink levels. The frame-hash
goldens structurally cannot see it: `Sum256` hashes palette indices and is
deliberately palette-independent.

Fix the convention so drawn ink lands on defined entries, add the assertion
that was missing, regenerate `internal/game/scenes/testdata/boot-frames.json`,
and record a deviation for the fact that the palette is still project-authored
rather than ROM-sourced (retired by readiness F2).

Alongside it: the audit check that asserts the readiness summary against the
counted rows, so error class (1) above cannot recur silently.

## Recommended next command

`/implement-discovery` is the wrong shape — this is a defect fix, not a
discovery. Work it directly, then `/run-quality-gates` and `/checkpoint`.
