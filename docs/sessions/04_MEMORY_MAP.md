# Memory Map

Authoritative list of every address this project has evidence for.
Notation: WRAM-relative offsets like `$2EB5` (SNES CPU `$7E2EB5`); ROM CPU
addresses like `$C10DF3`. See
[01_REVERSE_ENGINEERING_RULES.md](01_REVERSE_ENGINEERING_RULES.md).

## WRAM

| Range (WRAM-relative) | Size | Status | Description |
| --- | --- | --- | --- |
| `$2E72–$2E76` | 5 | Tentative | Copied in 5-byte units to `$602D+` by an unexplored routine near `$C101CC` (noticed in the Session 002 caller dump; not investigated). |
| `$2E78–$2EA7` | `$30` | **Confirmed** (layout) | Source region for CopyCharacterFields: six parallel arrays of four 16-bit entries — `$2E78` current HP (slots 0–2 matched display), `$2E80` max HP (candidate), `$2E88`/`$2E90` unknown (24/24 for slot 0; MP candidates), `$2E98` unknown (8,8,8 observed), `$2EA0` unknown (bit 13 drives `$61AD`). Layout in [05_DATA_STRUCTURES.md](05_DATA_STRUCTURES.md). |
| `$2EB5` | 2 | **Confirmed** | Current HP, first displayed party slot. [03_DISCOVERED_VARIABLES.md](03_DISCOVERED_VARIABLES.md). |
| `$2EB5–$2F34` | `$80` | **Confirmed** (fields `+0..+$B`) | Four destination records, stride `$20`, bases `$2EB5`, `$2ED5`, `$2EF5`, `$2F15`; records 0–2 verified against displayed party slots. Bytes `+$C..+$1F` of each record hold live data written by other, unidentified code. |
| `$61AD` | 1 | **Confirmed** (writer) | 4-bit per-slot mask output of CopyCharacterFields; consumer unknown. |
| `$602D–…` | ? | Tentative | Destination of the 5-byte copy routine noticed near `$C101CC`; indexed by a counter at `$64DA`. Not investigated. |
| `$3BF4–$3BFB` (`WRAM:+$3BF4`) | 8 | **Confirmed** (writers); authority Strong hypothesis | Per-slot 16-bit array, stride 2, written by battle delta stores (`ROMCPU:$C2134A/$C2133B/$C21399` logged PCs), init writers (`$C223F6`, `$C227B4`), `$FF`-filled by `$C0567B`, zeroed by `$C206BC/BF` and bank-`$C3` `$C36A58/5E`. Candidate authoritative current HP; slot-0 value observed propagating to `+$2E78`. [SESSION_003](SESSION_003.md). |
| `$628D` | 1 | **Confirmed** (use) | Masking gate flag; meaning unknown. |
| `$64DA` | 1? | Tentative | Counter (low nibble used ×5 as index), incremented by the unexplored 5-byte copy routine. |
| `$E9EF` | 1 | **Confirmed** (use) | Masking gate / per-frame branch flag; meaning unknown. |

Note: `$2EB5` lies above the low-RAM mirror (`$0000–$1FFF`), so absolute
addressing reaches WRAM only when `DB = $7E`. **Confirmed** in Session 002:
`DB = $7E` at entry and at the store. Direct page `D = $0000`; the routine
uses `$10`/`$12` as direct-page scratch.

## ROM

| Address (ROM CPU) | Status | Description |
| --- | --- | --- |
| `$C10DF3–$C10E66` | **Confirmed** | CopyCharacterFields — fully disassembled. [02_DISCOVERED_FUNCTIONS.md](02_DISCOVERED_FUNCTIONS.md). |
| `$C10E11` | **Confirmed** | Copy-loop head (branch target of `BNE`). |
| `$C10E14` | **Confirmed** | `STA $2EB5,Y` — Session 001's original breakpoint hit. |
| `$C10E4D` | **Confirmed** | Mask-loop head. |
| `$C101FB–$C10225` | **Confirmed** (code) | PerFrameBattleUpdate (candidate name) — calls CopyCharacterFields at `$C10200`. |
| `$C26425` | Strong hypothesis | `JSL $C101FB` call site (from stack return `$C2:6429`); surrounding code unexplored. |
| `$C11A24`, `$C14504`, `$C12F79`, `$C102CA`, `$C144BE`, `$C193E3` | Unknown | Bank-`$C1` subroutines called by PerFrameBattleUpdate; unexplored. |
| `$C2BF53`, `$C2B41A`, `$C20003` | Unknown | Bank-`$C2` long-called subroutines in the per-frame path; `$C20003` runs only when `$E9EF == 0`. |
| `~$C101CC–$C101FA` | Tentative | Unexplored routine: copies 5 bytes `$2E72→$602D+`, counter `$64DA`, ends `RTL`. Noticed in the caller dump. |
| `ROMCPU:$C21338`, `$C21347`, `$C21396` | Strong hypothesis | Battle delta store sites (heal clamp / damage / death zero) inferred from +3 post-instruction write-callback PCs; routine entries `$C21323/$C21350/$C21390` claimed in `battle.go`, dump not preserved. [SESSION_003](SESSION_003.md). |
| `ROMCPU:$C212FF` | Strong hypothesis | Dispatching `JSR (abs,X)` (stack return `$1302` at every delta store capture). |
| `ROMCPU:$C0567B` | **Confirmed** (writes) | `$FF` region fill hitting `WRAM:+$2E78` and `+$3BF4` regions (boot and battle teardown; `DB=$00`, large X counters). Routine bounds unexplored. |
| `ROMCPU:$C223F6`, `$C227B4`, `$C22408` | **Confirmed** (writes) | Battle-init writers of `WRAM:+$3BF4`/`+$2E78` arrays (slot 3 sentinels `$0000`/`$FFFF` for a 3-member party). Contexts unexplored. |
| `ROMCPU:$C206BC/$C206BF`, `$C36A58/$C36A5E`, `$C25D33` | Tentative | Additional zeroing/clearing writers of the HP arrays; contexts unknown (bank `$C3` pair fired outside the battle frame range, `DB=$00`). |
