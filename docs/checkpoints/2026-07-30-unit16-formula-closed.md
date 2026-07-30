# Checkpoint 2026-07-30 — Unit 16 complete (base formula closed; overnight wrap)

## Current question
Next frontier is a three-way choice (CURRENT_FOCUS): #28 variance hunt
(live), #27 stat-operand meanings (live), #26 physical-path decode
(dump-only).

## State
Overnight session 2 wrapping at a major natural boundary: **the standard
damage pipeline is decoded byte-exact and numerically closed end to
end** — `base = power×4 + (power×$11AE×$11AF)>>5` → defense
`(×(255−def))/256+1` → halvings → element response → ×1.5 chain →
9999-cap accumulator → delta engine → unified arrays → display copier →
HUD. The one absent ingredient anywhere in the decoded chain: RNG
(question #28).

**Mesen 2.1.1 running** with bridge; injected globals `dlog` + `bseen`
(EXP-0014's code, preserved in its record). Game mid/post-battle in the
mines. Kill with `kill -9` or reuse.

## Work completed (Units 15–16)
EXP-0014 (writer census + numeric closure) and EXP-0015 (formula decode
+ `$C20DD1` helper + store-site verification done programmatically).
Go: `BaseAmountStandard` with the live golden vector and the full-chain
golden test (450 → ~346).

## Tests and quality gates
gofmt clean; build/vet pass; `go test ./...` pass (5 packages; battle
suite now includes `TestDamageChainGolden`). Run at wrap.

## Git status
`main`; after this commit the overnight-2 stack is 3 commits
(fce1e70 + this + none pending). Not pushed (operator reviews).

## Evidence paths
`rom_C22B40_192.hex` (`05f3eae9…`), `rom_C20DD1_32.hex` (`905eaccd…`),
BASE-WRITER lines in `events.log`, `exp9.log`, `exp9-battle.png`.

## Unresolved decisions
None for the operator.

## Blockers
None hard.

## Exact next action
EXP-0016 (#28): from `checkpoint3-mines.mss`, provoke the same attack
twice with different input-frame timing (savestate reload between
trials); diff the BASE-WRITER/QUEUE values. Any divergence localizes
the RNG; identical values push the variance question upstream to
hit/crit rolls.

## Recommended next command
`/resume-session` in a fresh context; continue with EXP-0016 or the
dump-only #26 if the emulator is down.
