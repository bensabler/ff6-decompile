# Checkpoint 2026-07-30 — EXP-0024/AUD-0001: audio vertical proof (Unit 5)

## Current question
None open in this unit. Next: final project-state synchronization
(operator plan item 6).

## State
Headless Mesen live (probes EXP-0023/0024 armed, plus an eval-armed
DSP dumper). Evidence frozen at `local_artifacts/experiments/
EXP-0024/`. All five operator rebalance units complete.

## Work completed
EXP-0024 + AUD-0001: battle confirm SFX vertical — press/no-press
delta trials isolate one port write (`$21`→`$2140` from
`ROMCPU:$C117CC`, press+2 frames) and DSP voice 7 (rel 76–79); DSP
snapshot mid-SFX gives SRCN=$05 → directory `ARAM:$1B00` → sample at
`ARAM:$48D8` (2 BRR blocks, 18 bytes); the SFX pack `ARAM:$4800-$491F`
is byte-identical to `ROMFILE:0x051EC9-0x051FE8`. Go: `brr.Decode`
(filters 0–3, clamp+15-bit clip; synthetic tests + fuzz) and
`ff6lab brr info`; captured click decodes cleanly. M3 complete; M7 in
progress.

## Tests and quality gates
Run at commit: gofmt clean, build/vet pass, `go test ./...`, audit
clean.

## Git status
`main`, 10 ahead of origin after this commit. Not pushed.

## Blockers
None hard. Soft: SPC driver dispatch untraced; `$21` command semantics
unknown; MP savestate; GUI/testrunner parity.

## Exact next action
Final synchronization: sweep dashboards/indexes for stale statements,
re-run all gates from a tracked-only clean checkout (this session
touched CI-relevant sources: cmd, internal packages, bridge, probes),
write the closing checkpoint, and commit. Then stop — the operator's
five-unit plan is fulfilled.
