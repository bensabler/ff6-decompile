# AUDIT-0002 — telemetry and registry revision (Phase 8)

Corrects AUDIT-0001's telemetry design. **Design only; nothing built.**

## Two systems, kept separate

**Static workflow inventory and validation** answers "what exists and what
reaches it." **Event telemetry** answers "what was invoked." AUDIT-0001 merged
them, which is how a derived artifact became invocation proof.

## Retraction: derived events cannot confirm invocation

AUDIT-0001 proposed inferring **Confirmed** invocations from artifact
appearance, calling it the strongest mechanism because "nothing has to
remember it."

**AUDIT-0001's own record refutes it.** The same artifact-existence logic
produced one true classification (`/bootstrap-ghidra`, genuinely invoked) and
one false one (`/correlate-static-runtime`, never invoked). The method cannot
distinguish them, and AUDIT-0001 could not tell which of its own two
conclusions was sound.

Every event therefore carries:

```json
{
  "observation_source": "direct_log|hook|explicit_record|derived_artifact|reconstructed",
  "confidence": "Confirmed|Probable|Possible",
  "evidence_refs": []
}
```

**`derived_artifact` may never yield `confidence: Confirmed`.** A checkpoint
file appearing is `derived_artifact` → `Probable`.

## What actually works: the harness transcript

The strongest available invocation evidence is the per-session JSONL transcript
the harness already writes. It records every tool call with tool name, input
selector and ISO-8601 timestamp, is not authored by the agent being measured,
and requires no cooperation.

This is what let AUDIT-0002 refute four AUDIT-0001 classifications. It should
be the **primary** telemetry source, with `observation_source: direct_log`.

**Its limits, which the design must carry:**

- **Coverage is partial and not retroactive.** Four transcripts exist, all
  from 2026-08-02. Sessions 001-004 have none.
- **It is mutable and outside the repository.** Extracts must be frozen and
  hashed at capture time; custody begins then.
- **Privacy is non-negotiable.** Only tool name, timestamp and selector may be
  extracted. Never prompts, responses or conversation content.
- **It cannot attribute generated files** to a command.

## Hooks: untested, not reliable

AUDIT-0001 asserted hooks "cannot be skipped." No capability test was ever
run. **Classification: `available but untested`.** Before any design depends on
them, determine supported types, configuration location, project-local
authority, behaviour on normal stop / cancellation / crash / context
exhaustion, duplicate handling, whether they may run `ff6lab`, and whether
failure blocks or warns.

## Registry authority

**One authority model: structured resource front matter generates the
registry.** Rules already carry `paths:` frontmatter and skills carry
`name:`/`description:`, so the pattern exists. A registry generated from the
resources themselves cannot drift from them.

**Do not seed the registry from AUDIT-0001's baseline.** Use
`AUDIT-0002-corrected-baseline.json`.

## Corpus separation — the self-contamination fix

`AUDIT-0001-baseline.json` enumerates all 135 resource ids, so replaying the
orphan analysis against the closure tree returns **zero** orphans. The audit's
own output destroyed the metric.

**Any inventory validator must exclude audit and report outputs from the corpus
it scans**, by explicit path exclusion (`docs/workflows/`, generated indexes)
rather than by heuristic. Without this, orphan detection silently always
returns 0.

**And it must scan repository-root files**, which AUDIT-0001's corpus omitted —
the reason `PACKAGE_MANIFEST.json` was missed and the orphan count was wrong by
a factor of four.

## Metric definitions must be explicit

AUDIT-0001 reported one "orphan" number that silently answered a different
question from the one it implied. The registry must define and report
separately:

| Metric | Meaning |
|---|---|
| `textual_inbound` | any file mentions it — includes packaging manifests and narrative |
| `routing_bearing_inbound` | a command, skill, playbook, rule, template or the constitution names it as something to use |
| `auto_activated` | applies by path scope; naming is not the mechanism |

A rule with zero routing-bearing inbound is **correct**, not orphaned.

## Identity

`session_id` + `unit_id` + `unit_key` (`SESSION-006/UNIT-001`). Unit numbers
collide globally — `ACTIVITY_LOG.md` reached Unit 35 while Session 005
restarted at Unit 10, so Units 12, 17 and 18 each name two unrelated pieces of
work. **`unit_key` is the only safe key.**

Plus `schema_version`, `event_id`, `invocation_id`,
`parent_invocation_id`, `canonical_resource_id`, `alias_invoked_as`,
`event_type`, timestamp, monotonic sequence, status, outputs, commit, and
handling for start/finish pairing, abandoned invocations, crash recovery,
duplicate events, idempotency, file locking, concurrent writers, session
continuation, partial commits and clock skew.

## Build order

1. **`ff6lab workflow inventory` / `validate`** — static only, no telemetry.
   Needs no behaviour change and would have caught the root-file omission, the
   corpus contamination, and the missing consumer references.
2. **Transcript extraction** — `direct_log` events from harness transcripts,
   frozen and hashed, metadata only.
3. **Hooks** — only after capability testing.
4. **Lifecycle contract** — last, and never as the primary mechanism.

AUDIT-0001 ordered derived events first. That ordering is withdrawn.
