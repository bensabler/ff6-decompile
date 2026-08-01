# Checkpoint 2026-08-01 — EXP-0039: Whelk reached (breadth pass, Unit 41)

## Current question
None open in the unit. EXP-0039 stopped at a stable gameplay boundary
on context limits, with Whelk reached but not yet defeated.

## State
**Whelk has been reached, entered and observed** — the scenario's
furthest point to date. This was a deliberate **breadth-first** unit
(explore widely, register briefly, continue forward), not a depth
investigation, and no format, AI, provenance or RNG decoding was
performed.

## Work completed
- **Route map** from milestone 05 to Whelk with every branch and dead
  end enumerated; the mines↔exterior transition is **bidirectional**
  (south from `(26,1C)` exits to the Narshe exterior at `(1F,18)`;
  re-entry produced no shaft dialogue, corroborating the EXP-0037
  flag inventory).
- **Encounter discriminator:** a third mines encounter (formation 14
  again) fired at `(26,15)`, nine tiles from the `(26,0B)` tile of
  both scheduled runs, after a different step count — **fixed-tile
  triggering is refuted**, resolving an alternative EXP-0038 left
  open. The producer remains unlocated (CEN-WORLD-0006).
- **Magitek ability list captured** (CEN-MAGIC-0001, long-standing
  gap): Fire Beam, Bolt Beam, Ice Beam, Heal Force. Battle commands
  are character-specific (leader has Magic; another member does not).
- **Whelk (B17/B18 first contact):** scripted beat at `(2A,09)`, then
  **contact-triggered** battle from `(2A,07)`; **formation 432
  (`$01B0`)** with its staged record captured; multi-box
  introduction/warning dialogue observed; **shell counterattack
  confirmed behaviorally** (killed one member outright, then wiped
  the party). First attempt **ended in defeat**, capturing the defeat
  flow (CEN-BATTLE-0007, a registered gap).
- New census entries CEN-EVENT-0010 (Whelk beat + invocation) and
  CEN-EVENT-0011 (post-Whelk guard/Esper beat, sequence position
  unresolved); six existing entries updated.
- Scenario record + manifest (B11, B12, B15, B17, B18, B19),
  dashboards and indexes synchronized.

## What remains uncertain
- **Whelk's monster ids are Unknown.** Reading formation 432's record
  bytes 2-7 as ids gives record 0 — the opening guard — which is
  implausible for a boss, so the id field must carry a
  high-bit/extension not yet decoded (FF6 exceeds 256 monsters). New
  bounded question blocking B18/B14.
- Head vs shell target modelling, counter selection and damage, boss
  AI — all unobserved beyond the single counter event.
- The post-Whelk guard/Esper beat's true sequence position
  (CEN-EVENT-0011) — it was reached after a defeat-and-reload.
- Whether the mines zone draws formations other than 14 (three rolls,
  all 14).

## Active instrumentation and evidence
**No background processes running** — verified with `jobs -l` and
`pgrep` after terminating GUI Mesen (it ignores TERM; escalated to
KILL per the documented behavior) and reaping the job. Ephemeral
`.srm` removed so the next run boots virgin. Evidence: 30 artifacts
under `local_artifacts/experiments/EXP-0039/` (recon transcript,
bridge logs, 20+ screenshots, two savestates, `hashes.sha256`) and
`local_artifacts/scenarios/SCN-0001/07-pre-whelk/` (pre-Whelk
savestate + screenshot + hashes).

## Tests and quality gates
gofmt clean; build/vet clean; `go test ./...` green; `ff6lab audit`
clean; census clean; `archive verify` 8/8 clean; restricted-extension
scan clean.

## Git status
main; one coherent unit committed and pushed.

## Exact next action
**EXP-0040 — Whelk victory attempt (branch A).** Reload
`local_artifacts/scenarios/SCN-0001/07-pre-whelk/pre-whelk-recon.mss`
in visible GUI Mesen, use **Heal Force** to restore the party from
26/19/56 *before* engaging, then attack **only while the head is
extended**, verifying head state from a screenshot between actions
rather than mashing A. On victory, capture milestone
`10-whelk-victory` (WRAM + screenshot + savestate) and the first
stable post-battle state (B19). Branch B (deliberate shell attack) is
already partly recorded by this pass's defeat and should be re-run
cleanly afterwards.
