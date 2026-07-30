# FF6 Reverse-Engineering Session 003 (reconstructed record)

- **Date:** 2026-07-29, ~19:34–20:47 local (bridge log header 19:34:37;
  `battle/` package written 20:47)
- **Investigator:** Benjamin Sabler + Claude (automated via Lua bridge)
- **ROM identity:** SHA-256
  `0f51b4fca41b7fd509e4b8f9d543151f68efa5e97b08493e4b2a0c06f5d8d5e2`
  ([ROM_IDENTITY.md](../research/ROM_IDENTITY.md)) — recorded at migration,
  same local file the bridge loaded
- **Mesen:** 2.1.1 macOS x86_64 (per [SESSION_002](SESSION_002.md))
- **Goal (inferred):** answer open question #1 — find the code that
  *computes* HP changes upstream of the CopyCharacterFields display copy.

> **Provenance notice.** This record was reconstructed on 2026-07-29,
> *after* the session, from surviving raw evidence. The session was
> interrupted by the Version 4 installation before any record was written.
> Two consequences, stated once here and reflected in every confidence
> below:
>
> 1. The `NEW-3BF4-WRITER` write watch that produced the central log lines
>    is **not** in the preserved `mesen/bridge.lua` (which watches
>    `WRAM:+$2E78`–`+$2E7F`); it was injected at runtime, almost certainly
>    via the bridge `eval` command (`mesen/out/resp.txt` still holds an
>    `eval` acknowledgment). Its exact code and watch range are lost.
> 2. The byte-exact disassembly claims embedded in
>    `internal/game/battle/battle.go` (routines `ROMCPU:$C21323`,
>    `$C21350`, `$C21390`, `$C213A7`) have **no preserved ROM dumps**.
>    They are treated as claims pending re-verification (see
>    "Required next experiments").

## Starting state

Battle save state loaded twice early in the log (Narshe sequence; slot-0
display HP `$2A` = 42), then later loads continuing into the Narshe mines
(`mesen/out/checkpoint1.mss` 19:47, `checkpoint2.mss` 20:13,
`checkpoint3-mines.mss` 20:46 — SHA-256s in
[V4_MIGRATION_REPORT.md](../migrations/V4_MIGRATION_REPORT.md) §3).

## Method

`mesen/bridge.lua` (exec log at `ROMCPU:$C10DF3`, write watch on
`WRAM:+$2E78`–`+$2E7F`) plus a runtime-injected write watch on the
`WRAM:+$3BF4` region. Battles were played (damage taken, heals cast, party
deaths) while every previously-unseen writing PC was logged with registers
and a stack snapshot.

## Raw observations (verbatim from `mesen/out/events.log`)

All lines preserved in the log (SHA-256
`bcfc7f4c…a99d03`). Selected, with `WRAM:` addresses decoded:

| Frame | Logged PC | Wrote | Value | Registers of note |
|---|---|---|---|---|
| 40 | `ROMCPU:$C0567B` | `WRAM:+$2E78` | `$FF` | `A=$00FF X=$3188 PS=$24 DB=$00` |
| 20930 | `ROMCPU:$C25D33` | `WRAM:+$2E7E` (slot 3) | `$00` | `A=$0000 X=$00FE PS=$17` |
| 28180 | `ROMCPU:$C0567B` | `WRAM:+$3BF4` | `$FF` | `A=$00FF X=$240C PS=$24 DB=$00` |
| 36566 | `ROMCPU:$C223F6` | `WRAM:+$3BFA` (slot 3) | `$00` | `A=$0012 X=$01DA Y=$0010 PS=$05` |
| 36566 | `ROMCPU:$C22408` | `WRAM:+$2E7E` (slot 3) | `$FF` | `A=$FFFF X=$046E Y=$0010 PS=$05` |
| 36567 | `ROMCPU:$C227B4` | `WRAM:+$3BF8` (slot 2) | `$43` | `A=$0043 X=$0004 Y=$022B PS=$04` |
| 46913 | `ROMCPU:$C2134A` | `WRAM:+$3BF4` (slot 0) | `$22` | `A=$0022 X=$0000 Y=$0000 PS=$15` |
| 47092 | (exec log) | — | — | `wram2E78=$0022 wram2EB5=$0022` |
| 107643 | `ROMCPU:$C2133B` | `WRAM:+$3BF8` (slot 2) | `$5D` | `A=$005D X=$0000 Y=$0004 PS=$15` |
| 142917 | `ROMCPU:$C21399` | `WRAM:+$3BF6` (slot 1) | `$00` | `A=$0000 X=$0000 Y=$0002 PS=$17` |
| 142917 | `ROMCPU:$C206BC`, `$C206BF` | `WRAM:+$3BF6`, `+$3BF7` | `$00`, `$00` | `A=$0388 PS=$B4` |
| 160467 | `ROMCPU:$C36A5E`, `$C36A58` | `WRAM:+$3BF4`, `+$3BF5` | `$00`, `$00` | `X≈$03AA Y≈$0116 PS=$05 DB=$00` |

Stack snapshot at the `$C2134A` and `$C2133B` and `$C21399` captures begins
`02 13 …` — a 16-bit return address `$1302` (bank `$C2` context).

The per-frame exec log also tracked `WRAM:+$2E78`/`+$2EB5` equal at every
sampled hit, declining through damage (`$2A→$22→$18→$14→$04→$00` across the
session's battles).

Unlabeled artifacts: `mesen/out/diff.txt`/`diff2.txt` (19:56–19:57) hold
`offset:old>new` byte diffs with offsets `$0014`–`$00B0`; what was diffed
was not recorded. Preserved but unusable as evidence.

## Interpretation (with alternatives)

1. **A store instruction near `ROMCPU:$C21347` decrements slot HP in a
   16-bit array at `WRAM:+$3BF4`.** The logged PC `$C2134A` is interpreted
   as the *post-instruction* PC of a 3-byte `STA` at `$C21347` — the same
   +3 pattern fits `$C2133B`→`$C21338` (heal path) and `$C21399`→`$C21396`
   (death zero), and matches the addresses independently claimed in
   `battle.go`. *Alternative:* the callback reports the exact store PC and
   `battle.go`'s addresses are wrong. **Discriminator:** ROM dump of
   `$C21300–$C21410` (next experiment).
2. **`WRAM:+$3BF4` is the gameplay-authoritative current-HP array** (slot
   stride 2): damage, heal, and death writers all target it, and its slot-0
   value `$22` appeared in the display-source array `WRAM:+$2E78` by frame
   47092. *Alternative:* both arrays are copies of a third, deeper store.
   The copier `$3BF4`→`$2E78` was **not** identified this session.
3. **Slot dispatch:** `Y` held `slot×2` at the delta stores (`0`, `4`,
   `2`); the `$1302` stack return is consistent with a dispatching
   `JSR (abs,X)` at `ROMCPU:$C212FF` (a 3-byte indirect-indexed `JSR`
   pushes `$C21301`). Matches `battle.go`'s claimed `JSR ($131F,X)`
   dispatch. *Alternative:* plain `JSR` at `$C212FF`; indirection table
   address unverified without the dump.
   > **Corrected by [EXP-0001](../experiments/EXP-0001-c2-delta-engine-dump.md)
   > (same night):** the pushed-return arithmetic above was off by one —
   > a 3-byte `JSR` pushing `$1302` sits at `ROMCPU:$C21300`, and the
   > dump confirms `JSR ($131F,X)` there. Every other claim in
   > interpretations 1–3 was verified byte-exact; see
   > [02_DISCOVERED_FUNCTIONS.md](02_DISCOVERED_FUNCTIONS.md).
4. **Lifecycle writers:** `ROMCPU:$C0567B` (`DB=$00`, large X counters)
   fills both regions with `$FF` — region-fill at battle teardown/boot.
   `ROMCPU:$C223F6`/`$C227B4`/`$C22408` initialize the arrays at battle
   start (slot 3 of a 3-member party gets `$0000`/`$FFFF` sentinels).
   `ROMCPU:$C206BC`/`$C206BF` (8-bit pair) and bank-`$C3`
   `ROMCPU:$C36A58`/`$C36A5E` (`DB=$00`, non-battle frame range) also
   zero HP entries — contexts unknown.

## Confidence

- Writers' logged PCs, written addresses, values, registers: **Confirmed**
  (raw log preserved).
- Store instructions at `$C21338`/`$C21347`/`$C21396`: **Strong
  hypothesis** (three independent +3-consistent captures + battle.go's
  claim; no dump preserved).
- `WRAM:+$3BF4` = authoritative battle current-HP array: **Strong
  hypothesis** (all delta writers target it; one observed propagation into
  `+$2E78`; propagation mechanism unidentified).
- Full routine disassembly, MP array `+$3C08`, max arrays `+$3C1C`/`+$3C30`,
  status `+$3EE4` bit 1, `+$3C95` bit 0, delta-fetch `$C213A7`,
  `+$33E4`/`+$33D0` sentinels, `+$3A89`: **claimed in `battle.go` only;
  Unknown pending re-dump.** No surviving artifact evidences them.

## Go changes (made during the original session)

`battle` package (now `internal/game/battle/battle.go`): `PartySlots`,
`ApplyHPDelta`, `ApplyMPDelta`, death handling. Written from the lost
live context; its heal-clamp/overflow/zero-floor arithmetic is consistent
with everything the surviving log shows but exceeds it. **No tests were
written during the session** — added with this record.

## Verification

Original session: not recorded. At documentation time (2026-07-29):
gofmt clean; `go build`/`go vet`/`go test ./...` pass.

## Required next experiments

1. **Re-dump `ROMCPU:$C21300–$C21410`** (and `$C213A7` tail through
   `$C213E0`) via the bridge; verify every instruction `battle.go` claims;
   this discriminates interpretation 1 and most of the Unknown block.
2. Watch `WRAM:+$3C08`/`+$3C1C`/`+$3C30` during MP spend/heal to evidence
   the MP/max arrays.
3. Find the `$3BF4`→`$2E78` copier (write watch on `+$2E78` during battle
   with stack capture — the preserved bridge already does this).
4. Identify `ROMCPU:$C206B9`± and bank-`$C3` `$C36A5x` zeroing contexts.
