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
| `$3BF4–$3C07` (`WRAM:+$3BF4`) | `$14` | **Confirmed** (writers, code role, authority, 10-entry extent) | Unified battle current-HP array: entries 0–3 party, 4–9 enemies (EXP-0005: enemy HP 24/35 observed at entries 4–5, damaged and zeroed by the same delta-engine stores `ROMCPU:$C21338/$C21347/$C21396`). Init `$C223F6`/`$C227B4`; `$FF` fill `$C0567B`; zeroing `$C206BC/BF`, bank-`$C3` `$C36A58/5E`, `$C22CCE` (unexplored). Authoritative: display derives from it via `$C25D26`. [SESSION_003](SESSION_003.md), EXP-0003/0004/0005. |
| `WRAM:+$3C08`, `+$3C1C`, `+$3C30` | 8 each | **Confirmed** (code roles, EXP-0001) | Per-slot arrays operated by the delta engine: `+$3C08` current pool of the MP-flagged path (clamped to `+$3C30`); `+$3C1C` clamp ceiling of the HP path. Semantic labels (MP/max-HP/max-MP) Strong hypothesis from battle context and parallel shape. |
| `WRAM:+$3C95` | 8 | **Confirmed** (code role) | Per-slot; bit 0 = call death handler when the `+$3C08` pool hits zero ("dies at zero MP" candidate). Other bits unknown. |
| `WRAM:+$3EE4` | 8 | **Confirmed** (code role) | Per-slot 16-bit; bit 1 suppresses the death event in the handler (`BIT #$0002` at `ROMCPU:$C2139C`). Also masked into display records (Session 002 `$2EA0` link unestablished). Other bits unknown. |
| `WRAM:+$33D0–$33E3`, `+$33E4–$33F7` | `$14` each | **Confirmed** (code roles); 10-entry extent Strong hypothesis (EXP-0004) | Pending-delta arrays read by the fetch at `ROMCPU:$C213A7`; `$FFFF` = none; delta = `+$33E4` − `+$33D0` (secondary gated). Transient: set by `ROMCPU:$C20C9B` (mid-battle setter, ×12 in one battle) and swept back to `$FFFF` by `ROMCPU:$C2638E`/`$C26391` (bulk, entry 9 observed); initialized by `$C22408`. |
| `ROMCPU:$C20C9B` | — | **Confirmed** (write site) | Pending-delta setter (wrote `+$33D4` ← `$0004`, slot 2); caller chain returns `$0436`/`$0C2A` — damage-formula layer entry, undumped. |
| `ROMCPU:$C21406–$C21411`, `$C2069B` | — | Strong hypothesis | Post-fetch driver sequence (EXP-0001 dump): `JSR $629B / JSR $069B / JSR $4C5B / JSR $1429`; the `JSR $069B` at `$C21409` is the consistent stack-top of every steady refresh call — `$C2069B` likely tail-jumps into PartyDisplaySourceRefresh (`$C25D26`). Undumped. |
| `WRAM:+$11A2` | 1 | **Confirmed** (code role) | Bit 7 selects the MP path in the delta dispatch (`ROMCPU:$C212F9`). Writers unknown. |
| `WRAM:+$3A89`, `+$3A3C`, `+$3A81/+$3A82`, `+$327C`, `+$3419` | ? | Tentative | Touched by the delta engine (cleared / gates / tail writes); meanings unknown. |
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
| `ROMCPU:$C212F5–$C213D2` | **Confirmed** (byte-exact, EXP-0001) | Battle delta engine: dispatch wrapper (`$C212F5`, selector `WRAM:+$11A2` bit 7, `JSR ($131F,X)` at `$C21300`, table `$C2131F` = `$1323/$1350`), HP routine `$C21323` (stores `$C21338`/`$C21347`), MP routine `$C21350`, death handler `$C21390` (store `$C21396`, `JMP $C20E32`), delta fetch `$C213A7`. Full listing in [02_DISCOVERED_FUNCTIONS.md](02_DISCOVERED_FUNCTIONS.md). |
| `ROMCPU:$C0567B` | **Confirmed** (writes) | `$FF` region fill hitting `WRAM:+$2E78` and `+$3BF4` regions (boot and battle teardown; `DB=$00`, large X counters). Routine bounds unexplored. |
| `ROMCPU:$C223F6`, `$C227B4`, `$C22408` | **Confirmed** (writes) | Battle-init writers of `WRAM:+$3BF4`/`+$2E78` arrays (slot 3 sentinels `$0000`/`$FFFF` for a 3-member party). Contexts unexplored. |
| `ROMCPU:$C206BC/$C206BF`, `$C36A58/$C36A5E` | Tentative | Additional zeroing/clearing writers of the HP arrays; contexts unknown (bank `$C3` pair fired outside the battle frame range, `DB=$00`). |
| `ROMCPU:$C25D26–$C25D56` | **Confirmed** (byte-exact, EXP-0003) | PartyDisplaySourceRefresh: copies all six authoritative battle arrays (`+$3BF4/+$3C1C/+$3C08/+$3C30/+$3EE4/+$3EF8`) into the display-source arrays (`+$2E78/+$2E80/+$2E88/+$2E90/+$2E98/+$2EA0`), slots 3..0. Caller unknown; event-driven (~42 calls per observed battle). |
| `WRAM:+$3EF8` | **Confirmed** (code role) | Per-slot 16-bit; copied to display `+$2EA0`; bit 13 drives the `+$61AD` slot mask downstream. Meaning unknown. |
