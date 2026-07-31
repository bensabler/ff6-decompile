# Checkpoint 2026-07-31 — EXP-0027: spell database extracted

## Current question
None open in this unit. Per the prioritization rules, three
consecutive magic-adjacent units have run (EXP-0025/26/27) — the next
selection rotates domains after a coverage review.

## State
Lab processes shut down. Evidence frozen at
`local_artifacts/experiments/EXP-0027/` (17 files, hashes.sha256).
Census: 46 entries; first EXTRACTED_COMPLETE entry (spell name
table). ROM ledger: 21 regions (spell names + esper names known).

## Work completed
EXP-0027: name-table boundary Confirmed at exactly 54 entries; esper
name table discovered and format-decoded (27×8 at 0x26F6E1;
stride-~10 ability names follow, candidate); record byte 5 = MP cost
Confirmed behaviorally (field Cure cast: '5 MP Needed' gate, heal
34→77/77, deduction 24→19); field character-record block located by
single-byte WRAM diffs (HP +$1609, MP +$160D); defeat flow and
formation names (Were-Rat, Repo Man) registered; all 54 spell
records bulk-extracted locally, ids+numbers mirrored into
data/census/spells.json.

## Tests and quality gates
Run at commit: gofmt clean, build/vet pass, `go test ./...`
(10 packages), `ff6lab audit` clean (census checks included).

## Git status
`main`, 1 ahead of origin after this commit. Not pushed (push only on
request).

## Blockers
None hard. Soft items unchanged.

## Exact next action
Domain rotation per `ff6lab coverage gaps`: recommended — **monster
stat-record source trace** (write-watch the battle-init population of
the +$3B18/+$3B2C/+$3B68 per-slot tables from a fresh battle entry;
the ROM source it reveals is the monster database, unlocking the
MONSTER domain the way EXP-0026 unlocked MAGIC). Alternates: HUD font
load-path trace (GFX), SPC dispatch trace (AUDIO).
