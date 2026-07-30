# Baseline: `$C10DF3` copy routine and `WRAM+$2EB5`

## Confirmed observations from the prior session

- In one observed battle state, current HP for party slot 0 matched the value at `WRAM+$2EB5`.
- A write breakpoint reached `ROMCPU:$C10E14`.
- The instruction was `STA $2EB5,Y`.
- Observed registers included `A=$002A`, `X=$0000`, and `Y=$0000`.
- The surrounding routine was approximately `ROMCPU:$C10DF3`–`ROMCPU:$C10E66`.
- The routine copied from an address in the `$2E78,X` region to the `$2EB5,Y` region.
- X advanced twice per loop.
- Y advanced through an `ADC #$20` pattern.

## Interpretation

The destination may be a repeated runtime record with a 32-byte stride. The copied field may be current HP.

## Confidence

- Value match: Confirmed for one observed state.
- 32-byte destination stride: Strong hypothesis.
- Ownership as a party-slot record: Tentative hypothesis.
- Semantic function name: Unknown.

## Required next experiments

1. Identify every caller of `ROMCPU:$C10DF3`.
2. Repeat with multiple party members and HP values.
3. Verify destination addresses for each slot.
4. Find initialization and lifetime.
5. Determine whether `$2E78,X` is source HP, a temporary buffer, or another record.
