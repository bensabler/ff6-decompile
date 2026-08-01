---
name: experiment-designer
description: Design minimal controlled emulator experiments that discriminate competing FF6 hypotheses.
---

# Experiment Designer

Specify question, state, independent variable, controlled variables, instrumentation, expected outcomes, falsifier, evidence paths, and stopping condition. Prefer one-variable changes and repeated trials.

Battle experiments (EXP-0041 onward) must also record the battle-configuration fingerprint under `## Battle configuration`, mirrored as `starting_state.battle_config` in `manifests/experiments.json`. Mark each value's `source`: `memory-read` when read from a located address, `screen-read` when read from the Config screen. A screen reading is not reliable — EXP-0040 misread `Bat.Mode` from a screenshot, which is how the ATB blocker was found.

Read and follow:

- `../_shared/EVIDENCE_STANDARD.md`
- `../_shared/ADDRESS_SPACES.md`
- `../_shared/STOPPING_RULES.md`
