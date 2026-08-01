# Latest Checkpoint

**[2026-08-01 — EXP-0040: Whelk attempt stopped on ATB blocker](2026-08-01-exp0040-whelk-stopped-atb-blocker.md)**

State: **Whelk was attacked correctly and NOT defeated.** Two piloted
GUI attempts from the preserved pre-Whelk state, both under
`Bat.Mode = Wait`. The branch-A premises were verified — **Whelk is two
battle slots** (shell 50000 HP, head 1600 HP; new CEN-BATTLE-0009),
striking the visibly extended head reduces the head and **never** the
shell (six hits, 162-186), head/shell state is **visually
classifiable**, and a **field healing route** exists (Tonic ×4 →
76/105/106). Also corrected: MagiTek sets are character-specific
(leader eight, escorts four — EXP-0039's list was a non-leader's), and
the guard/Esper beat **precedes** Whelk on a clean run (CEN-EVENT-0011
resolved). **Milestone `10-whelk-victory` and B19 remain open.**

Stopped by operator directive on a **hard methodological blocker: the
project has no ATB model** (ACTIVE/WAIT semantics, qualifying submenu
pause states, timer domains, action-queue ordering). All head/shell
timing collected is **menu-pause-contaminated and unusable**;
ACTIVE and WAIT are separate experimental conditions. **Whelk gameplay
must not resume before the ATB research program.**

45 evidence artifacts preserved with hashes and unambiguous savestate
lineage. No background processes running; SRAM still virgin. All gates
clean (gofmt/build/vet/test, `ff6lab audit`, census, archive verify 8/8,
restricted-file scan). Exact next action: **start a new session, resume
from this checkpoint, audit existing battle infrastructure and ATB
evidence, and propose the first bounded ATB experiment before operating
Mesen.**
