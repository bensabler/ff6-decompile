# Latest Checkpoint

**[2026-08-01 — EXP-0042: battle-entry configuration sampling](2026-08-01-exp0042-battle-entry-config-sampling.md)**

State: the ATB program has its **staging rule**. Configuration sampling
across a battle entry is **mixed, and the split falls exactly where it
matters**. `ROMCPU:$C22472` reads the two config bytes **once each** at
the entry frame and decomposes them:

| Setting | Battle-local cell | Value |
|---|---|---|
| Bat.Mode | `WRAM:+$3A8F` | `01` = Wait, `00` = Active |
| Bat.Speed | `WRAM:+$3A90` | `255 − 24 × speed` (Fast `$FF` … Slow `$87`) |
| Cmd.Set | `WRAM:+$2F2E` | cleared when Window |
| Gauge | `WRAM:+$2021` | cleared when Off |
| `+$1D4E` bits 0-2 | `WRAM:+$2F34` | at `$C10FF7` |

Neither Bat.Mode nor Bat.Speed is re-read for timing during the battle.
`Msg.Speed` and `Cursor` **are** read live, by `$C198AC` (message-delay
table `ROMCPU:$C19872`) and `$C159D6` (clears the `$5C`-byte
cursor-memory block at `+$890F` when Cursor = Reset — the mechanism
behind EXP-0040's `Cursor = Memory` observation).

The `$C22472` arithmetic was decoded statically and then used to
**predict** `$3A8F`/`$3A90` for a second, differently-configured run;
both matched exactly. Two live encounters, both formation 14,
reproducing EXP-0038.

**Consequence:** ACTIVE/WAIT and Battle Speed must be established
**before battle entry**, or injected directly at `+$3A8F`/`+$3A90`.

**The hard ATB blocker remains open** — no timer domain, pause condition
or queue semantics is known. But there is now a concrete way in:
`+$3A90`'s consumer is unlocated and is the sharpest lead into ATB rate.
Registered as a Tentative hypothesis, not a conclusion.

Whelk not resumed; its savestates not reloaded. 8 evidence artifacts
preserved with hashes (including two **live battle savestates**); no
background processes; SRAM virgin. All gates clean (gofmt/build/vet/test,
`ff6lab audit`, census sync 63 entries, restricted-file scan).

Exact next action: **EXP-0043 — locate the consumer of `WRAM:+$3A90` and
the ATB gauges.** Read-watch `+$3A8F`–`+$3A90` across a battle and follow
the reader; expect convergence with the eight undumped callees of
`ROMCPU:$C101FB`. Start from
`local_artifacts/experiments/EXP-0042/in-battle-formation14.mss`.
