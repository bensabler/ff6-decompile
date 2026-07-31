---
name: census-observer
description: Register newly observed systems in the content census without expanding the active experiment's scope.
---

# Census Observer

The breadth counterpart to every depth-first experiment. The operating
rule is:

```text
Observe broadly. Register briefly. Investigate narrowly.
```

## When this runs

Before closing any experiment (see `../_shared/STOPPING_RULES.md`),
and on demand via `/census-observations`.

## Procedure

1. Review the experiment's evidence (screenshots, logs, dumps,
   decodes) and ask:
   - What unrelated systems became visible?
   - What new tables, routines, or memory regions must exist for what
     was observed?
   - Were new assets, commands, statuses, flags, menus, or
     transitions observed?
   - Does each belong in `manifests/content-census.json`?
2. For each new observation, add a census entry (or update an
   existing one) with: honest `reconstruction_status` /
   `runtime_status`, evidence citations, unknowns, and one bounded
   `next_action`. New ROM knowledge goes to
   `manifests/rom-regions.json` (never invent ownership).
3. Run `ff6lab census sync` (regenerates indexes + COVERAGE.md and
   validates).
4. **Return to the active research question.** Registration is not
   permission to investigate: do not trace, decode, or implement
   anything newly observed inside the current unit. Promotion above
   OBSERVED/CANDIDATE_LOCATION requires its own experiment.

## Rules

- Visible English names are observations, not proof of table
  structure.
- Never promote hypotheses to Confirmed during registration.
- Follow `../_shared/EVIDENCE_STANDARD.md` for confidence and
  `docs/research/CONTENT_TAXONOMY.md` for domains and status ladders.
- Bulk ROM-derived content stays in `local_artifacts/`
  (`docs/legal/ASSET_POLICY.md`).
