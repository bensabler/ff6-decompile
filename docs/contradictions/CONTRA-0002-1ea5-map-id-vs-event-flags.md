# CONTRA-0002: `WRAM:+$1EA5` — "map id" vs "map-load target" (both refuted)

- **Status:** **Resolved 2026-08-01** — both competing claims are
  refuted by a single static decode. The byte is neither: it is one
  byte of an **event-flag bit array**.
- **ROM revision:** SHA-256 `0f51b4fca41b7fd509e4b8f9d543151f68efa5e97b08493e4b2a0c06f5d8d5e2`
- **Method:** static ROM decode only. No emulator run was required, so
  no emulator identity is in play.
- **Dependent work frozen during resolution:** `internal/scenario/route`
  (`AddrMapContextCandidate`, the `UntilMapChange`/`MapValue` leg
  kind, `State.MapContext`) and the `MAPC` logging in
  `mesen/probes/EXP-0036.lua`. Nothing on the shipped golden route
  depended on the *semantics* — EXP-0036 had already moved transition
  detection to the player-position jump — so the freeze cost nothing
  and no route result is invalidated.

## The competing claims

**Claim A — "candidate map id" (EXP-0035).**
`+$1EA5` reads `$00` on opening-approach states, `$05` on all four
Narshe-exterior free-walk states, and `$0D` on both mines-interior
states, across eight independently produced savestates. Recorded at
Tentative, with the acknowledged tension that milestone 02 is visually
the Narshe exterior yet reads `$00`.

**Claim B — "map-load target / pending-map or event-state value"
(EXP-0036).** During the scheduled mines transition the byte reaches
`$0D` while the party is still visibly standing on the exterior and has
not moved, so it cannot denote the *current* map. Proposed instead as
the map the event has *decided on*, written by `ROMCPU:$C0B5B6` ahead
of the visible transition.

Both claims agree the byte is map-related and differ only on *when* it
denotes the map. That framing is what made them look like the only two
options.

## Conflict type

Per the resolver's taxonomy this is an **interpretation** conflict, and
the shared error is a **sampling artifact**: EXP-0035 sampled only
*settled* states, where a location value and a progress value are
indistinguishable; EXP-0036 sampled *during* a transition, which
refuted "current" but kept the map framing rather than questioning it.
No ROM, address-space, or emulator-version difference is involved.

## Discriminator

The smallest discriminating test was not another emulator run but a
**static decode of the writer**, plus a scan for every instruction that
references the address. Two facts fell out immediately:

1. **No instruction anywhere in the ROM references `$1EA5` directly.**
   A scan for absolute/long operands equal to `$1EA5` under every
   plausible opcode returns six sites, none of which is the observed
   writer and none of which is a plain `LDA $1EA5`. A standalone
   variable would be addressed directly; this one never is.
2. **The observed writer is an indexed read-modify-write.** The probe
   reports the PC *after* the storing instruction, so the live
   `$C0B5B6` corresponds to the store at `ROMCPU:$C0B5B3`:

   ```text
   $C0B5B1  LDA $EB               ; flag number
   $C0B5B3  JSR $BAED             ; -> Y = flag/8, X = flag&7
   $C0B5AC  LDA $1EA0,Y
   $C0B5AF  ORA $C0BAFC,X         ; OR in a single-bit mask
   $C0B5B3  STA $1EA0,Y           ; <- the observed write
   ```

## Evidence

| Item | Value |
|---|---|
| Set-mask table `ROMCPU:$C0BAFC` | `01 02 04 08 10 20 40 80` — single-bit set masks |
| Clear-mask table `ROMCPU:$C0BB04` | `FE FD FB F7 EF DF BF 7F` — exact complements |
| Index decoder `ROMCPU:$BAED` | `REP #$20 / TAX / LSR ×3 / TAY / TDC / SEP #$20 / TXA / AND #$07 / TAX / RTS` — flag number → Y = number/8, X = number&7 |
| Array bases | `WRAM:+$1E80`, `+$1EA0`, `+$1EC0` — three parallel arrays, `$20` bytes apart (256 flags each) |
| Access sites | `LDA $1EA0,Y` ×5, `STA $1EA0,Y` ×4; symmetric counts for the `$1E80` and `$1EC0` arrays |

**The decisive observation:** the recorded values **accumulate bits
monotonically** rather than replacing a value.

```text
$00 = 0b00000000   no bits
$01 = 0b00000001   bit 0
$05 = 0b00000101   bits 0, 2
$0D = 0b00001101   bits 0, 2, 3
```

A map identifier does not gain a bit each time the story advances. A
bit-flag byte does, and that is exactly what the set routine performs.

## Resolution

`WRAM:+$1EA5` is **byte 5 of the event-flag bit array based at
`WRAM:+$1EA0`**, covering flags `$28`–`$2F`. At milestone 05 the bits
set are 0, 2 and 3 — flags `$28`, `$2A`, `$2B`.

**Both Claim A and Claim B are refuted**, including EXP-0036's own
proposal. The apparent correlation with location was real but
incidental: event flags accumulate with **story progress**, and in a
linear opening story progress correlates almost perfectly with where
the party is. Milestone 02's `$00` — the tension EXP-0035 flagged and
could not explain — is simply a point before any of these flags had
been set, and is now fully explained rather than tolerated.

- **Confidence:** **Confirmed.** Instruction-level decode of the live
  writer, both mask tables read from ROM, the index decoder decoded in
  full, and the value progression independently consistent with
  bit-setting. No competing reading survives.

## Supersession

| Record | Change |
|---|---|
| EXP-0035 | Its `+$1EA5` "candidate map id" reading is **superseded**; the record now points here. Its *position* findings (`+$00AF`/`+$00B0`) are untouched and stand. |
| EXP-0036 | Its "map-load target / event-state" reading is **superseded** by this record. Its route results, milestone 05, and the decision to detect the transition by player position all stand — that decision was correct, and is now correct for the right reason. |
| CEN-WORLD-0007 | The map-id candidate is removed; the entry keeps only the position bytes and alias findings. |
| CEN-EVENT-0008 | New entry registering the event-flag system located by this resolution. |

## Consequences for dependent work

- `internal/scenario/route` identifiers that encoded the refuted
  reading are renamed to evidence-neutral ones: the leg condition
  becomes a generic "watched byte equals value" rather than a
  "map change", and the address constant is re-pointed at the array
  base with its true meaning.
- No golden-route result changes. Milestone 05 remains established on
  three byte-identical runs.
- **Unfrozen** on publication of this record.

## Exact next action

The event-flag system is now located, which is a larger prize than the
map question that started this. Follow up with a bounded unit that
watches `$1EA0,Y` writes across the scheduled route to map **which
flag numbers the opening sets and when** — that directly serves the
scenario's event-flag, treasure-flag and persistence requirements
(B16, B19) and gives the event engine (CEN-EVENT-0001) a concrete
anchor. The map header / tileset question that motivated EXP-0037 is
now **untouched by this resolution** and still needs its own experiment;
`+$1EA5` was never going to answer it.
