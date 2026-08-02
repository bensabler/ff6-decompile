# AUDIT-0001 — Telemetry design

**Audit date:** 2026-08-02 · **Status: design only. Nothing here is built.**

## The problem this solves

Historical reconstruction produced **101 `Unknown`** classifications out of 135
resources. That is not a usage failure — it is a measurement failure. The
project has never recorded which resource did which work, so the value of three
quarters of its orchestration layer cannot be assessed.

The design goal is narrow: make invocation recoverable **without depending on
Claude remembering to report it**, and without storing prompts, conversation
content, ROM bytes or personal paths.

## Constraint discovered during the audit — unit IDs collide

`dashboards/ACTIVITY_LOG.md` numbers units globally and reached Unit 35 by
2026-08-01 (Unit 49 in the EXP-0047 entry). Session 005 restarted at Unit 10 on
2026-08-02 and ran to Unit 18.

**Unit 12, 17 and 18 each name two unrelated pieces of work.**

The brief's proposed event shape carries a bare `unit_id`. Adopting it unchanged
would silently merge those units and corrupt the first metric the telemetry is
meant to produce. **`unit_id` must be qualified by `session_id`**, and the
registry must treat the pair as the key.

This is the kind of defect the telemetry exists to catch, found before the
telemetry was built.

## Event shape

Written to the ignored `local_artifacts/workflow-telemetry/` as newline-delimited
JSON, one file per session.

```json
{
  "timestamp": "2026-08-02T15:05:10Z",
  "session_id": "SESSION-006",
  "unit_id": "SESSION-006/UNIT-01",
  "event": "start",
  "resource_type": "command",
  "resource_id": "/trace-dma",
  "invoked_by": "user",
  "task_class": "graphics-provenance",
  "status": null,
  "outputs": [],
  "commit": null,
  "reason_code": null,
  "duration_seconds": null
}
```

Changes from the brief's proposal, each with a reason:

| Field | Change | Reason |
|---|---|---|
| `unit_id` | **Session-qualified** (`SESSION-006/UNIT-01`) | Unit IDs collide globally; see above |
| `invoked_by` | Add `"orchestrator"` as a distinct value from a resource id | Agent invocations this session came from the orchestrator, not from any resource — the routing gap must be visible in the data |
| `resource_type` | Add `"tool"` | `ff6lab` subcommands are the real backends; without this, `ff6lab state origin` usage stays invisible |
| `alias_invoked_as` | **New, optional** | The brief requires aliases normalize to their canonical command *while preserving the invoked spelling*. `resource_id` carries the canonical id; this carries what was typed. |
| `superseded_by` | **New, optional** | Lets a `bypass` event name what was done instead |

Privacy: no field carries free text from the conversation. `task_class` and
`reason_code` are **closed vocabularies** validated by the registry schema, so
neither can leak content.

## Registry

`manifests/workflow-registry.json`, validated by
`schemas/workflow-registry.schema.json`.

The registry is the declared inventory — one entry per resource, mirroring the
shape already proven in `AUDIT-0001-baseline.json`: id, path, type, declared
inputs/outputs, required skills and playbooks, expected artifact classes,
implementation status, routing status.

`AUDIT-0001-baseline.json` is the schema's **first real input**. The registry
must be able to load it and diff against it, so that "what changed since the
audit" is a command rather than a re-audit.

The schema must avoid the defect found in `schemas/experiment.schema.json`,
which defines seven properties while `CLAUDE.md` requires seven *different*
fields — three constitutional fields have no schema property at all, and
`domain` appears in 12 manifest entries the schema never defines. The workflow
registry schema should set `additionalProperties: false` so that drift fails
loudly instead of accumulating silently.

## `internal/workflow` and the `ff6lab workflow` surface

Model the package on `internal/audit`, which is the convention already proven
here: `audit.Run(root)` composes independent check functions against an
explicitly passed repository root, and every check is tested against fixture
trees. `CheckManifests`, `CheckExperimentIndexSync` and
`GenerateExperimentIndex` are the direct precedents.

```text
ff6lab workflow inventory   declared resources; diff against the registry
ff6lab workflow validate    registry <-> filesystem <-> reference graph
ff6lab workflow start       append a start event
ff6lab workflow finish      append a finish event with status and outputs
ff6lab workflow bypass      record that a resource was skipped, and why
ff6lab workflow report      aggregate events into the metrics below
```

`inventory` and `validate` are the load-bearing pair: they subsume most of what
AUDIT-0001 did by hand — orphan detection, broken-reference detection, doc
coverage, declared-versus-actual comparison — and they need **no telemetry at
all** to be useful. They should be built first and wired into `ff6lab audit`.

### Metrics `report` must support

Invocation counts · first and latest confirmed use · success/partial/failed/
abandoned counts · output artifacts · downstream commits · agent routing ·
bypass reasons · orphan detection · broken references · implementation status ·
per-session orchestration coverage · research-to-implementation conversion ·
implementation-to-demo conversion · declared-versus-actual comparison.

**Stale-resource warnings fire only after a configurable number of instrumented
units** (proposed default: 20). Before that threshold the correct report is
`insufficient data`, never `unused`. This is the schema-level enforcement of
the rule that absence of evidence is `Unknown`, not zero — the same rule
AUDIT-0001 applied by hand to 101 resources.

## Not depending on Claude to remember

Self-reporting is the weakest link, and this audit has direct evidence: of 13
agents given an explicit, unambiguous output contract in their own prompts,
**one complied fully, eleven partially, one violated it.** A design that
assumes reliable voluntary reporting is contradicted by measurement taken this
session.

Three mechanisms, in order of reliability:

**1. Derived events — no cooperation required (strongest).**
`ff6lab workflow validate` can reconstruct much of the signal from artifacts
that already exist: a new file under `docs/checkpoints/` implies `/checkpoint`;
a new `docs/correlations/CORR-*.md` implies `/correlate-static-runtime`; a new
`local_artifacts/static-analysis/exports/` artifact implies
`/export-ghidra-symbols`. This is exactly the inference AUDIT-0001 used to
reach `Probable`, promoted from manual reading to a tested function. It cannot
be forgotten because nothing has to remember it.

**2. A `SessionStart`/`Stop` hook pair.** Hooks are executed by the harness, not
by Claude, so they cannot be skipped. A `SessionStart` hook opens the session's
event file and stamps the session id; a `Stop` hook closes it and runs
`ff6lab workflow report --session`. This is the only mechanism that reliably
captures *session-level* orchestration coverage.

**3. A shared command lifecycle contract.** A new
`.claude/skills/_shared/WORKFLOW_LIFECYCLE.md`, referenced by every command the
way `EVIDENCE_STANDARD.md` and `STOPPING_RULES.md` already are, requiring a
`workflow start` before and a `workflow finish` after. This is the weakest link
and should be treated as a supplement to (1) and (2), never the primary
mechanism.

The ordering matters. **Build (1) first.** It requires no behaviour change from
anyone, it validates the registry against reality, and it would have caught the
`SESSION_004` index gap, the six orphaned commands, and the unimported
`internal/validate` — all three of which were found by hand in this audit and
none of which needed telemetry.
