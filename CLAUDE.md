# FF6 Reconstruction — Project Constitution

This file is the always-loaded project contract. Detailed procedures,
skills, rules, templates, schemas, and playbooks live under `.claude/`,
`docs/`, `schemas/`, and `manifests/`. This file points to them; it does
not duplicate them.

## Session startup

At the start of every fresh or resumed session:

1. Inspect repository state:
   - `git status --short --branch`
   - `git log --oneline --decorate -10`
   - `git diff`
   - `git diff --cached`

2. Read `docs/checkpoints/LATEST.md` and its linked checkpoint.
3. Read:
   - `dashboards/CURRENT_FOCUS.md`
   - `dashboards/BLOCKERS.md`
   - `dashboards/RESEARCH_QUEUE.md`

4. Read `.claude/README.md`.
5. Resume incomplete work before selecting anything new.

Never discard, reset, overwrite, or reorganize uncommitted work until its
purpose has been established from the checkpoint, records, diff, and Git
history.

## Source-of-truth precedence

Use this precedence:

1. Direct reproducible runtime evidence.
2. Reproducible static ROM evidence.
3. Canonical records:
   - `docs/experiments/`
   - `docs/discoveries/`
   - `docs/sessions/`

4. Tested implementation and golden vectors.
5. Current checkpoint.
6. Schemas and machine-readable manifests.
7. Generated indexes.
8. Dashboards.
9. Hypotheses and external references.
10. Chat memory.

Dashboards and indexes summarize canonical records; they never introduce
new facts. The repository—not the conversation—is project memory.

When two sources disagree, preserve the contradiction and resolve it with
evidence. Do not silently choose the newer-looking document.

## Evidence and confidence

Use only:

- Confirmed
- Strong hypothesis
- Tentative hypothesis
- Unknown

Follow:

- `docs/research/EVIDENCE_STANDARD.md`
- `docs/research/CONFIDENCE_POLICY.md`
- `docs/research/ADDRESS_NOTATION.md`
- `.claude/rules/research.md`

Every research record must separate:

- Observation
- Interpretation
- Alternatives
- Falsifier
- Confidence
- Exact next action

Confidence changes only when supported by cited evidence.

Addresses always include their domain prefix. Never invent semantic names,
fields, formats, boundaries, counts, or behavior. Unverified semantics use
`Unknown`, `Candidate`, or equivalent evidence-safe names.

External walkthroughs, wikis, databases, disassemblies, and guides may
provide leads, but never override evidence from the verified ROM revision.

## Legal and restricted-file boundary

Never commit ROMs, savestates, ROM-derived binaries or images, raw ROM
slices, extracted commercial assets, or absolute personal paths.

Local evidence lives under `local_artifacts/` and `mesen/out/`, both
ignored. Tracked records carry metadata, addresses, sizes, and SHA-256
hashes — never restricted bytes.

See `docs/legal/ASSET_POLICY.md` and
`.claude/skills/_shared/LEGAL_BOUNDARY.md`.

## Experiments

Before operating Mesen, create a falsifiable experiment record under
`docs/experiments/` containing:

- Question
- Starting state
- Controlled variables
- Expected outcomes
- Falsifier
- Required evidence
- Stopping condition

Preserve raw evidence and injected instrumentation before interpretation.

Follow `.claude/skills/_shared/STOPPING_RULES.md`.

Negative and inconclusive results are valid results. Do not broaden a
bounded experiment into adjacent research merely because another system
became visible.

## Content census

After every experiment, gameplay reconnaissance session, static ROM
discovery, implementation unit, or newly observed subsystem, invoke the
`ff6-content-census` skill or `/census-observations`.

Its rule is:

> Observe broadly. Register briefly. Investigate narrowly.

Register newly visible systems, assets, tables, mechanics, events, menus,
audio, graphics, maps, and persistent state without interrupting the active
bounded investigation.

Track reconstruction completeness separately from runtime coverage.

Do not rely on eventual gameplay discovery to guarantee completeness.

## Research workflow

The normal evidence chain is:

```text
experiment
→ discovery
→ Go implementation
→ tests
→ census and ROM-ownership synchronization
```

Implement only Confirmed behavior unless an explicitly isolated prototype
is labeled noncanonical.

Tests should encode evidence through:

- Golden vectors from live captures
- Exact byte-level transformations
- Regression cases for edge behavior
- Differential comparisons where practical

Every implementation must link back to its supporting discovery and
experiments. Every confirmed discovery should identify its implementation
status.

## Research balance

Do not remain indefinitely inside one subsystem.

Prefer work that unlocks:

- Entire tables
- Pointer systems
- Script interpreters
- Shared engines
- Compression formats
- Reusable asset formats
- Complete content families

After several consecutive experiments in one domain, review global coverage
before choosing the next target.

Balance work across:

- Behavior
- Static content
- Graphics
- Audio
- Maps and events
- Menus
- Persistence
- Compatibility and quirks

A generic SNES decoder does not count as FF6-specific reconstruction until
it is linked to an FF6 ROM region, consuming routine, runtime evidence, and
deterministic comparison.

## Documentation synchronization

A completed unit is not complete until all affected canonical and generated
state agrees.

Update as applicable:

- Experiments
- Discoveries
- Sessions
- Functions
- Variables
- Structures
- Content census
- ROM ownership
- Data inventories
- Manifests
- Indexes
- Dashboards
- Tests
- Checkpoint

Regenerate generated files through project tooling rather than editing them
manually when generators exist.

Do not allow a dashboard or index to retain claims contradicted by newer
canonical evidence.

## Quality gates

Before every commit, run:

```bash
gofmt -l .
go build ./...
go vet ./...
go test ./...
go build ./cmd/ff6lab
go run ./cmd/ff6lab audit
```

`gofmt -l .` must produce no output.

Also run applicable:

- Schema validation
- Census validation
- Manifest/index synchronization
- Internal-link checks
- Evidence metadata checks
- ROM-region overlap checks
- Archive verification (`ff6lab archive verify`, needs `FF6_ROM`)
- Restricted-file check (no asset extensions in `git ls-files`)

After changes to `.gitignore`, CI, repository structure, required source
files, manifests, or generated indexes, verify the project from a clean
tracked-only checkout. A working local tree is not sufficient proof that a
fresh clone works.

## Asset policy

The project is **extractor-first**: the public repository carries the
reconstruction (complete named inventories, IDs, statistics, mechanics,
addresses, layouts, hashes, opcode definitions, Go code, extractors,
manifests), and substantial graphics, audio, dialogue, maps, scripts and
raw binary material are generated locally under `local_artifacts/archive/`
from a user-supplied verified ROM.

Names, labels, and complete structured gameplay databases are **permitted
in Git**. The policy may not be used to hide ordinary reconstruction data
behind placeholders; that is enforced by a test.

See `docs/legal/ASSET_POLICY.md`.

## Checkpoints and commits

Use one coherent unit per commit.

After every completed, blocked, or interrupted unit:

1. Synchronize affected records, indexes, manifests, and dashboards.
2. Write a checkpoint under `docs/checkpoints/`.
3. Update `docs/checkpoints/LATEST.md`.
4. Record:
   - What was completed
   - What remains uncertain
   - Current Git state
   - Active instrumentation and evidence
   - Quality-gate results
   - Exact next action

5. Commit the coherent unit when its gates pass.

Do not push unless the operator explicitly asks.

## Interruption and resumption

Use `/resume-session`.

Resume from repository evidence rather than chat history. The resumption
process must:

1. Inspect Git state and uncommitted changes.
2. Read the latest checkpoint.
3. Identify the active unit.
4. Determine what was completed.
5. Continue the exact recorded next action.
6. Avoid duplicating completed work.

Do not begin a new research target while a recoverable active unit remains
unfinished.

### Static analysis

- `/bootstrap-ghidra`
- `/correlate-static-runtime`
- `/export-ghidra-symbols`

Ghidra output is a static lead. Mesen execution and bounded experiments are required before semantic confirmation.
