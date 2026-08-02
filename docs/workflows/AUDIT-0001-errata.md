# AUDIT-0001 — errata

Corrections established by AUDIT-0002. **AUDIT-0001's six records and its
baseline are preserved unmodified** as historical evidence. Nothing there was
rewritten; this file records what changed and why.

Claim-level detail: `AUDIT-0002-claim-ledger.json` (26 claims — 13 refuted,
7 confirmed, 4 partially supported, 2 unverifiable).

## Headline number: 31 orphans → 7

AUDIT-0001 reported **31 of 135 resources with zero inbound references**.

Replayed against its true input tree (`581ddbc`), the generator reproduces that
figure **byte-identically** — it is deterministic and was not miscomputed.

But the matcher's corpus walks only ten subdirectories and **never scans
repository-root files**. `PACKAGE_MANIFEST.json` sits at the root and lists
**24 of the 31 by full path**.

| Metric | Value |
|---|---|
| AUDIT-0001 reported | 31 |
| **Corrected textual orphans** | **7** — 6 recovery commands + `static-analysis-researcher` |
| **Routing-bearing orphans** | **38** — a different and more useful question |

Textual inbound was never the right proxy. A packaging manifest listing a file
is not something *routing* to it. Measured on inbound from a resource that
would actually use it — a command, skill, playbook, rule, template, or the
constitution — the figure is 38, and it points the opposite way from the
textual correction.

## Invocation classifications — wrong in both directions

Harness session transcripts record every tool call with a timestamp. Four exist,
**all spanning 2026-08-02 only**; nothing about sessions 001-004 is verifiable.
Five slash-command invocations exist in all preserved history.

| Command | AUDIT-0001 | Corrected | Basis |
|---|---|---|---|
| `/run-quality-gates` | **Confirmed** | **Refuted** | No `Skill` call in any session, including AUDIT-0001's in full. Shell commands were run instead |
| `/correlate-static-runtime` | **Confirmed** | **Unverifiable** | Absent from transcripts; CORR-0001 predates coverage |
| `/bootstrap-ghidra` | **Confirmed** (artifact exists) | **Confirmed** | Genuinely invoked 00:50:05Z. **Correct answer, invalid method** — the artifact-existence reasoning stays condemned |
| `/checkpoint` | Probable | **Confirmed** | Invoked twice (01:25:54, 14:32:14) |
| `/session-summary` | Probable | **Confirmed** | Invoked once (14:33:45) |
| 13 agents, live use | omitted from totals | **Confirmed** | 13 `Agent` calls, 13 distinct subagent types |
| `documentation-reviewer` historical | Unknown | **Confirmed** | Invoked in session `46404e8c`, before AUDIT-0001 |

AUDIT-0001 both over-claimed and under-claimed. A single systematic bias would
be easier to correct; the real pattern is that its evidence rule was applied
inconsistently.

## Scores and recommendations — withdrawn

**Every composite numeric score is withdrawn.** `build_baseline.py` assigns
`scope_clarity=4`, `output_contract=4`, `capability_honesty=5` and
`maintenance_burden=4` to every command regardless of evidence; derives
`stop_conditions` from backend existence; and derives `artifact_production`
from usage class. These are hard-coded defaults and invalid proxies, not
measurements.

**All 28 playbook `Keep` recommendations are withdrawn.** They came from an
unconditional `if t == "playbook": return "Keep"`. No playbook was assessed.

The scorecard's own disclosure that skill and playbook composites "must not
drive decisions" was correct but understated: the command and agent composites
rest on the same hard-coded inputs.

## Rules — category error

AUDIT-0001 reported **3 of 4 rules orphaned** and recommended `Clarify Routing`.

All four rules carry `paths:` frontmatter (`**/*.go`, `dashboards/**/*.md`,
`manifests/**/*.json`, `docs/research/**/*.md`). They are **auto-loaded by glob
scope and never invoked by name.** Being unnamed is their design.

**Recommendation withdrawn.** The matcher cannot observe path-scope activation
at all.

## Templates — real defect, wrong metric

"9 of 12 unreferenced" is false as stated: all nine appear in
`PACKAGE_MANIFEST.json`. **10 of 12 have no consumer reference**, which is the
defect AUDIT-0001 meant. The mechanism it identified — `/checkpoint` names its
template and `/session-summary` does not — survives.

## Gate capture — understated

The brief alleged missing exit capture for two gates. **4 of 12 gates had
reliable capture at Phase 0, and 4 of 12 at closure.** At closure the three
build gates were collapsed into a block printing `ff6lab ok / ff6demo ok /
gui ok` with no status at all.

## Broken references — over-generalised

"Zero broken references" measured two categories: markdown links and backtick
paths. It says nothing about missing named commands, skills, agents or
playbooks, or invalid declared output paths. A zero in two categories was
reported as a zero overall.

## Telemetry design — one mechanism refuted

AUDIT-0001's telemetry "mechanism 1" proposed deriving **Confirmed** command
invocations from artifact appearance. AUDIT-0001's own record is the
counterexample: the same artifact logic produced one true classification
(`/bootstrap-ghidra`) and one false one (`/correlate-static-runtime`).
**Derived artifacts may produce `Probable` associations only.**

Its claim that hooks "cannot be skipped" was asserted without any capability
test and is `unverifiable`.

## Findings that survive unchanged

- `internal/validate` has **zero importers**, non-test and test alike, while
  `VALIDATE_GRAPHICS.md:16` prescribes indexed comparison — re-verified.
- `/trace-dma`'s four surfaces and the probe's self-declared `UNEXERCISED`
  status — re-verified.
- `/audit-project` has no report-only mode — re-verified.
- `internal/audio/doc.go` claims SPC700, DSP and sequence code; only `brr`
  exists — re-verified.
- `ff6lab state origin` and `state sprites` implemented, absent from help —
  re-verified.
- Unit IDs collide globally — re-verified.
- **Explicit agent routing is absent** — and now measured correctly rather than
  by proxy: all 13 agents have zero routing-bearing inbound.
- AUDIT-0001's closure allowlist held: 12 files, no deviation.

## New findings AUDIT-0001 could not have seen

1. **AUDIT-0001 destroyed the metric it reported.** Replayed against the
   closure tree the matcher reports **zero** orphans, because
   `AUDIT-0001-baseline.json` names all 135 resources. Any future run of the
   same method returns 0 regardless of reality.
2. **`/checkpoint` cannot be invoked inside a bounded audit.** It declares it
   updates `dashboards/CURRENT_FOCUS.md` — a substantive dashboard write.
   This structurally explains why AUDIT-0001 wrote its closure records by hand.
3. **`/session-summary` declares no output path at all.**
4. **Automatic specialist selection does not exist.** Omitting `subagent_type`
   yields `general-purpose`, never a project agent.
5. **AUDIT-0001 wrote no `hashes.sha256`** of its own evidence.
