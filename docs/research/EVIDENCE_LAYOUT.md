# Evidence Layout

Adopted 2026-07-30 (repository-maintenance unit). Applies to EXP-0020
onward; historical experiments keep their recorded paths (do not move
old evidence — provenance outranks tidiness).

Each experiment stores its raw evidence in an immutable per-experiment
directory outside Git:

```text
local_artifacts/experiments/EXP-NNNN/
    experiment.json      # metadata: ROM identity, Mesen version, probe, timestamps
    commands.log         # full bridge command/response transcript
    events.log           # instrumentation output (probe-labeled)
    rom-slices/          # raw dumps (never committed)
    screenshots/
    hashes.sha256        # sha256 of every file above, written at close
```

The tracked experiment record (`docs/experiments/EXP-NNNN-*.md`) cites
paths, sizes, and SHA-256 hashes from `hashes.sha256` — metadata only,
never restricted bytes. Once `hashes.sha256` is written the directory is
frozen; later work goes in a new experiment directory.

Historical note: Sessions 001–003 and EXP-0001..0019 predate this
layout and cite shared files under `mesen/out/` (with archives under
`mesen/out/session003/` where a hash was recorded). Those citations are
preserved as-is; the recorded hashes identify the exact artifacts.
