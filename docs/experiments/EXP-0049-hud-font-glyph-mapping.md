# EXP-0049: How does the fixed-width text encoding index the battle HUD font block?

- **Status:** completed
- **Date:** 2026-08-02
- **Domain:** GFX / MENU
- **Demo requirement unblocked:** DEMO-0001 readiness **T2** (glyph →
  character mapping), which gates **T3** (text drawing) and **B17** (battle
  HUD).

## Question

`internal/game/textenc` maps encoding bytes to characters. GFX-0001 records a
257-tile 2bpp font block at `ROMFILE:0x047FB0-0x048FBF`. Nothing connected the
two: CEN-GFX-0001 lists "glyph-to-character mapping" as an unknown field, so a
renderer had no way to turn a byte into a tile.

**Which tile does encoding byte `b` select?**

## Starting state

No emulator run. This experiment is **static and fully reproducible from
tracked evidence** already preserved by EXP-0023.

## ROM identity

SHA-256 `0f51b4fca41b7fd509e4b8f9d543151f68efa5e97b08493e4b2a0c06f5d8d5e2`,
3,145,728 bytes, headerless.

## Emulator identity

Not applicable — no emulator was run. The runtime evidence consumed here was
captured by EXP-0023 (Mesen 2.1.1, `mesen/out/exp10-battle.mss`).

## Battle configuration

Not applicable — no battle was operated. The EXP-0023 capture this analysis
reads was taken during a battle, but no timing-sensitive claim is made.

## Independent variable

The candidate mapping function from encoding byte to VRAM tile index.

## Controlled variables

The ROM revision, the EXP-0023 evidence files, and the `textenc` byte→character
table as it stood before this experiment (letters, digits, `-`, spaces).

## Instrumentation

None injected. Two tracked evidence files were re-read:

- `local_artifacts/experiments/EXP-0023/rom_046FC0_8192.hex` — 8192 ROM bytes
  from the font anchor, i.e. dump tile index `n` is VRAM tile `n`.
- `local_artifacts/experiments/EXP-0023/exp23-tilemap.hex` — 2048 bytes, the
  32×32 BG3 tilemap, little-endian 16-bit entries (bits 0-9 tile index).

Decoding used `ff6lab tiles decode2bpp` and the tracked `tile2bpp` decoder.

## Expected outcomes

EXP-0023 recorded that the HUD tilemap references 37 distinct tiles, **all in
`$180-$1FF`**. `textenc`'s mapped bytes span `$80-$FF`. If the encoding indexes
the block directly, the only affine map carrying `$80-$FF` onto `$180-$1FF` is:

```text
VRAM tile = $100 + b
```

Under that relation the tilemap should decode to **coherent English game
text**. Under any other relation it should decode to noise.

## Falsifying outcome

Any of:

- a byte whose tile does not render the character `textenc` names;
- a tilemap that decodes to non-text under the relation;
- referenced tiles outside `$180-$1FF`;
- a glyph requiring a different offset than `$100`.

## Evidence paths

- `local_artifacts/experiments/EXP-0023/rom_046FC0_8192.hex`
- `local_artifacts/experiments/EXP-0023/exp23-tilemap.hex`
- `local_artifacts/experiments/EXP-0023/hashes.sha256` (freezes both)
- Generated output: `data/graphics/hud-font-glyphs.json` (tracked; hashes only)

## Trials

### Trial 1 — direct glyph decode (8 bytes)

| byte | `textenc` says | VRAM tile | dump tile | decoded shape |
|---|---|---|---|---|
| `$80` | A | `$180` | 384 | A |
| `$81` | B | `$181` | 385 | B |
| `$99` | Z | `$199` | 409 | Z |
| `$9A` | a | `$19A` | 410 | a |
| `$B3` | z | `$1B3` | 435 | z |
| `$B4` | 0 | `$1B4` | 436 | 0 |
| `$BD` | 9 | `$1BD` | 445 | 9 |
| `$C4` | - | `$1C4` | 452 | - |

8 of 8 agree.

### Trial 2 — whole-tilemap decode (the discriminating trial)

Every tilemap entry was decoded through `b = tile − $100`. Result, with `[XX]`
marking bytes `textenc` had no value for:

```text
 1 [EE][EE]Were-Rat    [BF][BF][BF][BF][BF]     0[F9][F1][F0][F0][F0][FA][EE]
 3 [EE][EE]Repo Man    WEDGE     0[F9][F1][F0][F0][F0][FA][EE]
 5 [EE][EE]            VICKS    29[F9][F8][F8][F8][F8][FA][EE]
 9 [EE][EE]  MagiTek   [BF][BF][BF][BF][BF]     0[F9][F1][F0][F0][F0][FA][EE]
11 [EE][EE]            WEDGE     0[F9][F1][F0][F0][F0][FA][EE]
13 [EE][EE]            VICKS    29[F9][F8][F8][F8][F8][FA][EE]
15 [EE][EE]  Item                       [EE]
17 [EE]  Row  Tek   [BF][BF][BF][BF][BF]     0[F9][F1][F0][F0][F0][FA][EE]
19 [EE]             WEDGE     0[F9][F1][F0][F0][F0][FA][EE]
21 [EE][EE]            VICKS    29[F9][F8][F8][F8][F8][FA][EE]
23 [EE][EE]  Item                       [EE]
```

All 37 referenced tiles fall in `$182-$1FF`. The text is coherent FF6 battle
HUD content: two monster names, two party names, four command names, and HP
digits, laid out in the column structure a battle HUD has.

### Trial 3 — whole-block classification (256 bytes)

Each byte's tile was decoded and tested for all-zero pixels:

| Class | Count | Bytes |
|---|---|---|
| Character (named by `textenc`) | 64 | `$80-$BD`, `$BF`, `$C4` |
| Blank (all pixels index 0) | 140 | `$00-$7F`, `$D0`, `$D1`, `$EB-$EF`, `$FB-$FF` |
| Structural (gauge pieces) | 5 | `$F0`, `$F1`, `$F8`, `$F9`, `$FA` |
| Unidentified (non-blank, unnamed) | 47 | `$BE`, `$C0-$C3`, `$C5-$CF`, `$D2-$EA`, `$F2-$F7` |

## Observations

1. Every one of the 8 directly decoded glyphs renders the character `textenc`
   names, at `VRAM tile = $100 + b`.
2. The 32×32 HUD tilemap decodes, under that relation and nothing else, to
   readable game text.
3. All 37 tilemap-referenced tiles lie in `$182-$1FF`, inside the block.
4. Bytes `$00-$7F` decode to entirely blank tiles — the block's lower half
   carries no glyphs.
5. `$BF` renders a question-mark glyph and occupies **all five cells** of the
   third party member's name slot, in the same capture where `WEDGE` and
   `VICKS` render normally in that column.
6. `$EE` is an all-zero tile used 324 times as this layer's background, and is
   distinct from `$FE`/`$FF`, which are also blank.
7. `$F9`/`$FA` render rounded left/right end caps; `$F0` renders an empty
   two-rule segment, `$F8` a filled one, `$F1` a partially filled one. They
   appear only in a fixed six-cell run immediately after the HP digits, and the
   run is fully filled on the row reading `29` and nearly empty on the rows
   reading `0`.

## Interpretation

The fixed-width encoding **is** a direct index into the HUD font block:

```text
VRAM tile       = $100 + b
ROMFILE offset  = 0x047FC0 + 16*b
sheet index     = b + 1        (the extracted sheet starts at VRAM $FF)
```

Observation 5 identifies `$BF` as `?`. Observation 7 identifies `$F0`, `$F1`,
`$F8`, `$F9`, `$FA` as the pieces of the HP/ATB gauge — a **role**, established
from position and fill correlation, not a character.

Observation 4 partially answers CEN-GFX-0001's open "tiles `$00-$FE` compose
region" question for the half that is in this block: VRAM `$100-$17F` is blank.

## Alternatives

- **A different offset than `$100`.** Refuted: any other offset shifts every
  glyph and destroys the text in Trial 2.
- **A lookup table rather than an affine relation.** Not refuted in general,
  but unnecessary: a single affine relation accounts for all 37 referenced
  tiles and all 8 direct decodes with no exceptions. A table would be
  indistinguishable *here* while making a strictly larger claim.
- **`$BF` is some other symbol that merely resembles `?`.** Weakly held. The
  shape and the five-cell name-slot context agree, but this project has not
  observed the naming screen, so "the unnamed protagonist's placeholder" is the
  interpretation, and `?` is the glyph.
- **The gauge tiles are characters this analysis failed to recognise.**
  Refuted by position: they occur only in the fixed six-cell run after the HP
  digits, never in a name or command field.

## Result

**The relation is established.** The falsifier did not fire: no referenced tile
fell outside the block, no decoded glyph disagreed with `textenc`, and the
tilemap decoded to text rather than noise.

Secondary results: `$BF` = `?` added to `textenc`; five gauge tiles identified
as structural; the blank set enumerated; 47 bytes recorded as non-blank and
**deliberately unidentified**.

## Confidence

- Relation `VRAM tile = $100 + b` — **Confirmed**. Two independent lines
  (direct glyph decode; whole-tilemap coherence over 37 tiles), no exceptions.
- The 64 named characters — **Confirmed**. `textenc`'s prior derivation
  (EXP-0026/0027) is now corroborated against *rendered output* rather than
  only against a menu tilemap.
- `$BF` = `?` — **Confirmed** as the glyph.
- Gauge role of `$F0`/`$F1`/`$F8`/`$F9`/`$FA` — **Confirmed** as a role;
  their exact fill semantics are **Unknown**.
- Blank classification of the 140 blank bytes — **Confirmed** by decode.
- The 47 unidentified bytes — **Unknown**. Shapes are visible in the local
  sheet, but a shape without a rendered context is not a character
  identification, and none is recorded.

## Stopping condition

Met: the question was "which tile does byte `b` select", and it is answered
with a falsifiable relation verified two ways. Identifying the remaining 47
glyphs is a separate unit needing captures that render them.

## Next action

Consume the relation: `internal/game/hudfont` carries the model and
`data/graphics/hud-font-glyphs.json` the generated table. DEMO-0001 Unit 4
builds the runtime font loader and text drawer on it.

Registered, not investigated (census): the gauge tiles are the battle HUD's
HP/ATB bar and are directly relevant to readiness row B17.
