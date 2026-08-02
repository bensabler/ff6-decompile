# DEMO-0001 acceptance criteria

- **Milestone:** [DEMO-0001](DEMO-0001-new-game-to-whelk.md)
- **Updated:** 2026-08-02 (Unit 0 — program start)
- **Overall:** **NOT PASSING** — 0 of 17 progression steps, 0 of 6 gates

DEMO-0001 is complete only when every required criterion passes. This file is
the scorecard. A step flips to PASS only when it has a recorded evidence path;
"it looked right" is not a pass.

## Gates

| Gate | Requirement | Status | Evidence |
|---|---|---|---|
| G1 | A documented command launches a standalone Go application | **FAIL** | `cmd/ff6demo` does not exist |
| G2 | The acceptance run requires neither Mesen nor Ghidra running | **FAIL** | no run exists |
| G3 | The run begins at New Game without loading a savestate or injecting internal state | **FAIL** | no run exists |
| G4 | Assets are generated locally from the operator's verified ROM by a documented workflow | **PARTIAL** | `ff6lab extract all` + `archive verify` 8/8 OK; `dialogue`/`maps`/`animations`/`scripts` categories empty |
| G5 | No prohibited files are committed | **PASS** | CI restricted-extension scan + `ff6lab audit` `CheckRestrictedTracked` green |
| G6 | All deviations from the original are recorded, not hidden by weakened tests | **PASS** (vacuously) | [DEMO-0001-DEVIATIONS.md](DEMO-0001-DEVIATIONS.md) — 1 open row |

## Progression steps

The continuous acceptance run, in order. Each step needs an evidence path under
`local_artifacts/` plus an automated assertion.

| # | Step | Beat | Status | Blocked by (readiness row) |
|---|---|---|---|---|
| 1 | Launches successfully | — | **FAIL** | E1 |
| 2 | Begins at New Game | B01 | **FAIL** | F14 |
| 3 | Presents the opening sequence | B02 | **FAIL** | F1, F5, T4 |
| 4 | Enters Narshe | B06 | **FAIL** | F1, F7 |
| 5 | Executes required scripted movement and events | B03, B05 | **FAIL** | F8, F4 |
| 6 | Presents all required dialogue | B05 | **FAIL** | F10, F11 |
| 7 | Enters and exits required scripted battles | B07 | **FAIL** | F13, B17, B18 |
| 8 | Allows field movement where the original does | B08 | **FAIL** | F3, F4 |
| 9 | Traverses required Narshe and mines areas | B09, B11 | **FAIL** | F1, F3 |
| 10 | Handles required map transitions | B10 | **FAIL** | F7 |
| 11 | Triggers required random or scripted encounters | B12, B13 | **FAIL** | F12 |
| 12 | Reaches Whelk | B17 | **FAIL** | steps 9–11 |
| 13 | Runs the Whelk introduction | B17 | **FAIL** | F8, F10 |
| 14 | Supports a functional Whelk battle | B18 | **FAIL** | B13, B15 |
| 15 | Whelk can be defeated using valid game mechanics | B18 | **FAIL** | B13, B10, B11 |
| 16 | Exits the battle correctly | B18 | **FAIL** | B18 |
| 17 | Reaches the agreed post-Whelk checkpoint | B19 | **FAIL** | step 16 |

## Required validation evidence

| Kind | Requirement | Status |
|---|---|---|
| Unit tests | parsers and deterministic logic | partial — existing packages covered; engine not yet written |
| Integration tests | subsystem boundaries | **none** |
| Scenario tests | route progression under deterministic input scripts | **none** for the demo; `internal/scenario/route` covers the *emulator* route |
| Framebuffer comparison | demo output vs Mesen at named milestones | **none** |
| Dialogue/event transcript comparison | demo vs recorded event-flag timeline | **none** — `data/scenarios/opening-event-flags.json` is the target |
| Asset hashes and provenance | every runtime asset traceable to a ROM span | **PASS** for the 8 extracted assets |
| Audio cue order and timing | demo cue sequence vs captured | **none** |
| Battle-state assertions | HP/MP/ATB/queue vs Mesen | **none** |
| Full acceptance evidence | complete New Game → Whelk run | **none** |

## Policy on scaffolding

Temporary scaffolding may exist only when it is isolated behind a clear
interface, marked provisional in code, recorded in
[DEMO-0001-DEVIATIONS.md](DEMO-0001-DEVIATIONS.md) with an exact replacement
requirement, and **not used to claim parity**.

Specifically prohibited by this milestone:

- fake predetermined battle outcomes;
- the route hardcoded as one bespoke cinematic;
- arbitrary timers substituted for verified battle logic;
- manually cropped screenshots used as runtime assets;
- tests weakened to hide a difference.
