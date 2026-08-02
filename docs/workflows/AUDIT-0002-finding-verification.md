# AUDIT-0002 — finding verification (Phase 6)

Each finding independently re-checked in the worktree. Provisional findings
received **contract analysis before any repair was proposed**.

## Capability-honesty concerns — all re-verified

| # | Finding | Status | Re-verification |
|---|---|---|---|
| 1 | `internal/validate` has no operational consumer | **confirmed** | `go list` — **zero** non-test importers **and zero test importers** outside the package. `VALIDATE_GRAPHICS.md:16` still reads "Compare indexed pixels when possible" — the failure mode `CompareResolved` exists to prevent |
| 2 | `internal/audio/doc.go` overstates its package | **confirmed** | Doc claims "SPC700, DSP, sequence, and sample reconstruction code". Only subpackage: `brr`. `internal/audio/{sequence,spc700,dsp}` all **absent** |
| 3 | Audio/SPC commands lack usable backends | **confirmed** | `/recover-sequence`, `/trace-spc-command`, `/validate-audio` have no backing package. `internal/validate` is framebuffer-only |
| 4 | `/trace-dma` is four surfaces | **confirmed** | `snesdma.go` + `trace.go` + both tests present; `dma-trace.lua` header still declares `STATUS: UNEXERCISED` |
| 5 | `/audit-project` has no report-only mode | **confirmed** | Command authorises mechanical fixes; playbook's Required outputs mandate updated records, index updates and a checkpoint |
| 6 | `ff6lab state origin`/`sprites` undiscoverable | **confirmed** | Both present in `cmd/ff6lab/state.go`; **0** occurrences in help text |
| 7 | Explicit agent routing absent | **confirmed, measurement corrected** | All 13 agents have zero routing-bearing inbound |

## Provisional findings — contract analysed

Twelve findings AUDIT-0001 asserted as defects without establishing the
intended contract. **Four are withdrawn or downgraded.**

| Finding | Verdict | Contract analysis |
|---|---|---|
| 9 templates unreferenced | **partially supported** | All 9 appear in `PACKAGE_MANIFEST.json`. "Unreferenced" is false; **10 of 12 lack a consumer reference**, which is the real defect |
| Rules with zero inbound are unreachable | **REFUTED** | All four carry `paths:` frontmatter and are **auto-loaded by glob scope**. Never being named is their design. AUDIT-0001's `Clarify Routing` recommendation withdrawn |
| `LATEST.md` duplication is a defect | **REFUTED as a defect** | The convention is deliberate restart-safety: a resuming session reads one file and gets the full checkpoint without a second hop. The real risk is *uncontrolled divergence*, not duplication. Recommend a consistency check, not removal |
| `experiment.schema.json` omits constitutional fields | **downgraded** | The schema validates `manifests/experiments.json`, a **compact index**, not the complete experiment record. Records live in Markdown under `docs/experiments/`. Requiring the full seven fields in an index schema would be a category error. The genuine gap: nothing validates that the *records* carry them |
| `domain` is invalid | **REFUTED** | The schema sets no `additionalProperties: false`, so `domain` is an **allowed unvalidated extension**, not a violation. It is used by `CheckBattleExperimentConfig` |
| `ACTIVITY_LOG.md` must list every unit | **not established** | No contract requires per-unit completeness. Session 005's absence is real; whether it is a defect depends on an unstated cadence |
| `OPEN_HYPOTHESES.md` must hold all hypotheses | **not established** | Its own text scopes it to tracked hypothesis IDs, not to every checkpoint-level alternative |
| `VALIDATE_GRAPHICS.md` prescribes indexed-only | **confirmed** | Step 4 read in full context still says "Compare indexed pixels when possible" and never mentions resolved-colour comparison |
| "Workflows" category implies a missing directory | **partially supported** | It may be conceptual. But it sits in a list where all six siblings are directories, and probe P2 independently read it as a missing directory. Ambiguous documentation is the defect |
| Local gates must duplicate CI | **REFUTED** | CI enforces the demo binary, restricted files, generated assets and the GUI build. Local gates need not reproduce platform-specific CI work. The defect is that `RUN_QUALITY_GATES.md` never *mentions* these gates exist |
| `ASSET-GFX-0002` census link is wrong | **confirmed, unrepaired** | Narshe tileset asset links `CEN-GFX-0006`, "Mines-interior PPU configuration". Research records are outside AUDIT-0002's write scope |
| Discovery→demo must be universal | **REFUTED as stated** | Not every discovery belongs in DEMO-0001. The requirement is an explicit **applicability decision**, not universal representation |

## Commands promising unsupported behaviour

Re-derived from backends rather than from AUDIT-0001's classification:

| Command | Promise | Backend |
|---|---|---|
| `/validate-audio` | DSP/voice/waveform comparison | **absent** |
| `/recover-sequence` | sequence format recovery | **absent** |
| `/trace-spc-command` | SPC700 port/driver tracing | **absent** |
| `/trace-dma` | live DMA tracing | **written, never executed** |
| `/validate-graphics` | graphics validation | **present but unreachable** — zero importers |
| `/audit-project` | assessment | **present but unusable for assessment** — mutation-only |

`/trace-dma` remains the honesty model: its probe header, playbook and command
all state their limits plainly. High honesty with incomplete backing is a
different condition from low honesty, and the two must not share a score.

## Playbook dispositions — all 28, individually assessed

AUDIT-0001 returned `Keep` for all 28 from an unconditional `return`. That is
withdrawn. Each playbook is now assessed on **content** — contract
specificity, backend truth, overlap, outcome evidence, maintenance risk — not
structure. Structure could not discriminate them: all 28 carry the same three
sections, 6-9 steps, 27-57 lines, which is precisely why the old scorer emitted
an identical 4.67 for every one.

Per-playbook basis and confidence: `AUDIT-0002-corrected-baseline.json`.

| Disposition | n | Playbooks |
|---|---|---|
| **Keep** | 14 | incl. `TRACE_DMA` and `RECOVER_COMPRESSION` — cited as the honesty model |
| **Keep and Improve** | 5 | `RECOVER_TEXT`, `RECOVER_EVENT_OPCODE`, `IMPLEMENT_DISCOVERY`, `PREPARE_RELEASE`, `RECOVER_BACKGROUND` |
| **Repair** | 4 | `ORCHESTRATE_RESEARCH`, `AUDIT_PROJECT`, `VALIDATE_GRAPHICS`, `RUN_QUALITY_GATES` |
| **Narrow Scope** | 3 | `VALIDATE_AUDIO`, `RECOVER_SEQUENCE`, `TRACE_SPC_COMMAND` — backends absent |
| **Merge Candidate** | 1 | `CAPTURE_GRAPHICS` — overlaps `CAPTURE_FRAME` |
| **Deferred** | 1 | `WEEKLY_REVIEW` — fully specified, but no identifiable output |

The four `Repair` cases share a pattern: the playbook prescribes something the
backing reality contradicts. `ORCHESTRATE_RESEARCH` step 6 is the single word
"Delegate" at the project's highest-frequency entry point; `AUDIT_PROJECT`'s
required outputs make report-only execution impossible; `VALIDATE_GRAPHICS`
step 4 prescribes the failure mode its own validator exists to prevent;
`RUN_QUALITY_GATES` names no build step while CI enforces four.

## Capability axes are contract-relative

`specification_status` records **specification quality**. Compatibility with a
particular audit's restrictions is recorded separately as
`audit_scope_compatibility`.

- **`/checkpoint` — `specification_status: complete`.** Its contract is
  coherent: it names its skill, its template, and all three write targets
  explicitly. Its `audit_scope_compatibility` is `incompatible`, because one of
  those targets is a dashboard a bounded audit may not write. **A command can
  be well specified and still be unusable inside a restricted scope.**
- **`/session-summary` — `specification_status: partial`**, not `missing`. It
  has meaningful procedural content and names its skill; it lacks a declared
  output path. `audit_scope_compatibility: unverifiable`.

## Not carried forward

AUDIT-0001's per-command implementation-status classification is retained as a
**lead**, not as a verified axis. It was produced by a hand-authored script
(`classify_capabilities.py`, zero repository reads); only the seven concerns
above and the six commands here were independently re-derived.
