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
    trigger?** Write-count arithmetic says event-driven (~42 calls per
    battle), not per-frame. *Experiment:* exec watch at `$C25D26` with
    stack capture during one attack resolution.

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

13. **What writes the pending-delta arrays `WRAM:+$33E4`/`+$33D0`?**
    The battle-logic layer above the delta engine.
    *Experiment:* write watch on `+$33E4`–`+$33F3` during one attack
    command; walk captured stacks.

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
