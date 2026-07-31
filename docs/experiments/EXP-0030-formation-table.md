# EXP-0030: The +$3F46 producer — locating the ROM formation table

- **Status:** completed (2026-07-31)
- **Question:** which code writes the staged formation monster-id
  list at `WRAM:+$3F46` (and `+$3F52`) at encounter entry, and from
  which ROM structure does it read — the battle formation table?
- **Starting state:** headless bridge; write-watch armed on
  `+$3F40-+$3F5F`; mines encounter walk (EXP-0028 pattern).
- **Method:** probe `mesen/probes/EXP-0030.lua` (watchwrites +
  the walk helper); on the init burst, take the writer PC, dump its
  routine from the ROM file, decode the source operand; verify by
  reading the formation record for the known encounter (ids 77/78)
  from the ROM file.
- **Expected outcomes:** writer decoded to a ROM table read
  (formation base + stride) verified against ids 77/78; or a staging
  chain one hop deeper (recorded, bounded out).
- **Falsifying outcome:** no writes to +$3F46 at encounter entry.
- **Raw evidence paths:** `local_artifacts/experiments/EXP-0030/`.
- **Result:** **the formation table is located and verified** — the
  falsifier is not met.
  - Writers at encounter entry: `$C22624` (+$3F43 clear) and
    **`$C2315C`** — the staging copy loop. Its decode (from the ROM
    file): `LDA $11E0; ASL x4; SEC; SBC $11E0` → **X = formation_id
    x 15**; then `LDA $CF6200,X / STA $3F44,Y` copying 16 bytes.
    Live X=$0294 → **formation id 44**. A parallel **4-byte-stride
    flags table at `$CF5900`** loads to `$2F48/$2F4A` (id x4), and
    the formation id itself comes from `WRAM:+$11E0` (the encounter
    roll's output — its producer is the next hop).
  - **Formation record 44** (`ROMFILE:0x0F6494`):
    `B0 03 | 13 4D FF FF FF FF | CA 5D …` — first word $03B0
    (matches the live staged write), then the six monster-id bytes:
    **{19, 77}** with $FF empties.
  - **Correction to EXP-0028's coincidence:** the mines pair is
    monsters **19 and 77** (powers 13 and 19: monster 19 = power 13,
    HP 24), not 77/78 — record 78's power-13 match was coincidental.
    Census and inventories corrected.
- **Confidence:** formation table base `$CF6200`, 15-byte stride,
  id-x15 addressing, and the 16-byte staging copy — **Confirmed**
  (code decode + live X + record-content verification against two
  independently observed powers). Flags table `$CF5900` (4-byte) —
  Confirmed location, semantics Unknown. `+$11E0` as the encounter
  roll output — Strong hypothesis (consumer-side only). Formation
  bytes +$08.. (positions/mold?) — Unknown.
- **Next action:** `+$11E0`'s producer (the encounter roll — zone
  data + step counter, CEN-WORLD-0006's target); formation bytes
  +$00/01 and +$08.. semantics.
