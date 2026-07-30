# EXP-0003: Is the routine at `ROMCPU:$C25D33` the `$3BF4`→`$2E78` copier?

- **Status:** running (2026-07-29, overnight session)
- **Question:** Open question #1's missing piece — what writes the
  display-source current-HP array `WRAM:+$2E78` during battle? Candidate:
  the routine containing the store whose write-callback PC is
  `ROMCPU:$C25D33` (168 in-battle writes in EXP-0002's encounter; also
  captured in Session 003 and at battle init). Does it read from
  `WRAM:+$3BF4`?
- **Starting state:** irrelevant — deterministic ROM read
  (identity per [ROM_IDENTITY.md](../research/ROM_IDENTITY.md)).
- **Observation method:** bridge `read cpu C25CC0 176`
  (`ROMCPU:$C25CC0`–`$C25D6F`), hand-disassembly with arithmetic branch
  verification. Store expected at `$C25D30` (callback PC − 3, the
  convention EXP-0001 confirmed).
- **Expected outcomes:**
  - *Supports:* a 3-byte store to the `$2E78` region at `$C25D30` fed by
    a load from the `$3BF4` region → the authoritative→display chain is
    code-complete.
  - *Refutes:* the routine reads some other source (then that source
    becomes the new question).
- **Falsifying outcome:** no store to `$2E78`-region at `$C25D30`, or the
  loaded source is not `$3BF4`-region.
- **Raw evidence paths:** `mesen/out/rom_C25CC0_176.hex` (+ SHA-256),
  EXP-0002 `writers` snapshots in `mesen/out/exp4.log`.
- **Result:** **Confirmed, and more.** Dump
  `mesen/out/rom_C25CC0_176.hex` (SHA-256 `76861e12…a5dc7`). The routine
  at `ROMCPU:$C25D26`–`$C25D56` (PHP / REP #$20 / SEP #$10 / LDY #$06,
  loop DEY×2 to 0) copies **six** per-slot 16-bit values per slot:
  `$3BF4,Y→$2E78,Y` (store at `$C25D30` = callback PC −3 ✓),
  `$3C1C,Y→$2E80,Y`, `$3C08,Y→$2E88,Y`, `$3C30,Y→$2E90,Y`,
  `$3EE4,Y→$2E98,Y`, `$3EF8,Y→$2EA0,Y`. Every CharacterFieldsSource
  array is a copy of an authoritative battle array. EXP-0002 write counts
  (168 watched writes ÷ 4 slots) imply ~42 invocations over the ~6600
  battle frames — event-driven, not per-frame. Preceding code
  (`$C25CAA?`–`$C25D25`, partially in dump) is a separate Y≤`$0C` loop
  over `+$3388`/`+$200D`/`+$2015`/`+$3ECA` with gates `+$3C88`/`+$3C9D`
  (six-entry — enemy-side candidate); following routine `$C25D57`+
  zeroes `+$2E99+X` (X=6..0) and `+$2F35`–`+$2F40` — both unexplored.
- **Status:** completed (2026-07-29)
- **Confidence:** Copier code and source/destination pairing —
  **Confirmed** (byte-exact dump; store site matches live captures from
  two battles plus Session 003). Caller and trigger of `$C25D26` —
  Unknown. Semantic upgrades recorded in 03/04/05 and hypotheses.
- **Next action:** promote records; find `$C25D26`'s caller (new open
  question); live MP observation remains for the MP labels.
