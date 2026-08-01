# Latest Checkpoint

**[2026-08-01 — EXP-0045: queued work resolves past the gate](2026-08-01-exp0045-queued-work-past-gate.md)**

State: the ATB pause model is now precise. `ROMCPU:$C21124` stops the
scheduler dead — tick `+$3A3E` and every `+$3AB4` gauge frozen across
**1 245 gated frames with zero exceptions** — but **an action already
pending when the gate engages still completes**, clearing `+$3AA0` bit 6
and advancing that slot's `+$3218` by exactly `$0100`.

| Trace | Pending at gate engage | Gated-frame changes |
|---|---|---|
| 1 | slot 6 | 1, at +78 frames |
| 2 | **none** | **0** across 438 gated frames |
| 3 | slot 8 — **predicted** | 1, at +119 frames |

Trace 3 was a prediction test: traces 1 and 2 differ in exactly one
respect, which predicts that engaging the gate while a slot is pending
produces a completion. It did, on a different slot at a different delay.
Arming moved into Lua because bridge round-trips are hundreds of frames
apart — which is exactly why EXP-0044 could not settle this.

**`+$3AA0` bit 6 is the pending-action marker**, resolving a semantics
question EXP-0043 and EXP-0044 both left open. It also **vindicates
EXP-0040**: "queued actions resolved out of issue order" while menus were
open was real system behaviour, not operator error.

Matrix additions: **Item** is a qualifying submenu; the **victory
presentation** is not. Magic, Row and Defend are unreachable from a
Magitek battle and are left unsampled rather than guessed.

Carried forward: the driver of deferred completion is unnamed; multiple
pending slots never co-occurred; `+$3218`'s `$0100` step rests on two
observations; party slots were never pending at a gate engage.

12 evidence artifacts with hashes; no background processes; SRAM virgin.
All gates clean.

Exact next action: **EXP-0046 — the action-queue execution path.**
Read-watch `+$3AA0` and `+$3218` through a gated interval arranged as
trace 3 arranged it, and capture the writing PC. That routine sits
outside the per-frame scheduler.
