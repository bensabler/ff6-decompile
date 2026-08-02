# Codex Workflow Compatibility Review

```text
review date: 2026-08-02
reviewing environment: Codex
branch reviewed: codex/workflow-foundation
reviewed HEAD: 78a8936ed747650370fb957e11e2a4de0dfecd13
review mode: read-only
repository modifications during review: none
```

## Executive conclusion

**Migration verdict: ready after named corrections.**

The workflow foundation contains useful work, but provisional R12 and R14 are
not yet trustworthy enough to resume FF6 reconstruction. The preserved import
snapshot is a sound archive, not a safe wholesale migration commit. It should
be migrated selectively only after the contract, approval, evidence-scoping,
output-enforcement, and closure defects identified below are corrected.

Both provisional implementation units receive the same disposition:

- R12: **retain with revisions**.
- R14: **retain with revisions**.

## Repository and handoff verification

The reviewed repository state was:

```text
path: /Users/benjaminsabler/Learning/GitHub/FF6-Reverse-Engineering/ff6-decompile
branch: codex/workflow-foundation
HEAD: 78a8936ed747650370fb957e11e2a4de0dfecd13
upstream: origin/feature/workflow-contract
working tree: clean
78a8936 is an ancestor of HEAD: yes
```

The worktrees observed during the review were:

| Worktree | HEAD | Branch |
|---|---|---|
| `ff6-decompile` | `78a8936` | `codex/workflow-foundation` |
| `ff6-audit2` | `69d99ae` | `maintenance/workflow-observability-audit2` |
| `ff6-r12` | `78a8936` | `feature/workflow-contract` |
| `ff6-replay-581` | `581ddbc` | detached |
| `ff6-replay-93f` | `93f7d03` | detached |

Live remote-ref verification found:

```text
feature/workflow-contract                  78a8936
maintenance/workflow-observability-audit2 69d99ae
```

### Commit chain and changed-file scope

The parent chain is linear and correct:

| Commit | Parent | Changed-file scope |
|---|---|---|
| `69d99ae` | `b1b0f4a` | Modified `dashboards/STATISTICS.md`, `docs/checkpoints/LATEST.md`, `docs/workflows/AUDIT-0002-status.md`, and `indexes/SESSIONS.md` |
| `1d20edf` | `69d99ae` | Added R12 contract/reconciliation Go code, tests, and schema: five files |
| `f01817f` | `1d20edf` | Added R14 evidence/store/CLI code and tests; modified CLI wiring/help: seven files |
| `4d12c29` | `f01817f` | Added only `docs/workflows/R11-command-surface-specification.md` |
| `78a8936` | `4d12c29` | Modified `dashboards/CURRENT_FOCUS.md`; added `docs/handoffs/CODEX-WORKFLOW-HANDOFF.md` |

### Status disagreements

The observable lifecycle facts are:

- `docs/workflows/AUDIT-0002-status.md` records `complete`, completed phases
  `0` through `11`, no remaining phases, and a verified closure candidate.
- `docs/workflows/AUDIT-0002-claim-ledger.json` retains top-level
  `completion_status: partial`.
- `docs/workflows/AUDIT-0002-corrected-baseline.json` retains
  `completion_status: partial`, completed phases `0` through `10`, and remaining
  phase `11`.

The **most likely interpretation, not verified history**, is that the two JSON
artifacts are generation-time snapshots whose lifecycle metadata was not
synchronized after Phase 11. Repository evidence proves the disagreement; it
does not explicitly prove why the artifacts were left in that state. Their
claim and resource data remain useful, but their top-level lifecycle fields
must not be treated as current without resolving that provenance question.

Additional disagreements and stale surfaces:

- `docs/handoffs/CODEX-WORKFLOW-HANDOFF.md` records current HEAD as `4d12c29`,
  although the document was added by `78a8936`. Its current-branch and worktree
  table no longer describe the reviewed state.
- The handoff's remote-ref subsection is explicitly historical. By contrast,
  `pushed: false` in `AUDIT-0002-status.md` and "Not pushed" in
  `AUDIT-0002-self-compliance.md` need clear closure-time qualification if they
  are to coexist with the branches now present on the remote.
- `dashboards/CURRENT_FOCUS.md` says "R11 ... does not exist," while an R11
  specification does exist. The accurate statement is that R11 is not
  implemented.
- Claim totals agree across the reviewed narrative surfaces: 26 claims -- 13
  refuted, 7 confirmed, 4 partially supported, and 2 unverifiable.

## R12 review

**Recommendation: retain with revisions.**

### Strengths

- Establishes an explicit contract lifecycle and frozen content hash.
- Separates artifact presence from resource invocation.
- Uses exact specialist matching; `general-purpose` cannot satisfy a named
  agent requirement.
- Treats a blank backend exit status as unverifiable, never as success.
- Computes verdicts instead of accepting them from callers.
- Preserves amendment lineage and rejects obvious post-freeze content changes.
- Includes important negative cases and restart-oriented reconciliation tests.

### Defects

#### Approval is represented, not evidenced

`state: approved` can be supplied in JSON or reached through `Advance()`
without an actor, timestamp, event identity, approval artifact, or other proof
that the operator saw and accepted the displayed contract. `ApprovedBy` on an
amendment is caller-supplied text and does not solve this problem.

#### The freeze invariant has bypasses

`Advance()` can move `approved` to `frozen` without calling `Freeze()`, leaving
the contract in a frozen state without a frozen hash. Conversely,
`VerifyFrozen()` checks for a matching hash but does not require a lifecycle
state consistent with a frozen or executing contract. Not every API path
therefore preserves the documented invariant.

#### Go validation and schema validation diverge

Go validation does not enforce all constraints represented by the JSON schema,
including the workflow enum, workflow-ID pattern, resource type, resource
type/execution-mode compatibility, timestamp structure, and output semantics.

#### Amendments can weaken material requirements

`weakeningCheck()` protects only required-resource presence, necessity, and a
small subset of mode changes. It permits weakening or changing:

- failure policy;
- execution mode in other ways;
- resource type and evidence rule;
- scope and workflow;
- required outputs and validations;
- stopping conditions and approval boundaries;
- conditional applicability and not-applicable decisions.

The amendment record also does not itself prove renewed operator approval.

#### Conditional and context-only semantics are unsafe

Conditional applicability is free text and is never evaluated. Any recorded
not-applicable reason is accepted. A required `context_only` requirement is
also accepted even though the source comments say that mode must not represent
genuinely required execution.

Unverifiable conditional requirements do not affect the verdict, even when
their failure policy is `block_completion` or `degrade_to_partial`.

#### Outputs and validations do not affect completion

`required_outputs` are not reconciled. `ObsArtifactPresent` has no execution
mode that consumes it, and required validation has no distinct enforcement
path. Collecting or hashing an artifact therefore does not make its presence a
condition of `complete`.

#### Evidence matching is too weak

Evidence-rule prose is not mechanically enforced. Backend reconciliation
accepts the first matching selector, so an earlier passing observation can
hide a later failure. Observations also carry no workflow, contract, run,
session, repository, working-directory, branch, or event binding.

### R12 conclusion

The lifecycle and reconciliation model are worth retaining, but the current
meanings of "approved," "frozen," "conditional," and "required output" promise
more than the implementation proves.

## R14 review

**Recommendation: retain with revisions.**

### Strengths

- Persists frozen contracts and mutable execution state separately.
- Refuses to start an unfrozen contract and detects stored contract-hash
  substitution.
- Treats missing evidence sources as incomplete rather than confirmed absence.
- Preserves real backend exit status, including a nil/blank status as
  unverifiable.
- Does not accept a caller-supplied verdict.
- Rejects a receipt whose claimed verdict disagrees with reconciliation.

### Run-scoping defects

An invocation from an unrelated earlier session can satisfy a current run if
its selector reaches the collector. Neither `Observation` nor `RunState` binds
evidence to:

- `workflow_id`;
- contract hash or immutable run ID;
- session and turn identity;
- repository identity or repository root;
- working directory;
- branch and HEAD;
- run start/end window;
- event or tool-use identity;
- evidence content hash.

Gate logs are located under a workflow-named directory, but their individual
rows do not prove which contract, session, repository, or run produced them.

### Claude transcript assumptions

Read-only metadata inspection of the configured Claude transcript tree found:

```text
53 JSONL files
0 direct JSONL children of ~/.claude/projects
7 project directories
```

The collector reads only direct `.jsonl` children of `~/.claude/projects`, so
it sees none of that real layout. Making the scan recursively global would be
unsafe because it would admit unrelated repositories and historical sessions.

The parser also assumes Claude-specific records:

- `message.content[].type == "tool_use"`;
- tool names `Skill` and `Agent`;
- input fields `skill` and `subagent_type`;
- Claude's JSONL organization.

Malformed JSON records are silently skipped without marking the evidence
source incomplete. Evidence references contain only the transcript basename,
which can collide. The tests explicitly treat a transcript containing a
malformed line as complete evidence, so this is current expected behavior, not
only an untested edge case.

### Codex compatibility

#### Agent invocation

Codex hook documentation exposes `SubagentStart` and `SubagentStop` events with
agent identity, and tool hooks can observe agent-spawn operations. This can
support agent evidence if the event is written into the workflow's own
run-scoped ledger.

See:

- <https://learn.chatgpt.com/docs/hooks>
- <https://developers.openai.com/codex/multi-agent>

#### Skill invocation

Codex documents skills as implicitly or explicitly selected, but the reviewed
official hook surface does not document a dedicated, stable skill-invocation
event. Exact skill invocation must therefore be classified **unverifiable**
unless a deterministic workflow wrapper emits its own event.

See:

- <https://developers.openai.com/codex/concepts/customization>
- <https://learn.chatgpt.com/docs/hooks>

#### Shell, backend, and tool execution

Pre- and post-tool hooks can expose tool name, input, response, tool-use ID,
session, turn, and working-directory context for many local function, shell,
patch, and MCP operations. Hosted or specialized execution paths may not be
covered and require a capability test. Deterministic backends should write
their own normalized run record and captured exit status rather than rely on
transcript interpretation.

#### Operator approval

A Codex permission request is a tool-permission event and can itself be allowed
or denied automatically. It does not prove that a human saw and approved a
workflow contract. Contract approval requires a distinct event bound to the
displayed contract hash.

#### Identity and timestamps

Codex provides session, turn, and working-directory information. It does not
provide a complete canonical repository identity as a single field; the
recorder must resolve and persist repository root, remote identity, branch,
and HEAD. A trusted timestamp is not documented as a universal hook field, so
the run recorder should timestamp and hash events when appending them.

The Codex app-server event and approval surface is a better prospective
integration boundary than parsing an unstable transcript format:

- <https://learn.chatgpt.com/docs/app-server>

### Required outputs and validations

Output files are collected and hashed, but the resulting artifact observations
never change the verdict. No required validation collector or reconciliation
rule exists. Required outputs and validations therefore do not currently block
completion.

### Closure immutability

`Close()` does not reject `PhaseClosed`. Closing the same run again overwrites
the evidence, reconciliation, verdict, and update timestamp. `SaveState()` is
also public and performs no phase-transition or state-hash validation.

Contract and state creation are two non-atomic writes. A failure after writing
the tracked contract leaves a half-created run. Mutable state under ignored
`local_artifacts` survives a process restart on the same filesystem, but is not
durable across a clean checkout, worktree loss, or evidence cleanup.

A partial, failed, or unverifiable reconciliation is persisted as
`PhaseClosed` before the CLI returns its non-zero result.

### Receipt behavior

A human receipt cannot override the deterministic verdict because no close
path accepts a verdict as input. However:

- a missing receipt is explicitly accepted;
- receipt validation occurs after closure state is stored;
- a disagreeing receipt therefore causes an error after mutation;
- re-closing can subsequently replace the recorded reconciliation.

### R14 conclusion

The store and deterministic-verdict direction are worth preserving, but the
current evidence collector is Claude-specific, globally scoped, and unusable
as trustworthy Codex evidence. Closure is mutable and required outputs are not
requirements in practice.

## Preserved import snapshot

The snapshot branch is a valid archive and must not be cherry-picked wholesale.

The diff from `93f7d03` contains:

```text
1 AGENTS.md
35 skill/shared files
13 agent TOML files
1 snapshot manifest
50 added files total
49 imported resources
```

All 35 skill/shared blobs are byte-identical to their corresponding
`.claude/skills/` source files.

### Critical translation defects

- Root `AGENTS.md` contains broken uppercase `.Codex/` references.
- It names slash commands, rules, templates, playbooks, and a README that were
  not imported.
- The omitted tracked source layer comprised 43 commands, 28 playbooks, 4
  rules, 12 templates, and `.claude/README.md`.
- No tracked hooks or configuration existed to provide the missing enforcement;
  ignored local settings were intentionally excluded from the snapshot.
- Every source agent had an explicit Claude `tools:` boundary. All 13 imported
  TOML files omit those boundaries.
- Agent instructions rely on prose saying "Use the skill" instead of explicit
  Codex skill and capability configuration.
- The import predates AUDIT-0002 corrections and this R12/R14 review.

No resource is assigned `drop`. The archive should retain the original content
even where activation is deferred or blocked.

### Root instruction disposition

`AGENTS.md`: **rewrite for Codex**. Path corrections alone are insufficient
because its command surface, invocation model, and enforcement assumptions are
also obsolete.

### Skill and shared-file dispositions

Paths below are relative to `codex/import-snapshot:.agents/skills/`.

| Disposition | Count | Resources |
|---|---:|---|
| **retain** | 8 | `65816-analyst`, `_shared/ADDRESS_SPACES.md`, `_shared/ASSET_PROVENANCE.md`, `_shared/EVIDENCE_STANDARD.md`, `_shared/LEGAL_BOUNDARY.md`, `experiment-designer`, `mesen-operator`, `static-runtime-correlator` |
| **retain with path corrections** | 1 | `census-observer` -- replace the omitted `/census-observations` dependency |
| **rewrite for Codex** | 3 | `_shared/STOPPING_RULES.md`, `context-manager`, `research-orchestrator` |
| **merge** | 7 | `asset-validator` into verification; `ff6-content-census` into the smaller census skill; `function-recovery`, `struct-recovery`, and `variable-recovery` into static/runtime analysis; `go-architect` and `go-implementer` into one Go implementation skill |
| **convert to AGENTS.md guidance** | 1 | `documentation-curator` |
| **convert to deterministic backend** | 2 | `asset-cataloger` mechanical hash/manifest checks; `quality-auditor` mechanical gates |
| **defer** | 9 | `audio-archaeologist`, `background-reconstructor`, `brr-reconstructor`, `contradiction-resolver`, `graphics-archaeologist`, `palette-researcher`, `release-manager`, `sprite-reconstructor`, `tileset-reconstructor` |
| **block because backend is absent** | 3 | `dsp-validator`, `sequence-reconstructor`, `spc700-analyst` |
| **unknown pending behavioral test** | 1 | `dma-tracer` -- the live probe exists but is explicitly unexercised |

### Agent dispositions

Paths below are relative to `codex/import-snapshot:.codex/agents/`.

| Disposition | Count | Agents |
|---|---:|---|
| **rewrite for Codex** | 6 | `assembly-analyst`, `experiment-planner`, `go-implementation-engineer`, `research-manager`, `static-analysis-researcher`, `verification-engineer` |
| **merge** | 2 | `documentation-reviewer` and `quality-reviewer` into the verification/review responsibility |
| **defer** | 4 | `asset-librarian`, `audio-researcher`, `graphics-researcher`, `release-engineer` |
| **unknown pending behavioral test** | 1 | `dma-researcher` |

Every retained or rewritten agent needs explicit sandbox, tool, MCP, and skill
boundaries before activation.

## Ghidra and Mesen coverage

Ghidra and Mesen are already first-class tracked project components:

- three Ghidra/static-analysis documents;
- `scripts/check-ghidra-setup.sh`;
- one full static/runtime correlation record, `CORR-0001`;
- 24 tracked Mesen Lua files: the bridge plus 23 probes/helpers;
- 52 experiment records plus the experiment README;
- a detailed Mesen capability matrix;
- offline savestate tooling described by the project;
- imported 65816, experiment, Mesen, DMA, function-recovery, and
  static/runtime-correlation guidance.

The intended evidence loop appears in both documentation and canonical
evidence:

```text
Mesen runtime observation
-> Ghidra static analysis
-> bounded hypothesis
-> Mesen verification or refutation
-> canonical repository evidence
-> Go reconstruction
```

Remaining integration issues:

- Static-analysis documents still invoke Claude slash commands without Codex
  equivalents.
- The Ghidra setup script checks paths, Java, and ROM hash but leaves loader,
  language, and mapped-address verification manual.
- `mesen/probes/dma-trace.lua` is explicitly unexercised.
- Several Mesen GUI capabilities remain unknown.
- `CORR-0001` demonstrates the evidence loop but records that no separate
  experiment document was created before operating Mesen. That disclosed
  exception should be prevented by the future workflow contract.
- Exact emulator identity must be recorded per run; historical records refer
  to different Mesen builds.

### Later narrowly scoped access requirements

A later, explicitly authorized workflow would need:

- read/execute access to the Ghidra installation under `../tools/ghidra/`;
- read access to loader source under `../tools/ghidra-snes-source/`;
- read/write access only to the external Ghidra project under
  `../workspaces/ghidra/`;
- read-only access to the single verified ROM under `../private/roms/`;
- permission to execute the exact Mesen build;
- write access only to ignored `mesen/out/` and `local_artifacts/` for raw
  evidence;
- bounded writes to specifically approved canonical repository records and Go
  files;
- no permission to copy ROM bytes, Ghidra databases, savestates, captures, or
  extracted commercial assets into tracked paths.

None of those external local components was inspected or operated during this
review.

## Minimum trustworthy resumption foundation

The smallest sufficient foundation is:

1. A concise rewritten root `AGENTS.md` containing repository authority, legal
   boundaries, evidence precedence, bounded-work rules, and current Codex paths.
2. The five shared contracts: address spaces, evidence, legal boundary,
   provenance, and stopping rules.
3. Corrected R12/R14 behavior providing:
   - evidenced operator approval;
   - immutable freezing and amendment;
   - run-scoped events;
   - enforceable conditional and failure-policy semantics;
   - required-output and validation reconciliation;
   - immutable closure;
   - durable continuation state.
4. One rewritten research-orchestration skill plus `mesen-operator`,
   `experiment-designer`, `65816-analyst`, and `static-runtime-correlator`.
5. Six bounded agents: research manager, experiment planner, assembly analyst,
   static-analysis researcher, verifier, and Go implementer, each with explicit
   capability boundaries.
6. A `/research` equivalent that displays and obtains approval for a bounded
   contract, and a `/continue` equivalent that reloads the open run. These do
   not need to be slash commands until the Codex mechanism is tested.
7. A deterministic quality-gate backend whose run-bound exit records and
   required outputs directly affect the verdict.

Everything else may be deferred until a selected research unit needs it.

The intended trustworthy lifecycle is:

```text
operator requests bounded research
-> contract is displayed and explicitly approved
-> required specialists are explicitly invoked
-> Mesen and/or Ghidra work is selected appropriately
-> run-scoped evidence is recorded
-> required outputs and validation are enforced
-> missing required work blocks completion
-> immutable state supports continuation in a later session
```

## Migration readiness

**ready after named corrections**

The named prerequisites are:

- prove approval rather than trusting `state: approved`;
- close every freeze and amendment bypass;
- enforce outputs, validations, conditional requirements, and failure policies;
- replace Claude transcript scanning with Codex-compatible run events;
- bind evidence to run, contract, session, repository, directory, branch, time,
  and event identity;
- make closure immutable and validate receipts before committing terminal
  state;
- rewrite the root instructions and essential agent capability boundaries.

The full import snapshot is not approved for wholesale application.

## Proposed next bounded action

Implement one R14 unit: **replace global Claude transcript discovery with an
append-only, hashed run-event ledger bound to workflow ID, contract hash,
immutable run ID, session/turn, repository, working directory, branch/HEAD,
timestamp, and event ID; add negative tests proving malformed, historical,
cross-run, and cross-repository events cannot satisfy a requirement.**

This action was proposed by the review. It was not performed during the
read-only review.

## Review validation status

No tests or builds were run during the compatibility review because it was
explicitly read-only. Repository documentation and implementation were
inspected without modifying files, changing branches, applying the migration
snapshot, operating Mesen or Ghidra, inspecting the private ROM, or resuming
FF6 reconstruction.
