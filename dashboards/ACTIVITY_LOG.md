# Activity Log

- 2026-08-01 (headless) — **EXP-0044: the ACTIVE/WAIT pause matrix; the
  ATB blocker is discharged.** `WRAM:+$2F41` is the **battle submenu
  flag** — resting `$00`, cleared per-frame at `ROMCPU:$C17A92` (`STZ`),
  raised at `ROMCPU:$C17C01` (`INC`) when a qualifying submenu opens,
  with a second clear at `ROMCPU:$C14434`. `ROMCPU:$C21124` ANDs it with
  the Wait flag `+$3A8F` and skips the **entire** per-frame battle
  update. All four located domains — tick `+$3A3E`, gauges `+$3AB4`,
  flags `+$3AA0`, accumulator `+$3218` — froze and resumed **together**
  across nine trials; no independent clock was found among them.
  **The pause is narrower than the folk model**: the main battle command
  window is on screen for most of a WAIT battle and does **not** pause,
  and action animations do **not** pause either (`$2F41` reads `00` in
  both). Only the ability list and target selection paused. Proven by an
  **in-place one-variable flip** of `+$3A8F` inside a single savestate —
  same submenu, `$2F41` = `01` throughout, frozen at Wait and resuming at
  Active — and cross-checked against **genuinely configured WAIT**
  (falsifier 2, tick stuck at `$17C7`). One unresolved signal recorded
  rather than smoothed over: a settling transient just after the gate
  engages, which this sampling resolution cannot separate from the last
  un-gated frame. Six matrix rows honestly marked `not sampled`.
  **Consequence: EXP-0040's Whelk timing can now be scoped rather than
  dismissed**, and Whelk is no longer blocked by an absent model. New
  CEN-BATTLE-0012; `/battle-baseline` added first as its own commit.
- 2026-08-01 (headless) — **EXP-0043: the ATB layer is located.** Gauge
  array **`WRAM:+$3AB4`** (10 entries, **stride 2**, 16-bit; party 0-3,
  enemies 4-9), advanced by `$3AC8,X >> 1` at
  `ROMCPU:$C21195-$C2119D` — observed climbing by exactly `$4E` =
  `$9C >> 1` and wrapping. Per-slot increment **`WRAM:+$3AC8`** built at
  `ROMCPU:$C209E0` from the slot Speed at `+$3B19,X`. Scheduler flags
  **`WRAM:+$3AA0`**, battle tick counter **`WRAM:+$3A3E`** (one per
  non-gated frame). `+$3A8F`'s consumer is **`ROMCPU:$C21124`** —
  `LDA $2F41 / AND $3A8F / BNE` gates the *entire* per-frame battle
  update, so that instruction is **where ACTIVE and WAIT diverge**;
  `$2F41` is the untested other half. **Battle Speed scales enemy gauges
  only**: the `CPX #$08 / BCC` branch at `$C209F6` skips the `$3A90`
  multiply for party slots, and measurement agrees — party increments
  byte-identical at Bat.Speed 3 and 6 (318/330/336) while enemy
  increments went 240 → 156. Three instruments converged (read watch,
  static decode, live sampling); no falsifier fired. **Caution recorded:
  the ATB family is stride 2, not the `$14` of the HP/stat family** — so
  DISC-0001's unified layout governs slot assignment, not stride. New
  CEN-BATTLE-0011. Blocker is now narrow rather than total.
- 2026-08-01 (headless) — **EXP-0042: configuration is sampled at battle
  entry; first battle-local timing cells found.** Answered **mixed**, and
  the split falls exactly where it matters. `ROMCPU:$C22472` reads the
  two config bytes **once each** at the battle-entry frame and
  **decomposes** them: Bat.Mode (bit 3) → `WRAM:+$3A8F` (`01` = Wait,
  `00` = Active), Bat.Speed (bits 0-2) → `WRAM:+$3A90` = `255 − 24 ×
  speed` (Fast `$FF` … Slow `$87`), Cmd.Set → `+$2F2E`, Gauge →
  `+$2021`, and `+$1D4E` bits 0-2 → `+$2F34` at `$C10FF7`. Neither
  Bat.Mode nor Bat.Speed is re-read for timing during the battle.
  `Msg.Speed` and `Cursor` **are** read live, by `$C198AC` (delay table
  `ROMCPU:$C19872`) and `$C159D6` (clears the `$5C`-byte cursor-memory
  block at `+$890F` when Cursor = Reset — the mechanism behind
  EXP-0040's `Cursor = Memory` observation). Both derived values were
  **predicted from the disassembly before the second run** and matched
  exactly, across two configurations and two live encounters (formation
  14, twice, reproducing EXP-0038). Establishes the ATB program's
  **staging rule**: ACTIVE/WAIT and Battle Speed must be set before
  battle entry, or injected at `+$3A8F`/`+$3A90`. `+$3A90`'s consumer is
  unlocated and is the sharpest lead into ATB rate — registered as a
  Tentative hypothesis, not a conclusion. New CEN-BATTLE-0010; new
  `watchreads`/`watchdump` in `probes/common.lua`. Blocker still open.
- 2026-08-01 (headless) — **ATB program opened; EXP-0041: battle
  configuration located and decoded.** All nine Config settings mapped
  to bits across `WRAM:+$1D4D`, `+$1D4E` and `+$1D54` — the block is
  **not contiguous**, which a three-byte read window initially hid
  (Controller looked like a null result until the read was widened).
  `Bat.Mode` = bit 3 of `+$1D4D`, `Bat.Speed` = bits 0-2; both speed
  fields swept to their clamps (0-5, displayed 1-6). The Config screen
  marks the **selected** option with tile attribute `$20` and the
  unselected with `$28` — the inverse of the intuitive reading, and the
  precise cause of EXP-0040's `Bat.Mode` misread, whose correction is
  now **independently confirmed from memory**. Configuration is **not**
  SRAM-backed before a save. Two enabling changes shipped first: the
  battle-configuration fingerprint became a required, audited record
  field (from EXP-0041 onward), and `ff6lab state` now reads work RAM
  and save RAM out of preserved `.mss` files with no emulator — trial 0
  used it to pull the `+$1D4E` Cursor candidate out of EXP-0040's
  savestate pair *before Mesen was launched*, and the live falsifier run
  reproduced it exactly. New CEN-MENU-0007. Blocker remains open: no
  timer domain or queue semantics known yet. 17 artifacts preserved with
  hashes; no background processes; SRAM still virgin.
- 2026-08-01 (piloted) — **EXP-0040: Whelk victory attempted, NOT
  achieved; stopped on an ATB methodological blocker.** Two attempts
  from the preserved pre-Whelk state, both under `Bat.Mode = Wait`.
  Confirmed: **Whelk occupies two battle slots** — shell slot 4
  = 50000/50000 HP, head slot 5 = 1600/1600 HP (new CEN-BATTLE-0009);
  **head-only targeting is correct**, six measured hits (162-186) all
  reduced the head while the shell's 50000 never moved, and no shell
  strike occurred; head/shell state is **visually classifiable**; a
  **field healing route** exists before contact (Tonic ×4 →
  76/105/106); MagiTek sets are character-specific with **eight**
  abilities for the leader and four for the escorts, **correcting
  EXP-0039's four-entry record**; and the guard/Esper beat
  (CEN-EVENT-0011) fires at `(2A,07)` **before** Whelk on a clean run,
  **resolving** its sequence position and correcting the earlier
  "post-Whelk" reading. Two operator errors were caught and corrected
  in-session: a misread of the Config screen (Bat.Mode was **already
  `Wait`**; the hand cursor is not a selection marker — the only
  config change made was `Cursor: Reset → Memory`), and blind
  multi-press batches that desynchronized from menu/dialogue
  transitions. **Blocker recorded:** the project has no ATB model
  (ACTIVE/WAIT semantics, qualifying submenu pause states, timer
  domains, action-queue ordering), so the battle could not be operated
  reliably and **all head/shell timing collected is menu-pause
  contaminated and unusable**. Milestone `10-whelk-victory` and B19
  remain open; Whelk gameplay must not resume before the ATB research
  program. 45 artifacts preserved; Mesen terminated; SRAM still virgin.

- 2026-08-01 (autonomous, breadth) — **EXP-0039: Whelk reached.** A
  deliberate breadth-first GUI pass explored the mines from milestone
  05 to the Whelk battle. Findings registered without decoding: the
  mines↔exterior transition is **bidirectional**; the entry corridor
  is linear with one turn plus a one-tile dead-end nub; a third
  encounter (formation 14 again) fired at `(26,15)` rather than
  `(26,0B)`, **refuting fixed-tile encounter triggering**; the
  **Magitek ability list** was captured (Fire Beam, Bolt Beam, Ice
  Beam, Heal Force) and battle commands proved character-specific;
  **Whelk** is formation 432, contact-triggered from `(2A,07)` after a
  scripted beat, with its introduction dialogue and **shell
  counterattack** observed. The first attempt **ended in defeat**,
  capturing the defeat flow (CEN-BATTLE-0007). New bounded question:
  the formation record's monster-id field needs a high-bit extension
  decode (Whelk's ids read as record 0, implausibly). New census
  entries CEN-EVENT-0010/0011; pre-Whelk savestate preserved. Next:
  EXP-0040, the branch-A victory attempt.

- 2026-08-01 (autonomous) — **EXP-0038: milestone `06-random-encounter`
  established.** Two scheduled power-on runs reproduce the mines random
  encounter at frame 51 307 (formation 14 = three of monster record 19,
  leg 19, near tile `(26,0B)`) with byte-identical milestone WRAM and
  identical event-flag timelines. Corridor
  `(26,1C)→(26,0B)→(28,0B)→(28,09)` mapped by GUI recon; scripted event
  at ~`(2A,09)` registered (CEN-EVENT-0009) and not entered; seventh
  verification of EXP-0030's formation table; mines yields two
  formations (14, 44); traversal and encounters write **no** event
  flags. Two probe defects (`shot`/`mkstate` dropped when EXP-0037 was
  derived from EXP-0036) found and fixed; artifact writing is now
  `pcall`-guarded. A background-task audit found and cleaned three
  processes (one orphaned monitor, one hung run, one redundant
  monitor). Next: EXP-0039, a deliberate **breadth-first** recon pass.

- 2026-08-01 (autonomous) — **EXP-0037: opening event-flag inventory
  complete.** All writes to `$1E80`/`$1EA0`/`$1EC0` captured across
  the scheduled route (one visible GUI pass + two headless evidence
  runs): 20 flags touched — 11 latched story flags, 4 transient, 5
  engine working bits — 162 value-changing writes, byte-identical
  across runs on every channel, final WRAM = milestone 05 (now five
  byte-identical runs). Every writer PC statically decoded; the
  16-handler script-command family found over eight bases; event
  interpreter anchored at candidate `$C09B5C` (CEN-EVENT-0001);
  GUI/testrunner parity verified for this schedule; `$1EA5`'s
  `$00→$01→$05→$0D` reproduced as `EVF-1EA0-$28/$2A/$2B`. New:
  DISC-0008, `internal/game/eventflags` (+tests), ROM-0027..0032,
  `data/scenarios/opening-event-flags.json`. Next: EXP-0038 (mines
  traversal to milestone 06).

- 2026-08-01 (audit) — Project audit focused on route tables vs raw
  recon logs vs `internal/scenario`: recon command transcript archived
  into EXP-0035 evidence (was only in mutable mesen/out), "A x3" and a
  stale ($1E,$27) trigger coordinate corrected against the raw log,
  probe-sync guard extended to timeouts/durations, CONTRA-0002
  propagated into the scenario manifest/record. One substantive item
  reported for approval: the "contact-triggered" battle-5 claim rests
  on wrong-tile evidence.

- 2026-08-01 (autonomous) — CONTRA-0002: **`WRAM:+$1EA5` is byte 5 of
  the event-flag bit array at `+$1EA0`** — both the EXP-0035 "map id"
  and EXP-0036 "map-load target" readings refuted by a static decode
  of the writer (`ORA $C0BAFC,X / STA $1EA0,Y`; decoder $BAED; values
  accumulate bits $00→$01→$05→$0D). Event-flag system registered
  (CEN-EVENT-0008); dependent route identifiers renamed and unfrozen.

- 2026-08-01 (autonomous) — Units 36/36b / EXP-0036: **milestone
  `05-mines-entry` established** — 17-leg state-driven route
  controller (position targets, battle edges, per-leg timeouts that
  name the earliest divergent leg); three power-on runs byte-identical
  at ($26,$1C) in the mines. Battle 5 = formation 84 {27,27,0,0}
  (ROM-verified 0x0F66EC); new pre-Whelk monster record 27 (115 HP /
  30 MP). Found and corrected an EXP-0035 condensed-table omission (a
  dropped `up` leg) after two self-naming leg timeouts. Go model +
  tests in `internal/scenario/route` with a Lua/Go probe-sync guard.

- 2026-08-01 (autonomous) — Unit 35 / EXP-0035: route from milestone
  04 to the mines interior mapped leg-by-leg via live position reads;
  **player tile bytes `+$00AF`/`+$00B0` located** (blind WRAM diff);
  fifth scripted battle registered (CEN-EVENT-0007); milestone 05
  deliberately not claimed from interactive recon.

- 2026-07-31 (autonomous) — Units 31-34 / EXP-0031..0034: **golden
  route power-on → free movement, deterministic** — title requires
  edge-toggled start+a (holds ignored; $4219 verified latched); real
  opening auto-runs to an input-waiting box; exactly four scripted
  battles (formations 2, 1, 2, 41 — every staged record matching the
  ROM table); battle-1 rewards 32 EXP / 96 GP with field HP/MP
  writeback ($C2496E/$C24979 → +$1609/+$160D); milestones 00-04, WRAM
  byte-identical across paired power-on runs. Lab controls set:
  RamPowerOnState=AllZeros, virgin SRAM (originals backed up).
  Corrected EXP-0032's formation-2 misread (two of monster 0, not one
  monster 12). SCN-0001 program record + manifest created.

- 2026-07-31 (autonomous) — Unit 30 / EXP-0030: **the formation table
  is located and verified — ROMFILE:0x0F6200, 15-byte records
  (id x15 via the ASLx4-minus-id idiom), monster ids at bytes +2..+7,
  staged to +$3F44 by $C2315C; flags table at 0x0F5900 (4-byte).**
  Formation 44 = the mines encounter = monsters **{19, 77}** —
  correcting EXP-0028's coincidental record-78 match; powers {13,19}
  reconcile exactly with EXP-0018's live captures (monster 19: power
  13, HP 24). Encounter-roll output identified at +$11E0 (producer =
  next hop). MONSTER domain now has its three core tables (stats,
  formations, flags) located and cross-verified.

- 2026-07-31 (autonomous) — Unit 29 / EXP-0029 (static): **the
  formation loader chain is decoded.** Per-slot loader $C22C30
  (A=id, Y=slot) scales the monster id x4/x8/x32 into the bank-$CF
  tables ($CF8400 = per-monster 16-bit attribute table -> $3254,Y);
  the sole caller $C22F22 loops six enemy slots reading ids from
  **WRAM:+$3F46 (the staged formation id list, $FF sentinel)** with
  an event-battle alternate source ($0206). The ROM formation record
  is one write-watch away (+$3F46's producer). The census overlap
  audit caught and corrected a region-boundary mistake in the same
  unit.

- 2026-07-31 (autonomous) — Unit 28 / EXP-0028: **the monster database
  is located — ROMFILE:0x0F0000, 32-byte records, Confirmed.** A mines
  free-walk triggered a live encounter; the battle-init write-watch
  burst named the populate routine ($C22CA3-$C22D8C), which reads
  `LDA $CF0000+off,X` straight into the Confirmed battle arrays
  (+$08 HP, +$0A MP, +$01 battle power, +$02/+$10 staging, 13 mapped
  offsets total). Stride 32 Confirmed by a dual cross-check: records
  77/78 hold powers 19/13 — exactly the EXP-0018 live values — so the
  mines formation is ids 77+78. EXP-0004's question 19b closed
  ($C22CCE = the max-HP store). Five aux bank-$CF tables registered
  as leads; random-encounter triggering registered (CEN-WORLD-0006).

- 2026-07-31 (resume) — Unit 27 / EXP-0027: **spell database
  extraction complete.** Name-table boundary Confirmed at exactly 54
  entries (class icons $E9x24/$EAx21/$E8x9; the menu bullet IS the
  record's class icon); the **esper name table** found at 0x26F6E1
  (27x8, clean stride decode) with a stride-~10 ability-name table
  following (candidate). Battle-menu cost anchor unreachable (party
  at 0 HP; the wait produced the 'Annihilated' defeat flow + on-screen
  formation names Were-Rat/Repo Man — both registered); the **field
  Cure cast** anchored the cost behaviorally: '5 MP Needed' gate,
  heal 34→77/77, deduction 24→19 — **record byte 5 = MP cost
  Confirmed** for id 45. Single-byte WRAM diffs located the field
  character-record block (HP +$1609, MP +$160D slot 0). All 54
  records extracted locally; ids+numbers mirrored into
  data/census/spells.json. First EXTRACTED_COMPLETE census entry.

- 2026-07-31 (census unit) — **Content census system established**
  (5 commits): 12-domain taxonomy with separate reconstruction and
  runtime status ladders; content-census + ROM-ownership manifests
  with schemas; `internal/census` validation/coverage/gap tooling
  wired into `ff6lab` (census validate|sync, coverage
  summary|gaps|domain, rom gaps) and the audit; 19 data-family
  inventories. **EXP-0025** registered 25 opening-sequence systems
  from a headless screenshot sweep (field menu, dialogue, maps,
  collision, '?????' pre-naming state, sprites, formations).
  **EXP-0026** censused the magic system: field-menu navigation to
  Terra's Magic list; text encoding derived from the tilemap
  (A-Z=$80+n, a-z=$9A+n); spell name table located
  (ROMFILE:0x26F567, 7-byte records, Fire=0, Cure=45); the $C46AC0
  table shown to be the global spell database with byte 5 = MP cost
  (Strong hypothesis); spell availability located at WRAM:+$1A6E
  (unique 128 KB search hit, $FF-marked ids {0,45} for Terra).
  Workflow: census-observer skill + /census-observations,
  /register-system, /update-coverage; stopping rules and the
  orchestrator's prioritization updated (observe broadly, register
  briefly, investigate narrowly). Coverage after the unit: 43 census
  entries across 10 of 12 domains; ROM known 10,753 bytes (0.34%).

- 2026-07-30 (resume, late) — Unit 24 / EXP-0024 + AUD-0001: **first
  audio vertical proof complete — the battle confirm SFX chain is
  Confirmed end to end.** Press-vs-no-press delta trials: exactly one
  port write in P (`$21`→`$2140` from `ROMCPU:$C117CC`, press+2), zero
  in N; voice bitmap XOR isolates **DSP voice 7** (rel 76-79);
  eval-armed DSP snapshot: SRCN=$05 → directory $1B00 → sample at
  `ARAM:$48D8` (2 BRR blocks). **SFX pack `$4800-$491F` byte-identical
  to `ROMFILE:0x051EC9-0x051FE8`.** Go: `brr.Decode` (S-DSP
  semantics, tests+fuzz) + `ff6lab brr info`. M3 complete; M7 in
  progress. Music replay determinism across trials confirmed on the
  APU side too.

- 2026-07-30 (resume, late) — Unit 23 / EXP-0023 + GFX-0001: **first
  graphics vertical proof complete — the battle HUD font is a raw ROM
  copy.** Atomic capture at anchor+120 (probe `EXP-0023.lua`; bridge
  `read` extended with vram/cgram/oam/aram/rom/dsp types): mode 1, BG3
  chr word $5000, 37 HUD tiles all in $180-$1FF, BG palette 0. ROM
  search: 15 distinctive tiles, one hit each, unanimous base
  **ROMFILE 0x046FC0**; identical run = tiles $FF-$1FF (257 tiles).
  Go: `internal/graphics/tile2bpp` (tests+fuzz) +
  `ff6lab tiles decode2bpp`; decoded VRAM == decoded ROM on sampled
  glyphs. M2 complete; M6 in progress. Load path/glyph semantics
  queued.

- 2026-07-30 (resume, late) — Unit 22 / EXP-0022: **attack-record table
  cross-checked against the live ROM.** 256 records dumped from
  `ROMCPU:$C46AC0` and decoded via the new `ff6lab attackdata scan`
  (first real-ROM use of the package). Record 238 (the enemy Fight
  record EXP-0021 saw loaded 8×): power byte 0 ✓ (matches EXP-0018's
  `v=0` MVN writes) and **physical-formula flag set** ✓ (matches
  EXP-0017's decoded path). Fire Beam (power 60, fire element,
  standard formula) narrowed to two candidates: indices 5 and 131
  (Tentative). CLI help text synchronized with the actual commands.
  MP verification recorded as savestate-blocked (BLOCKERS).

- 2026-07-30 (resume, late) — EXP-0021: **question #30 resolved for the
  tested window — the action layer is deterministic given the input
  frame schedule.** Two frame-exact trials (first press anchor+72 vs
  +270) produced byte-identical matched-ordinal content: all 8 record
  loads index 238, enemy powers 13/0/19/0 with the miss at the same
  ordinal; only cluster timing (+2/+69 frames) and two outcome-neutral
  press residues (`+$3A71` $04-vs-$00, POP3 carry) varied. GUI-era
  variance (EXP-0016/0018/0020) attributed to wall-clock harness jitter
  (Strong hypothesis). Operationally: the lab now runs **headless**
  (`Mesen --testrunner --timeout=7200` + `FF6_OUT` env + frame-scheduled
  input in `mesen/probes/EXP-0021.lua`) because the locked display
  breaks GUI Mesen (Avalonia RenderTimer −6661 / script window never
  opens); testrunner default ~2-min timeout discovered and overridden.
  Battle path pauses at this natural boundary per the operator's
  rebalance directive.

- 2026-07-30 (resume) — Manifest-completeness repair: EXP-0019 was
  missing from `manifests/experiments.json` (and therefore from the
  generated index); entry added from the record, index regenerated,
  and `ff6lab audit` gained the reverse-completeness check
  (`CheckExperimentRecordsInManifest`) that would have caught it.
  STATISTICS refreshed (20 experiments, 7 test-bearing packages).

- 2026-07-30 (late) — EXP-0020 + bridge v2 validation: **`+$3A70`
  refuted as the timing-varying state** (identical matched-ordinal
  reads under divergent timing; it is a per-action flag written just
  before use). What diverges is the number/timing of actions →
  scheduling interpretation (Tentative), RNG existence still Unknown;
  question #30 refined to EXP-0021. Bridge v2 validated live: probe
  loading, id'd atomic responses, duplicate-id suppression, transcripts;
  three defects found and fixed during validation (nil-leading candidate
  list, CWD-relative probe path, probe-to-probe dofile).

- 2026-07-30 (maintenance) — Repository repair unit: **fixed the
  `.gitignore` collision that had silently untracked
  `cmd/ff6lab/main.go` and `local_artifacts/README.md`** (the published
  repo could not build the CLI from a clean clone); CI now builds the
  CLI, checks required tracked sources, and runs `ff6lab audit`; root
  `CLAUDE.md` constitution added; `internal/audit` package +
  `ff6lab audit`/`ff6lab indexes generate` implemented with tests;
  `indexes/EXPERIMENTS.md` generated from the manifest; seven discovery
  records created (DISC-0001..0007) closing the experiment→discovery→
  implementation gap; SESSION_003 citations repointed to the immutable
  archive; PACKAGE_MANIFEST/BUILD_VERIFICATION labeled historical;
  per-experiment evidence layout adopted (EVIDENCE_LAYOUT.md); bridge
  v2 (relative paths, id'd atomic responses, probe files, fixed
  loadstate logging); `math.MinInt32` delta overflow fixed with
  regression tests; golden damage test tightened to the exact 347;
  dashboards/indexes/06_BATTLE_SYSTEM synchronized through EXP-0019.

- 2026-07-30 (late) — Unit 20 / EXP-0019: **the attack/spell data table
  discovered** — MVN loads 14-byte records from `ROMCPU:$C46AC0`
  (`ROMFILE:0x046AC0`) into the action block; every decoded `$11Ax`
  meaning becomes a record field. Fight-command staging decoded (six
  more per-slot stat tables; family ≈19). Miss/variance reinterpreted:
  divergence is enemy-AI action selection (#30). Go:
  `internal/game/attackdata` (Decode/RecordAt + typed accessors,
  synthetic fixtures, fuzz). First ROM data format implemented.

- 2026-07-30 (evening) — Unit 19 / EXP-0018: **the hit/act gate
  localized** — `$C2297D` clears the action block, `$C229D4` populates
  power (values deterministic: 13/19 per enemy); identical-state trials
  diverge only in *which* clears get populated. A miss =
  cleared-but-never-populated. The RNG consumer sits between the two,
  caller bracket `JSR` at `~$C2319D`. Question #29 → dump bracket
  (EXP-0019).

- 2026-07-30 (evening) — Unit 18 / EXP-0017: **physical base formula
  fully decoded — and the entire damage arithmetic is RNG-free.**
  Vigor²-shaped: t=power(×4 enemy)(×1.75 gated)+statA, ×statB²/256;
  party tail ×1.5+power; `+$3C58,X` halve/¾ flags (family member #13).
  Misses arrive as power=0 → RNG lives in the action-setup/hit-roll
  layer (new #29). `+$11AF` producer found (`+$3B18,X` per-slot table).
  Go: `BaseAmountPhysical` + `PhysicalFlags` with the EXP-0016 base-7
  closure as a golden test.

- 2026-07-30 (evening) — Unit 17 / EXP-0016: **variance confirmed,
  timing-dependent** — identical-savestate trials with different input
  timing produced miss-vs-hit on the first enemy action and varying
  bases (7/0/12). RNG consumer narrowed to the enemy/physical base
  path; miss-path writer `$C22C02` found. Local battle savestate
  technique established (`createSavestate` via queued exec callback).

- 2026-07-30 (overnight 2) — Units 15–16 / EXP-0014+0015: **the base
  damage formula is decoded and numerically closed** — `base = power×4 +
  (power×$11AE×$11AF)>>5` (live: 60/28/4 → 450 → ×(255−58)/256+1 →
  346, the observed Fire Beam). The `$C20DD1` 24-bit-shift helper
  explains the `$C247B7` wrapper's memory layout. Go:
  `BaseAmountStandard` + golden chain test. The standard damage
  pipeline is now byte-exact from formula to HUD. New: #26 physical
  path, #27 stat meanings, #28 variance hunt.

- 2026-07-30 (daytime) — Units 13–14 / EXP-0012+0013: arithmetic helpers
  decoded — `$C2370B` is a ×1.5-per-count chain (**randomness hypothesis
  refuted**); `$C24781` is the SNES hardware 8×8 multiply, making the
  defense scaling exactly `(amount×(255−def))/256 + 1`. Go:
  `Scale256`, `ApplyDefense`, `ChainBoost` + tables + a composition
  fuzz target. Questions #24/#25 closed; #23 (`+$11B0` producer) is the
  innermost frontier.

- 2026-07-30 (daytime) — Unit 12 / EXP-0011: **base-amount routines
  decoded** — variant A consumes a precomputed base (`+$11B0`) and
  applies defense ((255−def)/256 shape, pair at `+$3BB8,Y`), flag
  halvings, party-vs-party halving, and a final `$C2370B` transform;
  variant B = fraction-of-HP with min 1; `+$11A3` bit 7 retargets X/Y
  by `$14` (HP→MP). Family census reaches **twelve** `$14`-stride
  arrays (05 master table). No Go this unit (helpers unverified).
  Questions #23–25 opened.

- 2026-07-30 (daytime) — Units 10–11 / EXP-0009+0010: formula frame
  verified (`JSR $0B83` at `$C23469`, ten-slot target loop) and the
  **elemental-modifier block decoded byte-exact** (`$C20B83`–`$C20C2C`):
  nullify/flip-to-heal/zero/halve/double transforms over per-target
  16-bit element masks in two more `$14`-stride family arrays
  (`+$3BCC`, `+$3BE0`). Go: `battle.ApplyElementResponse` (behavior-
  derived names) + 14-case table test. Base-amount routines
  (`$C20C9E`/`$C20D87`) opened as question #22 — the innermost formula.

- 2026-07-30 (daytime) — Unit 8 / EXP-0007: **DP `$F0` = final per-hit
  damage, Confirmed on three anchors** (array arithmetic 9+4+2+6=21;
  HUD mid-values; captured popup "6"). A 346-damage Fire Beam killed the
  Were-Rat on camera. X/Y = attacker/target slot×2 (strong); `$F2`
  direction observation recorded. Question #21 narrowed: the full
  formula runs upstream of `$C20C28`; deeper-return lead `$6B16`.

- 2026-07-30 (daytime) — Unit 7 / EXP-0006: **pending-delta accumulator
  decoded byte-exact** (`ROMCPU:$C20C76`): amounts accumulate per slot
  with `$FFFF`-sentinel init and a **9999 cap** (`$270F`) — the damage
  cap located in code. One routine serves both pending arrays via a
  `Y += $14` retarget. Caller gate block at `$C20C2D` identified;
  EXP-0004's `$C20434` caller inference corrected (stack misparse).
  Go: `battle.AccumulatePending` + table tests. FN-0008 pending index.
  New question #21 (the actual damage formula computing DP `$F0`).

- 2026-07-30 (daytime) — Unit 6 / EXP-0005: **enemy slots confirmed** —
  enemy HP observed at entries 4–5 of the unified arrays (24/35), damaged
  and zeroed by the same delta-engine stores and death handler while
  VICKS won the encounter. H-BATTLE-0008 resolved Confirmed. Go model
  refactored: `battle.PartySlots` → `battle.BattleSlots` (10 entries,
  slot-uniform engine), tests extended (enemy-slot case). New writer
  `ROMCPU:$C22CCE` queued as question 19b.

- 2026-07-29 (overnight) — Unit 5 / EXP-0004: pending-delta producers
  identified at the PC level (setter `ROMCPU:$C20C9B`, sweepers
  `$C2638E`/`$C26391`); arrays shown transient (`$FFFF` between events);
  **structural discovery: the battle arrays are 10 entries wide**
  (`$14` stride family `+$3BF4/+$3C08/+$3C1C/+$3C30`; entry-9 write
  observed) — 4 party + 6 enemy candidate. Refresh trigger localized to
  the post-fetch driver (`JSR $069B` at `$C21409` → `$C2069B` → tail
  path). Questions #13/#19 narrowed; battle.go's 4-slot model flagged as
  party-view-only pending confirmation.

- 2026-07-29 (overnight) — Unit 4 / EXP-0003: open question #1 **answered**.
  PartyDisplaySourceRefresh (`ROMCPU:$C25D26`) decoded byte-exact: copies
  all six authoritative battle arrays into the `+$2E78` display family.
  H-BATTLE-0002 and H-BATTLE-0006 resolved (Confirmed); H-BATTLE-0004
  upgraded; `+$2E98` identified as the status-word copy (`+$3EE4`) and
  `+$2EA0` as the `+$3EF8` copy. New questions 19–20.

- 2026-07-29 (overnight) — Unit 3 / EXP-0002: PerFrameBattleUpdate fires
  only in battle across all tested contexts (≈175k non-battle frames, zero
  fires; two battles positive). Corrections: "every frame" softened to
  phase-dependent ≈0.23–0.87/frame; second one-shot caller `JSR
  ROMCPU:$C11090` discovered at battle entry. Hazard: Mesen auto-save
  destroyed the `_11.mss` battle state (recorded in capability matrix).
  Battle-init writers reproduced in an independent encounter.

- 2026-07-29 (overnight) — Unit 2 / EXP-0001: launched Mesen 2.1.1 with the
  bridge, dumped `ROMCPU:$C212F0–$C2141F`, verified every battle.go
  disassembly claim byte-exact; FN-0003..0006 upgraded to Confirmed
  (code); new unknowns recorded (`+$11A2` selector, dispatch tail, fetch
  gates); ROM mapping upgraded to Confirmed HiROM via Mesen header parse;
  Session 003 evidence archived under `mesen/out/session003/`.

- 2026-07-29 — Sessions 001–002: HP variable discovery, full decode of
  CopyCharacterFields (`ROMCPU:$C10DF3`), caller identification
  (`ROMCPU:$C101FB`), `chardata` Go package with tests. Records in
  `docs/sessions/`.
- 2026-07-29 — Session 003 (undocumented): bank-`$C2` HP/MP delta engine
  investigation; raw evidence in `mesen/out/`, implementation in
  `battle/battle.go`; interrupted by the V4 install before documentation.
- 2026-07-29 — Version 4 scaffold created.
- 2026-07-29 — `/audit-project`: 23 broken links fixed (V4 restructure),
  `mesen/bridge.lua` output path repaired, `*.mss` added to `.gitignore` and
  CI reject-list, stale doc paths in `chardata` fixed. Substantive findings
  in the audit report (chat) and migration report.
- 2026-07-29 — `/bootstrap-v4`: environment and ROM identity recorded,
  Mesen capability matrix filled from Session 002 evidence, indexes and
  dashboards initialized from canonical records, quality gates run,
  migration report and checkpoint written.
- 2026-07-29 (overnight) — Unit 1: SESSION_003 reconstructed record written
  from raw evidence; battle delta engine promoted to 02/04/05/06/08 with
  honest confidence split (stores Confirmed, addresses Strong hypothesis,
  disassembly Unknown); `internal/game/battle` gained table-driven tests;
  battle.go provenance comment corrected.
- 2026-07-29 — Post-migration decisions executed: module path restored to
  `github.com/bensabler/ff6-decompile`; `battle`/`chardata` moved under
  `internal/game/`; notation conflict resolved (CONTRA-0001); repository
  initialized in git and published publicly to GitHub.
