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

- `$C10DF3` fired once per frame, every frame, while the battle ran
  (menu open, actions resolving) — and zero times at the title screen.
- The return address was `$C1:0203` on every capture (thousands of frames).
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

## Not yet established

- Whether `$2E78` or some deeper structure is the authoritative HP store.
- What consumes the `$2EB5` records (display code? menu code?) and the
  `$61AD` mask.
- What the other eight per-frame subroutines do.
- Whether PerFrameBattleUpdate also runs outside battles.

See [08_OPEN_QUESTIONS.md](08_OPEN_QUESTIONS.md).
