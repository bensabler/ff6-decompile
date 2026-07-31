# Automatic research selection

## Required inputs

- exact target or queued question;
- ROM identity;
- emulator identity when runtime evidence is used;
- current checkpoint;
- evidence and stopping standards.

## Procedure

1. Read state and interrupted work.
2. List candidate questions.
3. Score them.
4. Select one bounded unit.
5. Write experiment.
6. Delegate.
7. Integrate results.
8. Update queue and checkpoint.

## Prioritization rules (depth vs breadth)

1. Finish a bounded active experiment before switching.
2. Register all newly visible systems (`/census-observations`).
3. After at most **three consecutive experiments in one subsystem**,
   review `ff6lab coverage summary` before selecting the next unit.
4. Prefer questions that unlock entire tables, pointer systems,
   script formats, or reusable engines over single-value questions.
5. Alternate among behavior, static content, graphics, audio,
   events/maps, and persistence when scores are close.
6. A generic SNES decoder counts as FF6 coverage only once it is
   connected to an actual FF6 asset and a runtime consumer.
7. Prefer one representative vertical proof, then bulk extraction.
8. Revisit runtime scenarios after bulk extraction to validate the
   interpretation.

The repeating cycle:

```text
focused experiment -> observation registration -> table/ownership
update -> implementation -> tests -> global coverage review -> next
highest-leverage target
```

## Required outputs

- updated canonical record;
- raw evidence references and hashes;
- confidence and alternatives;
- relationship/index updates;
- exact next action;
- checkpoint when the unit stops.
