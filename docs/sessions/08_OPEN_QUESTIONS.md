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

1. **What produces `$2E78–$2EA7`, and is it the authoritative party state?**
   *Partially answered in [SESSION_003](SESSION_003.md):* an upstream
   per-slot array at `WRAM:+$3BF4` receives the battle damage/heal/death
   stores, and its slot-0 value was observed propagating into `+$2E78`;
   init writers `ROMCPU:$C22408` (fills `+$2E7E` with `$FFFF`) and
   `$C25D33` also hit the region. **Still open:** the in-battle copier
   `+$3BF4` → `+$2E78`. *Experiment:* the preserved bridge's `+$2E78`
   write watch during one enemy hit; walk the captured stack.

1b. **Verify the Session 003 disassembly claims** (`ROMCPU:$C21323`,
   `$C21350`, `$C21390`, `$C213A7`; store sites `$C21338/$C21347/$C21396`;
   dispatch table at `$131F`; arrays `+$3C08/+$3C1C/+$3C30/+$3EE4/+$3C95`;
   sentinels `+$33E4/+$33D0`).
   *Experiment:* bridge ROM dump `$C21300–$C21410`, hand-disassemble,
   compare against `internal/game/battle/battle.go` claim by claim.

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
