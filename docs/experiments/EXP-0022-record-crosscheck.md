# EXP-0022: Attack-record table ROM cross-check

- **Status:** completed (2026-07-30)
- **Question (research-queue P1 semantic debt):** does the ROM
  attack-record table at `ROMCPU:$C46AC0` (EXP-0019, Confirmed base and
  stride) contain entries matching independently observed action
  parameters — specifically (a) the Fire Beam signature (power 60,
  element bit 0, both live-verified in EXP-0011/0015 from the action
  block), and (b) record 238's power byte 0 (EXP-0018 observed the MVN
  writing `v=0` to `+$11A6`; EXP-0021 confirmed every enemy Fight load
  is index 238)?
- **Starting state:** headless bridge session (testrunner, EXP-0021
  environment still live). ROM reads are state-independent.
- **Method:** bridge `read cpu C46AC0 3584` (records 0–255, a bounded
  scan window that covers the observed index 238; the table's true
  length is unknown). Archive the dump + SHA-256 under
  `local_artifacts/experiments/EXP-0022/`. Decode every record with
  `attackdata.RecordAt` via a new `ff6lab attackdata scan` subcommand
  (first consumer of the package against real ROM bytes). Enumerate
  power-60 entries and their element bytes; decode record 238.
- **Expected outcomes:**
  - *Supports the table reading:* at least one record has power 60 with
    element bit 0 (Fire Beam candidate — index becomes a Strong
    hypothesis, value-coincidence only); record 238 has power byte 0.
  - *Refutes / complicates:* no power-60+fire entry in 0–255 → the
    party Magitek action's power 60 was not MVN-loaded from this table
    slice (populate-side source, or table longer than 256); record 238
    power ≠ 0 → the EXP-0018 v=0 interpretation needs rework.
- **Falsifying outcome (for the cross-check):** record 238's power
  byte ≠ 0. (The Fire Beam scan can only support or stay open — a
  miss narrows, but does not refute, the table reading.)
- **Raw evidence paths:** `local_artifacts/experiments/EXP-0022/`
  (table_C46AC0_3584.hex SHA-256 `e493beae…3612f`, commands.log
  `694fcf93…72cd4`, scan-output.txt `6a56ee85…20b7f`, hashes.sha256).
- **Result:** both planned checks pass, plus one unplanned
  confirmation. 256 records decoded cleanly by
  `ff6lab attackdata scan` (first use of the package against real ROM
  bytes; the falsifier is **not** met).
  - **Record 238** (`43 00 01 22 00 00 00 00 00 FF 00 00 00 00`):
    **power byte (+6) = 0** — cross-checks EXP-0018's observed MVN
    write `v=0` to `+$11A6` for enemy Fight loads (EXP-0021: all enemy
    loads are index 238). **Bonus: flags2 (+2) bit 0 is set** — the
    physical-formula select (EXP-0017's live decode showed enemy Fight
    actions run the physical path). Element 0, mode 0.
  - **Fire Beam signature (power 60, element bit 0, standard formula):
    two candidates — indices 5 and 131** (`61 01 00 28 00 14 3C…` and
    `43 01 00 20 20 00 3C…`; both flags2 bit 0 clear, matching the
    EXP-0015 standard-path observation). Five other power-60 records
    (15, 124, 134, 190, 223) have different element bytes.
  - *Shape observation:* records 0–2 carry ascending single-bit
    elements `$01/$02/$04` with powers 21/22/20 and otherwise
    identical bytes — a coherent elemental-spell family at the table
    head.
- **Confidence:**
  - Record 238 power=0 and physical-flag bit — **Confirmed
    (byte-exact, live ROM dump; converges with EXP-0017/0018/0021 live
    behavior)**.
  - Table base/stride and the `attackdata` decoder against real ROM —
    **Confirmed** (clean 256-record decode; field values coincide with
    live pipeline observations).
  - "Index 5 or 131 is the Magitek Fire Beam record" — **Tentative
    hypothesis** (value coincidence only; the party action's index was
    never traced — EXP-0021's press-following solo loads were all 238,
    so the Fire Beam confirm was not reached or stages differently).
  - Table length > 256 records — Unknown (only 0–255 dumped).
- **Next action:** to disambiguate the Fire Beam index, capture the
  `+$11A0`–`+$11AD` block (or the MVN `X` source) during an actual
  Magitek Fire Beam confirm — needs a menu-navigation press script
  against the battle savestate; queue as a small follow-up unit.
