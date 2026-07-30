# Operations Manual

This manual defines how the reverse-engineering lab is operated.

## Operating principle

Claude is an engineer working inside the repository. Claude's transient context is never accepted as durable project memory. Every meaningful result must be written to a canonical file before the session ends.

## Session states

### Clean start
The repository is known, dashboards are current, no interrupted experiment exists, and the emulator state is documented.

### Active experiment
One bounded question, one deterministic starting state, one written experiment plan, and one defined stopping condition.

### Implementation
Evidence has reached the required confidence threshold and is being converted into Go code and tests.

### Interrupted
A session, emulator run, or experiment stopped before completion. Resume the exact recorded next action before selecting new work.

### Maintenance
No new reverse engineering. Documentation, schemas, tests, hashes, indexes, and contradictions are audited.

## Standard daily sequence

```text
/resume-session
/orchestrate
```

The orchestrator must select one bounded work unit. It must not start multiple unrelated investigations.

When the unit reaches a proven implementation boundary:

```text
/implement-discovery <discovery ID>
```

Before stopping:

```text
/checkpoint
```

## First run after installing Version 4

```text
/bootstrap-v4
/resume-session
/audit-project
/orchestrate
```

Do not run `/bootstrap-v4` every day. It is an installation/migration command.

## Safe interruption procedure

1. Stop advancing the emulator.
2. Run `/checkpoint`.
3. Confirm that the checkpoint records:
   - exact question;
   - emulator version;
   - ROM hash;
   - save state or reproducible path;
   - active breakpoints;
   - last observation;
   - files changed;
   - tests run;
   - next exact action.
4. Commit or stash work according to the Git policy.
5. End Claude.

## Experiment discipline

Every experiment must have:

- a stable ID;
- a question;
- starting state;
- controlled variables;
- observation method;
- expected outcomes;
- falsifying outcome;
- raw evidence paths;
- result;
- confidence;
- next action.

Playing the game without a written question is reconnaissance, not an experiment.

## Git policy

- `main` should remain buildable.
- Use narrow branches such as:
  - `research/hp-copy-c10df3`
  - `graphics/menu-cursor`
  - `audio/brr-decoder`
  - `tooling/asset-manifest`
- One commit should represent one coherent research or implementation unit.
- Never commit ROM-derived bytes.
- Before merging:
  - `gofmt -w .`
  - `go test ./...`
  - `go vet ./...`
  - manifest validation
  - documentation update
  - provenance audit

## Failure recovery

If Claude loses context:

```text
/resume-session
```

If the emulator state is lost:

1. Read the current experiment file.
2. Recreate the deterministic starting state.
3. Compare ROM and save-state hashes.
4. Reapply documented breakpoints.
5. Repeat the last confirmed step.
6. Never infer missing observations from memory.

If documentation conflicts:

```text
/resolve-contradiction <topic>
```

No implementation may proceed while a load-bearing contradiction remains unresolved.
