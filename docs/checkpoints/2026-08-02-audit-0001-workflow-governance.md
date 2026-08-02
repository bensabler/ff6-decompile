# Checkpoint — AUDIT-0001, workflow governance and orchestration utilization

**Date:** 2026-08-02 · **Branch:** `maintenance/workflow-observability`
**Parent:** `581ddbc` on `demo/whelk-content-parity`

## What this session was

**A workflow audit. No game research occurred.** No Mesen session was opened at
any point. No Ghidra tooling was run. No experiment was designed or started. No
ROM was read. No remediation was implemented.

**Every finding in the AUDIT-0001 records is a proposal awaiting operator
review.**

## Current question

None. AUDIT-0001 is complete and its remediation plan is unapproved. No
research or implementation unit is in flight.

## State

Branch `maintenance/workflow-observability`, created from
`demo/whelk-content-parity` at `581ddbc`, which remains synchronized with
`origin`. The audit branch is local-only and **was not pushed**.

No emulator running. No background processes. No resident instrumentation.

## Work completed

Ten phases, all complete. `AUDIT-0001-baseline.json` records
`audit_status: complete`, `completed_phases: [0..10]`, `remaining_phases: []`.
The Phase 4 context-safety boundary was assessed and **did not fire** — the
audit finished in one session, so one session record and one commit, as the
plan's single-session case provides.

### Files created — seven tracked records

```text
docs/workflows/AUDIT-0001-orchestration-inventory.md
docs/workflows/AUDIT-0001-capability-honesty.md
docs/workflows/AUDIT-0001-usage-baseline.md
docs/workflows/AUDIT-0001-effectiveness-scorecard.md
docs/workflows/AUDIT-0001-telemetry-design.md
docs/workflows/AUDIT-0001-remediation-plan.md
docs/workflows/AUDIT-0001-baseline.json
```

Plus this checkpoint, `docs/sessions/SESSION_006.md`, and the permitted closure
synchronization: `docs/checkpoints/LATEST.md`, `indexes/SESSIONS.md`, and the
single session-count row in `dashboards/STATISTICS.md`.

Raw evidence and the generator scripts are ignored under
`local_artifacts/workflow-audit/AUDIT-0001/`.

### Headline results

**135 resources** inventoried — 43 commands, 30 skills, 5 shared contracts,
13 agents, 28 playbooks, 4 rules, 12 templates. The figure of 126 used during
planning was an arithmetic error and is withdrawn.

**Zero broken references.** **31 orphans**, of which 12 are agents.

**Agent routing** — stated as two separate facts. Reachability is **Confirmed**:
all 13 agents were invoked successfully. Explicit project routing is
**Confirmed absent**: no command, skill, playbook, rule or template names any
agent, and `ORCHESTRATE_RESEARCH.md` step 6 is the single word "Delegate."
Three independent agents converged on this, including `quality-reviewer` given
free choice of target.

**Capability honesty** — 16 Implemented and Verified, 10 Partially Implemented,
10 Orchestration Only, 4 Implemented but Unexercised, 3 Missing Backend.
Nothing Broken, nothing Documentation Only.

**Historical usage** — 6 Confirmed, 16 Probable, 12 Possible, **101 Unknown**.
Unknown means no evidence was found; it is never reported as zero use.

**Zero removal, deprecation or merge candidates.**

### Findings the brief did not anticipate

- **`/audit-project` has no report-only mode.** Its playbook's Required outputs
  mandate updated records and a checkpoint, so the fix-nothing override could
  not be guaranteed. **The command was not invoked**; its eight steps were
  performed manually, read-only, and the missing mode is recorded as a
  capability-honesty finding.
- **`internal/validate` has zero importers.** `CompareResolved` — written to
  catch the palette-index defect — is reachable only from its own tests, while
  `VALIDATE_GRAPHICS.md:16` prescribes indexed comparison, the exact failure
  mode it prevents.
- **Unit IDs collide globally.** `ACTIVITY_LOG.md` reached Unit 35 by 08-01;
  Session 005 restarted at Unit 10. Unit 12, 17 and 18 each name two unrelated
  units. The proposed telemetry `unit_id` must be session-qualified.
- **`schemas/experiment.schema.json` omits three constitutional fields** —
  controlled variables, required evidence, stopping condition — and `domain`
  appears in 12 manifest entries the schema never defines.
- **`/session-summary` is missing from `docs/WORKFLOW_COMMANDS.md`**, and does
  not name `SESSION.md`, which is why that template is orphaned.
- **`ff6lab state origin` and `state sprites` are absent from `ff6lab help`**
  despite being cited as durable instrumentation.
- **`ASSET-GFX-0002`** (Narshe tileset) links **`CEN-GFX-0006`** (mines
  interior). Not repaired — research records are outside this audit's scope.

## Uncertain

- 101 of 135 resources have no recoverable usage evidence. This is a
  measurement failure, not a usage failure, and it is the floor of what
  historical reconstruction can establish.
- Skill and playbook composite scores do not discriminate (28 of 30 skills at
  exactly 4.00; all 28 playbooks at 4.67). The dimensions that would separate
  them need invocation telemetry. **Those composites must not drive decisions**
  — stated in the scorecard.
- Whether `/capture-graphics` genuinely overlaps `/capture-frame` or has a real
  distinct scope is a scoping call for the operator.

## Contradiction preserved, not resolved

`CLAUDE.md` requires a completed unit to synchronize dashboards; the AUDIT-0001
brief forbids modifying them. Resolved for this session as generated/required
synchronization only. **The contradiction is recorded in the inventory record**
and a future governance decision should state which authority wins for
audit-class sessions.

## Active emulator state

None. No emulator was launched.

## Breakpoints/watchers

None.

## Evidence paths

`local_artifacts/workflow-audit/AUDIT-0001/` — `inventory.json`,
`reference-graph.json`, `command-capabilities.json`, `agent-smoke-tests.json`,
`historical-usage.json`, `broken-references.json`, `gates-baseline.txt`,
`gates-closure.txt`, and the four generator scripts. All gitignored; none
ROM-derived.

## Tests and quality gates

Run **twice** — a pre-audit baseline at Phase 0 and again at closure. **Results
identical.**

gofmt clean; `go build`, `go vet`, `go test` pass; `ff6lab`, `ff6demo` and the
`-tags gui` build all compile; `ff6lab audit` clean (eleven checks);
`census validate` clean; restricted-file scan clean; `AUDIT-0001-baseline.json`
parses and reports 135 resources.

**`archive verify` NOT RUN** — `FF6_ROM` is unset. Recorded as not run, never
as passing.

## Git status

Branch `maintenance/workflow-observability`. One audit-only commit. Not pushed,
as directed.

## Unresolved decisions

Whether to approve, revise or reject the AUDIT-0001 remediation plan. Nothing
in it has been implemented.

## Blockers

None.

## Exact next action

**Review and approve, revise, or reject AUDIT-0001's remediation plan**
(`docs/workflows/AUDIT-0001-remediation-plan.md`).

If approved, the recommended order is P0-1 agent routing → P0-3 documentation
sync with its gate → P0-2 `/audit-project` report-only mode → P1-1 wire
`internal/validate`. Items 1–4 are text and small Go changes with no migration
cost and address the three highest-risk findings. Nothing in the plan requires
an emulator except the `/trace-dma` live validation, which is deferred to an
authorized operator session.

DEMO-0001 remains paused and unchanged. When content work resumes, the
tactical-pause checkpoint's priority list still stands.
