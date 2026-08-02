# AUDIT-0001 — Orchestration inventory and reference graph

**Audit date:** 2026-08-02 · **Branch:** `maintenance/workflow-observability`
**Repository HEAD at audit start:** `581ddbc`
**Machine-readable baseline:** `AUDIT-0001-baseline.json` (all 135 resources)

This is a workflow audit. No game research, no Mesen, no Ghidra, no
remediation. Every finding below is a proposal.

## Scope

Every orchestration resource under `.claude/`, plus the documentation entry
points, supporting Go packages, probes, scripts, schemas and manifests that
those resources claim to stand on.

## Resource counts — 135 total

| Category | Path | Count |
|---|---|---|
| Commands | `.claude/commands/*.md` | 43 |
| Specialist skills | `.claude/skills/*/SKILL.md` | 30 |
| Shared contracts | `.claude/skills/_shared/*.md` | 5 |
| Agents | `.claude/agents/*.md` | 13 |
| Playbooks | `.claude/playbooks/*.md` | 28 |
| Rules | `.claude/rules/*.md` | 4 |
| Templates | `.claude/templates/*.md` | 12 |
| **Total** | | **135** |

All 135 are inventoried, including every orphan. Orphans are the audit's
primary subjects, not exclusions.

## Reference graph method

The graph matches three spellings per resource — full path, bare name, and
type-specific forms (`/command`, `NAME.md`, `playbooks/NAME`, `_shared/NAME`) —
across every `.md`, `.go`, `.lua`, `.json` and `.sh` file in `.claude/`,
`docs/`, `dashboards/`, `indexes/`, `manifests/`, `schemas/`, `internal/`,
`cmd/`, `mesen/` and `scripts/`.

A single-spelling search is not sufficient, and one was corrected during this
audit: an initial path-only pass reported 8 unreferenced templates. The
generated graph resolves the true figure to **9**.

## Broken references — zero

No resource under `.claude/` contains a Markdown link or backtick path that
fails to resolve. `ff6lab audit`'s `CheckMarkdownLinks` was also clean at the
pre-audit baseline.

## Orphans — 31 of 135 have zero inbound references

| Type | Orphaned | Of | Which |
|---|---|---|---|
| Command | 6 | 43 | `/recover-background`, `/recover-compression`, `/recover-event-opcode`, `/recover-sprite`, `/recover-text`, `/recover-tileset` |
| Skill | 0 | 30 | — |
| Shared contract | 1 | 5 | `ASSET_PROVENANCE` |
| **Agent** | **12** | **13** | all but `dma-researcher` |
| Playbook | 0 | 28 | — |
| Rule | 3 | 4 | `assets`, `dashboards`, `go` |
| Template | 9 | 12 | `AUDIO_ASSET`, `CONTRADICTION`, `DISCOVERY`, `FUNCTION`, `GRAPHICS_ASSET`, `SESSION`, `STRUCT`, `VALIDATION_REPORT`, `VARIABLE` |

Orphan status is a routing fact, not a quality judgement. The six orphaned
commands include the **best-written** commands in the set: `/recover-compression`
opens by warning the operator not to assume compression exists, citing EXP-0050;
`/recover-text` distinguishes the Confirmed menu encoding from the unstudied
dialogue stream; `/recover-event-opcode` resumes precisely from CORR-0001's
unresolved predecessor. They are orphaned because they were created on
2026-08-02 in commit `7969d50` and nothing has referenced them since.

## Finding A — Agent routing: reachable, but not routed

**12 of 13 agents have zero inbound references from any file in the
repository.** The thirteenth, `dma-researcher`, appears only in two checkpoint
narratives and one Go source comment — description, never routing.

No command, skill, playbook, rule or template names any agent.
`ORCHESTRATE_RESEARCH.md` step 6 is the single word **"Delegate."**

This must be stated as two separate facts, because they are separate:

- **Reachability is CONFIRMED.** All 13 agents were invoked successfully in
  Phase 4. None failed to load.
- **Explicit project routing is CONFIRMED ABSENT.** Every one of those 13
  invocations worked only because the orchestrator supplied the agent name
  directly. Selection relies on semantically matching each agent's one-line
  `description` field at invocation time.

`.claude/README.md` states the project's own rule: *"Do not depend on
[automatic skill invocation] alone for critical procedures. Each command names
the skills and playbook it requires."* That discipline was applied to skills
and playbooks — 0 of 30 skills and 0 of 28 playbooks are orphaned — and never
extended to agents.

Three independent agents converged on this. `quality-reviewer`, given free
choice of target, selected it unprompted as the most important defect, and
explicitly dismissed the command aliases as "self-documented and benign."

## Finding B — Documentation drift

| Document | Commands listed | Missing | Phantom entries |
|---|---|---|---|
| `.claude/README.md` | 28 | **15** | 0 |
| `docs/WORKFLOW_COMMANDS.md` | 30 | **13** | 0 |

Missing from `.claude/README.md`: `/battle-baseline`, `/bootstrap-ghidra`,
`/capture-frame`, `/census-observations`, `/correlate-static-runtime`,
`/export-ghidra-symbols`, `/recover-background`, `/recover-compression`,
`/recover-event-opcode`, `/recover-map`, `/recover-sprite`, `/recover-text`,
`/recover-tileset`, `/register-system`, `/update-coverage`.

Missing from `docs/WORKFLOW_COMMANDS.md`: the same list minus the three Ghidra
commands and `/recover-background`, plus **`/session-summary`** — which the
audit brief did not anticipate.

**The two documents also disagree with each other.** `.claude/README.md` omits
the Ghidra trio that `docs/WORKFLOW_COMMANDS.md` documents;
`docs/WORKFLOW_COMMANDS.md` omits `/session-summary` that `.claude/README.md`
documents. Twelve commands are missing from both.

Neither document names a command that does not exist.

`.claude/README.md` also documents a **"Workflows"** component category —
"Multi-stage lifecycle procedures that cross specialists" — for which no
directory exists.

## Finding C — Aliases are deliberate and correct

`/recover-background`, `/recover-sprite` and `/recover-tileset` each open by
declaring themselves aliases and explaining the reason: the playbooks are named
`RECOVER_*`, so both spellings resolve rather than one silently failing. They
carry their canonical command's full skill list.

**These are not accidental duplication and must not be counted as such.** They
inherit their canonical command's implementation status in the baseline.

## Finding D — Census overlap

`ff6-content-census` is the only skill with no inbound reference from any
command. `/census-observations`, `/register-system` and `/update-coverage` all
route to `census-observer`.

`ff6-content-census` is reachable from the constitution — `CLAUDE.md` names it
directly — but not from the command layer. Two skills carry the same operating
rule ("Observe broadly. Register briefly. Investigate narrowly.") with no
declared routing distinction, and three commands share one of them with no
declared distinction between the three.

## Finding E — Templates: 9 of 12 unreferenced

Only `CHECKPOINT`, `EXPERIMENT` and `STATIC_CORRELATION` are referenced.

The mechanism is visible by contrast:

- `/checkpoint` says "Use the context-manager skill and
  `.claude/templates/CHECKPOINT.md`" — the template is reachable.
- `/session-summary` says "Use the documentation-curator skill" and **never
  names `SESSION.md`** — the template is orphaned.

A template exists for `/session-summary` and the command cannot reach it. The
same pattern explains `DISCOVERY`, `FUNCTION`, `VARIABLE`, `STRUCT`,
`CONTRADICTION`, `VALIDATION_REPORT`, `AUDIO_ASSET` and `GRAPHICS_ASSET`.

## Finding F — Twelve commands reference no playbook

`/battle-baseline`, `/bootstrap-ghidra`, `/bootstrap-v4`,
`/census-observations`, `/checkpoint`, `/correlate-static-runtime`,
`/export-ghidra-symbols`, `/register-system`, `/resume-session`,
`/session-summary`, `/trace-caller`, `/update-coverage`.

Several are legitimately playbook-free (`/checkpoint` is a template-driven
procedure). `/trace-caller` is the weakest: no playbook, no template, no named
output location.

## Finding G — Unit IDs collide globally

Not anticipated by the brief, and it directly affects the telemetry design.

`dashboards/ACTIVITY_LOG.md` numbers units globally, reaching **Unit 35** by
2026-08-01 and Unit 49 in the EXP-0047 entry. Session 005 then restarted at
**Unit 10** on 2026-08-02 and ran to Unit 18.

**Unit 12, Unit 17 and Unit 18 each name two unrelated pieces of work.** Any
telemetry keyed on `unit_id` would silently merge them. See the telemetry
design record.

## Finding H — Dashboard staleness

| Dashboard | Last commit | Status |
|---|---|---|
| `ACTIVITY_LOG.md` | 2026-08-01 | **Stale.** Newest entry is 2026-08-01; Session 005 ran 16 commits, four experiments and nine units on 08-02, none of them logged. |
| `OPEN_HYPOTHESES.md` | 2026-07-31 | **Stale.** Holds only `H-BATTLE-*`. The tactical-pause checkpoint carries three live hypotheses (33-byte table as entrance/warp records; `$39`/`$17` as map ids; 128-byte runs as animated tiles) and three live alternatives, none registered. |
| `COVERAGE.md` | 2026-08-02 | Generated, current, but carries no as-of date. |
| Others | 2026-08-01/02 | Current. |

The `OPEN_HYPOTHESES.md` gap is the more serious: hypotheses live in
checkpoint prose while the dashboard built to track them is not updated, so the
project's live hypothesis set is only recoverable by reading checkpoints.

## Finding I — Constitution and schema disagree on experiment fields

`CLAUDE.md` requires every experiment record to separate seven things:
Question, Starting state, Controlled variables, Expected outcomes, Falsifier,
Required evidence, Stopping condition.

`schemas/experiment.schema.json` defines exactly seven properties and requires
all seven — but they are `schema_version`, `id`, `question`, `starting_state`,
`expected_outcomes`, `falsifying_outcome`, `status`.

**Three constitutional fields have no schema property at all:** controlled
variables, required evidence, stopping condition. They survive only in the
Markdown records, outside machine validation.

Separately, `domain` appears in 12 of 52 manifest entries and is not defined in
the schema — schema and manifest have already drifted from each other.

## Finding J — A provenance defect the tooling does not catch

`manifests/assets.json` → `ASSET-GFX-0002` ("Narshe field BG tileset, first
block", `ROMFILE:0x208460-0x20A45F`, runtime consumer "SCN-0001 milestones 02
and 04") links `census_refs: ["CEN-GFX-0006"]`.

`CEN-GFX-0006` is **"Mines-interior PPU configuration and VRAM image"**, whose
`related_experiments` is `EXP-0035`.

A Narshe asset points at a mines census record. `ff6lab audit` passes, because
it validates that the reference resolves, not that it is the right target.

**Not repaired.** Research records are outside this audit's write scope.

## Finding K — `LATEST.md` duplicates the checkpoint body verbatim

`docs/checkpoints/LATEST.md` is not a pointer. It carries a two-line header
followed by a **byte-identical copy** of the linked checkpoint's entire body —
182 lines at the time of this audit.

This is the "records duplicated across files" class the audit brief asks about.
Two copies of the same content diverge the moment either is edited, and the
duplication is not machine-checked: `ff6lab audit` verifies that the link
resolves, not that the bodies agree.

The convention was **followed rather than corrected** at this session's
closure, because changing it would be remediation. A check asserting that
`LATEST.md`'s body matches its linked checkpoint — the shape
`CheckReadinessSummary` already uses to assert a summary against the rows it
summarizes — would make the duplication safe, or the duplication could be
removed in favour of a genuine pointer. Either is a decision for the operator.

## Contradiction recorded, not resolved

`CLAUDE.md` requires that a completed unit synchronize all affected canonical
and generated state, dashboards included. The AUDIT-0001 brief forbids
modifying dashboards.

Resolution applied at closure: generated/required synchronization only —
`docs/checkpoints/LATEST.md`, `indexes/SESSIONS.md`, and the single
session-count row in `dashboards/STATISTICS.md`, which SESSION_006 makes
mechanically stale. No substantive dashboard change; no research fact changed
this session.

**This contradiction is preserved rather than silently resolved.** A future
governance decision should state which authority wins for audit-class sessions.
