# AUDIT-0001 — Capability honesty

**Audit date:** 2026-08-02 · Companion to `AUDIT-0001-orchestration-inventory.md`

Every classification below was made from **backing code, probes and scripts** —
never from Markdown alone. Per-command evidence is in
`AUDIT-0001-baseline.json`.

## Vocabulary

`Implemented and Verified` · `Implemented but Unexercised` ·
`Partially Implemented` · `Orchestration Only` · `Documentation Only` ·
`Missing Backend` · `Broken` · `Unverified`

**`Orchestration Only` is not a defect.** A command that is a disciplined
procedure over Claude's own reading and writing of records is a legitimate
design. It becomes a defect only when the command's text claims tooling that
does not exist.

## Distribution — all 43 commands

| Status | Count |
|---|---|
| Implemented and Verified | 16 |
| Partially Implemented | 10 |
| Orchestration Only | 10 |
| Implemented but Unexercised | 4 |
| Missing Backend | 3 |
| Documentation Only | 0 |
| Broken | 0 |
| Unverified | 0 |

Nothing is Broken. Nothing is Documentation Only. The orchestration layer is in
better condition than the absence of telemetry suggests — but three commands
prescribe procedures against tooling that does not exist.

## `/trace-dma` — four surfaces, four statuses

The brief required these be kept distinct. Reporting one status would be
dishonest.

| Surface | Status | Evidence |
|---|---|---|
| Offline DMA register-state decoding | **Real, tested, used** | `internal/platform/snesdma/snesdma.go`; EXP-0050 used it |
| Trace parser | **Real, tested — but never cross-validated** | `internal/platform/snesdma/trace.go` + tests. Fixtures are hand-authored synthetic lines shaped to match the probe's comment spec. Parser and probe have never seen each other's output. |
| Live DMA-enable tracing | **Written, never run** | `mesen/probes/dma-trace.lua` — its own header says `STATUS: UNEXERCISED` |
| Source-provenance validation | **Refuted as a shortcut** | EXP-0050: the mines savestate's channel 0 claims `$7EC180 → $2118`, but VRAM matches that buffer on only **18%** of bytes |

**Configuration at rest is not a trace.** Savestate channel state must never be
reported as a live DMA trace.

`/trace-dma` therefore scores **5 on capability honesty while only 3 on backend
completeness.** Honesty and completeness are different axes. The probe header,
the playbook and the command all state their limits plainly — this is the
project documenting its own gap correctly, and it is the model the rest of the
orchestration layer should follow.

The one closing step: run `probe dma-trace` over a map load and feed one
resulting line through `ParseTraceLine`. That single run validates surfaces 2,
3 and 4 together. **Not performed — this audit operates no emulator.**

## Finding — `/audit-project` has no report-only mode

The command says *"Fix only safe mechanical defects; report substantive changes
for approval."* `AUDIT_PROJECT.md`'s **Required outputs** are: updated canonical
record, raw evidence references and hashes, confidence and alternatives,
relationship/index updates, exact next action, checkpoint when the unit stops.

Every documented execution path mutates the repository. There is no
report-only mode to fall back on.

**Consequence for this audit:** the brief's fix-nothing override could not be
guaranteed from within the command's own contract, so AUDIT-0001 **did not
invoke it**. Its eight procedure steps were performed manually, read-only.

A command whose only execution path mutates the repository cannot be used to
assess the repository. That is a governance defect in its own right, and it is
the reason `/audit-project` scores **1 on capability honesty** despite having
eleven genuine checks behind `ff6lab audit`.

## Finding — three commands prescribe absent tooling

| Command | Claim | Reality |
|---|---|---|
| `/validate-audio` | DSP/voice/waveform comparison | `internal/validate` contains **framebuffer comparison only**. No DSP, APU or audio comparison code exists anywhere in `internal/`. |
| `/recover-sequence` | sequence format recovery | No sequence package exists. |
| `/trace-spc-command` | SPC700 port/driver tracing | No SPC700 tooling beyond `mesenstate`'s raw ARAM read. |

The `dsp-validator` and `sequence-reconstructor` skills exist and are
command-reachable. The chain runs command → skill → playbook and then stops.

## Finding — a Go doc comment overstates its package

`internal/audio/doc.go` reads:

> Package audio contains SPC700, DSP, sequence, and sample reconstruction code.

The package contains one subpackage: `brr`. This is capability dishonesty
inside Go source, a class a Markdown-only audit would not have found. Located
by the `audio-researcher` smoke test and verified directly.

## Finding — real capability that is undiscoverable

`ff6lab state origin` and `ff6lab state sprites` are implemented
(`cmd/ff6lab/state.go`), and the tactical-pause checkpoint cites them as
durable tracked instrumentation. **Neither appears in `ff6lab help`.**

`/recover-compression` instructs the operator to "Run `ff6lab state origin`
before assuming a format exists" — a command the tool's own help does not
mention.

This is the inverse of the usual defect: not a documented capability that does
not exist, but an existing capability that is not documented. Both break the
same contract.

`ff6lab help` also advertises "Planned command groups: evidence, validate,
report" — three groups with no implementation.

## Finding — the validator nothing calls

The sharpest result of the audit.

`internal/validate/framebuffer.go` provides `CompareResolved`, `Visibility` and
`PaletteUse`. They exist specifically because hashing palette *indices* is not
the same identity as comparing visible *colour* — the defect that rendered the
demo's text black on black (fixed in `6d36e42`). The file documents the lesson.
`TestResolvedCatchesInvisibleInk` proves it works.

**`go list` shows `internal/validate` has zero importers — no non-test
importers, and no test importers outside its own package.**

Meanwhile `VALIDATE_GRAPHICS.md` step 4 reads:

> Compare indexed pixels when possible.

The playbook prescribes the exact failure mode the code was written to prevent,
and the code that prevents it is wired to nothing.

This is scored as a **genuine zero** on downstream integration — not unknown.
The evidence exists and the measured value is zero.

## Finding — release gates are stronger than the playbooks admit

An inspection limited to `.claude/` would have reported the demo binary,
restricted files, generated assets and the GUI build as ungated. That would
have been **wrong**.

`.github/workflows/ci.yml` has a `test` job that builds both binaries, runs
`ff6lab audit`, and runs two inline `git ls-files | grep` steps rejecting
tracked ROM/audio/state files and rendered graphics; plus a separate
`gui-build` job that installs windowing headers, verifies pinned module
checksums, and runs `go build -tags gui` and `go vet -tags gui`.

All four are gated. But `RUN_QUALITY_GATES.md` names none of them — it has no
build step at all, no `ff6demo`, no `gui` tag, no asset scanning. An operator
following only the playbook would not know these gates exist.

Residual gap: `restrictedExt` and `renderedExt` are **closed extension lists**.
`CheckTrackedBinaries` exists because an extensionless `ff6demo` binary reached
the remote once already, which is precisely the blind spot a closed list
leaves.

Found by the `release-engineer` smoke test, which corrected the auditor.

## Finding — the discovery-to-demo link is prose only

`IMPLEMENT_DISCOVERY.md`'s required outputs stop at "relationship/index
updates" and "checkpoint." Nothing requires updating
`docs/demo/DEMO-0001-READINESS.md`.

The link exists because a human copies the DISC-id into the readiness matrix's
"Known evidence" column. `CheckReadinessSummary` validates the matrix's
internal consistency — status counts, totals, vocabulary, duplicate ids — not
that any discovery is represented in it.

"Did this discovery reach `ff6demo`?" is answerable today by grep, and only by
grep discipline that nothing enforces.

## Commands whose backends are real and exercised

For balance — 16 commands are Implemented and Verified, and several are strong:

- **`/battle-baseline`** is the project's only complete
  command → rule → automated-enforcement chain. Battle experiments lacking a
  configuration fingerprint **fail `ff6lab audit`**
  (`CheckBattleExperimentConfig`). It exists because EXP-0040 misread `Bat.Mode`
  off a screenshot, and the gate makes that class of error impossible to repeat
  silently. This is the model for every remediation in this audit.
- `/run-quality-gates` — every gate it names is real; run twice this session.
- `/census-observations` and family — schema, manifest, validator, sync.
- `/reconstruct-tileset`, `/recover-brr`, `/recover-palette`,
  `/capture-frame` — real decoders with tests and shipped output.
