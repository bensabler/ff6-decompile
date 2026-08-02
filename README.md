# FF6 Ultimate Reconstruction Project

An evidence-driven, production-oriented framework for reconstructing **Final Fantasy III (USA) / Final Fantasy VI for the SNES** as an idiomatic Go project.

The repository is designed to preserve four distinct products:

1. **Verified research** — addresses, routines, structures, formats, experiments, and confidence.
2. **Original Go implementation** — maintainable code that recreates verified behavior and formats.
3. **Reconstruction tooling** — ROM readers, tracers, decoders, compositors, validators, and catalog tools.
4. **Local asset workspace** — legally obtained ROMs, emulator states, runtime captures, and extracted assets that remain outside Git.

The repository—not a chat session—is the source of truth.

## What this project is not

- It is not a ROM-distribution project.
- It is not a one-click “decompile the entire game” promise.
- It does not treat an emulator screenshot, PNG export, WAV capture, or disassembly guess as a proven source representation.
- It does not permit speculative names or structures to silently become facts.

## Core pipeline

```text
Controlled Mesen experiment
        ↓
Raw runtime evidence
        ↓
Address-space and provenance analysis
        ↓
ROM format / routine reconstruction
        ↓
Idiomatic Go implementation
        ↓
Automated validation against evidence
        ↓
Canonical documentation and checkpoint
```

## Running the demo

DEMO-0001 is the playable deliverable: a Go reconstruction of New Game through
Whelk victory. It is **early** — see
[`docs/demo/DEMO-0001-READINESS.md`](docs/demo/DEMO-0001-READINESS.md) for the
honest per-requirement state.

Generate the assets from your own verified ROM (nothing ROM-derived ships in
this repository):

```bash
export FF6_ROM=/path/to/your/verified.sfc
go run ./cmd/ff6lab extract all
```

Then run it:

```bash
go run -tags gui ./cmd/ff6demo
```

Or headless, which needs no display and is the authoritative mode for
automated comparison:

```bash
go run ./cmd/ff6demo -headless -frames 120 -capture-last
```

The demo binary cannot read a ROM: `internal/rom` is barred, transitively, from
`cmd/ff6demo`, `internal/content` and `internal/engine`, and `ff6lab audit`
walks the import graph to enforce it.

### Dependencies

The default build has **no third-party dependencies**, and every quality gate
runs without system libraries. The windowed host is the single exception:
Ebitengine, confined to `internal/engine/ebitenhost` and gated behind the `gui`
build tag, because it requires cgo on macOS and Linux. Rationale and
alternatives: [ADR-0001](docs/decisions/ADR-0001-rendering-host.md).

## Quick start (research workflow)

1. Merge this package into the root of the FF6 reconstruction repository.
2. Keep the ROM and all extracted commercial material under `local_artifacts/`.
3. Open Claude Code from the repository root.
4. In a fresh session, run:

```text
/bootstrap-v4
/resume-session
/audit-project
/orchestrate
```

5. Before ending or interrupting the session:

```text
/checkpoint
```

Read these next:

- [`OPERATIONS.md`](OPERATIONS.md)
- [`.claude/README.md`](.claude/README.md)
- [`docs/WORKFLOW_COMMANDS.md`](docs/WORKFLOW_COMMANDS.md)
- [`docs/FRESH_SESSION_SEQUENCE.md`](docs/FRESH_SESSION_SEQUENCE.md)
- [`ARCHITECTURE.md`](ARCHITECTURE.md)
- [`docs/research/EVIDENCE_STANDARD.md`](docs/research/EVIDENCE_STANDARD.md)

## Established baseline

The following observations are preserved as initial evidence, not as a completed conclusion:

- Current HP was observed at `WRAM+$2EB5` for party slot 0 in one battle state.
- A write breakpoint reached ROM CPU `$C10E14`: `STA $2EB5,Y`.
- The surrounding routine was approximately `$C10DF3`–`$C10E66`.
- The observed loop advanced X by 2 and Y by `$20`.
- This supports, but does not yet prove, a 32-byte destination record.

See `docs/research/baseline/C10DF3_HP_COPY.md`.
