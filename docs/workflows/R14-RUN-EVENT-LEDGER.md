# R14 provider-neutral run-event ledger

This document describes the atomic R14 correction that binds reconciliation
evidence to one immutable workflow run. It does not activate a provider
adapter, implement operator approval, enforce required outputs, generate
receipts, or correct the broader R12 contract model.

## Security boundary

The following are three separate properties:

```text
hash-chain integrity
does not prove event authenticity
and does not automatically satisfy a workflow requirement
```

The hash chain detects changes to the stored sequence. There is no signature,
trusted timestamp authority, or equivalent authentication mechanism in this
unit. A structurally valid, untampered event may still be self-reported,
manually imported, legacy evidence, unsupported, or ineligible for the
requirement whose selector it names.

Reconciliation therefore has two boundaries:

1. Verify the complete run identity and ledger. No event from an invalid ledger
   becomes an observation.
2. Evaluate each verified event's provenance, trust basis, event kind, and
   provider-supplied session/turn identity before it may satisfy a requirement.

Selector equality is necessary for invocation evidence but never sufficient.
Artifact and output observations remain artifact evidence only; they cannot
prove an agent, skill, command, or backend invocation.

## Immutable run identity

`workflow start` generates a 128-bit random `run_id` distinct from the
human-readable `workflow_id`. The immutable `identity.json` records:

```text
schema_version        workflow_id          run_id
contract_hash         created_at           repository_root
repository_identity   starting_branch      starting_head
```

Repository identity does not use the checkout directory name. The resolver
prefers a credential-safe normalized `origin` URL. If `origin` is absent but
other remotes exist, it hashes the sorted normalized remote set. A repository
with no configured remote is recorded honestly as its canonical Git common
directory, which is stable across linked worktrees but intentionally not
presented as a portable clone identity.

The canonical compact JSON encoding of the typed identity record, in its
declared field order, is SHA-256 hashed into durable state. A standalone
replacement of `identity.json` therefore invalidates verification. Linked
worktrees retain distinct `repository_root` values while sharing the same
no-remote common-directory identity.

The identity includes no ROM path, prompt, response, conversation content, or
private artifact content.

## Storage and permissions

Tracked contracts remain separate from ignored mutable and raw evidence:

```text
docs/workflows/runs/<workflow-id>/contract.json
local_artifacts/workflows/<workflow-id>/state.json
local_artifacts/workflows/<workflow-id>/<run-id>/identity.json
local_artifacts/workflows/<workflow-id>/<run-id>/events.jsonl
```

The runtime directories use mode `0700`; identity, state, ledger, and lock
files use `0600`. Creation is exclusive. Durable state anchors the canonical
identity hash as well as the ledger tail. A successful append is synced before
returning, followed by an atomic, synced update of the durable tail sequence
and hash in `state.json`.

## Normalized event and hash chain

Every JSONL event uses the schema in
`schemas/workflow-run-event.schema.json`. Optional provider values are explicit
empty strings or `null`, not omitted. Branch and HEAD are recorded on every
event; they may advance after the run starts.

Sequence 1 uses 64 zeroes as `previous_hash`. Later events use the preceding
event's `event_hash`. To calculate `event_hash`, the writer:

1. sets `event_hash` to the empty string;
2. marshals the typed event as compact JSON in its declared field order;
3. computes SHA-256 over those exact canonical bytes;
4. stores the lowercase hexadecimal digest.

Verification decodes each line with unknown fields rejected, recreates those
canonical bytes, and checks schema version, sequence, duplicate IDs, identity,
time bounds, repository identity, working directory, branch/HEAD shape,
previous-hash linkage, event hash, and the durable tail anchor.

The recognized normalized event kinds are:

```text
agent_started          agent_finished          skill_selected
tool_started           tool_finished           backend_finished
output_observed        validation_finished     operator_approval
```

Recognition is not a trust claim. `operator_approval`, for example, is
recognized structurally but is not eligible in this unit because operator
approval is explicitly out of scope.

## Provenance and eligibility

The normalized provenance categories are:

```text
provider_hook          deterministic_backend   operator_record
manual_import          legacy_transcript       unknown
```

Current eligibility is deliberately narrow:

- a collector-observed provider-hook `agent_started` event can satisfy an exact
  agent invocation requirement when its selector and provider session and turn
  identity are present. An `agent_finished` event remains valid lifecycle
  evidence but cannot independently satisfy an invocation requirement;
- collector-observed `skill_selected` events can satisfy exact skill
  invocation requirements under the same identity constraint;
- a provider-hook `tool_finished` event can become a backend observation only
  when it is collector-observed, carries provider session and turn identity,
  names the exact executed command in `selector`, and has a non-empty
  `tool_use_id` binding it to the distinct provider tool call. A captured
  non-null command exit status governs pass or failure. A bound completion
  with a null status remains ordering-visible only as unverifiable evidence;
  generic provider tool completion is never inferred to mean command success;
- deterministic `backend_finished`/`validation_finished` events can become
  backend observations; a null exit status remains unverifiable, never
  success;
- output events become artifact observations only;
- manual imports, legacy transcripts, unknown provenance, self-reported hook
  events, and operator records cannot satisfy invocation requirements here.

`collector_observed` and `backend_exit_status` are recorded trust bases, not
cryptographic proof. A future adapter must be reviewed and capability-tested
before its claims should be relied upon.

Each converted observation retains its verified ledger sequence. For an exact
deterministic-backend selector, the eligible observation with the greatest
sequence is the sole governing completion: exit zero satisfies, non-zero
fails, and null is unverifiable. Earlier matching completions are superseded.
No match is unsatisfied. Reconciliation does not use event timestamps, slice
position, or map iteration as a substitute for ledger order.

The current package-private deterministic executor obtains its command from an
exact `resource_id` in the active frozen contract, runs it as one `/bin/sh -c`
argument from the immutable repository root, and records the shell's actual
`ProcessState.ExitCode()`. It allocates sequence and revalidates the run,
identity, ledger, contract, and requirement under the append lock after
execution; stdout, stderr, and inherited environment are not ledger evidence.
Its complete boundary and platform limitations are documented in
[`R14-DETERMINISTIC-BACKEND.md`](R14-DETERMINISTIC-BACKEND.md).

## Damaged or incomplete evidence

The verifier never skips, repairs, reorders, or discards malformed history.
Changed, deleted, reordered, inserted, duplicated, incorrectly hashed,
malformed, truncated, identity-mismatched, time-invalid, or unsupported-schema
records invalidate the complete ledger. No events from that ledger are
converted. Relevant evidence channels are marked incomplete, integrity errors
are stored in state, and reconciliation cannot become `complete`.

The durable tail anchor detects deletion of the final complete event as well as
middle-record deletion. After terminal closure, appends and a second close are
rejected before state mutation.

## Creation and append failure behavior

Start prepares the ignored runtime bundle first and publishes the tracked
contract last. If identity creation succeeds but ledger creation fails, or a
later ordinary start step fails, the newly created runtime bundle and any new
contract file are removed. An existing run is never overwritten or repaired.

Under the per-run append lock, the writer first verifies the complete immutable
identity: its canonical hash must equal `state.identity_hash`, and its workflow,
run, and contract fields must equal durable state. A replaced or mismatched
identity is rejected before any ledger bytes are appended. Any append rejected
during validation leaves both the ledger and durable tail state unchanged. The
event record is appended and synced before the tail anchor is updated
atomically. If the process or filesystem fails after only one of those durable
writes, the ledger and anchor disagree; verification then marks the channel
invalid rather than guessing which side to repair. A partially written final
JSONL line is also invalid and remains available for external diagnosis.

## Default close path and compatibility

`ff6lab workflow close` now asks the store to verify the current run ledger and
uses only observations converted from it. It does not inspect
`~/.claude/projects`, recursively scan any provider-global history, interpret
Codex conversations, read a legacy gate TSV, or infer invocation from live
artifact existence.

No public command was added for asserting that an agent, skill, command, or
backend ran. The deterministic executor and event append method are
package-private until a reviewed adapter defines a narrower capability. Tests
append internal fixture events while exercising the same
identity, integrity, provenance, and eligibility checks future adapters must
use.

## Known limitations

- Events are hash-chained and tamper-evident, not signed or authenticated.
- There is no external signed anchor. A party able to rewrite the ledger,
  durable state anchor, identity, and tracked contract together can manufacture
  a self-consistent history; this unit detects damage, not a fully privileged
  forger.
- No Codex hook, Claude adapter, app-server integration, or transcript import
  is implemented.
- Operator approval remains structurally recognized but ineligible.
- Required-output and validation contract enforcement remains a separate R12
  correction.
- One human workflow ID still maps to one tracked contract in the existing
  store; this unit adds immutable run identity without redesigning that store.
- The exclusive append lock fails closed if a process dies while holding it;
  recovery requires external review rather than automatic lock removal.
