# FF6 Reconstruction Session 004 — EXP-0040 (Whelk victory attempt, stopped)

- **Date:** 2026-08-01, ~11:46–13:0x local (bridge loaded 11:46:15)
- **Investigator:** Benjamin Sabler + Claude (piloted, Lua bridge)
- **ROM identity:** SHA-256
  `0f51b4fca41b7fd509e4b8f9d543151f68efa5e97b08493e4b2a0c06f5d8d5e2`
- **Emulator:** Mesen 2.1.1 macOS, **GUI (operator-visible)**, one
  instance, foreground, normal speed.
- **Starting checkpoint:**
  [2026-08-01 — EXP-0039](../checkpoints/2026-08-01-exp0039-whelk-breadth-recon.md)

## Objective

Complete EXP-0040: defeat Whelk from the preserved pre-Whelk state by
healing the party and attacking **only while the head is exposed**, then
capture milestone `10-whelk-victory` and the first stable post-battle
state (B19).

## What the run attempted

Two piloted GUI attempts, both from the preserved pre-Whelk state and
both under `Bat.Mode = Wait`:

| Attempt | Origin | Attack | Head damage | Outcome |
|---|---|---|---|---|
| 1 | `pre-whelk-recon.mss` → field-healed | Bolt Beam | 1600 → 909 (4 hits) | stopped to inspect Config |
| 2 | `pre-whelk-healed-wait.mss` | Fire Beam | 1600 → 1246 (2 hits) | **frozen by operator directive** |

Before contact, the party was healed on the field with Tonic ×4 from
the EXP-0039 carry-in state **26 / 19 / 56** to **76/77, 105/105,
106/107**.

## Why it was stopped

Operator directive, on a **methodological blocker**: the project has no
model of FF6's ATB — ACTIVE/WAIT semantics, which submenu states
qualify as pausing, the relevant timer domains, or action-queue
ordering and resolution. Without one the battle could not be operated
efficiently or its timing interpreted.

Concretely, action selection repeatedly desynchronized from game state:
the head changed between opening a submenu and reaching target
selection, queued actions resolved out of issue order, heals landed on
unintended allies, and once the target cursor was resting on the
**shell** at the moment of confirmation and had to be cancelled. Those
are operator/tooling failures against an unmodelled system, not facts
about Whelk.

## What can be concluded (Confirmed)

- **Whelk occupies two battle slots** (new **CEN-BATTLE-0009**): shell
  slot 4 = **50000/50000 HP**, head slot 5 = **1600/1600 HP**, read from
  the DISC-0001 arrays; MP candidates 120/120 and 1000/1000.
- **Head-only targeting is correct.** Six measured hits on the visibly
  extended head (162, 168, 171, 177, 181, 186) each reduced slot 5;
  slot 4 held 50000 throughout. **Falsifier 3 not triggered.** No shell
  strike occurred in either attempt.
- **Head/shell state is visually classifiable** at 4× upscale.
  **Falsifier 2 not triggered.** The white target cursor occupies the
  same screen region and must not be mistaken for the head.
- **A field healing route exists before contact** — field menu
  reachable, inventory Tonic ×4 (recovers 50 HP) + Potion ×1.
- **MagiTek ability sets are character-specific and EXP-0039's record
  was incomplete**: Terra shows **eight** (Fire Beam, Bolt Beam, Ice
  Beam, Bio Blast, Heal Force, Confuser, X-fer, TekMissile); Wedge and
  Vicks show **four**. Updates CEN-MAGIC-0001.
- **The ally target cursor defaults to the caster's own slot.**
- **`Cursor = Memory`** reopens a character's ability list on their
  last-used ability.
- **The guard/Esper beat (CEN-EVENT-0011) precedes Whelk** at `(2A,07)`
  on clean, never-defeated runs, reproduced twice. This **resolves** its
  sequence position and **corrects** the earlier "post-Whelk" reading.
- Staged `+$3F44` formation record reproduced byte-identically (third
  independent confirmation).

## What cannot be concluded

This run does **not** support, and no downstream record may assert:

- that Whelk uses any particular frame timer;
- that WAIT pauses every battle clock;
- that a visible submenu necessarily activates WAIT behavior;
- that observed enemy inactivity proves the enemy ATB is frozen;
- that holding the MagiTek menu open demonstrates anything about
  Whelk's timer;
- that the healing behavior chosen here was optimal;
- that the fight is impossible or unusually long under correct ATB
  operation.

**All head/shell timing observed is menu-pause-contaminated** and must
not be used to characterize the natural cycle. ACTIVE-mode and
WAIT-mode timing are separate experimental conditions.

## Operator errors caught and corrected in-session

1. **Config misread.** This session initially reported `Bat.Mode` as
   `Active` and claimed to switch it to `Wait`. Re-inspection of the
   same capture (`34-config-before.png`) shows **"Active" grey, "Wait"
   white** — Wait was already selected; the hand cursor marks the row,
   not the selection. The right-press was a no-op. **The only
   configuration changed this session was `Cursor: Reset → Memory`.**
   Both attempts ran under WAIT.
2. **Blind multi-press batches.** Early actions fired 4–6 presses with
   fixed sleeps, which raced menu and dialogue transitions. Switched to
   verify-then-press with a screenshot before every confirmation. Two
   artifacts were renamed after the mode correction
   (`33-attempt1-end.png`, `attempt1-end.mss`).

## Run classification

- **Aborted due to missing ATB model** (primary)
- **Partial contextual capture**
- **Failed tactical attempt**

## Evidence preserved

`local_artifacts/experiments/EXP-0040/` — **45 artifacts, 1.1 MB**,
`hashes.sha256` written over all of them:

- 39 screenshots (pre-Whelk load, field menu, item list and each Tonic
  use, approach, guard beat, battle entry, full intro dialogue
  including the lightning/shell warnings, MagiTek lists for leader and
  escort, head-extended and head-retracted frames, target-cursor
  frames, Config before/after, the "Slime" action, frozen final frame).
- 4 savestates with recorded lineage:
  `pre-whelk-healed.mss` → `pre-whelk-healed-wait.mss` →
  `attempt1-end.mss`, `attempt2-frozen-f182248.mss`.
- `bridge-commands.log`, `bridge-events.log`.

**Frozen state (attempt 2):** frame **182248**; formation `B0 01` (432);
party HP Terra **51/77**, Wedge **105/105**, Vicks **107/107**; Terra MP
24/29; head **1246/1600**; shell **50000/50000**; head **retracted**;
menu depth = Terra's **MagiTek submenu open, cursor on Fire Beam**, no
target selection.

Savestate lineage: `pre-whelk-recon.mss` (sha256 `852c82a0…858`) →
field healing → `pre-whelk-healed.mss` → Config (`Cursor` → `Memory`) →
`pre-whelk-healed-wait.mss` → guard beat → Whelk attempt 2 → frozen.

## Process cleanup

One background process existed: GUI Mesen **PID 62332** (parent 61294,
the launching shell), purpose "visible GUI Mesen + bridge for EXP-0040",
owning `mesen/out/*`. Bridge output is written with open/append/close per
line, so it was already flushed; logs were copied into the experiment
directory before shutdown. Mesen ignores SIGTERM, so it was terminated
with `kill -9` and confirmed gone; the background task reaped with exit
code 1 as expected. `pgrep` for `Mesen|mesen|testrunner|analyze` is
clean. **SRAM directory verified still empty** (virgin boot preserved —
a `kill -9`'d GUI writes no `.srm`).

## Repository changes

- **New:** `docs/experiments/EXP-0040-whelk-victory.md` (pre-registered
  design preserved, plus amendments, observations, result, explicit
  not-claimed list, blocker).
- **New:** `docs/sessions/SESSION_004.md` (this file).
- **Updated:** `manifests/content-census.json` — new **CEN-BATTLE-0009**;
  CEN-EVENT-0011 renamed and promoted to `confirmed`; CEN-MAGIC-0001 and
  CEN-EVENT-0010 corrected/extended (61 entries).
- **Updated:** `manifests/experiments.json` — EXP-0040 registered with
  status `blocked`.
- **Updated:** `data/scenarios/opening-to-whelk.json` (B17/B18/B19) and
  `docs/scenarios/SCN-0001-opening-to-whelk.md`.
- **Updated:** `dashboards/BLOCKERS.md` (new hard blocker),
  `CURRENT_FOCUS.md`, `RESEARCH_QUEUE.md`, `ACTIVITY_LOG.md`.
- **Regenerated:** `indexes/EXPERIMENTS.md`, `indexes/CONTENT_CENSUS.md`,
  `indexes/ROM_REGIONS.md`, `dashboards/COVERAGE.md`.
- No Go source changed. Raw evidence stays under ignored
  `local_artifacts/`; nothing restricted is tracked.

## Quality results

| Check | Result |
|---|---|
| `gofmt -l .` | clean (no output) |
| `go build ./...` | clean |
| `go vet ./...` | clean |
| `go test ./...` | all packages ok |
| `go build ./cmd/ff6lab` | ok |
| `ff6lab audit` | **clean** (after registering EXP-0040 + regenerating indexes) |
| `ff6lab census sync` | **clean** (61 entries) |
| `ff6lab archive verify` (`FF6_ROM` set) | rom verified, **8/8 ok, clean** |
| restricted-extension scan of `git ls-files` | clean |

Two checks failed on first run and were fixed within preservation
scope: `census sync` reported `unknown experiment EXP-0040` (fixed by
registering it in `manifests/experiments.json`), and `audit` reported
`indexes/EXPERIMENTS.md missing EXP-0040` (fixed by
`ff6lab indexes generate`).

## Remaining worktree state

Committed in one session-closing commit. No uncommitted or incomplete
items were left behind; no pre-existing unrelated worktree changes
existed at session start (the tree was clean on `main`).

## Methodological blocker

**Further Whelk execution is deferred until the project establishes a
usable ATB model**, including ACTIVE/WAIT behavior, qualifying submenu
pause states, relevant timer domains, and action-queue behavior. Whelk
gameplay must not resume before that research. Milestone
`10-whelk-victory` and B19 remain open.

## Exact next action

Start a new session, resume from this checkpoint, audit existing battle
infrastructure and ATB evidence, and propose the first bounded ATB
experiment before operating Mesen.
