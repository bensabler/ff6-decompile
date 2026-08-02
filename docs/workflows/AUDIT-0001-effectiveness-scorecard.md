# AUDIT-0001 — Effectiveness scorecard

**Audit date:** 2026-08-02 · Full per-resource scores in `AUDIT-0001-baseline.json`

## Method, and its limits

Each scored resource carries 0–5 values on the dimensions the brief specifies
for its type. **Inapplicable dimensions are `null`, never 0**, and are excluded
from the composite.

```text
composite = arithmetic mean of applicable (non-null) dimensions
```

Every composite in the baseline publishes its formula, its
`applicable_dimension_count`, its `total_dimension_count`, and the list of
`null_dimensions` that were excluded. 114 of 135 resources are scored; rules,
templates and shared contracts have no dimension list in the brief and carry
`not_applicable` rather than a fabricated score.

### Where the composite is informative — and where it is not

**This is the scorecard's most important disclosure.**

| Type | n | Range | Informative? |
|---|---|---|---|
| Command | 43 | 2.90 – 4.10 | **Yes.** Real spread driven by verifiable inputs: documentation coverage, backend status, playbook wiring, usage class. |
| Agent | 13 | 3.38 – 3.86 | **Partly.** Spread comes almost entirely from live smoke-test results, which is the only hard agent evidence that exists. |
| Skill | 30 | 2.00 – 4.00 | **No.** 28 of 30 score exactly 4.00. The scorecard cannot currently distinguish them. |
| Playbook | 28 | 4.67 – 4.67 | **No.** All 28 score identically. |

Skills and playbooks are uniform because the available signals are uniform:
every playbook has Required inputs, Procedure and Required outputs sections;
every skill is command-reachable and shared-contract compliant. **The dimensions
that would discriminate them — procedural completeness, implementation truth,
demonstrated downstream value — cannot be scored without invocation
telemetry**, so they are `null`.

Do not use the skill or playbook composites to make decisions. They are
recorded for completeness and to establish the diff baseline; the two
exceptions that *did* separate (`ff6-content-census` 2.00, `census-observer`
3.00) did so on overlap, which is observable without telemetry.

**Score confidence** is `medium` where usage is Confirmed or Probable and `low`
otherwise — which is 101 of 135 resources.

## Commands — the informative ranking

### Highest

| Composite | Command | Why |
|---|---|---|
| 4.10 | `/recover-tileset` | Real tested decoders both bit depths; HUD font ships in `ff6demo` |
| 4.08 | `/bootstrap-v4` | Confirmed use, documented, artifact exists |
| 4.08 | `/resume-session` | Confirmed, documented, most-used session control |
| 4.08 | `/run-quality-gates` | Every gate real; run twice this session |
| 4.00 | `/correlate-static-runtime` | Confirmed via CORR-0001; names its template |

### Lowest

| Composite | Command | Why |
|---|---|---|
| 2.90 | `/validate-audio` | Prescribes DSP comparison; no audio comparison code exists |
| 3.00 | `/trace-spc-command` | No SPC700 tooling |
| 3.00 | `/recover-sequence` | No sequence package |
| 3.10 | `/recover-map` | No map backend — but correctly scoped to the largest open area |
| 3.27 | `/recover-text` | Menu encoding real; dialogue stream unbacked |

**The floor is 2.90.** Nothing is worthless. The weakest commands are weak
because their *subsystems* are unstudied, not because the commands are badly
written — and `/recover-map` in particular is a well-specified command pointing
at the project's biggest genuine gap.

### Highest-value resources

Ranked by leverage rather than composite:

1. **`/battle-baseline`** — the only complete
   command → rule → automated-enforcement chain in the project.
   `CheckBattleExperimentConfig` fails `ff6lab audit` for battle experiments
   without a configuration fingerprint. It converts a discipline into a gate,
   which is what every other remediation in this audit is trying to achieve.
2. **`/run-quality-gates`** — the only command whose full output is verifiable
   on demand, and the audit's own control.
3. **`/checkpoint`** — 61 records; the project's actual memory.
4. **`/census-observations` family** — schema, manifest, validator and sync
   behind them.
5. **`/trace-dma`** — highest *honesty*, and its four-surface disclosure is the
   model the rest of the layer should copy.

### Highest-risk resources

1. **`/audit-project`** — mutation-only, capability honesty 1. The project's
   designated self-assessment tool cannot be run in assessment mode. It is also
   one of only two commands with Confirmed historical use, so it is both
   trusted and unsafe for the purpose this audit needed.
2. **`/validate-graphics`** — genuine zero on downstream integration. The
   validator exists, is tested, and is imported by nothing, while its playbook
   prescribes the failure mode the validator prevents.
3. **`/validate-audio`, `/recover-sequence`, `/trace-spc-command`** — prescribe
   procedures against tooling that does not exist.
4. **The six orphaned recovery commands** — the best-written commands in the
   set, invisible from both documentation entry points.

## Agents — scored on live evidence

All 13 were invoked once. Scores are driven by
`output_contract_compliance` and `integration_value`, both measured this
session. `routing_coverage` is **0 for 12 of 13** — a genuine zero: the
evidence exists and the measured value is zero.

| Composite | Agent | Note |
|---|---|---|
| 3.86 | `assembly-analyst` | Only agent with **full** contract compliance; found a defect in its own contract |
| 3.75 | `quality-reviewer` | Unprompted, independently selected agent routing as the top defect and correctly dismissed the aliases as benign |
| 3.75 | `release-engineer` | Corrected the auditor by finding CI enforcement the playbooks hide |
| 3.75 | `verification-engineer` | Found the unimported validator — the audit's sharpest finding |
| 3.38 | `audio-researcher` | Real finding, worst contract compliance |

**Contract compliance across 13 agents: 1 fully compliant, 11 partial
(rationale over the 150-word cap), 1 violated** (`audio-researcher` appended a
second structured result). Every agent respected scope: read-only, no Mesen, no
Ghidra, no experiments, and `git status` was clean after every batch — including
for the four agents holding Write/Edit/Bash authority.

**Three agents made errors the auditor caught:** `documentation-reviewer`
counted 42 command files (43 exist); `go-implementation-engineer` reported a
broken link to `docs/indexes/` when the reference is to `indexes/` at the
repository root; `graphics-researcher` reported no backend for
`/reconstruct-tileset` while `tile2bpp`/`tile4bpp` are real and tested.
Agent output required verification in 3 of 13 cases — a useful base rate.

## Recommendations — full distribution

| Type | Keep | Repair | Clarify Routing | Validate Backend |
|---|---|---|---|---|
| Command (43) | 23 | 14 | 1 | 5 |
| Skill (30) | 28 | — | 2 | — |
| Agent (13) | — | — | **13** | — |
| Playbook (28) | 28 | — | — | — |
| Rule (4) | 1 | — | 3 | — |
| Template (12) | 3 | — | 9 | — |
| Shared contract (5) | 4 | — | 1 | — |
| **Total** | **87** | **14** | **29** | **5** |

**Zero Removal Candidates. Zero Deprecation Candidates. Zero Merge
Candidates.**

No resource became a removal candidate on absent telemetry alone, and none
qualified on other grounds either. The 14 `Repair` recommendations are
documentation-coverage fixes; the 29 `Clarify Routing` are the agents, the
orphaned templates, the three orphaned rules, `ASSET_PROVENANCE`, the census
pair, and `/capture-graphics`.

`/capture-graphics` was the one plausible Merge Candidate — it overlaps
`/capture-frame` on evidence set. It is scored `Clarify Routing` instead
because `CAPTURE_GRAPHICS.md` adds layers and viewer captures that
`CAPTURE_FRAME.md` does not, so the distinction may be intentional but is
stated nowhere. Deciding that is a scoping call for the operator, not a
conclusion this audit can reach from evidence.
