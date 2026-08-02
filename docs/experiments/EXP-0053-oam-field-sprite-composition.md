# EXP-0053: What does OAM hold on the Narshe exterior, and where do its tiles come from?

## Question

Readiness F6 (field sprites, walking frames) has had **no evidence of any
kind**. Every preserved savestate carries `ppu.oamRam`, and the Mesen bridge
has had `read oam` wired since it was written, but **no experiment has ever
read it** (CEN-GFX-0008).

**Which sprites does the Narshe exterior actually show, how are they composed,
and do their tiles come from the ROM verbatim?**

## Starting state

- **Emulator:** none. Static analysis of a preserved savestate and the ROM.
- **ROM identity:** `0f51b4fca41b7fd509e4b8f9d543151f68efa5e97b08493e4b2a0c06f5d8d5e2`.
- **State:** SCN-0001 milestone `04-free-movement`, `run1`, the Narshe exterior
  under player control, frozen and hashed.

## Independent variable

None. This is a single-state observation, deliberately: the question is what
one capture contains, and nothing here is compared across scenes.

## Controlled variables

One ROM revision, one capture, one decode path.

## Instrumentation

`internal/platform/snesoam` and `ff6lab state sprites`, written for this
experiment; `internal/mesenstate` for OAM, VRAM and the PPU object mode.

## Expected outcomes

- A small number of on-screen sprites, most of OAM parked off-display.
- If field sprite graphics are uncompressed, their VRAM tiles appear in the
  ROM verbatim.

## Falsifying outcome

For the uncompressed reading: sprite tiles absent from the ROM.

## Evidence paths

`local_artifacts/scenarios/SCN-0001/04-free-movement/run1-04.mss`, unmodified,
covered by that directory's `hashes.sha256`. Reproduce with:

```bash
ff6lab state sprites local_artifacts/scenarios/SCN-0001/04-free-movement/run1-04.mss
```

## Observations

Object mode **3** (16×16 / 32×32), object name base `VRAM:$C000`.
**5 of 128 sprites are on screen**; the other 123 are parked at Y=240.

| # | X | Y | Tile | Size | Pal | Pri | Flip |
|---|---|---|---|---|---|---|---|
| 53 | 120 | 81 | `$000` | 16×16 | 2 | 2 | — |
| 54 | 112 | 87 | `$1AC` | 16×16 | 7 | 2 | — |
| 55 | 128 | 87 | `$1AC` | 16×16 | 7 | 2 | **H** |
| 101 | 112 | 103 | `$1AE` | 16×16 | 7 | 2 | — |
| 102 | 128 | 103 | `$1AE` | 16×16 | 7 | 2 | **H** |

Two compositions are visible in that table.

**Sprite 53 is a single 16×16 at screen centre**, palette 2, tile `$000`.
Decoding its four tiles out of VRAM and rendering them as characters produces a
recognisable humanoid figure — head, torso, arms, feet. It is the player
character.

**Sprites 54/55 and 101/102 are a mirrored 32×32 object.** Each pair shares one
tile name with the right-hand sprite H-flipped, and the pairs sit 16 pixels
apart vertically. Four 16×16 sprites, two tile names, built by mirroring.

### ROM provenance

Per 8×8 tile, comparing VRAM against the ROM:

| Sprite | Tile | VRAM | ROMFILE |
|---|---|---|---|
| player | `$000` | `$C000` | `0x150180` |
| player | `$001` | `$C020` | `0x1501A0` |
| player | `$010` | `$C200` | `0x150240` |
| player | `$011` | `$C220` | `0x150260` |
| object | `$1AD` | `$F5A0` | `0x1841A0` |
| object | `$1BC` | `$F780` | `0x184380` |
| object | `$1BD` | `$F7A0` | `0x1843A0` |
| object | `$1AE` | `$F5C0` | `0x1841C0` |
| object | `$1AF` | `$F5E0` | `0x1841E0` |

Across the whole sprite tile area `VRAM:$C000-$FFFF`, **220 of 296 distinctive
8×8 tiles (74 %) are found in the ROM verbatim.**

The player's tiles are in bank `$D5` around `0x150180`; the mirrored object's
are in bank `$D8` around `0x1841A0`. EXP-0050 found battle sprite runs in the
same bank `$D8` region.

## Interpretation

**Field sprite graphics are largely uncompressed**, like the map tile blocks
EXP-0050 found. 74 % is a floor, not a ceiling: the measure counts only tiles
distinctive enough to search for, and misses any tile the engine assembles.

The player character's field sprite is at **`ROMFILE:0x150180`**, which is the
first ROM location the project has for any character graphic.

The tile-to-ROM mapping is **not linear**. `$000`→`0x150180` and
`$001`→`0x1501A0` are contiguous, but `$010`→`0x150240` is `0xC0` further on
rather than the `0x40` a straight copy would give. The name table is 16 tiles
wide while the ROM sheet is arranged differently, so the loader copies rows.
The row stride is **not** established here — two rows is not a pattern.

## A correction, recorded because it nearly became a finding

The first pass probed 128-byte blocks and reported the player's tiles as **not
present in the ROM**, which would have made field character sprites the
strongest compression candidate on the route. That was wrong. A 16×16 sprite's
tiles are `n, n+1, n+16, n+17`, so `$000`-`$001` and `$010`-`$011` are 14 tiles
apart in VRAM — a 128-byte probe spans the gap and matches nothing. Re-probed
per 8×8 tile, all four are present.

The lesson is the same one this session keeps producing: the probe geometry has
to match the data's geometry, and a negative from the wrong geometry looks
exactly like a negative from the right one.

## Alternatives

- **The mirrored object is an NPC.** Possible, untested. Palette 7, priority 2
  and a symmetric 32×32 build are consistent with a decorative object or a
  large NPC; nothing here distinguishes them.
- **Sprite 53 is not the player.** Unlikely — it is the only sprite on a
  different palette, sits at screen centre, and renders as a humanoid — but
  "the sprite the player controls" was not tested by moving and re-capturing.
- **The unmatched 26 % is compressed.** Not established, and the same caveat as
  EXP-0050 applies: a verbatim search is defeated by runtime assembly and
  bit-plane reordering too.

## Result

**Confirmed:** the OAM composition of milestone 04; object mode 3 with base
`VRAM:$C000`; the player's field sprite tiles at `ROMFILE:0x150180`; the
mirrored 32×32 object's tiles around `0x1841A0`; and that 74 % of the sprite
tile area is verbatim ROM.

**Unknown:** the sprite sheet's ROM layout and row stride, which character the
sheet belongs to, the walking-frame set and animation timing, and what the
mirrored object is.

## Confidence

- OAM composition and object mode: **Confirmed** (decoded, reproducible).
- Per-tile ROM addresses: **Confirmed** (byte-exact).
- Sprite 53 is the player character: **Strong hypothesis**.
- The mirrored object's identity: **Unknown**.
- Sheet layout and row stride: **Unknown** — two rows is not a pattern.

## Stopping condition

Reached. The question was what OAM holds and where the tiles come from; both
are answered. Recovering the walking-frame set is a larger question and is
deliberately not started here.

## Next action

Capture a second field state with the player facing a different direction, then
diff OAM and the sprite tile region. That names the walking-frame set and the
sheet's row stride in one comparison, and it is what F6 needs next.
