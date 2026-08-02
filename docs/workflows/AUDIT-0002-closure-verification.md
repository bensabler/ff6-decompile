# AUDIT-0002 — Phase 0-1 closure verification

**Status: PARTIAL.** Phases 0 and 1 complete; Phases 2-11 not started.
**Branch:** `maintenance/workflow-observability-audit2`, worktree
`../ff6-audit2`, from verified closure commit
`93f7d03a199f0c396886ec6497e8401cef88ff54`.

AUDIT-0002 implements no remediation. Every statement below is evidence-bounded
per the vocabulary in the approved plan.

## Phase 0a — slash-command feasibility preflight: **PASS**

The harness writes a per-session JSONL transcript recording every `tool_use`
with an ISO-8601 timestamp. A slash command invoked via the `Skill` tool is
recorded as `{"tool":"Skill","input":{"skill":"<name>"}}`; a subagent as
`{"tool":"Agent","input":{"subagent_type":"<name>"}}`.

**Invocation-level evidence exists and is preservable.** It is
harness-generated, not auditor testimony — materially stronger than
`auditor_recollection`, though the source file is `mutable` rather than a git
object, so the frozen extract's custody begins at AUDIT-0002.

Only tool-call **metadata** was extracted — tool name, timestamp, selector. **No
prompts, no responses, no conversation content.**

### Coverage — the decisive limitation

Four transcripts exist. **All span 2026-08-02 only** (00:35 → 16:47). There is
**no transcript for sessions 001-004** (2026-07-29 → 08-01). Any invocation
claim about work before 2026-08-02 remains `historically_unverifiable`.

### Every slash-command invocation in preserved history

```text
2026-08-02T00:50:05.251Z   /bootstrap-ghidra    session 46404e8c
2026-08-02T01:25:54.587Z   /checkpoint          session 46404e8c
2026-08-02T14:32:14.412Z   /checkpoint          session e426272f
2026-08-02T14:33:45.085Z   /session-summary     session e426272f
2026-08-02T15:04:48.114Z   /resume-session      session 9cde5fa3 (AUDIT-0001)
```

Five invocations, four distinct commands, across 999 recorded tool calls.

## Corrections to AUDIT-0001, from transcript evidence

Recorded here; formal ledger verdicts follow in Phase 2.

| # | AUDIT-0001 said | Transcript shows | Correction |
|---|---|---|---|
| 1 | `/run-quality-gates` **Confirmed**, "invoked under direct observation at Phase 0 and Phase 10" | **No `Skill` call for it in any preserved session**, including the full AUDIT-0001 session | **Refuted.** Charge 1b promoted from `auditor_recollection` to transcript-supported, bounded to 2026-08-02 coverage |
| 2 | `/bootstrap-ghidra` **Confirmed** because its artifact exists | **Actually invoked** 2026-08-02T00:50:05Z | **Conclusion correct, method invalid.** Charge 6 stands undiminished — a bad method that happens to reach a true answer is still a bad method |
| 3 | `/checkpoint` **Probable** | **Invoked twice** (00:25:54, 14:32:14) | **Under-classified.** AUDIT-0001 erred toward under-claiming here — the opposite direction from its `/run-quality-gates` error |
| 4 | `/session-summary` **Probable** | **Invoked once** (14:33:45) | **Under-classified** |
| 5 | `/correlate-static-runtime` **Confirmed** because CORR-0001 exists | Absent from all transcripts; CORR-0001 predates coverage | **Unsupported.** Reverts to `unverifiable` |
| 6 | All 13 agents `Unknown` historical use | Session `9cde5fa3` records exactly 13 `Agent` calls with 13 distinct project subagent types | AUDIT-0001 live use **corroborated by harness evidence**, no longer attestation alone |
| 7 | All 13 agents `Unknown` historical use | `documentation-reviewer` invoked in session `46404e8c`, hours **before** AUDIT-0001 | **Contradicted for this agent** — it has pre-audit live use |

AUDIT-0001 erred in **both** directions: it over-claimed
`/run-quality-gates` and `/correlate-static-runtime`, and under-claimed
`/checkpoint` and `/session-summary`. A single systematic bias would have been
easier to correct; the actual pattern is that its evidence rule was applied
inconsistently.

## Phase 0b — state re-observation

| Check | Result | Basis |
|---|---|---|
| Closure commit | `93f7d03a199f…` resolves unambiguously, type `commit` | `live_observation` |
| Branch / HEAD | `maintenance/workflow-observability` at `93f7d03` | `live_observation` |
| `…-audit2` branch | Absent locally and on live remote before creation | `live_observation` |
| Worktrees | One, the primary checkout | `live_observation` |
| Commits after `93f7d03` | None | `live_observation`; "no remediation occurred" adds `operator_attestation` |
| Uncommitted work | None outside ignored audit evidence | `live_observation` |
| Processes | None capable of mutating evidence | `live_observation`, point-in-time |
| Remote heads | Unchanged from the planning observation | `live_observation` — stop condition 6 not triggered |

**No stop condition triggered.** Worktree created with
`git worktree add ../ff6-audit2 -b … 93f7d03`; nothing forced, reset, deleted
or reused.

## Phase 0d — gate baseline

**Two separate evidence streams, as required.**

1. **Direct `/run-quality-gates` invocation** — performed; recorded in the live
   session transcript and re-verified at closure.
2. **Independent shell harness** — `run-gates.sh`, isolated run directory
   `gate-logs/phase-00-20260802T164950Z/`, initialized `status.tsv` and
   `failures`, per-gate NUL-delimited argv, `manifest.json` with run id, phase,
   purpose, HEAD, cwd, timestamps, prerequisites and verdict.

**Result: 11 required gates, all `rc=0`. AGGREGATE PASS, script exit 0.**
`archive-verify` recorded `not_run: FF6_ROM unset` — never inferred as a pass.

### The harness was itself tested

A negative control proved the properties AUDIT-0001's harness lacked, rather
than asserting them:

| Property | Result |
|---|---|
| Real non-zero status captured | `failing` → `rc=42`, not masked |
| `set -e` does not abort before recording | Confirmed — command runs as an `if` condition |
| Optional failure recorded, not escalated | `optfail` → `rc=7`, absent from `failures` |
| Aggregate fails on required failure | `AGGREGATE: FAIL`, exit 1 |
| Argv boundaries preserved | `["bash","-c","exit 0","arg with spaces","second arg"]` round-trips exactly; `$*` collapses them |

## Phase 1 — evidence frozen

19 files copied to `frozen/` and hashed; `shasum -c` verifies **19/19 OK**.
Covers all 12 AUDIT-0001 evidence files, the 4 transcript metadata extracts,
and the `93f7d03` versions of `LATEST.md`, `SESSIONS.md` and `STATISTICS.md`.

Per amendment 3, evidence from the previously mutable ignored AUDIT-0001
artifacts is **now promotable from `artifact_verified` to
`directly_verified`**. It was not promoted before this point.

**AUDIT-0001 wrote no `hashes.sha256` of its own** (`absence_observed`).
Whether that is a governance defect remains `deferred_pending_contract` for
Phase 6.

## Phase 1 — closure allowlist: **no committed deviation**

12 files committed; **every one inside AUDIT-0001's approved allowlist.**
`dashboards/STATISTICS.md` changed exactly two lines — the as-of line and the
session-count row — matching its stated scope.

Transient staging that was reverted is **`not_tested`**: no contemporaneous
evidence preserves it.

## Phase 1 — gate comparison: charge 2 **understated the problem**

The brief alleged missing exit capture for `go test` and `census validate`.
The frozen logs show the failure was far broader.

| Gate | Phase 0 exit | Closure exit | Capture |
|---|---|---|---|
| `gofmt -l .` | 0 | 0 | reliable |
| `go build ./...` | 0 | 0 | reliable |
| `go vet ./...` | 0 | 0 | reliable |
| `go test ./...` | **BLANK** | **BLANK** | unreliable |
| `go build ./cmd/ff6lab` | 0 | **absent** | unreliable |
| `go build ./cmd/ff6demo` | 0 | **absent** | unreliable |
| `go build -tags gui …` | 0 | **absent** | unreliable |
| `ff6lab audit` | 0 | 0 | reliable |
| `ff6lab census validate` | **BLANK** | **absent** | unreliable |
| restricted-file scan | absent | absent | no exit line either run |
| baseline JSON parse | absent | absent | no exit line either run |
| `archive verify` | absent | absent | not run — `FF6_ROM` unset |

**7 of 12 gates captured an exit status at Phase 0; 4 of 12 at closure. Only 4 gates support a reliable Phase-0-to-closure comparison** — the intersection. (Corrected in Phase 11: an earlier draft reported 4 at Phase 0, conflating per-run capture with cross-run comparability.) At
closure the three build gates were collapsed into a `builds` block printing
`ff6lab ok / ff6demo ok / gui ok` with no status at all.

**Blank fields remain blank.** The visible output (`audit: clean`,
`census: clean`, `ff6lab ok`) supports a **separate inference** that these
gates passed, and that inference is recorded as such — it does not retroactively
become a captured zero.

AUDIT-0002's own Phase 0 run captured **11 of 11** required gates reliably.

## Open item requiring operator approval

**`/checkpoint` and `/session-summary` cannot be invoked under the current
allowlist.** Their contracts, established here as the named-command policy
requires:

- **`/checkpoint`** declares it will "Create a timestamped checkpoint under
  `docs/checkpoints/`, update `docs/checkpoints/LATEST.md`, and update
  **`dashboards/CURRENT_FOCUS.md`**." The third target is **not on the
  AUDIT-0002 allowlist** and is a substantive dashboard change the plan
  prohibits.
- **`/session-summary`** declares **no output path at all** — only "produce a
  factual session summary." Its `documentation-curator` skill mentions updating
  indexes, implying `indexes/SESSIONS.md`, but the path is never stated.
  Compatibility **cannot be verified** from the contract.

Neither may be invoked without explicit operator approval, and no
slash-command invocation credit is claimed for continuation. Continuation state
lives in `AUDIT-0002-status.md`, as planned.

This also explains, structurally, why AUDIT-0001 wrote its closure records by
hand: `/checkpoint` would have written a prohibited dashboard. AUDIT-0001 never
recorded that reasoning — it simply wrote the files and claimed neither
invocation, which is why its `Probable` classifications for these two commands
were honest even though the transcript shows the commands were invoked in
*other* sessions.
