# Checkpoint 2026-08-01 — EXP-0037: opening event-flag inventory (Unit 39)

## Current question
None open in the unit. EXP-0037 reached its stopping condition; the
next unit is EXP-0038 (mines traversal to milestone 06).

## State
The opening's event-flag behavior is inventoried and deterministic.
All writes to the three verified arrays (`WRAM:+$1E80/+$1EA0/+$1EC0`)
across the scheduled route (power-on → milestone `05-mines-entry`) are
captured, attributed, and reproduced: **20 flags touched, 162
value-changing writes, identical across one GUI pass and two headless
runs** at frame+address+value+PC granularity. The two evidence runs
are byte-identical on every channel (48 065-line write log, 62
snapshots, final WRAM), and the final WRAM equals the milestone-05
hash `c26453d3…` — five byte-identical runs of that milestone now
exist, and the instrumentation provably perturbs nothing.

## Work completed
- **EXP-0037** (record: `docs/experiments/EXP-0037-...md`): probe
  `mesen/probes/EXP-0037.lua` (EXP-0036 route unchanged + 96-byte
  write watch with shadow integrity check, JSONL timeline, array
  snapshots; guarded by the extended `probe_sync_test.go`).
  11 latched story flags (all via the script set-handlers), 4
  transient, 5 engine working bits; boot clear at frame 21; zero flag
  changes during battles; `$1EA5`'s `$00→$01→$05→$0D` reproduced as
  `EVF-1EA0-$28/$2A/$2B` (EXP-0036 cross-check passed).
- **Static decode of every live writer PC**, which found the
  **16-handler script-command family** (`$C0B593-$C0B6D2`, eight
  bases — five beyond the verified three), the flag-test family, the
  boot clear loops, and the **event-interpreter anchor
  `ROMCPU:$C09B5C`** (every handler tail; CEN-EVENT-0001 promoted to
  CANDIDATE_LOCATION).
- **DISC-0008** + implementation `internal/game/eventflags`
  (decoder/masks/Ref/FlagAt, golden-vector tests).
- **Tracked inventory** `data/scenarios/opening-event-flags.json`
  (per-flag events with era/leg/battle context, run provenance,
  evidence hashes).
- ROM regions ROM-0027..ROM-0032; census updates (CEN-EVENT-0008
  NORMAL_PATH_VERIFIED with 20 found records, CEN-EVENT-0001 anchor,
  CEN-QUIRK-0001 testrunner uninit-read addresses registered).
- **GUI/testrunner input parity verified for this schedule**
  (BLOCKERS updated); run1's clean-exit `.srm` shown hash-identical
  to the backed-up virgin image.
- Dashboards, indexes, scenario record + manifest (B10/B16/B19),
  RESEARCH_QUEUE, ACTIVITY_LOG synchronized.

## What remains uncertain
- Flag **meanings** (deliberately unassigned; identifiers stay
  address-anchored).
- Whether the boot-burst flags (frames 2 516–2 528) are
  input-triggered or timed (constant schedule cannot discriminate).
- SRAM backing of the arrays; runtime use of the statically-decoded
  extra bases (`$1EE0-$1F5F`, `$1DC9`); the `$1E40-$1E6F` region
  cleared at boot; the `$0200/$0205` source bit mirrored into
  `EVF-1EC0-$FF`; flag-test callers.

## Active instrumentation and evidence
No Mesen running; no stale processes; Saves dir empty (virgin SRAM).
Evidence under `local_artifacts/experiments/EXP-0037/` (three runs ×
flags.jsonl / snapshots.log / events.log / wram-final.bin +
testrunner stdout + `hashes.sha256`). Milestone artifacts untouched.

## Tests and quality gates
gofmt clean; build/vet clean; `go test ./...` green (14 packages with
`FF6_ROM`, ROM-dependent tests skip cleanly without); `ff6lab audit`
clean; `census validate` clean; `archive verify` 8/8 clean;
restricted-extension scan clean; probe-sync test now guards both
EXP-0036 and EXP-0037 encodings.

## Git status
main; one coherent unit committed and pushed (this checkpoint's
commit).

## Exact next action
**EXP-0038 — golden route segment 5: mines traversal to milestone
`06-random-encounter`.** Extend the state-driven controller from
milestone 05 into the mines along EXP-0035's recon notes, reach the
first random encounter reproducibly, and capture the encounter
trigger context (`+$11E0` producer lead, CEN-WORLD-0006) without
decoding the encounter system itself. Create the experiment record
first; run headless with the established lab controls.
