# FF6 Reverse-Engineering Session 002

- **Date:** 2026-07-29
- **Investigator:** Benjamin Sabler + Claude (automated via Lua bridge)
- **ROM identity/checksum:** `Final Fantasy III (USA).sfc` (Desktop copy);
  checksum not recorded
- **Mesen CE version:** Mesen 2.1.1 (macOS x64 Intel)
- **Goal:** Answer open question #1 — who calls `$C10DF3` and in what
  context — and validate the Session 001 structure hypotheses.

## Starting state

Save state `Final Fantasy III (USA)_11.mss`: Narshe intro battle, 3-member
Magitek party (`?????`/WEDGE/VICKS), Magitek command menu open. HP
53/39/33.

## Method

New tooling: [mesen/bridge.lua](../../mesen/bridge.lua) — a Lua
script loaded with the ROM from the command line. It logs exec-callback
captures at `$C10DF3` (registers + stack snapshot) to
`tools/mesen/out/events.log` and polls `out/cmd.txt` for commands
(`screenshot`, `read`, `press`, `loadstate`, `eval`), so the session can be
driven and observed programmatically while the GUI stays visible.
Requirements discovered: `AllowIoOsAccess` must be enabled in Mesen's
script settings; `emu.loadSavestate` must run inside a main-CPU exec
callback; use `emu.getCpuState(emu.cpuType.snes)` for registers (the
generic `getState()` does not expose them).

## Experiment

1. Launch Mesen with ROM + bridge; confirm zero `$C10DF3` hits at title.
2. Load the battle save state; log exec captures at `$C10DF3` and `$C10E14`.
3. Dump ROM `$C10DF3` (116 bytes) and `$C101C0` (112 bytes); disassemble by
   hand; verify every branch offset arithmetically.
4. Live-read `$628D`, `$E9EF`, `$61AD`, `$2E78–$2EA7`, `$2EB5+`.
5. Play: cast Heal Force (input injection), let enemies act; re-read memory
   and screenshot the HP window.

## Raw observations

- Title screen: 0 hits. In battle: 1 hit/frame, `ret=$C1:0203` every time.
- Entry: `A=$0000 X=$0002 Y=$000C SP=$15E7 PS=$26 (m=1,x=0,e=false) K=$C1
  DB=$7E D=$0000`.
- Stack `SP+1..`: `02 02 | 0C 00 02 00 | 28 64 C2 | …` — RTS return
  `$C10203`, pushed Y/X, then JSL return `$C2:6429`.
- At `$C10E14` (three consecutive iterations): `PS=$04 (m=0,x=0)`;
  `A=$0035 X=0 Y=0`, `A=$0027 X=2 Y=$20`, `A=$0021 X=4 Y=$40`.
- ROM dumps: see disassemblies in
  [02_DISCOVERED_FUNCTIONS.md](02_DISCOVERED_FUNCTIONS.md).
- Live WRAM: `$628D=$00 $E9EF=$00 $61AD=$0F`; source/destination dumps in
  [05_DATA_STRUCTURES.md](05_DATA_STRUCTURES.md).
- Battle HP window: `?????`=53, WEDGE=39, VICKS=33 — equal to `$2E78[0..2]`.
- After Heal Force on slot 0: `$2E78[0]` 53→63 = `$2E80[0]` exactly.
  Enemy attacks: slot 2 33→29→3, slot 1 39→29, mirrored in records next
  frame. Record 2 read `03 00 46 00 …` (HP=3, max=70).
- Record 0 bytes `+$0C..$1F`: `08 00 00 00 00 00 00 00 FF FF FF FF FF 0E 96
  84 83 86 84 FF` — written by something other than this routine.

## Findings

1. **Caller identified (Confirmed):** `JSR $0DF3` at `$C10200` inside
   `$C101FB` (PerFrameBattleUpdate, candidate name), entered via JSL from
   bank `$C2`; once per frame during battle, never at title.
2. **Routine fully decoded (Confirmed):** six-field copy + conditional
   masking (`$628D`/`$E9EF` gates; masks `$0038`/`$0000`) + `$61AD` slot
   mask (bit n = NOT bit13 of `$2EA0[n]`). Session 001 had seen only the
   first four field copies.
3. **Slot correspondence (Confirmed for 0–2):** records and source arrays
   are the displayed party slots, in order.
4. **`$2E80` = max HP (Strong hypothesis):** heal snapped current HP to
   exactly this value; all entries ≥ current HP; HP gauges displayed.
5. **Flags/width/DB (Confirmed):** `m=0,x=0` at stores; X/Y start 0;
   `DB=$7E`; `D=$0000`; native mode.

## Documentation changes

- `02`, `03`, `04`, `05`, `06`, `08` all updated; this note added.

## Go changes

- `chardata` extended: six-field `CharacterFieldsSource`, `PartySlotRecord`
  with `+$0C..$1F` preserved, `CopyMode` (two flags), `CopyCharacterFields`
  now applies masking and returns the `$61AD` slot mask. Tests updated,
  including a case encoding the live battle capture.

## Verification

```text
gofmt -l .:     clean
go build ./...: pass
go vet ./...:   pass
go test ./...:  pass (chardata)
```

## Open questions

- See [08_OPEN_QUESTIONS.md](08_OPEN_QUESTIONS.md) — five answered,
  twelve now open (producer of `$2E78`, consumers of records/`$61AD`,
  flag meanings, the eight sibling subroutines, bank-`$C2` driver, etc.).

## Next experiment

Write breakpoint on `$2E78` (WRAM) during an enemy damage event, capture
the call stack — that finds the code that *computes* HP changes (the
authoritative battle logic upstream of this display copy). Second choice:
read breakpoint on `$2EB5` to find the records' consumer.
