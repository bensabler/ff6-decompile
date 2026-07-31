# Battle System

Current picture of the battle architecture, from Session 002's live
instrumentation of the Narshe intro battle (3-member Magitek party vs
guards).

## Per-frame battle update chain

**Status: strong hypothesis** (single battle context observed so far).

```text
bank $C2 frame driver (unexplored; JSL at ~$C26425)
  └─ JSL $C101FB  PerFrameBattleUpdate — once per frame during battle
       ├─ JSR $C11A24   unexplored
       ├─ JSR $C10DF3   CopyCharacterFields — refresh party-slot records
       ├─ JSR $C14504   unexplored
       ├─ JSR $C12F79   unexplored
       ├─ JSR $C102CA   unexplored
       ├─ JSR $C144BE   unexplored
       ├─ JSL $C2BF53   unexplored
       ├─ JSR $C193E3   unexplored
       ├─ JSL $C2B41A   unexplored
       └─ JSL $C20003   only when $E9EF == 0
```

Observed facts supporting this:

- `$C10DF3` fired repeatedly while battles ran and zero times at the
  title screen and in field/event contexts (EXP-0002: ≈175k non-battle
  frames, zero fires across four contexts).
- Rate (EXP-0002 correction): up to ~1/frame when the battle is settled
  (~0.87 observed), lower during entry/heavy phases (≈0.23–0.7) — cause
  undiscriminated (lag frames vs conditional dispatch).
- The return address was `$C1:0203` on every steady-state capture
  (thousands of hits); one additional one-shot caller at battle entry:
  `JSR` at `ROMCPU:$C11090` (EXP-0002).
- The stack shows entry via `JSL` from bank `$C2` (return `$C2:6429`).

## Party-slot data flow

**Status: confirmed** for the copy itself; the producer of `$2E78` and the
consumer of `$2EB5` records are not yet identified.

```text
gameplay logic (unidentified)
  └─ writes $2E78–$2EA7   six parallel per-slot arrays (HP confirmed)
       └─ CopyCharacterFields (per frame)
            ├─ writes $2EB5+n*$20, offsets +0..+$B (records 0–2 = displayed slots 0–2)
            └─ writes $61AD 4-bit slot mask (consumer unidentified)
```

Battle events observed flowing through this path: enemy damage decrementing
slot HP values, and a Heal Force cast setting slot 0's current HP to exactly
the `$2E80` (max HP candidate) value — all mirrored into the records on the
next frame.

## Authoritative layer (Session 003 + EXP-0001)

**Status: engine code Confirmed byte-exact** (EXP-0001 dump); the array's
authoritative role and the semantic labels remain strong hypotheses;
pending-delta producers unidentified.

```text
pending-delta producers (unidentified)
  └─ WRAM:+$33E4 / +$33D0 per-slot deltas ($FFFF = none)
       └─ fetch @ ROMCPU:$C213A7 (delta = $33E4 − gated $33D0)
            └─ dispatch @ ROMCPU:$C21300 JSR ($131F,X);
               WRAM:+$11A2 bit7 selects entry ($C2131F: $1323 HP / $1350 MP)
                 ├─ HP @ $C21323: heal clamp→$3C1C @ $C21338, damage @ $C21347
                 ├─ MP @ $C21350: same over $3C08/$3C30; $3C95 bit0 → death;
                 │                exits LDA #$0080 / JMP $C2464C (unexplored)
                 └─ death @ $C21390: clears $3A89, zeroes HP @ $C21396,
                    $3EE4 bit1 suppresses, else JMP $C20E32 A=$0080 (unexplored)
                      └─ writes WRAM:+$3BF4 per-slot array (Y = slot×2)
                           └─ PartyDisplaySourceRefresh @ ROMCPU:$C25D26
                              (event-driven; caller unknown; EXP-0003)
                              copies all six arrays:
                              $3BF4→$2E78  $3C1C→$2E80  $3C08→$2E88
                              $3C30→$2E90  $3EE4→$2E98  $3EF8→$2EA0
                                └─ CopyCharacterFields → WRAM:+$2EB5 records
Lifecycle: $FF fill @ ROMCPU:$C0567B (boot/teardown);
           init @ $C223F6/$C227B4/$C22408 (battle start)
```

## Not yet established

- Whether `$2E78` or some deeper structure is the authoritative HP store.
- What consumes the `$2EB5` records (display code? menu code?) and the
  `$61AD` mask.
- What the other eight per-frame subroutines do.
- Whether PerFrameBattleUpdate also runs outside battles.

See [08_OPEN_QUESTIONS.md](08_OPEN_QUESTIONS.md).


## Damage pipeline (EXP-0006..0019, 2026-07-30)

**Status: Confirmed byte-exact** for every stage below; semantic labels
at their recorded hypothesis levels; RNG consumption localized to the
AI/action-selection layer (question #30, open).

```text
AI / action selection (RNG consumer - undecoded, #30)
  └─ action setup: 14-byte attack record MVN'd from ROMCPU:$C46AC0+14×n
     into WRAM:+$11A0-+$11AD (spells), or fight populate from per-slot
     stat tables (+$3B18/+$3B2C/+$3B68/+$3B7C/+$3B90/+$3BA4) @ $C2299F
       └─ target loop @ ~$C23442 (Y=$12→0, +$3018,Y bits vs DP $A4)
            └─ base amount @ $C20B83: standard ($C22B69:
               power×4 + power×$11AE×$11AF>>5) or physical ($C22BA6:
               vigor² shape) per $11A2 bit 0; +$11B0 staging
                 └─ defense (255−def)/256+1 · halvings · party-vs-party
                    halving · element response (+$3BCC/+$3BE0 masks)
                    · ×1.5 chain ($BC) → DP $F0
                      └─ accumulator @ $C20C76 → pending arrays
                         (+$33D0/+$33E4, $FFFF none, cap 9999)
                           └─ delta engine @ $C21300 dispatch →
                              +$3BF4-family arrays (10 slots)
                                └─ display copier @ $C25D26 →
                                   +$2E78 family → CopyCharacterFields
                                   @ $C10DF3 → +$2EB5 records → HUD
```

Numeric closure: power 60 / stats 28,4 → base 450 → defense ≈58 →
346-observed Fire Beam (EXP-0007/0014/0015). Misses arrive as power 0
(cleared-but-never-populated action blocks, EXP-0018/0019).
