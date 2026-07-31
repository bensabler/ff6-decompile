# EXP-0029: The populate routine's X source — formation→monster mapping

- **Status:** completed (2026-07-31)
- **Question:** how is the monster-record offset (X=$09A0 observed)
  computed, and from what structure does the monster id come — i.e.,
  where are the battle formation records?
- **Starting state:** fully static: the EXP-0028 archived probe
  stacks (return addresses of the populate stores) and the local ROM
  file (code bytes at ROMFILE 0x2xxxx = ROMCPU $C2xxxx under HiROM).
  No emulator operation.
- **Method:** (1) parse the archived stack windows of the populate
  hits for the JSR return chain; (2) decode the caller code from the
  ROM file: the X-computation (expected: id×32 via shifts) and the
  id's source read (the formation structure); (3) if the source is a
  ROM table, bound its record shape with the known encounter (ids
  77/78) as the anchor.
- **Expected outcomes:** the id×stride computation decoded (stride
  proof) and a formation-record structure located at candidate level
  or better.
- **Falsifying outcome (for the stride-32 claim):** the computation
  is not a ×32 scaling of an id.
- **Raw evidence paths:** `local_artifacts/experiments/EXP-0029/`
  (code slices, decode notes, hashes.sha256).
- **Result:** the chain is decoded; the falsifier is not met (the
  scaling is exactly id×32 for the stat records).
  - **Per-slot loader `$C22C30`** (entry: A = monster id, Y = slot
    offset): `JSR $C22D71` stores the 16-bit `$CF8400[2×id]`
    attribute to `$3254,Y` (A restored), then ASL chains scale the
    id ×4 (`$CF3000`), ×8 (the `$CFC050`-region pair loop), and ×32
    (the `$CF0000` stat records, matching EXP-0028's X=$09A0 =
    77×32).
  - **The only caller, `$C22F22`**, loops X=5→0 / Y=$12→$04 (six
    enemy slots) reading the monster id from **`WRAM:+$3F46,X` — the
    formation's staged id list** ($FF = empty sentinel; an alternate
    id source `$0206` under the `$3A97` flag path — event-battle
    candidate). `+$3F52` (×4-scaled into DP `$EE`) is a second
    formation field.
  - `$CF8400` is therefore a **per-monster attribute table**, not
    the id source (the earlier id-pair search misfired for exactly
    this reason — negative result preserved).
- **Confidence:** loader entry/scaling and the `+$3F46` staging read
  — Confirmed (code decode anchored by the live EXP-0028 hit).
  `+$3F46` as "formation monster-id list" — Strong hypothesis
  (semantics from the consumer side only). `$0206`/`$3A97`
  event-battle path and `+$3F52` meaning — Tentative/Unknown.
- **Next action:** write-watch `+$3F46`/`+$3F52` at the next
  encounter entry to capture the ROM formation-record reader — that
  writer's source operand is the formation table.
