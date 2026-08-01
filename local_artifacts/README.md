# Local Artifacts

Ignored by Git. Holds everything the asset policy keeps out of the
public repository (`docs/legal/ASSET_POLICY.md`).

## Regenerating the archive

```sh
export FF6_ROM="/path/to/Final Fantasy III (USA).sfc"
ff6lab extract all        # regenerate + refresh manifests/assets.json
ff6lab archive verify     # prove it matches the tracked manifest
ff6lab archive inventory  # list tracked assets (no ROM needed)
```

Extraction refuses any ROM whose SHA-256 is not
`0f51b4fca41b7fd509e4b8f9d543151f68efa5e97b08493e4b2a0c06f5d8d5e2`, and
never overwrites a differing file without `-force`.

## Layout

```text
archive/           regenerated from the ROM; safe to delete
  graphics/  audio/  dialogue/  maps/  animations/  scripts/  raw/
experiments/       PRESERVED evidence — never regenerated, never
                   overwritten by extraction
scenarios/         golden-route milestone states and dumps
backups/           originals of lab settings/SRAM changed by experiments
```

`experiments/` is preserved research evidence: extraction refuses to
write there.
