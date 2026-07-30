# EXP-0017: Complete the enemy/physical base path; where is the RNG?

- **Status:** running (2026-07-30)
- **Question (#26 + #28):** Decode the `$11A2`-bit-0-selected base path
  (`ROMCPU:$C22B9D` entry, store `$C22BE9`, party-tail past `$C22BFF`,
  miss writer `$C22C02`) precisely, and enumerate its reads — the
  timing-varying one is the RNG state.
- **Starting state:** deterministic ROM reads; `rom_C22B40_192.hex`
  already in hand; supplemental `read cpu C22BFF 48` for the tail.
- **Method:** hand-disassembly anchored at the verified `$C22BE9`
  store and the `$C22B9C` RTS; read census of the whole path.
- **Expected outcomes:** either an RNG-state read inside the path
  (locates #28's consumer), or a fully deterministic read set
  (pushing RNG consumption upstream into hit/AI action setup — then
  the `$C22C02` miss writer's surrounding routine becomes the next
  bracket).
- **Falsifying outcome (for "RNG in the damage path"):** all reads are
  the known action-state bytes (`$11A2/$11A6/$11AE/$11AF/$B2/$11B0`).
- **Raw evidence paths:** `rom_C22B40_192.hex`,
  `mesen/out/rom_C22BFF_48.hex`.
- **Result:** (artifact `rom_C22BFF_48.hex` SHA-256 `4af2d686…d8774c`,
  joining `rom_C22B40_192.hex`) — **the physical path is fully decoded
  and contains no RNG read.**
  - Shared entry `$C22B9D`: `$11A2` bit 0 → physical; else
    `JMP $2B69` (standard).
  - Physical: `t = power` (×4 for enemy attackers, `CPX #$08`);
    ×1.75 unless `$B2` bit 14 (`(t/2+t)/2+t` stack shape); `+ $11AE`
    (8-bit add with carry propagation); then the `$E8/$E9/$EA` juggle
    reconstructs the **low 16 bits of the full product** `t×$11AF`,
    scaled `×$11AF/256` again → **`t = ((t×statB)&0xFFFF)×statB/256`**
    (statB², with statB=16 as identity); store `$C22BE9`.
  - Party-attacker tail (`$C22BF0`–`$C22C01`): `base = ((power×2 +
    base)>>1) + base` = **base×1.5 + power**; then attacker flags
    `$3C58,X` (bit 0 → halve; bit 3 → ×0.75 via `-(base>>2)`).
    Enemy attackers skip the tail and flags (`BCS $C22C1F`).
  - The EXP-0016 "miss writer `$C22C02`" was this tail store at
    `$C22BFF` (+3 rule); on the miss action the inputs were zero —
    **misses are decided upstream** (power arrived 0).
  - **Read census (whole path): `$11A2, $11A6, $11AE, $11AF, $B2,
    $11B0, $3C58,X` — all deterministic action-state. No RNG.**
    Variance (EXP-0016) therefore enters in the action-setup/hit-roll
    layer that populates the `$11A1`–`$11B0` block per action.
  - Bonus: the following routine (`$C22C21`+) writes `$11AF` from
    per-slot table `+$3B18,X` — a `$11AF` producer (question #27
    lead). `$3C58 = $3C44+$14` — family member #13.
  - Numeric closure: enemy base 7 = power 1 ×4 ×1.75 with statB=16
    (identity) — matches the EXP-0016 captures.
- **Status:** completed (2026-07-30)
- **Confidence:** Path decode — **Confirmed (byte-exact, anchored at
  both verified stores and the RTS)**. "RNG upstream of the damage
  arithmetic" — **Confirmed** (complete read census). statB² reading —
  Confirmed arithmetic; `$3B18` = statB producer — Confirmed (code);
  labels (vigor/level) — Unknown.
- **Next action:** implement `BaseAmountPhysical`; hunt the action-setup
  layer via a `+$11A6` write watch (question #29, the true RNG
  consumer).
