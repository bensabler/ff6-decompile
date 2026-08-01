# EXP-0043: The consumer of `WRAM:+$3A90`, and the ATB gauges

## Question

**What reads `WRAM:+$3A90`, and what does it advance?**

EXP-0042 established that Battle Speed is transformed at battle entry
into `+$3A90` = `255 − 24 × speed` and never re-read. That makes `+$3A90`
the project's sharpest lead into battle timing: a battle-local value
derived from a timing setting, whose consumer is unlocated.

The secondary question, expected to be answered by the same evidence:
**where are the ATB gauges**, and do they follow the `$14`-stride
ten-slot layout DISC-0001 established for the other per-battler arrays?

## Starting state

`local_artifacts/experiments/EXP-0043/` works from
`local_artifacts/experiments/EXP-0042/in-battle-active-speed6.mss` — a
**live battle** preserved by EXP-0042, formation 14, configured
`Bat.Mode = Active` and `Bat.Speed = 6` (`+$3A8F` = `00`,
`+$3A90` = `$87`). ACTIVE mode with no menu open is deliberate: nothing
should be pausing, so any gauge should advance freely.

`in-battle-formation14.mss` (Wait, Bat.Speed 3, `+$3A90` = `$CF`) is the
contrast state, giving a second speed value against the same formation.

No route replay is needed and no Whelk state is involved.

## ROM identity

`Final Fantasy III (USA).sfc`, sha256
`0f51b4fca41b7fd509e4b8f9d543151f68efa5e97b08493e4b2a0c06f5d8d5e2`.

## Emulator identity

Mesen 2.1.1 macOS, **headless** (`--testrunner --timeout=7200`,
`FF6_OUT` set), bridge v2, one instance. `Snes.RamPowerOnState=AllZeros`,
virgin SRAM.

## Battle configuration

| Setting | Value | Source |
|---|---|---|
| Bat.Mode | `Active` | memory-read |
| Bat.Speed | 6 | memory-read |
| Msg.Speed | 3 | memory-read |
| Cmd.Set | `Window` | memory-read |
| Gauge | `On` | memory-read |
| Sound | `Stereo` | memory-read |
| Cursor | `Reset` | memory-read |
| Reequip | `Optimum` | memory-read |
| Controller | `Single` | memory-read |

`WRAM:+$1D4D` = `$25`, `WRAM:+$1D4E` = `$00`, `WRAM:+$1D54` = `$00`;
battle-local `WRAM:+$3A8F` = `$00`, `WRAM:+$3A90` = `$87`.

Per EXP-0042's staging rule the configuration is already baked into the
battle-local cells in this savestate, so it cannot drift mid-run.

## Independent variable

None in the primary pass — this is an observational structural unit
against a fixed, already-established configuration. The contrast state
(`+$3A90` = `$CF` versus `$87`) provides a second point if the first pass
locates a candidate accumulator.

## Controlled variables

Same ROM, emulator, lab settings. Configuration fixed and unmodifiable
mid-battle (EXP-0042). No input supplied during the sampling intervals,
so no menu opens and no action is issued — the battle is left to run.

## Instrumentation

Two independent instruments on the same question, so their answers can be
required to converge:

1. **Structural.** `watchreads` (added in EXP-0042) over
   `WRAM:+$3A8F`–`+$3A90`, giving the reading PCs. The reader is then
   dumped from ROM and decoded to find what it writes.
2. **Differential.** Three savestates at equal frame intervals with no
   input, compared with `ff6lab state diff`. A gauge should appear as a
   cell advancing **monotonically by a similar amount** across both
   intervals, and — if DISC-0001's layout holds — as one member of a
   `$14`-stride ten-entry family.

Convergence requirement: the address the decoded reader writes must also
appear in the differential set. Either instrument alone is weaker —
the read watch names a routine but not necessarily a gauge, and the diff
names changing cells but not their owner.

## Expected outcomes

- *Direct:* one reading PC consumes `+$3A90`, and the routine adds it to
  (or subtracts it from) a per-slot accumulator that the differential
  pass independently shows advancing. Gauge storage and stride recorded.
- *Partial:* the reader is found but its target is not obviously a
  per-slot array — for example a single global counter — which would be
  a real finding about the architecture and would redirect the gauge
  question.
- *Refuted:* nothing reads `+$3A90` during a running battle, which would
  mean it is consumed only on a path this state does not exercise, and
  would send the unit to a scripted or freshly-entered battle instead.

## Falsifying outcome

1. Zero reads of `+$3A90` across a running ACTIVE-mode battle.
2. A reader exists but writes nothing that the differential pass sees
   advancing — the two instruments fail to converge, in which case
   neither result may be promoted.
3. Candidate gauge cells advance identically in the `$87` and `$CF`
   states, which would refute `+$3A90` having any rate role and reduce it
   to a value of unknown purpose.

## Evidence paths

`local_artifacts/experiments/EXP-0043/` — read tables, interval
savestates, diffs, ROM dumps of the decoded regions,
`bridge-commands.log`, `bridge-events.log`, `experiment.json`, and
`hashes.sha256` written at close, after which the directory is frozen.

## Trials

1. **Structural** — read-watch over `+$3A8F`–`+$3A90` on a free-running
   ACTIVE battle, ~3 500 frames, no input.
2. **Static** — ROM dumps of both reading regions, hand-decoded.
3. **Differential** — interval savestates and direct sampling of the
   candidate arrays.
4. **Contrast** — the same formation at Bat.Speed 3 (`$3A90` = `$CF`)
   against Bat.Speed 6 (`$3A90` = `$87`).

## Observations

### Both consumers, immediately

| Reading PC | Address | Count | Over |
|---|---|---|---|
| `$C21127` | `$7E3A8F` | 3 402 | 3 529 frames — **once per frame** |
| `$C209FD` | `$7E3A90` | 12 | 3 280 frames — sparse |

### `$3A8F`'s consumer is the per-frame scheduler gate

`read cpu C21110 128`, decoded from `$C2111B`:

```asm
C2111B  08           PHP
C2111C  E2 30        SEP #$30
C2111E  20 1F 4D     JSR $C24D1F
C21121  AD 41 2F     LDA $2F41
C21124  2D 8F 3A     AND $3A8F        ; the Wait flag
C21127  D0 62        BNE $C2118B      ; -> skip the entire update
C21129  AD 6C 3A     LDA $3A6C
C2112C  A2 02        LDX #$02
C2112E  C5 0E        CMP $0E
C21130  F0 5E        BEQ $C21190
C21132  1A / CA / D0 F8                ; small loop
C21136  EE 3E 3A     INC $3A3E        ; battle tick counter, low
C21139  D0 03        BNE $C2113E
C2113B  EE 3F 3A     INC $3A3F        ; ...high
C2113E  20 83 5A     JSR $C25A83
C21141  A2 12        LDX #$12         ; 18 = slot 9 x 2
C21143  EC E2 3E     CPX $3EE2
        ...
C2114B  BD A0 3A     LDA $3AA0,X      ; per-slot flags
C2114E  4A           LSR
C2114F  90 33        BCC $C21184
        ... status/gate tests ...
C21179  20 93 11     JSR $C21193      ; gauge advance
C2117C  BD 19 32     LDA $3219,X
C2117F  F0 03        BEQ $C21184
C21181  20 BB 11     JSR $C211BB      ; second accumulator
C21184  CA / CA      DEX DEX          ; stride 2
C21186  10 BB        BPL $C21143      ; loop over 10 slots
```

The whole per-frame battle update is gated on `$2F41 AND $3A8F`. `$3A8F`
is the Wait flag (EXP-0042); `$2F41` is zeroed at battle entry by
`$C22498` (also EXP-0042) and read `00` throughout this free-running
battle, so the gate never fired here.

`WRAM:+$3A3E` is a 16-bit counter incremented once per non-gated frame —
a **battle tick counter**. Measured `$2427` → `$24B8` = 145 increments
over 290 emulator frames, i.e. ≈0.5 per frame.

### The gauge advance

```asm
C21193  C2 20        REP #$20
C21195  BD C8 3A     LDA $3AC8,X      ; per-slot increment
C21198  4A           LSR              ; >> 1
C21199  18           CLC
C2119A  7D B4 3A     ADC $3AB4,X      ; += gauge
C2119D  9D B4 3A     STA $3AB4,X
C211A0  E2 20        SEP #$20
C211A2  B0 06        BCS $C211AA
C211A4  EB           XBA
C211A5  DD 2C 32     CMP $322C,X
C211A8  90 10        BCC $C211BA
C211AA  A9 FF        LDA #$FF
C211AC  9D 2C 32     STA $322C,X
C211AF  20 77 4E     JSR $C24E77
C211B2  A9 20        LDA #$20
C211B4  1D A0 3A     ORA $3AA0,X
C211B7  9D A0 3A     STA $3AA0,X
C211BA  60           RTS
```

and a second accumulator on a different gate:

```asm
C211BB  C2 21        REP #$21
C211BD  BD 18 32     LDA $3218,X
C211C0  7D C8 3A     ADC $3AC8,X
C211C3  9D 18 32     STA $3218,X
```

### The arrays, sampled live

All are **10 entries, stride 2, 16-bit**, with party in slots 0–3 and
enemies in 4–9 — the DISC-0001 slot convention, but at stride 2 rather
than the `$14` stride of the HP/stat family.

```
+$3AC8 (increment): 3E 01  4A 01  50 01  0000 0000 0000  9C 00  9C 00  9C 00  0000
+$3AB4 (gauge):     0000   0000   0000   0000 0000 0000  9A 00  5A 00  1E 00  0000
+$3AA0 (flags):     8F 00  8F 00  8F 00  0000 0000 0000  01 01  01 01  01 01  0000
+$3B18 (source):    04 21  03 23  03 24  0000 0000 0000  04 1E  04 1E  04 1E  0000
```

Formation 14 places its three enemies in slots **6–8**.

`+$3AB4` advanced by exactly `$4E` = 78 = `$9C >> 1` per increment, the
predicted `$3AC8,X >> 1`. Slot 8 was caught wrapping `$00B6` → `$0004`,
and separately caught advancing `$001E` → `$006C` in the same sample
where its `+$3AA0` low byte went `$01` → `$41` and back to `$01` at the
next sample.

Party gauges read `0000` and did not move: with no input supplied they
had become ready early and parked awaiting a command, which is also why
the party was eventually annihilated. The defeat screen continues to
increment `+$3A3E`.

### The contrast — where `$3A90` actually lands

`$C209FA LDA $3A90` sits inside the routine that computes `+$3AC8,X`,
behind a branch:

```asm
C209F0  BD 19 3B     LDA $3B19,X      ; per-slot Speed
C209F3  69 14        ADC #$14         ; + 20
C209F5  EB           XBA
C209F6  E0 08        CPX #$08
C209F8  90 06        BCC $C20A00      ; X < 8 (party) -> skip
C209FA  AD 90 3A     LDA $3A90        ; enemies only
C209FD  20 81 47     JSR $C24781      ; multiply
C20A00  68           PLA
C20A01  20 81 47     JSR $C24781
C20A04  C2 20        REP #$20
C20A06  4A 4A 4A 4A  LSR x4
C20A0A  9D C8 3A     STA $3AC8,X
```

X is the stride-2 slot index, so `X < 8` is slots 0–3 — the party.
Measured against the same formation:

| Slot | Bat.Speed 6 (`$3A90` = `$87`) | Bat.Speed 3 (`$3A90` = `$CF`) |
|---|---|---|
| 0–2 (party) | 318 / 330 / 336 | **318 / 330 / 336** — identical |
| 6–8 (enemy) | 156 / 156 / 156 | **240 / 240 / 240** |

The per-slot Speed source `+$3B19,X` reads 33 / 35 / 36 for the party and
30 for each enemy. The party increments fit `(Speed + 20) × 6` exactly at
all three slots.

## Interpretation

**The ATB layer is located.**

- `WRAM:+$3AB4` — the **ATB gauge**, 10 entries, stride 2, 16-bit,
  advanced by `ROMCPU:$C21195`–`$C2119D`.
- `WRAM:+$3AC8` — the **per-slot increment**, computed near
  `ROMCPU:$C209E0` from the slot's Speed at `+$3B19,X`.
- `WRAM:+$3AA0` — the **per-slot flag word** the scheduler tests before
  advancing a slot, and writes on the threshold path.
- `WRAM:+$3A3E` — a **battle tick counter**, one per non-gated frame.
- `ROMCPU:$C21124` — the **gate on the entire per-frame battle update**,
  `$2F41 AND $3A8F`. Since `$3A8F` is the Wait flag, this is where
  ACTIVE/WAIT is enforced, with `$2F41` as the partner condition.

**Battle Speed scales enemy gauges only.** The `CPX #$08 / BCC` branch
skips the `$3A90` multiply for party slots, and the measurement confirms
it: party increments were byte-identical at two Battle Speeds while enemy
increments changed by a factor of ~1.54, tracking `$3A90`'s ratio
207/135 ≈ 1.53. In this ROM revision Battle Speed is not a global pace
control; it changes how fast enemies fill relative to the party.

This also explains EXP-0042's sparse read count. `$3A90` is consulted
only when the increment table is recomputed — 12 reads for three enemies
is four recomputations — not per tick.

### Contextual observations, not conclusions

- `$2F41` is the other half of the pause gate and read `00` throughout.
  It is the obvious candidate for the menu-state flag that makes WAIT
  actually wait, but **nothing here tests that** — no menu was opened.
  This is the next unit's question.
- The tick counter advanced ≈0.5× per emulator frame in the one window
  measured. That is adjacent to open question #18 (the sub-1.0 per-frame
  dispatch rate) but was not the object of study, and one window is not a
  rate measurement.
- `+$3AA0` bit 6 was seen set transiently at a gauge advance, while the
  decoded threshold path at `$C211B2` ORs `#$20` (bit 5). Both are
  observations; reconciling them needs finer sampling than this unit ran.
- `(Speed + 20) × 6` fits all three party increments, but three points
  from one party do not establish a formula.

## Alternatives

- **The `$14`-stride assumption did not hold.** The ATB family is
  **stride 2**, not the `$14` stride of the HP/stat arrays, so DISC-0001's
  layout governs slot *assignment* but not stride. This unit's evidence
  is direct, but it means "unified battle layout" must not be applied to
  new arrays without checking.
- **`+$3AB4` could be one of several accumulators.** `$3218,X` receives
  the full `$3AC8,X` under a different gate at `$C211BB`, and `$3219,X`
  gates it. Which is "the" ATB gauge and which is a charge or delay timer
  is not settled here; `+$3AB4` is the one whose overflow reaches the
  flag write.
- **Party gauges were never observed climbing**, only sitting at zero
  while ready. The increment path is shared code and the party increments
  are populated, so the inference is sound, but the direct observation
  covers enemy slots only.
- **Two Battle Speed values, one formation, one party.** The asymmetry is
  supported structurally by the branch as well as numerically, which is
  what raises it above correlation — but the middle speeds are untested.
- `+$3B19,X` reading 33/35/36/30 is consistent with a Speed stat and is
  used as one, but this unit did not verify it against any other source.

## Result

**The primary question is answered.** `WRAM:+$3A90` is consumed at
`ROMCPU:$C209FA`, inside the routine that builds the per-slot ATB
increment table `WRAM:+$3AC8` — and only for enemy slots. The gauges are
`WRAM:+$3AB4`, advanced by `$3AC8,X >> 1` at `ROMCPU:$C21195`.

The secondary question is answered too, and a third came free:
`WRAM:+$3A8F`'s consumer is `ROMCPU:$C21124`, the gate on the whole
per-frame battle scheduler — the mechanism by which ACTIVE and WAIT
differ.

All three instruments converged: the read watch named the routines, the
static decode named the arrays, and direct sampling showed `+$3AB4`
advancing by exactly the predicted amount. No falsifier fired.

## Confidence

- `WRAM:+$3AB4` is the ATB gauge array, 10 entries stride 2, advanced by
  `$3AC8,X >> 1` at `ROMCPU:$C21195`–`$C2119D`: **Confirmed** (decoded
  statically, observed advancing by exactly `$4E` = `$9C >> 1`, wrap
  observed).
- `WRAM:+$3AC8` is the per-slot increment array: **Confirmed**.
- `WRAM:+$3A90` consumed at `ROMCPU:$C209FA` for enemy slots only:
  **Confirmed** (branch decoded; party increments byte-identical and
  enemy increments changed across two Battle Speeds).
- Battle Speed scales enemy ATB and not party ATB: **Confirmed for the
  tested formation and the two tested speeds**; middle speeds, other
  formations and other party compositions untested.
- `ROMCPU:$C21124` gates the per-frame battle update on
  `$2F41 AND $3A8F`: **Confirmed** (code role).
- `WRAM:+$3A3E` is a battle tick counter: **Confirmed** (code role and
  observed incrementing).
- `WRAM:+$3AA0` is the per-slot scheduler flag word: **Confirmed** (code
  role); individual bit semantics **Unknown**.
- `+$3B19,X` as the per-slot Speed stat: **Strong hypothesis**.
- `$2F41` as the menu/pause partner: **Tentative hypothesis** — untested.
- The exact increment formula, threshold, overflow and reset behaviour:
  **Unknown**, deliberately out of scope.

## Stopping condition

The unit ends when `+$3A90`'s consumer is identified and its target
either located or explicitly bounded.

Out of scope, each belonging to a later unit:

- The full increment formula, overflow and threshold behaviour.
- The ACTIVE/WAIT pause matrix and which menu states gate which clock.
- Haste, Slow, Stop and other status modifiers.
- The action queue, readiness arbitration and execution lifecycle.
- Any Whelk work.

Broad capture is permitted and encouraged; anything outside the primary
question is recorded as a contextual observation, not a conclusion.

## Next action

**EXP-0044 — the ACTIVE/WAIT pause matrix.** `ROMCPU:$C21124` gates the
entire per-frame battle update on `$2F41 AND $3A8F`. `$3A8F` is known;
`$2F41` is not. Find what sets and clears `$2F41` — a write-watch across
opening and closing each battle submenu — and then build the matrix of
menu/presentation states against timer domains, reading the domains
directly now that they are named: `+$3AB4` (gauges), `+$3A3E` (tick
counter), `+$3AA0` (flags), `+$3218`.

This is the unit the master ATB program has been aiming at, and it is now
cheap: every quantity in the matrix has an address. It is also the first
unit that genuinely warrants `/battle-baseline` and parallel read-only
observers over a frozen package.

Note the staging rule from EXP-0042 still applies for the *configuration*
variable, but `$3A8F` and `$3A90` can now be patched directly mid-battle,
which makes ACTIVE-versus-WAIT a one-variable comparison inside a single
savestate lineage rather than two separate route runs.

Deferred, each with a known entry point:

- The exact increment formula, threshold, overflow and gauge reset.
- `+$3AA0` bit semantics, and reconciling the observed bit 6 with the
  decoded bit 5 write at `$C211B2`.
- `$3218`/`$3219` — the second accumulator and its gate.
- `$322C,X` — the comparand on the threshold path.
- Verifying `+$3B19,X` as Speed against an independent source.
- Battle Speed's middle values, and other formations and parties.
- Whether the ≈0.5/frame tick rate relates to open question #18.
