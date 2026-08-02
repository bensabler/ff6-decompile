# AUDIT-0001 — Historical usage baseline

**Audit date:** 2026-08-02 · Per-resource detail in `AUDIT-0001-baseline.json`

## The rule this record enforces

**Absence of evidence is `Unknown`, never zero.**

The project has never instrumented resource invocation. Everything below is
reconstruction from git history, commit messages, checkpoints, session records,
experiment records, `dashboards/ACTIVITY_LOG.md`, generated artifacts, probe
headers and documentation provenance. It is the best-supported reading of the
evidence, and it is incomplete by construction. That is the argument for the
telemetry design.

### Classification rules applied

| Class | Requires |
|---|---|
| **Confirmed** | An outcome-bearing record attributes work to the resource **by name**, or AUDIT-0001 invoked it under direct observation. |
| **Probable** | Its artifact class exists in quantity and the workflow demonstrably ran, but no record names it as the agent. |
| **Possible** | Mentioned only in planning or next-action context; no artifact is attributable to it. |
| **Unknown** | No evidence either way. **Not zero use.** |

Also applied: existence is not use; a documentation reference is not
invocation; matching output format is not proof; new commands cannot be
credited to sessions predating their creation commit; aliases normalize to
their canonical command while preserving the invoked spelling.

## Result — 135 resources

| Class | Count | Share |
|---|---|---|
| Confirmed | 6 | 4% |
| Probable | 16 | 12% |
| Possible | 12 | 9% |
| **Unknown** | **101** | **75%** |

Three quarters of the orchestration layer has no recoverable usage evidence.
This is a **measurement failure, not a usage failure** — the project simply
never recorded which resource did which work.

## Raw textual mention counts are not usage

Occurrences of each command's literal `/name` in `docs/sessions`,
`docs/checkpoints`, `docs/experiments`, `docs/discoveries`, `docs/correlations`
and `dashboards`:

`/checkpoint` 25 · `/orchestrate` 10 · `/battle-baseline` 9 ·
`/resume-session` 8 · `/audit-project` 4 · `/run-quality-gates` 3 ·
**24 commands at zero.**

**These are raw textual mention counts and must never be reported as
invocation counts.** A targeted search for outcome-bearing, past-tense
attributions ("ran", "invoked", "executed", "`/command`: <result>") across all
sessions, checkpoints and dashboards returns exactly **three** commands.

The overwhelming majority of `/command` mentions sit in **"Exact next action"**
and **"Recommended next command"** sections. Those are *plans*. The tactical-pause
checkpoint's "Recommended next command: `/audit-project`, then `/orchestrate`"
contributes one mention each to two commands and records zero invocations.

Conflating the two would have inflated `/orchestrate` from Probable to
Confirmed on the strength of recommendations that were never acted on within
the record.

## The six Confirmed

| Resource | Basis |
|---|---|
| `/audit-project` | `ACTIVITY_LOG.md` 2026-07-29 attributes 23 broken-link fixes, a `bridge.lua` path repair and `.gitignore` changes to it **by name**. |
| `/bootstrap-v4` | `ACTIVITY_LOG.md` 2026-07-29 attributes environment/ROM identity recording, index initialization and a migration report to it **by name**. |
| `/bootstrap-ghidra` | `local_artifacts/static-analysis/ghidra-environment.md` exists — the exact artifact its step 10 specifies — and the 2026-08-01 checkpoint attributes it. |
| `/correlate-static-runtime` | `docs/correlations/CORR-0001-C09B5C.md` exists in the `STATIC_CORRELATION.md` shape it mandates, with `mesen/probes/CORR-0001.lua` and hashed local evidence. |
| `/resume-session` | Invoked under direct observation, AUDIT-0001 Phase 0. |
| `/run-quality-gates` | Invoked under direct observation, AUDIT-0001 Phase 0 and Phase 10. |

Only **two** of these — `/audit-project` and `/bootstrap-v4` — are confirmed by
the project's own historical records naming the command. Both entries were
written on the same day, 2026-07-29, and the practice was not continued.

## Notable Probable cases

- **`/checkpoint`** — 61 checkpoint records and a maintained `LATEST.md`. The
  most-used workflow in the project by any measure, and still only Probable:
  not one checkpoint says it was produced by `/checkpoint`.
- **`/battle-baseline`** — 9 mentions and a machine gate that battle
  experiments pass, implying the discipline was applied.
- **`/trace-dma`** — offline decoding demonstrably used by EXP-0050; the live
  surface has never run. Split status, not a single one.

## Resources added after the work they describe

Eight commands were created in `7969d50` on 2026-08-02: `/capture-frame`,
`/recover-background`, `/recover-compression`, `/recover-event-opcode`,
`/recover-map`, `/recover-sprite`, `/recover-text`, `/recover-tileset`.

Their content is unusually good precisely because it was written *from*
completed work — `/recover-compression` cites EXP-0050's refutation,
`/recover-event-opcode` resumes from CORR-0001's unresolved predecessor. But
that same fact means **none of the work they describe can be credited to
them.** They are classified `Possible` at best, and six of the eight are
orphaned and undocumented.

The commit that created them is titled "Put implementations under the DMA
doctrine, and **create the missing commands**" — the project noticed the gap
and filled it. What it did not do was add them to either documentation entry
point.

## Bypass patterns

Distinguishing the categories the brief requires:

| Pattern | Category | Evidence |
|---|---|---|
| Experiments EXP-0050–0053 completed with records but **no preserved evidence directories**, caught only at the tactical pause | **Resource incomplete** — no stopping rule required evidence preservation before closing | Tactical-pause checkpoint: four experiments harvested retroactively; instrumentation for EXP-0051/0052 was ad-hoc and "would not have survived" |
| Map-descriptor work done directly rather than through `/recover-map` | **Resource added after the work** | `/recover-map` created 2026-08-02 in `7969d50`, after EXP-0051/0052 |
| Graphics validation done without `internal/validate` | **Resource undiscoverable** — the playbook prescribes indexed comparison and never mentions the resolved comparison | `VALIDATE_GRAPHICS.md:16`; zero importers of `internal/validate` |
| Agent specialists never routed explicitly | **Resource undiscoverable** — nothing names them | 12 of 13 agents orphaned |
| `SESSION_004` never indexed | **Resource incomplete** — no gate checked session-index synchronization | Commit `581ddbc` fixed it and named the gap rather than hiding it |
| `ACTIVITY_LOG.md` not updated for Session 005 | **Resource ignored** — the log exists, is reachable, and was simply not appended to | Newest entry 2026-08-01; 16 commits followed on 08-02 |

Only one of these is a resource being *ignored*. The rest are structural:
incomplete stopping rules, undiscoverable capability, or resources created
after the fact. **That distinction matters for remediation** — writing more
process discipline would not have prevented five of the six.

## Why unit IDs cannot key telemetry

`ACTIVITY_LOG.md` numbers units globally and reached Unit 35 by 2026-08-01
(Unit 49 in the EXP-0047 entry). Session 005 restarted at Unit 10 on
2026-08-02 and ran to Unit 18.

**Unit 12, Unit 17 and Unit 18 each name two unrelated pieces of work.** The
proposed telemetry event shape carries a `unit_id`; without a session
qualifier it would silently merge them. See the telemetry design record.

## What this baseline is for

`AUDIT-0001-baseline.json` is the tracked artifact future sessions diff
against. It is deliberately conservative: 101 `Unknown` classifications is the
honest reading, and every one of them is a resource whose value the project
currently cannot measure.

**No resource in this audit is a removal candidate on absent telemetry
alone.**
