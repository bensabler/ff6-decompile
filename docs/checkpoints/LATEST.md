# Latest Checkpoint

**[2026-08-01 — extractor-first archival architecture](2026-08-01-extractor-first-migration.md)**

State: the project is now extractor-first. The asset policy permits
names, labels and complete structured gameplay databases in Git, and
forbids using the policy to hide ordinary reconstruction data; a test
enforces it. All 54 spell names restored (cross-checked against
EXP-0027; the `$FE` narrow-space rendering discrepancy resolved with
EXP-0026 evidence). New `internal/rom`, `internal/game/textenc` and
`internal/extract` packages, 8 deterministic extractors, and
`ff6lab extract` / `ff6lab archive verify|inventory`. A fresh clone plus
the verified ROM regenerates the archive byte-for-byte. Exact next
action: EXP-0037 — inventory the opening's event flags
(`WRAM:+$1E80`/`+$1EA0`/`+$1EC0`).
