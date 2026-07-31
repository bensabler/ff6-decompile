# Checkpoint 2026-07-31 — EXP-0030: formation table located (Unit 30)

## Current question
None open. Next: +$11E0 encounter-roll producer, or domain rotation
(GFX/AUDIO) per the prioritization rules.

## State
Lab to be shut down at session close. Evidence at
local_artifacts/experiments/EXP-0030/. Census 48 entries; ROM ledger
26 regions. MONSTER domain: stats + formations + flags tables all
located and cross-verified.

## Work completed
EXP-0030: formation staging writer $C2315C decoded (id x15 ->
$CF6200 reads, 16-byte copy to +$3F44; flags $CF5900 x4). Formation
44 verified = mines encounter = monsters {19,77}; corrected
EXP-0028's coincidental 77/78 claim. +$11E0 = the formation id input
(encounter-roll output).

## Tests and quality gates
Run at commit: gofmt/build/vet/test (10 packages)/audit clean.

## Git status
main, 4 ahead of origin after this commit. Not pushed.

## Exact next action
Either: (a) EXP-0031 — +$11E0 producer write-watch during the
encounter walk (zone data + step state -> the encounter system), or
(b) rotate to GFX (HUD font load path) / AUDIO (SPC dispatch) per
rule 3. A fresh session may pick either from ff6lab coverage gaps.
