# Checkpoint 2026-07-30 — Repository repair and state synchronization

## Current question
None mid-flight. Next: EXP-0020 (question #30, AI/RNG layer) per
RESEARCH_QUEUE; then MP verification; then graphics/audio vertical
proofs (operator brief, Units 2–5).

## State
Maintenance unit complete. Mesen still running with the v1 bridge
resident (v2 written to disk; takes effect on next launch — validate
then). All injected watches from EXP-0014..0018 remain live in the
running instance.

## Work completed
See ACTIVITY_LOG (2026-07-30 maintenance entry) — gitignore collision,
CI, CLAUDE.md, audit tooling, discovery layer, evidence citations,
historical labels, evidence layout, bridge v2, Go fixes, dashboard sync.

## Tests and quality gates
gofmt clean; build/vet pass; `go test ./...` pass (7 packages);
`go build ./cmd/ff6lab` pass; `ff6lab audit` clean.

## Git status
`main`, 5 commits ahead of origin after this commit (Units 17–20 + this).
Not pushed (operator pushes after review).

## Blockers
Soft only — see BLOCKERS.md (bridge v2 unvalidated; Mesen version note;
MP labels pending verification).

## Exact next action
EXP-0020: probe-based watches (bridge v2 `probe` command) on the
`$C23190` caller chain during identical-state trials from
`mesen/out/exp10-battle.mss`; restart Mesen first to load bridge v2 and
validate it; use the new `local_artifacts/experiments/EXP-0020/`
evidence layout.

## Recommended next command
Continue autonomous session (EXP-0020).
