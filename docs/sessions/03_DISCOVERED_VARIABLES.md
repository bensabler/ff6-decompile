# Discovered Variables

### PartySlotCurrentHP (first displayed party slot)

- **Address:** `$2EB5` (WRAM-relative); SNES CPU `$7E2EB5`
- **Address space:** WRAM-relative
- **Kind:** Variable (16-bit unsigned)
- **Status:** Confirmed
- **Previous names:** `PartySlot[0].CurrentHP` (Session 001 label)
- **Observed behavior:** Holds the value displayed as the first party slot's
  current HP; refreshed once per frame by CopyCharacterFields from
  `$2E78[0]`.
- **Evidence:**
  - Session 001: converging 16-bit exact-value search; write breakpoint at
    `$C10E14` (`STA $2EB5,Y`).
  - Session 002: live reads matched the displayed HP for slots 0–2
    (53/39/33 on screen = `$35/$27/$21` in `$2E78` and in records 0–2),
    tracked through damage and a heal.
- **Interpretation:** Field `+0` of destination record 0 — a per-frame copy
  of `$2E78[0]`. The array at `$2E78` is upstream of it; which one gameplay
  logic treats as authoritative is still open.
- **Alternative explanations:** none remaining for identity; authority
  (display copy vs. gameplay value) unresolved.
- **Validation experiment:** freeze/edit `$2EB5` vs `$2E78[0]` separately
  and observe display and damage resolution.
- **Go representation:** `PartySlotRecord.CurrentHP` in
  [internal/game/chardata/chardata.go](../../internal/game/chardata/chardata.go)
- **Related discoveries:** CopyCharacterFields,
  [05_DATA_STRUCTURES.md](05_DATA_STRUCTURES.md)
- **First observed in session:** [SESSION_001](SESSION_001.md)
- **Last updated:** 2026-07-29 ([SESSION_002](SESSION_002.md))

### UnknownFlag628D

- **Address:** `$628D` (WRAM-relative); SNES CPU `$7E628D`
- **Address space:** WRAM-relative
- **Kind:** Variable (8-bit flag; tested for zero/nonzero)
- **Status:** Confirmed as a masking gate in CopyCharacterFields; gameplay
  meaning Unknown
- **Previous names:** None
- **Observed behavior:** When nonzero, CopyCharacterFields never masks
  fields `+8`/`+$A`. Read `$00` throughout the observed battle.
- **Evidence:** Disassembly `$C10DFA–$C10DFD`; live read.
- **Interpretation / Alternatives:** Unknown. Could be a mode flag (e.g.,
  a battle subtype) — no evidence.
- **Validation experiment:** Write breakpoint on `$628D`; catch what sets it.
- **Go representation:** `CopyMode.UnknownFlag628D`
- **Related discoveries:** CopyCharacterFields, UnknownFlagE9EF
- **First observed in session:** [SESSION_002](SESSION_002.md)
- **Last updated:** 2026-07-29

### UnknownFlagE9EF

- **Address:** `$E9EF` (WRAM-relative); SNES CPU `$7EE9EF`
- **Address space:** WRAM-relative
- **Kind:** Variable (8-bit flag; tested for zero/nonzero)
- **Status:** Confirmed as a masking gate in CopyCharacterFields and a
  branch condition in PerFrameBattleUpdate; gameplay meaning Unknown
- **Previous names:** None
- **Observed behavior:** When nonzero (and `$628D` zero), CopyCharacterFields
  masks field `+8` with `$0038` and zeroes field `+$A`. PerFrameBattleUpdate
  skips `JSL $C20003` when it is nonzero. Read `$00` throughout the observed
  battle.
- **Evidence:** Disassembly `$C10DFF–$C10E0B` and `$C1021A–$C10222`; live
  read.
- **Interpretation / Alternatives:** Unknown. Its use in two places in the
  per-frame path suggests a battle-wide mode.
- **Validation experiment:** Write breakpoint on `$E9EF`; catch what sets it.
- **Go representation:** `CopyMode.UnknownFlagE9EF`
- **Related discoveries:** CopyCharacterFields, PerFrameBattleUpdate
- **First observed in session:** [SESSION_002](SESSION_002.md)
- **Last updated:** 2026-07-29

### UnknownSlotMask61AD

- **Address:** `$61AD` (WRAM-relative); SNES CPU `$7E61AD`
- **Address space:** WRAM-relative
- **Kind:** Variable (8-bit, low 4 bits used)
- **Status:** Confirmed as CopyCharacterFields' output; consumer and meaning
  Unknown
- **Previous names:** None
- **Observed behavior:** Rewritten once per frame: bit n (n = 0..3) is set
  when bit 13 of source field `$2EA0[n]` is clear. Read `$0F` with `$2EA0`
  all zero — formula verified live.
- **Evidence:** Disassembly `$C10E46–$C10E65`; live read.
- **Interpretation:** A per-slot condition mask derived from an unknown
  status-like field. Consumer not yet identified.
- **Alternative explanations:** none until the consumer is found.
- **Validation experiment:** Read breakpoint on `$61AD` to find consumers.
- **Go representation:** return value of `chardata.CopyCharacterFields`
- **Related discoveries:** CopyCharacterFields
- **First observed in session:** [SESSION_002](SESSION_002.md)
- **Last updated:** 2026-07-29
