# Discovered Functions

### CopyCharacterFields

- **Address:** `$C10DF3`
- **Address space:** ROM CPU
- **Kind:** Function
- **Status:** Confirmed (behavior; full byte-exact disassembly, verified
  against live memory). Higher-level *purpose* of the copied buffer remains
  a strong hypothesis (battle display/staging data).
- **Previous names:** `CopyCharacterBattleData` (Session 001);
  "candidate CopyCharacterFields" (Session 001 promotion). Session 002's
  full disassembly confirmed the copy behavior and extended it from four
  copied fields to six plus a mask output.
- **Observed behavior:** Full disassembly (Session 002, dumped via Lua
  bridge, every branch offset and observed address verified byte-exact):

  ```asm
  ; entry: m=1 (8-bit A), x=0 (16-bit X/Y), DB=$7E, D=$0000, native mode
  C10DF3  TDC             ; A = 0
  C10DF4  TAX             ; X = 0
  C10DF5  DEX             ; X = $FFFF
  C10DF6  STX $10         ; mask08 = $FFFF   (16-bit store, x=0)
  C10DF8  STX $12         ; mask0A = $FFFF
  C10DFA  LDA $628D       ; flag byte 1 ($7E628D)
  C10DFD  BNE $C10E0C     ; nonzero -> never mask
  C10DFF  LDA $E9EF       ; flag byte 2 ($7EE9EF)
  C10E02  BEQ $C10E0C     ; zero -> no masking
  C10E04  INX             ; X = 0
  C10E05  STX $12         ; mask0A = $0000
  C10E07  LDX #$0038
  C10E0A  STX $10         ; mask08 = $0038
  C10E0C  REP #$20        ; m=0: 16-bit A (explains m=0 captured at $C10E14)
  C10E0E  TDC
  C10E0F  TAX             ; X = 0
  C10E10  TAY             ; Y = 0
  ; ---- copy loop: six 16-bit fields per slot, four slots ----
  C10E11  LDA $2E78,X
  C10E14  STA $2EB5,Y     ; +0 current HP (Session 001 write breakpoint)
          LDA $2E80,X
          STA $2EB7,Y     ; +2
          LDA $2E88,X
          STA $2EB9,Y     ; +4
          LDA $2E90,X
          STA $2EBB,Y     ; +6
          LDA $2E98,X
          AND $10
          STA $2EBD,Y     ; +8, masked
          LDA $2EA0,X
          AND $12
          STA $2EBF,Y     ; +$A, masked
          INX
          INX             ; X += 2
          TYA
          CLC
          ADC #$0020
          TAY             ; Y += $20
          CPX #$0008
          BNE $C10E11
  ; ---- mask loop: derive 4-bit per-slot mask from source +$A high bytes ----
  C10E46  TDC
          SEP #$20        ; m=1: 8-bit A
          STZ $10
          TDC
          TAX             ; X = 0
  C10E4D  LDA $2EA1,X     ; high byte of $2EA0 entry
          AND #$20        ; isolate bit 5 (bit 13 of the 16-bit value)
          EOR #$20        ; invert
          LSR
          ORA $10
          LSR
          STA $10         ; accumulate: slot n -> bit n
          INX
          INX
          CPX #$0008
          BNE $C10E4D
  C10E61  LDA $10
  C10E63  STA $61AD       ; publish 4-bit slot mask ($7E61AD)
  C10E66  RTS
  ```

- **Evidence:**
  - Session 001: converging memory search on `$2EB5`; write breakpoint at
    `$C10E14`.
  - Session 002: ROM byte dump `$C10DF3–$C10E66`; exec-callback captures at
    entry (`PS=$26`: m=1, x=0; `DB=$7E`, `D=$0000`, native) and at the store
    (`PS=$04`: m=0, x=0; iterations `X=0/Y=0`, `X=2/Y=$20`, `X=4/Y=$40` with
    A = 53, 39, 33 — the three displayed party HP values);
    live WRAM reads matching every copied field and `$61AD = $0F` matching
    the mask formula for an all-zero `$2EA0` array.
- **Interpretation:** Per-frame refresh of a party-slot record buffer from
  six parallel source arrays, with conditional masking of the last two
  fields, plus derivation of a per-slot bitmask (bit n set when bit 13 of
  the slot's `$2EA0` field is clear).
- **Inputs:** flag bytes `$628D`, `$E9EF`; source arrays `$2E78–$2EA7`;
  DB must be `$7E`.
- **Outputs / modified memory:** records at `$2EB5+n*$20` offsets `+0..+$B`;
  `$61AD`; direct-page scratch `$10`, `$12`; registers A, X, Y, flags.
- **Callers:** `JSR $0DF3` at `$C10200`, inside PerFrameBattleUpdate
  (`$C101FB`, below) — the steady-state caller during battle. **Second
  caller (EXP-0002):** a one-shot `JSR` at `ROMCPU:$C11090` (return
  `$C1:1093`) fired once at battle entry, before the per-frame stream
  began; context unexplored. Rate note (EXP-0002): the per-frame stream
  runs at up to ~1/frame when the battle is settled but drops during
  entry/heavy phases (≈0.23–0.7/frame) — "every frame" was an
  overgeneralization. Never observed at the title screen or in field/event
  contexts (≈175k non-battle frames, zero fires).
- **Callees:** none.
- **Alternative explanations:** The destination buffer's consumer is
  unidentified; "display/staging data" is inferred from the per-frame
  refresh pattern, not proven.
- **Validation experiment:** Trigger the masked mode (find when
  `$628D`/`$E9EF` become nonzero) and verify fields `+8`/`+$A` get masked
  in live memory.
- **Go representation:** `chardata.CopyCharacterFields`
  ([internal/game/chardata/chardata.go](../../internal/game/chardata/chardata.go)) — includes the mask
  modes and returns the `$61AD` slot mask.
- **Related discoveries:** PerFrameBattleUpdate,
  [05_DATA_STRUCTURES.md](05_DATA_STRUCTURES.md),
  [03_DISCOVERED_VARIABLES.md](03_DISCOVERED_VARIABLES.md)
- **First observed in session:** [SESSION_001](SESSION_001.md)
- **Last updated:** 2026-07-29 ([SESSION_002](SESSION_002.md))

### Battle HP/MP delta engine (candidate) — `ROMCPU:$C21323` / `$C21350` / `$C21390`

- **Address:** dispatch wrapper `ROMCPU:$C212F5`–`$C2131E`, table
  `ROMCPU:$C2131F` (= `$1323, $1350`), HP routine `ROMCPU:$C21323`, MP
  routine `ROMCPU:$C21350`, death handler `ROMCPU:$C21390`, delta fetch
  `ROMCPU:$C213A7`–`$C213D2`
- **Address space:** ROM CPU
- **Kind:** Function cluster
- **Status:** **Confirmed (code, byte-exact)** —
  [EXP-0001](../experiments/EXP-0001-c2-delta-engine-dump.md) dumped
  `ROMCPU:$C212F0`–`$C2141F` and every instruction below was verified
  arithmetically (branch offsets, operands). Live store captures
  (Session 003) match the stores at `$C21338`/`$C21347`/`$C21396`
  (+3 post-instruction callback PCs, now Confirmed). **Semantic labels**
  (HP vs MP, max, "death") remain Strong hypothesis / per-field — see
  [05_DATA_STRUCTURES.md](05_DATA_STRUCTURES.md).
- **Previous names:** `ApplyHPDelta`/`ApplyMPDelta` (battle.go,
  Session 003); "claimed" entry in this file between Unit 1 and EXP-0001.
- **Observed behavior:** Verified disassembly (registers per live
  captures: 16-bit A, 8-bit X/Y, `DB=$7E`, `Y = slot×2`):

  ```asm
  C212F5  PHX
  C212F6  PHP
  C212F7  LDX #$02        ; default: table entry 1 (MP)
  C212F9  LDA $11A2       ; selector byte: bit 7 set -> MP path
  C212FC  BMI $C21300
  C212FE  LDX #$00        ; entry 0 (HP)
  C21300  JSR ($131F,X)   ; pushed return $1302 = Session 003 stacks
  C21303  SEP #$20
  C21305  BCC $C2131C     ; C clear on return -> skip tail
  C21307  LDA $02,S       ; ---- post-dispatch tail: purpose unknown ----
  C21309  TAX
  C2130A  STX $EE
  C2130C  JSR $362F       ; $C2362F - unexplored
  C2130F  CPY $EE
  C21311  BEQ $C2131C
  C21313  STA $327C,Y
  C21316  LDA $3018,Y
  C21319  TRB $3419
  C2131C  PLP
  C2131D  PLX
  C2131E  RTS
  C2131F  .dw $1323,$1350 ; dispatch table
  ; ---- HP delta ----
  C21323  JSR $13A7       ; fetch delta; Z = zero, C = non-negative
  C21326  BEQ $C2133B     ; zero -> CLC/RTS
  C21328  BCC $C2133D     ; negative -> damage
  C2132A  CLC
  C2132B  ADC $3BF4,Y     ; heal: delta + current
  C2132E  BCS $C21335     ; 16-bit overflow -> clamp
  C21330  CMP $3C1C,Y
  C21333  BCC $C21338     ; sum < max -> keep sum
  C21335  LDA $3C1C,Y     ; clamp to max
  C21338  STA $3BF4,Y     ; heal store (capture pc $C2133B)
  C2133B  CLC
  C2133C  RTS
  C2133D  EOR #$FFFF      ; damage: with C=0, SBC below adds the +1
  C21340  STA $EE
  C21342  LDA $3BF4,Y
  C21345  SBC $EE         ; current - magnitude
  C21347  STA $3BF4,Y     ; damage store (capture pc $C2134A)
  C2134A  BEQ $C21390     ; exactly zero -> death handler
  C2134C  BCS $C2133C     ; positive -> RTS
  C2134E  BRA $C21390     ; underflow -> death handler
  ; ---- MP delta (same shape over $3C08/$3C30) ----
  C21350  JSR $13A7
  C21353  BEQ $C2133B
  C21355  BCC $C2136B
  C21357  CLC
  C21358  ADC $3C08,Y
  C2135B  BCS $C21362
  C2135D  CMP $3C30,Y
  C21360  BCC $C21365
  C21362  LDA $3C30,Y
  C21365  STA $3C08,Y
  C21368  CLC
  C21369  BRA $C2138A     ; -> exit tail
  C2136B  EOR #$FFFF
  C2136E  STA $EE
  C21370  LDA $3C08,Y
  C21373  SBC $EE
  C21375  STA $3C08,Y
  C21378  BEQ $C2137C
  C2137A  BCS $C2138A
  C2137C  TDC
  C2137D  STA $3C08,Y     ; zero-floor MP
  C21380  LDA $3C95,Y
  C21383  LSR             ; bit 0 -> carry
  C21384  BCC $C21389
  C21386  JSR $1390       ; dies-at-zero-MP -> death handler
  C21389  SEC
  C2138A  LDA #$0080      ; MP exit tail; purpose unknown
  C2138D  JMP $464C       ; $C2464C - unexplored
  ; ---- death handler ----
  C21390  SEC
  C21391  TDC
  C21392  TAX
  C21393  STX $3A89       ; clear $3A89 (8-bit X)
  C21396  STA $3BF4,Y     ; zero HP (capture pc $C21399)
  C21399  LDA $3EE4,Y
  C2139C  BIT #$0002
  C2139F  BNE $C2133C     ; bit 1 set -> suppress, RTS
  C213A1  LDA #$0080
  C213A4  JMP $0E32       ; $C20E32 - death event, unexplored
  ; ---- delta fetch ----
  C213A7  LDA $33D0,Y     ; secondary pending delta
  C213AA  INC
  C213AB  BEQ $C213BC     ; $FFFF sentinel -> 0
  C213AD  LDA $3018,Y
  C213B0  BIT $3A3C
  C213B3  BEQ $C213B9
  C213B5  TDC
  C213B6  STA $33D0,Y     ; gate hit -> cancel secondary
  C213B9  LDA $33D0,Y
  C213BC  STA $EE
  C213BE  LDA $3A81       ; 16-bit, overlapping reads as written
  C213C1  AND $3A82
  C213C4  BMI $C213C8
  C213C6  STZ $EE         ; gate clear -> drop secondary
  C213C8  LDA $33E4,Y     ; primary pending delta
  C213CB  INC
  C213CC  BEQ $C213CF     ; $FFFF sentinel -> 0
  C213CE  DEC
  C213CF  SEC
  C213D0  SBC $EE         ; delta = primary - secondary
  C213D2  RTS
  ```

- **Evidence:** EXP-0001 dump `mesen/out/rom_C212F0_304.hex` (SHA-256
  `2800f34b…d5d56a`); Session 003 live write captures
  (`mesen/out/session003/events.log`, SHA-256 `bcfc7f4c…a99d03`).
- **Inputs:** `WRAM:+$11A2` (bit 7 selector), pending-delta arrays
  `+$33E4`/`+$33D0` (per-slot, `$FFFF` = none), gates `+$3A3C`,
  `+$3A81`/`+$3A82`, `+$3018,Y`; arrays `+$3BF4`/`+$3C1C` (HP path),
  `+$3C08`/`+$3C30` (MP path), `+$3C95` bit 0, `+$3EE4` bit 1.
- **Outputs / modified memory:** `+$3BF4,Y` or `+$3C08,Y`; `+$3A89`
  cleared and `JMP $C20E32` (A=`$0080`) on unsuppressed death;
  `+$33D0,Y` cancelled under the `$3A3C` gate; tail may write
  `+$327C,Y` and `TRB $3419`; scratch `$EE`.
- **Alternative explanations:** none remaining for the code; semantic
  labels rest on Session 002/003 battle context and parallel structure.
- **Slot universality (EXP-0005, 2026-07-29):** the engine is
  slot-uniform — its damage store (`$C21347`) and death-handler zero
  (`$C21396`) were observed writing **enemy** entries (indexes 4–5 of the
  10-entry arrays) when the party attacked, with on-screen defeat
  matching the zeroed entries. Party and enemies share one engine.
- **Validation experiment:** live MP observation (spend/heal MP watching
  `+$3C08`); find `$11A2` writers; explore `$C20E32` and the `$464C`
  tail.
- **Go representation:**
  [internal/game/battle/battle.go](../../internal/game/battle/battle.go) —
  arithmetic now Confirmed byte-exact; delta fetch, `$3A89`, `$C20E32`,
  and the exit tail deliberately unmodeled.
- **Related discoveries:** CopyCharacterFields (downstream display copy),
  battle-array lifecycle writers ([04_MEMORY_MAP.md](04_MEMORY_MAP.md))
- **First observed in session:** [SESSION_003](SESSION_003.md)
- **Last updated:** 2026-07-29 (EXP-0001)

### PendingDeltaAccumulate (candidate) — `ROMCPU:$C20C76`

- **Address:** `ROMCPU:$C20C76`–`$C20C9D` (store `$C20C98`); observed
  caller path `JSR $0C2D` at `ROMCPU:$C20C28` → `$C20C2D` gate block
- **Address space:** ROM CPU
- **Kind:** Function
- **Status:** **Confirmed (code, byte-exact — EXP-0006)**; polarity/
  gate semantics and the DP `$F0` amount provenance Unknown.
- **Observed behavior:** adds the amount in DP `$F0` to the slot's
  pending delta (`+$33D0,Y`, or `+$33E4,Y` when the prelude's
  `ROL/EOR $F2/LSR` carry retargets `Y += $14`), treating the `$FFFF`
  sentinel as 0, **clamping the accumulated total at `$270F` (9999)** —
  the damage-number cap emerges here. Full listing in
  [EXP-0006](../experiments/EXP-0006-delta-setter.md).
- **Evidence:** EXP-0006 dumps (SHA-256s in the record); live store
  captures (`$C20C9B` ×12, EXP-0004) with matching PHP/Y/return stack
  bytes.
- **Interpretation:** the queueing API of the battle formula layer:
  damage and healing are accumulated per slot into opposing pending
  arrays, consumed later by the delta-engine fetch (`$C213A7`,
  delta = `+$33E4` − `+$33D0`).
- **Amount semantics (EXP-0007, Confirmed):** DP `$F0` holds the final
  per-hit amount — it equals both the applied HP delta and the on-screen
  damage popup (three-anchor correlation: array arithmetic, HUD values,
  captured popup "6"). At entry, `X` = attacker slot×2 and `Y` = target
  slot×2 (strong hypothesis, 5/5 captures); DP `$F2` read `$20` on a
  party→enemy hit vs `$00` on enemy→party (meaning unresolved).
- **Alternatives:** which array is "damage" vs "heal" is inferred from
  the fetch's subtraction direction; unverified live (a Heal Force cast
  with this watch would settle it).
- **Validation experiment:** heal-cast logging (polarity); dump around
  `ROMCPU:$C26B10` (stack hint for the formula frame).
- **Go representation:** `battle.AccumulatePending`
  ([internal/game/battle/battle.go](../../internal/game/battle/battle.go)).
- **Related discoveries:** Battle HP/MP delta engine (consumer),
  EXP-0004 writer census.
- **First observed in session:** EXP-0004 (write PC); decoded 2026-07-30
  (EXP-0006).
- **Last updated:** 2026-07-30

### DamageAmountPipeline / elemental-modifier block (candidate) — `ROMCPU:$C20B83`

- **Address:** `ROMCPU:$C20B83`–`$C20C2C` (joins the gate block and
  accumulator); base-amount callees `ROMCPU:$C20C9E` / `$C20D87`
  (selected by `+$11A4` bit 7) undumped
- **Address space:** ROM CPU
- **Kind:** Function (per-target amount post-processing)
- **Status:** **Confirmed (code, byte-exact — EXP-0010 + EXP-0006 dumps
  joined)**; semantic labels Strong hypothesis; base formulas Unknown
  (question #22).
- **Observed behavior:** full verified listing in
  [EXP-0010](../experiments/EXP-0010-formula-body.md). Shape: `$11A6`
  gate → base-amount JSR (variant by `+$11A4` bit 7) → polarity byte
  DP `$F2` from `+$11A4` with status-driven flips (`+$3EE4` bits,
  `+$3C95` bit 7, `+$11AA`) → element block vs `+$11A1`:
  battle-wide nullify (`~+$3EC8 & +$11A1 == 0` → zero), then per-target
  first-match-wins `+$3BCC,Y` flip-to-heal / `+$3BCD,Y` zero /
  `+$3BE1,Y` halve / `+$3BE0,Y` double (with `$8000` overflow guard) →
  falls into the queue gate block.
- **Evidence:** `rom_C20B83_141.hex` (SHA-256 `3a9034d3…c0a39f`) +
  `rom_C20C10_48.hex`; entry `PHP` matches the live stack PS byte;
  called from the `$C23469` target loop (EXP-0009).
- **Interpretation:** the elemental response system: per-slot 16-bit
  mask entries at `+$3BCC` (flip|zero packed) and `+$3BE0`
  (double|halve packed) — two more `$14`-stride 10-entry family arrays.
  Candidate labels: absorb/immune/weak/resist; `+$11A1` = attack
  element byte; `+$3EC8` = battle-wide element-nullify byte.
- **Alternatives:** labels could permute (e.g. `$3BCD` "immune" vs
  another negation concept); discriminable live with known
  absorb/resist targets.
- **Validation experiment:** fight an element-absorbing target and
  watch the flip; dump `$C20C9E`/`$C20D87` (question #22).
- **Go representation:** `battle.ApplyElementResponse` +
  `battle.ElementResponse` (behavior-derived names per naming rules).
- **Related discoveries:** PendingDeltaAccumulate (downstream),
  `$C23469` target loop (upstream), unified battle arrays (ST-0003).
- **First observed in session:** 2026-07-30 (EXP-0010).
- **Last updated:** 2026-07-30

### PartyDisplaySourceRefresh (candidate) — `ROMCPU:$C25D26`

- **Address:** `ROMCPU:$C25D26`–`$C25D56`
- **Address space:** ROM CPU
- **Kind:** Function
- **Status:** **Confirmed (code, byte-exact — EXP-0003)**; caller and
  trigger Unknown; "display source" purpose Strong hypothesis (downstream
  consumer is CopyCharacterFields, whose own consumer is unidentified).
- **Previous names:** none (writer PC `$C25D33` in Session 003 /
  EXP-0002 captures).
- **Observed behavior:**

  ```asm
  C25D26  PHP
  C25D27  REP #$20        ; 16-bit A
  C25D29  SEP #$10        ; 8-bit X/Y
  C25D2B  LDY #$06        ; slots 3..0, Y = slot*2
  C25D2D  LDA $3BF4,Y
  C25D30  STA $2E78,Y     ; current HP (capture pc $C25D33)
  C25D33  LDA $3C1C,Y
  C25D36  STA $2E80,Y     ; heal-clamp ceiling -> display max HP
  C25D39  LDA $3C08,Y
  C25D3C  STA $2E88,Y     ; MP-path pool -> display
  C25D3F  LDA $3C30,Y
  C25D42  STA $2E90,Y     ; MP-path ceiling -> display
  C25D45  LDA $3EE4,Y
  C25D48  STA $2E98,Y     ; status word -> display (masked $0038 downstream)
  C25D4B  LDA $3EF8,Y
  C25D4E  STA $2EA0,Y     ; -> display; bit 13 drives $61AD downstream
  C25D51  DEY
  C25D52  DEY
  C25D53  BPL $C25D2D
  C25D55  PLP
  C25D56  RTS
  ```

- **Evidence:** EXP-0003 dump `mesen/out/rom_C25CC0_176.hex` (SHA-256
  `76861e12…a5dc7`); live store captures at `$C25D33` across Session 003
  and both EXP-0002 battles; EXP-0002 `writers` counts (168 watched
  writes ≈ 42 invocations over ~6600 battle frames → event-driven).
- **Interpretation:** the missing link of open question #1 — every
  CharacterFieldsSource array is a copy of an authoritative battle
  array. Full chain: delta engine → `WRAM:+$3BF4` family →
  this copier → `WRAM:+$2E78` family → CopyCharacterFields →
  `WRAM:+$2EB5` records + `+$61AD`.
- **Alternative explanations:** none for the copy itself; invocation
  policy (event-driven vs sampled) rests on write-count arithmetic only.
- **Validation experiment:** exec watch at `$C25D26` with stack capture
  → caller and trigger events.
- **Go representation:** none yet (pure copy; model after caller known).
- **Related discoveries:** Battle HP/MP delta engine (upstream),
  CopyCharacterFields (downstream), CharacterFieldsSource.
- **First observed in session:** [SESSION_003](SESSION_003.md) (store
  PC); decoded 2026-07-29 (EXP-0003).
- **Last updated:** 2026-07-29

### PerFrameBattleUpdate (candidate)

- **Address:** `$C101FB`
- **Address space:** ROM CPU
- **Kind:** Function
- **Status:** Strong hypothesis (per-frame battle update dispatcher)
- **Previous names:** None
- **Observed behavior:** Disassembly from ROM dump (Session 002):

  ```asm
  C101FB  PHX
  C101FC  PHY
  C101FD  JSR $1A24       ; $C11A24 — unexplored
  C10200  JSR $0DF3       ; CopyCharacterFields
  C10203  JSR $4504       ; $C14504 — unexplored
  C10206  JSR $2F79       ; $C12F79 — unexplored
  C10209  JSR $02CA       ; $C102CA — unexplored
  C1020C  JSR $44BE       ; $C144BE — unexplored
  C1020F  JSL $C2BF53     ; bank $C2 — unexplored
  C10213  JSR $93E3       ; $C193E3 — unexplored
  C10216  JSL $C2B41A     ; bank $C2 — unexplored
  C1021A  LDA $E9EF       ; same flag byte CopyCharacterFields tests
  C1021D  BNE $C10223
  C1021F  JSL $C20003     ; only when $E9EF == 0
  C10223  PLY
  C10224  PLX
  C10225  RTL
  ```

- **Evidence:** ROM dump `$C101C0–$C1022F`; return address `$C1:0203` in
  every steady-state exec capture at `$C10DF3`; stack snapshot showing
  pushed X/Y then a long-call return address `$C2:6429` (so a
  `JSL $C101FB` at `$C26425`); EXP-0002: zero fires across title and
  three field/event contexts (≈175k frames), fires within ~30 frames of a
  random-encounter entry and continuously thereafter (≈0.23–0.87/frame
  by phase).
- **Interpretation:** A fixed sequence of per-frame battle subsystem
  updates, entered long from bank `$C2` each frame. CopyCharacterFields is
  its second step.
- **Alternative explanations:** Could run in non-battle contexts not yet
  observed (menus, events); "battle" is inferred from the single observed
  context.
- **Validation experiment:** Exec-log `$C101FB` on the world map / in the
  menu to test whether it is battle-only; exec-log `$C26425` to find the
  bank-`$C2` frame loop.
- **Go representation:** none yet (callees unexplored).
- **Related discoveries:** CopyCharacterFields
- **First observed in session:** [SESSION_002](SESSION_002.md)
- **Last updated:** 2026-07-29
