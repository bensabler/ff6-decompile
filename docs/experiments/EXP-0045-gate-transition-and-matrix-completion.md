# EXP-0045: The gate transition frame by frame, and the rest of the matrix

## Question

Two bounded parts on the same gate.

**A — the settling transient.** When `WRAM:+$2F41` goes `0 → 1` under
WAIT, does any timer domain advance on a frame where the gate is
*already set*? EXP-0044 saw one sample pair in which a scheduler flag
cleared and `+$3218` advanced by `$0100` just after the gate engaged,
then nothing for 920 frames. At bridge-sampling resolution that cannot be
told apart from work completing on the last un-gated frame.

**B — the unsampled matrix rows.** EXP-0044 left six states marked
`not sampled`. Which of them raise `+$2F41`?

## Starting state

`local_artifacts/experiments/EXP-0042/in-battle-active-speed6.mss` — live
battle, formation 14, party alive, enemies in slots 6-8. Same state
EXP-0044 used, so the two units are directly comparable.

## ROM identity

`Final Fantasy III (USA).sfc`, sha256
`0f51b4fca41b7fd509e4b8f9d543151f68efa5e97b08493e4b2a0c06f5d8d5e2`.

## Emulator identity

Mesen 2.1.1 macOS, headless (`--testrunner --timeout=7200`, `FF6_OUT`
set), bridge v2, one instance. `Snes.RamPowerOnState=AllZeros`, virgin
SRAM.

## Battle configuration

Per `/battle-baseline`, `source: memory-read`: `WRAM:+$1D4D` = `$25`,
`+$1D4E` = `$00`, `+$1D54` = `$00`; battle-local `+$3A8F` = `$00`,
`+$3A90` = `$87`. Bat.Mode Active, Bat.Speed 6, Msg.Speed 3, Cmd.Set
Window, Gauge On, Sound Stereo, Cursor Reset, Reequip Optimum,
Controller Single.

`+$3A8F` is patched to `$01` (Wait) for part A, exactly as in EXP-0044,
whose trial 9 validated that patch against genuinely configured WAIT.

## Independent variable

Part A: none — an observational per-frame trace across a transition the
operator causes but does not otherwise vary.
Part B: the menu or presentation state, one at a time.

## Controlled variables

Same ROM, emulator, lab settings, savestate, formation, party,
configuration. `+$3A90` never patched.

## Instrumentation

`mesen/probes/EXP-0045.lua`:

- An `endFrame` callback recording **one sample per emulated frame** into
  a Lua table when armed — frame, `+$2F41`, `+$3A8F`, tick `+$3A3E`, and
  the enemy slots of `+$3AB4`, `+$3AA0`, `+$3218`. Buffered in memory and
  flushed to a file on demand, because per-frame file I/O would distort
  the very timing under study.
- `watchwrites` over `WRAM:+$2F41`, carried over from EXP-0044 so the
  writer PCs can be re-checked against the trace.

Bridge round-trips are **not** used for part A: they are hundreds of
frames apart and are exactly why the transient was unresolved.

## Expected outcomes

- *Last-un-gated-frame:* every domain advance in the trace occurs on a
  frame where `+$2F41` reads `0`. The transient was the final normal
  frame before the gate engaged, and nothing resolves past it.
- *Work resolves past the gate:* at least one domain advances on a frame
  where `+$2F41` reads `1`, meaning `$C21124` does not gate everything
  and some already-queued work completes anyway.
- *Something else:* the gate flickers, or a domain advances on gated
  frames only in the first few frames after the transition — a bounded
  drain rather than either clean answer.

## Falsifying outcome

1. For "last un-gated frame": any domain advancing on a frame with
   `+$2F41` = `1`.
2. For "work resolves past the gate": no advance on any gated frame in
   the trace.
3. If the trace shows `+$2F41` toggling every frame rather than holding,
   EXP-0044's steady-state readings were aliased and its matrix must be
   re-examined.

## Evidence paths

`local_artifacts/experiments/EXP-0045/` — per-frame trace log, per-state
samples, screenshots, `bridge-commands.log`, `bridge-events.log`,
`experiment.json`, `hashes.sha256` written at close.

## Trials

| # | Setup | Pending at gate engage | Gated-frame changes |
|---|---|---|---|
| 1 | submenu opened by hand | **slot 6** (`+$3AA0` bit 6 set) | **1**, at +78 frames |
| 2 | submenu opened by hand | **none** | **0** across 438 gated frames |
| 3 | **prediction test** — frame-accurate trigger | **slot 8** | **1**, at +119 frames |
| 4 | Item submenu | none | 0 across 370 frames |
| 5 | victory presentation | n/a (gate never set) | n/a |

## Observations

### Trace 1 — the transient is not a boundary artifact

450 frames, gate engaging at f=60560. Six domain changes before the gate,
all with `gate=00`. Then:

```
f=60559 gate=00 tick=0896 g=00AE,0086,00A0 fl=0141,0101,0101 ac=0013,7133,3CDB
f=60560 gate=01 tick=0896 g=00AE,0086,00A0 fl=0141,0101,0101 ac=0013,7133,3CDB
   ... 77 frames, byte-identical ...
f=60637 gate=01 tick=0896 g=00AE,0086,00A0 fl=0141,0101,0101 ac=0013,7133,3CDB
f=60638 gate=01 tick=0896 g=00AE,0086,00A0 fl=0101,0101,0101 ac=0113,7133,3CDB
   ... 349 frames, byte-identical ...
f=60987 gate=01 tick=0896 g=00AE,0086,00A0 fl=0101,0101,0101 ac=0113,7133,3CDB
```

At **f=60638, with the gate set**, slot 6's `+$3AA0` low byte went
`$41` → `$01` (bit 6 cleared) and its `+$3218` entry went `$0013` →
`$0113` (+`$0100`). The tick counter and every gauge stayed frozen.

77 clean frames precede it and 349 follow it, so this is neither a
boundary artifact nor an ongoing clock. **Falsifier 1 fired**: the
"last un-gated frame" explanation is refuted.

### Trace 2 — and it does not always happen

Same procedure, 460 frames. Eleven domain changes, **all** with
`gate=00`; **zero** gated-frame changes across 438 gated frames. At the
moment the gate engaged the flags read `0101,0109,0101` — **no slot had
bit 6 set**.

### Trace 3 — the prediction test

Traces 1 and 2 differ in exactly one respect: whether any slot had
`+$3AA0` bit 6 set when the gate engaged. That yields a prediction —
*engage the gate while a slot is pending and a deferred completion will
fire during the gated interval.*

Because bridge round-trips are hundreds of frames apart, the arming was
moved into Lua: an `endFrame` callback watched the enemy flag bytes and,
the moment any bit 6 appeared, started the trace and injected the press
itself.

It fired with **slot 8** pending (`0101,0101,0141`), a different slot
from trace 1. 119 frames into the gated interval:

```
f=82650 gate=01 tick=0D9A g=002C,0070,008A fl=0101,0101,0141 ac=F617,7D27,0083
f=82651 gate=01 tick=0D9A g=002C,0070,008A fl=0101,0101,0101 ac=F617,7D27,0183
```

Slot 8's bit 6 cleared and its `+$3218` gained exactly `$0100`, tick and
gauges frozen. **Predicted, and reproduced on a different slot with a
different delay.**

### Part B — the remaining states

- **Item submenu**: `+$2F41` = `01`, and 370 traced frames with **zero**
  changes in any domain. A qualifying submenu.
- **Victory presentation**: `+$2F41` = `00` across all 429 in-battle
  frames of the trace, with the tick advancing on 195 of them. **Does not
  pause.** Battle memory tore down at the 430th frame.
- **Magic, Row, Defend**: **not reachable** from this state — the Magitek
  command set offers only MagiTek and Item. They belong to a
  normal-party battle, not to this savestate.
- **Dialogue and damage display**: covered by EXP-0044's action-resolution
  row, where `+$2F41` reads `00`; not separated further here.
- **Defeat presentation**: not sampled with gate instrumentation.
  EXP-0043 incidentally observed `+$3A3E` still incrementing on the
  "Annihilated" screen, which is consistent with `advances` but was not
  measured against the gate.

## Interpretation

**Queued work resolves past the gate; nothing new starts.**

`ROMCPU:$C21124` stops the scheduler — the tick counter and every ATB
gauge froze in all three traces without exception. But an action already
*pending* when the gate engages still completes during the gated
interval, clearing `+$3AA0` bit 6 and advancing that slot's `+$3218` by
exactly `$0100`.

**`+$3AA0` bit 6 is the pending-action marker.** EXP-0043 saw it set
transiently at a gauge advance and left its meaning Unknown; EXP-0044
could not reconcile it with the `ORA #$20` at `$C211B2`. Its behaviour
here settles the reading: set when a slot has an action outstanding,
cleared when that action completes — and the completion is driven by
something *outside* `$C21124`, because it happens while that gate is
shut.

The delay differed between traces (78 versus 119 frames), so whatever
drives the completion is not a fixed constant, and this unit does not
identify it.

This is the refinement the master programme anticipated almost verbatim:
enemy ATB accumulation pauses while the menu flag is active, yet queued
enemy actions may still resolve during the interval.

### What it means for EXP-0040

EXP-0040 reported that "queued actions resolved out of issue order" while
menus were open, and treated that as an operator failure against an
unmodelled system. It was a real behaviour of the system: actions pending
when a submenu opened continued to resolve while the gauges were frozen.
That record's observation is vindicated; its interpretation of the cause
is now available.

## Alternatives

- **The completion driver is unidentified.** Something advances a pending
  action while `$C21124` is shut; this unit shows *that* it happens, not
  *what* does it. A read-watch on `+$3218`/`+$3AA0` during a gated
  interval would name it.
- **One event per trace.** Traces 1 and 3 each showed a single deferred
  completion. Whether several pending slots would all drain, or only one,
  was not tested — no trace had two slots pending at once.
- **`+$3218` advancing by exactly `$0100`** was seen twice. Whether that
  is a fixed step or an artifact of these two actions is unknown; the
  accumulator's low byte was unchanged in both.
- **Enemy slots only.** All three traces watched slots 6-8. No party slot
  was pending at a gate engage during this unit.
- **Trace 2 is a negative**, and absence of a gated change over 438
  frames is strong but not proof that nothing could ever fire; it is
  reported as the observation it is.

## Result

**Part A is answered: work already pending when the gate engages does
resolve past it.** The tick counter and ATB gauges freeze without
exception, but a pending slot's action completes — `+$3AA0` bit 6 clears
and `+$3218` advances by `$0100` — 78 and 119 frames into the gated
interval in two independent traces, the second of them **predicted in
advance** from the discriminator identified in the first two.

`+$3AA0` bit 6 is the pending-action marker, resolving a semantics
question EXP-0043 and EXP-0044 both left open.

**Part B**: Item joins the qualifying submenus; the victory presentation
does not pause; Magic, Row and Defend are unreachable from a Magitek
battle and are correctly left unsampled.

### Updated matrix

| State | ACTIVE | WAIT |
|---|---|---|
| Battle running, no menu | advances | advances |
| Main command window open | advances | advances |
| MagiTek / ability list open | advances* | **paused**, queued work still resolves |
| **Item submenu open** | advances* | **paused** |
| Target selection | advances* | **paused**, queued work still resolves |
| Action resolving / animation | advances* | advances |
| **Victory presentation** | advances* | **advances** |
| Defeat presentation | not sampled | tick observed advancing (EXP-0043, un-instrumented) |
| Magic / Row / Defend | unreachable here | unreachable here |

\* structurally implied — the gate is an `AND`, so `+$3A8F` = `$00`
makes it zero regardless of `+$2F41`.

## Confidence

- Pending work completes during a gated interval: **Confirmed** —
  observed twice, the second time predicted in advance on a different
  slot, with a controlled negative (trace 2) in between.
- The tick counter and ATB gauges never advance while gated:
  **Confirmed** — 1 245 gated frames across three traces, zero exceptions.
- `+$3AA0` bit 6 is the pending-action marker: **Strong hypothesis** —
  perfect correlation across three traces plus a mechanistic reading, but
  the routine that clears it has not been located.
- `+$3218` advancing by exactly `$0100` on completion: **Bounded** — two
  observations.
- Item is a qualifying submenu: **Confirmed**.
- The victory presentation does not pause: **Confirmed** for this battle.
- What drives a pending action to completion while the gate is shut:
  **Unknown**.
- Whether multiple pending slots all drain: **Unknown** — never observed.

## Stopping condition

Ends when part A is answered at Confirmed or explicitly bounded, and the
reachable states in part B are sampled. States that cannot be held open
or reached from this savestate are recorded as `not sampled` rather than
guessed — Row and Defend are not offered by the Magitek command set, so
they belong to a normal-party battle, not this one.

Out of scope: the increment formula and threshold, the action queue,
status modifiers, battle types other than random encounters, and Whelk.

## Next action

**EXP-0046 — name the driver of deferred completion.** Something advances
a pending action to completion while `ROMCPU:$C21124` is shut. Read-watch
`+$3AA0` and `+$3218` through a gated interval arranged the way trace 3
arranged it, and capture the writing PC. That routine is the action-queue
execution path, which is the next major piece of the ATB model and the
last one EXP-0040's "actions resolved out of issue order" observation
depends on.

Deferred, each with a known entry point:

- Whether several pending slots all drain during one gated interval.
- Whether `+$3218` always advances by `$0100`, and what its low byte is
  for.
- Magic, Row and Defend submenus — need a normal-party battle, not a
  Magitek one.
- The defeat presentation, with gate instrumentation this time.
- Party slots pending at a gate engage; only enemy slots were seen here.
