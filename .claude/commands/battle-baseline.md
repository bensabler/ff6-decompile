Use the mesen-operator and experiment-designer skills.

Capture and validate the battle-configuration fingerprint before a battle
experiment operates Mesen. Read it from memory, never from the Config
screen — EXP-0040 misread `Bat.Mode` off a screenshot and that is how the
ATB blocker was found.

1. Read the persistent settings: `WRAM:+$1D4D` (bits 0-2 Bat.Speed 0-5,
   bit 3 Bat.Mode 1=Wait, bits 4-6 Msg.Speed 0-5, bit 7 Cmd.Set),
   `WRAM:+$1D4E` (bit 4 Reequip, 5 Sound, 6 Cursor, 7 Gauge),
   `WRAM:+$1D54` (bit 7 Controller). EXP-0041.
2. In a battle, also read the battle-local cells the entry sampler
   derived: `WRAM:+$3A8F` (Wait flag) and `WRAM:+$3A90`
   (`255 - 24 x Bat.Speed`). These, not the persistent bytes, are what
   battle timing runs on. EXP-0042.
3. Record ROM sha256, Mesen build and mode, savestate origin and hash,
   scenario and battle identifier, run identifier, and evidence-package
   identifier.
4. Assign the evidence directory per `docs/research/EVIDENCE_LAYOUT.md`
   and confirm it will not overwrite a frozen package.
5. Confirm required instrumentation is armed and no other Mesen instance
   is running.
6. **Reject an ambiguous start**: if the persistent bytes and the
   battle-local cells disagree, or a value cannot be read, stop and say
   so rather than proceeding.

Write the result into the experiment record's `## Battle configuration`
section and `starting_state.battle_config` in `manifests/experiments.json`
with `source: memory-read`. `internal/audit.CheckBattleExperimentConfig`
enforces its presence from EXP-0041 onward.

`mesen/probes/common.lua` provides `battleconfig()`, which returns the
whole fingerprint decoded.
