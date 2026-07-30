# Checkpoint 2026-07-30 — Unit 12 complete (base-amount routines decoded)

## Current question
Questions #23–25 (base producer at `+$11B0`; helpers `$C2370B`,
`$C247B7`).

## State
Supervised daytime session. Mesen live with bridge + injected watches.
No Go changes this unit (helpers unverified — threshold not met).

## Work completed
EXP-0011: both base-amount routines decoded byte-exact where covered;
family census = 12 arrays (05 master table); HP→MP retarget mechanism
found (`$C20DDD`); questions #22 closed, #23–25 opened.

## Last raw observation
`rom_C20C9E_233.hex` (`c8c895a8…`), `rom_C20D87_128.hex` (`95e1f088…`).

## Tests and quality gates
No Go changes; last green at Unit 11 commit.

## Git status
`main`, 16 commits ahead of origin; committing Unit 12 now.

## Exact next action
Unit 13 / EXP-0012: dump `$C247B0`–`$C24800` (#25 multiplier — expect
SNES ALU registers) and `$C23700`–`$C23760` (#24 final transform —
randomness candidate); both are deterministic dumps and unlock Go
implementation of variant A.

## Recommended next command
Continue autonomous session (Unit 13).
