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
  (`$C101FB`, below). Observed firing once per frame during battle; never
  at the title screen. No other callers observed yet.
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
  every exec capture at `$C10DF3`; stack snapshot showing pushed X/Y then a
  long-call return address `$C2:6429` (so a `JSL $C101FB` at `$C26425`);
  fires once per frame in battle, never at the title screen.
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
