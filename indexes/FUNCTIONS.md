# Functions

| ID | Title | Status | Confidence | Record |
|---|---|---|---|---|
| FN-0001 | CopyCharacterFields — `ROMCPU:$C10DF3`–`$C10E66` | Documented; implemented in `chardata` with tests | Confirmed (behavior, full disassembly); purpose of destination buffer strong hypothesis | [02_DISCOVERED_FUNCTIONS.md](../docs/sessions/02_DISCOVERED_FUNCTIONS.md) |
| FN-0002 | PerFrameBattleUpdate (candidate) — `ROMCPU:$C101FB`–`$C10225` | Documented; callees unexplored | Strong hypothesis (dispatcher role); code Confirmed | [02_DISCOVERED_FUNCTIONS.md](../docs/sessions/02_DISCOVERED_FUNCTIONS.md) |
| FN-0003 | HP delta routine — `ROMCPU:$C21323` (stores `$C21338`/`$C21347`) | Documented; implemented + tested (`internal/game/battle`) | **Confirmed (code, byte-exact — EXP-0001)**; "HP" label Strong hypothesis | [02_DISCOVERED_FUNCTIONS.md](../docs/sessions/02_DISCOVERED_FUNCTIONS.md) |
| FN-0004 | MP delta routine — `ROMCPU:$C21350` (arrays `+$3C08`/`+$3C30`) | Documented; implemented + tested (`internal/game/battle`) | **Confirmed (code, byte-exact — EXP-0001)**; "MP" label Strong hypothesis | [02_DISCOVERED_FUNCTIONS.md](../docs/sessions/02_DISCOVERED_FUNCTIONS.md) |
| FN-0005 | Death handler — `ROMCPU:$C21390` (store `$C21396`) | Documented; implemented + tested (`internal/game/battle`) | **Confirmed (code, byte-exact — EXP-0001)**; "death" label Strong hypothesis | [02_DISCOVERED_FUNCTIONS.md](../docs/sessions/02_DISCOVERED_FUNCTIONS.md) |
| FN-0006 | Delta dispatch wrapper + fetch — `ROMCPU:$C212F5`/`$C213A7` | Documented | **Confirmed (code, byte-exact — EXP-0001)**; tail purpose and delta producers Unknown | [02_DISCOVERED_FUNCTIONS.md](../docs/sessions/02_DISCOVERED_FUNCTIONS.md) |
| FN-0007 | PartyDisplaySourceRefresh (candidate) — `ROMCPU:$C25D26` | Documented | **Confirmed (code, byte-exact — EXP-0003)**; caller/trigger Unknown | [02_DISCOVERED_FUNCTIONS.md](../docs/sessions/02_DISCOVERED_FUNCTIONS.md) |
| FN-0008 | PendingDeltaAccumulate (candidate) — `ROMCPU:$C20C76` | Documented; implemented + tested (`battle.AccumulatePending`) | **Confirmed (code, byte-exact — EXP-0006)**; polarity semantics and amount provenance Unknown | [02_DISCOVERED_FUNCTIONS.md](../docs/sessions/02_DISCOVERED_FUNCTIONS.md) |
