# Codex workflow handoff

**Written 2026-08-02.** State frozen for review. No further remediation was
performed after this document was written.

## Repository identity

```text
repository:      git@github.com:bensabler/ff6-decompile.git
current branch:  feature/workflow-contract
current HEAD:    4d12c29  (R11 specification, documentation only)
parent chain:    4d12c29 -> f01817f -> 1d20edf -> 69d99ae (audit2 HEAD)
                 69d99ae -> ... -> 93f7d03 (AUDIT-0001 closure)
                 93f7d03 -> ... -> 581ddbc (demo/whelk-content-parity)
```

**Worktrees** (all on one clone; the replay worktrees are disposable):

| Path | Checked out | Purpose |
|---|---|---|
| `ff6-decompile` | `93f7d03` `maintenance/workflow-observability` | AUDIT-0001 branch |
| `ff6-audit2` | `69d99ae` `maintenance/workflow-observability-audit2` | AUDIT-0002 branch |
| `ff6-r12` | `4d12c29` `feature/workflow-contract` | R12/R14 implementation |
| `ff6-replay-581` | `581ddbc` detached | AUDIT-0002 generator replay |
| `ff6-replay-93f` | `93f7d03` detached | AUDIT-0002 contamination check |

**Remote refs before this handoff pushed anything:**

```text
297ba88  refs/heads/main
581ddbc  refs/heads/demo/whelk-content-parity
caf6ff3  refs/heads/demo/new-game-to-whelk
```

Neither audit nor feature branch existed on the remote.

## Why FF6 work was paused

Decompilation stopped because the project's own operating layer could not be
trusted, and that had to be fixed before more long autonomous runs.

- **Claude was not reliably using the project's commands, skills and agents.**
  135 resources exist. Most were never demonstrably invoked.
- **Automatic specialist routing was not functioning.** Pre-registered probes
  showed that omitting `subagent_type` yields `general-purpose`, never a
  project specialist. No command, skill, playbook, rule or template names any
  of the 13 agents.
- **Some advertised capabilities had no operational backend.**
  `/validate-audio`, `/recover-sequence` and `/trace-spc-command` prescribe
  procedures against code that does not exist. `internal/validate` has zero
  importers while its playbook prescribes the failure mode it was written to
  prevent.
- **Artifact existence was used as invocation evidence.** The same reasoning
  produced one true and one false conclusion in AUDIT-0001, and no method
  present at the time could tell them apart.
- **Command and lifecycle behaviour was hard to audit.** Determining what
  actually ran required reading raw tool activity by hand.

The operator's own framing: with 135 resources, they should not have to audit
raw tool activity to learn whether the right ones were used.

## Audit history

### AUDIT-0001 — `93f7d03`, branch `maintenance/workflow-observability`

Inventoried the orchestration layer and found real defects. It also committed
several of the failures it was built to detect, and certified itself.

### AUDIT-0002 — `69d99ae`, branch `maintenance/workflow-observability-audit2`

Became necessary because AUDIT-0001's method could not distinguish its sound
conclusions from its unsound ones. Ran from an isolated worktree at the
verified AUDIT-0001 closure commit.

```text
Phase 11:            RAN (freeze, two bounded reviews, six disputes, two-stage closure)
completion_status:   complete   (per docs/workflows/AUDIT-0002-status.md at 69d99ae)
```

**Note the discrepancy with the handoff request**, which anticipated
`partial` / Phase 11 not run. Phase 11 did run and AUDIT-0002 recorded
`complete` at `69d99ae`. Codex should treat `AUDIT-0002-status.md` on the audit
branch as authoritative and this sentence as the reconciliation.

What "complete" means here is narrow: the audit's own eleven-phase process
finished and its closure was verified. It does **not** mean the remediation is
done — most of it is not.

### Major AUDIT-0001 claims, corrected

Claim-level detail: `docs/workflows/AUDIT-0002-claim-ledger.json`
(26 claims — 13 refuted, 7 confirmed, 4 partially supported, 2 unverifiable).

| AUDIT-0001 claim | AUDIT-0002 |
|---|---|
| 31 orphaned resources | **Refuted.** 7 textual orphans; the matcher never scanned repository-root files, and `PACKAGE_MANIFEST.json` lists 24 of the 31 |
| `/run-quality-gates` invoked (Confirmed) | **Refuted.** Absent from every preserved session transcript |
| `/correlate-static-runtime` invoked (Confirmed) | **Unverifiable.** Predates transcript coverage |
| `/bootstrap-ghidra` invoked (Confirmed via artifact) | **Confirmed** — but by a method that was invalid. Right answer, wrong reasoning |
| `/checkpoint`, `/session-summary` Probable | **Refuted as under-claims.** Both actually invoked |
| 3 of 4 rules orphaned | **Refuted.** Category error: rules are auto-activated by `paths:` glob |
| 9 templates unreferenced | **Partially supported.** All 9 are in `PACKAGE_MANIFEST.json`; 10 of 12 lack a *consumer* reference |
| All composite scores | **Refuted.** Hard-coded defaults and invalid proxies |
| All 28 playbooks `Keep` | **Refuted.** Produced by an unconditional `return` |
| Explicit agent routing absent | **Confirmed**, and better measured: 13 of 13 have zero routing-bearing inbound |
| `/trace-dma` four surfaces, `internal/validate` unimported, `/audit-project` mutation-only, unit-ID collisions | **Confirmed**, re-verified |

Left unverifiable: AUDIT-0001's Phase 4 agent process metrics (no transcripts
preserved); all invocation history before 2026-08-02 (no transcript coverage);
hook capability (never tested).

## Remediation implemented early

**Remediation began before the audit's remediation plan was reviewed by
anyone other than the operator.** Both commits are on
`feature/workflow-contract`, not on an audit branch.

### `1d20edf` — R12, execution contract

**Adds** `internal/workflow` (`contract.go`, `reconcile.go`) and
`schemas/workflow-contract.schema.json`.

**Enforces** a contract lifecycle — draft, displayed, approved, frozen and
hashed, executing, reconciled. After freezing, a requirement may not be
weakened, removed or reclassified; the amendment path preserves the original
and demands a reason plus renewed approval. `warn_only` is illegal on a
required requirement. Evidence rules: an artifact never proves an invocation;
`general-purpose` never satisfies a named specialist; a mentioned or
context-loaded skill earns nothing; a blank exit status is unverifiable, not a
pass.

**Tests** 38 cases, including negative acceptance tests N1, N2, N3, N4, N5, N7,
N8.

**Known limitations.** Pure logic only — gathers no observations, owns no
state, enforces nothing at runtime by itself.

**Status: provisional, awaiting Codex review.**

### `f01817f` — R14, deterministic reconciliation

**Adds** `internal/workflow/{evidence.go,store.go}`,
`cmd/ff6lab/workflow.go`, and a `workflow` entry in `ff6lab help`. Modifies
`cmd/ff6lab/main.go` and `internal/project/project.go`.

**Enforces** that the verdict is computed, never asserted: no subcommand
accepts a verdict as input, and `workflow close` exits non-zero on anything but
`complete`. Three evidence sources each report their own blindness rather than
returning an empty set. Transcript reading takes tool-call metadata only —
never prompts or conversation content. Durable state is restart-safe; the
approved contract is tracked, mutable state is ignored.

**Tests** 57 in `internal/workflow` plus 17 in `cmd/ff6lab`, including N6 and
N9. All nine negative acceptance tests now pass across R12 and R14.

**Known limitations.** Nothing yet *produces* contracts, so an orchestrator can
still work without opening one. Enforcement is real but opt-in.

**Status: provisional, awaiting Codex review.**

### `4d12c29` — R11 specification, no implementation

Documentation only, self-identifying as unimplemented. Retained because it
records a capability question honestly rather than guessing at it.

**The complete operator workflow does not exist.**

## Current limitations

- **R11 is not implemented.**
- **The seven argument-free workflow commands do not exist.** `/research`,
  `/extract`, `/reconstruct`, `/implement`, `/validate`, `/continue`,
  `/review` are proposals only.
- **The visible Claude command surface has not been reduced.** 43 command
  files, 43 visible.
- **No command-surface hiding or migration mechanism has been verified.**
  Command files carry no frontmatter; `.claude/commands/` has no
  subdirectories; no project document describes a hiding mechanism.
  Subdirectory namespacing is a common convention but is **untested here**.
- **R13 receipts are not implemented.** R14 validates a receipt's claimed
  verdict; nothing generates one.
- **Workflow contracts are not the mandatory default path.**
- **The Claude-to-Codex migration has not occurred.**
- **FF6 decompilation remains paused** pending review and migration decisions.

## Claude-to-Codex migration boundary

- **`.claude/` remains the source Claude configuration.** It is authoritative
  for Claude and is the input to migration, not the output.
- **`.claude/` resources are not assumed to work natively in Codex.** No
  compatibility claim is made for any command, skill, agent, playbook, rule or
  template.
- **Source Claude files must be preserved during migration.** Translate into
  new artifacts; do not convert `.claude/` in place.
- **Migration must use the reviewed migration/translation process**, not an
  ad-hoc copy.
- **Migrated resources require semantic validation**, not merely successful
  format conversion. A translated skill that loads is not a skill that works.
- **Do not migrate blindly.** As classified in
  `docs/workflows/AUDIT-0002-corrected-baseline.json`: 1 deprecated,
  3 aliases, 3 blocked with absent backends, 11 unverified, 13 partial.
  Backend-absent resources must not arrive in Codex looking operational.
- **Do not add the seven proposed workflows on top of every existing command.**
  The point of R11 is a reduced surface. Adding seven to 43 produces 50 and
  makes the original problem worse.

## Exact next action for Codex

**The first Codex task is read-only. Modify nothing.**

1. Inspect the pushed branches `maintenance/workflow-observability-audit2` and
   `feature/workflow-contract`.
2. Compare the pre-remediation audit state at `69d99ae` against R12 (`1d20edf`)
   and R14 (`f01817f`).
3. Review R12 and R14 for correctness and compatibility — including whether the
   evidence rules match how Codex actually records invocations, which may
   differ from the Claude harness the transcript reader was written against.
4. Inspect `.claude/` and the AUDIT-0002 workflow findings under
   `docs/workflows/`.
5. Produce a Claude-to-Codex migration plan.
6. Recommend whether R12 and R14 should be **retained, revised, or reverted
   through new commits** — never by rewriting `1d20edf` or `f01817f`.
7. **Stop before modifying anything.**

After approval, Codex may run the official Claude-to-Codex migration process on
a **new migration branch**.

## Reading order for a fresh agent

1. This document.
2. `docs/workflows/AUDIT-0002-status.md` — audit outcome.
3. `docs/workflows/AUDIT-0001-errata.md` — what changed and why.
4. `docs/workflows/AUDIT-0002-remediation-plan-v2.md` — R1-R14, with the
   confirmation-risk disclosure on R11-R14.
5. `docs/workflows/AUDIT-0002-self-compliance.md` — plan versus actual,
   including nine errors and where each was caught.
6. `docs/workflows/R11-command-surface-specification.md` — the untested
   capability question.

## Standing caveats

- **The audits were self-audits.** AUDIT-0002 re-audited work by the same
  author, and its Phase 11 reviewers ran on the same underlying model.
  Correlated review, not independent review. Weight accordingly.
- **Transcript coverage is 2026-08-02 only.** Invocation history before that
  date is unverifiable in principle from what survives.
- **`archive verify` has not run in any audit session** — `FF6_ROM` was never
  set. It is recorded as not run, never inferred as passing.
