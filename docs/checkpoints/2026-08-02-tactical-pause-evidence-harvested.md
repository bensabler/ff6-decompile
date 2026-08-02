# Checkpoint 2026-08-02 — Tactical pause, evidence harvested and hashed

## Current question

None. This is a deliberate pause at an atomic boundary, not a stopping point
forced by a blocker. **No unit is in flight.** The next session performs a
project-wide workflow and orchestration audit.

## State

Branch `demo/whelk-content-parity`, **HEAD `38a62a1`**, fifteen commits ahead of
`main` at `297ba88`. Worktree clean. Pushed through `8591aef`; `38a62a1` is
local-only at the moment this checkpoint is written.

No emulator running. No background processes. No resident instrumentation. No
Mesen session was opened at any point this session.

## Confirmed before this session

The merged DEMO-0001 foundation (units 0–9), ROM identity
`0f51b4fca41b7fd509e4b8f9d543151f68efa5e97b08493e4b2a0c06f5d8d5e2` unchanged.

## Work completed in this pause

**Evidence harvested and hashed.** Four experiments from this session had
records but no preserved evidence directories, and their derived artifacts were
sitting in a session-temporary scratchpad that would have vanished. All four now
follow the project's per-experiment convention:

| Directory | Contents | Files |
|---|---|---|
| `EXP-0050/` | `origin-*.txt`, `ppu-*.txt` for five milestones | 11 |
| `EXP-0051/` | `search-pointer-encodings.py` (replays all four negative encodings) | 2 |
| `EXP-0052/` | three WRAM images, `records-33byte.txt`, `wram-discriminator.py` | 6 |
| `EXP-0053/` | `oam-04.bin`, `vram-04.bin`, `sprites-04.txt`, `ppu-04.txt` | 5 |

Each carries a `README.md` naming its source savestate, its regeneration
command, and its expected headline numbers, plus `hashes.sha256`, verified.
All gitignored — every file is ROM-derived.

**Instrumentation preserved.** The analysis that produced EXP-0051's negative
result and EXP-0052's finding was ad-hoc and would not have survived. Both are
now scripts in their evidence directories, and **both were replayed and
reproduce their recorded numbers exactly**: 0 pointer-table runs, 0
co-occurrences, 19 198 discriminating WRAM bytes, the `$39`/`$17` identifiers,
and the 33-byte lattice.

The durable instrumentation is already tracked: `ff6lab state origin|ppu|sprites`,
`internal/romorigin`, `internal/platform/snesoam`, `internal/platform/snesdma`,
and `mesen/probes/dma-trace.lua`.

**One record strengthened, from the replay.** EXP-0052 claimed a 33-byte
lattice at `0x2D9173` + 20 and + 22 records but had verified `record[28]` only
on the two records that matched the ROM directly. Milestone 02's match begins
one byte earlier because the preceding byte coincides, so its start was
inferred from arithmetic alone. Checked during the harvest: at `0x2D9173`,
`record[28]` is `$39` — the same map id that capture's WRAM holds. The
inference is now confirmed by the field it predicts.

That strengthens the interpretation without changing it, and makes the
entrance/warp reading the leading alternative: all three records carry a map id
at offset 28, and the two naming the same map (`$39`, milestones 02 and 04) are
*different* records. The conclusion stands and is still negative — this table
is not the map header, because a structure indexed by map id would give 02 and
04 the same record.

## Map-descriptor investigation — consolidated status

**Confirmed observations**

- Field BG tile blocks, byte-exact: Narshe `ROMFILE:0x208460` (8192 B),
  `0x223000` (8192 B), `0x224F00` (4096 B); mines `0x20DFA0` (8192 B),
  `0x23C0C0` (4096 B). Uncompressed.
- Milestones 02 and 04 load an identical tileset; milestone 05 does not.
- Three blocks are common to all three field scenes and are therefore not map
  tilesets: `0x226700`, `0x0487C0`, `0x182FFF`.
- `0x0487C0` + 2048 = `0x048FC0`, exactly the end of HUD font block ROM-0016.
- A map identifier: `$39` (Narshe) / `$17` (mines) at `WRAM:+$1305`, `+$13E2`,
  `+$1F80`.
- A 33-byte record table in bank `$ED`, copied verbatim to `WRAM:+$0520`, with
  three captures on one lattice and the map id at `record[28]` on all three.

**Hypotheses, with confidence**

- The 33-byte table is entrance/warp records — **leading alternative**,
  untested.
- `$39`/`$17` are *map* ids specifically rather than area or tileset ids —
  **Tentative**; milestones 02 and 04 agree on both, so the readings are not
  separated.
- The 128-byte runs from bank `$E6` are animated tiles — **Tentative**.

**Alternatives still live**

- The pointer encoding space is not exhausted (four tried).
- **The block boundaries may be wrong.** `ff6lab state origin` anchors a matched
  run where the *image* begins, so a ROM block starting before `VRAM:$0000`'s
  source is reported from its midpoint and its true start has never been
  searched for. This would invalidate all four encodings at once and remains
  untested.
- The header may not be copied into WRAM verbatim at all, but decomposed into
  fields on load — which would make the verbatim method structurally unable to
  find it.

**Unresolved questions**

What selects the tile blocks; what indexes the 33-byte table; the tilemap/layout
source; the sprite sheet's ROM row stride; whether any of the unmatched VRAM is
compressed.

**Stopping rationale**

Six meaningfully different approaches across EXP-0051 and EXP-0052 failed. The
instrument is wrong: this is a question for a disassembler or a trace.

## Last raw observation

Replay of `wram-discriminator.py` against the harvested images: 19 198
discriminating bytes, `$39`/`$17` at three WRAM addresses, verbatim ROM copies
at `WRAM:+$0520` on the 33-byte lattice — all matching the recorded values.

## Active emulator state

None.

## Breakpoints/watchers

None.

## Evidence paths

`local_artifacts/experiments/EXP-0050/`, `EXP-0051/`, `EXP-0052/`, `EXP-0053/`,
each with `README.md` and verified `hashes.sha256`. Sources are the frozen
SCN-0001 milestone savestates, unmodified.

## Files changed

`docs/experiments/EXP-0052-map-id-and-33-byte-record.md` only. No code changed
in this pause.

## Tests and quality gates

gofmt clean; `go build`, `go vet`, `go test` pass on both build variants;
`ff6lab audit` clean (eleven checks); `census validate` clean; `archive verify`
9/9; restricted-file scan clean; `local_artifacts/` correctly untracked.

## Git status

Branch `demo/whelk-content-parity`, HEAD `38a62a1`, worktree clean, one commit
ahead of `origin`.

## Unresolved decisions

Whether to promote the WRAM-discrimination script to a tracked `ff6lab`
subcommand. It produced a real finding and is currently preserved only as an
untracked script. Deferred deliberately — it is implementation work, and this
is a pause.

## Blockers

None blocking. The map-descriptor line is stopped by choice under the
three-attempts rule, not by an external dependency.

## Exact next action

**Project-wide workflow and orchestration audit**, as directed. No research or
implementation unit should begin before it.

When content work resumes, in priority order:

1. **Test the block-boundary alternative.** Nearly free, needs no instrument,
   and would invalidate all four pointer encodings if it holds. Probe the ROM
   immediately before `0x208460` and `0x20DFA0` for the true block starts, then
   re-run the search against those.
2. **Diff the sprite tile region between milestones 02 and 04.** Both are the
   Narshe exterior; if the player faces differently, the walking-frame set and
   the sheet's row stride fall out with no new capture (F6).
3. Needing one operator session: **run `probe dma-trace` over a map load**. It
   would name the source address and the routine that set it up in a single
   observation, closing the map-descriptor question. The probe is written and
   **unexercised**.

## Recommended next command

`/audit-project`, then `/orchestrate`.
