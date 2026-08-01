# Checkpoint 2026-08-01 — EXP-0043: the ATB layer is located (Unit 45)

## Current question

None open in the unit. EXP-0043 is **completed**: both the primary
question (`+$3A90`'s consumer) and the secondary one (gauge location) are
answered, and the ACTIVE/WAIT divergence point came free.

## State

**The ATB layer is located.** Every timer domain the master program asked
for now has an address, and the instruction where ACTIVE and WAIT differ
is identified. The blocker that stopped EXP-0040 is **narrow rather than
total**: what remains is the pause *condition*, the exact arithmetic, and
the action queue.

Whelk was not resumed and its savestates were not reloaded.

## Confirmed before this session

Configuration storage and bit layout (EXP-0041, CEN-MENU-0007);
battle-entry sampling and the staging rule (EXP-0042, CEN-BATTLE-0010).

## Work completed

| Address | Role | Confidence |
|---|---|---|
| `WRAM:+$3AB4` | **ATB gauge array**, 10 × 2 bytes, party 0-3 / enemies 4-9 | Confirmed |
| `WRAM:+$3AC8` | per-slot **increment**; gauge += `$3AC8,X >> 1` | Confirmed |
| `WRAM:+$3AA0` | per-slot scheduler **flags** | Confirmed (code role); bits Unknown |
| `WRAM:+$3A3E`-`+$3A3F` | 16-bit **battle tick counter** | Confirmed |
| `WRAM:+$3218` | second per-slot accumulator, gated on `$3219,X` | Confirmed (code role) |
| `WRAM:+$3B19,X` | per-slot **Speed** feeding the increment | Strong hypothesis |
| `WRAM:+$2F41` | other half of the pause gate | Tentative — **untested** |
| `ROMCPU:$C2111B-$C21192` | per-frame battle scheduler; **gate at `$C21124`** | Confirmed |
| `ROMCPU:$C21193-$C211BA` | gauge advance | Confirmed |
| `ROMCPU:$C209E0-$C20A0E` | increment builder; consumes `+$3A90` | Confirmed |

Two results worth stating plainly:

1. **`ROMCPU:$C21124` is where ACTIVE and WAIT diverge.**
   `LDA $2F41 / AND $3A8F / BNE` skips the entire per-frame battle
   update. `$3A8F` is the Wait flag from EXP-0042; `$2F41` is zeroed at
   battle entry and read `00` throughout a free-running battle, so the
   gate never fired here. Finding what sets it is the next unit.

2. **Battle Speed scales enemy gauges only.** The `CPX #$08 / BCC` branch
   at `$C209F6` skips the `$3A90` multiply for party slots (X is the
   stride-2 slot index, so X < 8 is slots 0-3). Measured against the same
   formation: party increments **byte-identical** at Bat.Speed 3 and 6
   (318/330/336); enemy increments 240 versus 156, tracking `$3A90`'s
   207/135 ratio. Structural and numerical evidence agree.

Method: three instruments were required to converge before anything was
claimed — the read watch named the routines, the static ROM decode named
the arrays, and live sampling showed `+$3AB4` advancing by exactly the
predicted `$4E` = `$9C >> 1`, including a wrap. No falsifier fired.

**Correction carried forward:** the ATB family uses **stride 2**, not the
`$14` stride of the HP/stat family. DISC-0001's unified layout governs
slot *assignment* but not stride; new battle arrays must not assume `$14`.

## Last raw observation

Frame 66038, formation 14, ACTIVE, `$3A90` = `$87`: `+$3AB4` slot 8 went
`$001E` → `$006C` (+`$4E`, exactly `$3AC8` slot 8 `$009C >> 1`) while its
`+$3AA0` low byte went `$01` → `$41`, returning to `$01` at the next
sample.

## Active emulator state

**None.** One headless `--testrunner` instance; logs harvested before
termination, killed with `kill -9`, absence confirmed by `pgrep`.
`jobs -l` empty. **SRAM directory verified still empty.**

## Breakpoints/watchers

None left installed. `mesen/probes/EXP-0043.lua` armed a read watch on
`WRAM:+$3A8F`-`+$3A90` and a WRAM window sampler; both die with the
process.

## Evidence paths

`local_artifacts/experiments/EXP-0043/` — **8 artifacts, 400 KB**, with
`hashes.sha256`: interval savestates `t0`/`t1`/`t2`, a screenshot,
`read-table.txt`, `bridge-commands.log`, `bridge-events.log`,
`experiment.json`.

Starting state: `local_artifacts/experiments/EXP-0042/in-battle-active-speed6.mss`,
with `in-battle-formation14.mss` as the Battle Speed contrast. No route
replay was needed — EXP-0042's preserved live battles were sufficient.

## Files changed

- New: `docs/experiments/EXP-0043-atb-gauges-and-speed-consumer.md`,
  `mesen/probes/EXP-0043.lua`, this checkpoint.
- Updated: `docs/sessions/04_MEMORY_MAP.md` (twelve entries),
  `manifests/experiments.json`, `manifests/content-census.json`,
  `dashboards/` (CURRENT_FOCUS, BLOCKERS, RESEARCH_QUEUE, ACTIVITY_LOG),
  `docs/checkpoints/LATEST.md`.
- Regenerated: `indexes/EXPERIMENTS.md`, `indexes/CONTENT_CENSUS.md`,
  `indexes/ROM_REGIONS.md`, `dashboards/COVERAGE.md`.

## Tests and quality gates

`gofmt -l .` clean; `go build ./...`, `go vet ./...`, `go test ./...` ok;
`ff6lab audit` clean; `ff6lab census sync` clean (64 entries);
restricted-extension scan clean. No Go source changed this unit.

## Git status

`main`, one coherent unit committed, worktree clean.

## Unresolved decisions

- **What sets and clears `WRAM:+$2F41`** — the untested half of the
  ACTIVE/WAIT gate, and the next unit's target.
- The exact increment formula, threshold, overflow and gauge reset.
- `+$3AA0` bit semantics: bit 6 was observed set where the decoded
  threshold path at `$C211B2` writes bit 5.
- The relationship between `+$3AB4` and the second accumulator `+$3218`.
- `+$322C,X`, the comparand on the threshold path.
- `+$3B19,X` as Speed is unverified against an independent source.
- Battle Speed's middle values; other formations and party compositions.
- Whether the ≈0.5-per-frame tick rate relates to open question #18.

## Blockers

**ATB blocker — now narrow, not total.** Timer domains are named and
readable. Still unknown: the pause condition, the increment arithmetic,
and the action queue.

**Whelk gameplay must not resume before the pause matrix exists.**
EXP-0040's head/shell timing remains menu-pause-contaminated and can only
be reinterpreted once the pause semantics are known.

This unit does **not** support, and no downstream record may assert:
that opening any particular menu pauses anything (`$2F41` was never
observed non-zero); that `+$3AB4` is the only ATB-like accumulator;
that the gauge threshold is any particular value; that `+$3AA0` bit 5 or
bit 6 means "ready"; or that Battle Speed behaves the same way at its
untested middle values or in other battle types.

## Exact next action

**EXP-0044 — the ACTIVE/WAIT pause matrix.** Write-watch
`WRAM:+$2F41` across opening and closing each battle submenu to find what
sets and clears it, then build the matrix of menu and presentation states
against timer domains, reading `+$3AB4`, `+$3A3E`, `+$3AA0` and `+$3218`
directly. `$3A8F`/`$3A90` can be patched in place, so ACTIVE versus WAIT
becomes a one-variable comparison inside a single savestate lineage
rather than two route runs. Start from
`local_artifacts/experiments/EXP-0042/in-battle-formation14.mss`.

This is the first unit that genuinely warrants `/battle-baseline` and
bounded parallel read-only observers over a frozen evidence package.

## Recommended next command

`/orchestrate`
