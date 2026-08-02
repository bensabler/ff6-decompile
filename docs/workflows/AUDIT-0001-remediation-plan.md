# AUDIT-0001 — Remediation plan

**Audit date:** 2026-08-02 · **Status: proposals only. Nothing was implemented.**

Every proposal below states its problem, the existing resource considered
first, why that resource is insufficient, the exact files affected, migration
impact, test strategy, telemetry behaviour, rollback and priority.

## Design stance

The project already contains the pattern every proposal here imitates:
**`/battle-baseline` → the fingerprint discipline → `CheckBattleExperimentConfig`
→ `ff6lab audit` fails.** A discipline became a gate, and that class of error
cannot silently recur.

Every remediation below tries to convert a discipline into a gate rather than
adding another document asking someone to remember. Where a proposal cannot do
that, it says so.

**Also: prefer fixing existing resources over adding new ones.** Of the nine
candidate new resources the brief listed, this plan recommends **two**, defers
five, and rejects two.

---

## P0-1 — Route the agents explicitly

**Problem.** 12 of 13 agents have zero inbound references. No command, skill,
playbook, rule or template names any agent. `ORCHESTRATE_RESEARCH.md` step 6 is
the single word "Delegate." Three independent smoke-test agents converged on
this as the top defect.

**Existing resource considered.** `ORCHESTRATE_RESEARCH.md` and
`research-orchestrator/SKILL.md` already exist and are the natural home.

**Why insufficient as-is.** They instruct delegation without naming a target.
`.claude/README.md` states the project's own rule — *"Each command names the
skills and playbook it requires"* — precisely because automatic matching is
unsafe for critical procedures. That discipline was applied to skills (0 of 30
orphaned) and playbooks (0 of 28 orphaned) and never extended to agents.

**Proposal.** Add a domain → agent routing table to
`.claude/playbooks/ORCHESTRATE_RESEARCH.md`, and name the relevant agent in each
domain command the way skills are already named. **No new resource.**

**Files affected.** `.claude/playbooks/ORCHESTRATE_RESEARCH.md`;
`.claude/skills/research-orchestrator/SKILL.md`; optionally one line in each of
~13 domain commands.

**Migration impact.** None. Additive text; no existing behaviour changes.

**Test strategy.** Extend `ff6lab audit` with an orphan check: every file under
`.claude/agents/` must be named by at least one command, skill or playbook.
This is the same shape as `CheckRequiredTracked` and would fail today.

**Telemetry.** `invoked_by` distinguishes `orchestrator` from a named resource,
so the routing gap stays measurable after the fix.

**Rollback.** Delete the table; the check is the only thing that would fail.

**Priority: P0.** Highest leverage, lowest risk, no new resource.

---

## P0-2 — Give `/audit-project` a report-only mode

**Problem.** Its only documented execution path mutates the repository.
AUDIT-0001 could not invoke the project's own audit command and performed its
eight steps by hand.

**Existing resource considered.** `ff6lab audit` is already read-only and has
eleven real checks.

**Why insufficient.** It covers manifests, links, boundaries and binaries — not
the substantive review steps (unsupported facts, address ambiguity, stale
dashboards, provenance gaps) that `AUDIT_PROJECT.md` specifies.

**Proposal.** Add an explicit report-only mode to the command and playbook —
`/audit-project report` — whose Required outputs are a findings list and
nothing else. Keep the existing fixing mode unchanged as the default.

**Files affected.** `.claude/commands/audit-project.md`;
`.claude/playbooks/AUDIT_PROJECT.md`.

**Migration impact.** None; default behaviour unchanged.

**Test strategy.** No automated test is possible for a Claude-executed
procedure. The honest check is the next audit-class session: it should be able
to invoke the command. Record that as the acceptance criterion.

**Telemetry.** `reason_code` distinguishes `report` from `fix` runs.

**Rollback.** Remove the mode.

**Priority: P0.** This audit was itself the blocked case.

---

## P0-3 — Synchronize command documentation, then gate it

**Problem.** 15 commands missing from `.claude/README.md`, 13 from
`docs/WORKFLOW_COMMANDS.md`, 12 from both; the two documents disagree with each
other; `.claude/README.md` documents a "Workflows" category with no directory.

**Existing resource considered.** `ff6lab audit`'s `CheckMarkdownLinks` already
walks every `.md`.

**Why insufficient.** It validates that links resolve, not that the command set
is covered. Nothing detects a command that exists and is documented nowhere.

**Proposal.** Add the missing entries to both documents, remove or implement
the phantom "Workflows" category — **and add a `CheckCommandDocumentation` gate
to `internal/audit`** so the drift cannot recur. Documentation fixes without the
gate will drift again; that is exactly how the current state arose.

**Files affected.** `.claude/README.md`; `docs/WORKFLOW_COMMANDS.md`;
`internal/audit/audit.go` + a new `commands.go` and `commands_test.go`.

**Migration impact.** None.

**Test strategy.** Table-driven fixture-tree tests matching the existing
`internal/audit` pattern: a fixture with an undocumented command must produce a
finding.

**Telemetry.** None needed — this is a static check.

**Rollback.** Remove the check from `audit.Run`'s list.

**Priority: P0.** The gate is the deliverable; the text fix alone is not.

---

## P1-1 — Wire `internal/validate` into the graphics procedure

**Problem.** `CompareResolved` exists, is tested, and encodes the palette-index
lesson that produced a visible defect. **`go list` shows zero importers.**
`VALIDATE_GRAPHICS.md:16` says "Compare indexed pixels when possible" —
prescribing the failure mode the code prevents.

**Existing resource considered.** The playbook and the package both already
exist. Nothing new is needed.

**Why insufficient.** They have never been connected.

**Proposal.** Rewrite `VALIDATE_GRAPHICS.md` step 4 to require resolved-colour
comparison and name `internal/validate.CompareResolved`; expose it as
`ff6lab validate framebuffer` so the playbook has a runnable target. Note that
`ff6lab help` already advertises a planned `validate` command group — this
fills a slot the tool has already promised.

**Files affected.** `.claude/playbooks/VALIDATE_GRAPHICS.md`;
`cmd/ff6lab/main.go` + a new `validate.go`.

**Migration impact.** None; `internal/validate`'s API is unchanged.

**Test strategy.** `internal/validate` is already tested, including
`TestResolvedCatchesInvisibleInk`. Add a `cmd/ff6lab` subcommand test in the
existing `main_test.go` style. An importer count above zero is itself the
regression check.

**Telemetry.** `resource_type: "tool"` makes `ff6lab validate` usage visible.

**Rollback.** Revert the playbook; the subcommand is additive.

**Priority: P1.** The highest-value single wiring fix in the repository.

---

## P1-2 — `ff6lab workflow inventory` and `validate`

**Problem.** AUDIT-0001 built the orphan detection, reference graph, doc
coverage and declared-versus-actual comparison **by hand**, in throwaway
scripts under `local_artifacts/`. The next audit would rebuild them.

**Existing resource considered.** `internal/audit`.

**Why insufficient.** It has no model of the orchestration layer at all — no
registry, no reference graph, no concept of a command or agent.

**Proposal.** Build `internal/workflow` with `manifests/workflow-registry.json`
and `schemas/workflow-registry.schema.json`, exposing
`ff6lab workflow inventory` and `ff6lab workflow validate`. **Build these two
first and defer the event-recording subcommands** — they need no telemetry and
no behaviour change, and they subsume most of this audit.

`AUDIT-0001-baseline.json` is the first input; `validate` must be able to diff
against it.

**Files affected.** New `internal/workflow/`; new manifest and schema;
`cmd/ff6lab/main.go`; wire `validate` into `audit.Run`.

**Migration impact.** New manifest must be generated once and kept in sync — the
same maintenance shape as `manifests/experiments.json`.

**Test strategy.** Fixture-tree tests per the `internal/audit` convention: a
fixture with an orphaned agent, an undocumented command and a broken reference
must produce exactly three findings.

**Telemetry.** This *is* the telemetry substrate.

**Rollback.** Remove the subcommand and the `audit.Run` entry.

**Priority: P1.**

---

## P1-3 — Reconcile the experiment schema with the constitution

**Problem.** `CLAUDE.md` requires seven experiment fields. The schema has no
property at all for **controlled variables, required evidence, or stopping
condition**. `domain` appears in 12 of 52 manifest entries and is undefined in
the schema.

**Existing resource considered.** `CheckManifests` already validates manifests.

**Why insufficient.** It validates against a schema that omits three
constitutional fields, so it cannot detect their absence.

**Proposal.** Add `controlled_variables`, `required_evidence`,
`stopping_condition` and `domain` to `schemas/experiment.schema.json`, set
`additionalProperties: false`, and backfill the 52 existing entries.

**Files affected.** `schemas/experiment.schema.json`;
`manifests/experiments.json`; possibly `.claude/templates/EXPERIMENT.md` for
the `Falsifier`/`Falsifying outcome` naming drift.

**Migration impact.** **Non-trivial** — 52 entries need real values, not
placeholders, and the values must come from the Markdown records rather than be
invented. Make the new fields optional first, backfill, then require them.

**Test strategy.** Existing `CheckManifests` tests extend naturally.

**Telemetry.** None.

**Rollback.** Revert the schema; entries with extra fields still validate while
`additionalProperties` is absent.

**Priority: P1**, staged.

---

## P2 — Smaller repairs

| # | Problem | Proposal | Files |
|---|---|---|---|
| P2-1 | `/session-summary` cannot reach `SESSION.md`; 9 of 12 templates orphaned | Name the template in each record-producing command, exactly as `/checkpoint` does | 9 command files |
| P2-2 | `ff6lab state origin`/`sprites` absent from `ff6lab help`; "Planned command groups" advertises three unbuilt groups | Add the two subcommands to help; mark planned groups as planned or remove | `cmd/ff6lab/main.go` |
| P2-3 | `internal/audio/doc.go` claims SPC700/DSP/sequence code that does not exist | Correct the doc comment to state actual coverage | `internal/audio/doc.go` |
| P2-4 | `RUN_QUALITY_GATES.md` names no build step, no `ff6demo`, no `gui` tag, no asset scanning, though CI enforces all four | Bring the playbook in line with `.github/workflows/ci.yml` and `CLAUDE.md` | `.claude/playbooks/RUN_QUALITY_GATES.md` |
| P2-5 | `ASSET-GFX-0002` (Narshe tileset) links `CEN-GFX-0006` (mines interior) | Correct or add the census reference; consider a check that asset and census scenarios agree | `manifests/assets.json` |
| P2-6 | `OPEN_HYPOTHESES.md` stale since 2026-07-31; live hypotheses exist only in checkpoint prose. `ACTIVITY_LOG.md` missing all of Session 005 | Register the three live map-descriptor hypotheses; append Session 005 | 2 dashboards |
| P2-7 | `/correlate-static-runtime` never names `docs/correlations/` as its output location | State the output location | 1 command file |
| P2-8 | Unit IDs collide (Unit 12/17/18 each name two units) | Adopt session-qualified unit ids | `.claude/templates/CHECKPOINT.md`, `SESSION.md` |
| P2-9 | Restricted/rendered extension lists are closed sets | Consider content sniffing alongside extensions | `internal/audit/audit.go` |

---

## Candidate resources — recommended, deferred, rejected

### Recommended (2 of 9)

| Candidate | Verdict | Reason |
|---|---|---|
| `WORKFLOW_LIFECYCLE.md` shared contract | **Recommend, after P1-2** | The `_shared/` pattern is proven — `EVIDENCE_STANDARD.md` and `STOPPING_RULES.md` are referenced and non-orphaned. But a lifecycle contract with no `ff6lab workflow` behind it is a document asking Claude to remember, which this audit measured as unreliable. |
| `/validate-capabilities` | **Recommend as a subcommand, not a command** | The need is real. It is `ff6lab workflow validate` (P1-2), which is testable and cannot be forgotten. A `.claude` command would be a third way to ask for the same thing. |

### Deferred (5 of 9)

`/workflow-audit` and `/workflow-report` — premature. Until `ff6lab workflow`
exists and has been run, a command wrapping it would be orchestration over
nothing. Revisit after P1-2.

`/start-unit` and `/close-unit` — the need (lifecycle boundaries) is real, but
hooks (telemetry design, mechanism 2) capture the same boundaries **without
depending on anyone invoking a command**. Prefer the mechanism that cannot be
skipped. Revisit only if hooks prove unavailable.

`workflow-telemetry` skill — subsumed by `WORKFLOW_LIFECYCLE.md` plus the
tooling. Adding a skill for it would create the same overlap the audit found
between `census-observer` and `ff6-content-census`.

### Rejected (2 of 9)

`workflow-governance-auditor` agent — **rejected.** The audit's central finding
is that 12 of 13 existing agents have no routing. Adding a fourteenth
unrouted agent would deepen the exact defect. Route the existing agents first
(P0-1); reconsider only if a governance specialist is still missing afterward.

`capability-validator` skill — **rejected.** Capability honesty is a property
of code, checked by `ff6lab workflow validate`. A skill cannot verify that a Go
package exists more reliably than a Go program can.

---

## The ten required items, addressed

1. **Agent routing** — P0-1.
2. **README/documentation synchronization** — P0-3, gated.
3. **`census-observer` vs `ff6-content-census`** — *No merge recommended.*
   `CLAUDE.md` names `ff6-content-census` directly, so merging would break the
   constitution's own reference. The defect is the absent routing distinction,
   not the duplication. Fix: state in both skills which is entry-point
   (`ff6-content-census`, constitution-level, breadth strategy) and which is
   the per-experiment procedure (`census-observer`), and state how
   `/census-observations`, `/register-system` and `/update-coverage` differ.
   Included in P2 scope; **not a Merge Candidate.**
4. **`/trace-dma` live validation** — one operator session: run
   `probe dma-trace` over a map load, feed one line through `ParseTraceLine`,
   confirm against a known transfer. Validates three of four surfaces at once.
   **Requires an emulator; explicitly out of scope for this audit.**
5. **Command aliases** — **Keep unchanged.** Deliberate, self-documenting, and
   correct. They must be normalized in telemetry via `alias_invoked_as`, not
   removed.
6. **Recently added recovery commands** — the six orphans are the best-written
   commands in the set. Fix is documentation coverage (P0-3), not the commands.
7. **Command lifecycle telemetry** — telemetry design record; build order
   P1-2 → hooks → lifecycle contract.
8. **Historical usage uncertainty** — 101 `Unknown` is the honest floor.
   Derived events (telemetry mechanism 1) can promote some retroactively.
   **No resource may be removed on this uncertainty.**
9. **Workflow-to-demo outcome measurement** — the discovery → demo link is
   prose only and unenforced. Require `IMPLEMENT_DISCOVERY.md` to name the
   readiness row it advances, and extend `CheckReadinessSummary` to verify every
   discovery is referenced by at least one row. Folds into P1-2.
10. **Automatic detection of duplicated facts and stale summaries** — partly
    solved: `CheckReadinessSummary` already asserts a summary against the rows
    it summarizes (`54a5200`). Generalize that pattern to
    `indexes/SESSIONS.md` (which would have caught the `SESSION_004` gap) and
    to `dashboards/STATISTICS.md`.

---

## Recommended implementation order

1. **P0-1** agent routing — highest leverage, no new resource
2. **P0-3** documentation sync **plus its gate**
3. **P0-2** `/audit-project` report-only mode
4. **P1-1** wire `internal/validate`
5. **P1-2** `ff6lab workflow inventory` + `validate`
6. **P1-3** experiment schema reconciliation, staged
7. **P2** repairs, cheapest first
8. Telemetry events, hooks, `WORKFLOW_LIFECYCLE.md` — only after P1-2
9. `/trace-dma` live validation — next authorized operator session

Items 1–4 are text and small Go changes with no migration cost, and they
address the three highest-risk findings. **Nothing in this plan requires an
emulator except item 9**, which is explicitly deferred.
