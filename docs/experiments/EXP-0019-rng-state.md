# EXP-0019: Decode the action clear/populate routine; find the RNG state

- **Status:** running (2026-07-30)
- **Question (#29 endgame):** The hit/act decision sits between the
  action-block clear (`ROMCPU:$C2297D`) and the power populate
  (`$C229D4`), reached via a `JSR` at `~$C2319D`. Decode that stretch,
  enumerate its reads, and identify the timing-varying one — the
  battle RNG state.
- **Starting state:** deterministic ROM reads; live verification via
  `exp10-battle.mss` trials afterward if a candidate emerges.
- **Observation method:** bridge `read cpu C22950 160` and
  `read cpu C23190 48`; artifacts `rom_C22950_160.hex`,
  `rom_C23190_48.hex` + SHA-256s; hand-disassembly anchored at the two
  verified write sites (`$C2297D`, `$C229D4` via the +3 rule:
  instruction ends at −3).
- **Expected outcomes:** a conditional between clear and populate whose
  operand is a WRAM state cell (RNG candidate) or a call to a known
  random helper; *alternative:* the gate is in the `$C2319D` caller.
- **Falsifying outcome:** no conditional path between clear and
  populate in the dumped window (bracket error).
- **Raw evidence paths:** the two `.hex` artifacts.
- **Result:** (artifacts `rom_C22950_160.hex` SHA-256 `9c411a90…f503d7`,
  `rom_C23190_48.hex` `b261e1a0…141acf`) — **the "clear" writer is an
  MVN table load; the action block is ROM data.**
  - `$C22966` routine: `XBA; LDA #$0E; JSR $4781` (index × 14);
    `REP #$31; ADC #$6AC0; TAX; LDY #$11A0; LDA #$000D;`
    **`MVN $C4,$7E`** (`$C2297A`, callback PC `$C2297D` ✓) — copies a
    **14-byte record from `ROMCPU:$C46AC0 + 14×index`
    (`ROMFILE:0x046AC0+`) into `WRAM:+$11A0`–`+$11AD`**; then
    `ASL $11A9`-based flag post-processing. The EXP-0018 "clear v=0"
    events were MVN bytes landing on `+$11A6` — table entries whose
    power byte is 0 at that moment of capture.
  - `$C2299F` routine (entry of the `JSR $299F` at `$C2319D` ✓):
    per-attacker stat staging — `+$11AE` ← `+$3B2C,X`; `+$11AF` ←
    `+$3B18,X` (via `$C22C21`, cross-linking EXP-0017); flags from
    `+$3C45,X`; then the fight-command populate: **power `+$11A6` ←
    `+$3B68,X`** (`$C229D1` store, callback `$C229D4` ✓, with a
    `$B6`-bit7 X+1 quirk), `+$11A1` ← `+$3B90,X` (element),
    `+$11A8` ← `+$3B7C,X`, plus `+$3BA4,X` bit manipulation into
    DP `$B3`.
  - **Reinterpretation of EXP-0016/0018 divergence:** the identical-
    state trials diverged in *which action-setup path ran* — table
    (MVN) actions vs fight (populate) actions — i.e. **enemy AI action
    selection**. The RNG consumer is in the AI/selection layer above
    `$C23190` (which reads `+$3A70`+1 bit 0 — alternator candidate —
    and DP `$B5` before calling the setup).
  - **Field map for the 14-byte attack record** (entry byte k ↔
    `WRAM:+$11A0+k`, semantics from the pipeline decodes): +1 element
    byte; +2 flags (bit 0 physical-formula, bit 7 MP-path dispatch);
    +3 (bit 7 MP retarget); +4 mode (bit 7 base-variant select, bit 1
    flip chain, bit 0 party-halving gate); +6 power; +9 flag byte
    (ASL'd post-load); +10 abort bits (`#$82`). Bytes 0, 5, 7, 8, 11,
    12, 13 unknown.
- **Status:** completed (2026-07-30)
- **Confidence:** MVN table load, base address `$C4:6AC0`, 14-byte
  stride, destination block — **Confirmed (byte-exact)**. Fight-setup
  stat-table reads — **Confirmed (code)**; table labels
  (battle-power/element/etc.) Strong hypothesis. "RNG in AI selection
  layer" — Strong hypothesis (divergence localization; consumer still
  unread). Field-map semantics — Confirmed per the earlier
  pipeline verifications; unknown bytes stay Unknown.
- **Next action:** implement the attack-record format in Go (ROM
  reader + typed accessors, synthetic-fixture tests); AI-layer RNG
  hunt becomes question #30.
