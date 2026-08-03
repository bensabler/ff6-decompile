# R14 Codex hook capability probe

This report characterizes the Codex hook surface needed by the provider-neutral
R14 run-event ledger. It is a diagnostic capability probe, not a production
adapter. Nothing in this unit writes to the workflow ledger, asserts trusted
`provider_hook` provenance, implements contract approval, or adds agent
completion requirements. The follow-up hardening changes only which verified
agent lifecycle event may become invocation evidence.

```text
probe date: 2026-08-02
foundation commit: df3a0c6aa4e8c078b5795bba34490a1b3b63c599
probe working branch: codex/r14-codex-hook-probe
installed CLI: codex-cli 0.146.0-alpha.9.2
installed app: 26.727.51351 (build 6119)
setup-task hook definitions recognized or loaded: not established
setup-task hook definitions trusted: no
setup-task hook handler execution observed: no
controlled live-probe matching handler execution: locally observed
controlled sanitized summary record count: 12
preserved archive record count: 414
post-probe spillover record count: 402 by count, assuming append-only operation
preserved archive filename: local_artifacts/codex-hook-probe/events-probe-and-postprobe-20260802.jsonl
preserved archive SHA-256 (operator supplied): 86d00525ae7f5d4661ed82e8a3f2c6c5c289fa7c40d40435689b1a77238bbda8
repository hook registration: retired after the probe
live hook probe: completed; do not rerun
```

The operator reports that the controlled probe was completed after the setup
commit. The accepted capability conclusions come only from the original
12-record sanitized summary. Those supplied matching sanitized records
establish handler execution; they do not independently record the exact
hook-trust UI state, which remains unknown. The probe must not be rerun, and the
preserved raw JSONL must not be inspected in this unit.

## Evidence vocabulary

This report keeps these classifications distinct:

- **officially documented**: stated by current public Codex documentation;
- **locally observed**: returned by the installed runtime or present in a real
  sanitized hook capture;
- **inference**: a bounded conclusion from documented and locally observed
  facts, not direct evidence by itself;
- **unknown**: the available documented and local evidence does not resolve the
  question;
- **unsupported**: the documented or directly tested surface does not provide
  the required capability; every such statement names its scope;
- **not observed**: the controlled capture did not produce the evidence;
- **not testable**: the current documented/configurable surface offers no
  corresponding probe point.

Synthetic test fixtures establish sanitizer behavior only. They are never
classified as live hook evidence.

## Setup-time recovered state and completed phases

This section is a historical record of the setup task, before the later
controlled live probe. Recovery found the branch at the required foundation
commit with three untracked setup files and no probe documentation or sanitized
output. No pre-existing repository hook configuration was overwritten.

The completed work before recovery was:

1. foundation branch, commit, worktree, and remote-ref verification;
2. read-only runtime, configuration-layer, trust, and hook-capability review;
3. creation of the diagnostic collector, its synthetic tests, and the
   repository-local hook configuration;
4. an initial focused synthetic test run.

The hook definitions were not established as loaded and were not trusted in
that already-running setup task. No live `SessionStart`, subagent, Bash success,
Bash failure, or skill-selection event was captured in that task. The setup
report was then validated and committed for the later operator-controlled probe
whose sanitized results are recorded below.

## Runtime, configuration, trust, and reload

| Question | Classification | Result |
|---|---|---|
| Installed hook feature | locally observed | `codex features list` reports `hooks` as stable and enabled. |
| User base configuration | locally observed | `~/.codex/config.toml` loads; it does not explicitly disable hooks. It was not modified. |
| Repository trust | locally observed | The user configuration marks this repository trusted. |
| Existing local hook definitions | locally observed | None existed in the repository, user hook file, user inline configuration, enabled plugin hook paths, or active Git hooks before this setup. |
| Local managed files | locally observed | No `/etc/codex/requirements.toml`, `/etc/codex/managed_config.toml`, or readable macOS MDM hook policy was found. |
| Cloud/session policy | unknown | Local file inspection cannot prove that no cloud-delivered or session-only policy exists. |
| Project-layer loading | officially documented | Project `.codex` configuration loads only for a trusted project. |
| Command-hook trust | officially documented | A non-managed command hook is skipped until its exact definition is reviewed and trusted. Project trust alone is insufficient. |
| Hook hot reload | unknown | The documented app-server API lists `hooks/list`, but no hook reload method. The setup task did not establish hot reload. |
| Full app restart | inference | A restart was a conservative fallback if newly trusted definitions still appeared skipped, but the documentation did not establish that it was always required. |

Current references are the official [Hooks](https://developers.openai.com/codex/hooks),
[Advanced configuration](https://developers.openai.com/codex/config-advanced),
[Managed configuration](https://developers.openai.com/codex/enterprise/managed-configuration),
and [App Server](https://developers.openai.com/codex/app-server) documentation.

## Diagnostic setup and retirement

The historical tracked setup consisted only of:

```text
.codex/hooks.json
.codex/hooks/codex_hook_probe.py
.codex/hooks/test_codex_hook_probe.py
docs/workflows/R14-CODEX-HOOK-PROBE.md
```

After the controlled probe and subsequent diagnostic spillover, this retirement
unit removed `.codex/hooks.json` to prevent further automatic diagnostic
collection. No replacement repository, plugin, or user hook is introduced by
this retirement unit, and no other tracked repository file actively registers
the collector. The collector and its synthetic tests remain tracked as inactive
reference material only; tracked repository configuration no longer registers
them for automatic execution. Hook retirement does not modify production R14
ledger behavior.

During the active diagnostic phase, `.codex/hooks.json` registered observational
command hooks for only the required events:

| Hook event | Matcher | Purpose |
|---|---|---|
| `SessionStart` | `startup|resume` | session payload shape |
| `SubagentStart` | `.*` | subagent start identity shape |
| `SubagentStop` | `.*` | subagent stop identity shape |
| `PreToolUse` | `^Bash$` | Bash identity and exact controlled input |
| `PostToolUse` | `^Bash$` | Bash completion and response shape |

While the configuration was active, each handler invoked the same
repository-local collector, returned success, and provided no hook decision or
additional context. Matching handlers could run concurrently; the collector
made no inference from arrival order.

When invoked during the active diagnostic phase, the collector read at most one
1 MiB UTF-8 JSON object from standard input and appended one sanitized record
beneath the already ignored path:

```text
local_artifacts/codex-hook-probe/events.jsonl
```

That was the collector's historical configured write path. The accumulated
archive was later preserved under the normalized filename recorded below.

The directory is forced to mode `0700` and the regular, non-symlink JSONL file
to `0600`. Appends use an inter-process exclusive lock and `fsync`. Each record
has an independent random diagnostic capture ID and collector timestamp; it
has no workflow identity, sequence, ordering claim, or provider timestamp
claim.

The collector records only:

- hook event name;
- top-level field names and JSON types;
- selected-field presence, type, byte length, and non-path identifier digest;
- bounded `tool_input` and `tool_response` field/type/length shapes;
- an explicit exit-status field's path, type, and numeric value, if one really
  exists in the payload;
- the collector-generated UTC timestamp and diagnostic capture ID.

Prompts, responses, transcripts, environment contents, credentials, ROM or
private contents, complete tool output, arbitrary string values, and arbitrary
MCP payload values are not retained. The only commands eligible for verbatim
retention are:

```text
printf 'codex-hook-probe-success\n'
sh -c 'printf "codex-hook-probe-failure\n" >&2; exit 7'
```

Every other command is represented only by its UTF-8 byte length and SHA-256
digest. Other tool strings retain length but no content digest. Tool events
whose `tool_name` is not exactly `Bash` are rejected even if the configuration
matcher were bypassed. Input containing production fields such as `workflow_id`, `run_id`,
`contract_hash`, provenance/trust claims, observations, reconciliation, or a
verdict is rejected. Collection errors return success without steering Codex
and do not create a production event.

## Sanitized live-capture summary

The accepted capability evidence is the operator-supplied sanitized summary of
the original controlled 12-record checkpoint. This report does not derive any
event facts from the preserved raw JSONL.

| Hook event | Record count |
|---|---:|
| `SessionStart` | 1 |
| `SubagentStart` | 1 |
| `SubagentStop` | 2 |
| `PreToolUse` | 4 |
| `PostToolUse` | 4 |
| **Total** | **12** |

The supplied findings characterize the two controlled command pairs identified
below. They do not assign a purpose or outcome to the other `PreToolUse` and
`PostToolUse` records, and this report makes no inference about them.

## Official payload matrix versus local evidence

Codex documents these common hook fields: `session_id`, `transcript_path`,
`cwd`, `hook_event_name`, and `model`. Turn-scoped events add `turn_id`. The
required events also document `permission_mode`. Transcript fields are never
read by this probe.

| Event | Additional officially documented fields | Locally observed fields | Status |
|---|---|---|---|
| `SessionStart` | `source`; no documented `turn_id` | `session_id`, `cwd`, `permission_mode` | locally observed |
| `SubagentStart` | `turn_id`, `agent_id`, `agent_type` | `session_id`, `turn_id`, `agent_id`, `agent_type` | locally observed |
| `SubagentStop` | `turn_id`, `agent_id`, `agent_type`, `agent_transcript_path`, `stop_hook_active`, `last_assistant_message` | `session_id`, `turn_id`, `agent_id`, `agent_type` | locally observed |
| `PreToolUse` for Bash | `turn_id`, `tool_name`, `tool_use_id`, `tool_input`; Bash input documents `command` | `session_id`, `turn_id`, `tool_name`, `tool_use_id`, exact safe command | locally observed |
| `PostToolUse` for Bash | the PreToolUse identity/input fields plus `tool_response` | `session_id`, `turn_id`, `tool_name`, `tool_use_id`, exact safe command; no explicit exit-status field | locally observed |

Documentation states that `PostToolUse` runs after a tool completes, including
when the underlying command exits nonzero. It does not document a stable Bash
exit-status field inside `tool_response`. Therefore generic hook completion is
not evidence that the command succeeded.

## Required capability questions

### Agent invocation

| Required value | Official surface | Local capture |
|---|---|---|
| `session_id` | documented | observed on start and stop |
| `turn_id` | documented for subagent events | observed on start and stop |
| `agent_id` | documented | observed on start and stop |
| `agent_type` | documented | observed on start and stop |
| `SubagentStart` identity | documented by `hook_event_name` plus agent fields | observed |
| matching start/stop lifecycle | both events document the same identity fields | observed |

One complete lifecycle was locally observed: a `SubagentStart` and
`SubagentStop` had matching `session_id`, `turn_id`, `agent_id`, and
`agent_type` values. A second `SubagentStop` had the same agent ID and type but
a different turn ID and no corresponding captured start. Its cause is
**unknown**. It is neither an independently proven invocation nor skill
evidence.

Agent start identity and a matching start/stop lifecycle are therefore locally
observable. For R14 eligibility, a verified `agent_started` event may satisfy
an explicit agent-invocation requirement when all existing selector, provider
identity, provenance, and trust requirements are met. An `agent_finished` event
remains valid lifecycle evidence but cannot independently satisfy invocation;
no stop event is required merely to prove that invocation began. This
diagnostic result does not itself write or authenticate a production ledger
event.

### Backend execution

| Required value | Official surface | Local capture |
|---|---|---|
| `session_id` | documented | observed |
| `turn_id` | documented | observed |
| `tool_use_id` | documented | observed and matched across each controlled pair |
| exact predetermined command | documented as Bash `tool_input.command` | observed for both safe commands |
| completion event | documented as `PostToolUse` | observed for both safe commands |
| real underlying command exit status | no stable field documented | explicit field absent from both controlled completions |

The success command's `PreToolUse` and `PostToolUse` records had identical
tool-use identity. The exit-7 command's records did too. Thus exact safe command
invocation and Pre/Post completion correlation are locally observable.

The success command was expected to exit 0 and the failure command was expected
to exit 7. Neither `PostToolUse` record exposed an explicit exit-status field,
so the real underlying status was not locally observable through the hook
payload. Hook-based deterministic backend pass/fail is therefore
**unsupported in the observed runtime**. A `PostToolUse` event proves hook
completion, not command success or failure.

### Sanitized Bash response shapes

| Controlled command | Documented PostToolUse shape | Sanitized live shape | Expected outcome / observed status |
|---|---|---|---|
| success command | `tool_response` is a JSON value; exact members are unspecified | completion observed; explicit exit-status field absent | expected 0 / explicit status absent |
| exit-7 command | `tool_response` is a JSON value; exact members are unspecified | completion observed; explicit exit-status field absent | expected 7 / explicit status absent |

Synthetic fixtures exercise a response object containing string, number, and
exit-status members. They prove only that arbitrary response strings become
length-only shape, that arbitrary output is omitted, and that a genuinely
present numeric `exit_code` or `exit_status` can be retained and compared with
the safe command's expected value. They do not establish a field that was
absent from the real Codex response shape.

### Skill invocation

The official hook event list has no dedicated skill-selection or
skill-invocation event. That capability is **unsupported by the documented hook
surface**. No dedicated skill-selection event was locally observed in the
supplied controlled capture. Because there is no configurable skill hook event
to attach the collector to, exact skill invocation is **not testable through
this hook configuration**. The unmatched `SubagentStop` is not skill evidence
and must not be reclassified as such.

Reading `SKILL.md`, mentioning a skill, selecting it in prose, or observing
generic tools used by a skill cannot be upgraded into `skill_selected`
evidence. A future workflow that requires exact skill invocation needs a
reviewed deterministic wrapper or a future dedicated provider event.

### Operator approval

`PermissionRequest` is officially a tool-permission hook. It is not configured
by this probe and, even if captured, it would not prove that an operator saw
and approved a frozen workflow contract. R14 contract approval still requires
a distinct event bound to the displayed contract hash.

## Coverage gaps and unknowns

- Hosted `WebSearch` is documented as outside normal hook coverage, and some
  specialized tools may opt out.
- One complete agent start/stop identity tuple was locally observed, but an
  unmatched stop also occurred. The cause and broader frequency of unmatched
  or repeated stops remain unknown, and a stop is not independent invocation
  evidence.
- `PostToolUse` means that the tool completed; it does not define command
  success, and neither controlled completion exposed a real exit-status member.
- The hook surface supplies no dedicated skill event.
- Hook payloads do not document a universal provider timestamp. The collector
  timestamp is diagnostic local time, not authenticated provider time.
- Multiple matching hook handlers can run concurrently. Capture arrival order
  is not event order and cannot become ledger sequence.
- Local-file inspection cannot exclude cloud/session policy, and hook hot
  reload remains unknown.
- This probe does not inspect transcripts, provider-global history, MCP
  payload contents, the ROM, Mesen, or Ghidra.

## Controlled live-probe status

| Controlled event | Status |
|---|---|
| fresh `SessionStart` | performed; one record captured |
| one captured `SubagentStart` and `SubagentStop` lifecycle | performed; one matching lifecycle plus one unmatched stop captured |
| `printf 'codex-hook-probe-success\n'` | performed; Pre/Post records correlated, no exit-status field exposed |
| `sh -c 'printf "codex-hook-probe-failure\n" >&2; exit 7'` | performed; Pre/Post records correlated, no exit-status field exposed |
| dedicated skill-selection event | not observed; exact skill invocation remains unavailable |

```text
live hook probe: completed
live hook probe rerun: prohibited in this unit
```

## Setup validation (historical)

The final setup validation is recorded here before the single setup commit:

```text
focused collector tests: pass — 12 tests, exit 0
hook configuration parse: pass — exit 0
raw capture ignore boundary: pass — ignored, no capture present, nothing tracked
git diff --cached --check: pass — exit 0
tracked private/restricted/rendered/binary scans: pass — no prohibited match
gofmt -l .: pass — exit 0, no output
go build ./...: pass — exit 0
go vet ./...: pass — exit 0
go test ./...: pass — exit 0
go build ./cmd/ff6lab: pass — exit 0
go build ./cmd/ff6demo: pass — exit 0
go run ./cmd/ff6lab audit: pass — exit 0, audit clean
go run ./cmd/ff6lab census validate: pass — exit 0, census clean
go mod verify: pass — exit 0, all modules verified
go build -tags gui ./cmd/ff6demo: pass — exit 0
go vet -tags gui ./...: pass — exit 0
go test -tags gui ./...: pass — exit 0
tracked-only candidate checkout gates: pass — full source and GUI set, exit 0
archive verify: not run — FF6_ROM unavailable
```

The first sandboxed invocations of `go build ./...`, `go vet ./...`,
`go test ./...`, and census validation exited 1 because the sandbox denied
writes to the existing external Go caches. Each was rerun with cache access and
exited 0 as recorded above; the initial exit 1 was an environment restriction,
not a source-gate result. The restricted-file searches also use exit 1 to mean
"no match," which is their passing result.

No validation command listed above was a live hook event in the historical
setup task: the definitions were not trusted, loading was not established, no
handler execution was observed, and no capture was produced in that task.
These results are not current-unit validation, and this report does not invent
current-unit gate outcomes. The historical `no capture present` result predates
the controlled probe and later archive.

## Probe-retirement validation

The retirement unit was validated without opening, reading, parsing, copying,
or hashing the raw archive:

```text
.codex/hooks.json absent from working tree: pass — exit 0
tracked active collector registration search: pass — no match, exit 1
preserved archive ignore boundary: pass — ignored, exit 0
preserved archive tracked-file check: pass — not tracked, expected exit 1
focused collector tests: pass — 12 tests, exit 0
production R14 ledger source diff: pass — no diff, exit 0
git diff --check: pass — exit 0
gofmt -l .: pass — exit 0, no output
go test ./...: pass — exit 0
go vet ./...: pass — exit 0
go build ./...: pass — exit 0
go build ./cmd/ff6lab: pass — exit 0
go build ./cmd/ff6demo: pass — exit 0
go run ./cmd/ff6lab audit: pass — exit 0, audit clean
go run ./cmd/ff6lab census validate: pass — exit 0, census clean
go mod verify: pass — exit 0, all modules verified
go build -tags gui ./cmd/ff6demo: pass — exit 0
go vet -tags gui ./...: pass — exit 0
go test -tags gui ./...: pass — exit 0
tracked ROM/audio/state extension scan: pass — exit 0
tracked rendered-asset extension scan: pass — exit 0
archive verify: not run — FF6_ROM unavailable
```

The build commands exited 0 after completing successfully. The sandbox denied
best-effort writes to the external Go module stat cache and emitted non-fatal
warnings; this did not change the recorded command statuses.

## Live-capture preservation boundary

The controlled live probe is complete and must not be rerun. The preserved raw
JSONL remains ignored, private, untracked, and uncommitted and must not be
inspected or copied into tracked documentation. The operator supplied this
archive metadata without requiring inspection of its contents:

```text
archive filename: local_artifacts/codex-hook-probe/events-probe-and-postprobe-20260802.jsonl
total archived records: 414
SHA-256: 86d00525ae7f5d4661ed82e8a3f2c6c5c289fa7c40d40435689b1a77238bbda8
```

The archive is not an isolated 12-record capture. The diagnostic hook remained
active during later Codex sessions, so the archive grew to 414 records. By
count, 402 records accumulated after the original 12-record checkpoint,
assuming the collector continued its documented append-only behavior. Those
post-probe spillover records were not inspected, classified, or used as
workflow evidence. No conclusion may be drawn from them without a new,
separately approved evidence review, which is outside this unit's scope.

The archive metadata above is operator supplied. It must not be recalculated or
augmented from the raw file in this unit. This retirement unit did not open,
read, parse, summarize, copy, or hash the archive, and the supplied digest was
not independently verified.

## Production-adapter readiness

The accepted controlled 12-record evidence is **not sufficient to activate a
production Codex adapter**. The 402 unreviewed spillover records add no accepted
evidence. Agent start identity and one matching lifecycle are locally
observable, but no reviewed adapter writes those events into the production run
ledger. An unmatched stop also shows why finish-only evidence cannot establish
invocation.

Backend identity and Pre/Post correlation are locally observable, but
hook-based deterministic backend pass/fail is unsupported in the observed
runtime because neither controlled completion exposed a real process exit
status. Dedicated skill evidence remains unsupported by the documented hook
surface and unavailable through the tested configuration. Operator approval
and coverage gaps remain separate unsolved boundaries.

A future implementation must keep two responsibilities separate:

- a Codex hook adapter may supply session, turn, agent, tool, and tool-use
  identity evidence; and
- a deterministic backend executor must execute each frozen-contract command
  itself, capture the real process exit status, and write
  `deterministic_backend` events. It must never accept a caller-provided claim
  that the command passed.

Neither production component is implemented by this unit. Diagnostic records
must not be translated into trusted workflow events automatically.
