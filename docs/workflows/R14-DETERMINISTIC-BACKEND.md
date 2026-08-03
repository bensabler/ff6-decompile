# R14 package-private deterministic backend executor

This unit adds one narrow production path: execute the command selected by an
exact `resource_id` in the active frozen workflow contract, observe the local
process status, and append a normalized completion to the private R14 run-event
ledger. It is not a public command runner, approval system, hermetic sandbox,
output collector, validation engine, or authenticity mechanism.

## Frozen command selection

The package-private caller supplies only a workflow ID and an exact requirement
ID. The caller does not supply an executable command, exit status, pass/fail
claim, event, provenance, sequence, or workflow verdict.

Before launching a process, the executor:

1. loads durable state and binds the operation to its current `run_id`;
2. acquires the per-run ledger lock briefly, reloads state, and requires the
   same run to be in the `executing` phase;
3. verifies the immutable run identity and the complete ledger against the
   durable identity and tail anchors;
4. loads the stored contract, requires `state: frozen`, runs the existing
   frozen-hash verification, and requires `contract.frozen_hash` to equal
   `state.contract_hash`;
5. finds exactly one requirement whose `resource_id` equals the caller's
   requirement ID byte-for-byte; and
6. requires that requirement's `execution_mode` to be
   `deterministic_backend`.

The selected requirement's exact `resource_id` becomes both the command string
and the eventual event `selector`. Backend eligibility is not inferred from
`resource_type`, command spelling, a filename, or `evidence_rule` prose.
Missing, duplicate, ambiguous, non-backend, unfrozen, hash-mismatched, closed,
identity-invalid, ledger-invalid, or otherwise non-executable runs fail before
the runner is invoked.

The preparation lock is released before external execution. The ledger lock is
never held while the command runs.

## Process and shell semantics

The production runner constructs this exact argument vector:

```text
argv[0] = /bin/sh
argv[1] = -c
argv[2] = <the exact frozen resource_id string>
```

The command is passed as one `-c` argument. Go does not split it with
`strings.Fields`, add quoting, or accept additional caller-provided arguments.
`/bin/sh` interprets that string, so quoting, expansion, redirection, pipelines,
and command lists have the shell's normal meaning. The observed process is the
launched `/bin/sh`; its status reflects the shell's result for the command it
interpreted.

The working directory is exactly the immutable run identity's canonical
`repository_root`. The executor verifies the current repository root and
credential-safe repository identity against that binding before launch and
again before append. It does not substitute the caller's current directory.
Branch and HEAD are allowed to advance during a run and are resolved again for
the completion event.

The `exec.Cmd` environment is left unset, so the child inherits the executor
process's environment. No environment snapshot, variable, value, or credential
is written to `state.json` or the event ledger. This makes execution compatible
with ordinary developer tooling, but not hermetic or reproducible solely from
the ledger.

A `context.Context` controls the OS process through `exec.CommandContext`.
This unit adds no public timeout policy. Package-internal callers may route
stdout and stderr to supplied `io.Writer` destinations. The executor neither
stores nor parses those streams.

## Process outcomes and infrastructure failures

After `Start`, the runner calls `Wait` and reads
`cmd.ProcessState.ExitCode()`. It never derives status from output, error text,
known command names, or the success of a wrapper operation.

| Condition | Executor result | Ledger behavior |
| --- | --- | --- |
| Shell starts and exits 0 | observed status `0` | append `backend_finished` with `exit_status: 0` |
| Shell starts and exits nonzero | observed nonzero status; this is not an infrastructure error | append `backend_finished` with the actual nonzero status |
| Shell cannot be started | launch error; no observed terminal status | append nothing and fabricate no status |
| Process starts but no terminal `ProcessState` is available | execution error; status is unknown | append nothing |
| Process terminates, but post-execution verification or persistence fails | return the observed status together with a clear no-credit error | do not present the result as trustworthy ledger credit |
| Precondition or integrity check fails | error before invocation | append nothing |

A recorded nonzero completion is valid evidence that the backend ran and
failed; the executor itself returns it as an observed result rather than a
launch failure. Reconciliation applies the existing latest-sequence rule:
status zero satisfies the matching deterministic-backend requirement, any
nonzero status fails it, and a missing status is unverifiable. Misleading output
such as `PASS`, `success`, or `exit 0` cannot change that decision.
A command-not-found error produced by an already launched shell is likewise a
real nonzero shell outcome (commonly status 127), not a shell launch failure.

If persistence fails before any bytes are appended, there is no completion
record. A write or sync failure may instead leave a partial or unanchored line;
if the record is appended and synced but the durable tail update fails, the
ledger and state anchor disagree. In every such case, the executor returns an
error saying that no trustworthy credit was recorded. Verification fails
closed, and any residual bytes remain available for external diagnosis.

## Output is not evidence in this unit

stdout and stderr are execution streams only. Their bytes are not:

- copied into the `RunEvent`;
- written to `state.json`;
- parsed for pass/fail or exit status;
- transformed into `output_observed` events; or
- written to a new command-log format.

The executor does not check declared required outputs after the command. Output
existence and artifact evidence remain separate R12/R14 work.

## Completion event and eligibility

For a shell that reaches an observable terminal state and passes post-execution
verification, the executor constructs these fields from code rather than from
caller assertions:

| Field | Value or source |
| --- | --- |
| `schema_version` | current typed event schema version |
| `event_id` | newly generated private event ID |
| `workflow_id`, `run_id`, `contract_hash` | verified immutable identity and durable state |
| `observed_at` | UTC time sampled after process termination |
| `provider` | `local-process` |
| `source_kind` | `deterministic_backend` |
| `collector_id` | `internal/workflow.deterministic-backend` |
| `trust_basis` | `backend_exit_status` |
| `repository_identity`, `cwd` | reverified repository identity and immutable repository root |
| `branch`, `head` | repository state resolved at append time |
| `event_kind` | `backend_finished` |
| `selector` | exact frozen requirement `resource_id` that was executed |
| `exit_status` | `ProcessState.ExitCode()` from the launched shell |
| `session_id`, `turn_id`, `tool_use_id`, `evidence_ref` | explicit empty strings |

The ledger writer assigns `sequence`, `previous_hash`, and `event_hash`; the
executor cannot preselect them. No provider session, turn, agent, skill, or tool
identity is manufactured. Command spelling never changes the event kind to
`validation_finished`.

After full-ledger verification, the existing eligibility boundary converts
this event to a backend observation because it has deterministic-backend
provenance, captured-exit-status trust basis, a recognized backend-completion
kind, and a non-empty selector. The executor always supplies a non-null status.
Selector equality remains necessary for reconciliation, and the greatest
eligible ledger sequence for that exact selector governs.

These fields record the implementation's observation and trust basis. They are
not a signature, trusted timestamp, remote attestation, or cryptographic proof
that the command was authentic. The hash chain and durable tail make mutation
detectable under the R14 threat model; they do not prevent a fully privileged
party from rewriting all local records consistently.

## Sequence allocation and post-execution revalidation

Once the shell terminates, the executor reacquires the per-run ledger lock and
uses the sequence-allocating append helper. Under that lock, the helper:

1. requires the expected `run_id` still to be current and the run not to be
   closed;
2. revalidates immutable identity against durable state;
3. verifies the complete ledger and its durable tail anchor;
4. requires the phase still to be `executing`, the contract binding to match
   the pre-execution plan, and the immutable repository root to be unchanged;
5. reloads and verifies the frozen contract and exact deterministic-backend
   requirement;
6. re-resolves repository identity, branch, and HEAD;
7. rejects duplicate event IDs; and
8. selects `tail_sequence + 1` and the verified tail hash, validates and hashes
   the event, appends and syncs it, then atomically updates durable tail state.

Consequently, an unrelated event appended while the backend runs receives its
own sequence first; the backend completion takes the next sequence from the
then-current tail. A second completion for the same selector receives a later
sequence and supersedes the earlier completion for reconciliation.

If the run closes before this append boundary, if the run ID changes, if the
identity or ledger becomes invalid, or if the active contract/requirement
differs from the one executed, no backend completion is credited. The returned
error states that the command may have run and reports its observed status, but
that no trustworthy ledger credit was recorded. A transient local replacement
restored to the same verified bytes and bindings before revalidation has no
separate history record in this unit.

## Deliberate boundaries

This executor remains package-private. There is currently no
`ff6lab workflow backend` or `ff6lab workflow run` command and no public Store
method that accepts an arbitrary command, event, exit status, pass/fail claim,
or verdict.

This unit also does not:

- require or record operator approval before execution;
- sandbox commands or filter the inherited environment;
- enforce declared required outputs;
- enforce validation declarations or emit `validation_finished`;
- expose a hook endpoint or Codex lifecycle adapter; or
- import JSON events, transcripts, or caller-authored evidence.

Codex `PostToolUse` is not used to decide backend pass/fail. The retired R14
hook probe observed command completion but no explicit child exit-status field;
tool completion therefore could not distinguish the controlled exit-0 and
exit-7 cases reliably. A `PostToolUse` record can establish that the hook ran,
not the status required here. The local executor instead launches the frozen
command and reads the OS process state directly. See
`docs/workflows/R14-CODEX-HOOK-PROBE.md` for the bounded probe findings; its
diagnostic hooks remain retired.

## Platform and signal limitations

- The production rule requires `/bin/sh` and is intended for the repository's
  macOS and Linux developer environments. Native Windows execution is not
  implemented.
- `/bin/sh` implementations and installed external commands may differ between
  macOS and Linux. Contracts should use portable shell syntax when both are
  required; the ledger does not capture the shell version or toolchain.
- `ProcessState.ExitCode()` returns `-1` when the launched shell was terminated
  by a signal. A surviving shell may instead encode a nested command's signal
  as a conventional nonzero shell status such as `128 + signal`. The event
  model records only the integer, not signal identity, core-dump state, or a
  portable explanation; reconciliation treats every representation as
  nonzero failure.
- Context cancellation uses Go's default `CommandContext` behavior for the
  launched shell. This unit adds no process-group management or guarantee that
  descendants started by shell syntax are terminated.
- Inherited environment, filesystem state, repository changes, and external
  services can affect the result. Exact frozen selection is deterministic; the
  execution environment is not claimed to be.
