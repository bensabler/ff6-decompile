# AUD-0001: Battle-menu confirm SFX (click)

## Cue/trigger
A-press (menu confirm) in the battle savestate `exp10-battle.mss`.
Press at anchor+72 → APU port write 2 frames later → DSP voice 7 keys
on 2 frames after that, sounding 4 frames (rel 76–79). Established by
press-vs-no-press delta over identical windows (EXP-0024 trials P/N:
1 vs 0 port writes; voice-7 activity in P only).

## Local paths and hashes
`local_artifacts/experiments/EXP-0024/` — exp24-ports.log
(`dced842c…d935a5`), exp24-voices.log (`ec86140d…53a2d0`),
exp24-dsp.hex (`1db54893…a13aed`), sample5_48D8.hex
(`cae3d195…65480b`), aram_4800_288.hex (`cd3c7dc8…d6a861`),
directory_1B00.txt (`e45e5bfc…adb064`), brr-info.txt
(`c89e7d6b…361547`), hashes.sha256.

## CPU/APU protocol
Confirm SFX command: **one write of `$21` to `$2140` (APU port 0)
from `ROMCPU:$C117CC`**. Background traffic observed separately
(periodic `$E4→$2142, id→$2141, $18→$2140` triples from
`$C11811/$C11818/$C1181F`, and `$28→$2140` from `$C117E8`) — music/
engine heartbeat, uninterpreted (Unknown). The `$21` command byte's
encoding (direct SFX id vs indexed command) is Unknown pending driver
decode.

## SPC driver path
Not traced (bounded out per plan) — the SPC-side dispatch from port 0
to voice allocation is a follow-up unit.

## Sequence
Not applicable at this level (single-shot SFX; no sequence decoded).

## Instruments
Sample directory at ARAM `$1B00` (DSP `DIR=$1B`). Entry 5 (SFX):
start=loop=`ARAM:$48D8`.

## BRR samples
Sample 5: **2 BRR blocks, 18 bytes** (`ARAM:$48D8-$48E9`), block 2
header `$97` (range 9, filter 1, LOOP+END). Decoded by
`internal/audio/brr.Decode` / `ff6lab brr info`: 32 PCM samples,
silent first block, second block sweeps to min −9582 / max 1536 — a
short click. **ROM provenance: the contiguous SFX sample pack
`ARAM:$4800-$491F` (samples 0–7, 288 bytes verified) is byte-identical
to `ROMFILE:0x051EC9-0x051FE8`** (sample 5 at `0x051FA1`), found by
unique byte search + full-span comparison.

## DSP state
At mid-SFX (rel 77): voice 7 `VOLL/VOLR=$22/$22, PITCH=$2AA2,
SRCN=$05, ADSR=$FFE0, ENVX=$7F, OUTX=$DD`. Music voices 2/5 active
concurrently (SRCN `$23/$22`). Full 128-byte register file archived.

## Timing
Press rel 72 → port write rel 74 → voice-7 active rel 76–79 (4
frames). Frame-exact (headless testrunner, frame-scheduled input).

## Reconstruction recipe
1. Read 18 bytes at `ROMFILE:0x051FA1` (or the pack from `0x051EC9`).
2. Decode with `brr.Decode` (filter/range semantics per S-DSP).
3. Play at the captured pitch (`$2AA2`) and volume for the in-game
   sound.

## Validation
Press/no-press delta isolates trigger and voice; ARAM↔ROM byte
identity over the 288-byte pack; Go decode of the captured sample
(brr-info.txt archived). EXP-0024 falsifier (no delta) not met.

## Confidence
Trigger chain (press → `$21`@`$2140` from `$C117CC` → voice 7 →
SRCN 5 → `ARAM:$48D8`) — **Confirmed** (delta trials + DSP snapshot).
ROM provenance of the SFX sample pack — **Confirmed** (byte-exact,
288-byte span). `$21` command semantics, SPC driver dispatch, and the
background port protocol — Unknown. Music-voice attribution —
Tentative (concurrent activity, not delta-isolated).
