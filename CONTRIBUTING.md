# Contributing

## Required standards

A contribution must be:

- evidence-linked;
- reproducible;
- address-space explicit;
- confidence-rated;
- tested when it changes code;
- free of redistributed ROM-derived bytes.

## Research contribution checklist

- [ ] Stable experiment/discovery ID
- [ ] Exact ROM identity
- [ ] Exact emulator identity
- [ ] Reproduction steps
- [ ] Raw evidence references
- [ ] Observation separated from interpretation
- [ ] Alternatives considered
- [ ] Falsifying outcome documented
- [ ] Canonical indexes updated
- [ ] Checkpoint updated

## Code contribution checklist

- [ ] Idiomatic Go
- [ ] Package has a narrow responsibility
- [ ] Parser bounds checks are explicit
- [ ] Errors are wrapped with context
- [ ] Unit tests added
- [ ] Fuzz test added for untrusted binary format
- [ ] Synthetic fixtures only
- [ ] `gofmt`, tests, and vet pass
- [ ] Architecture and format docs updated

## Naming

Names must reflect confidence:

- `unknownC10DF3` is acceptable during investigation.
- `copyPartyHPRecords` requires evidence.
- `PartySlot.CurrentHP` requires ownership, offset, stride, and lifetime evidence.

Do not encode speculation in durable public names.
