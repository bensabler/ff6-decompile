# Checkpoint 2026-08-01 — CONTRA-0002 resolved: `+$1EA5` is an event-flag byte (Unit 37)

## Current question
None open. The contradiction is resolved and dependent work is
unfrozen. Next: EXP-0037, mapping the opening's event flags.

## State
`WRAM:+$1EA5` had two competing readings in canonical records —
"candidate map id" (EXP-0035) and "map-load target / event-state"
(EXP-0036). **Both are refuted.** The byte is **byte 5 of an
event-flag bit array based at `WRAM:+$1EA0`**. The resolution needed
no emulator run: a static decode of the writer settled it.

The contradictions dashboard had also been asserting "none open" while
CEN-WORLD-0007 carried an unresolved contradiction field; that gap is
fixed and the dashboard now carries a note to check census-embedded
contradictions at each resolution.

## Work completed
Record: `docs/contradictions/CONTRA-0002-1ea5-map-id-vs-event-flags.md`.

- **Froze** dependent work: `internal/scenario/route`
  (`AddrMapContextCandidate`, `UntilMapChange`/`MapValue`,
  `State.MapContext`) and the probe's `MAPC` logging. Nothing on the
  shipped route depended on the semantics, so no route result was at
  risk.
- **Discriminator (static, no run):** scanned the ROM for every
  instruction referencing `$1EA5` — **none exists**; a standalone
  variable would be addressed directly. Decoded the live writer
  instead: the probe reports the PC after the store, so `$C0B5B6`
  corresponds to `STA $1EA0,Y` at `ROMCPU:$C0B5B3`, preceded by
  `LDA $1EA0,Y` / `ORA $C0BAFC,X`.
- **Confirmed the flag system:** set masks `$C0BAFC` =
  `01 02 04 08 10 20 40 80`; clear masks `$C0BB04` =
  `FE FD FB F7 EF DF BF 7F` (exact complements); index decoder
  `$BAED` = `REP #$20 / TAX / LSR ×3 / TAY / TDC / SEP #$20 / TXA /
  AND #$07 / TAX / RTS` → Y = flag/8, X = flag&7. Three parallel
  arrays at `+$1E80`, `+$1EA0`, `+$1EC0`, `$20` bytes apart.
- **The decisive observation:** the recorded values *accumulate bits*
  (`$00` → `$01` → `$05` → `$0D` = bits 0, then 0+2, then 0+2+3). A
  map id does not gain a bit as the story advances; a flag byte does.
- **Explained the tension EXP-0035 could not:** milestone 02 reads
  `$00` because no flags in that byte were set yet. The location
  correlation was incidental — flags accumulate with story progress,
  which tracks location in a linear opening.
- **Unfroze** by renaming the identifiers that encoded the refuted
  claim: `AddrMapContextCandidate` → `AddrEventFlagsBase` (`0x1EA0`),
  `UntilMapChange`/`MapValue`/`MapContext` →
  `UntilWatchedByte`/`WatchValue`/`WatchByte` (deliberately generic —
  the condition now carries no semantics of its own). Probe renamed
  `MAPC` → `EVFLAG`; the tag change is noted in the probe header so
  pre-resolution logs stay interpretable.
- Supersession links written into EXP-0035 and EXP-0036 (both keep
  their original readings visible, marked superseded). CEN-WORLD-0007
  stripped of the map-id claim; **CEN-EVENT-0008** registers the flag
  system. Dashboards and indexes updated.

## What remains uncertain
- Which array holds which flag family (event / treasure / map-state).
- The meaning of individual flag numbers, including `$28`/`$2A`/`$2B`
  observed set at milestone 05.
- Whether the arrays are SRAM-backed.
- **CEN-WORLD-0004 (map header / tileset load path) is now genuinely
  unstarted** — `+$1EA5` was a false lead and its removal leaves that
  question with no active thread.

## Tests and quality gates
gofmt clean; build/vet clean; `go test ./...` green including the
route package's 11 test groups and the probe-sync guard; `ff6lab
audit` clean; `census validate` clean.

## Git status
main; unit committed and pushed.

## Active instrumentation and evidence
No emulator running. Resolution evidence is the ROM decode recorded
inline in CONTRA-0002 (mask tables, decoder, access-site counts) plus
the live write log in `local_artifacts/experiments/EXP-0036/`.

## Exact next action
EXP-0037 (write the record first): arm write-watches on all three flag
arrays (`WRAM:+$1E80`, `+$1EA0`, `+$1EC0`, `$20` bytes each) across the
existing scheduled route, decoding each write back to a flag number
via the byte index and the bit that changed. Produce an ordered list of
**which flags the opening sets, at which frame, and at which route
beat**. Bound it to the route already established (power-on →
milestone 05); do not chase individual flag meanings in the same unit.
That inventory is what B16 and B19 need, and it gives CEN-EVENT-0001
its first concrete anchor.
