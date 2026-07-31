# DISC-0003: Element-response transform

## Status
Confirmed; implemented and tested.

## Supporting experiments
EXP-0010 (records under `docs/experiments/`).

## Discovery
Per-target element masks in two family arrays (`WRAM:+$3BCC` flip|zero, `+$3BE0` double|halve) modify the pending amount first-match-wins after a battle-wide nullify test (`~+$3EC8 & +$11A1`): flip-to-heal, zero, halve, double (with $8000 guard). Byte-exact at `ROMCPU:$C20BD3-$C20C1D`.

## Go implementation
`battle.ElementResponse, battle.ApplyElementResponse`

## Tests
`TestApplyElementResponse`

## Confidence and residue
Behavior Confirmed byte-exact per the cited experiments; semantic labels
beyond the verified behavior (stat/element names, MP identity) remain at
their recorded hypothesis levels in `docs/sessions/` and the hypotheses
dashboard.
