# Latest Checkpoint

**[2026-08-02 — Units 15-16, a negative result and the missing implementations](2026-08-02-units15-16-negative-result-and-tooling.md)**
(preceding: [Units 10-14, the compression premise tested and refuted](2026-08-02-units10-14-content-parity-session.md))

State: branch `demo/whelk-content-parity`, nine commits ahead of `main` at
`297ba88`, worktree clean, not pushed. No emulator running, no background
processes.

**Unit 15 is a negative result.** The map header that selects the tile blocks
EXP-0050 located does not exist in the form the search assumed. Four
meaningfully different encodings were tried — 3-byte LE pointers in the upper
window, a structural scan for aligned pointer runs, 2-byte offsets with an
implied bank requiring Narshe and mines to co-occur, and the lower HiROM mirror
in both byte orders. None found a table. The one relaxed-scan "hit" is noise:
452 such runs exist by chance and it holds neither Narshe block. **The loader
computes these addresses**, and the next instrument is the routine, not a
search.

**EXP-0050 carried a claim that was never measured, and now carries the
correction.** It said `0x208460` and `0x223000` are "shared with the mines
interior". They are not — the mines uses `0x20DFA0` and `0x23C0C0`, verified on
both mines savestates. The claim came from running `state ppu` on the mines
state and `state origin` on the Narshe exterior and conflating the two. What is
actually shared is milestone 02 with 04, the same town, which is a weak
discriminator rather than the strong one Unit 15 was planned around.

Two measured by-products: three blocks are common to all three field scenes and
so are **not** map tilesets; and `0x0487C0` + 2048 = `0x048FC0`, exactly the
end of the HUD font block, so the field loads its **last 128 tiles** to
`VRAM:$B800` while the battle loads its start to `VRAM:$AFF0`. T1's load path
now has both destinations.

**Unit 16 put implementations under doctrine that had none.** `/trace-dma`, the
`dma-tracer` skill, the `dma-researcher` agent and `TRACE_DMA.md` had all
existed while nothing read a DMA register. `internal/platform/snesdma` is the
half testable without an emulator, and it encodes two hardware rules that
change conclusions: a raw size of **zero means 65536** — `mesenstate` had this
wrong and `ff6lab state ppu` printed `0` for a full 64 KB VRAM upload — and the
16-bit source address wraps **within its bank**, so `SourceSpan` refuses
bank-crossing transfers rather than pointing a provenance search at a region
the hardware never touched. `mesen/probes/dma-trace.lua` captures registers on
`$420B`/`$420C` writes and is marked **UNEXERCISED** everywhere it is
referenced; writing a probe is not running one.

Five commands that did not exist were created with playbooks —
`recover-compression`, `recover-map`, `recover-text`, `recover-event-opcode`,
`capture-frame` — and three that existed only as `reconstruct-*` are aliased so
both spellings resolve. The playbooks carry this session's lessons:
`RECOVER_COMPRESSION` step 1 is "establish that the data is compressed at all".

Gates: gofmt clean; build/vet/test on both variants including two new fuzz
targets; `ff6lab audit` clean; census clean; restricted scan clean.

Exact next action: **test EXP-0051's "block boundaries are wrong" alternative
— it is nearly free and needs no instrument.** Probe the ROM immediately before
`0x208460` and `0x20DFA0` to find where each block really starts, then re-run
the pointer search against the true starts. `state origin` anchors runs at the
image start and cannot see past it, so this alternative has never been tested
and would invalidate all four of Unit 15's searches if it holds.

Fallback: **read OAM from milestone 04** — every savestate carries it,
`ff6lab state oam` reaches it, and no experiment has ever looked. Needing the
operator: **run `probe dma-trace` over a map load.**
