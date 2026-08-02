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

### R5 — `ff6lab workflow inventory` / `validate` — **static only** (P1)

**Scope, stated explicitly:** static inventory and validation **only**. Durable
workflow state, reconciliation and closure verdicts belong to **R14**, not
here.

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

---

# The operator-level problem

Everything above corrects the record. None of it changes what the operator has
to do. With 135 resources, they should not have to read raw tool activity to
learn whether the right ones were used.

**Measured result: of the 43 existing commands, `operator_workflow` count is
zero.** Every one is an internal helper, a diagnostic tool, an alias, or
deprecated. There is no outcome-named surface at all, which is why sequencing
falls to the operator.

`outcome_status` is `no_evidence` for **130 of 135** resources — one
`validated`, three `partial`, one `failed`. Nothing ties a durable artifact to
the resource that produced it. That is the evidentiary basis for R12-R14.

## Command classification — two independent dimensions

`surface_role` and `operability` are separate: a command may be a correct alias
that works, or a diagnostic tool whose backend is absent. Full per-command
detail with basis is in `AUDIT-0002-corrected-baseline.json`.

| `surface_role` | n | | `operability` | n |
|---|---|---|---|---|
| `operator_workflow` | **0** | | `verified` | 15 |
| `internal_helper` | 33 | | `partial` | 13 |
| `diagnostic_manual` | 6 | | `unverified` | 11 |
| `alias` | 3 | | `blocked_backend_absent` | **3** |
| `deprecated` | 1 | | `not_applicable` | 1 |

**Blocked, and must never present as operational:** `/validate-audio`,
`/recover-sequence`, `/trace-spc-command`.
**Deprecated:** `/bootstrap-v4`. **Aliases:** the three `/recover-*` spellings.

## R11 — Operator-facing workflow surface (new, P0)

Seven argument-free interactive commands. **No public `/plan`, `/execute`,
`/verify` or `/close`** — those are internal lifecycle stages owned by the
workflow engine, never operator decisions.

| Command | Kind | Behaviour |
|---|---|---|
| `/research` `/extract` `/reconstruct` `/implement` `/validate` | outcome | Interactively resolve bounded scope → build and **display** the contract → **request approval** → freeze → execute |
| `/continue` | lifecycle | Resumes the **currently open** contract. Shows workflow, completed and remaining requirements, unresolved failures, exact next authorised action. Never asks the operator to reconstruct domain or phase |
| `/review` | lifecycle | **Read-only.** Reports state, evidence, execution compliance, missing requirements, recommended next action. Never advances the workflow |

### Measurable surface reduction

Not seven commands on top of 43. Each class has a stated disposition:

| Class | Disposition | Before | After |
|---|---|---|---|
| Primary operator workflow | visible, manually invocable | 0 | **7** |
| Advanced diagnostic | moved to an identified advanced surface | 6 | 6 (advanced) |
| Internal helper | hidden from the normal operator surface | 33 | 0 visible |
| Alias | hidden or removed after a documented migration period | 3 | 0 visible |
| Deprecated | removed after compatibility review | 1 | 0 |
| Blocked, backend absent | visibly blocked; cannot present as working | 3 | 3 (blocked) |
| **Visible to the operator by default** | | **43** | **7** |

**Acceptance criterion.** After remediation the default surface lists exactly 7
commands; 6 more appear only on the advanced surface; 33 helpers are not
listed; 3 aliases and 1 deprecated command are gone from the default surface;
and the 3 blocked commands are shown as blocked with their absent backend named.

### `/extract` acceptance specification

Specification for later remediation. **Not executed.**

1. Operator invokes `/extract` with **no arguments**.
2. Claude interactively resolves a bounded point-A-to-point-B range.
3. Claude **displays** the planned required and conditional specialists.
4. Operator **approves** the plan.
5. Exact project agents are invoked by **exact `subagent_type`**.
6. Required extraction backends run.
7. Outputs receive provenance and hashes.
8. Transcript evidence is reconciled against the frozen contract.
9. Missing mandatory work **blocks completion**.
10. A concise execution receipt is produced.
11. The workflow **does not drift** into reconstruction, implementation or demo
    integration without separate authorisation.

### Negative acceptance tests — enforcement must reject these

Each reproduces a failure mode actually observed in AUDIT-0001 or AUDIT-0002.

| # | Scenario | Required outcome |
|---|---|---|
| N1 | Required specialist not invoked | Cannot become `complete` |
| N2 | `general-purpose` invoked instead of the required specialist | Requirement unsatisfied; cannot become `complete` |
| N3 | Required skill appears in prose or context only | No invocation credit; cannot become `complete` |
| N4 | Conditional specialist marked not applicable | Applicability reason preserved; verifier accepts **only** when the contract's applicability rule supports it |
| N5 | Required backend exits non-zero | `partial` or `failed` per the declared failure policy |
| N6 | Required invocation cannot be verified | `unverifiable` or `partial`; **never `complete`** |
| N7 | Frozen contract modified without an approved amendment | Reconciliation **rejected** |
| N8 | Expected artifact exists, invocation evidence absent | Artifact-compatible activity recorded; **invocation not credited**; completion blocked where invocation was required |
| N9 | Receipt says `complete`, `ff6lab` says `partial`/`failed` | Receipt validation **fails**; the `ff6lab` verdict is authoritative |

N8 is the direct enforcement of AUDIT-0001's `/bootstrap-ghidra` and
`/correlate-static-runtime` errors — the same artifact logic produced one true
and one false conclusion, and no method could tell them apart.

## R12 — Frozen execution contract and reconciliation (new, P0)

**Lifecycle.** `draft → displayed to operator → operator approved → frozen and
hashed → executed → reconciled → complete | partial | failed | unverifiable`.

**Requirements may not be weakened, removed or changed after execution
begins.** Any post-approval change must preserve the original contract, create
a formal amendment with a reason, receive renewed operator approval, and be
re-frozen and re-hashed before execution continues.

**A required agent, skill, command, backend operation, output or validation
that was skipped may never be retroactively reclassified optional or
not-applicable.**

### Per-requirement declaration

Every required or conditional resource declares:

```json
{
  "resource_id": "graphics-researcher",
  "resource_type": "agent",
  "requirement": "required",
  "execution_mode": "explicit_agent_call",
  "evidence_rule": "matching transcript Agent call using the exact subagent_type",
  "failure_policy": "block_completion"
}
```

`execution_mode` ∈ `explicit_agent_call` · `explicit_skill_call` ·
`deterministic_backend` · `operator_action` · `context_only` ·
`not_applicable`. Conditional requirements additionally declare an
applicability rule.

**Evidence rules.**

- Required agents must be invoked by **exact `subagent_type`**.
  **`general-purpose` does not satisfy a named-specialist requirement.**
- A skill earns invocation credit **only** from a transcript or other approved
  direct record proving the exact invocation. Available, mentioned or
  context-loaded earns nothing.
- A deterministic backend earns credit only from the preserved command, its
  exit status, and applicable output evidence.
- **Artifact existence alone never proves invocation.**
- A low-level slash command may be mandatory **only** where direct invocation
  is supported and preserved as evidence. **Where nested slash-command
  invocation cannot be reliably proven, move the procedure into an internal
  skill or deterministic backend — do not reproduce the steps manually and
  credit the command.** This is the rule AUDIT-0001 broke.

R1's routing table becomes an **input** to this: each workflow declares which
agents it requires. R1 is retained as a finding and **subsumed by R12**.

## R13 — Execution receipt (new, P0)

Human-readable, one per run: planned versus invoked agents; planned versus
observed skills and commands; backend operations with exit statuses; required
outputs; evidence references and hashes; validation results; missing, skipped,
failed or not-applicable requirements; the **deterministic** workflow verdict;
one proposed next action.

**Purpose: the operator audits exceptions and decisions, not routine tool
calls.**

**The receipt reflects the `ff6lab` verdict — it never declares completion
itself, and its own existence is never proof that execution occurred.**

Tracked, per run:

```text
docs/workflows/runs/<workflow-id>/contract.json
docs/workflows/runs/<workflow-id>/receipt.md
```

Raw transcript extracts, detailed reconciliation data, diagnostic logs and
bulky evidence stay ignored under `local_artifacts/`. **Proposed only — no
`runs/` directory is created by this audit.**

## R14 — `ff6lab workflow` is authoritative for closure (new, P0)

**Claude may propose and execute work. Claude may not decide in prose that a
workflow is complete.**

`ff6lab workflow` deterministically reconciles the frozen contract against
transcript evidence, explicit agent calls, explicit skill calls where
observable, backend commands and exit statuses, expected outputs, evidence
hashes, validation results, and approved not-applicable decisions — then
computes the verdict:

```text
complete | partial | failed | unverifiable
```

**Claude reports that verdict and may not override it.**

Durable state: `start`; planned requirements; current phase; expected-versus-
observed reconciliation; `partial` and `failed` status; restart-safe
continuation; closure validation. Deterministic and testable in Go, modelled on
`internal/audit`.

The public slash command stays argument-free even though it drives
deterministic `ff6lab` operations internally.

**Acceptance:** all nine negative tests N1-N9 pass, and the verifier's verdict
disagreeing with a receipt causes receipt validation to fail (N9).

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

## Recommended order — all 14 items

Ordered so the operator-facing outcome arrives before the cosmetic corrections.

| # | Item | Why here |
|---|---|---|
| 1 | **R12** frozen contract, execution modes, reconciliation | The enforcement substrate. Subsumes R1 |
| 2 | **R14** `ff6lab workflow` authoritative closure | Without a deterministic verdict, R12 is prose |
| 3 | **R11** operator-facing surface, 43 → 7 visible | The operator-value item; needs R12/R14 beneath it |
| 4 | **R13** execution receipt | Depends on R14's verdict |
| 5 | **R8** command lifecycle contracts | Declared write sets are required by R12 |
| 6 | **R3** canonical command inventory + gate | Feeds R11's visibility classes |
| 7 | **R9** corpus separation | Cheap; prevents the self-contamination defect recurring |
| 8 | **R5** `ff6lab workflow inventory`/`validate`, static only | Natural companion to R9 |
| 9 | **R2** `/audit-project` report-only | Unblocks the next audit-class session |
| 10 | **R4** wire `internal/validate` | Highest-value single wiring fix |
| 11 | **R6** capability-honesty corrections | Text-level, independent |
| 12 | **R10** `CheckMarkdownLinks` skips ignored paths | Small Go fix; unblocks evidence preservation |
| 13 | Telemetry: transcript extraction, then hooks **only after** capability testing | Depends on R14 |
| 14 | **R7** `/trace-dma` live validation | **Requires an emulator and an authorised operator session** |

**R1 does not appear as a separate item — it is subsumed by R12**, which
requires exact-`subagent_type` invocation and verifies it, rather than merely
listing agent names in a playbook.

Items 5-12 are text and small Go changes. **Only item 14 requires an
emulator.**
