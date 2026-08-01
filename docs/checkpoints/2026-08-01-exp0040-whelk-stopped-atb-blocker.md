# Checkpoint 2026-08-01 — EXP-0040: Whelk attempt stopped on ATB blocker (Unit 42)

## Current question

None open in the unit. EXP-0040 is **stopped, not completed**: Whelk was
attacked correctly but never defeated, and further execution is deferred
behind a methodological blocker.

## State

**Whelk was attacked correctly and was NOT defeated.** Two piloted GUI
attempts from the preserved pre-Whelk state, both under
`Bat.Mode = Wait`. The branch-A premises were **verified** — Whelk is
two battle entities, the head is damageable while extended, the shell
never takes damage from head-targeted attacks, and the head/shell state
is visually classifiable — but the battle could not be **operated**
reliably because the project has no ATB model.

**Milestone `10-whelk-victory` is NOT established. B19 remains
uncaptured.**

Run classification: **aborted due to missing ATB model** (primary),
**partial contextual capture**, **failed tactical attempt**.

## Confirmed before this session

Golden route power-on → milestone `06-random-encounter`; event-flag
system (DISC-0008); unified battle arrays (DISC-0001); formation table;
Whelk reached, entered and lost (EXP-0039).

## Work completed

- **Whelk is two battle slots** — new **CEN-BATTLE-0009**: shell slot 4
  = **50000/50000 HP**, head slot 5 = **1600/1600 HP** (DISC-0001
  arrays); MP candidates 120/120 and 1000/1000.
- **Head-only targeting Confirmed**: six measured hits (162, 168, 171,
  177, 181, 186) all reduced slot 5; slot 4 held 50000 for the whole
  unit. No shell strike occurred, so EXP-0039's counter was not
  re-triggered and nothing new is claimed about it.
- **Head/shell state is visually classifiable** at 4× upscale
  (falsifier 2 not triggered); the white target cursor sits in the same
  region and must not be confused with the head.
- **Field healing route found**: field menu reachable before contact,
  inventory **Tonic ×4 / Potion ×1**; four Tonics took the party from
  **26/19/56** to **76/77, 105/105, 106/107**.
- **MagiTek sets are character-specific; EXP-0039's list was
  incomplete** — Terra **eight** (Fire Beam, Bolt Beam, Ice Beam, Bio
  Blast, Heal Force, Confuser, X-fer, TekMissile), escorts **four**.
  CEN-MAGIC-0001 corrected.
- **Ally target cursor defaults to the caster's slot**; `Cursor=Memory`
  reopens a character's list on their last-used ability.
- **CEN-EVENT-0011 resolved and renamed**: the guard/Esper beat fires at
  `(2A,07)` **before** Whelk on clean, never-defeated runs (reproduced
  twice) — normal progression, not a post-defeat artifact.
- Intro dialogue records that Whelk **eats lightning and stores the
  energy in its shell** and says not to attack the shell; enemy action
  **"Slime"** observed. Contextual only; **no elemental behavior
  tested**.

### Operator errors caught and corrected

1. **Config misread.** Initially reported `Bat.Mode = Active` and a
   switch to Wait. Re-inspection of `34-config-before.png` shows
   "Active" grey / "Wait" white — **Wait was already selected**; the
   hand cursor marks the row, not the selection. The right-press was a
   no-op. **Only `Cursor: Reset → Memory` was actually changed.** Both
   attempts ran under WAIT.
2. **Blind multi-press batches** raced menu/dialogue transitions;
   switched to verify-then-press with a screenshot before every
   confirmation. Artifacts renamed accordingly
   (`33-attempt1-end.png`, `attempt1-end.mss`).

## Last raw observation

Attempt 2 frozen at **frame 182248**: formation `+$11E0` = `B0 01`
(432); party HP **51/77**, **105/105**, **107/107**; Terra MP 24/29;
head **1246/1600**; shell **50000/50000**; head **retracted**; Terra's
**MagiTek submenu open with the cursor on Fire Beam**, no target
selection open. Field position bytes read `00 10` (not field-meaningful
during battle).

## Active emulator state

**None.** GUI Mesen **PID 62332** (parent 61294) was the only background
process; logs were harvested first, then it was terminated with
`kill -9` (it ignores TERM) and confirmed gone. The background task
reaped with exit code 1 as expected. `jobs -l` empty; `pgrep -f
"Mesen|mesen|testrunner|analyze"` clean. **SRAM directory verified still
empty** — virgin boot preserved.

## Breakpoints/watchers

None left installed. The unit used only bridge reads
(`loadstate`/`press`/`read`/`screenshot`/`eval`) plus the bridge's
standing `$C10DF3` exec and `$2E78-$2E7F` write callbacks, which live in
`mesen/bridge.lua` and die with the process. No new watches were added.

## Evidence paths

`local_artifacts/experiments/EXP-0040/` — **45 artifacts, 1.1 MB**, with
`hashes.sha256` over all of them: 39 screenshots, 4 savestates,
`bridge-commands.log`, `bridge-events.log`.

Savestate lineage (unambiguous):
`local_artifacts/scenarios/SCN-0001/07-pre-whelk/pre-whelk-recon.mss`
(sha256 `852c82a0…858`) → field healing → `pre-whelk-healed.mss` →
Config (`Cursor` → `Memory`) → `pre-whelk-healed-wait.mss` → guard beat
→ Whelk battle attempt 2 → `attempt2-frozen-f182248.mss`.
`attempt1-end.mss` branches from `pre-whelk-healed.mss`.

## Files changed

- New: `docs/experiments/EXP-0040-whelk-victory.md`,
  `docs/sessions/SESSION_004.md`, this checkpoint.
- Updated: `manifests/content-census.json` (61 entries; new
  CEN-BATTLE-0009; CEN-EVENT-0011 renamed + `confirmed`;
  CEN-MAGIC-0001, CEN-EVENT-0010 corrected),
  `manifests/experiments.json` (EXP-0040, status `blocked`),
  `data/scenarios/opening-to-whelk.json` (B17/B18/B19),
  `docs/scenarios/SCN-0001-opening-to-whelk.md`,
  `dashboards/BLOCKERS.md`, `CURRENT_FOCUS.md`, `RESEARCH_QUEUE.md`,
  `ACTIVITY_LOG.md`, `docs/checkpoints/LATEST.md`.
- Regenerated: `indexes/EXPERIMENTS.md`, `indexes/CONTENT_CENSUS.md`,
  `indexes/ROM_REGIONS.md`, `dashboards/COVERAGE.md`.
- No Go source changed.

## Tests and quality gates

`gofmt -l .` clean; `go build ./...` clean; `go vet ./...` clean;
`go test ./...` all ok; `go build ./cmd/ff6lab` ok; `ff6lab audit`
**clean**; `ff6lab census sync` **clean** (61 entries);
`ff6lab archive verify` with `FF6_ROM` set — **rom verified, 8/8 ok,
clean**; restricted-extension scan of `git ls-files` clean.

Two checks failed on first run and were fixed within preservation scope:
`census sync` → `unknown experiment EXP-0040` (registered it in
`manifests/experiments.json`); `audit` → `indexes/EXPERIMENTS.md missing
EXP-0040` (ran `ff6lab indexes generate`).

## Git status

`main`, one coherent session-closing unit committed. Worktree clean; no
incomplete or unresolved changes carried forward. The tree was clean at
session start, so no unrelated pre-existing changes were involved.

## Unresolved decisions

- Whelk's monster ids remain **Unknown** — the formation record's id
  field needs a high-bit/extension decode (FF6 exceeds 256 monsters).
  Unchanged by this unit; still a bounded future question.
- Which slot the shell counter is attributed to, and Whelk's elemental
  behavior (lightning absorption is asserted only by in-game dialogue,
  **untested**).
- The true post-victory handoff beat is still unobserved.

## Blockers

**HARD BLOCKER — no ATB model.** Further Whelk execution is deferred
until the project establishes a usable ATB model, including ACTIVE/WAIT
behavior, which submenu states qualify as pausing, the relevant timer
domains, and action-queue ordering/resolution.

**Whelk gameplay must not resume before that research.**

Evidence constraint carried forward: **every head/shell transition
observed in EXP-0040 is menu-pause-contaminated** and may not be used to
characterize Whelk's natural head/shell timing. ACTIVE-mode and
WAIT-mode timing are **separate experimental conditions**.

This run does **not** support, and no downstream record may assert: that
Whelk uses any particular frame timer; that WAIT pauses every battle
clock; that a visible submenu necessarily activates WAIT; that enemy
inactivity proves the enemy ATB is frozen; that holding the MagiTek menu
open proves anything about Whelk's timer; that the healing chosen was
optimal; or that the fight is impossible or unusually long under correct
ATB operation.

## Exact next action

**Start a new session, resume from this checkpoint, audit existing
battle infrastructure and ATB evidence, and propose the first bounded
ATB experiment before operating Mesen.**

Next-session entry instruction: use the planned **ATB-master-session
prompt**. Do not resume Whelk gameplay, do not reload the Whelk
savestates for another attempt, and do not operate Mesen before the
first bounded ATB experiment is designed and recorded.

## Recommended next command

`/resume-session`
