# SESSION_006 — AUDIT-0001, workflow governance and orchestration utilization

**Date:** 2026-08-02
**Branch:** `maintenance/workflow-observability` (from `581ddbc`)
**Type:** Audit and design. **No game research. No emulator. No Ghidra.**

## What this session did

Audited the project's operating system — commands, skills, shared contracts,
agents, playbooks, rules, templates, their supporting implementations, their
documentation, and their historical use — and designed, without building,
telemetry that would make future usage measurable.

Nothing was remediated. Every finding is a proposal.

## What was not done

- No experiment was designed or started.
- No Mesen session was opened; no probe was run, including `dma-trace.lua`.
- No Ghidra tooling was run.
- No ROM was read; `FF6_ROM` was never set.
- No command, skill, agent, playbook, rule, template, schema, manifest, Go
  source file, Mesen probe or research record was modified.
- `/audit-project` was **not invoked** — see below.
- DEMO-0001 was not resumed.
- Nothing was pushed.

## Method

Ten phases: resume and gate baseline; read-only project inspection; inventory
and reference graphs; capability-honesty classification; thirteen live agent
smoke tests; historical usage reconstruction; scorecard; workflow outcome
analysis; telemetry design; remediation design and closure.

The reference graph matched three spellings per resource — full path, bare
name, and type-specific forms — across every `.md`, `.go`, `.lua`, `.json` and
`.sh` file in the tracked tree. A single-spelling pass was insufficient and was
corrected: an early path-only search reported 8 unreferenced templates; the
generated graph resolves it to 9.

## Results

### Inventory — 135 resources

43 commands, 30 skills, 5 shared contracts, 13 agents, 28 playbooks, 4 rules,
12 templates. Zero broken references. 31 orphans.

A planning-stage total of 126 was an arithmetic error, caught by the operator
before execution and withdrawn.

### Capability honesty — 43 commands

16 Implemented and Verified · 10 Partially Implemented · 10 Orchestration Only
· 4 Implemented but Unexercised · 3 Missing Backend · 0 Broken · 0
Documentation Only.

`/trace-dma` was split into four surfaces with four different statuses, as
required: offline decoding real and used; trace parser real but never
cross-validated against the probe; live tracing written and never run;
source-provenance refuted as a shortcut by EXP-0050's 18% match.

### Historical usage — the honest floor

6 Confirmed · 16 Probable · 12 Possible · **101 Unknown**.

Raw textual mention counts (`/checkpoint` 25, `/orchestrate` 10, …) are
recorded as mentions and never as invocations. A targeted search for
outcome-bearing attributions returns three commands; most `/command` mentions
sit in "Exact next action" sections, which are plans.

### Agent smoke tests — 13 of 13

All invoked successfully. **Zero scope violations; `git status` clean after
every batch**, including for the four agents holding Write/Edit/Bash authority.

Output-contract compliance: 1 full, 11 partial, 1 violated. Three agents made
errors the auditor caught and corrected. Two agents corrected the *auditor*:
`release-engineer` found CI enforcement in `.github/workflows/ci.yml` that a
`.claude`-only reading would have reported as absent, and `dma-researcher`
sharpened the `/trace-dma` analysis.

## Findings worth carrying forward

1. **Explicit agent routing is absent** while agent reachability is confirmed.
   Twelve of thirteen agents are orphaned; `ORCHESTRATE_RESEARCH.md` step 6 is
   the word "Delegate." Three independent agents converged on this.
2. **`/audit-project` has no report-only mode.** Its playbook's Required
   outputs mandate mutation, so the audit could not use the project's own audit
   command and performed its steps by hand.
3. **`internal/validate` has zero importers.** The validator written to catch
   the palette-index defect is wired to nothing, while its playbook prescribes
   the failure mode it prevents.
4. **Unit IDs collide globally** — Unit 12, 17 and 18 each name two unrelated
   units, which would corrupt any telemetry keyed on `unit_id`.
5. **The experiment schema omits three constitutional fields.**
6. **Real capability is undiscoverable** — `ff6lab state origin`/`sprites` are
   implemented and cited but absent from `ff6lab help`.
7. **The six orphaned recovery commands are the best-written in the set.**
   Orphan status is a routing fact, not a quality judgement.

## Records produced

Seven tracked records under `docs/workflows/`, listed in the checkpoint. Raw
evidence and generator scripts are ignored under
`local_artifacts/workflow-audit/AUDIT-0001/`.

## Quality gates

Run twice — pre-audit baseline and closure — with **identical results**.

gofmt clean; build, vet, test pass; `ff6lab`, `ff6demo` and `-tags gui` all
compile; `ff6lab audit` clean; `census validate` clean; restricted-file scan
clean.

**`archive verify` NOT RUN** (`FF6_ROM` unset).

## Contradiction recorded

`CLAUDE.md` requires dashboard synchronization at unit closure; the audit brief
forbade dashboard modification. Resolved for this session as
generated/required synchronization only — `LATEST.md`, `indexes/SESSIONS.md`,
and the session-count row in `STATISTICS.md`. The contradiction is preserved in
`AUDIT-0001-orchestration-inventory.md` rather than silently resolved.

## Next action

Review and approve, revise, or reject
`docs/workflows/AUDIT-0001-remediation-plan.md`.
