# EXP-0024: Battle confirm SFX — trigger, APU command, voice, sample

- **Status:** completed (2026-07-30)
- **Question (Unit 5, audio vertical):** what does the main CPU write
  to the APU ports (`$2140-$2143`) when the battle-menu confirm SFX
  triggers, which DSP voice keys on for it, and which ARAM sample
  (SRCN → directory → BRR address) does it play — with ROM provenance
  if the sample is raw-stored?
- **Starting state:** headless bridge session; `exp10-battle.mss`
  reloaded per trial; frame-scheduled input (EXP-0021 machinery).
- **Method:** tracked probe `mesen/probes/EXP-0024.lua`:
  - Write watch on `snesMemory $00:2140-$00:2143` logging PC, value,
    frame (capped), active from load.
  - Per-frame DSP voice-activity sampling (ENVX per voice via the
    `dsp` read type) in a window around the press, logging activation
    transitions.
  - Trial P (press): load, press A 10 frames at anchor+72.
    Trial N (no press): load, no input, same window — the port-write
    and voice-activation **delta** between trials isolates the SFX
    from music traffic.
  - Then targeted reads: full DSP register file, the keyed voice's
    SRCN, sample-directory entry (`DIR*$100 + 4*SRCN` in ARAM), first
    BRR blocks; ROM byte-search of the BRR span for provenance.
- **Expected outcomes:**
  - *Clean trigger:* trial P shows a port-write burst near the press
    absent from trial N, plus a new voice activation; SRCN resolves to
    a directory entry and a BRR sample in ARAM; ROM search hits →
    full vertical (trigger → port command → voice → sample → ROM).
  - *Partial:* trigger and voice isolate but the sample bytes miss in
    ROM (transformed/loaded differently) → record the negative
    provenance result; the trigger/voice/ARAM chain still stands.
- **Falsifying outcome (for the trigger claim):** no port-write or
  voice-activation delta between P and N trials (the confirm SFX is
  not CPU-triggered in this window, or the press produced no SFX).
- **Bounds:** if isolating the SFX requires decoding the SPC driver's
  command protocol beyond direct port/DSP observation, stop and record
  the boundary (follow-up unit). Music traffic is expected background;
  only the P−N delta is interpreted.
- **Raw evidence paths:** `local_artifacts/experiments/EXP-0024/`
  (ports/voices logs, DSP snapshot, directory + sample + pack dumps,
  brr-info; hashes in the asset record).
- **Result:** **clean trigger; the full vertical closes.** Detail in
  [AUD-0001](../audio/AUD-0001-battle-confirm-sfx.md).
  - Trial P (press at anchor+72): exactly **one** in-window port write
    — `$21 → $2140` from `ROMCPU:$C117CC` at rel 74. Trial N (no
    press): **zero**. Falsifier not met.
  - Voice delta (per-frame ENVX bitmaps, P xor N): **voice 7 only**,
    active rel 76–79. Music voices replayed identically across trials
    — APU-side schedule determinism, consistent with EXP-0021.
  - DSP snapshot at rel 77 (third trial, eval-armed dumper): voice 7
    SRCN=$05, PITCH=$2AA2, VOL $22/$22, ENVX=$7F.
  - Directory (DIR=$1B → ARAM $1B00): entry 5 start=loop=$48D8; the
    sample is 2 BRR blocks / 18 bytes (block 2 header $97:
    range 9, filter 1, LOOP+END).
  - **ROM provenance Confirmed:** unique byte-search hit at
    `ROMFILE:0x051FA1`; the contiguous SFX pack `ARAM:$4800-$491F`
    (288 bytes, samples 0–7) is byte-identical to
    `ROMFILE:0x051EC9-0x051FE8`.
  - Go proof: `brr.Decode` implemented (S-DSP semantics, tests +
    fuzz) + `ff6lab brr info`; the captured click decodes to 32 PCM
    samples (silent lead-in, sweep to −9582).
- **Confidence:** trigger chain and pack ROM provenance —
  **Confirmed**; `$21` command semantics, SPC driver dispatch,
  background port protocol — Unknown; music-voice attribution —
  Tentative.
- **Next action:** follow-ups queued — SPC-side dispatch trace (port 0
  `$21` → voice allocation), background `$E4/id/$18` protocol,
  remaining pack samples' cues.
