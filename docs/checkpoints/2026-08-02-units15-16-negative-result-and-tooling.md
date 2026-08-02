# Checkpoint 2026-08-02 — Units 15-16, a negative result and the missing implementations

## Current question

None open. Two questions were closed — one with a negative result, one with
code — and the next action needs one operator session.

## State

Branch `demo/whelk-content-parity`, nine commits ahead of `main` at `297ba88`.
Worktree clean, not pushed. No emulator running, no background processes.

## Work completed

| Unit | Commit | Result |
|---|---|---|
| 15 | `2673065` | **EXP-0051** — no pointer table selects the map tile blocks; EXP-0050 corrected |
| 16 | `7969d50` | Implementations put under the DMA doctrine; five missing commands created |

### Unit 15 — a negative result, and a correction to my own record

The unit set out to find the map header that selects the tile blocks EXP-0050
located. **It does not exist in the form the search assumed.** Four
meaningfully different encodings were tried:

1. 3-byte LE CPU pointers, upper window — the Narshe blocks do not occur in the
   ROM at all; the mines blocks occur 3 and 2 times in unrelated contexts.
2. Structural scan for runs of ≥8 consecutive 3-byte pointers into banks
   `$E0-$E4`, 32-byte aligned — **zero runs anywhere in the image**.
3. 2-byte offsets with an implied bank, requiring a Narshe and a mines offset
   within 256 bytes — zero co-occurrences out of 40 and 65 occurrences.
4. Lower HiROM mirror `$40-$7F` and big-endian order — scattered singles only.

A relaxed scan returns 452 plausible-looking runs, of which one holds two mines
addresses and neither Narshe address; with ~50 % of random triples passing the
bank test, that is noise and is recorded as noise.

**The loader computes these addresses.** Recorded per the three-attempts
stopping rule rather than by guessing further encodings. Two live alternatives
are named rather than buried: the encoding space is not exhausted, and — more
likely — the block boundaries may be wrong, because `state origin` anchors a
run where the *image* begins, so a ROM block starting before `VRAM:$0000`'s
source would be reported from its midpoint and its true start never searched
for.

**EXP-0050 carried a claim that was never measured.** It stated that
`0x208460` and `0x223000` are "shared with the mines interior". That is false —
the mines uses `0x20DFA0` and `0x23C0C0`, verified on both mines savestates.
The claim came from running `state ppu` on the mines state and `state origin`
on the Narshe exterior and conflating the two. EXP-0050 now carries the
correction inline. What is actually shared is milestone 02 with milestone 04,
the same town, which is a weak discriminator rather than the strong one Unit 15
was planned around.

Two measured by-products: three blocks are common to all three field scenes and
are therefore **not** map tilesets (`0x226700`, `0x0487C0`, `0x182FFF`); and
`0x0487C0` + 2048 = `0x048FC0`, exactly the end of the HUD font block ROM-0016,
so the field loads that block's **last 128 tiles** to `VRAM:$B800` while the
battle loads its start to `VRAM:$AFF0`. T1's load path now has both
destinations; only the copying routine is missing.

### Unit 16 — implementations under the doctrine

`/trace-dma`, the `dma-tracer` skill, the `dma-researcher` agent and
`TRACE_DMA.md` had all existed for a long time with **nothing reading a DMA
register**, in Lua or in Go.

`internal/platform/snesdma` is the half testable without an emulator. Two
hardware rules are encoded because both are easy to get wrong and both change
conclusions:

- **A raw size of zero means 65536.** The hardware decrements before testing.
  `mesenstate` had this wrong, and `ff6lab state ppu` printed `0` for what is a
  full 64 KB VRAM upload. It now prints 65536.
- **The source address is 16-bit and the bank does not increment**, so a
  transfer past `$FFFF` wraps within the bank. Those bytes are not contiguous,
  so `SourceSpan` refuses them rather than reporting a range into the next
  bank — a wrong span sends a provenance search somewhere the hardware never
  touched.

`mesen/probes/dma-trace.lua` captures registers on writes to `$420B`/`$420C`
and logs raw bytes rather than a decoded summary, so a log can be re-decoded
later. It is marked **UNEXERCISED** in its own header, in `TRACE_DMA.md`, in
the capability matrix and in the census. Writing a probe is not running one.

Five commands that did not exist were created with playbooks: `recover-compression`,
`recover-map`, `recover-text`, `recover-event-opcode`, `capture-frame`. Three
more existed only as `reconstruct-*` while their playbooks were named
`RECOVER_*`; both spellings now resolve. The playbooks carry this session's
lessons rather than generic procedure — `RECOVER_COMPRESSION` step 1 is
"establish that the data is compressed at all", `RECOVER_MAP` requires two maps
so a field can be told from a constant.

## Last raw observation

`ff6lab state ppu local_artifacts/experiments/EXP-0035/recon-mines-inside.mss`
— channel 0, `$2118`, source `$7EC180`, **65536** bytes (previously printed as
`0`).

## Active emulator state

None.

## Evidence paths

No new evidence captured. Unit 15 read the frozen SCN-0001 savestates and the
ROM; Unit 16 wrote code and doctrine.

## Tests and quality gates

`gofmt -l .` clean. Build, vet and test pass on both build variants, including
two new fuzz targets in `snesdma`. `ff6lab audit` clean. `census validate`
clean. Restricted scan clean.

## Git status

Branch `demo/whelk-content-parity`, nine commits ahead of `main`, worktree
clean, **not pushed**.

## Blockers

**One, and it needs the operator.** Both remaining paths to the map header run
through instruments that need a live session or a disassembler:

- `mesen/probes/dma-trace.lua` is written and unexercised. One session running
  it over a map load would name the source address and the routine that set it
  up in a single observation, closing EXP-0051's open question.
- Ghidra is bootstrapped and needs no session, but does need a bounded
  static-analysis unit.

## Exact next action

**Test EXP-0051's "block boundaries are wrong" alternative first — it is
nearly free and needs neither instrument.** Probe the ROM immediately before
`0x208460` and `0x20DFA0` to find where each block really starts, then re-run
the pointer search against the true starts. `state origin` cannot see past the
image start, so this alternative has never been tested and would invalidate all
four of Unit 15's searches if it holds.

If that fails: **read OAM from milestone 04** (CEN-GFX-0008). Every savestate
carries `ppu.oamRam`, `ff6lab state oam` reaches it, no experiment has ever
looked, and readiness F6 depends on it.

Needing the operator: **run `probe dma-trace` over a map load.**

## Recommended next command

Direct static work, then `/run-quality-gates` and `/checkpoint`. Use
`/recover-map` once the header is located.
