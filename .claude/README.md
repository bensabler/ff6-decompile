# Claude Operator Guide

This directory is the agent operating system for the FF6 reconstruction lab.

## Components

### Commands
User-invoked workflow entry points in `.claude/commands/`.

Example:

```text
/orchestrate
/reconstruct-sprite Terra idle
/checkpoint
```

### Skills
Specialist procedures in `.claude/skills/*/SKILL.md`. Claude should load them when the task matches their description. Commands explicitly name important skills so the behavior does not rely solely on automatic selection.

### Agents
Focused subagents in `.claude/agents/`. Agents have constrained roles and should return evidence or implementation results to the main orchestrator.

### Playbooks
Repeatable laboratory procedures in `.claude/playbooks/`. A playbook defines required evidence, steps, stopping conditions, and deliverables.

### Rules
Path-scoped requirements for Go files, research records, manifests, and documentation.

### Templates
Canonical forms for experiments, discoveries, functions, variables, structures, graphics, audio, sessions, and checkpoints.

### Workflows
Multi-stage lifecycle procedures that cross specialists.

## Command categories

### Session control
- `/bootstrap-v4`
- `/resume-session`
- `/checkpoint`
- `/session-summary`

### Automatic research
- `/orchestrate`
- `/weekly-review`
- `/audit-project`

### CPU and behavior
- `/investigate-function`
- `/investigate-variable`
- `/reconstruct-struct`
- `/trace-caller`
- `/trace-dma`
- `/validate-hypothesis`
- `/resolve-contradiction`

### Graphics
- `/capture-graphics`
- `/reconstruct-sprite`
- `/reconstruct-tileset`
- `/reconstruct-background`
- `/recover-palette`
- `/validate-graphics`

### Audio
- `/investigate-audio`
- `/recover-brr`
- `/recover-sequence`
- `/trace-spc-command`
- `/validate-audio`

### Implementation
- `/implement-discovery`
- `/run-quality-gates`
- `/prepare-release`

## Fresh-session sequence

For a new Claude conversation in an already initialized repository:

```text
/resume-session
/orchestrate
```

For the first conversation after installing Version 4:

```text
/bootstrap-v4
/resume-session
/audit-project
/orchestrate
```

Before stopping:

```text
/checkpoint
```

## Arguments

Commands use the text following the slash command as `$ARGUMENTS`.

Examples:

```text
/investigate-variable WRAM+$2EB5
/reconstruct-sprite "Terra battle idle frame"
/investigate-audio "menu cursor sound"
/implement-discovery DISC-BATTLE-0001
```

## Automatic skill invocation

Skills may be selected automatically when their descriptions match the work. Do not depend on that alone for critical procedures. Each command names the skills and playbook it requires.

## Working agreement

Claude must never:

- run exploratory gameplay indefinitely;
- replace evidence with plausible narrative;
- label viewer output as original ROM source without provenance;
- commit extracted commercial assets;
- implement a speculative format as final;
- leave a session without a checkpoint after meaningful work.
