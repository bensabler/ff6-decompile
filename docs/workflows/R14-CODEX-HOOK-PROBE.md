# R14 Codex hook capability probe

This report characterizes the Codex hook surface needed by the provider-neutral
R14 run-event ledger. It is a diagnostic capability probe, not a production
adapter. Nothing in this unit writes to the workflow ledger, asserts trusted
`provider_hook` provenance, implements contract approval, or changes
reconciliation.

```text
probe date: 2026-08-02
foundation commit: df3a0c6aa4e8c078b5795bba34490a1b3b63c599
working branch: codex/r14-codex-hook-probe
installed CLI: codex-cli 0.146.0-alpha.9.2
installed app: 26.727.51351 (build 6119)
repository hook definitions recognized or loaded in this task: not established
repository hook definitions trusted in this task: no
repository hook handler execution observed in this task: no
sanitized live captures recovered: none
live hook probe: not run
```

The setup requires operator review of the exact hook definitions and a fresh
session. A fresh session is necessary even if hook loading can occur without a
full app restart, because this task's `SessionStart` has already passed.

## Evidence vocabulary

This report keeps these classifications distinct:

- **officially documented**: stated by current public Codex documentation;
- **locally observed**: returned by the installed runtime or present in a real
  sanitized hook capture;
- **inferred**: a bounded conclusion from documented and local facts;
- **not observed**: the live probe has not produced the evidence;
- **not testable**: the current documented/configurable surface offers no
  corresponding probe point;
- **unsupported**: the documented surface explicitly does not provide the
  capability.

Synthetic test fixtures establish sanitizer behavior only. They are never
classified as live hook evidence.

## Recovered state and completed phases

Recovery found the branch at the required foundation commit with three
untracked setup files and no probe documentation or sanitized output. No
pre-existing repository hook configuration was overwritten.

The completed work before recovery was:

1. foundation branch, commit, worktree, and remote-ref verification;
2. read-only runtime, configuration-layer, trust, and hook-capability review;
3. creation of the diagnostic collector, its synthetic tests, and the
   repository-local hook configuration;
4. an initial focused synthetic test run.

The hook definitions were not established as loaded and were not trusted in
this already-running task.
No live `SessionStart`, subagent, Bash success, Bash failure, or skill-selection
event was captured. The first incomplete step was therefore to finish this
report, validate the setup, commit it once, and stop for operator action.

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
| Hook hot reload | unknown | The documented app-server API lists `hooks/list`, but no hook reload method. The current task did not establish hot reload. |
| Full app restart | inferred fallback | A restart is conservative if newly trusted definitions still appear skipped, but the documentation does not establish that it is always required. |

Current references are the official [Hooks](https://developers.openai.com/codex/hooks),
[Advanced configuration](https://developers.openai.com/codex/config-advanced),
[Managed configuration](https://developers.openai.com/codex/enterprise/managed-configuration),
and [App Server](https://developers.openai.com/codex/app-server) documentation.

## Diagnostic setup

The tracked setup consists only of:

```text
.codex/hooks.json
.codex/hooks/codex_hook_probe.py
.codex/hooks/test_codex_hook_probe.py
docs/workflows/R14-CODEX-HOOK-PROBE.md
```

`.codex/hooks.json` registers observational command hooks for only the required
events:

| Hook event | Matcher | Purpose |
|---|---|---|
| `SessionStart` | `startup|resume` | session payload shape |
| `SubagentStart` | `.*` | subagent start identity shape |
| `SubagentStop` | `.*` | subagent stop identity shape |
| `PreToolUse` | `^Bash$` | Bash identity and exact controlled input |
| `PostToolUse` | `^Bash$` | Bash completion and response shape |

Each handler invokes the same repository-local collector, returns success, and
provides no hook decision or additional context. Matching handlers may run
concurrently; the collector makes no inference from arrival order.

The collector reads at most one 1 MiB UTF-8 JSON object from standard input and
appends one sanitized record beneath the already ignored path:

```text
local_artifacts/codex-hook-probe/events.jsonl
```

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

## Official payload matrix versus local evidence

Codex documents these common hook fields: `session_id`, `transcript_path`,
`cwd`, `hook_event_name`, and `model`. Turn-scoped events add `turn_id`. The
required events also document `permission_mode`. Transcript fields are never
read by this probe.

| Event | Additional officially documented fields | Locally observed fields | Status |
|---|---|---|---|
| `SessionStart` | `source`; no documented `turn_id` | none | not observed — live probe not run |
| `SubagentStart` | `turn_id`, `agent_id`, `agent_type` | none | not observed — live probe not run |
| `SubagentStop` | `turn_id`, `agent_id`, `agent_type`, `agent_transcript_path`, `stop_hook_active`, `last_assistant_message` | none | not observed — live probe not run |
| `PreToolUse` for Bash | `turn_id`, `tool_name`, `tool_use_id`, `tool_input`; Bash input documents `command` | none | not observed — live probe not run |
| `PostToolUse` for Bash | the PreToolUse identity/input fields plus `tool_response` | none | not observed — live probe not run |

Documentation states that `PostToolUse` runs after a tool completes, including
when the underlying command exits nonzero. It does not document a stable Bash
exit-status field inside `tool_response`. Therefore generic hook completion is
not evidence that the command succeeded.

## Required capability questions

### Agent invocation

| Required value | Official surface | Local capture |
|---|---|---|
| `session_id` | documented | not observed |
| `turn_id` | documented for subagent events | not observed |
| `agent_id` | documented | not observed |
| `agent_type` | documented | not observed |
| `SubagentStart` or `SubagentStop` identity | documented by `hook_event_name` plus agent fields | not observed |

These values were **not actually observed together**. One subagent lifecycle
occurred while preparing this setup, but loading was not established and no
handler execution was observed; it produced no capture and earns no evidence credit. Agent invocation is
plausibly observable from the official contract but is not yet reliable in
this environment.

### Backend execution

| Required value | Official surface | Local capture |
|---|---|---|
| `session_id` | documented | not observed |
| `turn_id` | documented | not observed |
| `tool_use_id` | documented | not observed |
| exact predetermined command | documented as Bash `tool_input.command` | not observed |
| completion event | documented as `PostToolUse` | not observed |
| real underlying command exit status | no stable field documented | not observed |

These values were **not actually observed together**. The real command exit
status is currently unknown. Until a live success capture contains `0` and the
live exit-7 capture contains `7` in an explicit provider field, the hook
surface cannot reliably satisfy R14 deterministic-backend evidence.

### Sanitized Bash response shapes

| Controlled command | Documented PostToolUse shape | Sanitized live shape | Real exit status |
|---|---|---|---|
| success command | `tool_response` is a JSON value; exact members are unspecified | not observed | unknown |
| exit-7 command | `tool_response` is a JSON value; exact members are unspecified | not observed | unknown |

Synthetic fixtures exercise a response object containing string, number, and
exit-status members. They prove only that arbitrary response strings become
length-only shape, that arbitrary output is omitted, and that a genuinely present numeric
`exit_code` or `exit_status` can be retained and compared with the safe
command's expected value. They do not establish the real Codex response shape.

### Skill invocation

The official hook event list has no dedicated skill-selection or
skill-invocation event. That capability is **unsupported by the documented hook
surface**. No explicit skill attempt was captured locally, so a dedicated
runtime event is also **not observed**. Because there is no configurable skill
hook event to attach the collector to, the stronger runtime question is **not
testable through this hook configuration**.

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
- Agent tooling has provider-specific lifecycle hooks, but the exact local
  identity tuple remains unobserved.
- `PostToolUse` means that the tool completed; it does not define command
  success or document a real exit-status member.
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
| fresh `SessionStart` | not performed |
| one captured `SubagentStart` and `SubagentStop` lifecycle | not performed |
| `printf 'codex-hook-probe-success\n'` | not performed |
| `sh -c 'printf "codex-hook-probe-failure\n" >&2; exit 7'` | not performed |
| one explicit available-skill selection attempt | not performed |

```text
live hook probe: not run
```

## Validation

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

No validation command is a live hook event: the definitions were not trusted,
loading was not established, no handler execution was observed, and no capture
was produced in this task.

## Exact operator action required

After checking out the setup commit, perform the following from the repository
root:

1. Start the bundled CLI against the repository:

   ```sh
   /Applications/ChatGPT.app/Contents/Resources/codex -C "$(git rev-parse --show-toplevel)"
   ```

2. Open `/hooks`. Review and trust the exact repository command-hook
   definition or definitions for the five configured events. Repository trust
   by itself is not enough.
3. Exit that startup task. An untrusted `SessionStart` may already have been
   skipped before the definitions were approved.
4. Start a fresh Codex app task rooted at this repository. If the definitions
   still appear skipped or absent, restart the app as a conservative fallback,
   reopen the repository, and check `/hooks` again.
5. In that fresh task only, allow the automatic startup capture, ask Codex to
   spawn one available subagent and let it finish, then ask Codex to execute
   each command exactly once:

   ```sh
   printf 'codex-hook-probe-success\n'
   sh -c 'printf "codex-hook-probe-failure\n" >&2; exit 7'
   ```

6. Submit this exact harmless explicit skill-selection prompt once:

   ```text
   $anthropic-skills:learn Acknowledge this explicit skill selection in one sentence. Do not call tools.
   ```

   This checks only whether any dedicated event appears. The configured hook
   surface has no skill event, so absence remains `not observed` and `not
   testable`; it is not additional proof beyond the documented unsupported
   surface. Do not infer invocation from the prompt, a skill-file read, or
   generic tool events.
7. End the task before inspecting `local_artifacts/codex-hook-probe/events.jsonl`
   from an ordinary terminal, so inspection commands do not add probe events.
   Keep the raw diagnostic file ignored and uncommitted.
8. Review only the sanitized field/type/length/digest records. If any private
   value appears, stop and quarantine the diagnostic file rather than copying
   it into tracked documentation.

## Production-adapter readiness

The current evidence is **not sufficient to design or activate the production
Codex adapter**. The official contract makes agent evidence promising, but the
required local identity tuple has not been captured. Backend evidence remains
blocked on a real, explicit command exit status. Dedicated skill evidence is
unsupported by the documented hook surface. Operator approval and coverage
gaps remain separate unsolved boundaries.

After the operator-controlled live capture, a new review may characterize the
actual sanitized wire shapes and decide whether a narrow adapter is feasible.
That later review must not translate diagnostic records into trusted workflow
events automatically.
