# EXP-0052: What in WRAM distinguishes one field map from another?

## Question

EXP-0051 failed to find a pointer table selecting the tile blocks. This asks
the same question from the other side: **what does WRAM hold that differs
between two field maps and is the same within one?** A value that behaves that
way is a map identifier or is derived from one, and it is what a header table
would be indexed by.

## Starting state

- **Emulator:** none. Static analysis of preserved savestates and the ROM.
- **ROM identity:** `0f51b4fca41b7fd509e4b8f9d543151f68efa5e97b08493e4b2a0c06f5d8d5e2`.
- **States:** SCN-0001 milestones 02 and 04 (Narshe, same tileset per EXP-0050)
  and 05 (mines interior), all previously frozen and hashed.

## Independent variable

The map the capture was taken on.

## Controlled variables

One ROM revision. Three captures from one golden route. The discriminator is
**02 == 04 and 04 != 05**, which holds the two Narshe states fixed against the
mines state rather than comparing two arbitrary maps.

## Instrumentation

`ff6lab state wram`, and byte searches over the ROM.

## Expected outcomes

- A map identifier appears as a small value repeated at a few WRAM addresses.
- If a header record is copied into WRAM, some differing WRAM run appears in
  the ROM verbatim, and the two maps' copies sit at a fixed stride.

## Falsifying outcome

For the identifier: no small value satisfies `02 == 04 != 05`. For the copied
record: no differing WRAM run is found in the ROM for both maps.

## Evidence paths

`local_artifacts/scenarios/SCN-0001/{02-narshe-entry,04-free-movement,05-mines-entry}/run1-*.mss`,
unmodified, covered by each directory's `hashes.sha256`.

## Observations

19 198 bytes satisfy `02 == 04 != 05`, in 2 037 runs. Two results survive that
noise.

### A map identifier

Three separate WRAM addresses hold the same pair of values:

| Address | Narshe (02 and 04) | Mines (05) |
|---|---|---|
| `WRAM:+$1305` | `$39` | `$17` |
| `WRAM:+$13E2` | `$39` | `$17` |
| `WRAM:+$1F80` | `$39` | `$17` |

Milestones 02 and 04 agree, which is consistent with EXP-0050's finding that
they load an identical tileset.

### A 33-byte record copied from ROM

Scanning every map-dependent WRAM window in the low 8 KB for a verbatim ROM
match found exactly one region, `WRAM:+$0520`, and its ROM sources lie on a
single lattice:

| Capture | ROM source | Offset from `0x2D9173` |
|---|---|---|
| 02 narshe-entry | `0x2D9173` | 0 records |
| 04 free-movement | `0x2D9407` | **20 records** |
| 05 mines-entry | `0x2D9449` | **22 records** |

660 = 20 × 33 and 726 = 22 × 33 exactly. **Stride 33.**

And every record carries the identifier at the same offset. `record[28]` is:

| Capture | Record start | `record[28]` | WRAM id |
|---|---|---|---|
| 02 narshe-entry | `0x2D9173` (inferred from the lattice) | `$39` | `$39` |
| 04 free-movement | `0x2D9407` (matched directly) | `$39` | `$39` |
| 05 mines-entry | `0x2D9449` (matched directly) | `$17` | `$17` |

Milestone 02's record start is **inferred**: its verbatim match begins one byte
earlier, at `0x2D9172`, because the preceding byte happens to agree. The
lattice fixes the true start at `0x2D9173`, and `record[28]` there is `$39` —
so the inference is independently confirmed by the field it predicts, rather
than resting on the arithmetic alone.

(Verified during the Unit 19 evidence harvest; the original revision of this
record checked only the two directly-matched records.)

## Interpretation

`WRAM:+$0520` receives a 33-byte record copied verbatim from a table in bank
`$ED`, and that record contains the value that also appears at `+$1305`,
`+$13E2` and `+$1F80`.

**This is not the map header.** Milestones 02 and 04 share a map identifier and
a tileset but load *different* records, twenty apart, **and both records carry
`$39` at offset 28**. A structure indexed by map id would give them the same
record. Whatever the table is indexed by, it is not the identifier.

Two distinct records referencing one map is what strengthens the
entrance/warp reading below: a map has several entrances, and each would need
its own record naming the map it leads to.

The most that is established is: a 33-byte static record, selected per capture
by something not yet known, referencing a map by the same id WRAM holds.

## Alternatives

- **The table is entrance/warp records** — "go to map `$39` at position X,Y" —
  which would explain two records for one map, since a map has several
  entrances. This is now the **leading** alternative: all three records carry a
  map id at offset 28, and the two that name the same map are different
  records. Still untested — nothing here checks the other 32 bytes against a
  position, and no third map has been examined.
- **`$39`/`$17` are not map ids** but area or tileset ids, and the identical
  value across 02 and 04 reflects a shared tileset rather than a shared map.
  EXP-0050 cannot separate these, because 02 and 04 agree on both.
- **The record is NPC or object state.** Less likely given it is a verbatim ROM
  copy of a fixed size, but not excluded.
- **The header is not copied into WRAM verbatim at all**, but decomposed into
  fields on load. This would explain why the search found only one region, and
  it would make the verbatim method structurally unable to find the header.

## Result

**Partial, and negative on the primary question.** A map identifier is located
with three independent addresses. A 33-byte ROM record table is located and its
stride confirmed by three captures on one lattice. **Neither selects the tile
blocks**, and the header remains unlocated after six meaningfully different
approaches across EXP-0051 and this record.

## Confidence

- `$39`/`$17` differ per map and repeat at three addresses: **Confirmed**.
- Those values are *map* identifiers specifically: **Tentative hypothesis** —
  02 and 04 share a tileset as well as a value, so the two readings are not
  separated.
- The 33-byte record, its lattice and `record[28]`: **Confirmed** (three
  captures, exact arithmetic).
- What the table is indexed by, and what the record means: **Unknown**.

## Stopping condition

Reached. Six approaches across two experiments have not found the selector, and
the instrument is wrong: this is a question for a disassembler or a trace.

## Next action

Read the loader routine. Ghidra is bootstrapped and needs no session;
`mesen/probes/dma-trace.lua` would answer it in one observation but has never
been run.

Cheapest unstarted alternative, still untouched: **read OAM from milestone 04**
(CEN-GFX-0008), which readiness F6 depends on.
