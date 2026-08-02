# EXP-0051: Are the map tile block addresses stored as pointers in the ROM?

## Question

EXP-0050 located the ROM spans that supply a field map's BG tile data. The next
question for readiness F1 is what *selects* them: a map header, a tileset
record, a pointer table.

**Are the observed block addresses present in the ROM as pointers at all?** If
they are, the table that holds them is the map header's target and F1 has a
structure to decode. If they are not, the loader computes them, and the next
instrument is a disassembler or a trace rather than a search.

## Starting state

- **Emulator:** none. Static analysis of preserved savestates and the ROM.
- **ROM identity:** SHA-256
  `0f51b4fca41b7fd509e4b8f9d543151f68efa5e97b08493e4b2a0c06f5d8d5e2`,
  unchanged.
- **States:** SCN-0001 milestones 02, 04 and 05 (`run1`), plus
  `EXP-0035/recon-mines-inside.mss`, all previously frozen and hashed.

## Independent variable

The pointer encoding searched for.

## Controlled variables

One ROM revision. The same block addresses throughout, taken from
`ff6lab state origin` output. Every search run over the whole 3 MiB image.

## Instrumentation

`ff6lab state origin` (EXP-0050) for the block addresses; direct byte searches
over the ROM image for the encodings.

## Expected outcomes

- If a tileset pointer table exists, at least one encoding places two or more of
  the observed addresses at a fixed stride in one region.
- If the addresses are computed, no encoding finds them together.

## Falsifying outcome

Any encoding placing the Narshe blocks and the mines blocks in one table.

## Evidence paths

`local_artifacts/scenarios/SCN-0001/{02-narshe-entry,04-free-movement,05-mines-entry}/run1-*.mss`
and `local_artifacts/experiments/EXP-0035/recon-mines-inside.mss`, unmodified.
Reproduce the block sets with:

```bash
ff6lab state origin <file.mss> vram
```

## Observations

### The block sets, measured

| Scene | VRAM `$0000` | second | third |
|---|---|---|---|
| Narshe entry (02) | `0x208460` (8192) | `0x223000` (8192) | `0x224F00` (4096) |
| Narshe exterior (04) | `0x208460` (8192) | `0x223000` (8192) | `0x224F00` (4096) |
| Mines interior (05) | `0x20DFA0` (8192) | `0x23C0C0` (4096) | — |

Milestones 02 and 04 carry **identical** tile blocks. They are the same town,
so this is expected and is a weak discriminator, not a strong one.

Three further blocks appear in **all three** field scenes and are therefore not
map tilesets:

| VRAM | ROMFILE | Length |
|---|---|---|
| `$5800-$5BFF` | `0x226700` | 1024 |
| `$B800-$BFFF` | `0x0487C0` | 2048 |
| `$E3FF-$FFFF` | `0x182FFF` | ~7169 |

### Four encodings searched, none found

| # | Encoding | Result |
|---|---|---|
| 1 | 3-byte LE CPU pointer, upper window `$C0-$FF` | Narshe blocks: **no occurrence at all**. Mines blocks: 3 and 2 scattered single hits, in unrelated contexts |
| 2 | Structural scan: runs of ≥8 consecutive 3-byte pointers into banks `$E0-$E4`, 32-byte aligned, stride 3 | **zero runs anywhere in the ROM** |
| 3 | 2-byte offsets with an implied bank, Narshe and mines within 256 bytes of each other | **zero co-occurrences** (40 and 65 occurrences respectively, never near) |
| 4 | Lower HiROM mirror `$40-$7F`, and big-endian order, for every block | scattered singles only; no two of our addresses in one region |

A relaxed structural scan — any run of ≥16 consecutive byte triples whose high
byte is a plausible bank, no alignment constraint — returns 452 runs across the
ROM, of which exactly one contains two of our addresses (`0x27166E`, entries 7
and 63). It holds **neither Narshe block**, and with ~50 % of random triples
passing the "plausible bank" test, runs of that length arise by chance
throughout the image. It is noise.

### A verified by-product

`0x0487C0` + 2048 = `0x048FC0`, which is **exactly the end of the HUD font
block** ROM-0016 (`0x047FB0-0x048FBF`). The field scenes load the block's last
**128 tiles** to `VRAM:$B800`.

The battle scene loads the same block's start: EXP-0050 found
`VRAM:$AFE8 ← 0x047FA8`.

## Interpretation

**The observed block addresses are not stored as plain pointers in the ROM**,
under any of the four encodings tested. The loader computes them.

That is a useful negative. It rules out the cheapest structure — a pointer
table a map header indexes — and redirects the next attempt from *searching for
data* to *reading the routine that produces it*.

The three shared blocks are a second useful result: they are common assets
loaded on every field map, so they are **not** what a map header selects, and
any future search for the header's output should exclude them. One of the three
is now identified.

## Alternatives

- **The encoding is one not tested.** Live. Addresses could be stored divided,
  biased, as bank + 16-bit split across two tables, as an index into a base, or
  packed with other fields. Four encodings is not the space.
- **The block boundaries are wrong.** Live, and the more likely of the two.
  `ff6lab state origin` anchors a run where the *image* begins, so a ROM block
  that starts earlier than `VRAM:$0000`'s source would be reported from its
  midpoint, and the true block start — which is what a pointer would hold —
  would never be searched for. Nothing here rules that out.
- **The blocks are assembled from smaller records** whose individual starts are
  pointed to, with the 8192-byte contiguity an artifact of ROM layout. Encoding
  1 would have found those starts and did not, which weakens but does not kill
  this.

## Result

**Negative result, recorded as such.** No pointer table selects these blocks in
any encoding tested. F1 does not gain a structure from this experiment.

**Confirmed by-products:** the per-scene block sets; that milestones 02 and 04
share a complete tileset while 05 does not; that three blocks are common to all
three field scenes; and that the HUD font block's last 128 tiles load to
`VRAM:$B800` on the field.

## Correction to EXP-0050

EXP-0050 states that `0x208460` and `0x223000` are "**shared with the mines
interior**". **That is false.** The mines interior uses `0x20DFA0` and
`0x23C0C0`, on both mines savestates.

The claim was written without running the trace on a mines state: `state ppu`
had been run on `recon-mines-inside.mss`, and `state origin` on the Narshe
exterior, and the two were conflated. It is the same failure this session
documented three times over — an unverified assertion recorded as if measured —
and it is corrected in EXP-0050 rather than quietly dropped.

What survives is the part that was measured: the field tile blocks are
uncompressed and contiguous. The pressure ranking that followed from EXP-0050
does not depend on the shared-blocks claim.

## Confidence

- Per-scene block sets: **Confirmed** (byte-exact, reproducible, two independent
  mines states agree).
- No pointer table under encodings 1-4: **Confirmed as a negative for those
  encodings.** Not "the addresses are computed" — that is the interpretation,
  and it is a **Strong hypothesis**.
- HUD font tail at `VRAM:$B800`: **Confirmed** (exact arithmetic against
  ROM-0016).
- `0x182FFF` as field character sprites: **Unknown.** It is a shared block near
  the OAM tile region and nothing more has been established.

## Stopping condition

Reached, by the project's three-attempts rule: four meaningfully different
searches for the same object failed. Continuing to guess encodings is the wrong
instrument.

## Next action

**Find the routine, not the pointer.** Two options, in order of cost:

1. **Static.** Ghidra is bootstrapped (`local_artifacts/static-analysis/`).
   Disassemble outward from the VRAM upload path and read how the source
   address is formed. `/correlate-static-runtime` exists for exactly this.
2. **Runtime.** A DMA trace would name the source address and the routine that
   set it up in one observation. That instrument did not exist when this
   experiment ran; `mesen/probes/dma-trace.lua` was written immediately
   afterward and is **unexercised** — running it is the fastest path to closing
   this question, and it needs one operator session.

Also worth doing first because it is nearly free: re-run `state origin` with
backward extension past the image start disabled, or probe the ROM immediately
before `0x208460` and `0x20DFA0`, to test the "block boundaries are wrong"
alternative.
