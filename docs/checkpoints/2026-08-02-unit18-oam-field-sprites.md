# Checkpoint 2026-08-02 — Unit 18, OAM read for the first time

## Current question

None open. The next action is a single bounded capture-and-diff.

## State

Branch `demo/whelk-content-parity`, thirteen commits ahead of `main`, pushed.
Worktree clean, no emulator running, no background processes.

## Work completed

**EXP-0053 read `ppu.oamRam`, which no experiment had ever done** despite every
preserved savestate carrying it and the Mesen bridge having `read oam` wired
since it was written. Readiness F6 had no evidence of any kind; it now has a
composition, a ROM address, and a named next step.

Milestone 04, Narshe exterior under player control: object mode **3**
(16×16 / 32×32), name base `VRAM:$C000`, **5 of 128 sprites on screen** — the
other 123 parked at Y=240, which is why the on-screen filter exists.

- **Sprite 53** is a single 16×16 at screen centre, palette 2, tile `$000`.
  Its four tiles decode to a recognisable humanoid figure. Tiles at
  **`ROMFILE:0x150180`** — the project's first ROM location for any character
  graphic.
- **Sprites 54/55 and 101/102** are a mirrored 32×32 object: each pair shares
  one tile name with the right-hand sprite H-flipped. Tiles around `0x1841A0`
  in bank `$D8`, the same region EXP-0050 found battle sprite runs in.

**74 % of `VRAM:$C000-$FFFF` is verbatim ROM**, so field sprite graphics are
largely uncompressed — the same shape as the map tile blocks.

The tile-to-ROM map is **not linear**: `$000`→`0x150180` and `$001`→`0x1501A0`
are contiguous, but `$010`→`0x150240` is `0xC0` further rather than `0x40`. The
row stride is recorded Unknown, because two rows is not a pattern.

### A correction, recorded because it nearly became a finding

The first pass probed 128-byte blocks and reported the player's tiles absent
from the ROM — which would have made field character sprites the strongest
compression candidate on the route. Wrong: a 16×16 sprite's tiles are
`n, n+1, n+16, n+17`, so `$000`-`$001` and `$010`-`$011` sit 14 tiles apart in
VRAM and a 128-byte probe spans the gap. Re-probed per 8×8 tile, all four are
present.

**The probe geometry has to match the data's geometry, and a negative from the
wrong geometry looks exactly like a negative from the right one.**

## Last raw observation

`ff6lab state sprites .../04-free-movement/run1-04.mss` — 5 on-screen sprites,
mode 3, base `VRAM:$C000`.

## Active emulator state

None.

## Evidence paths

No new evidence captured. Read the frozen milestone-04 savestate and the ROM.

## Tests and quality gates

gofmt clean; build, vet and test pass on both build variants, including a new
fuzz target in `snesoam`; `ff6lab audit` clean; `census validate` clean;
restricted scan clean.

The readiness-count check fired again during this unit and caught the summary
drifting when F6 moved — the third time a check written this session has caught
an error before commit.

## Git status

Branch `demo/whelk-content-parity`, thirteen commits ahead of `main`, pushed,
worktree clean.

## Blockers

None for the next unit.

## Exact next action

**Capture a second field state with the player facing a different direction,
then diff OAM and the sprite tile region.** That names the walking-frame set
and the sheet's row stride in one comparison, which is exactly what F6 needs
next, and it converts "the player's sprite is at `0x150180`" into an
extractable sheet.

This needs one operator session, because the corpus has only one player-
controlled field state.

Offline alternatives if no session is available:

- **Diff the sprite tile region between milestones 02 and 04.** Both are the
  Narshe exterior; if the player faces differently between them, the frame set
  falls out with no new capture.
- **Run `probe dma-trace` over a map load** — still the fastest path to the map
  header, and still unexercised.

## Recommended next command

`/capture-frame` for the second field state, then `/reconstruct-sprite`.
