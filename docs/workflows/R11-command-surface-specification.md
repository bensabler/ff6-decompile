# R11 — operator-facing command surface: two-stage specification

**Status: specification only. Neither stage is implemented.**

This document adds no command-surface manifest, does not modify `ff6lab`, moves
no command file, and creates none of the seven proposed workflow commands.
Those are remediation actions and are out of scope here.

## Current state and target state, kept separate

Conflating these is the error this specification exists to avoid, and an earlier
draft of it made exactly that mistake by describing the repository as having
"50 commands" — adding proposed commands to existing ones.

```text
current command files:              43
current visible surface:            43
proposed primary workflow commands:  7   (do not exist)
target primary visible surface:      7
```

**No command file has been added.** The seven names below are proposals.

### Current 43, classified

From `AUDIT-0002-corrected-baseline.json`, on two independent dimensions.

| `surface_role` | n | | `operability` | n |
|---|---|---|---|---|
| `operator_workflow` | **0** | | `verified` | 15 |
| `internal_helper` | 33 | | `partial` | 13 |
| `diagnostic_manual` | 6 | | `unverified` | 11 |
| `alias` | 3 | | `blocked_backend_absent` | **3** |
| `deprecated` | 1 | | `not_applicable` | 1 |

`operator_workflow` is zero **against R11's own bar**, which R11 defines. That
circularity is disclosed in the remediation plan and is not re-argued here.

### Proposed 7 (do not exist)

`/research` · `/extract` · `/reconstruct` · `/implement` · `/validate` —
outcome commands. `/continue` · `/review` — lifecycle commands. All
argument-free and interactive.

No public `/plan`, `/execute`, `/verify` or `/close`: those are internal
lifecycle stages owned by the workflow engine, not operator decisions.

## The capability this specification cannot assume

**How a command is hidden from the operator surface is unknown.**

Established by inspection:

- Command files carry **no frontmatter** — they are plain Markdown bodies.
- `.claude/commands/` has **no subdirectories**.
- No project document describes a hiding, namespacing or visibility mechanism.

A common convention namespaces commands in subdirectories (`/internal:foo`),
but that is **untested in this environment** and cannot be verified from the
repository. AUDIT-0001 asserted that hooks "cannot be skipped" without testing
hook capability; that claim was later classified `unverifiable`. The same
mistake is available here and is refused.

**Stage 2 is therefore gated on a capability test, specified below.**

## Stage 1 — classification and enforcement

No file moves. No behavioural change to any existing invocation.

**Deliverables.** A tracked classification of every command on `surface_role`
and `operability`; an `ff6lab` gate; documentation surfaces presenting the
seven proposed commands as primary, the diagnostics separately, and the
backend-absent commands as blocked.

**Gate acceptance.** The gate fails when: any command file is unclassified; any
classified id has no command file; any command with
`operability: blocked_backend_absent` carries `surface_role: operator_workflow`;
or any command claims `operator_workflow` without a contract-producing path.
It must fail today on an unclassified command added by hand, and pass after
classification — demonstrated, not asserted.

**Why a gate and not documentation.** Two hand-maintained inventories already
disagreed with each other and with the filesystem; R3 exists because of it. A
classification without machine enforcement would drift the same way.

**Rollback.** Delete the manifest and the gate registration. Nothing else
changes, because nothing else moved.

**What Stage 1 does not achieve.** The raw slash-command list still shows 43.
Stage 1 makes the classification true, enforced and visible in documentation;
it does not shorten the list. **Reporting Stage 1 as "43 → 7" would be false.**

## Stage 2 — visible-surface reduction, gated on a capability test

### The capability test, which must run and pass first

Before any file is moved, determine and record:

1. Does a command in `.claude/commands/<dir>/name.md` appear at all?
2. Under what identifier — `/name`, `/dir:name`, or something else?
3. Does the flat surface still list it?
4. Does an existing `/name` keep working after the move?
5. Are subdirectory commands presented differently from top-level ones?

Method: move **one** low-risk command — `/bootstrap-v4`, already
`deprecated` / `not_applicable` — into a subdirectory, restart a session, and
record what the surface shows. Preserve the observation.

Classify the result `verified` / `partially verified` /
`available but untested` / `unavailable` / `unknown`. **Stage 2 proceeds only
on `verified`.**

If namespacing is unavailable, Stage 2 is **withdrawn**, not worked around: the
honest outcome is that the surface cannot be shortened by this mechanism, and
Stage 1's classification is what the project gets.

### If verified

Move 33 `internal_helper` commands to an internal namespace and 6
`diagnostic_manual` commands to an advanced namespace. Retire the 3 aliases
after a documented migration period. Remove the 1 deprecated command after a
compatibility review. Leave the 3 blocked commands visible **and visibly
blocked**, naming the absent backend — a blocked command must never read as
operational.

**Acceptance.** The default surface lists exactly the primary workflow
commands; the advanced set appears only on the advanced surface; helpers are
absent from the default surface; every moved command remains invocable under
its new identifier; and every documentation reference resolves.

**Migration impact — not zero.** 39 files change path. Every reference in
`.claude/README.md`, `docs/WORKFLOW_COMMANDS.md`, playbooks, checkpoints and
session records must be updated or knowingly left as historical. Operator
muscle memory for existing names breaks. Describing this as "no migration
impact" would repeat the error AUDIT-0002 found in AUDIT-0001's routing
proposal.

**Rollback.** Move the files back. Recorded per file, so the inverse is exact.

## Dependencies

Stage 1 is independent. **Stage 2 depends on the capability test.** Both depend
on the seven commands existing, which is separate work: a surface reduction
that hides 39 commands before their replacements exist would leave the operator
with fewer capabilities, not better ones.

Sequence: seven commands (contract-driven, on R12/R14) → Stage 1 → capability
test → Stage 2 or withdrawal.

## Open questions for the operator

1. **Do the seven commands get built before Stage 1, or alongside it?** This
   specification assumes before; hiding helpers first would strand the operator.
2. **Is `/bootstrap-v4` acceptable as the capability-test subject?** It is
   deprecated and non-applicable, so the blast radius is minimal, but it is
   still a real command file.
3. **What migration period do the three aliases get?** They are deliberate and
   self-documenting; removing them early breaks a compatibility promise the
   project made on purpose.
