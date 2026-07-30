# Checkpoint 2026-07-30 — Unit 10 complete (formula frame verified)

## Current question
Question #21. Unit 11 / EXP-0010: dump and decode the bracketed formula
body `ROMCPU:$C20B83`–`$C20C10`.

## State
Supervised daytime session. Mesen live with bridge + injected watches
(code in EXP records). No Go changes since Unit 7; gates green there.

## Work completed
EXP-0009: `JSR $0B83` confirmed at `$C23469` (dump `rom_C23430_128.hex`,
`a83c86d3…`); formula body bracketed to `$C20B83`–`$C20C2C`; ten-slot
target-loop context recorded (Tentative); `+$3018` promoted (per-slot
bit table, Strong).

## Exact next action
EXP-0010: `read cpu C20B83 141`; hand-decode toward the DP `$F0`
computation; promote findings; implement in Go only if a complete,
verified arithmetic unit emerges.

## Git status
`main`, 13 commits ahead of origin.
