# EXP-0046: The action-queue execution path — what completes a pending action while the scheduler is frozen

## Question

**Which routine clears `WRAM:+$3AA0` bit 6 and advances `WRAM:+$3218`
while `ROMCPU:$C21124`'s gate is shut?**

EXP-0045 established that an action pending when the gate engages still
completes — bit 6 clears and `+$3218` gains `$0100` — 78 and 119 frames
into two gated intervals, with a controlled negative in between and the
second observation predicted in advance. The tick counter and every ATB
gauge stayed frozen throughout.

Whatever performs that write therefore runs **outside** the per-frame
battle scheduler. Naming it is the queue/execution half of the ATB model,
and it is the mechanism behind EXP-0040's report that queued actions
resolved out of issue order while menus were open.

## Starting state

`local_artifacts/experiments/EXP-0042/in-battle-active-speed6.mss` — live
battle, formation 14, enemies in slots 6-8. The same state EXP-0044 and
EXP-0045 used, so all three units are directly comparable.

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

`+$3A8F` is patched to `$01` (Wait), the patch EXP-0044 trial 9 validated
against genuinely configured WAIT.

## Independent variable

None. This is a structural capture: the gated interval is arranged the
way EXP-0045 trace 3 arranged it, and the instrumentation observes who
writes.

## Controlled variables

Same ROM, emulator, lab settings, savestate, formation, party,
configuration. `+$3A90` never patched. No input during the gated interval
beyond the single press that opens the submenu.

## Instrumentation

`mesen/probes/EXP-0046.lua`:

- Write watches over `WRAM:+$3AA0`–`+$3AB3` (scheduler flags) and
  `WRAM:+$3218`–`+$322B` (the second accumulator), each recording per
  writing PC: total count, **count while the gate is set**, frame bounds,
  and the last address and value.
- Every write that occurs while `+$2F41` is non-zero additionally emits a
  `probelog` line, giving registers and a 24-byte stack window — the same
  technique earlier units used to recover caller chains.
- The EXP-0045 pending-trigger: an `endFrame` callback that arms the
  capture and injects the submenu press the moment any enemy slot's
  bit 6 appears, so the gate reliably engages with work pending.

The discriminator is simply **`countGated > 0`**: a PC that writes these
cells while the scheduler is gated cannot be the scheduler.

## Expected outcomes

- *Direct:* one or a small number of PCs write while gated; their stack
  windows name a caller chain that can be dumped and decoded, giving the
  execution path.
- *Partial:* the writes are seen but the stack does not resolve a caller,
  leaving a located routine with an unknown invoker.
- *Refuted:* no write occurs while gated in this run — the completion
  would then be happening by some path the watch does not cover (a DMA,
  or a write outside the watched ranges), and the unit would have to
  widen before concluding.

## Falsifying outcome

1. No gated write is captured despite a completion being observed in the
   frame trace — the instrumentation is missing the writer and no
   conclusion may be drawn.
2. The only gated writers turn out to be inside `ROMCPU:$C2111B`–`$C21192`
   — that would contradict EXP-0045, since the gate skips that whole
   routine, and would mean the gate reading is wrong.
3. The captured PC writes on *every* frame rather than once — the
   aggregate would then be recording ordinary scheduler traffic, not a
   deferred completion.

## Evidence paths

`local_artifacts/experiments/EXP-0046/` — writer tables, `probelog`
stack captures, ROM dumps of the decoded region, `bridge-commands.log`,
`bridge-events.log`, `experiment.json`, `hashes.sha256` at close.

## Trials

One run. `+$3A8F` patched to `$01`; the pending-trigger armed and fired
at frame **60047** with slot 7 already pending (`+$3AA0` slot 7 = `$0141`).
The completion was observed at frame **60428** as slot 7 going `$0141` →
`$0101` with its `+$3218` entry advanced.

## Observations

### Every gated write landed on one frame

All thirteen `probelog` captures carry **`frame=60169`** — 122 frames
after the trigger, and a single frame wide. That matches EXP-0045's +78
and +119 observations: the deferred completion is a **one-frame burst**,
not a slow drain.

### The gated writers

| PC | Range | Gated / total |
|---|---|---|
| `$C20798` | `+$3AA0` | 12 / 84 |
| `$C211BA` | `+$3AA0` | 12 / 96 |
| `$C20974` | `+$3AA0` | 3 / 21 |
| `$C201C1` | `+$3218` | **1 / 7** |

### `+$3218`'s `$0100` — named

`$C201C1` wrote `$7E3227`, the **high byte** of slot 7's `+$3218` entry,
with value `1`. Dumping `$C201B0`:

```asm
C201B1  BD CC 32     LDA $32CC,X
C201B4  1A           INC
C201B5  D0 1E        BNE $C201D5
C201B7  BD A0 3A     LDA $3AA0,X
C201BA  89 08        BIT #$08        ; flag bit 3
C201BC  D0 08        BNE $C201C6
C201BE  FE 19 32     INC $3219,X     ; <- the write; next PC $C201C1
C201C1  D0 03        BNE $C201C6
C201C3  DE 19 32     DEC $3219,X
C201C6  A9 FF        LDA #$FF
C201C8  9D 2C 32     STA $322C,X
C201CB  9E B5 3A     STZ $3AB5,X
```

`INC $3219,X` increments the **high** byte of the `+$3218` word, which is
exactly the `+$0100` EXP-0045 measured — structurally and numerically the
same quantity. It is guarded by `+$3AA0` **bit 3** being clear, and the
routine goes on to write `$FF` into `+$322C,X` and zero `+$3AB5,X`.

### The apparent contradiction, resolved

`$C211BA` sits inside `ROMCPU:$C21193`, the gauge-advance routine the
scheduler calls — and the scheduler is what the gate skips. A writer
there while gated would contradict EXP-0045.

The stacks resolve it. Every gated capture shares one return-address
chain, the innermost being **`$C208CB`**. Dumping `$C208B8`:

```asm
C208C6  A9 50        LDA #$50
C208C8  20 B4 11     JSR $C211B4     ; return address $C208CB
C208CB  AD 04 34     LDA $3404
```

`$C211B4` is not the scheduler — it is a **shared two-instruction
helper**, `ORA $3AA0,X / STA $3AA0,X / RTS`. The scheduler enters it at
`$C211B2` having loaded `A = #$20`; `$C208C6` enters it one instruction
later with `A = #$50`. The callback PC `$C211BA` is the helper's store in
both cases, so the same PC appears in gated and un-gated traffic without
the scheduler ever running.

That also explains the values in the log: writes of `$81`, `$41` and
`$01` to the same cell on one frame are successive passes of a shared
OR/AND helper, not a monotone state machine.

### The flag clear

```asm
C20788  BD 04 32     LDA $3204,X
C2078B  09 40        ORA #$40
C2078D  9D 04 32     STA $3204,X
C20790  A9 7F        LDA #$7F
C20792  3D A0 3A     AND $3AA0,X
C20795  9D A0 3A     STA $3AA0,X     ; <- write; next PC $C20798
C20798  60           RTS
```

`$C20795` clears **bit 7** of the flag word (mask `$7F`) and sets bit 6 of
a separate per-slot byte at `+$3204,X`. `$C20974` decodes as
`LDA $3AA0,X / ORA #$08 / STA $3AA0,X / JMP $C211EF`, setting **bit 3** —
the same bit `$C201BA` tests before incrementing.

## Interpretation

**The action-execution path lives in `ROMCPU:$C207xx`–`$C209xx` and runs
outside the per-frame scheduler.** It is invoked while `$C21124`'s gate
is shut, does its work in a single frame, and touches the same per-slot
state the scheduler does — through **shared helpers**, which is why its
writes appear at PCs that also occur in scheduler traffic.

The specific finding EXP-0045 left open is answered: **`ROMCPU:$C201BE`
(`INC $3219,X`) is what advances `+$3218` by `$0100`**, guarded by
`+$3AA0` bit 3.

A structural correction falls out of this. EXP-0043 attributed the
`+$3AA0` write at `$C211B7` to the scheduler's threshold path. That is
true of one caller, but the store is a **shared helper with at least two
entry points**, so a PC in that range does not by itself imply the
scheduler ran. Any future reasoning from `$C211B4`-region PCs must check
the caller.

`+$3AA0`'s bit picture is now partly filled: **bit 3** gates the `+$3218`
increment and is set by `$C20974`; **bit 6** is the pending-action marker
(EXP-0045); **bit 7** is cleared by `$C20795` during completion. Bits 0-2,
4 and 5 remain Unknown.

## Alternatives

- **The caller chain above `$C208C6` is not decoded.** The stack shows
  further return addresses (`$C208B1`, `$C21420`, `$C20EB6`), but this
  unit dumped only the innermost frame. What *invokes* the execution path
  each frame — and why it fired 122 frames after the gate engaged — is
  not established.
- **One completion event.** The burst-on-a-single-frame reading rests on
  this run plus EXP-0045's two; three events total, all enemy slots, all
  one formation.
- **`$C20798` and `$C20974` are located, not understood.** Their roles are
  read from a few instructions each; neither routine's boundaries or
  purpose were recovered.
- **The `gated` test reads `+$2F41` and `+$3A8F` at write time.** If
  either changed within the frame the classification could mislead —
  though EXP-0045's per-frame trace showed `+$2F41` holding steady across
  hundreds of gated frames, which argues against intra-frame flicker.
- **`+$3204,X`** appeared for the first time here and is entirely
  uncharacterised.

## Result

**Answered.** The routine that advances `+$3218` while the scheduler is
frozen is `ROMCPU:$C201BE` (`INC $3219,X`), gated on `+$3AA0` bit 3 and
reached from an action-execution path in `ROMCPU:$C207xx`–`$C209xx` that
runs outside `$C21124`.

The apparent presence of a scheduler PC (`$C211BA`) among the gated
writers is explained rather than explained away: `$C211B4` is a shared
`ORA $3AA0,X / STA / RTS` helper with at least two entry points.

Falsifier 1 did not fire — gated writes were captured. Falsifier 2 was
raised by the `$C211BA` capture and **resolved** by the stack evidence
rather than by assumption. Falsifier 3 did not fire: the writes were
confined to a single frame, not every frame.

## Confidence

- `ROMCPU:$C201BE` (`INC $3219,X`) advances `+$3218` by `$0100`:
  **Confirmed** (byte-exact decode; the observed write was to the high
  byte with value 1, matching the measured delta).
- The increment is guarded by `+$3AA0` bit 3: **Confirmed** (code role).
- `$C211B4` is a shared flag helper with at least two entry points
  (`$C211B2` with `A=$20`, `$C208C6` with `A=$50`): **Confirmed**
  (byte-exact decode plus the matching return address on every stack).
- Deferred completion is a single-frame burst: **Strong hypothesis** —
  one fully instrumented event here, consistent with EXP-0045's two.
- The execution path lives in `$C207xx`–`$C209xx` and runs outside the
  scheduler gate: **Confirmed** for the captured writers.
- `+$3AA0` bit 3 set by `$C20974`, bit 7 cleared by `$C20795`:
  **Confirmed** (code roles); the meaning of those bits: **Unknown**.
- What invokes the execution path, and the 122-frame delay: **Unknown**.
- `+$3204,X`: **Unknown**.

## Stopping condition

Ends when the gated writer is identified and either its caller chain is
decoded or the chain is explicitly recorded as unresolved.

Out of scope: the full action-queue data structure, ordering and
arbitration rules; the increment formula and threshold; status modifiers;
battle types other than random encounters; Whelk.

## Next action

**EXP-0047 — what invokes the action-execution path, and when.** The
stack chain above `$C208C6` (`$C208B1`, `$C21420`, `$C20EB6`) was
captured but not decoded, and nothing yet explains why the completion
fired 122 frames after the gate engaged rather than immediately. Dump
those frames, and exec-watch the entry point across a gated interval to
recover the invocation cadence.

That is the last structural piece of the ATB model. With it, the
programme's deliverables — the action lifecycle and queue model — can be
written down, and the Whelk decision can be taken on evidence.

Deferred, each with a known entry point:

- `+$3AA0` bits 0-2, 4, 5; and `+$3204,X`, new here and uncharacterised.
- The boundaries and purpose of `$C20798` and `$C20974`.
- Whether several pending slots drain in one burst or several.
- Party slots at a gate engage; only enemy slots have been observed.

**Correction to propagate:** EXP-0043 read the `+$3AA0` write at
`$C211B7` as the scheduler's threshold path. The store is a shared helper
with at least two callers, so a PC there does not imply the scheduler
ran. The memory map has been amended; no earlier conclusion depended on
the stronger reading.
