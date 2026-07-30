# Open Questions

Ordered roughly by how much each answer would unlock.

## Answered in Session 002 (kept for the record)

- ~~Who calls `$C10DF3`?~~ → `JSR` at `$C10200` inside `$C101FB`
  (PerFrameBattleUpdate), entered via `JSL` from bank `$C2` (~`$C26425`),
  once per frame during battle, never at the title screen.
- ~~Do records 1–3 hold party slots?~~ → Records 0–2 matched displayed
  slots 0–2 (names/HP 53/39/33) exactly. Record 3 unverified (3-member
  party).
- ~~Accumulator/index width at `$C10E14`?~~ → `m=0, x=0` (16-bit both);
  entry is `m=1` and the prologue's `REP #$20` switches.
- ~~Initial X/Y?~~ → Both zeroed by `TDC/TAX/TAY` after the prologue.
- ~~Data bank?~~ → `DB = $7E` confirmed at entry and store; `D = $0000`.

## Open

1. ~~**What produces `$2E78–$2EA7`, and is it the authoritative party
   state?**~~ → **Answered** ([SESSION_003](SESSION_003.md) +
   [EXP-0003](../experiments/EXP-0003-2e78-producer.md)):
   PartyDisplaySourceRefresh (`ROMCPU:$C25D26`) copies all six
   authoritative battle arrays (`+$3BF4` family, mutated by the delta
   engine) into the `+$2E78` display-source family. The `+$2E78` family
   is *not* authoritative. Follow-ups: #19 (copier's caller/trigger),
   #2 (records' consumer) unchanged.

19. **Who calls PartyDisplaySourceRefresh (`ROMCPU:$C25D26`) and on what
    trigger?** *Partially answered (EXP-0004):* event-driven confirmed
    (13 calls in ~80 s); every steady stack tops with the `JSR $069B` at
    `ROMCPU:$C21409` (post-fetch driver) → `$C2069B` → (tail path,
    unverified) → `$C25D26`. **Still open:** the `$C2069B` code and the
    non-steady first-call path. *Experiment:* dump `$C20690`–`$C20700`.

19b. **What is the writer at `ROMCPU:$C22CCE`?** Wrote enemy-HP entries
    4 times around the EXP-0005 victory (cleanup candidate).
    *Experiment:* dump `$C22CA0`–`$C22D00`; exec-log during a victory.

20. **What are the routines flanking the copier** — the six-entry
    `Y≤$0C` loop over `+$3388`/`+$200D`/`+$2015` (enemy-side candidate,
    entry before `$C25CC0`) and the `$C25D57`+ zeroer of `+$2E99+X` and
    `+$2F35`–`+$2F40`? *Experiment:* dump `$C25C60–$C25CC0` and
    `$C25D57–$C25DA0`; exec watches.

1b. ~~**Verify the Session 003 disassembly claims.**~~ → Answered by
   [EXP-0001](../experiments/EXP-0001-c2-delta-engine-dump.md)
   (2026-07-29): every claim verified byte-exact; full listing in
   [02_DISCOVERED_FUNCTIONS.md](02_DISCOVERED_FUNCTIONS.md). New
   follow-ups: 13–16 below.

13. ~~**What writes the pending-delta arrays?**~~ → **Answered**
    (EXP-0004 + [EXP-0006](../experiments/EXP-0006-delta-setter.md)):
    PendingDeltaAccumulate (`ROMCPU:$C20C76`) accumulates DP `$F0` into
    the slot's pending array with `$FFFF`-sentinel init and a **9999
    cap**; sweepers `$C2638E`/`$C26391`; init `$C22408`. (EXP-0004's
    "`$C20434` caller" was a stack misparse — corrected in EXP-0006.)
    Remaining piece is #21.

21. **What computes the amount in DP `$F0` before queueing** — the
    actual damage/heal formula? *Narrowed by
    [EXP-0007](../experiments/EXP-0007-amount-correlation.md):* `$F0`
    holds the final per-hit number (equals popup and applied delta),
    so the complete formula runs upstream of the `$C20C28` JSR.
    Attacker/target arrive as `X`/`Y` (slot×2). Staging candidates:
    `$C20420`+ (DP `$F4`–`$FC`), conditional doubling at `$C20C1A`.
    *Progress:* the `$6B16` lead was refuted (EXP-0008); the verified
    frame is the `JSR $0B83` at `$C23469` (EXP-0009), bracketing the
    formula body to `ROMCPU:$C20B83`–`$C20C2C` inside a ten-slot
    target loop gated by `+$3018,Y` bits vs DP `$A4`.
    *Progress (EXP-0010):* the post-processing is fully decoded — the
    elemental-modifier block (`$C20B83`–`$C20C2C`) with per-target
    response masks (`+$3BCC`/`+$3BE0` family arrays) and status-driven
    polarity flips. The remaining unknown is #22.

22. ~~**What do the base-amount routines compute?**~~ → **Decoded**
    ([EXP-0011](../experiments/EXP-0011-base-formulas.md)): variant A =
    defense/halving post-processing of a precomputed base (`+$11B0`);
    variant B = fraction-of-HP with min 1; `+$11A3` bit 7 = HP→MP index
    retarget. Remaining unknowns split into #23–#25.

23. **What computes the base amount at `WRAM:+$11B0`** (battle power ×
    level/stat layer, upstream of the whole pipeline)?
    *Experiment:* write watch on `+$11B0` during one attack; walk the
    stack.

24. ~~**What is the final transform `ROMCPU:$C2370B`?**~~ → **Answered,
    hypothesis refuted**
    ([EXP-0012](../experiments/EXP-0012-arithmetic-helpers.md)): it is
    a ×1.5-per-count chain (count in DP `$BC` from the target loop's
    `+$3A54` gate, clamped `$FFFF`, consumed on use) — not randomness.
    Where the damage variance comes from folds into #23/#25.

25. ~~**What is the core multiply `ROMCPU:$C24781`?**~~ → **Answered**
    ([EXP-0013](../experiments/EXP-0013-core-multiply.md)): the SNES
    hardware 8×8 multiply of A's own bytes (`$4202/$4203` via one
    16-bit store, product from `$4216`); the `$C247B7` wrapper is
    algebraically `floor(value×$E8/256)`, making the defense scaling
    exactly `(amount×(255−def))/256 + 1`. Bonus: `$C24792` = hardware
    divide helper. Remaining innermost frontier: #23 (`+$11B0`
    producer).

14. **What does `WRAM:+$11A2` mean** (bit 7 = MP-path selector)?
    *Experiment:* write watch on `+$11A2` around command execution.

15. **What are the death continuation `ROMCPU:$C20E32` (A=`$0080`) and
    the MP exit tail `ROMCPU:$C2464C` (A=`$0080`)?**
    *Experiment:* dump and disassemble both entry regions.

16. **What is the post-dispatch tail** (`$C21307`–`$C2131B`: `$C2362F`
    call, `+$327C,Y` store, `TRB $3419`) **and when does the engine
    return carry set?**
    *Experiment:* exec watch at `$C21307` with register capture.

17. **What calls CopyCharacterFields from `ROMCPU:$C11090` at battle
    entry** (one-shot `JSR`, EXP-0002), and what else does that init path
    do? *Experiment:* dump around `$C11080–$C110C0`; exec-log the region
    during an encounter transition.

18. **Why does the per-frame stream skip frames during battle entry/heavy
    phases** (≈0.23–0.7/frame vs ~0.87 settled, EXP-0002)? Lag frames vs
    conditional dispatch in the bank-`$C2` driver.
    *Experiment:* compare fire count against Mesen's lag-frame counter
    over the same window.

2. **What consumes the `$2EB5` records and the `$61AD` mask?**
   *Experiment:* read breakpoints on `$2EB5` and `$61AD`; classify readers.

3. **What do the two flag bytes mean?** `$628D` (suppresses masking) and
   `$E9EF` (enables masking; also gates `JSL $C20003` in the per-frame
   chain). Both were 0 all battle.
   *Experiment:* write breakpoints on both; play until they change (try
   different battle types, menus, events).

4. **What are fields `$2E98` (observed 8,8,8 then `$0208` mid-battle) and
   `$2EA0` (bit 13 → `$61AD`)?** The `$0038` mask isolates bits 3–5 of
   `$2E98` values in masked mode.
   *Experiment:* log both arrays each frame across status changes
   (poison/sleep etc. once available) and diff against events.

5. **Who writes record bytes `+$0C..$1F`** (live non-zero data, e.g.
   `FF FF FF FF FF 0E 96 84 83 86 84 FF` in record 0)?
   *Experiment:* write breakpoint on `$2EC1–$2ED4`.

6. **What are the other eight subroutines in PerFrameBattleUpdate?**
   (`$C11A24`, `$C14504`, `$C12F79`, `$C102CA`, `$C144BE`, `$C2BF53`,
   `$C193E3`, `$C2B41A`, `$C20003`.)
   *Experiment:* dump and disassemble one at a time, starting with
   `$C11A24` (runs immediately before the copy).

7. **What drives the per-frame chain from bank `$C2`?**
   *Experiment:* dump around `$C26425`; exec-log it in and out of battle.

8. **Is `$2EB5` the true record start?** (Could be an interior field.)
   *Experiment:* write breakpoints on `$2EA8–$2EB4`.

9. **Do `$2E88`/`$2E90` hold current/max MP?** (24/24 for slot 0 only.)
   *Experiment:* battle with an MP-spending party member; watch `$2E88`.

10. **Does record 3 / slot 3 behave identically with a 4-member party?**
    *Experiment:* progress to a 4-member party; verify `$2F15`.

11. **Confirm max HP.** Current evidence: heal set current = `$2E80` value;
    gauges imply a max. One more independent observation (level-up or a
    max-HP-boosting effect) would upgrade it to Confirmed.

12. **What is the 5-byte copy routine near `$C101CC`** (`$2E72→$602D+`,
    counter `$64DA`)? Noticed in the caller dump; unexplored.
