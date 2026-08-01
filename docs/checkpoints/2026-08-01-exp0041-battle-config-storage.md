# Checkpoint 2026-08-01 — EXP-0041: battle configuration storage (Unit 43)

## Current question

None open in the unit. EXP-0041 is **completed**: all nine Config
settings are located and bit-decoded, and no falsifier fired.

## State

The ATB research program is **open and has produced its first result**.
The hard blocker from EXP-0040 **remains open** — no timer domain, pause
condition, or action-queue semantics is known — but its configuration
half is closed: `Bat.Mode` and `Bat.Speed` can now be read and set from
memory, so ACTIVE and WAIT can be established as controlled conditions
instead of eyeballed off a screenshot.

Whelk gameplay was **not** resumed and the Whelk savestates were **not**
reloaded, per the standing directive.

## Confirmed before this session

Golden route power-on → milestone `06-random-encounter`; event-flag
system (DISC-0008); unified battle arrays (DISC-0001); formation table;
Whelk reached, entered, not defeated (EXP-0039/0040).

## Work completed

Three narrow commits.

1. **Workflow preparation.** `## Battle configuration` added to the
   experiment template; `starting_state.battle_config` declared in
   `schemas/experiment.schema.json` with a required `source` of
   `memory-read` / `screen-read` / `mixed`; and
   `internal/audit.CheckBattleExperimentConfig` registered in
   `audit.Run` to enforce it for battle experiments **from EXP-0041
   onward** (the same cutover convention `EVIDENCE_LAYOUT.md` uses, so no
   historical entry needed retrofitting). Schemas in this repo are
   validated by nothing, so the audit is where the requirement became
   real. Verified end-to-end against the live manifest: deleting
   EXP-0041's `battle_config` produces the finding, restoring it clears
   it.

2. **`ff6lab state`.** New `internal/mesenstate` parses Mesen `.mss`
   files — signature, zlib streams, and a flat sequence of
   NUL-terminated name / little-endian length / data blocks — exposing
   `memoryManager.workRam` (131072 B) and `cart.saveRam` (8192 B).
   Subcommands `list|wram|sram|read|diff`. `SRAM:+$xxxx` registered in
   `ADDRESS_NOTATION.md` because the tool emits those addresses and no
   prefix existed. This makes ~20 preserved savestates readable without
   an emulator.

3. **EXP-0041.** The configuration map:

   | Byte | Contents | Default |
   |---|---|---|
   | `WRAM:+$1D4D` | bits 0-2 Bat.Speed (0-5 → displayed 1-6), bit 3 **Bat.Mode** (0 = Active, 1 = Wait), bits 4-6 Msg.Speed (0-5), bit 7 Cmd.Set | `$2A` |
   | `WRAM:+$1D4E` | bit 4 Reequip, bit 5 Sound, bit 6 Cursor, bit 7 Gauge | `$00` |
   | `WRAM:+$1D54` | bit 7 Controller | `$00` |

   - The block is **not contiguous**. A three-byte read window made
     Controller look like a null result until it was widened; the
     full-WRAM diff is what caught it.
   - A cleared bit is always the left-hand screen option.
   - Both speed fields were **swept to both clamps**, so 0-5 is measured,
     not inferred from the displayed range.
   - The Config screen marks the **selected** option's text cells with
     tile attribute `$20` and the unselected with `$28` — the inverse of
     the intuitive reading. This is the exact mechanism behind EXP-0040's
     `Bat.Mode` misread, and it **independently confirms that
     correction**: `Wait` is where a new game arrives on this revision.
   - Configuration is **not SRAM-backed** before a save: `cart.saveRam`
     was byte-identical across a config change and virgin throughout.
   - Trial 0 was retrospective and used no emulator — `ff6lab state diff`
     over EXP-0040's preserved single-config-change savestate pair
     yielded `WRAM:+$1D4E` `$00` → `$40` as the Cursor candidate. The
     live falsifier run, from a clean unrelated savestate with only the
     Cursor row touched, reproduced that transition exactly.

   New **CEN-MENU-0007**. Memory map extended with the three storage
   bytes and the eighteen Config-screen text-cell addresses.

## Last raw observation

Falsifier run, fresh reload of `05-mines-entry.mss` → field menu →
Config → Cursor row → one RIGHT press: `WRAM:+$1D4D`-`+$1D54` read
`2A 40 00 12 34 56 06 00`; the screen's "Reset" cells took attribute
`$28` and "Memory" `$20`.

## Active emulator state

**None.** Two headless `--testrunner` instances ran in sequence (the
second after a relaunch with a corrected working directory); logs were
harvested before termination, both killed with `kill -9`, absence
confirmed by `pgrep`. `jobs -l` empty. **SRAM directory verified still
empty** — virgin boot preserved, as expected for `kill -9` (a clean
`emu.stop()` exit would have written the `.srm`).

## Breakpoints/watchers

None installed. The unit used bridge commands only
(`loadstate`/`press`/`read`/`screenshot`/`eval`) plus the bridge's
standing `$C10DF3` exec and `$2E78-$2E7F` write callbacks, which die with
the process. No probe was written — savestate capture plus
`ff6lab state diff` covered the whole-WRAM net more cheaply.

## Evidence paths

`local_artifacts/experiments/EXP-0041/` — **17 artifacts, 616 KB**, with
`hashes.sha256` over all of them: 8 screenshots, 4 savestates,
`trial0-cursor-diff.txt`, `bridge-commands.log`, `bridge-events.log`,
`experiment.json`.

Starting state: `local_artifacts/scenarios/SCN-0001/05-mines-entry/05-mines-entry.mss`
(milestone 05, frame 51 578, EXP-0036, three byte-identical runs).
Trial 0 read EXP-0040's `pre-whelk-healed.mss` /
`pre-whelk-healed-wait.mss` **read-only and offline**; neither was
loaded into an emulator.

## Files changed

- New: `internal/mesenstate/` (+ tests), `cmd/ff6lab/state.go` (+ tests),
  `docs/experiments/EXP-0041-battle-config-storage.md`, this checkpoint.
- Updated: `.claude/templates/EXPERIMENT.md`,
  `schemas/experiment.schema.json`, `internal/audit/audit.go` (+ tests),
  `internal/project/project.go`, `.claude/skills/experiment-designer/`
  and `mesen-operator/SKILL.md`, `.claude/skills/_shared/ADDRESS_SPACES.md`,
  `docs/research/ADDRESS_NOTATION.md`, `docs/sessions/04_MEMORY_MAP.md`,
  `manifests/experiments.json`, `manifests/content-census.json`,
  `dashboards/` (CURRENT_FOCUS, BLOCKERS, RESEARCH_QUEUE, MILESTONES,
  ACTIVITY_LOG), `docs/checkpoints/LATEST.md`.
- Regenerated: `indexes/EXPERIMENTS.md`, `indexes/CONTENT_CENSUS.md`,
  `indexes/ROM_REGIONS.md`, `dashboards/COVERAGE.md`.

Also corrected in passing: `dashboards/MILESTONES.md` row S1 still
claimed milestones 00-05, predating EXP-0038's milestone 06.

## Tests and quality gates

`gofmt -l .` clean; `go build ./...`, `go vet ./...`, `go test ./...` all
ok; `ff6lab audit` clean; `ff6lab census sync` clean (62 entries);
restricted-extension scan of `git ls-files` clean. No savestate,
screenshot, or ROM-derived byte is tracked.

## Git status

`main`. Three coherent commits, worktree clean.

## Unresolved decisions

- Whether battle code reads `+$1D4D`/`+$1D4E` directly or a copy taken at
  battle entry — **the next experiment**, and it gates how every later
  ATB unit is staged.
- The routine that writes the configuration bytes, and its callers.
- `+$1D4E` bits 0-3 and `+$1D54` bits 0-6: untouched by the nine
  settings, meaning unknown.
- SRAM persistence, still blocked on capturing a save event
  (CEN-SAVE-0001).
- Whether the `$20`/`$28` attribute pair is a general menu-renderer
  convention or specific to this screen.

## Blockers

**HARD BLOCKER — no ATB model. Still open.** Timer domains, which submenu
states pause which clocks, and action-queue ordering all remain Unknown.
EXP-0041 closed only the configuration prerequisite.

**Whelk gameplay must not resume before that research.** The evidence
constraint carried forward is unchanged: every head/shell transition
observed in EXP-0040 is menu-pause-contaminated and may not be used to
characterize Whelk's natural timing; ACTIVE and WAIT are separate
experimental conditions.

This unit does **not** support, and no downstream record may assert:
that battle code reads these bytes directly; that changing configuration
mid-battle takes effect; that `Bat.Speed` affects any particular clock;
or anything about what ACTIVE or WAIT actually pause.

## Exact next action

**EXP-0042 — battle-entry configuration sampling.** Install a *read*
watch on `WRAM:+$1D4D`-`+$1D4E` and capture the reading PCs across a
battle entry, to determine whether the battle routines consult these
bytes directly or a copy taken at entry. Use the mines random encounter
(milestone 06, EXP-0038), which is reproducible on demand and needs no
Whelk state. Note `probes/common.lua` has `watchwrites` but no read
equivalent; add one there rather than creating parallel infrastructure.

## Recommended next command

`/orchestrate`
