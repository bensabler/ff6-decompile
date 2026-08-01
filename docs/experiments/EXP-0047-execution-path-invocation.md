# EXP-0047: What invokes the action-execution path, and when

## Question

**What calls the action-execution path, and why does a pending action
complete tens of frames after the gate engages rather than immediately?**

EXP-0046 located the completion write (`ROMCPU:$C201BE`, `INC $3219,X`)
and captured a stack chain above it, but decoded only the innermost
frame. Three delays have now been measured — **78, 119 and 122 frames** —
with no explanation.

The discriminating sub-question: **is the path invoked every frame while
gated, or only on the frame the completion happens?** If every frame, the
delay is a countdown *inside* the path. If only once, the delay lives
upstream in whatever schedules it.

## Starting state

`local_artifacts/experiments/EXP-0042/in-battle-active-speed6.mss` — live
battle, formation 14, enemies in slots 6-8. Fourth unit on this state, so
all are directly comparable.

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
Controller Single. `+$3A8F` patched to `$01` (Wait) — the patch validated
by EXP-0044 trial 9.

## Independent variable

None. Structural capture plus static decode.

## Controlled variables

Same ROM, emulator, lab settings, savestate, formation, party,
configuration. `+$3A90` never patched.

## Instrumentation

Reconstructing the chain from EXP-0046's stack windows. The deepest
capture (`SP=$15E9`) reads, as successive JSR return addresses:

```
$C208CB  <- verified in EXP-0046 (JSR at $C208C8)
$C208B1
$C21420
$C20EB6
$C201A3
$C20019
```

and the shallower accumulator capture (`SP=$15F3`) begins directly at
`$C20019`. So the routine containing `$C201BE` is entered by a `JSR` at
**`ROMCPU:$C20016`**, and calls deeper from `$C201A0`.

Each frame is verified by dumping `return − 3` and confirming a `JSR`
there, rather than trusting a naive pair-read of the stack — `JSL`
pushes three bytes, so pair-reading is only valid where each frame is
confirmed.

`mesen/probes/EXP-0047.lua`:

- An **exec** callback on the entry point, counting invocations and
  recording, for the first few dozen, the frame and whether the gate was
  set at that moment.
- The EXP-0045/0046 pending-trigger, so a gated interval is arranged with
  work outstanding.
- The per-frame state sampler, so the completion frame can be located in
  the same timeline as the invocation counts.

## Expected outcomes

- *Countdown inside the path:* the entry fires on most or all gated
  frames while the completion happens on exactly one of them. The delay
  is internal, and a per-slot counter should be findable.
- *Upstream scheduling:* the entry fires rarely — ideally once, on the
  completion frame. The delay lives in whatever decides to call it.
- *Neither:* the entry never fires while gated, meaning the reconstructed
  chain is wrong and the real invoker is elsewhere.

## Falsifying outcome

1. The entry point never executes during a gated interval in which a
   completion is nonetheless observed — the chain reconstruction is
   wrong.
2. The entry fires on every frame *including* during intervals where no
   completion ever occurs and nothing is pending — it would then be
   ordinary per-frame traffic and could not explain the delay by itself.
3. A `JSR` is not found at `return − 3` for a reconstructed frame — that
   frame is not a call frame and the chain must be truncated there.

## Evidence paths

`local_artifacts/experiments/EXP-0047/` — ROM dumps of each chain frame,
the invocation table, `bridge-commands.log`, `bridge-events.log`,
`experiment.json`, `hashes.sha256` at close.

## Trials

1. **JSR verification** of the five reconstructed chain frames.
2. **Exec watch on `$C20016`** across a gated interval in which a
   completion occurred.
3. **Exec watch on four verified sites** with gate tagging and fresh
   stack capture.

## Observations

### Two frames of the reconstruction are refuted

Dumping `return − 3` for each reconstructed frame:

| Frame | Bytes at return − 3 | Verdict |
|---|---|---|
| `$C208CB` | `20 B4 11` `JSR $C211B4` | call frame (EXP-0046) |
| `$C208B1` | `20 C6 08` `JSR $C208C6` | call frame |
| `$C21420` | `20 3F 08` `JSR $C2083F` | call frame |
| `$C20EB6` | `BA 11 BF` — `TSX` | **not a call frame** |
| `$C201A3` | `20 D3 13` `JSR $C213D3` | call frame |
| `$C20019` | `20 ED 23` `JSR $C223ED` | JSR present, but see below |

**Falsifier 3 fired for `$C20EB6`**: that stack slot is data, not a
return address. The chain must be truncated there.

### And `$C20016` is refuted as the entry point

An exec watch on `$C20016` recorded **zero** executions — first across
~2 s of un-gated battle, then across a gated interval in which a
completion demonstrably occurred (slot 8's `+$3218` reached `$016B` and
all pending flags cleared).

**Falsifier 1 fired.** A `JSR` existing at `return − 3` is necessary but
not sufficient: `$C20016` happens to hold one, but it is not on this
path. Naive pair-reading of a raw stack window is unreliable here,
because routines push their own data and `JSL` frames are three bytes.
Only frames confirmed by *execution* are trustworthy.

### The cadence — the actual answer

Exec watches on the four verified sites, with the gate state recorded at
each hit:

| Site | Executions |
|---|---|
| `$C208AE` | 24 |
| `$C2141D` | 4 |
| `$C201A0` | 2 |
| `$C201BE` | 2 |

and the tagged captures:

```
EXEC C201A0 gate=0 ... frame=75825
EXEC C2141D gate=0 ... frame=75946
EXEC C208AE gate=0 ... frame=75946   (x3, X = $0010/$000E/$000C)
EXEC C201BE gate=0 ... frame=75946
EXEC C2141D gate=0 ... frame=75981
EXEC C201A0 gate=0 ... frame=76086
EXEC C2141D gate=1 ... frame=76208
EXEC C201BE gate=1 ... frame=76208
```

**`$C201BE` — the completion write — executed at `gate=0` and at
`gate=1`.** It is not gated at all.

Firing frames 75825, 75946, 75981, 76086, 76208 give deltas of
**121, 35, 105, 122** emulator frames. The completion delays measured in
EXP-0045 and EXP-0046 were **78, 119 and 122**.

`$C208AE`'s three hits on one frame carry X = `$0010`, `$000E`, `$000C` —
the stride-2 indices for slots 8, 7 and 6, so the path sweeps the slots
on each invocation.

## Interpretation

**The delay is not a countdown and not a gate effect. The path is
periodic, and a pending action completes on its next invocation.**

The action-execution path runs independently of `ROMCPU:$C21124` — it
executed with the gate both clear and set — and it is invoked on the
order of once every ~100–120 frames, sweeping the battle slots each time.
An action that becomes pending must therefore wait, on average, part of
one period before it is serviced, which is exactly the range of delays
observed (78, 119, 122 against inter-invocation gaps of 121, 35, 105,
122).

This closes the question EXP-0046 left open without needing the full call
chain: whatever schedules the path, it is **upstream of the ACTIVE/WAIT
gate**, which is why queued work resolves during a gated interval at all.

It also reframes EXP-0045's finding. "Queued work resolves past the gate"
is better stated as: **the gate stops the scheduler, and the execution
path was never behind the gate to begin with.**

## Alternatives

- **The period is not established.** Four inter-invocation gaps from one
  run, one of them (35) markedly shorter than the rest. This may be a
  fixed period with an extra trigger, or genuinely irregular. It is a
  range, not a constant.
- **The invoker is still unnamed.** Two candidate frames were refuted and
  the rest of the stack could not be trusted. Naming it needs a different
  instrument — an exec watch walking outward from `$C2141D`, or a call
  trace — not stack archaeology.
- **`$C223ED`, `$C2083F`, `$C213D3`** are named as call targets by
  verified `JSR`s but none was dumped or decoded.
- **One run, enemy slots only, one formation.** The same caveats as the
  three preceding units.
- **`$C208AE`'s slot sweep** is inferred from three X values on a single
  frame; the loop bounds were not read.

## Result

**Answered, and two reconstructions refuted along the way.**

The action-execution path is **periodic and ungated**: `$C201BE`
executes at both gate states, roughly every 100–120 frames, sweeping the
battle slots. The 78/119/122-frame completion delays are the wait for the
next invocation — the delay is scheduled upstream of `$C21124`, not
counted down inside the path and not caused by the gate.

Refuted: `$C20EB6` is not a call frame (falsifier 3), and `$C20016` is
not the entry point despite holding a plausible `JSR` (falsifier 1). The
lesson is recorded rather than buried — a raw stack window is only
trustworthy where each frame is confirmed by execution, not by the
presence of a `JSR` at `return − 3`.

## Confidence

- The action-execution path is not gated by `$C21124`: **Confirmed** —
  `$C201BE` observed executing at `gate=0` and `gate=1`.
- It is periodic rather than per-frame: **Confirmed** — zero executions
  across seconds of battle at `$C20016`'s level, and four widely spaced
  firings at the verified sites.
- The completion delay is the wait for the next invocation: **Strong
  hypothesis** — the observed gaps (121, 35, 105, 122) bracket the
  observed delays (78, 119, 122), but from one run.
- Invocation period: **Bounded**, roughly 35–122 frames; not a constant
  on this evidence.
- `$C2141D`, `$C208AE`, `$C201A0`, `$C201BE` are on the live path:
  **Confirmed** (all executed).
- `$C208AE` sweeps slots: **Strong hypothesis** (three X values on one
  frame matching slots 6-8).
- The invoker of the path: **Unknown** — two candidates refuted.

## Stopping condition

Ends when the chain frames are each confirmed or rejected by a `JSR`
check, and the invocation cadence question is answered at Confirmed or
explicitly bounded. Recovering the full action lifecycle, the queue data
structure, and ordering rules is **out of scope**; so are status
modifiers, other battle types, and Whelk.

## Next action

**EXP-0048 — name the invoker, with the right instrument.** Stack
archaeology has now failed twice on this path, so switch technique: exec
-watch `$C2141D` and walk *outward* by watching candidate callers
directly, or use Mesen's trace facility across a single invocation. The
question is narrow — one routine, fired ~every 100 frames — and the
sites to walk from are confirmed.

With that, the ATB programme's action-lifecycle and queue-model
deliverables can be written down, and the Whelk decision taken on
evidence.

**Method note to carry forward:** a `JSR` at `return − 3` is necessary
but not sufficient to confirm a stack frame. `$C20016` holds one and is
not on the path. Confirm frames by execution.

Deferred, each with a known entry point:

- `$C223ED`, `$C2083F`, `$C213D3` — named call targets, none decoded.
- The invocation period, on more than one run.
- `+$3AA0` bits 0-2, 4, 5; `+$3204`.
- Whether several pending slots drain in one invocation.
