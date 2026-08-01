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
| `ROMCPU:$C20C76–$C20C9D` | — | **Confirmed** (byte-exact, EXP-0006) | PendingDeltaAccumulate: DP `$F0` amount += slot pending (`$FFFF`→0), **capped at `$270F` (9999)**, store `$C20C98`; prelude retargets `+$33D0,Y`→`+$33E4,Y` (`Y += $14`) via DP `$F2`/carry polarity. Caller: `JSR $0C2D` at `$C20C28` → gate block `$C20C2D` (`+$11A2` bit 0, `+$3A82`&`+$3A83`, `+$3EE4,X` bit 1). Amount provenance = question #21. |
| `ROMCPU:$C23469` | — | **Confirmed** (EXP-0009) | `JSR $0B83` on the accumulator call path (pushes `$346B`, matching all EXP-0007 stacks); sits in a ten-slot descending loop (`LDY #$12`, −2/iter) gated by `+$3018,Y` bits vs DP `$A4` — target-loop candidate (Tentative). |
| `ROMCPU:$C20B83–$C20C2C` | — | Strong hypothesis (bracketing) | The damage-formula body: called from `$C23469`, flows into the decoded gate region (`$C20C0B`–`$C20C2C`) that queues DP `$F0`. Undumped interior = EXP-0010. |
| `WRAM:+$3018` | `$14`? | Strong hypothesis | Per-slot bit table (bit n ↔ slot n candidate): read at the fetch gate (`$C213AD`), the dispatch tail (`$C21316`), and the target loop (`$C23442` region). Layout/writers unknown. |
| `WRAM:+$3BCC`, `+$3BE0` | `$14` each | **Confirmed** (code roles, EXP-0010) | Per-slot 16-bit element-response masks (`$14`-stride family): `+$3BCC,Y` low = flip-to-heal (absorb cand.), `+$3BCD,Y` high = zero (immune cand.); `+$3BE0,Y` low = double (weak cand.), `+$3BE1,Y` high = halve (resist cand.). Writers unknown. |
| `WRAM:+$11A1`, `+$11A4`, `+$11A6`, `+$11AA`, `+$3EC8` | 1 each | **Confirmed** (code roles) | Current-action state: `+$11A1` element byte; `+$11A4` mode/polarity source (bit 7 selects base-formula variant, bit 1 in the flip chain); `+$11A6` pipeline gate; `+$11AA` abort bits (`#$82`); `+$3EC8` battle-wide element-nullify byte. Meanings beyond code roles unknown. |
| `ROMCPU:$C22966–$C22989` | — | **Confirmed** (byte-exact, EXP-0019) | Attack-record loader: `MVN $C4,$7E` copies the 14-byte entry at `ROMCPU:$C46AC0 + 14×index` (`ROMFILE:0x046AC0+`) into `WRAM:+$11A0`–`+$11AD`. |
| `ROMCPU:$C2299F–$C229ED+` | — | **Confirmed** (code, EXP-0019) | Fight-command/attacker staging: `+$11AE`←`+$3B2C,X`, `+$11AF`←`+$3B18,X`, power `+$11A6`←`+$3B68,X`, element `+$11A1`←`+$3B90,X`, `+$11A8`←`+$3B7C,X`; gates from `+$3C45,X`/`+$3BA4,X`. |
| `WRAM:+$3B18, +$3B2C, +$3B68, +$3B7C, +$3B90, +$3BA4` (+`+$3C45` grid-adjacent) | `$14` each | **Confirmed** (code roles) | Per-slot battler stat tables feeding action setup (statB, statA, battle power, ?, element, mode bits). Labels (level/vigor/weapon-power…) unproven. Family ≈19 members. |
| `ROMCPU:$C20C9E`–`$C20D86`, `$C20D87`+ | — | **Confirmed** (decoded stretches, EXP-0011) | Base-amount routines: variant A consumes the precomputed base at `WRAM:+$11B0`, applies defense (`+$3BB8,Y` 16-bit pair, (255−def)/256 shape), flag halvings, party-vs-party halving, final `$C2370B` transform; boost sibling `$C20D4A` (+≈50%, `+$3C44,X` flags); variant B = fraction-of-HP (min 1). Helper `$C20DDD`: `+$11A3` bit 7 → `X/Y += $14` (HP→MP retarget). Helpers `$C247B7`/`$C2370B` unexplored (#24/#25). |
| `WRAM:+$11B0` | 2 | **Confirmed** | Base amount: written by the standard-path computation `ROMCPU:$C22B69`–`$C22B9C` (**base = `+$11A6`×4 + (`+$11A6`×`+$11AE`×`+$11AF`)>>5**, EXP-0015, numerically closed live: 60/28/4 → 450), enemy/physical path `$C22BE9` store, per-target staging `$C23422`, init writers `$C25561`/`$C20E8A`/`$C210E1`; boosted in place by `$C20D4A`. |
| `WRAM:+$11AE`, `+$11AF` | 1 each | **Confirmed** (code roles) | Stat operands of the base formula (28 and 4 live). Gameplay meanings and producers unknown — question #27. |
| `ROMCPU:$C20DD1` | — | **Confirmed** (byte-exact) | Shift helper: 24-bit full product (`$EA:$E9:$EA` composition from the `$C247B7` wrapper) >> (A+1), 16 result bits. `$C20DCB` entry = wrapper call + shift 4. |
| `WRAM:+$3BB8`, `+$3C44` | `$14` each | **Confirmed** (code roles, EXP-0011) | More `$14`-stride family arrays: `+$3BB8,Y` 16-bit defense pair (`$FF`=none; physical/magical by carry path); `+$3C44,X` attacker boost flags. Twelve family members now identified (see 05). |
| `ROMCPU:$C21406–$C21411`, `$C2069B` | — | Strong hypothesis | Post-fetch driver sequence (EXP-0001 dump): `JSR $629B / JSR $069B / JSR $4C5B / JSR $1429`; the `JSR $069B` at `$C21409` is the consistent stack-top of every steady refresh call — `$C2069B` likely tail-jumps into PartyDisplaySourceRefresh (`$C25D26`). Undumped. |
| `WRAM:+$11A2` | 1 | **Confirmed** (code role) | Bit 7 selects the MP path in the delta dispatch (`ROMCPU:$C212F9`). Writers unknown. |
| `WRAM:+$3A89`, `+$3A3C`, `+$3A81/+$3A82`, `+$327C`, `+$3419` | ? | Tentative | Touched by the delta engine (cleared / gates / tail writes); meanings unknown. |
| `$628D` | 1 | **Confirmed** (use) | Masking gate flag; meaning unknown. |
| `$64DA` | 1? | Tentative | Counter (low nibble used ×5 as index), incremented by the unexplored 5-byte copy routine. |
| `$E9EF` | 1 | **Confirmed** (use) | Masking gate / per-frame branch flag; meaning unknown. |
| `WRAM:+$1D4D` | 1 | **Confirmed** (EXP-0041) | Config settings byte 1. Bits 0–2 Bat.Speed (0–5, displayed 1–6); bit 3 Bat.Mode (0 = Active, 1 = **Wait**); bits 4–6 Msg.Speed (0–5, displayed 1–6); bit 7 Cmd.Set (0 = Window, 1 = Short). New-game default `$2A`. Both speed fields swept to their clamps. Writer routine unknown. |
| `WRAM:+$1D4E` | 1 | **Confirmed** (EXP-0041) | Config settings byte 2. Bit 4 Reequip (0 = Optimum), bit 5 Sound (0 = Stereo), bit 6 Cursor (0 = Reset, 1 = Memory), bit 7 Gauge (0 = On). Default `$00`. Bits 0–3 untouched by the nine Config settings; meaning unknown. |
| `WRAM:+$1D54` | 1 | **Confirmed** (EXP-0041) | Config settings byte 3. Bit 7 Controller (0 = Single, 1 = Multiple). Default `$00`. Bits 0–6 untouched by the nine Config settings; meaning unknown. Not adjacent to `+$1D4D`/`+$1D4E` — the configuration block is **not** contiguous. |
| `WRAM:+$39A6`, `+$39B6`, `+$3A2E`, `+$3A32`, `+$3AAE`, `+$3AB2`, `+$3B26`, `+$3B36`, `+$3BA6`, `+$3BB6`, `+$3C26`, `+$3C36`, `+$3CA6`, `+$3CB6`, `+$3D26`, `+$3D36`, `+$3DA6`, `+$3DB6` | — | **Confirmed** (EXP-0041) | Config screen text cells, character byte interleaved with a tile attribute byte. The **selected** option's cells carry attribute `$20`, the unselected `$28` — the inverse of the intuitive reading, and the reason EXP-0040 misread `Bat.Mode` from a screenshot. Renderer unmapped. |
| `WRAM:+$3A8F` | 1 | **Confirmed** (EXP-0042) | Battle-local **Bat.Mode** flag. Zeroed at battle entry and `INC`'d by `ROMCPU:$C2247E` iff `+$1D4D` bit 3 is set: `01` = Wait, `00` = Active, both observed. Sampled **once** at entry; not re-read from configuration during the battle. Consumers unknown. |
| `WRAM:+$3A90` | 1 | **Confirmed** (EXP-0042) | Battle-local **Bat.Speed** value, written at battle entry by `ROMCPU:$C22481-$C2248D` as `255 − 24 × BatSpeed` (stored speed 0–5): Fast → `$FF`, default 3 → `$CF`, Slow → `$87`. Arithmetic decoded statically, then predicted and matched at two configurations. **Consumer unknown** — the sharpest current lead into ATB rate; direction (larger = faster) is a Tentative hypothesis only. |
| `WRAM:+$2F2E`, `+$2021`, `+$2F34` | 1 each | **Confirmed** (code role, EXP-0042) | Further battle-entry destinations of configuration: `+$2F2E` cleared when Cmd.Set = Window (`$C22477`); `+$2021` cleared when Gauge = Off (`$C22495`); `+$2F34` ← `+$1D4E` bits 0–2 (`$C10FF7`, bits no Config setting touches). Consumers unknown. |
| `WRAM:+$890F–+$896A` | `$5C` | **Confirmed** (code role, EXP-0042) | Cursor-memory block, zeroed by `ROMCPU:$C159DE` when `+$1D4E` bit 6 (Cursor) is clear = Reset. The mechanism behind `Cursor = Memory` reopening a character's last-used ability (EXP-0040). |
| `ROMCPU:$C22472–$C22497` | — | **Confirmed** (byte-exact, EXP-0042) | Battle-entry configuration sampler: reads `+$1D4D`/`+$1D4E` once each and **decomposes** them into `$2F2E`, `$3A8F`, `$3A90`, `$2021`. |
| `ROMCPU:$C198AC–$C198C3`, `$C159D6–$C159EA` | — | **Confirmed** (byte-exact, EXP-0042) | The two in-battle readers of the persistent config bytes. `$C198AC` extracts Msg.Speed (`+$1D4D` bits 4–6) and indexes a delay table at `ROMCPU:$C19872`; `$C159D6` tests Cursor (`+$1D4E` bit 6) and clears `+$890F`. Presentation only — neither feeds battle timing. |

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
