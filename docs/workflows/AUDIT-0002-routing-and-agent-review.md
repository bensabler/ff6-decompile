# AUDIT-0002 — routing and agent review (Phase 7)

## AUDIT-0001's Phase 4 process metrics: unverifiable

`agent-smoke-tests.py` makes **zero repository reads** — it is hand-authored
prose. **No raw agent transcripts were preserved.** Contract-compliance ratings,
durations, word counts and per-agent quality judgements are therefore
**unverifiable from currently available repository evidence** — not refuted, and
not permanently unverifiable, since the findings those agents produced about the
repository remain independently checkable and several were re-verified here.

What *is* now verifiable: **13 `Agent` calls with 13 distinct project
`subagent_type` values** in session `9cde5fa3`, from harness evidence. The
invocations happened. Their qualitative assessment cannot be re-derived.

## Five routing dimensions, separately classified

| Dimension | Status | Evidence |
|---|---|---|
| Direct named invocation | **Confirmed** | 13 AUDIT-0001 calls + 2 AUDIT-0002 probes, all harness-recorded |
| Automatic selection | **Refuted for the tested mechanism** | Probes P1/P2 |
| Implicit description matching | **Refuted for the tested mechanism** | Probes P1/P2 |
| Command-triggered selection | **not_tested** | No command invocation was authorised to test it |
| Orchestrator selection without an agent name | **not_tested** | `/orchestrate` not invoked |

## The probes

Expectations were **pre-registered before running** (`agent-transcripts/pre-registration.md`).

The `Agent` tool's contract states that if `subagent_type` is omitted, the
general-purpose agent is used. If accurate, automatic selection of a project
specialist cannot occur — routing is explicit-by-orchestrator or absent, with
no third possibility.

| Probe | Task, no agent named | Predicted | Observed |
|---|---|---|---|
| **P1** | Verification task strongly matching `verification-engineer` | `general-purpose` | **`general-purpose`** |
| **P2** | Documentation-review task strongly matching `documentation-reviewer` | `general-purpose` | **`general-purpose`** |
| **P3** | Negative control requiring no delegation | no delegation | control, not delegated |

Both probes were read-only, bounded, prohibited nested delegation, and complied.

### Interpretation, bounded as pre-registered

**Automatic and implicit-description routing to project specialists do not
occur under the tested mechanism.** Omitting the name yields `general-purpose`
— a harness agent, not one of the 13.

This does **not** prove routing is globally absent. It proves that the one
mechanism by which a project resource could cause automatic selection does not
function that way. Command-triggered selection and `/orchestrate` behaviour
remain untested.

**Consequence for remediation:** adding a routing table to
`ORCHESTRATE_RESEARCH.md` helps the *orchestrator* choose. It cannot create
automatic selection, because no such mechanism exists. AUDIT-0001's P0-1
proposal is therefore correctly scoped but must not be described as enabling
automatic routing — and its claimed "no migration impact" is wrong, since it
changes delegation behaviour.

## Routing measured correctly

AUDIT-0001 used **textual inbound** as its orphan proxy. That proxy was wrong
in both directions: it missed `PACKAGE_MANIFEST.json` (under-counting
references) and it counted packaging and narrative mentions as if they were
routing (over-counting).

| Metric | Agents |
|---|---|
| Textual inbound = 0 | 12 of 13 |
| Textual inbound = 0, incl. root files | 1 of 13 |
| **Routing-bearing inbound = 0** | **13 of 13** |

Routing-bearing means inbound from a resource that would *use* it — a command,
skill, playbook, rule, template, or the constitution. **No command, skill,
playbook, rule or template names any of the 13 agents.**
`ORCHESTRATE_RESEARCH.md` step 6 remains the single word "Delegate."

The finding AUDIT-0001 reported survives and strengthens. Its measurement did
not.

## Probe by-products

Both probes independently corroborated repository facts while performing their
cover tasks:

- `RUN_QUALITY_GATES.md` has 8 numbered steps naming `gofmt`, `go test` and
  `go vet`, and **no `go build`** — corroborating the gate-coverage finding.
- `.claude/README.md` lists **28** commands and documents a **"Workflows"**
  component category with no corresponding directory — corroborating both the
  documentation-drift and phantom-category findings.

## Correlated, not independent

These probes and the Phase 11 reviews are **separate bounded reviews by agents
running on the same underlying model as the primary auditor.** They reduce
single-pass error and surface missed evidence. They are not independent model
review and are not described as such anywhere in AUDIT-0002's output.
