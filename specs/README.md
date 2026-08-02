# Shared semantic specifications (Lane 4) — not yet extracted

This directory is the home for **implementation-independent** descriptions of
FF6's formats and behavior: the descriptions that both the native Go port and a
future matching SNES reconstruction implement, and that neither one owns.

Lane defined by [ADR-0002](../docs/decisions/ADR-0002-project-product-lanes.md).

## Current state: the specifications exist, but embedded

They are not missing. They are distributed across two places, both of which are
inside an implementation or inside prose:

- **`docs/discoveries/`** — DISC-0001…0008 describe record layouts, arithmetic
  and bit arrays in prose precise enough to implement from, with citations.
- **Go doc comments and constants** — `internal/game/battledata`,
  `internal/game/atb`, `internal/game/hudfont`, `internal/game/textenc`,
  `internal/platform/snesaddr` and `internal/extract/extractors.go` each carry
  the authoritative form of some format, expressed as Go.

That arrangement is **adequate while there is exactly one implementation**, and
this file does not propose churn for its own sake.

## When it stops being adequate

The moment Lane 3 has any content. A specification that lives inside one
implementation cannot arbitrate between two: when the Go port and the assembly
disagree, the question "which is wrong?" has no independent answer, and the
implementation that happens to hold the doc comment wins by accident.

## Criterion for anything placed here

**Two independent implementations built from the specification alone must agree
byte-for-byte.** A description that cannot support that is a discovery record,
and belongs in `docs/discoveries/`, not here.

Concretely, a specification here should state:

- the exact byte layout, with field offsets and widths;
- the address domain of every address, per `docs/research/ADDRESS_NOTATION.md`;
- the transformation, in terms that do not assume a language or a word size;
- edge and error behavior, including what the original does with inputs it
  never receives in normal play;
- the evidence that establishes each claim, and its confidence;
- what remains **Unknown**, named as such rather than omitted.

The last point matters most. Every format this project has recovered is
partially decoded — DISC-0007's attack records have 6 of 14 bytes unmapped, the
monster record has 20-plus of 32. A specification that quietly omits the
unmapped fields would read as complete and is worse than no specification.

## Not a rewrite target

Extracting a specification does not mean deleting the prose or the Go comment
it came from. It means one of them becomes the citation and the others point at
it. Which direction that goes is a decision for the unit that does the
extraction, not a decision made in advance here.
