# Proper Fresh-Session Command Sequence

## First session after installing Version 4

Open Claude Code from the repository root and issue these commands one at a time:

```text
/bootstrap-v4
```

Review the migration report. Resolve any blocker it identifies.

```text
/resume-session
```

Confirm that Claude correctly identifies the current work and exact next action.

```text
/audit-project
```

This catches conflicts introduced by merging Version 4 with the existing repository.

```text
/orchestrate
```

Claude now selects and executes one bounded research unit.

When evidence becomes implementation-ready:

```text
/implement-discovery <DISCOVERY-ID>
```

Before ending:

```text
/run-quality-gates
/checkpoint
```

## Every later fresh session

```text
/resume-session
/orchestrate
```

Before ending:

```text
/checkpoint
```

## When Claude was interrupted mid-task

```text
/resume-session
```

Do **not** run `/orchestrate` until Claude completes or formally abandons the interrupted unit.

## When you want a specific target

Skip automatic selection:

```text
/resume-session
/investigate-variable WRAM+$2EB5
/checkpoint
```

or:

```text
/resume-session
/reconstruct-sprite "Terra battle idle"
/checkpoint
```

## Periodic maintenance

Every few substantial sessions:

```text
/weekly-review
/audit-project
/run-quality-gates
/checkpoint
```
