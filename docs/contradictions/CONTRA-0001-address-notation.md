# CONTRA-0001: Address notation standard (pre-V4 vs V4)

## Status
Resolved (2026-07-29, V4 migration follow-up).

## Claim A
[docs/sessions/01_REVERSE_ENGINEERING_RULES.md](../sessions/01_REVERSE_ENGINEERING_RULES.md)
(pre-V4): addresses may be written as contextual naked hex (`$2EB5`,
`$C10DF3`) provided the document declares which space each form denotes; it
also named a now-retired authority
(`.claude/skills/ff6-reconstruction-skill/SKILL.md`).

## Evidence A
All pre-V4 canonical documents (`docs/sessions/00–08`, SESSION_001/002)
follow this convention and declare it in-file (e.g. 04_MEMORY_MAP.md's
notation header).

## Claim B
[docs/research/ADDRESS_NOTATION.md](../research/ADDRESS_NOTATION.md) (V4):
every address must carry a domain prefix (`ROMCPU:$C10E14`, `WRAM:+$2EB5`);
"a naked hexadecimal number is not acceptable in canonical documentation."

## Evidence B
V4 scaffold standard; identical text in
`.claude/skills/_shared/ADDRESS_SPACES.md`, enforced by every V4 skill.

## Possible cause
Standard evolution: the V4 package introduced a stricter notation without
migrating or marking the pre-V4 rule set.

## Discriminating experiment
None — this is a policy conflict, not an empirical question.

## Resolution
Claim B (V4 domain-prefix notation) **supersedes** Claim A for all new and
substantively updated canonical documentation, effective 2026-07-29.
Pre-V4 documents remain valid historical records under their declared
in-file notation and are **not** mass-edited; whenever such a record is next
substantively updated, its addresses gain domain prefixes at that time.
Rationale: prefix notation removes per-document ambiguity at negligible
cost, while rewriting historical records wholesale would churn evidence
documents without adding information.

## Supersession updates
- 01_REVERSE_ENGINEERING_RULES.md: superseded-notice added to its intro and
  its Address notation section; retired-skill reference replaced with
  pointers to `docs/research/` and `.claude/skills/_shared/`.
- `indexes/CONTRADICTIONS.md`: CONTRA-0001 entry added.
- `dashboards/BLOCKERS.md`: notation blocker cleared.
