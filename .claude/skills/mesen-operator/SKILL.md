---
name: mesen-operator
description: Operate Mesen carefully for deterministic evidence collection without conflating viewer state with source data.
---

# Mesen Operator

Record the exact Mesen build. Use save states only with hash/provenance. Freeze frames deliberately. Document breakpoints, trace filters, viewers, memory ranges, and Lua scripts. Never depend on undocumented UI behavior without recording it locally.

Before any battle experiment, capture the in-game battle-configuration fingerprint alongside ROM and emulator identity. Prefer reading the settings from memory over reading the Config screen; when only a screen reading is available, record it as `screen-read` so the record shows its confidence. The hand cursor on that screen marks the row, not the selection.

Read and follow:

- `../_shared/EVIDENCE_STANDARD.md`
- `../_shared/ADDRESS_SPACES.md`
- `../_shared/STOPPING_RULES.md`
