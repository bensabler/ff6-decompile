# Asset Policy — extractor-first archival architecture

The project reconstructs the whole game. It does that **without
republishing the game** by splitting the work in two:

- the **public repository** carries the reconstruction — every fact,
  structure, mechanic, identifier and piece of code needed to understand
  and rebuild the systems;
- the **local archive** carries the substantial copyrighted material,
  regenerated on demand from a ROM the operator already owns.

A fresh clone plus the verified ROM reproduces the complete archive
deterministically (`ff6lab extract all`, proven by `ff6lab archive
verify`). Nothing of research value is lost by keeping the assets out of
Git, because the manifests describe every generated file exactly.

## Permitted in Git — track these directly

Track, in full, with real values:

- **Complete named inventories:** spells, items, equipment, monsters,
  characters, commands, formations, maps, shops, treasures, abilities,
  espers, music tracks, sound effects.
- **Identity and structure:** IDs, names, short labels, statistics,
  flags, formulas, mechanics, addresses, pointers, offsets, strides,
  record layouts, dimensions, relationships, counts, hashes, provenance.
- **Complete structured gameplay databases** — a full table of IDs,
  names and numeric fields is reconstruction data, not an asset dump.
- **Event and AI opcode definitions**, including operand shapes and
  semantics.
- **Project-authored behavioral descriptions** written in our own words.
- **Go reconstruction code**, extractors, converters, decoders,
  renderers, and their tests.
- **Asset manifests** describing every expected generated file.
- **Test vectors and fixtures** small enough to be identifying data
  rather than a substitute for the asset.

### Names and labels are permitted

An ordinary name — `Fire`, `Cure 2`, `Ramuh` — is a short factual label
that identifies a record. Names are **required** in the public
inventories: they are how a reader connects a record to observed
behavior.

## Generated locally — never committed

Generate under `local_artifacts/archive/` (ignored):

```text
graphics/   sprites, portraits, tiles, tilemaps, palettes,
            backgrounds, animation frames, font sheets
audio/      BRR samples, decoded PCM, sound effects, music
            sequences, rendered audio
dialogue/   complete dialogue and text extraction
maps/       tilemaps, collision, object and trigger data
animations/ frame sequences and timing
scripts/    event and AI script bodies
raw/        raw ROM slices and undecoded table dumps
```

Also local, under `local_artifacts/` generally: ROM images, savestates,
screenshots, emulator captures, raw binary evidence.

## The rule that closes the loophole

> **This policy may not be used to hide ordinary reconstruction data.**

A placeholder such as

```json
"name": "see local extraction (asset policy)"
```

in a permitted field is a **policy violation**, not policy compliance.
It withholds a short factual label the policy explicitly permits, and it
makes the public inventory useless for its purpose.

This is enforced mechanically:
`internal/census.TestNoPlaceholdersInPermittedNameFields` fails the build
if a permitted name field carries a placeholder marker.

Genuinely unresolved values remain honest: an empty value, `unknown`, or
`Unnamed` are all acceptable, because "not yet established" is a real
research state. What is forbidden is replacing a **known** value with a
policy excuse. Never invent a value to satisfy the test.

## Judgement calls

The distinction is **substantiality**, not format:

| Material | Where | Why |
|---|---|---|
| 54 spell names + stats | Git | identifying data; the inventory's purpose |
| Full 54-record raw byte dump | local | substantial reproduction of the table |
| Font sheet dimensions, palette, hashes | Git | structural metadata |
| The rendered font sheet PNG | local | the asset itself |
| Dialogue IDs, pointers, lengths, hashes, speaker metadata, short project-authored descriptions | Git | identification without reproduction |
| Dialogue text bodies | local | the copyrighted work |
| Opcode definitions and operand semantics | Git | interface description |
| Script bodies | local | substantial content |

When a case is genuinely ambiguous, record the reasoning in the relevant
record rather than resolving it silently in either direction.

## Extraction contract

`ff6lab extract` must, in order:

1. Verify the ROM's SHA-256 against `internal/rom.SupportedSHA256`.
2. Refuse unsupported revisions outright — never extract from an
   unverified image.
3. Generate assets under the ignored archive root.
4. Produce byte-deterministic output.
5. Record every output in `manifests/assets.json` with full provenance.
6. Compare generated hashes against the tracked manifest
   (`ff6lab archive verify`) and report missing, changed, drifted, or
   unknown files.
7. **Never silently overwrite** an existing file whose bytes differ, and
   never write outside the archive root or into preserved experiment
   evidence.

## Enforcement summary

| Gate | Enforces |
|---|---|
| `TestNoPlaceholdersInPermittedNameFields` | no placeholders in permitted fields |
| `TestSpellInventoryIsComplete` | 54 distinct spell names present |
| `TestExtractionIsDeterministic` | regeneration is byte-stable |
| `TestAssetsMatchTrackedManifest` | manifest describes reality |
| `TestManifestEntriesAreComplete` | every entry carries full provenance |
| `ff6lab archive verify` | on-disk archive matches the manifest |
| `git ls-files` restricted-extension scan | no assets committed |

See also `.claude/skills/_shared/LEGAL_BOUNDARY.md`.
