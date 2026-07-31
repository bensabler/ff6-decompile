# FF6 Reconstruction — Project Constitution

This file is the always-loaded contract. The detailed operating system
lives under `.claude/` (skills, playbooks, rules, templates) and `docs/`
— this file points there; it does not duplicate them.

## Startup reading (every fresh session)

1. `docs/checkpoints/LATEST.md` → the linked checkpoint.
2. `dashboards/CURRENT_FOCUS.md`, `BLOCKERS.md`, `RESEARCH_QUEUE.md`.
3. `.claude/README.md` for the command/skill/playbook system.

## Source-of-truth precedence

Runtime evidence > canonical records (`docs/sessions/`,
`docs/experiments/`, `docs/discoveries/`) > indexes (`indexes/`) >
dashboards (`dashboards/`) > chat memory. Dashboards summarize records;
they never introduce facts. The repository, not the conversation, is
the project memory.

## Evidence and confidence

Four levels only: Confirmed, Strong hypothesis, Tentative hypothesis,
Unknown (`docs/research/EVIDENCE_STANDARD.md`). Every record separates
observation / interpretation / alternatives / falsification. Confidence
changes only with cited evidence (`docs/research/CONFIDENCE_POLICY.md`).
Addresses always carry a domain prefix (`docs/research/ADDRESS_NOTATION.md`).
Never invent names, fields, formats, or behavior; unproven semantics
carry `Unknown`/`Candidate` markers (`.claude/rules/research.md`).

## Legal and restricted-file boundary

Never commit ROMs, savestates, ROM-derived binaries or images, raw ROM
slices, or absolute personal paths. Local evidence lives under
`local_artifacts/` and `mesen/out/` (both ignored). Tracked records
carry metadata and SHA-256 hashes, not restricted bytes.
See `docs/legal/ASSET_POLICY.md`.

## Experiments

Write the falsifiable plan (question, starting state, expected
outcomes, falsifier) in `docs/experiments/EXP-NNNN-*.md` **before**
operating Mesen. Preserve raw evidence (and injected instrumentation
code) before interpreting. Stop per `.claude/skills/_shared/STOPPING_RULES.md`.
Negative results are results.

## Workflow

experiment → discovery (`docs/discoveries/`) → Go implementation →
tests. Implement only Confirmed behavior; tests encode the evidence
(golden vectors from live captures where possible).

## Quality gates (before every commit)

`gofmt -l .` clean · `go build ./...` · `go vet ./...` ·
`go test ./...` · `go build ./cmd/ff6lab` · `go run ./cmd/ff6lab audit`.

## Checkpoints and commits

One coherent unit per commit. After every completed or blocked unit:
update the affected records/indexes/dashboards, write a checkpoint
(`docs/checkpoints/`, update `LATEST.md`), and record the exact next
action. Do not push unless the operator asks.

## Resuming after interruption

`/resume-session` — it reads the newest checkpoint and dashboards and
resumes interrupted work before selecting anything new.
