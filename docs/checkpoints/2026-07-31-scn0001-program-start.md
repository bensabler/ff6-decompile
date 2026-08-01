# Checkpoint 2026-07-31 — SCN-0001 program start (infrastructure + scenario spec)

## Current question
SCN-0001: reconstruct New Game → Whelk victory as a vertical slice
(operator directive, 2026-07-31). Active unit: EXP-0031 (golden route
segment 1: power-on → New Game).

## State
Two units committed this session:
1. Infrastructure (`2cddafc`): expanded CLAUDE.md constitution
   (restored the dropped legal-boundary section) + the
   `ff6-content-census` skill (fixed a stray post-frontmatter rule).
2. SCN-0001 master scenario record + machine manifest; dashboards
   synchronized (CURRENT_FOCUS, RESEARCH_QUEUE).

Scenario baseline: 19 beats defined (B01–B19), all PARTIAL; gaps and
evidence links recorded per beat in the scenario record. No power-on
route exists; event/map/encounter-roll systems unlocated; nothing
observed from Whelk introduction onward.

## Work completed
- Constitution + census skill validated, defects corrected, committed,
  **pushed** (operator explicitly requested push).
- `docs/scenarios/SCN-0001-opening-to-whelk.md` +
  `data/scenarios/opening-to-whelk.json` authored: boundary, golden
  route plan (11 milestones), Whelk branches A/B/C, beat matrices
  seeded from EXP-0001..0030 evidence, honest gap list.

## Tests and quality gates
gofmt/build/vet/test (10 packages)/ff6lab audit/census validate —
clean at both commits.

## Git status
main, in sync with origin (infra pushed; scenario unit committed
after this checkpoint, push per operator standing request this
session).

## Active instrumentation and evidence
None running. Archived states: mesen/out/checkpoint1/2/3-mines,
exp10-battle (.mss, local only). Bridge: mesen/bridge.lua; probes
mesen/probes/EXP-*.lua. Headless-only constraint per BLOCKERS.md.

## Exact next action
EXP-0031: write the experiment record (falsifiable, bounded), then a
frame-scheduled probe that boots the ROM from power-on with no
savestate, navigates title → New Game via scheduled presses, saves
milestone state `00-new-game` under
local_artifacts/scenarios/SCN-0001/00-new-game/, dumps the input
transcript + key WRAM assertions, and repeats the run to verify
byte-identical assertions (determinism).
