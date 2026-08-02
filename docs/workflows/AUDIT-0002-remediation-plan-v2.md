# AUDIT-0002 — corrected remediation plan v2 (Phases 9-10)

**Nothing here is implemented.** Every acceptance test below is a
**specification for later remediation**; AUDIT-0002 ran none of them, because
each depends on behaviour that does not yet exist.

## Verdicts on AUDIT-0001's proposals

| AUDIT-0001 proposal | Verdict | Reason |
|---|---|---|
| P0-1 agent routing table | **approve with revision** | Finding confirmed and better measured. But it cannot create *automatic* selection — no such mechanism exists — and its "no migration impact" claim is **wrong** |
| P0-2 `/audit-project` report-only | **approve unchanged** | Re-verified |
| P0-3 doc sync + gate | **approve with revision** | Prefer one canonical source generating both surfaces over two hand-maintained inventories |
| P1-1 wire `internal/validate` | **approve with revision** | Acceptance must be behavioural, not an importer count |
| P1-2 `ff6lab workflow` | **approve with revision** | Must exclude audit outputs from its corpus and scan root files |
| P1-3 experiment schema | **supersede** | Rests on a category error — the schema indexes a manifest, not the record |
| P2-6 dashboard updates | **defer pending evidence** | No contract establishes the required cadence |
| P2-8 unit ids in templates | **approve with revision** | Propagate `unit_key` to every surface, not only checkpoints |
| Discovery→demo linkage | **supersede** | Requires an applicability decision, not universal representation |
| Rules "Clarify Routing" | **reject** | Category error — rules are auto-loaded by path scope |
| Telemetry mechanism 1 (derived events) | **reject** | Refuted by AUDIT-0001's own two contradictory outcomes |
| All 28 playbook `Keep` | **reject** | Produced by an unconditional return; no playbook was assessed |
| All composite scores | **reject** | Hard-coded defaults and invalid proxies |

## Corrected remediation items

Each carries an acceptance test. **None was executed.**

### R1 — Agent routing (P0)

**Defect.** All 13 agents have zero routing-bearing inbound.
`ORCHESTRATE_RESEARCH.md` step 6 is the word "Delegate."

**Fix.** Domain → agent routing table in `ORCHESTRATE_RESEARCH.md`; name the
relevant agent in each domain command as skills already are.

**Acceptance test (specification).** A static name-presence check is
*necessary at most*. Behavioural acceptance requires: ordinary operator task →
command or orchestrator → correct agent selected **without the operator naming
it** → authority boundaries honoured → bounded result integrated. Measure
selection accuracy, unnecessary delegation, duplicate work, context cost,
elapsed time, scope violations, integration quality.

**Migration impact: behavioural, not zero.** Explicit routing changes which
agent handles which work.

**Known constraint.** Omitting `subagent_type` yields `general-purpose`. The
table helps the orchestrator choose; it cannot make selection automatic.

### R2 — `/audit-project` report-only mode (P0)

**Defect.** Every documented execution path mutates. AUDIT-0002 could not
invoke it and performed its steps manually.

**Fix.** `/audit-project` and `/audit-project report` → report-only;
`/audit-project fix` → explicit mutation authorisation. Compare against keeping
mutation as default and choose deliberately.

**Acceptance test.** The next audit-class session invokes it and completes
with a clean worktree.

### R3 — One canonical command inventory (P0)

**Defect.** `.claude/README.md` omits 15 commands, `docs/WORKFLOW_COMMANDS.md`
omits 13, they disagree with each other, and README documents a "Workflows"
category with no directory.

**Fix.** One canonical source generating both surfaces, plus a
`CheckCommandDocumentation` gate. Two hand-maintained inventories plus a
consistency gate would preserve the drift mechanism.

**Acceptance test.** Adding a command file without regenerating fails the gate;
regenerating passes; both surfaces agree by construction.

### R4 — Wire `internal/validate` (P1)

**Defect.** Zero importers, non-test and test alike, while
`VALIDATE_GRAPHICS.md:16` prescribes the failure mode `CompareResolved`
prevents.

**Fix.** Rewrite step 4 to require resolved-colour comparison and name the
function; expose `ff6lab validate framebuffer` (help already advertises a
planned `validate` group).

**Acceptance test.** **An importer count above zero is not sufficient.** Must
prove: indexed identity behaviour; resolved-colour mismatch detection;
invisible-ink regression; the CLI path; and a real demo or validation-pipeline
consumer.

### R5 — `ff6lab workflow inventory` / `validate` (P1)

**Fix.** Static inventory and validation, modelled on `internal/audit`. Kept
separate from event telemetry.

**Acceptance test.** Reproduces AUDIT-0002's corrected figures — 7 textual
orphans, 38 routing-bearing orphans, 13 of 13 agents unrouted — from a clean
run; **excludes `docs/workflows/` from its corpus**; **scans repository-root
files**; reports the three metrics separately; treats auto-activated rules as
correct rather than orphaned.

### R6 — Capability-honesty corrections (P2)

| Item | Fix | Acceptance |
|---|---|---|
| `internal/audio/doc.go` overstates | Correct the comment to actual coverage | Comment matches the package tree |
| `ff6lab state origin`/`sprites` absent from help | Add them; mark planned groups as planned | Every implemented subcommand appears in help |
| `/validate-audio`, `/recover-sequence`, `/trace-spc-command` | Mark backend-absent explicitly in the command text | No command promises tooling that does not exist |
| `RUN_QUALITY_GATES.md` names no build step | Align with CI and `CLAUDE.md` | Playbook names every gate CI enforces |
| `/correlate-static-runtime` output path | Name `docs/correlations/` | Declared output path exists |

### R7 — `/trace-dma` live validation (P1, needs an operator session)

One run of `probe dma-trace` over a map load, feeding one line through
`ParseTraceLine`, validates the live probe, the parser and provenance together.
**Requires an emulator. Out of scope for every audit so far.**

### R8 — Command lifecycle contracts (P2) — **new**

**Defect.** `/checkpoint` declares it writes `dashboards/CURRENT_FOCUS.md`,
making it unusable inside any bounded-scope session. `/session-summary`
declares **no output path at all**.

**Fix.** Declare an explicit, complete output path set for every
record-producing command, and give `/checkpoint` a bounded mode that omits the
dashboard write.

**Acceptance test.** Every command's declared write set is complete and
machine-checkable against what it actually writes.

### R9 — Corpus separation (P1) — **new**

**Defect.** AUDIT-0001's baseline enumerates all 135 resource ids, so replaying
its own analysis returns **zero** orphans. The audit destroyed its own metric.

**Fix.** Explicit path exclusion of audit outputs from any validator corpus.

**Acceptance test.** Adding a report that names every resource does not change
the orphan count.

### R10 — `CheckMarkdownLinks` walks gitignored paths (P2) — **new**

**Defect, found by the gate failing on AUDIT-0002's own evidence.**
`internal/audit.CheckMarkdownLinks` walks every `.md` in the tree, skipping only
`.git` and `node_modules`. It therefore walks `local_artifacts/`, which is
gitignored and is where the project's own convention says preserved evidence
lives.

Freezing a copy of `LATEST.md`, `SESSIONS.md` and `STATISTICS.md` as evidence
produced **15 findings** and failed `ff6lab audit`, purely because those copies'
relative links do not resolve from their new location. The content was correct;
the location was fatal.

**This penalises the evidence-preservation behaviour the project requires.**
AUDIT-0001 never hit it because it never froze evidence.

**Fix.** Skip gitignored paths in `CheckMarkdownLinks` (or at minimum
`local_artifacts/`), so preserved evidence cannot break the gate.

**Acceptance test.** A frozen copy of a canonical document placed under
`local_artifacts/` does not change `ff6lab audit`'s result.

**Interim workaround used here:** frozen documents are stored with a
`.md.frozen` extension so the walker does not see them. That is a workaround,
not a fix — it disguises evidence to satisfy a checker.

## Phase 10 — safeguards for long research runs

**Default position held: no new permanent agent, no new permanent skill.**

The lifecycle: plan → execute → **freeze evidence** → independent verification
→ synthesize → promote or reject → implement → integrate → validate →
checkpoint.

**The control AUDIT-0001 most conspicuously lacked is executor/verifier
separation.** It authored its findings, generated its own baseline, certified
its own closure, and pre-declared completion before the closing phase ran.

Can existing resources absorb these controls?

| Control | Owner | New resource needed? |
|---|---|---|
| Evidence frozen and hashed before interpretation | `EVIDENCE_STANDARD` | **No** — extend it |
| Stopping conditions, context boundary | `STOPPING_RULES` | **No** |
| Executor/verifier separation | `research-orchestrator` + `verification-engineer` | **No** — both exist and are unrouted, which R1 fixes |
| Falsifiable question, controlled variables | `experiment-designer` | **No** |
| Adversarial challenge | `quality-reviewer` | **No** |
| Globally unique `unit_key` | checkpoint template | **No** — add the field |
| Performed-actions-only checkpoints, no future-phase credit | checkpoint template | **No** — add a plan-vs-actual section |
| Restart-safe partial state | checkpoint template + status record | **No** |

**Every control lands on an existing resource.** The reason they did not
function is not absence — it is that `verification-engineer` and
`quality-reviewer` are two of the 13 agents nothing routes to. **R1 is the
prerequisite for the entire safeguard set**, which is why it is P0.

No new shared contract is proposed. None of the five conditions for one is met:
the rules can be absorbed by existing contracts, and a new unrouted contract
would become exactly the kind of unmeasured resource this audit exists to find.

## Recommended order

1. **R1** agent routing — prerequisite for every safeguard
2. **R3** canonical command inventory + gate
3. **R2** `/audit-project` report-only
4. **R4** wire `internal/validate`
5. **R9** corpus separation, with **R5** `ff6lab workflow`
6. **R6** capability-honesty corrections
7. **R8** command lifecycle contracts
8. Telemetry: transcript extraction, then hooks only after capability testing
9. **R7** `/trace-dma` live validation — next authorised operator session

Items 1-4 are text and small Go changes. **Only item 9 requires an emulator.**
