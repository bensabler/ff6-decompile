# EXP-0028: Battle-init enemy stat population — writers and ROM source

- **Status:** completed (2026-07-31)
- **Question (CEN-MONSTER-0001 next action, domain rotation):** which
  code populates the per-slot enemy stat tables
  (`WRAM:+$3B18/+$3B2C/+$3B68` family) at battle entry, and from
  which ROM records does it read — locating the monster database at
  candidate level or better?
- **Starting state:** headless bridge session. Primary path: arm
  write-watches on the stat-table family, load the free-walk mines
  state (`checkpoint3-mines.mss`), walk until a random encounter
  triggers battle init. Fallback if movement is event-locked there:
  advance the `checkpoint1.mss` scripted sequence (A presses) into a
  scripted battle. `exp10-battle.mss` cannot serve — its init already
  ran.
- **Method:** tracked probe `mesen/probes/EXP-0028.lua`:
  first-capture-per-PC write watch over `WRAM:+$3B18-+$3BB7` (the
  $14-stride family block: slots' battle stats), logging PC,
  address, value, registers. On battle init the burst of first-hits
  names the populate routine(s). Then: dump the writer routine(s)
  (`read cpu`) and hand-decode the source addressing — a long-address
  read or pointer walk exposes the monster-record base in ROM (or a
  WRAM staging area whose own producer is then watched in a bounded
  second pass).
- **Expected outcomes:**
  - *Direct:* a small writer set with a loop reading
    ROM (bank/pointer visible in the decode) → monster-record base
    at CANDIDATE_LOCATION/LOCATED; region registered.
  - *Indirect:* writers read from WRAM staging → the staging
    producer's first writer captured in the same session if cheap,
    else recorded as the bounded next experiment.
- **Falsifying outcome (for the "tables populated at init" model):**
  no writes to the family during a fresh battle entry (would mean the
  tables persist from elsewhere — itself a finding).
- **Bounds:** locate, do not extract; monster-record format decoding
  is a later unit. Encounter-walk budget: bounded press batches with
  probe-log polling; if no encounter triggers within the budget on
  both paths, record the negative and stop.
- **Raw evidence paths:** `local_artifacts/experiments/EXP-0028/`
  (events-slice.log `58be10d2…8978f3`, rom_C22C40_340.hex
  `fa14f416…78e6c8a`, rom_C22CB0_176.hex `3b23a544…dcdf160d`,
  m1.png, hashes.sha256).
- **Result:** **the monster database is located** — the falsifier is
  not met; the direct expected outcome landed.
  - The mines free-walk triggered a random encounter on walk loop 4
    (encounter mechanics themselves registered as a census
    observation). The write-watch burst at battle init (frames
    211052-211054) shows two clusters: `$C22899-$C2290F` writing
    party-slot stats (X=$0004 = slot 2) and `$C22CA3-$C22D8C`
    populating enemy slots (Y=slot offset) by **long-indexed reads
    from bank $CF**.
  - **Base: `ROMCPU:$CF0000` (`ROMFILE:0x0F0000`), X=$09A0 for this
    encounter's monster.** Field offsets reach +$1E and $09A0/32=77
    exactly → **32-byte record stride (Strong hypothesis), monster
    id 77 candidate**.
  - **Field map anchored to Confirmed arrays** (read offset → store):
    `+$08→$3BF4,Y and $3C1C,Y` (current+max HP), `+$0A→$3C08,Y and
    $3C30,Y` (current+max MP), `+$01→$3B68,Y` (the fight-populate
    power table, EXP-0019), `+$02→$3B7C,Y`, `+$10→$3B18,Y` with
    `+$00→$3B19,Y` (the EXP-0017 stat pair), `+$03/+$04→$3B54/55,Y`,
    `+$05→$3BB8,Y`, `+$07→$3B41,Y`, `+$0C→$3D84,Y`, `+$0E→$3D98,Y`,
    `+$13→$3C80,Y`, `+$1A→$3CA8,Y`, `+$1E→$3E4C,Y`.
  - **Question 19b (EXP-0004) closed:** the queued mystery writer
    `ROMCPU:$C22CCE` is this routine's max-HP store
    (`STA $3C1C,Y`).
  - **Five auxiliary bank-$CF tables** surfaced in the same routine
    (`$CF3000,X→$3308,Y`, `$CFC050,X`, `$CF37C0,X→$3C81,Y`,
    `$CF8400,X→$3254,Y`, `$CF8700,X`) — separate monster-adjacent
    tables, semantics Unknown, registered as leads.
- **Confidence:** monster-record base bank/address — **Confirmed**
  (long-address operands in the dumped, live-hit routine). 32-byte
  stride and id-77 identity — Strong hypothesis (divisibility +
  offset range; the X-computation upstream is undecoded). Field
  semantics — Confirmed where the destination array is Confirmed
  (HP/MP/power/staging), Candidate elsewhere. Aux tables — Unknown.
- **Next action:** decode the X-computation upstream (stride proof +
  formation→monster-id mapping — the formation record consumer);
  then bulk-extract candidate records and cross-check observed values
  (Were-Rat/Repo Man powers 13/19 from EXP-0018 should appear at
  +$01 of their records).
