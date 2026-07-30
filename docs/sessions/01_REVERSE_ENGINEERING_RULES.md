# Reverse Engineering Rules

These rules govern every discovery in this repository.

> **Superseded in part (2026-07-29,
> [CONTRA-0001](../contradictions/CONTRA-0001-address-notation.md)):** the
> authoritative standards now live in `docs/research/` (evidence, confidence,
> address notation) and `.claude/skills/_shared/`. The
> `ff6-reconstruction-skill` this file once referenced is retired. This file
> remains the historical pre-V4 rule set; its evidence classifications match
> the V4 standard, but its address-notation rule is superseded — see below.

## Evidence classifications

- **Confirmed** — directly demonstrated by repeatable debugger evidence or
  byte-for-byte behavior.
- **Strong hypothesis** — multiple independent observations support it, but
  an alternative remains possible.
- **Tentative hypothesis** — plausible interpretation based on limited
  evidence.
- **Unknown** — insufficient evidence.

Never present a hypothesis as confirmed.

## Naming

Names must match the evidence. Prefer `PartySlotCurrentHP`,
`UnknownCharacterField02`, `CopyCharacterFields`. Avoid ownership or purpose
claims (`TerraHP`, `ApplyDamage`, `InitializeBattle`) until demonstrated.
Uncertain names carry a marker such as `Unknown` or `Candidate`.

When a name improves: rename consistently in code and docs, record the old
name under **Previous names** in the discovery record, and state the evidence
that justified the rename.

## Address notation

> **Superseded (2026-07-29):** new and substantively updated canonical docs
> must use domain-prefixed notation per
> [ADDRESS_NOTATION.md](../research/ADDRESS_NOTATION.md) (`ROMCPU:$C10E14`,
> `WRAM:+$2EB5`). Pre-V4 documents keep the contextual notation below as
> historical records. See
> [CONTRA-0001](../contradictions/CONTRA-0001-address-notation.md).

Always distinguish, and never silently interchange:

- SNES CPU address — `$7E2EB5`
- WRAM-relative offset — `$2EB5`
- ROM CPU address — `$C10DF3`
- ROM file offset — when known

Hexadecimal uses `$` in reverse-engineering docs and `0x` in Go source.

## Workflow for each discovery

1. Read the relevant project documentation.
2. Inspect the debugger evidence.
3. Separate direct observations from interpretations.
4. Assign a confidence level to every interpretation.
5. Identify missing evidence and contradictions.
6. Propose the smallest experiment that could confirm or disprove it.
7. Update documentation before or alongside implementation.
8. Implement only behavior supported by sufficient evidence.
9. Add tests that encode the verified behavior.
10. Run `gofmt -w .`, `go test ./...`, `go vet ./...`.

## Document ownership

- Memory addresses → [04_MEMORY_MAP.md](04_MEMORY_MAP.md)
- Function records → [02_DISCOVERED_FUNCTIONS.md](02_DISCOVERED_FUNCTIONS.md)
- Struct layouts → [05_DATA_STRUCTURES.md](05_DATA_STRUCTURES.md)
- Chronological work → [SESSION_*.md files in this directory](.)
- Unresolved items → [08_OPEN_QUESTIONS.md](08_OPEN_QUESTIONS.md)

Facts live in exactly one owning file; other files cross-reference.
