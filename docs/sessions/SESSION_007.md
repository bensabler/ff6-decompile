# SESSION_007 — AUDIT-0002, independent re-audit of AUDIT-0001

**Date:** 2026-08-02
**Branch:** `maintenance/workflow-observability-audit2`, worktree `../ff6-audit2`
**Base:** `93f7d03` (AUDIT-0001 closure)
**Type:** Audit and design. **No game research. No emulator. No Ghidra.**

## What this session did

Re-audited AUDIT-0001 from an isolated worktree, corrected its record without
rewriting it, completed a capability baseline across eight evidence-bearing
axes for all 135 workflow resources, and produced a corrected remediation plan
aimed at the operator-level problem: with 135 resources, the operator should
not have to read raw tool activity to learn whether the right ones were used.

Nothing was remediated.

## What was not done

- No experiment designed or started; no Mesen session; no Ghidra tooling.
- No ROM read; `FF6_ROM` never set.
- No command, skill, agent, playbook, rule, template, schema, manifest, Go
  file or Mesen probe modified.
- **AUDIT-0001's six records and baseline unmodified** — corrections are in
  `AUDIT-0001-errata.md`.
- `/orchestrate`, `/audit-project`, `/checkpoint`, `/session-summary` not
  invoked. No `docs/workflows/runs/` created. No workflow command built.
- Nothing pushed.

## Method

Twelve phases. A slash-command feasibility preflight ran **first**, because the
whole correction depended on being able to prove invocation rather than recall
it — and AUDIT-0001's central defect was crediting a command for work done by
hand.

The preflight succeeded: the harness writes per-session transcripts recording
every tool call with a timestamp. Only **tool-call metadata** was extracted —
tool name, timestamp, selector. No prompts or conversation content.

Generators were replayed against a detached worktree at the **true historical
input tree**, never the evolving audit branch.

## Results

**Invocation, corrected in both directions.** Five slash-command invocations
exist in all preserved history. `/run-quality-gates` is in none of them — its
`Confirmed` is refuted. `/correlate-static-runtime` is unverifiable.
`/bootstrap-ghidra` *was* invoked, so AUDIT-0001 reached a true conclusion via
the invalid artifact-existence method. `/checkpoint` and `/session-summary` were
each invoked and had been under-classified. Coverage is 2026-08-02 only.

**Orphans: a measurement bug.** The generator replays byte-identically; its
corpus never scanned repository-root files, and `PACKAGE_MANIFEST.json` lists 24
of the 31. Corrected: 7 textual orphans. Measured on routing-bearing inbound —
the metric that matters — **38**, with **all 13 agents at zero**.

**AUDIT-0001 destroyed its own metric.** Replayed against the closure tree the
matcher returns zero orphans, because its baseline names all 135 resources.

**Automatic specialist selection does not exist.** Pre-registered probes:
omitting `subagent_type` yields `general-purpose`.

**Baseline completed.** 135 resources × 8 axes, validation PASS, 0 problems.
662 values mechanically derived, 418 manually adjudicated with rationales. All
composite scores withdrawn. All 28 playbooks individually assessed on content:
14 Keep, 5 Keep and Improve, 4 Repair, 3 Narrow Scope, 1 Merge Candidate,
1 Deferred.

**Remediation plan v2:** 14 items with nine negative acceptance tests, ordered
so the enforcement substrate and deterministic verdict precede the operator
surface they support.

## Phase 11 review

Findings were frozen and hashed before any review.
`verification-engineer` ran blind to the conclusions; `quality-reviewer` was
asked to attack them. Six disputes, all resolved, none escalated.

**Four corrections came from review, not from me** — gate capture counts, my
own hash file's working directory, the circularity of `operator_workflow == 0`,
and an artifact bar applied to resources that produce no artifacts.

The sharpest criticism is disclosed rather than removed: **R11-R14's diagnosis
and remedy share one rubric this audit defines.**

## Quality gates

Run twice with reliable exit capture, **identical both times**: eleven required
gates `rc=0`, AGGREGATE PASS. The harness was negative-control tested rather
than assumed to work — and it caught a real failure mid-audit, when frozen
evidence broke `ff6lab audit`.

**`archive verify` NOT RUN** — `FF6_ROM` unset.

## Standing limitations

Reviews are **correlated, not independent** — same underlying model as the
auditor. I re-audited my own work; that is mitigated, not eliminated.
Transcript coverage does not reach sessions 001-004. AUDIT-0001's Phase 4
process metrics are unrecoverable.

## Next action

Review and approve, revise, or reject
`docs/workflows/AUDIT-0002-remediation-plan-v2.md`.
