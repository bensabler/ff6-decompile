# Checkpoint 2026-07-30 — Unit 9 complete (formula lead refuted, refined)

## Current question
Question #21 (damage formula). Unit 10: verify the refined `$C23469`
JSR frame (EXP-0009 dump `ROMCPU:$C23430`–`$C234B0`).

## State
Supervised daytime session. Mesen live with bridge + injected watches
(pfCount, dseen, cseen, eseen, qseen, dlog). No Go changes since Unit 7.

## Work completed
EXP-0008: `$C26B14` JSR refuted by dump (`rom_C26AE0_128.hex`,
`9feb9ecb…`); stack model refined (pushed-PS byte) → new tentative frame
`JSR` at `~$C23469`.

## Exact next action
EXP-0009: `read cpu C23430 128`; verify a JSR occupying
`$C23469`–`$C2346B`; decode surroundings for DP `$F0` staging. Fallback:
live full-stack + DP snapshot at `$C20C76`.

## Git status
`main`, 9 commits ahead of origin. Gates last run green at Unit 7
(no Go changes since).
