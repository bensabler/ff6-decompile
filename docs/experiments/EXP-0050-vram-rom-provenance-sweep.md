# EXP-0050: Which captured VRAM regions are verbatim ROM copies?

## Question

Readiness row X1 asserts that FF6's compression format "gates F1, F6, B14,
B15, B16" — maps, field sprites, battle backgrounds, enemy graphics and party
battle sprites. That claim has never been tested. It was inferred from the
absence of records, not from evidence.

**Which parts of a captured VRAM image appear in the ROM verbatim, and which do
not?** A region that is a verbatim copy needs a slice, not a decompressor, and
is not gated by X1 at all.

## Starting state

- **Emulator:** none. Static analysis of preserved savestates and the ROM.
- **ROM identity:** SHA-256
  `0f51b4fca41b7fd509e4b8f9d543151f68efa5e97b08493e4b2a0c06f5d8d5e2`,
  re-verified this session, unchanged from `docs/research/ROM_IDENTITY.md`.
- **States:** the SCN-0001 golden-route milestone savestates, already frozen
  and hashed under `local_artifacts/scenarios/SCN-0001/`. No new capture was
  taken and no emulator was launched.

## Independent variable

The scene: which milestone savestate the VRAM image comes from.

## Controlled variables

One ROM revision. One search procedure with fixed parameters (32-byte probe,
64-byte minimum reported run, 3 distinct values required in a probe window).
The `run1` state of each milestone, all captured on the same route.

## Instrumentation

`internal/romorigin` and `ff6lab state origin`, written for this experiment.
The procedure walks the captured image, takes a probe window at each position,
searches the ROM for it with `bytes.Index`, extends any match backward and
forward to its full run, and resumes past the run.

`internal/mesenstate` supplies the VRAM image and the PPU configuration
(Unit 12). Nothing here reads the ROM at runtime in the demo; this is a
laboratory tool.

## Expected outcomes

- If FF6 stores map tile graphics compressed, as X1 assumes, a field scene's
  BG tile region will show **no** long verbatim runs.
- If it stores them uncompressed, the tile region will map to a small number
  of long contiguous ROM spans.

## Falsifying outcome

For the "uncompressed" reading: any field scene whose BG1/BG2 tile region
fails to map to contiguous ROM spans. For the "compressed" reading: any long
contiguous span at all.

## Evidence paths

Savestates under `local_artifacts/scenarios/SCN-0001/*/run1-*.mss`, unmodified
and covered by each directory's `hashes.sha256`. Reproduce with:

```bash
ff6lab state origin <file.mss> vram
```

## Observations

Verbatim coverage of the full 64 KB VRAM image, per scene:

| Milestone | Scene | BG mode | Verbatim | Spans |
|---|---|---|---|---|
| 04-free-movement | Narshe exterior | 1 | **51 %** (34 003 B) | 33 |
| 05-mines-entry | Mines interior | 1 | **52 %** (34 155 B) | 37 |
| 02-narshe-entry | Narshe entry | 1 | **47 %** (31 449 B) | 29 |
| 03-first-scripted-battle | Battle | 1 | **38 %** (25 082 B) | 52 |
| 01-opening-cinematic | Opening | **7** | **0 %** (0 B) | 0 |

The field scenes resolve into a consistent shape. For milestone 04, with BG1
and BG2 reading tiles from `VRAM:$0000` and BG3 from `VRAM:$6000`:

| VRAM | ROMFILE | ROMCPU | Length |
|---|---|---|---|
| `$0000-$1FFF` | `0x208460` | `$E08460` | 8192 |
| `$2000-$3FFF` | `0x223000` | `$E23000` | 8192 |
| `$4000-$4FFF` | `0x224F00` | `$E24F00` | 4096 |
| `$5000-$57FF` | banks `$E6` | — | ~18 runs of 128-171 B |
| `$5800-$5BFF` | `0x226700` | `$E26700` | 1024 |

Three long contiguous spans totalling **20 480 bytes** carry the static tile
set. The `$5000-$57FF` region is instead assembled from many short runs, most
of them exactly **128 bytes** — four 4bpp tiles — drawn from a narrow region of
bank `$E6` (`0x261B00-0x2635FE`).

> **Correction (EXP-0051, same day).** This section originally read: "The mines
> interior (milestone 05) shares `0x208460` and `0x223000` and the same
> bank-`$E6` short runs, and differs in its third block." **That is false.** The
> mines interior uses `0x20DFA0` (8192 B) and `0x23C0C0` (4096 B), verified on
> both mines savestates.
>
> The claim was never measured. `state ppu` had been run on the mines state and
> `state origin` on the Narshe exterior, and the two were conflated into an
> assertion. What *is* shared, and was measured, is milestone **02** with
> milestone 04 — the same town, so a weak result rather than a strong one.
>
> Nothing else in this record depends on it: the uncompressed contiguous spans
> are measured, and so is the pressure re-ranking that followed.

In the battle scene, `VRAM:$AFE8-$BFFF` maps to `ROMFILE:0x047FA8` — the HUD
font block, ROM-0016, `0x047FB0-0x048FBF`. Party and enemy sprite regions
(`$6180-$7829`) are short runs from banks `$D8` and `$E9`.

The Mode 7 opening produced **zero** matched bytes across the whole image.

## Interpretation

**The field map tile graphics on this route are stored uncompressed.** Three
contiguous ROM spans account for 20 KB of tile data, and a decompressor is not
involved in placing them.

The 128-byte runs are consistent with **animated tiles**: a small set of tiles
replaced per frame from a table of frames, four tiles at a time. That reading
is not established here — it is an interpretation of the run length and the
clustering, and this experiment did not test it.

The Mode 7 result at 0 % is the corpus's strongest signal that *something*
transforms graphics on the way to VRAM. It is also the case where an
alternative explanation is most available: Mode 7 interleaves tile and tilemap
bytes in VRAM, so an untransformed source would still fail a byte-for-byte
search. **The 0 % does not establish compression.**

Two independent by-products:

1. The HUD font's **destination** is now known: `VRAM:$AFF0` for tile `$FF`,
   which is the load address readiness T1 listed as missing evidence. The
   routine that copies it is still unlocated.
2. `0x208460` and `0x223000` appear in the Narshe exterior **and the Narshe
   entry**, which are the same town. Three *other* blocks — `0x226700`,
   `0x0487C0` and `0x182FFF` — appear in all three field scenes including the
   mines, and EXP-0051 identifies `0x0487C0` as the HUD font block's last 128
   tiles.

## Alternatives

- **The matched spans are coincidence.** Rejected for the long runs: an 8192-byte
  agreement between a 64 KB image and a 3 MB ROM does not occur by chance. It
  remains a live concern for the 128-byte runs, which is why the tool reports a
  minimum run length and why those runs are described rather than claimed.
- **The unmatched remainder is compressed.** Not established. A verbatim search
  is defeated by bit-plane reordering, by runtime composition from smaller
  records, by tilemaps built in WRAM, and by a single differing byte. The
  measured 47-52 % is an upper bound on "copied", not a lower bound on
  "compressed".
- **These scenes are unrepresentative.** Possible. Five scenes on one route
  were tested. Nothing is claimed about world-map graphics, other towns, other
  battle backgrounds, or any scene outside SCN-0001.

## Result

**Confirmed:** the BG tile data for milestones 02, 04 and 05 is present in the
ROM verbatim, in contiguous spans, and the HUD font is present verbatim in the
battle scene. Byte-exact, reproducible from tracked evidence with one command.

**Refuted:** the unqualified claim that a compression format gates *all* map
and battle graphics. For the map tile graphics on the DEMO-0001 route, it gates
nothing.

**Unknown:** what accounts for the unmatched 48-62 %, and whether any of it is
compressed.

## Confidence

- Verbatim spans and their addresses: **Confirmed** (byte-exact, reproducible,
  no emulator needed).
- Shared tile blocks between two field maps: **Confirmed**.
- HUD font destination `VRAM:$AFF0`: **Confirmed**.
- 128-byte runs are animated tiles: **Tentative hypothesis**.
- Anything about the unmatched remainder: **Unknown**.

## Stopping condition

Reached. The question was whether verbatim copies exist and where; both are
answered. Characterising the remainder is a different, larger question and is
deliberately not started here.

## Next action

Re-scope readiness X1 from "compression gates all graphics" to what the
evidence supports, and follow the ROM spans this experiment located toward the
map header that selects them — `ROMFILE:0x208460` and `0x223000` are now
concrete anchors for F1, which previously had none.
