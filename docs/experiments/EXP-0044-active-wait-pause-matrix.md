# EXP-0044: The ACTIVE/WAIT pause matrix — what `WRAM:+$2F41` gates

## Question

**What sets and clears `WRAM:+$2F41`, and which timer domains actually
pause in each menu and presentation state under ACTIVE versus WAIT?**

EXP-0043 found the gate: `ROMCPU:$C21124` runs
`LDA $2F41 / AND $3A8F / BNE skip`, which skips the *entire* per-frame
battle update. `$3A8F` is the Wait flag; `$2F41` was zeroed at battle
entry and read `00` throughout a free-running battle, so the gate was
never seen firing. Everything ACTIVE/WAIT turns on is therefore `$2F41`.

This is the unit the ATB blocker has been asking for since EXP-0040, and
it is the one that decides whether that experiment's contaminated Whelk
timing can ever be reinterpreted.

## Starting state

`local_artifacts/experiments/EXP-0042/in-battle-active-speed6.mss` — a
live battle, formation 14, **party alive** (32/30/43), enemies in slots
6-8. Configured `Bat.Mode = Active`, `Bat.Speed 6`
(`+$3A8F` = `$00`, `+$3A90` = `$87`).

`in-battle-formation14.mss` (genuinely configured WAIT, `+$3A8F` = `$01`)
is the cross-check for the patched-WAIT condition below.

## ROM identity

`Final Fantasy III (USA).sfc`, sha256
`0f51b4fca41b7fd509e4b8f9d543151f68efa5e97b08493e4b2a0c06f5d8d5e2`.

## Emulator identity

Mesen 2.1.1 macOS, **headless** (`--testrunner --timeout=7200`,
`FF6_OUT` set), bridge v2, one instance.
`Snes.RamPowerOnState=AllZeros`, virgin SRAM.

## Battle configuration

Captured per `/battle-baseline`, `source: memory-read`:
`WRAM:+$1D4D` = `$25`, `+$1D4E` = `$00`, `+$1D54` = `$00`;
battle-local `+$3A8F` = `$00` (Active), `+$3A90` = `$87` (Bat.Speed 6).

→ Bat.Mode Active, Bat.Speed 6, Msg.Speed 3, Cmd.Set Window, Gauge On,
Sound Stereo, Cursor Reset, Reequip Optimum, Controller Single.

## Independent variable

Two, varied one at a time:

1. **Menu / presentation state** — what the player has open.
2. **Battle mode** — ACTIVE versus WAIT, applied by patching
   `WRAM:+$3A8F` in place.

Patching rather than reconfiguring is deliberate. EXP-0042 established
that battle timing runs off `+$3A8F`, not the Config bytes, so writing
that cell *is* the mode switch — and it makes ACTIVE-versus-WAIT a
one-variable comparison inside a single savestate lineage instead of two
separate route runs with divergent RNG and party state.

### Runtime patch record

| Field | Value |
|---|---|
| Address | `WRAM:+$3A8F` |
| Address space | WRAM (SNES `$7E3A8F`) |
| Original value | `$00` (Active, from battle entry) |
| Patched value | `$01` (Wait) |
| Applied / removed | per trial; recorded in the trial table |
| Reason | isolate battle mode without re-entering battle |
| Expected effect | `$C21124`'s gate can now fire whenever `$2F41` is non-zero |
| Observed effect | gate fires whenever `$2F41` is non-zero; behaviour identical to configured WAIT (trial 9) |

The patch is validated against `in-battle-formation14.mss`, where WAIT
was reached by configuration rather than by patching.

## Controlled variables

Same ROM, emulator, lab settings, savestate lineage, formation, party and
configuration. Within a comparison only the menu state or `+$3A8F`
changes. `+$3A90` is never patched, so increment values stay fixed.

## Instrumentation

- `mesen/probes/EXP-0044.lua`: `watchwrites` over `WRAM:+$2F41` to catch
  what sets and clears it, plus a sampler returning every timer domain in
  one response.
- Timer domains sampled, all located by EXP-0043:

  | Domain | Address |
  |---|---|
  | battle tick counter | `WRAM:+$3A3E` (16-bit) |
  | ATB gauges | `WRAM:+$3AB4` (10 × 2) |
  | scheduler flags | `WRAM:+$3AA0` (10 × 2) |
  | second accumulator | `WRAM:+$3218` (10 × 2) |
  | gate flag | `WRAM:+$2F41` |

- Emulator frames come from `emu.getState().frameCount` and are reported
  alongside every sample, so "frames" is never ambiguous: **emulator
  frames** are the independent clock, **battle ticks** (`+$3A3E`) are a
  dependent one.

## Expected outcomes

- *Direct:* `$2F41` is set on entering some menu states and cleared on
  leaving them; under WAIT those states freeze the tick counter and all
  per-slot domains together, while under ACTIVE nothing freezes.
- *Partial:* `$2F41` gates the scheduler but some domain still advances
  in a paused state — for example a queued action resolving, or an
  animation counter outside this gate. That would be the most
  interesting result and the one the master programme explicitly warns to
  look for.
- *Refuted:* `$2F41` never becomes non-zero in any menu state, meaning
  ACTIVE/WAIT is enforced somewhere else entirely and `$C21124` is not
  the whole story.

## Falsifying outcome

1. No menu state sets `$2F41` → the `$C21124` gate is not the ACTIVE/WAIT
   mechanism, and EXP-0043's interpretation must be scoped.
2. Patched WAIT (`+$3A8F` ← `$01`) behaves differently from configured
   WAIT in `in-battle-formation14.mss` → the patch is not a valid
   substitute for the configuration, and every patched trial is void.
3. A domain advances identically in ACTIVE and WAIT for a state where
   `$2F41` is set → that domain is not gated by `$C21124`, and the matrix
   must record it as an independent clock rather than a paused one.

## Evidence paths

`local_artifacts/experiments/EXP-0044/` — writer table for `$2F41`,
per-state domain samples, savestates at the state boundaries,
`bridge-commands.log`, `bridge-events.log`, `experiment.json`, and
`hashes.sha256` written at close, after which the directory is frozen.

## Trials

Every sample reports the **emulator frame** alongside the **battle tick**
(`+$3A3E`), so the two clocks are never conflated.

| # | State | Mode | Result |
|---|---|---|---|
| 1 | command window open | ACTIVE (configured) | all advance |
| 2 | command window open | WAIT (patched) | all advance |
| 3 | MagiTek submenu open | WAIT (patched) | **all frozen** |
| 4 | same submenu, flipped in place | ACTIVE (patched) | all resume |
| 5 | same submenu, flipped back | WAIT (patched) | **all frozen** |
| 6 | submenu closed | WAIT (patched) | all resume |
| 7 | target selection | WAIT (patched) | **all frozen** |
| 8 | action resolving | WAIT (patched) | all advance |
| 9 | submenu open | **WAIT (configured)** | **all frozen** — falsifier 2 |

## Observations

### What writes `+$2F41`

Three writers over the session:

| PC | Instruction | Count |
|---|---|---|
| `$C17A92` | `STZ $2F41` | 11 285 |
| `$C17C01` | `INC $2F41` | 2 |
| `$C14434` | `STZ $2F41` | 1 |

`$C17C01` fired **exactly twice** — once per submenu opened. The
dominant writer is a clear, not a set: the flag's resting state is zero
and a submenu raises it.

### Trial 2 — WAIT with only the command window open

```
f=62768 tick=0CC8 gate2F41=00 wait3A8F=01 spd3A90=87
f=63068 tick=0D5E gate2F41=00 wait3A8F=01 spd3A90=87
gauge slot7  0022 -> 0070   (+$4E, one increment)
```

Tick advanced 150 over 300 emulator frames. **The main battle command
window does not set the gate**, so WAIT does not pause merely because a
command menu is on screen.

### Trials 3, 5 — WAIT with a submenu open

```
f=72108 tick=1184 gate2F41=01 wait3A8F=01
f=72318 tick=1184   ...
f=72558 tick=1184   ...
f=72788 tick=1184   ...
f=73028 tick=1184 gate2F41=01 wait3A8F=01
gauge3AB4 / flags3AA0 / accum3218: byte-identical throughout
```

**Every located domain frozen across 920 emulator frames.**

### Trial 4 — the one-variable flip

Same submenu, `$2F41` = `01` throughout, only `+$3A8F` written:

| `+$3A8F` | frames | tick | gauges |
|---|---|---|---|
| `$01` (Wait) | 310 | `$10E0` → `$10E0` | unchanged |
| `$00` (Active) | 300 | `$10EA` → `$117A` (+144) | slot 6 `$00E6` → `$079A`, slot 8 wrapped |
| `$01` (Wait) | 920 | frozen | unchanged |

### Trial 6 — closing the submenu

`$2F41` returned to `00`; tick resumed `$11C9` → `$125A` (+145 / 290
frames) and gauge slot 7 advanced `$000C` → `$005A`. The set/clear cycle
is reversible.

### Trial 8 — action resolution

After confirming the ability, sampled six times across 330 emulator
frames while the action played out:

```
f=80668 tick=152F gate2F41=00 wait3A8F=01
f=80998 tick=15C0 gate2F41=00 wait3A8F=01
```

`$2F41` reads **`00`** and the tick advances continuously. **Action
animations do not pause under WAIT.**

### Trial 9 — falsifier 2

`in-battle-formation14.mss`, WAIT reached by *configuration* rather than
by patching (`battleconfig()` reports `Bat.Mode=Wait`, `$3A90` = `$CF`,
"agrees"). With the MagiTek list open: `gate2F41=01`, tick frozen at
`$17C7` across three samples. **Identical to patched WAIT.**

### Contextual — a settling transient

Immediately after the gate engaged in trial 5, one sample pair showed
slot 6's `+$3AA0` low byte going `$41` → `$01` and `+$3218` slot 6
advancing by `$0100`, before everything settled and stayed frozen for the
following 920 frames. An action already in flight appears to finish after
the gate engages. **This is not resolved here**: at this sampling
resolution it cannot be distinguished from work completing on the last
un-gated frame before the patch landed.

### Party gauges

After the party member acted, `+$3AB4` slot 2 read `$0068` — non-zero and
refilling. Party slots use the same gauge array; earlier readings of
`0000` were ready-and-parked actors, not absence of storage.

## Interpretation

**`WRAM:+$2F41` is the submenu flag, and `ROMCPU:$C21124` is the whole of
ACTIVE/WAIT.**

The gate is a plain `AND` of two one-bit conditions:

- `+$3A8F` — the mode, fixed at battle entry (EXP-0042).
- `+$2F41` — raised while a qualifying **submenu** is open, cleared
  otherwise by a per-frame `STZ`.

Because it is an `AND`, ACTIVE can never pause: with `+$3A8F` = `$00` the
product is zero whatever the menu state. That is a structural certainty,
and trial 4 demonstrates it directly.

**All located timer domains share this single gate.** The tick counter,
the ATB gauges, the scheduler flags and the second accumulator froze and
resumed together in every trial. There is no evidence here of a domain
running independently of the scheduler — which is a real finding, since
the master programme specifically warned to look for one.

**The pause is narrower than the folk model.** Two states that intuition
would expect to pause do not:

- the **main battle command window**, which is on screen for most of a
  WAIT battle;
- **action animations and their resolution**.

Only a qualifying submenu — the ability list, and target selection
reached through it — raises `$2F41`.

### What this means for EXP-0040

EXP-0040's Whelk timing was called menu-pause-contaminated, and that
judgement stands, but it can now be *scoped* rather than merely asserted.
The contamination applies to intervals spent inside the MagiTek/ability
list and target selection. Intervals spent at the command window, or
watching an action resolve, were **not** paused. Whether enough of
EXP-0040's captured intervals fall outside the submenu states to
reconstruct anything is a question for the replay audit, not this unit;
the record is not reinterpreted here.

## Alternatives

- **"Qualifying submenu" is not yet enumerated.** MagiTek/ability list
  and target selection were tested. Item, Magic, Row and Defend were not,
  nor were dialogue, damage-number display, victory or defeat
  presentation. The `$C17C01` setter is a single site, which suggests one
  common path, but that is inference.
- **The settling transient** may be work genuinely running past the gate,
  or simply the final un-gated frame. Only a frame-stepped capture across
  the transition can separate them.
- **Only located domains were sampled.** Animation counters, AI script
  state, status timers and boss-state timers have no addresses yet, so
  "everything pauses together" means *everything this project can
  currently see*. A domain outside `$C21124` would not have appeared.
- **One formation, one party, one battle type.** Random encounter only;
  scripted, pincer, back and boss battles untested, and Whelk in
  particular is a boss with its own script.
- **The mode patch could in principle differ from configuration** in some
  way trial 9 did not probe — it matched on the frozen/not-frozen
  outcome, not on every downstream behaviour.

## Result

**The ACTIVE/WAIT pause matrix.** Rows are states, columns are the timer
domains located by EXP-0043; all four move together, so they are given as
one column.

| State | ACTIVE | WAIT |
|---|---|---|
| Battle running, no menu | advances | advances |
| **Main command window open** | advances | **advances** |
| MagiTek / ability list open | advances¹ | **paused** |
| Target selection | advances¹ | **paused** |
| Action resolving / animation | advances¹ | **advances** |
| Item / Magic / Row / Defend submenus | not sampled | not sampled |
| Dialogue, damage display, victory, defeat | not sampled | not sampled |

¹ Structurally implied — the gate is `$2F41 AND $3A8F`, and `$3A8F` = `0`
in ACTIVE makes the product zero regardless of `$2F41`. Directly verified
for the ability-list row (trial 4); the other two ACTIVE cells were not
sampled independently.

Domains covered: `+$3A3E` (battle tick), `+$3AB4` (ATB gauges), `+$3AA0`
(scheduler flags), `+$3218` (second accumulator).

No falsifier fired. Falsifier 1 was disposed of by observing `$2F41` = 1;
falsifier 2 by trial 9; falsifier 3 did not arise, since no domain
advanced while another was gated.

## Confidence

- `WRAM:+$2F41` is set at `ROMCPU:$C17C01` when a battle submenu opens
  and cleared per-frame at `ROMCPU:$C17A92`: **Confirmed** (writer PCs,
  counts matching submenu opens, decoded instructions).
- `ROMCPU:$C21124` gates all four located timer domains on
  `$2F41 AND $3A8F`: **Confirmed** (frozen/resumed together across nine
  trials, including an in-place one-variable flip).
- The main command window does **not** pause under WAIT: **Confirmed**.
- Action animations do **not** pause under WAIT: **Confirmed** for the
  observed action.
- The ability list and target selection **do** pause under WAIT:
  **Confirmed**.
- Patched WAIT is equivalent to configured WAIT: **Confirmed** for the
  frozen/not-frozen outcome (trial 9).
- The full set of qualifying submenus: **Unknown** — four states tested,
  the rest not sampled.
- Whether any unlocated domain (animation, AI script, status, boss-state)
  advances during a paused interval: **Unknown**.
- The settling transient after the gate engages: **Unresolved signal**.

## Stopping condition

The unit ends when `$2F41`'s writers are identified and the matrix is
filled for the reachable states under both modes, with each cell marked
`advances`, `paused`, `conditionally advances`, `not sampled` or
`unknown`.

Out of scope, each belonging to a later unit:

- The exact increment formula, threshold and gauge reset.
- The action queue, readiness arbitration and execution lifecycle.
- Status modifiers (Haste, Slow, Stop, Sleep, Stop-like states).
- Battle-type differences (scripted, pincer, back, boss).
- Whelk, which stays deferred until this matrix exists.

Presentation states that cannot be held open reliably from the bridge
are recorded as `not sampled` rather than guessed.

## Next action

**EXP-0045 — enumerate the qualifying submenus, and settle the transient.**
Two bounded pieces, one unit:

1. Walk every reachable battle menu state — Item, Magic, Row, Defend,
   plus dialogue, damage display and the victory/defeat presentations —
   sampling `+$2F41` in each. `ROMCPU:$C17C01` is a single setter, so the
   likely answer is one common path, but the matrix currently has six
   rows marked `not sampled` and they should not stay that way.
2. Frame-step across the gate transition to decide whether the settling
   transient is work genuinely resolving past the gate or the last
   un-gated frame. Use `emu.addEventCallback` on `endFrame` rather than
   bridge round-trips, which are far too coarse.

After that the ATB programme's remaining questions are the increment
formula and threshold, the action queue and readiness arbitration, and
status modifiers — none of which block the replay audit.

**Whelk is no longer blocked by an absent model.** EXP-0040's timing can
now be scoped: intervals inside the ability list and target selection
were paused, intervals at the command window and during action resolution
were not. Deciding whether to reinterpret those captures or re-run the
fight is an orchestration call, not this unit's.

Deferred, each with a known entry point:

- Domains with no address yet: animation, AI script, status, boss-state
  timers. "Everything pauses together" currently means everything the
  project can see.
- Battle types other than random encounters — Whelk is a boss with its
  own script.
- `ROMCPU:$C14434`, the second clear path, fired once and is unexplained.
